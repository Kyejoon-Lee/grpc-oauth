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
	GatewayIP:   "localhost",
	GatewayPort: "9091",

	ServerIP:   "localhost",
	ServerPort: "9090",
}

func GetConfig() *Config {
	return config
}
