package cmd

import (
	"fmt"
	"os"

	"github.com/haozibi/leetcode-badge/app"
	"github.com/haozibi/zlog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// ENVPrefix default env prefix
	ENVPrefix = "LCB"
)

// NewRunCommand new run command
func NewRunCommand() *cobra.Command {

	var opt app.Config
	vBasic := viper.New()
	vRedis := viper.New()
	vMySQL := viper.New()

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run Web",
		Run: func(cmd *cobra.Command, args []string) {

			var err error

			err = vBasic.Unmarshal(&opt)
			er(err)
			err = vRedis.Unmarshal(&opt.RedisConfig)
			er(err)
			err = vMySQL.Unmarshal(&opt.MySQLConfig)
			er(err)

			if opt.Debug {
				zlog.NewBasicLog(os.Stdout, zlog.WithDebug(true))
			}

			opt.CronSpec = "30 8 * * *"

			a := app.New(opt)
			err = a.Run()
			if err != nil {
				if opt.Debug {
					fmt.Printf("%+v\n", err)
				} else {
					fmt.Println(err)
				}
				os.Exit(1)
			}
			fmt.Println("See you")
		},
	}

	flag := cmd.Flags()

	// basic
	flag.BoolVarP(&opt.Debug, "debug", "d", false, "enable debug")
	flag.StringVarP(&opt.Address, "address", "", ":8080", "http listen address")
	flag.StringVarP(&opt.CacheType, "cache", "", "memory", "cache type, memory or redis")
	flag.StringVarP(&opt.StorageType, "storage", "", "mysql", "storage type, only mysql")

	// cache
	flag.StringVarP(&opt.RedisConfig.Address, "redis-address", "", "", "required when the cache type is redis")
	flag.StringVarP(&opt.RedisConfig.Password, "redis-password", "", "", "optional when the type is redis")

	// mysql
	flag.StringVarP(&opt.MySQLConfig.Host, "mysql-host", "", "", "required when the storage type is mysql")
	flag.StringVarP(&opt.MySQLConfig.DBName, "mysql-database", "", "", "required when the storage type is mysql")
	flag.StringVarP(&opt.MySQLConfig.User, "mysql-user", "", "", "required when the storage type is mysql")
	flag.StringVarP(&opt.MySQLConfig.Password, "mysql-password", "", "", "required when the storage type is mysql")
	flag.IntVarP(&opt.MySQLConfig.Port, "mysql-port", "", 3306, "required when the storage type is mysql")

	vBasic.BindPFlag("Debug", flag.Lookup("debug"))
	vBasic.BindPFlag("Address", flag.Lookup("address"))
	vBasic.BindPFlag("CacheType", flag.Lookup("cache"))
	vBasic.BindPFlag("StorageType", flag.Lookup("storage"))

	vRedis.BindPFlag("Address", flag.Lookup("redis-address"))
	vRedis.BindPFlag("Password", flag.Lookup("redis-password"))

	vMySQL.BindPFlag("Host", flag.Lookup("mysql-host"))
	vMySQL.BindPFlag("DBName", flag.Lookup("mysql-database"))
	vMySQL.BindPFlag("User", flag.Lookup("mysql-user"))
	vMySQL.BindPFlag("Password", flag.Lookup("mysql-password"))
	vMySQL.BindPFlag("Port", flag.Lookup("mysql-port"))

	vBasic.SetEnvPrefix(ENVPrefix)
	vRedis.SetEnvPrefix(ENVPrefix + "_REDIS")
	vMySQL.SetEnvPrefix(ENVPrefix + "_MYSQL")

	// LCB_DEBUG, LCB_ADDRESS ...
	vBasic.BindEnv(
		"Debug", "Address", "CacheType", "StorageType",
	)

	// LCB_REDIS_ADDRESS, LCB_REDIS_PASSWORD
	vRedis.BindEnv(
		"Address", "Password",
	)

	vMySQL.BindEnv(
		"Host", "DBName", "User", "Password", "Port",
	)

	vBasic.AutomaticEnv()
	vRedis.AutomaticEnv()
	vMySQL.AutomaticEnv()

	return cmd
}

func er(err error) {
	if err == nil {
		return
	}
	fmt.Println("error:", err)
	os.Exit(1)
}
