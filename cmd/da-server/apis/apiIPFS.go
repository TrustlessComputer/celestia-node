package apis

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	_ipfs "github.com/celestiaorg/celestia-node/cmd/da-server/internal/ipfs/funcs"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

const (
	NAMESPACE_IPFS = "ipfsda"
)

func ApiStoreIPFS(w http.ResponseWriter, r *http.Request) {

	data := RequestData{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rawDecodedText, err := base64.StdEncoding.DecodeString(data.Data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ipfsHash, err := _ipfs.UploadData(fmt.Sprintf("f-%d", time.Now().UnixMicro()), rawDecodedText)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = w.Write([]byte(fmt.Sprintf("/%s/%s", NAMESPACE_IPFS, ipfsHash)))
	return

}

func ApiGetIPFS(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	namespace := vars["namespace"]
	if len(namespace) > 10 {
		namespace = namespace[:10]
	}

	ipfsHash := vars["ipfsHash"]

	data, err := _ipfs.GetData(ipfsHash)
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
