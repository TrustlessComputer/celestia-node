package config

type Config struct {
	RpcUrl  string
	APIKey  string
	ChainID int
}

func GetConfig() Config {
	return Config{
		RpcUrl: "https://node.lighthouse.storage/api/v0/add",
		APIKey: "892a9858._",
	}
}
