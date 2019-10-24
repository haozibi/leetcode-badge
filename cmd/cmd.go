package cmd

import (
	"fmt"
	"os"

	"github.com/haozibi/leetcode-badge/app"

	"github.com/spf13/cobra"
)

var rootCMD = &cobra.Command{
	Use:   app.BuildAppName,
	Short: "leetcode badge CLI",
}

func init() {
	rootCMD.AddCommand(NewRunCommand())
	rootCMD.AddCommand(NewVersionCommand())
}

// Execute exec
func Execute() {
	if err := rootCMD.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
