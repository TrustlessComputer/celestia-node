package config

import "os"

type Config struct {
	RpcUrl       string
	RetrieveFile string
	APIKey       string
	ChainID      int
	PinURL       string
	Env          string
}

func GetConfig() Config {
	_env := os.Getenv("FILE_COIN_ENV")
	if _env == "" {
		_env = "mainnet"
	}

	return Config{
		RpcUrl:       "https://node.lighthouse.storage/api/v0/add",
		APIKey:       os.Getenv("FILE_COIN_API_KEY"),
		RetrieveFile: "https://gateway.lighthouse.storage",
		PinURL:       "https://api.lighthouse.storage/api",
		Env:          _env,
	}
}
