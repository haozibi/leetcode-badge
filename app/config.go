package app

type Config struct {
	Address      string
	DebugAddress int
	SqlitePath   string
	EnableCron   bool
}

type RedisConfig struct {
	Address  string
	Password string
}

type MySQLConfig struct {
	Host     string
	DBName   string
	User     string
	Password string
	Port     int
}
