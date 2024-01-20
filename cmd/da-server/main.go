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
	apiNearDa.HandleFunc("/get/{namespace}/{index}/{headerHash}", apis.ApiGetNearDA).Methods("GET")

	log.Fatal(http.ListenAndServe("0.0.0.0:22259", router))
}
