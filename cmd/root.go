// Package cmd provides add custom CTL command
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tianxingpan/go-ctl/pkg/config"
	"runtime"
)

var (
	buildVersion = "0.0.1"
	RootCmd = &cobra.Command{
		Use:   "go-ctl",
		Short: "a SSH tool.",
		Long:  `go-ctl is an SSH assistant tool.`,
	}
)

// 执行root命令
func Execute() error {
	var status bool
	err := RootCmd.Execute()
	if err == nil {
		status = true
	} else {
		status = false
	}

	config.RefreshConfig(status)

	return err
}

func init() {
	// 定义版本号
	RootCmd.Version = fmt.Sprintf("%s %s/%s", buildVersion, runtime.GOOS, runtime.GOARCH)
	// 注册初始化具柄
	cobra.OnInitialize(initConfig)
}

// 初始化配置文件
func initConfig() {
	if err := config.Init(""); err != nil {
		panic(err)
	}
}
