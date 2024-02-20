package funcs

import (
	b64 "encoding/base64"
	"github.com/celestiaorg/celestia-node/cmd/da-server/internal/arweave/config"
	"github.com/everFinance/goar"
	"github.com/everFinance/goar/types"
	"os"
)

func GetConfig() config.Config {
	return config.Config{
		RpcUrl:     "https://arweave.net",
		WalletFile: os.Getenv("ARWEAVE_WALLET"),
	}
}

func Wallet() (*goar.Wallet, error) {
	cnf := GetConfig()
	sEnc, err := b64.StdEncoding.DecodeString(cnf.WalletFile)
	if err != nil {
		return nil, err
	}
	wallet, err := goar.NewWallet(sEnc, cnf.RpcUrl)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func StoreData(data []byte) (*string, error) {
	wl, err := Wallet()
	if err != nil {
		return nil, err
	}

	tx, err := storeData(wl, data)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func storeData(wallet *goar.Wallet, data []byte) (*string, error) {

	tx, err := wallet.SendData(
		data, // Data bytes
		[]types.Tag{
			types.Tag{
				Name:  "testSendData",
				Value: "123",
			},
		},
	)

	if err != nil {
		return nil, err
	}

	return &tx.ID, nil
}

func GetData(hash string) ([]byte, error) {
	//TODO - implement me

	return nil, nil
}

func getData(wallet *goar.Wallet, hash string) ([]byte, error) {
	//TODO - implement me

	return nil, nil
}
