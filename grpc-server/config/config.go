package config

type Config struct {
	DBAdapter  string `env:"DB_ADAPTER"`
	DBHost     string `env:"DB_HOST"`
	DBPort     string `env:"DB_PORT"`
	DBUsername string `env:"DB_USER" yaml:"-"`
	DBPassword string `env:"DB_PASSWORD" yaml:"-"`
	DBName     string `env:"DB_NAME"`

	ClientID     string `env:"CLIENT_ID"`
	ClientSecret string `env:"CLIENT_SECRET"`

	TimeZone string `env:"TIMEZONE"`
}

var config = &Config{
	DBAdapter:  "postgres",
	DBHost:     "localhost",
	DBPort:     "5432",
	DBUsername: "joon",
	DBName:     "toy-db",
	DBPassword: "test1234",
	TimeZone:   "Asia/Seoul",
}

func GetConfig() *Config {
	return config
}
