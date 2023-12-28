package main

import (
	"fmt"
	"log"
	"testing"
)

func TestStoreBlob(t *testing.T) {
	url := "http://localhost:22258"
	blobKey, err := StoreBlob(url+"/store", []byte("hello"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("blobKey", blobKey)
}

func TestGetBlob(t *testing.T) {
	url := "http://172.17.0.1:22258"
	blobHeight, err := GetBlob(url + "/get/tcelestia/818210/ce4fb46ee9c363488fa2162e1d82eb142c1de45c159f6b629c619a0aff1d840d")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("blobData", string(blobHeight))
}
