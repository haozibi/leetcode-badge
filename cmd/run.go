package cmd

import (
	"fmt"
	"os"

	"github.com/haozibi/leetcode-badge/app"

	"github.com/spf13/cobra"
)

// NewRunCommand Run Web Command
func NewRunCommand() *cobra.Command {

	var opt app.Config

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run Web",
		Run: func(cmd *cobra.Command, args []string) {

			a := app.New(opt)
			err := a.Run()
			if err != nil {
				fmt.Printf("%+v\n", err)
				os.Exit(1)
			}
			fmt.Println("See you")
		},
	}

	flag := cmd.Flags()

	flag.StringVarP(&opt.Address, "address", "", ":8080", "http listen address")
	flag.StringVarP(&opt.SqlitePath, "sqlite-path", "", "/data/leetcode-badge/lc.db", "sqlite3 file path")
	flag.BoolVarP(&opt.EnableCron, "enable-cron", "", false, "if enable cron")

	// // basic
	// flag.BoolVarP(&opt.Debug, "debug", "d", false, "enable debug")
	// flag.StringVarP(&opt.Address, "address", "", ":8080", "http listen address")
	// flag.StringVarP(&opt.CacheType, "cache", "", "memory", "cache type, memory or redis")
	// flag.StringVarP(&opt.StorageType, "storage", "", "mysql", "storage type, only mysql")
	//
	// // cache
	// flag.StringVarP(&opt.RedisConfig.Address, "redis-address", "", "", "required when the cache type is redis")
	// flag.StringVarP(&opt.RedisConfig.Password, "redis-password", "", "", "optional when the type is redis")
	//
	// // mysql
	// flag.StringVarP(&opt.MySQLConfig.Host, "mysql-host", "", "", "required when the storage type is mysql")
	// flag.StringVarP(&opt.MySQLConfig.DBName, "mysql-database", "", "", "required when the storage type is mysql")
	// flag.StringVarP(&opt.MySQLConfig.User, "mysql-user", "", "", "required when the storage type is mysql")
	// flag.StringVarP(&opt.MySQLConfig.Password, "mysql-password", "", "", "required when the storage type is mysql")
	// flag.IntVarP(&opt.MySQLConfig.Port, "mysql-port", "", 3306, "required when the storage type is mysql")
	//
	// vBasic.BindPFlag("Debug", flag.Lookup("debug"))
	// vBasic.BindPFlag("Address", flag.Lookup("address"))
	// vBasic.BindPFlag("CacheType", flag.Lookup("cache"))
	// vBasic.BindPFlag("StorageType", flag.Lookup("storage"))
	//
	// vRedis.BindPFlag("Address", flag.Lookup("redis-address"))
	// vRedis.BindPFlag("Password", flag.Lookup("redis-password"))
	//
	// vMySQL.BindPFlag("Host", flag.Lookup("mysql-host"))
	// vMySQL.BindPFlag("DBName", flag.Lookup("mysql-database"))
	// vMySQL.BindPFlag("User", flag.Lookup("mysql-user"))
	// vMySQL.BindPFlag("Password", flag.Lookup("mysql-password"))
	// vMySQL.BindPFlag("Port", flag.Lookup("mysql-port"))
	//
	// vBasic.SetEnvPrefix(ENVPrefix)
	// vRedis.SetEnvPrefix(ENVPrefix + "_REDIS")
	// vMySQL.SetEnvPrefix(ENVPrefix + "_MYSQL")
	//
	// // LCB_DEBUG, LCB_ADDRESS ...
	// vBasic.BindEnv(
	// 	"Debug", "Address", "CacheType", "StorageType",
	// )
	//
	// // LCB_REDIS_ADDRESS, LCB_REDIS_PASSWORD
	// vRedis.BindEnv(
	// 	"Address", "Password",
	// )
	//
	// vMySQL.BindEnv(
	// 	"Host", "DBName", "User", "Password", "Port",
	// )
	//
	// vBasic.AutomaticEnv()
	// vRedis.AutomaticEnv()
	// vMySQL.AutomaticEnv()

	return cmd
}
