package apis

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/celestiaorg/celestia-node/cmd/da-server/internal/avail/config"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
)

const NAMESPACE_AVAIL = "avail"

func ApiStoreAvail(w http.ResponseWriter, r *http.Request) {
	type RequestData struct {
		Data string `json:"data"`
	}
	data := RequestData{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Use the data
	decodedBytes, err := base64.StdEncoding.DecodeString(data.Data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cnf := getConfig()
	hash, txIndex, err := submitData(cnf.Size, cnf.ApiURL, cnf.Seed, cnf.AppID, decodedBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = w.Write([]byte(fmt.Sprintf("/%s/%d/%s", NAMESPACE_AVAIL, *txIndex, *hash)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	return
}

func ApiGetAvail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	if len(namespace) > 10 {
		namespace = namespace[:10]
	}

	txIndexStr := vars["txIndex"]
	txIndex, err := strconv.Atoi(txIndexStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	blockHash := vars["blockHash"]
	d, err := queryData(blockHash, uint32(txIndex))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//TODO - verify here?
	_, err = w.Write([]byte(d.Root.Hex()))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// funcs
func getConfig() config.Config {
	return config.Config{
		Seed:   "verb jump guide coffee path squirrel hire verify gun robust rail fork",
		ApiURL: "wss://goldberg.avail.tools/ws",
		Size:   1000,
		AppID:  0,
		Dest:   "5H3qehpRTFiB3XwyTzYU43SpG7e8jW87vFug95fxdew76Vyi",
		Amount: 10,
	}
}

// submitData creates a transaction and makes a Avail data submission
func submitData(size int, ApiURL string, Seed string, AppID int, data []byte) (*string, *uint32, error) {
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

	c, err := types.NewCall(meta, "DataAvailability.submit_data", types.NewBytes(data))
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

	spew.Dump(string(keyringPair.PublicKey))
	spew.Dump(keyringPair.Address)
	spew.Dump(keyringPair.URI)

	key, err := types.CreateStorageKey(meta, "System", "Account", keyringPair.PublicKey)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot create storage key:%w", err)
	}

	var accountInfo types.AccountInfo
	nonce := uint32(1)
	ok, err := api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err == nil && !ok {
		nonce = uint32(accountInfo.Nonce)
	}

	o := types.SignatureOptions{
		BlockHash:   genesisHash,
		Era:         types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash: genesisHash,
		Nonce:       types.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion: rv.SpecVersion,
		Tip:         types.NewUCompactFromUInt(0),
		//AppID:              types.NewUCompactFromUInt(uint64(appID)),
		TransactionVersion: rv.TransactionVersion,
	}

	// Sign the transaction using Alice's default account
	err = ext.Sign(keyringPair, o)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot sign:%w", err)
	}

	// Send the extrinsic
	//block hash
	hash, err := api.RPC.Author.SubmitExtrinsic(ext)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot submit extrinsic:%w", err)
	}

	fmt.Printf("Data submitted: %v against appID %v  sent with hash %#x\n", data, AppID, hash)
	d := hash.Hex()
	txIndex := uint32(1)

	return &d, &txIndex, nil
}

type HeaderF struct {
	Root             types.Hash
	Proof            []types.Hash
	Number_Of_Leaves int
	Leaf_index       int
	Leaf             types.Hash
}

func queryData(blockHash string, txIndex uint32) (*HeaderF, error) {
	_config := getConfig()
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
