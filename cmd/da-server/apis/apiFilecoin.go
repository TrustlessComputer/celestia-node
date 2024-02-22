package apis

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	_fileCoin "github.com/celestiaorg/celestia-node/cmd/da-server/internal/filecoin/funcs"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	NAMESPACE_FILECOIN = "filecoin"
)

func ApiStoreFileCoin(w http.ResponseWriter, r *http.Request) {
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

	hash, err := _fileCoin.StoreData(rawData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = w.Write([]byte(fmt.Sprintf("/%s/%s", NAMESPACE_FILECOIN, *hash)))
	return

}

func ApiGetFileCoin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	namespace := vars["namespace"]
	if len(namespace) > 10 {
		namespace = namespace[:10]
	}

	hash := vars["cid"]
	data, err := _fileCoin.GetData(hash)
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
