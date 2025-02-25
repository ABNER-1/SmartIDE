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

package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/leansoftX/smartide-cli/pkg/common"
)

// 是否链接了docker compose文件
func (c SmartIdeConfig) IsLinkDockerComposeFile() bool {
	dockerFilePath := strings.TrimSpace(c.Workspace.DockerComposeFile)
	return len(dockerFilePath) > 0
}

// 获取本地链接的docker compose文件的路径
func (instance SmartIdeConfig) GetLocalLinkDockerComposeFile() (
	localLinkDockerComposeFilePath string, localLinkDockerComposeFileContent string) {
	return instance.getLinkDockerComposeFile(nil)

}

// 获取本地链接的docker compose文件的路径
func (instance SmartIdeConfig) GetRemoteLinkDockerComposeFile(sshRemote *common.SSHRemote) (
	localLinkDockerComposeFilePath string, localLinkDockerComposeFileContent string) {
	return instance.getLinkDockerComposeFile(sshRemote)

}

// 获取本地链接的docker compose文件的路径
func (instance SmartIdeConfig) getLinkDockerComposeFile(sshRemote *common.SSHRemote) (
	localLinkDockerComposeFilePath string, localLinkDockerComposeFileContent string) {

	// 如果没有配置链接docker-compose直接退出
	if !instance.IsLinkDockerComposeFile() {
		return localLinkDockerComposeFilePath, localLinkDockerComposeFileContent
	}

	// 确定工作目录
	workingDir := ""
	if instance.Workspace.DevContainer.workingDirectoryPath != "" { // 有指定的工作目录时
		workingDir = instance.Workspace.DevContainer.workingDirectoryPath
	} else {
		workingDir, _ = os.Getwd()
	}

	// 一定要获取，是因为docker-compose文件到路径是相对配置文件所在目录的
	configYamlFileDir := filepath.Dir(instance.Workspace.DevContainer.configRelativeFilePath)

	// 可能会有yaml配置文件是字符串传递的问题，这种情况下就直接用当前工作目录作为基准
	localLinkDockerComposeFilePath = common.FilePahtJoin4Linux(workingDir, configYamlFileDir, instance.Workspace.DockerComposeFile)

	// 获取文件内容
	if localLinkDockerComposeFilePath != "" {
		if sshRemote == nil || (sshRemote == &common.SSHRemote{}) { // 本地模式
			// read and parse
			localLinkDockerComposeFileContentBytes, err := os.ReadFile(localLinkDockerComposeFilePath)
			common.CheckError(err)
			localLinkDockerComposeFileContent = string(localLinkDockerComposeFileContentBytes)
		} else { // 远程主机模式
			// read and parse
			localLinkDockerComposeFileContent = sshRemote.GetContent(localLinkDockerComposeFilePath)
		}

	}

	if localLinkDockerComposeFileContent == "" {
		common.SmartIDELog.Error("link compose file is empty")
	}

	//

	return localLinkDockerComposeFilePath, localLinkDockerComposeFileContent
}

// 获取服务名称列表
func (c *SmartIdeConfig) GetServiceNames() (serviceNames []string) {

	if c.Workspace.LinkCompose != nil {
		for serviceName := range c.Workspace.LinkCompose.Services {
			serviceNames = append(serviceNames, serviceName)
		}
	} else {
		for serviceName := range c.Workspace.Servcies {
			serviceNames = append(serviceNames, serviceName)
		}
	}

	return serviceNames
}

// 获取yaml配置文件的路径
func (yamlFileConfig *SmartIdeConfig) GetConfigFileAbsolutePath() string {
	return yamlFileConfig.getConfigYamlFilePath()
}

// 获取本地配置文件所在的路径
func (c *SmartIdeConfig) getConfigYamlFilePath() string {
	return filepath.Join(c.Workspace.DevContainer.workingDirectoryPath, c.Workspace.DevContainer.configRelativeFilePath)
}
