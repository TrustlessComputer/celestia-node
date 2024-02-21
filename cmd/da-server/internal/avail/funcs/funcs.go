package funcs

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/celestiaorg/celestia-node/cmd/da-server/internal/avail/config"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"time"
)

// submitData creates a transaction and makes a Avail data submission
func SubmitData(data []byte) (*string, *uint32, error) {
	cnf := config.GetConfig()
	//Size := cnf.Size
	ApiURL := cnf.ApiURL
	Seed := cnf.Seed
	//AppID := cnf.AppID

	api, err := gsrpc.NewSubstrateAPI(ApiURL)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot create api:%w", err)
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return nil, nil, fmt.Errorf("cannot get metadata:%w", err)
	}

	// Set data and appID according to need
	//data, _ := extrinsics.RandToken(size)
	//appID := 0
	//
	//// if app id is greater than 0 then it must be created before submitting data
	//if AppID != 0 {
	//	appID = AppID
	//}

	subData := hex.EncodeToString(data)
	c, err := types.NewCall(meta, "DataAvailability.submit_data", types.NewBytes([]byte(subData)))
	if err != nil {
		return nil, nil, fmt.Errorf("cannot create new call:%w", err)
	}

	// Create the extrinsic
	ext := types.NewExtrinsic(c)

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot get block hash:%w", err)
	}

	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return nil, nil, fmt.Errorf("cannot get runtime version:%w", err)
	}

	keyringPair, err := signature.KeyringPairFromSecret(Seed, 42)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot create LeyPair:%w", err)
	}

	key, err := types.CreateStorageKey(meta, "System", "Account", keyringPair.PublicKey)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot create storage key:%w", err)
	}

	var accountInfo types.AccountInfo
	nonce := uint32(1)
	ok, err := api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err == nil && ok {
		nonce = uint32(accountInfo.Nonce)
	}

	o := types.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		AppID:              types.NewUCompactFromUInt(uint64(cnf.AppID)),
		TransactionVersion: rv.TransactionVersion,
	}

	// Sign the transaction using Alice's default account
	err = ext.Sign(keyringPair, o)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot sign:%w", err)
	}

	// Send the extrinsic
	//block hash
	sub, err := api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot submit extrinsic:%w", err)
	}

	defer sub.Unsubscribe()
	//timeout := time.After(100 * time.Second)
	timeout := time.After(600 * time.Second) //5 min
	h := ""

timeout_break:
	for {
		select {
		case status := <-sub.Chan():
			if status.IsInBlock {
				fmt.Printf("Txn inside block %v\n", status.AsInBlock.Hex())
			} else if status.IsFinalized {
				hash := status.AsFinalized
				err := getData(hash, api, subData)

				fmt.Printf("Txn inside finalized block with data: %s \n", subData)
				if err != nil {
					//panic(fmt.Sprintf("cannot get data:%v", err))
					return nil, nil, err
				}

				h = hash.Hex()
				break timeout_break
			}
		case <-timeout:
			break timeout_break
		}
	}

	if h == "" {
		err = errors.New("timeout of 100 seconds reached without getting finalized status for extrinsic")
		return nil, nil, err
	}

	return &h, &nonce, nil
}

func QueryData(blockHash string, nonce int64) ([]byte, error) {
	_config := config.GetConfig()
	api, err := gsrpc.NewSubstrateAPI(_config.ApiURL)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Transaction included in finalized block: %v\n", blockHash)
	h, _ := types.NewHashFromHexString(blockHash)

	data, err := getDataString(h, nonce, api)
	if err != nil {
		return nil, err
	}

	str := string(data)
	slice := str[2:]

	dataBytes, err := hex.DecodeString(slice)
	if err != nil {
		str = string(data)
		slice = str[1:]
		dataBytes, err = hex.DecodeString(slice)
		if err != nil {
			return nil, err
		}

		return dataBytes, nil

	}

	return dataBytes, nil
}

// getData extracts data from the block and compares it
func getData(hash types.Hash, api *gsrpc.SubstrateAPI, data string) error {
	block, err := api.RPC.Chain.GetBlock(hash)
	if err != nil {
		return fmt.Errorf("cannot get block by hash:%w", err)
	}

	_ = block
	//for _, ext := range block.Block.Extrinsics {
	//	// these values below are specific indexes only for data submission, differs with each extrinsic
	//	if ext.Method.CallIndex.SectionIndex == 29 && ext.Method.CallIndex.MethodIndex == 1 {
	//		arg := ext.Method.Args
	//		str := string(arg)
	//		slice := str[2:]
	//		fmt.Println("string value", slice)
	//		fmt.Println("data", data)
	//		if slice == data {
	//			fmt.Println("Data found in block")
	//		}
	//	}
	//}

	return nil
}

func getDataString(hash types.Hash, index int64, api *gsrpc.SubstrateAPI) ([]byte, error) {
	block, err := api.RPC.Chain.GetBlock(hash)
	if err != nil {
		return nil, fmt.Errorf("cannot get block by hash:%w", err)
	}
	for _, ext := range block.Block.Extrinsics {
		// these values below are specific indexes only for data submission, differs with each extrinsic
		if ext.Method.CallIndex.SectionIndex == 29 && ext.Method.CallIndex.MethodIndex == 1 && ext.Signature.Nonce.Int64() == index {
			arg := ext.Method.Args
			return arg, nil
		}
	}
	return []byte{}, nil
}
