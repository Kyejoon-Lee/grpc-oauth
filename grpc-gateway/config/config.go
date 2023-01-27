package config

type Config struct {
	GatewayIP   string `env:GATEWAY_IP`
	GatewayPort string `env:"GATEWAY_PORT"`

	ClientID     string `env:"CLIENT_ID"`
	ClientSecret string `env:"CLIENT_SECRET"`

	ServerIP   string `env:"SERVER_IP"`
	ServerPort string `env:"SERVER_PORT"`
}

var config = &Config{
	GatewayIP:    "localhost",
	GatewayPort:  "9091",
	ClientID:     "757ff2f0153f3958ee4269936d1d701b",
	ClientSecret: "fCgU42MO4BKBn5S3diz63R67h3RL9NGu",
	ServerIP:     "localhost",
	ServerPort:   "9090",
}

func GetConfig() *Config {
	return config
}
