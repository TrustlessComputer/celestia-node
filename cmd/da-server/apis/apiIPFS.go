package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	NAMESPACE_IPFS = "ipfsda"
)

func ApiStoreIPFS(w http.ResponseWriter, r *http.Request) {
	apiKey := ""
	type RequestData struct {
		Data string `json:"data"`
	}
	data := RequestData{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//TODO - implement me
	ipfsHash, err := UploadData(apiKey, "", []byte(data.Data))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = w.Write([]byte(fmt.Sprintf("/%s/%s", NAMESPACE_IPFS, ipfsHash)))
	return

}

func ApiGetIPFS(w http.ResponseWriter, r *http.Request) {
	//TODO - implement me
	return
}
