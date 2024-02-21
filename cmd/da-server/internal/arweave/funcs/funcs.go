package funcs

import (
	b64 "encoding/base64"
	"fmt"
	"github.com/celestiaorg/celestia-node/cmd/da-server/internal/arweave/config"
	"github.com/everFinance/goar"
	"github.com/everFinance/goar/types"
	"time"
)

func Wallet() (*goar.Wallet, error) {
	cnf := config.GetConfig()
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

	txId, err := storeData(wl, data)
	if err != nil {
		return nil, err
	}

	// waiting here for get completed tx
	errTx := goar.ErrPendingTx
	// try 10 times for making sure tx success
	for i := 0; i < 100; i++ {
		if errTx != nil {
			_, errTx = getData(wl, *txId)
			if errTx != nil {
				fmt.Printf("try get arwear txid=%s status %v", *txId, errTx)
			}
		} else {
			break
		}
		time.Sleep(6 * time.Second) // max time is 60s
	}
	if errTx != nil {
		fmt.Printf("err arwear txid=%s err %v", *txId, errTx)
		return nil, errTx
	}
	fmt.Printf("arwear txid=%s", *txId)
	return txId, nil
}

func storeData(wallet *goar.Wallet, data []byte) (*string, error) {

	tx, err := wallet.SendData(
		data, // Data bytes
		[]types.Tag{
			types.Tag{
				Name:  "bvm.network",
				Value: "arweave-da",
			},
		},
	)

	if err != nil {
		return nil, err
	}

	return &tx.ID, nil
}

func GetData(hash string) ([]byte, error) {
	wl, err := Wallet()
	if err != nil {
		return nil, err
	}

	base64Data, err := getData(wl, hash)
	if err != nil {
		return nil, err
	}

	rawDecodedText, err := b64.StdEncoding.DecodeString(string(base64Data))
	if err != nil {
		//trick the short data
		rawDecodedText, err = b64.StdEncoding.DecodeString(fmt.Sprintf("%s=", string(base64Data)))
		if err != nil {
			return nil, err
		}
	}

	return rawDecodedText, nil

}

func getData(wallet *goar.Wallet, hash string) ([]byte, error) {
	_b, err := wallet.Client.GetTransactionData(hash)
	if err != nil {
		return nil, err
	}

	return _b, nil
}
