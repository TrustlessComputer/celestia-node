package main

import (
	"fmt"
	"log"
	"testing"
)

func TestStoreBlob(t *testing.T) {
	url := "http://localhost:22258"
	blobKey, err := StoreBlob(url+"/store/abcde", []byte("hello"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("blobKey", blobKey)
}

func TestGetBlob(t *testing.T) {
	url := "http://localhost:22258"
	blobHeight, err := GetBlob(url + "/get/abcde/812235/e028045f77b2f204d547ee7e46fba108fcd8c77160b68ece8f424c7d2c5384ff")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("blobData", string(blobHeight))
}
