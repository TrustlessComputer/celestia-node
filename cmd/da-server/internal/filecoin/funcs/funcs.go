package funcs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/celestiaorg/celestia-node/cmd/da-server/internal/filecoin/config"
	"github.com/tendermint/tendermint/types/time"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func dealParams() DealParams {
	return DealParams{
		NumCopies:       3, //maximum
		RenewThreshold:  240,
		RepairThreshold: 28800,
	}
}

func dealParamsString() string {
	dps := dealParams()
	_b, err := json.Marshal(dps)
	if err != nil {
		return ""
	}

	return string(_b)
}

func createTheUploadedFile(data []byte) (*os.File, *string, error) {
	fn := fmt.Sprintf("file-%d", time.Now().UTC().UnixMicro())
	f, err := os.Create(fn)
	if err != nil {
		return nil, nil, err
	}

	_, err = f.Write(data)
	if err != nil {
		return nil, nil, err
	}

	return f, &fn, nil
}

func StoreData(data []byte) (*string, error) {
	cnf := config.GetConfig()

	//TODO - implement me
	//client := lighthouse.NewClient("892a9858.220d782af3e14bb6a05589a51e4898ab")
	token := cnf.RpcUrl
	apiUrl := cnf.RpcUrl

	//create the uploaded file
	_, fn1, err := createTheUploadedFile(data)
	if err != nil {
		return nil, err
	}
	fn := *fn1

	// Create a buffer to store the request body
	var buf bytes.Buffer

	// Create a new multipart writer with the buffer
	w := multipart.NewWriter(&buf)

	// Add a file to the request
	file, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a new form field
	fw, err := w.CreateFormFile("file", fn)
	if err != nil {
		return nil, err
	}

	// Copy the contents of the file to the form field
	if _, err := io.Copy(fw, file); err != nil {
		return nil, err
	}

	// Close the multipart writer to finalize the request
	w.Close()

	// Send the request
	req, err := http.NewRequest("POST", apiUrl, &buf)
	if err != nil {
		return nil, err
	}

	dps := dealParamsString()

	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("X-Deal-Parameter", dps)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fResp := &FileCoinResp{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, fResp)
	if err != nil {
		return nil, err
	}

	os.Remove(fn)
	return &fResp.Hash, nil
}

func GetData(cid string) ([]byte, error) {
	//TODO - implement me

	return nil, nil
}
