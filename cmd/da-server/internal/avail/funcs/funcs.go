package funcs

import (
	"fmt"
	"github.com/celestiaorg/celestia-node/cmd/da-server/internal/avail/config"
	"github.com/celestiaorg/celestia-node/cmd/da-server/internal/avail/extrinsics"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"time"
)

// funcs
func GetConfig() config.Config {
	return config.Config{
		Seed:   "verb jump guide coffee path squirrel hire verify gun robust rail fork",
		ApiURL: "https://goldberg.avail.tools/api",
		Size:   1000,
		AppID:  0,
		Dest:   "5H3qehpRTFiB3XwyTzYU43SpG7e8jW87vFug95fxdew76Vyi",
		Amount: 10,
	}
}

// submitData creates a transaction and makes a Avail data submission
func SubmitData() error {
	cnf := GetConfig()

	api, err := gsrpc.NewSubstrateAPI(cnf.ApiURL)
	if err != nil {
		return fmt.Errorf("cannot create api:%w", err)
	}

	// Set data and appID according to need
	data, _ := extrinsics.RandToken(cnf.Size)

	finalizedBlockCh := make(chan types.Hash, 1)
	go submitData(api, data, cnf.Seed, cnf.AppID, finalizedBlockCh)

	return nil
}

func QueryData(blockHash string, txIndex uint32) (*HeaderF, error) {
	_config := GetConfig()
	api, err := gsrpc.NewSubstrateAPI(_config.ApiURL)
	if err != nil {
		//panic(fmt.Sprintf("cannot create api client:%v", err))
		return nil, err
	}

	//var finalizedBlockCh = make(chan types.Hash)
	//go func() {
	//	err = extrinsics.SubmitData(api, "data", config.Seed, config.AppID, finalizedBlockCh)
	//	if err != nil {
	//		panic(fmt.Sprintf("cannot submit data:%v", err))
	//	}
	//}()

	// block hash to query proof
	//blockHash := <-finalizedBlockCh

	fmt.Printf("Transaction included in finalized block: %v\n", blockHash)
	h, _ := types.NewHashFromHexString(blockHash)
	transactionIndex := types.NewU32(txIndex)

	// query proof
	var data HeaderF
	err = api.Client.Call(&data, "kate_queryDataProof", transactionIndex, h)
	if err != nil {
		//panic(fmt.Sprintf("%v\n", err))
		return nil, err
	}
	fmt.Printf("Root:%v\n", data.Root.Hex())
	// print array of proof
	fmt.Printf("Proof:\n")
	for _, p := range data.Proof {
		fmt.Printf("%v\n", p.Hex())
	}

	fmt.Printf("Number of leaves: %v\n", data.Number_Of_Leaves)
	fmt.Printf("Leaf index: %v\n", data.Leaf_index)
	fmt.Printf("Leaf: %v\n", data.Leaf.Hex())

	return &data, nil
}

func submitData(api *gsrpc.SubstrateAPI, data string, seed string, appID int, finalizedBlockCh chan types.Hash) error {
	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return err
	}

	c, err := types.NewCall(meta, "DataAvailability.submit_data", types.NewBytes([]byte(data)))
	if err != nil {
		return fmt.Errorf("error creating new call: %s", err)
	}

	// Create the extrinsic
	ext := types.NewExtrinsic(c)

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return fmt.Errorf("error getting genesis hash: %s", err)
	}

	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return fmt.Errorf("error retrieveing runtime version: %s", err)
	}

	keyringPair, err := signature.KeyringPairFromSecret(seed, 42)
	if err != nil {
		return fmt.Errorf("error creating keyring pair: %s", err)
	}

	key, err := types.CreateStorageKey(meta, "System", "Account", keyringPair.PublicKey)
	if err != nil {
		return fmt.Errorf("cannot create storage key:%w", err)
	}

	var accountInfo types.AccountInfo
	ok, err := api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil || !ok {
		return fmt.Errorf("cannot get latest storage:%v", err)
	}

	nonce := uint32(accountInfo.Nonce)
	options := types.SignatureOptions{
		BlockHash:   genesisHash,
		Era:         types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash: genesisHash,
		Nonce:       types.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion: rv.SpecVersion,
		Tip:         types.NewUCompactFromUInt(100),
		//AppID:              types.NewUCompactFromUInt(uint64(appID)),
		TransactionVersion: rv.TransactionVersion,
	}

	// Sign the transaction using Alice's default account
	err = ext.Sign(keyringPair, options)
	if err != nil {
		return fmt.Errorf("cannot sign:%v", err)
	}

	// Send the extrinsic
	sub, err := api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		return fmt.Errorf("cannot submit extrinsic:%v", err)
	}

	defer sub.Unsubscribe()
	timeout := time.After(100 * time.Second)
	for {
		select {
		case status := <-sub.Chan():
			if status.IsInBlock {
				fmt.Printf("Txn inside block %v\n", status.AsInBlock.Hex())
			} else if status.IsFinalized {
				fmt.Printf("Txn inside finalized block\n")
				finalizedBlockCh <- status.AsFinalized
				return nil
			}
		case <-timeout:
			fmt.Printf("timeout of 100 seconds reached without getting finalized status for extrinsic")
			return nil
		}
	}
}
