package cmd

import (
	"fmt"

	"github.com/haozibi/leetcode-badge/app"

	"github.com/spf13/cobra"
)

var versionCMD = &cobra.Command{
	Use:   "version",
	Short: "show version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%ss %s %s",
			app.BuildAppName,
			app.BuildVersion,
			app.BuildTime,
		)
	},
}
