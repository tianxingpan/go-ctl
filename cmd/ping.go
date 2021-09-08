// Package cmd provides add custom CTL command
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tianxingpan/go-ctl/pkg/controller"
)

var ping = &cobra.Command{
	Use:                   "ping",
	Short:                 "Check network",
	Long:                  "Check the network delay between the local machine and the target machine",
	Example:               "go-ctl ping $IP",
	DisableFlagsInUseLine: true,
	Run:                   execPing,
}

// 执行ping
func execPing(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		// 未传入IP
		fmt.Println("No incoming IP or web address")
		_ = cmd.Help()
		return
	}
	count, _ := cmd.Flags().GetInt("count")
	port, _ := cmd.Flags().GetInt("port")
	timeout, _ := cmd.Flags().GetInt("timeout")
	packet, _ := cmd.Flags().GetInt("size")
	controller.Ping(args[0], packet, count, timeout, port)
}

func init() {
	// 设置参数
	ping.Flags().IntP("port", "P", 0, "Specifies the port to connect to on the remote machines")
	ping.Flags().IntP("timeout", "t", 1000, "Request timeout")
	ping.Flags().IntP("count", "c", 0, "Number of consecutive counts")
	ping.Flags().IntP("size", "s", 32, "Packet size")

	RootCmd.AddCommand(ping)
}