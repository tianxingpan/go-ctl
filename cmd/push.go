// Package cmd provides add custom CTL command
package cmd

/*
 * 上传文件或目录到远程机器
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

var push = &cobra.Command{
	Use:                   "push",
	Short:                 "Push file or directory",
	Long:                  "Push file to remote machines",
	Example:               "go-ctl push -s=./test.txt -d=/data/gongyi",
	DisableFlagsInUseLine: true,
	Run:                   execPush,
}

// 执行上传文件
func execPush(cmd *cobra.Command, args []string) {
	// 还得判断命令行参数是否有值，没有的话，直接不绑定关系
	_ = viper.BindPFlag("ras.hosts", cmd.Flags().Lookup("hosts"))
	_ = viper.BindPFlag("ras.port", cmd.Flags().Lookup("port"))
	_ = viper.BindPFlag("ras.user", cmd.Flags().Lookup("user"))
	_ = viper.BindPFlag("ras.password", cmd.Flags().Lookup("password"))
	_ = viper.BindPFlag("ras.source", cmd.Flags().Lookup("source"))
	_ = viper.BindPFlag("ras.destination", cmd.Flags().Lookup("destination"))

	// 绑定环境变量
	_ = viper.BindEnv("ras.hosts", "RAS_HOSTS")
	_ = viper.BindEnv("ras.user", "RAS_USER")
	_ = viper.BindEnv("ras.password", "RAS_PASSWORD")

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
	source := viper.Get("ras.source").(string)
	destination := viper.Get("ras.destination").(string)

	hostArray := strings.Split(hosts, ",")
	var wg sync.WaitGroup
	wg.Add(len(hostArray))
	for _, host := range hostArray {
		addr := fmt.Sprintf("%s:%d", host, port)
		go runGoroutine(&wg, addr, user, password, source, destination)
		//go func() {
		//	defer wg.Done()
		//	fmt.Println(user, host)
		//
		//	ipHost := fmt.Sprintf("%s:%d", host, port)
		//
		//	sftClient := controller.SFTP{}
		//	err := sftClient.Init(ipHost, user, password)
		//	if err != nil {
		//		fmt.Println(err)
		//		return
		//	}
		//	defer sftClient.Close()
		//
		//	err = sftClient.Push(source, destination)
		//	if err != nil {
		//		fmt.Println(err)
		//	}
		//	//fmt.Printf("[\033[1;33mOK\033[m][%s]\n\n", ipHost)
		//}()
	}
	wg.Wait()
}

func runGoroutine(wg *sync.WaitGroup, addr, user, password, source, destination string) {
	defer wg.Done()
	fmt.Println(addr, user)

	sftClient := controller.SFTP{}
	err := sftClient.Init(addr, user, password)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer sftClient.Close()

	err = sftClient.Push(source, destination)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Printf("[\033[1;33mOK\033[m][%s]\n\n", ipHost)
}

func init() {
	// 设置参数
	push.Flags().StringP("hosts", "H", "", "Connect to the remote machines on the given hosts separated by comma, can be replaced by environment variable 'H'")
	push.Flags().IntP("port", "P", 22, "Specifies the port to connect to on the remote machines, can be replaced by environment variable 'PORT'")
	push.Flags().StringP("user", "u", "", "Specifies the user to log in as on the remote machines, can be replaced by environment variable 'U'")
	push.Flags().StringP("password", "p", "", "The password to use when connecting to the remote machines, can be replaced by environment variable 'P'")

	push.Flags().StringP("source", "s", "", "Local source files uploaded, separated by comma")
	push.Flags().StringP("destination", "d", "", "Remote destination directory")

	// 必传字段设置
	_ = push.MarkFlagRequired("source")
	_ = push.MarkFlagRequired("destination")
	RootCmd.AddCommand(push)
}