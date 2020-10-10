package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/haozibi/leetcode-badge/app"
)

var rootCMD = &cobra.Command{
	Use:   app.BuildAppName,
	Short: "LeetCode Badge CLI",
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
