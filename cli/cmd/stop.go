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

package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	cmdCommon "github.com/leansoftX/smartide-cli/cmd/common"

	"github.com/leansoftX/smartide-cli/cmd/server"
	"github.com/leansoftX/smartide-cli/internal/apk/appinsight"
	"github.com/leansoftX/smartide-cli/internal/biz/workspace"
	"github.com/leansoftX/smartide-cli/internal/model/response"
	"github.com/leansoftX/smartide-cli/pkg/common"
	"github.com/spf13/cobra"
) // stopCmd represents the stop command

var stopCmd = &cobra.Command{
	Use:     "stop",
	Short:   i18nInstance.Stop.Info_help_short,
	Long:    i18nInstance.Stop.Info_help_long,
	Example: `  smartide stop {workspaceid} `,
	Run: func(cmd *cobra.Command, args []string) {

		mode, _ := cmd.Flags().GetString("mode")
		workspaceIdStr := cmdCommon.GetWorkspaceIdFromFlagsOrArgs(cmd, args)
		if strings.ToLower(mode) == "server" || strings.Contains(workspaceIdStr, "SWS") {
			serverModeInfo, _ := server.GetServerModeInfo(cmd)
			if serverModeInfo.ServerHost != "" {
				wsURL := fmt.Sprint(strings.ReplaceAll(strings.ReplaceAll(serverModeInfo.ServerHost, "https", "ws"), "http", "ws"), "/ws/smartide/ws")
				common.WebsocketStart(wsURL)

				if pid, err := workspace.GetParentId(workspaceIdStr, workspace.ActionEnum_Workspace_Stop, serverModeInfo.ServerToken, serverModeInfo.ServerHost); err == nil && pid > 0 {
					common.SmartIDELog.Ws_id = workspaceIdStr
					common.SmartIDELog.ParentId = pid
				}

			}
		}
		common.SmartIDELog.Info(i18nInstance.Stop.Info_start)

		// 检查错误并feedback
		var checkErrorFeedback = func(err error, workspaceInfo workspace.WorkspaceInfo) {
			if err != nil {
				server.Feedback_Finish(server.FeedbackCommandEnum_Stop, cmd, false, nil, workspaceInfo, err.Error(), "")
			}
			common.CheckError(err)
		}

		// 获取 workspace 信息
		common.SmartIDELog.Info(i18nInstance.Main.Info_workspace_loading)
		workspaceInfo, err := cmdCommon.GetWorkspaceFromCmd(cmd, args)
		entryptionKey4Workspace(workspaceInfo) // 申明需要加密的文本
		common.CheckError(err)

		if workspaceInfo.CliRunningEnv == workspace.CliRunningEvnEnum_Server { // cli 在服务器上运行
			// 远程主机上停止
			appinsight.SetCliLocalTrack(appinsight.Cli_Host_Stop, args, workspaceInfo.ID, "")
			err := stopRemote(workspaceInfo)
			checkErrorFeedback(err, workspaceInfo)

			// feeadback
			common.SmartIDELog.Info("反馈运行结果...")
			err = server.Feedback_Finish(server.FeedbackCommandEnum_Stop, cmd, err == nil, nil, workspaceInfo, "", "")
			common.CheckError(err)

		} else if workspaceInfo.Mode == workspace.WorkingMode_Remote &&
			workspaceInfo.CacheEnv == workspace.CacheEnvEnum_Server &&
			workspaceInfo.CliRunningEnv == workspace.CliRunningEnvEnum_Client { // 录入的是服务端工作区id

			// 当前用户信息
			currentServerAuth, err0 := workspace.GetCurrentUser() // 当前服务区用户，在mode server 模式下才会赋值
			common.CheckError(err0)

			// 触发stop
			err := server.Trigger_Action("stop", workspaceIdStr, currentServerAuth, make(map[string]interface{}))
			common.CheckError(err)

			// 轮询检查工作区状态
			isStop := false
			for !isStop {
				serverWorkSpace, err := workspace.GetWorkspaceFromServer(currentServerAuth, workspaceInfo.ID, workspace.CliRunningEnvEnum_Client)
				if serverWorkSpace == nil {
					common.SmartIDELog.Error("工作区数据查询为空！")
				}
				if err != nil {
					common.SmartIDELog.ImportanceWithError(err)
				}
				if serverWorkSpace.ServerWorkSpace.Status == response.WorkspaceStatusEnum_Stop ||
					serverWorkSpace.ServerWorkSpace.Status == response.WorkspaceStatusEnum_Error_Stop {
					isStop = true
				}

				time.Sleep(time.Second * 15)
			}

		} else { // 普通模式下
			// 判断是否有工作区数据
			if workspaceInfo.IsNil() {
				common.SmartIDELog.Error(i18nInstance.Main.Err_workspace_none)
			}

			// 执行对应的stop
			if workspaceInfo.Mode == workspace.WorkingMode_Local {
				appinsight.SetCliLocalTrack(appinsight.Cli_Local_Stop, args, workspaceInfo.ID, "")
				stopLocal(workspaceInfo)

			} else {
				appinsight.SetCliLocalTrack(appinsight.Cli_Host_Stop, args, workspaceInfo.ID, "")
				err := stopRemote(workspaceInfo)
				common.CheckError(err)

			}

		}

		common.SmartIDELog.Info(i18nInstance.Stop.Info_end)
		common.WG.Wait()
	},
}

// 停止本地容器
func stopLocal(workspace workspace.WorkspaceInfo) {
	// 校验是否能正常执行docker
	err := common.CheckLocalEnv()
	common.CheckError(err)

	// 本地执行docker-compose
	composeCmd := exec.Command("docker-compose", "-f", workspace.TempYamlFileAbsolutePath,
		"--project-directory", workspace.WorkingDirectoryPath, "stop")
	composeCmd.Stdout = os.Stdout
	composeCmd.Stderr = os.Stderr
	if composeCmdErr := composeCmd.Run(); composeCmdErr != nil {
		common.SmartIDELog.Fatal(composeCmdErr)
	}
}

// 停止远程容器
func stopRemote(workspaceInfo workspace.WorkspaceInfo) error {
	// ssh 连接
	common.SmartIDELog.Info(i18nInstance.Stop.Info_sshremote_connection_creating)

	sshRemote, err := common.NewSSHRemote(workspaceInfo.Remote.Addr, workspaceInfo.Remote.SSHPort, workspaceInfo.Remote.UserName, workspaceInfo.Remote.Password, workspaceInfo.Remote.SSHKey)
	if err != nil {
		return err
	}

	// 项目文件夹是否存在
	if !sshRemote.IsDirExist(workspaceInfo.WorkingDirectoryPath) {
		msg := fmt.Sprintf(i18nInstance.Stop.Err_env_project_dir_remove, workspaceInfo.ID)
		return errors.New(msg)
	}

	// 检查临时文件夹是否存在
	if !sshRemote.IsFileExist(workspaceInfo.TempYamlFileAbsolutePath) || !sshRemote.IsFileExist(workspaceInfo.ConfigYaml.GetConfigFileAbsolutePath()) {
		workspaceInfo.SaveTempFilesForRemote(sshRemote)
	}

	// 检查环境
	err = sshRemote.CheckRemoteEnv()
	if err != nil {
		return err
	}

	// 停止容器
	common.SmartIDELog.Info(i18nInstance.Stop.Info_docker_stopping)
	command := fmt.Sprintf("docker-compose -f %v --project-directory %v stop",
		common.FilePahtJoin4Linux(workspaceInfo.TempYamlFileAbsolutePath), common.FilePahtJoin4Linux(workspaceInfo.WorkingDirectoryPath))
	err = sshRemote.ExecSSHCommandRealTime(command)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	//stopCmd.Flags().StringVarP(&configYamlFileRelativePath, "filepath", "f", "", i18nInstance.Stop.Info_help_flag_filepath)

}
