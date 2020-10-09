package cmd

import (
	"fmt"

	"github.com/haozibi/leetcode-badge/app"
	"github.com/spf13/cobra"
)

// NewVersionCommand new version command
func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(app.Version())
		},
	}
	return cmd
}
