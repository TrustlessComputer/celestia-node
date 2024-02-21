package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"jackalda/apis"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Server start...")

	router := mux.NewRouter()

	//jackal:
	apiJackal := router.PathPrefix("/jackal").Subrouter()
	apiJackal.HandleFunc("/store", apis.ApiStoreJackal).Methods("POST", "GET")
	apiJackal.HandleFunc("/get/{namespace}/{fileName}", apis.ApiGetJackal).Methods("GET")

	log.Fatal(http.ListenAndServe("0.0.0.0:22259", router))
}
