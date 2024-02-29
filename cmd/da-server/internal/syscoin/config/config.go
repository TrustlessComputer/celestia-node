package config

import "os"

type Config struct {
	RpcURL   string
	Password string
	User     string
}

func GetConfig() (config Config) {
	config.RpcURL = os.Getenv("SYSCOIN_RPC")
	config.User = os.Getenv("SYSCOIN_USER")
	config.Password = os.Getenv("SYSCOIN_PASSWORD")
	return config
}
