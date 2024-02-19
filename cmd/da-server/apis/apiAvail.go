package apis

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/celestiaorg/celestia-node/cmd/da-server/internal/avail/config"
	"github.com/celestiaorg/celestia-node/cmd/da-server/internal/avail/extrinsics"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/davecgh/go-spew/spew"
	"net/http"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
)

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

	_ = decodedBytes
	//TODO - implement me
	cnf := getConfig()
	err = submitData(cnf.Size, cnf.ApiURL, cnf.Seed, cnf.AppID)
	if err != nil {
		panic(fmt.Sprintf("cannot submit data:%v", err))
	}

	return

}

func ApiGetAvail(w http.ResponseWriter, r *http.Request) {
	//TODO - implement me
	return
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
func submitData(size int, ApiURL string, Seed string, AppID int) error {
	api, err := gsrpc.NewSubstrateAPI(ApiURL)
	if err != nil {
		return fmt.Errorf("cannot create api:%w", err)
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return fmt.Errorf("cannot get metadata:%w", err)
	}

	// Set data and appID according to need
	data, _ := extrinsics.RandToken(size)
	appID := 0

	// if app id is greater than 0 then it must be created before submitting data
	if AppID != 0 {
		appID = AppID
	}

	c, err := types.NewCall(meta, "DataAvailability.submit_data", types.NewBytes([]byte(data)))
	if err != nil {
		return fmt.Errorf("cannot create new call:%w", err)
	}

	// Create the extrinsic
	ext := types.NewExtrinsic(c)

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return fmt.Errorf("cannot get block hash:%w", err)
	}

	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return fmt.Errorf("cannot get runtime version:%w", err)
	}

	keyringPair, err := signature.KeyringPairFromSecret(Seed, 42)
	if err != nil {
		return fmt.Errorf("cannot create LeyPair:%w", err)
	}

	spew.Dump(string(keyringPair.PublicKey))
	spew.Dump(keyringPair.Address)
	spew.Dump(keyringPair.URI)

	key, err := types.CreateStorageKey(meta, "System", "Account", keyringPair.PublicKey)
	if err != nil {
		return fmt.Errorf("cannot create storage key:%w", err)
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
		return fmt.Errorf("cannot sign:%w", err)
	}

	// Send the extrinsic
	hash, err := api.RPC.Author.SubmitExtrinsic(ext)
	if err != nil {
		spew.Dump(err)
		return fmt.Errorf("cannot submit extrinsic:%w", err)
	}
	fmt.Printf("Data submitted by Alice: %v against appID %v  sent with hash %#x\n", data, appID, hash)

	return nil
}
