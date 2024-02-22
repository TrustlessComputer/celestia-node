package config

type Config struct {
	RpcUrl     string
	GetInfoURL string
	APIKey     string
	ChainID    int
}

func GetConfig() Config {
	return Config{
		RpcUrl:     "https://node.lighthouse.storage/api/v0/add",
		APIKey:     "6e51f484.",
		GetInfoURL: "https://gateway.lighthouse.storage",
	}
}
