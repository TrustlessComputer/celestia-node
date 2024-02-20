package apis

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	_avail "github.com/celestiaorg/celestia-node/cmd/da-server/internal/avail/funcs"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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

	hash, txIndex, err := _avail.SubmitData(decodedBytes)
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
	d, err := _avail.QueryData(blockHash, uint32(txIndex))
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
