// Package cmd provides add custom CTL command
package cmd

/*
 * 远程执行命令
 * RAS--远程（遥控）服务器
 * RAS 远程（遥控）服务器 Remote Access Server
 * 环境变量设置时，前缀为RAS_
*/

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tianxingpan/go-ctl/pkg/controller"
	"os"
	"strconv"
	"strings"
	"sync"
)

var cli = &cobra.Command{
	Use:                   "cli",
	Short:                 "The command is executed on the remote machines.",
	Long:                  "The remote connection machine executes the command and returns the corresponding resultThe remote connection machine executes the command and returns the corresponding result",
	Example:               "go-ctl cli -c='ls /bin'",
	DisableFlagsInUseLine: true,
	Run:                   execCLI,
}

func checkCliParams(cmd *cobra.Command) bool {
	return false
}

// SSH命令执行
func execCLI(cmd *cobra.Command, args []string) {
	// 参数校验
	var port int
	envPort := os.Getenv("RAS_PORT")
	if envPort != "" {
		v, err := cmd.Flags().GetInt("port")
		if err != nil || v == 0 || v == 22 {
			port, _ = strconv.Atoi(envPort)
			viper.Set("ras.port", port)
		}
	} else {
		port = viper.Get("ras.port").(int)
		if port == 0 {
			fmt.Println("Error: required flag(s) \"port\" not set")
			_ = cmd.Help()
			return
		}
	}
	hosts := viper.Get("ras.hosts").(string)
	user := viper.Get("ras.user").(string)
	password := viper.Get("ras.password").(string)
	command := viper.Get("ras.command").(string)

	// 具体执行
	hostArray := strings.Split(hosts, ",")
	var wg sync.WaitGroup
	wg.Add(len(hostArray))
	for _, host := range hostArray {
		go func() {
			defer wg.Done()

			ipHost := fmt.Sprintf("%s:%d", host, port)
			sshClient := controller.SSH{}
			err := sshClient.Init(ipHost, user, password)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer sshClient.Close()
			err = sshClient.Execute(command)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("[\033[1;33mOK\033[m][%s]\n\n", ipHost)
		}()
	}
	wg.Wait()
}

func init() {
	// 设置参数
	cli.Flags().StringP("hosts", "H", "", "Connect to the remote machines on the given hosts separated by comma, can be replaced by environment variable 'H'")
	cli.Flags().IntP("port", "P", 22, "Specifies the port to connect to on the remote machines, can be replaced by environment variable 'PORT'")
	cli.Flags().StringP("user", "u", "", "Specifies the user to log in as on the remote machines, can be replaced by environment variable 'U'")
	cli.Flags().StringP("password", "p", "", "The password to use when connecting to the remote machines, can be replaced by environment variable 'P'")
	cli.Flags().StringP("command", "c", "", "The command is executed on the remote machines")

	// 还得判断命令行参数是否有值，没有的话，直接不绑定关系
	_ = viper.BindPFlag("ras.hosts", cli.Flags().Lookup("hosts"))
	_ = viper.BindPFlag("ras.port", cli.Flags().Lookup("port"))
	_ = viper.BindPFlag("ras.user", cli.Flags().Lookup("user"))
	_ = viper.BindPFlag("ras.password", cli.Flags().Lookup("password"))
	_ = viper.BindPFlag("ras.command", cli.Flags().Lookup("command"))

	// 绑定环境变量
	_ = viper.BindEnv("ras.hosts", "RAS_HOSTS")
	_ = viper.BindEnv("ras.user", "RAS_USER")
	_ = viper.BindEnv("ras.password", "RAS_PASSWORD")

	// 必传字段设置
	_ = cli.MarkFlagRequired("command")
	RootCmd.AddCommand(cli)
}
