package apis

import (
	"errors"
	"fmt"
	near "github.com/near/rollup-data-availability/gopkg/da-rpc"
	"net/http"
	"time"
)

const (
	DA_KEY      = "ed25519:5rruwJXodZu6phNsApFcAm9LFxSy7nYpwnCB8vQDAvJKVgDZ424uGyXQiHQGTM3sbeBkvVXKbGxFiuswXaNRWKjv"
	DA_CONTRACT = "8363d67bf00d1aea294285ad3fd1066378ac502d25995cf4a474326535a37201"
	DA_ACCOUNT  = "8363d67bf00d1aea294285ad3fd1066378ac502d25995cf4a474326535a37201"
)

func ApiTestNearDA(w http.ResponseWriter, r *http.Request) {

	config, err := near.NewConfig(DA_ACCOUNT, DA_CONTRACT, DA_KEY, 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("config", config)
	if config.Namespace.Id != 1 {
		err = errors.New("Expected namespace id to be equal")
	}
	if config.Namespace.Version != 0 {
		err = errors.New("Expected namespace version to be equal")
	}
	http.Error(w, err.Error(), http.StatusOK)

	return

}

func ApiStoreNearDA(w http.ResponseWriter, r *http.Request) {

	config, err := near.NewConfig(DA_ACCOUNT, DA_CONTRACT, DA_KEY, 1)
	if err != nil {
		fmt.Println("NewConfig err:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var candidateHex string = "0xfF00000000000000000000000000000000000000"

	// var data []byte = []byte("elvis")

	bytes := make([]byte, 64)
	copy(bytes, []byte("1111111111111111111111111111111122222222222222222222222222222222"))

	result, err := config.Submit(candidateHex, bytes)

	if err != nil {
		fmt.Println("submit err:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	time.Sleep(30 * time.Second)

	blob, err := config.Get(result, 0)
	if err != nil {
		fmt.Println("config.Get err:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("blob byte: ", blob)
	fmt.Println("submit result byte:", result)

	frameRef := near.FrameRef{}
	err = frameRef.UnmarshalBinary(result)
	if err != nil {
		fmt.Println("UnmarshalBinary err:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("frameRef.TxId", frameRef.TxId, "frameRef.TxCommitment", frameRef.TxCommitment)

	fmt.Println("frameRef.TxId.String()", string(frameRef.TxId), "frameRef.TxCommitment.String()", string(frameRef.TxCommitment))

	if string(frameRef.TxId) != "11111111111111111111111111111111" {
		err = errors.New("Expected id to be equal")
		fmt.Println(frameRef.TxId, " err:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}
	if string(frameRef.TxCommitment) != "22222222222222222222222222222222" {
		err = errors.New("Expected commitment to be equal")
		fmt.Println(frameRef.TxCommitment, " err:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	return

}

func ApiGetNearDA(w http.ResponseWriter, r *http.Request) {

	id := make([]byte, 32)
	copy(id, []byte("11111111111111111111111111111111"))
	commitment := make([]byte, 32)
	copy(commitment, []byte("22222222222222222222222222222222"))
	frameRef := near.FrameRef{
		TxId:         id,
		TxCommitment: commitment,
	}
	binary, err := frameRef.MarshalBinary()
	println("binary, id, commitment", binary, id, commitment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	config, err := near.NewConfig(DA_ACCOUNT, DA_CONTRACT, DA_KEY, 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	blob, err := config.Get(binary, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("blob: ", blob)
	return
}
