/*
SmartIDE - Dev Containers
Copyright (C) 2023 leansoftX.com

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package start

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	smartideServer "github.com/leansoftX/smartide-cli/cmd/server"
	"github.com/leansoftX/smartide-cli/internal/apk/appinsight"
	"github.com/leansoftX/smartide-cli/internal/biz/config"
	"github.com/leansoftX/smartide-cli/internal/biz/workspace"
	"github.com/leansoftX/smartide-cli/internal/model"
	"github.com/leansoftX/smartide-cli/internal/model/response"
	"github.com/leansoftX/smartide-cli/pkg/common"
	"github.com/leansoftX/smartide-cli/pkg/tunnel"
)

// 远程服务器执行 start 命令
func ExecuteServerVmStartByClientEnvCmd(workspaceInfo workspace.WorkspaceInfo,
	yamlExecuteFun func(yamlConfig config.SmartIdeConfig, workspaceInfo workspace.WorkspaceInfo, cmdtype, userguid, workspaceid string)) (
	workspace.WorkspaceInfo, error) {
	currentAuth, err := workspace.GetCurrentUser()
	if err != nil {
		return workspaceInfo, err
	}

	if currentAuth != (model.Auth{}) && currentAuth.Token != "" && currentAuth.Token != nil {
		wsURL := fmt.Sprint(strings.ReplaceAll(strings.ReplaceAll(currentAuth.LoginUrl, "https", "ws"), "http", "ws"), "/ws/smartide/ws")
		common.WebsocketStart(wsURL)
		if pid, err := workspace.CreateWsLog(workspaceInfo.ServerWorkSpace.NO, currentAuth.Token.(string), currentAuth.LoginUrl, "客户端启动工作区", "客户端启动工作区", common.SmartIDELog.TekEventId); err == nil {
			common.SmartIDELog.Ws_id = workspaceInfo.ServerWorkSpace.NO
			common.SmartIDELog.ParentId = pid
		}
	}
	yamlExecuteFun(workspaceInfo.ConfigYaml, workspaceInfo, appinsight.Cli_Host_Start, "", workspaceInfo.ID)
	//
	common.SmartIDELog.Info(i18nInstance.VmStart.Info_starting)
	// 检查工作区的状态
	if workspaceInfo.ServerWorkSpace.Status != response.WorkspaceStatusEnum_Start {
		if workspaceInfo.ServerWorkSpace.Status == response.WorkspaceStatusEnum_Pending ||
			workspaceInfo.ServerWorkSpace.Status == response.WorkspaceStatusEnum_Init {
			return workspaceInfo, errors.New("当前工作区正在启动中，请等待！")

		} else if workspaceInfo.ServerWorkSpace.Status == response.WorkspaceStatusEnum_Stop {
			return workspaceInfo, errors.New("当前工作区已停止！")

		} else {
			return workspaceInfo, errors.New("当前工作区运行异常！")

		}
	}

	//0. 连接到远程主机
	msg := fmt.Sprintf(" %v@%v:%v ...", workspaceInfo.Remote.UserName, workspaceInfo.Remote.Addr, workspaceInfo.Remote.SSHPort)
	common.SmartIDELog.Info(i18nInstance.VmStart.Info_connect_remote + msg)
	if workspaceInfo.Remote.IsNil() {
		return workspaceInfo, errors.New("关联 远程主机 信息为空！")
	}

	sshRemote, err := common.NewSSHRemote(workspaceInfo.Remote.Addr, workspaceInfo.Remote.SSHPort, workspaceInfo.Remote.UserName, workspaceInfo.Remote.Password, workspaceInfo.Remote.SSHKey)
	if err != nil {
		return workspaceInfo, err
	}

	//6. 当前主机绑定到远程端口
	common.SmartIDELog.Info(i18nInstance.VmStart.Info_tunnel_waiting) // log
	var addrMapping map[string]string = map[string]string{}
	var unusedLocalPort4IdeBindingPort int
	for i, pmi := range workspaceInfo.Extend.Ports {
		port := &workspaceInfo.Extend.Ports[i] // 指针，用于重新设定值
		if pmi.CurrentHostPort <= 0 {
			common.SmartIDELog.Importance(fmt.Sprintf("%v 绑定端口不正确 ", pmi.CurrentHostPort))
			continue
		}

		// 获取本地未使用的端口
		var unusedClientPort int
		unusedClientPort, err = common.CheckAndGetAvailableLocalPort(pmi.CurrentHostPort, 100)
		if err != nil {
			common.SmartIDELog.ImportanceWithError(err)

		}

		// 更新extend.ports的信息
		if port.OldClientPort == 0 { // 值为空的时候，直接赋值
			port.OldClientPort = unusedClientPort
		} else { // 值不为空的时候，使用之前的clientport字段赋值
			port.OldClientPort = port.ClientPort
		}
		port.ClientPort = unusedClientPort // 当前端口赋值
		unusedClientPortStr := strconv.Itoa(port.ClientPort)
		addrMapping["localhost:"+unusedClientPortStr] = "localhost:" + strconv.Itoa(pmi.CurrentHostPort)

		// 获取webide的本地端口
		if pmi.HostPortDesc != "" {
			unusedClientPortStr += fmt.Sprintf("(%v)", pmi.HostPortDesc)
			if strings.Contains(strings.ToLower(pmi.HostPortDesc), "webide") {
				unusedLocalPort4IdeBindingPort = workspaceInfo.Extend.Ports[i].ClientPort
			}
		}

		// 打印信息
		msg := fmt.Sprintf("localhost:%v -> %v:%v -> container:%v",
			unusedClientPortStr, workspaceInfo.Remote.Addr, pmi.CurrentHostPort, pmi.ContainerPort)
		common.SmartIDELog.Info(msg)

		port.IsConnected = true

	}
	workspaceInfo.UpdateSSHConfig()
	//6.2. 执行绑定
	tunnel.TunnelMultiple(sshRemote.Connection, addrMapping)

	//8. 打开浏览器
	var checkUrl string
	if workspaceInfo.ConfigYaml.Workspace.DevContainer.IdeType != config.IdeTypeEnum_SDKOnly {
		//vscode启动时候默认打开文件夹处理
		common.SmartIDELog.Info(i18nInstance.VmStart.Info_warting_for_webide + fmt.Sprintf(`: %v`, unusedLocalPort4IdeBindingPort))
		switch workspaceInfo.ConfigYaml.Workspace.DevContainer.IdeType {
		case config.IdeTypeEnum_VsCode:
			checkUrl = fmt.Sprintf("http://localhost:%v/?folder=vscode-remote://localhost:%v%v",
				unusedLocalPort4IdeBindingPort, unusedLocalPort4IdeBindingPort, workspaceInfo.GetContainerWorkingPathWithVolumes())
		case config.IdeTypeEnum_JbProjector:
			checkUrl = fmt.Sprintf(`http://localhost:%v`, unusedLocalPort4IdeBindingPort)
		case config.IdeTypeEnum_Opensumi:
			checkUrl = fmt.Sprintf(`http://localhost:%v/?workspaceDir=/home/project`, unusedLocalPort4IdeBindingPort)
		default:
			checkUrl = fmt.Sprintf(`http://localhost:%v`, unusedLocalPort4IdeBindingPort)
		}
	}

	go func() {
		if workspaceInfo.ConfigYaml.Workspace.DevContainer.IdeType != config.IdeTypeEnum_SDKOnly {
			isUrlReady := false
			// 检测浏览器
			for !isUrlReady {
				resp, err := http.Get(checkUrl)
				if (err == nil) && (resp.StatusCode == 200) {
					isUrlReady = true
					//common.OpenBrowser(checkUrl) // 这里不用打开，从server中点击即可
					common.SmartIDELog.InfoF(i18nInstance.VmStart.Info_open_brower, checkUrl)

				} else {
					msg := fmt.Sprintf("%v 检测失败", checkUrl)
					common.SmartIDELog.Debug(msg)

				}
			}
		}

		//9. 更新server端的extend字段
		err = smartideServer.FeeadbackExtend(currentAuth, workspaceInfo)
		if err != nil {
			common.SmartIDELog.ImportanceWithError(err)
		}
		common.SmartIDELog.Info("本地端口绑定信息 更新完成！")
	}()

	return workspaceInfo, nil
}
