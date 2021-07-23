// Package cmd provides add custom CTL command
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tianxingpan/go-ctl/pkg/utils"
	"github.com/tianxingpan/go-ctl/pkg/vars"
)

var upgradeCmd = &cobra.Command{
	Use:   "update",
	Short: "upgrade go-ctl",
	Long:  "upgrade go-ctl to latest version",
	Run: func(cmd *cobra.Command, args []string) {
		upCmd := fmt.Sprintf("GO111MODULE=on go get -u %s", vars.ProjectOpenSourceUrl)
		info, err := utils.RunCMD(upCmd, "")
		if err != nil {
			fmt.Printf("Upgrade failed, errmsg:%s\n", err.Error())
			return
		}
		fmt.Print(info)
	},
}

func init() {
	RootCmd.AddCommand(upgradeCmd)
}
