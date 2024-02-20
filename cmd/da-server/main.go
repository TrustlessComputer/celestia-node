package main

import (
	"fmt"
	"github.com/celestiaorg/celestia-node/cmd/da-server/apis"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {

	fmt.Println("Server start...")

	router := mux.NewRouter()

	// need to remove if updated code:
	router.HandleFunc("/store", apis.ApiStoreCelestia).Methods("POST")
	router.HandleFunc("/get/{namespace}/{height}/{commitment}", apis.ApiGetCelestia).Methods("GET")

	// elestia
	apiCelestia := router.PathPrefix("/celestia").Subrouter()
	apiCelestia.HandleFunc("/store", apis.ApiStoreCelestia).Methods("POST")
	apiCelestia.HandleFunc("/get/{namespace}/{height}/{commitment}", apis.ApiGetCelestia).Methods("GET")

	// eigenda
	apiEigenda := router.PathPrefix("/eigenda").Subrouter()
	apiEigenda.HandleFunc("/store", apis.ApiStoreEigenda).Methods("POST")
	apiEigenda.HandleFunc("/get/{namespace}/{index}/{headerHash}", apis.ApiGetEigenda).Methods("GET")

	// near da:
	apiNearDa := router.PathPrefix("/nearda").Subrouter()
	apiNearDa.HandleFunc("/store", apis.ApiStoreNearDA).Methods("POST", "GET")
	apiNearDa.HandleFunc("/get/{namespace}/{dataHex}", apis.ApiGetNearDA).Methods("GET")

	// Avail da: TODO
	apiAvail := router.PathPrefix("/avail").Subrouter()
	apiAvail.HandleFunc("/store", apis.ApiStoreAvail).Methods("POST", "GET")
	apiAvail.HandleFunc("/get/{namespace}/{txIndex}/{blockHash}", apis.ApiGetAvail).Methods("GET")

	// Jackal da: TODO
	//apiJackal := router.PathPrefix("/jackal").Subrouter()
	//apiJackal.HandleFunc("/store", apis.ApiStoreJackal).Methods("POST", "GET")
	//apiJackal.HandleFunc("/get/{namespace}/{dataHex}", apis.ApiGetJackal).Methods("GET")

	// Arweave da: TODO
	apiArweave := router.PathPrefix("/arweave").Subrouter()
	apiArweave.HandleFunc("/store", apis.ApiStoreArweave).Methods("POST", "GET")
	apiArweave.HandleFunc("/get/{namespace}/{dataHex}", apis.ApiGetArweave).Methods("GET")

	// IPFS - pinata
	apiIPFS := router.PathPrefix("/ipfs").Subrouter()
	apiIPFS.HandleFunc("/store", apis.ApiStoreIPFS).Methods("POST", "GET")
	apiIPFS.HandleFunc("/get/{namespace}/{ipfsHash}", apis.ApiGetIPFS).Methods("GET")

	//FileCoin da: TODO
	apiFileCoin := router.PathPrefix("/filecoin").Subrouter()
	apiFileCoin.HandleFunc("/store", apis.ApiStoreFileCoin).Methods("POST", "GET")
	apiFileCoin.HandleFunc("/get/{namespace}/{dataHex}", apis.ApiGetFileCoin).Methods("GET")

	//Syscoin da: TODO
	apiSysCoin := router.PathPrefix("/syscoin").Subrouter()
	apiSysCoin.HandleFunc("/store", apis.ApiStoreSysCoin).Methods("POST", "GET")
	apiSysCoin.HandleFunc("/get/{namespace}/{dataHex}", apis.ApiGetSysCoin).Methods("GET")

	log.Fatal(http.ListenAndServe("0.0.0.0:22258", router))
}
