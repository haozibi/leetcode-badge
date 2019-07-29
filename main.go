package main

import (
	"os"

	"github.com/haozibi/leetcode-badge/cmd"

	"github.com/haozibi/zlog"
)

func main() {

	zlog.NewBasicLog(os.Stdout, zlog.WithNoColor(true))

	cmd.Execute()
}
