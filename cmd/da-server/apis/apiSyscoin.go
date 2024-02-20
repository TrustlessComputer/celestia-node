package apis

import (
	"encoding/base64"
	"fmt"
	"github.com/celestiaorg/celestia-node/cmd/da-server/internal/syscoin/funcs"
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

	resultHex, err := funcs.UploadData("", decodedBytes)
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
	//TODO - implement me
	return
}
