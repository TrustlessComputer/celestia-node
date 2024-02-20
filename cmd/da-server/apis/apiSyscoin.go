package apis

import (
	"encoding/base64"
	"fmt"
	_syscoin "github.com/celestiaorg/celestia-node/cmd/da-server/internal/syscoin/funcs"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	NAMESPACE_SYSCOIN = "syscoinda"
)

func ApiStoreSysCoin(w http.ResponseWriter, r *http.Request) {
	//TODO - implement me
	data, err := DecodeReqBody(r)
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

	resultHex, err := _syscoin.UploadData(decodedBytes)
	if err != nil {
		fmt.Println("submit err:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = w.Write([]byte(fmt.Sprintf("/%s/%s", NAMESPACE_SYSCOIN, resultHex)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	return

}

func ApiGetSysCoin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	namespace := vars["namespace"]
	if len(namespace) > 10 {
		namespace = namespace[:10]
	}

	versionhash_or_txid := vars["versionhash_or_txid"]

	data, err := _syscoin.GetData(versionhash_or_txid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	return
}
