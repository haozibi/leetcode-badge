package app

type Config struct {
	Debug bool

	Address     string
	CronSpec    string
	CacheType   string
	StorageType string

	RedisConfig RedisConfig
	MySQLConfig MySQLConfig
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
