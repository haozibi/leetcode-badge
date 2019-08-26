package cmd

import (
	"fmt"
	"os"

	"github.com/haozibi/leetcode-badge/app"
	"github.com/haozibi/leetcode-badge/internal/cache"

	"github.com/spf13/cobra"
)

var (
	addr        string
	cacheType   string
	cacheAddr   string
	cachePasswd string
	debug       bool
)

func init() {

	runCMD.Flags().StringVarP(&addr, "address", "a", ":8080", "http run address")
	runCMD.Flags().StringVarP(&cacheType, "type", "t", "memory", "memory or redis")
	runCMD.Flags().StringVarP(&cacheAddr, "cacheaddr", "", "", "required when the type is redis")
	runCMD.Flags().StringVarP(&cachePasswd, "passwd", "", "", "optional when the type is redis")
	runCMD.Flags().BoolVarP(&debug, "debug", "d", false, "debug")
}

var runCMD = &cobra.Command{
	Use:   "run",
	Short: "Run Web",
	Run: func(cmd *cobra.Command, args []string) {

		if os.Getenv("LCHTTPAddr") != "" {
			addr = os.Getenv("LCHTTPAddr")
		}

		var ct cache.CacheType

		switch cacheType {
		case cache.CacheRedis.String():
			if cacheAddr == "" {
				fmt.Println("miss redis address")
				os.Exit(1)
			}
			ct = cache.CacheRedis
		case cache.CacheMemory.String():
			ct = cache.CacheMemory
		default:
			fmt.Println("cache type error not support")
			os.Exit(1)
		}

		c := &app.Config{
			ListenAddr:  addr,
			CacheType:   ct,
			CacheAddr:   cacheAddr,
			CachePasswd: cachePasswd,
			Debug:       debug,
		}

		a := app.New(c)
		err := a.Run()
		if err != nil {
			fmt.Println(err)
			if debug {
				fmt.Printf("%+v", err)
			}
			os.Exit(1)
		}
	},
}
