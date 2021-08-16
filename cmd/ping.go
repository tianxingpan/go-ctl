// Package cmd provides add custom CTL command
package cmd

import "github.com/spf13/cobra"

var ping = &cobra.Command{
	Use:                   "ping",
	Short:                 "Network detection command",
	Long:                  "Network detection command",
	Example:               "go-ctl ping $IP",
	DisableFlagsInUseLine: true,
	Run:                   execPing,
}

// 执行ping
func execPing(cmd *cobra.Command, args []string) {}

func init() {
	// 设置参数
	cli.Flags().IntP("port", "P", 22, "Specifies the port to connect to on the remote machines, can be replaced by environment variable 'PORT'")
	cli.Flags().IntP("timeout", "t", 1000, "Specifies the user to log in as on the remote machines, can be replaced by environment variable 'U'")
	cli.Flags().IntP("count", "c", 0, "The password to use when connecting to the remote machines, can be replaced by environment variable 'P'")

	RootCmd.AddCommand(ping)
}