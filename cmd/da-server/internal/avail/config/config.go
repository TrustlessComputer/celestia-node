package config

import "os"

type Config struct {
	Seed   string `json:"seed"`
	ApiURL string `json:"api_url"`
	Size   int    `json:"size"`
	AppID  int    `json:"app_id"`
	Dest   string `json:"dest"`
	Amount uint64 `json:"amount"`
}

func GetConfig() Config {
	return Config{
		Seed:   os.Getenv("AVAIL_SEED"),
		ApiURL: "wss://goldberg.avail.tools/ws",
	}
}
