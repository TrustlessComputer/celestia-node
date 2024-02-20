package config

import "os"

type Config struct {
	RpcUrl     string
	WalletFile string
}

func GetConfig() Config {
	return Config{
		RpcUrl:     "https://arweave.net",
		WalletFile: os.Getenv("ARWEAVE_WALLET"),
	}
}
