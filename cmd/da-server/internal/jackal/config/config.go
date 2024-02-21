package config

import "os"

type Config struct {
	Seed    string
	RPC     string
	ChainId string
}

func GetConfig() Config {
	config := Config{
		Seed:    "slim odor fiscal swallow piece tide naive river inform shell dune crunch canyon ten time universe orchard roast horn ritual siren cactus upon forum",
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
