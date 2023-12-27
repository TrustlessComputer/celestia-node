package main

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/celestiaorg/celestia-node/api/rpc/client"
	"github.com/celestiaorg/celestia-node/blob"
	"github.com/celestiaorg/celestia-node/share"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func storeBlob(rpc *client.Client, data []byte, ns []byte) ([]byte, uint64, error) {
	var nameSpace, _ = share.NewBlobNamespaceV0(ns)
	newBlob, err := blob.NewBlob(
		0,
		nameSpace,
		data,
	)
	if err != nil {
		return nil, 0, err
	}
	submit, err := rpc.Blob.Submit(context.Background(), []*blob.Blob{newBlob}, nil)
	if err != nil {
		return nil, 0, err
	}
	return newBlob.Commitment, submit, nil
}

func getBlob(rpc *client.Client, height uint64, commitment []byte, ns []byte) ([]byte, error) {
	var nameSpace, _ = share.NewBlobNamespaceV0(ns)
	data, err := rpc.Blob.Get(context.Background(), height, nameSpace, blob.Commitment(commitment))
	if err != nil {
		return nil, err
	}
	return data.Data, nil
}

func main() {
	var CELESTIS_NODE = "http://localhost:26658"
	var JWTToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJwdWJsaWMiLCJyZWFkIiwid3JpdGUiLCJhZG1pbiJdfQ.D1NtqUqJIX_FzFgiapuX3GMJoSeCT1-tg7XmxdlZmA0"
	fmt.Println("Server start...")
	client, err := client.NewClient(context.Background(), CELESTIS_NODE, JWTToken)
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/store/{namespace}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		namespace := vars["namespace"]
		if len(namespace) > 10 {
			namespace = namespace[:10]
		}

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

		commitment, height, err := storeBlob(client, decodedBytes, []byte(namespace))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = w.Write([]byte(fmt.Sprintf("%d/%s", height, hex.EncodeToString(commitment))))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		return
	}).Methods("POST")

	router.HandleFunc("/get/{namespace}/{height}/{commitment}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		namespace := vars["namespace"]
		if len(namespace) > 10 {
			namespace = namespace[:10]
		}

		heightStr := vars["height"]
		height, err := strconv.Atoi(heightStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		commitment, err := hex.DecodeString(vars["commitment"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println("namespace", namespace)
		data, err := getBlob(client, uint64(height), commitment, []byte(namespace))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println("data", string(data))
		_, err = w.Write(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		return
	}).Methods("GET")

	log.Fatal(http.ListenAndServe("0.0.0.0:22258", router))
}
