package funcs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"os"
)

type SyscoinConfig struct {
	rpcURL  string
	chainId string
}

func getSyscoinConfig() (config SyscoinConfig) {
	config.rpcURL = "https://rpc.tanenbaum.io"
	config.chainId = "5700"

	env := os.Getenv("api_env")

	if env == "mainnet" {
		config.rpcURL = "https://rpc.syscoin.org"
		config.chainId = "57"
	}
	return config
}

func UploadData(fileName string, data []byte) (string, error) {
	return "", nil
}

func GetData(hash string) ([]byte, error) {
	config := getSyscoinConfig()
	client := &http.Client{}
	requestData := `{"jsonrpc": "1.0", "id": "curltest", "method": "getnevmblobdata", "params": ["` + hash + `"]}`
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	req, err := http.NewRequest(
		"POST",
		config.rpcURL,
		bytes.NewReader(jsonData),
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	fmt.Println("body", string(body))
	return body, nil
}
