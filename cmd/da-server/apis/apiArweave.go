package apis

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	_arweave "github.com/celestiaorg/celestia-node/cmd/da-server/internal/arweave/funcs"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	NAMESPACE_ARWEAVE = "arweave"
)

func ApiStoreArweave(w http.ResponseWriter, r *http.Request) {

	type RequestData struct {
		Data string `json:"data"`
	}

	data := RequestData{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rawData, err := base64.StdEncoding.DecodeString(data.Data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hash, err := _arweave.StoreData(rawData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = w.Write([]byte(fmt.Sprintf("/%s/%s", NAMESPACE_ARWEAVE, hash)))
	return

}

func ApiGetArweave(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	namespace := vars["namespace"]
	if len(namespace) > 10 {
		namespace = namespace[:10]
	}

	hash := vars["dataHex"]
	data, err := _arweave.GetData(hash)
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
