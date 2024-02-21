package config

import "os"

type Config struct {
	Seed    string
	RPC     string
	ChainId string
}

func GetConfig() Config {
	config := Config{
		Seed:    "fiction stadium edge curious never romance enrich idea produce tennis witness struggle",
		RPC:     "https://jackal-testnet-rpc.polkachu.com:443",
		ChainId: "lupulella-2",
	}

	env := os.Getenv("api_env")

	if env == "mainnet" {
		config.Seed = os.Getenv("JACKAL_SEED_PHRASE")
		config.RPC = "https://rpc.jackalprotocol.com:443"
		config.ChainId = "jackal-1"
	}

	return config
}
