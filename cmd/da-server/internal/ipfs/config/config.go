package config

import "os"

type Config struct {
	JWT    string
	API    string
	GetAPI string
}

func GetConfig() Config {
	return Config{
		API:    "https://api.pinata.cloud",
		GetAPI: "https://nbc-alpha.mypinata.cloud/ipfs",
		JWT:    os.Getenv("PINATA_JWT"),
	}
}
