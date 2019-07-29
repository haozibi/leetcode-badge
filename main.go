package main

import (
	"os"

	"github.com/haozibi/leetcode-badge/cmd"
	"github.com/haozibi/leetcode-badge/static"

	"github.com/haozibi/zlog"
)

func main() {

	err := static.RestoreAssets("./", "static")
	if err != nil {
		panic(err)
	}

	zlog.NewBasicLog(os.Stderr, zlog.WithNoColor(true))
	cmd.Execute()
}
