package funcs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/celestiaorg/celestia-node/cmd/da-server/internal/filecoin/config"
	"github.com/tendermint/tendermint/types/time"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func dealParams(env string) DealParams {
	return DealParams{
		NumCopies:       2, //maximum
		RenewThreshold:  240,
		RepairThreshold: 28800,
		Network:         env,
		AddMockData:     2,
	}
}

// Default parameters set. All RaaS workers enabled, any miners can take the deal. 2 MiB mock file added.
func dealParamsDefault(env string) DealParams {
	return DealParams{
		NumCopies:       3, //maximum
		RenewThreshold:  240,
		RepairThreshold: 28800,
	}
}

func dealParamsMock(env string) DealParams {
	return DealParams{
		AddMockData: 4,
		Network:     env,
	}
}

func dealParamsString(dpsType string, env string) string {
	dps := DealParams{}
	switch dpsType {
	case "default":
		dps = dealParamsDefault(env)
	case "mock":
		dps = dealParamsMock(env)
	default:
		dps = dealParams(env)
	}

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

func pinFile(fileName, cid string) (interface{}, error) {
	cnf := config.GetConfig()

	urlLink := fmt.Sprintf("%s/lighthouse/pin", cnf.PinURL)
	reqBody := FileCoinReq{
		FileName: fileName,
		Cid:      cid,
	}

	_b, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"POST",
		urlLink,
		bytes.NewReader(_b),
	)

	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", cnf.APIKey))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	//fmt.Println("body", string(body))
	return body, nil

}

func StoreData(data []byte) (*string, error) {
	cnf := config.GetConfig()

	//TODO - implement me
	token := cnf.APIKey
	apiUrl := cnf.RpcUrl

	//create the uploaded file
	_, fn1, err := createTheUploadedFile(data)
	if err != nil {
		return nil, err
	}
	fn := *fn1
	defer os.Remove(fn)

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

	dps := dealParamsString("default", cnf.Env)
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
		err = errors.New(string(body)) // catch auth fail!!!
		return nil, err
	}

	pinFile(fResp.Name, fResp.Hash)

	return &fResp.Hash, nil
}

func GetData(cid string) ([]byte, error) {
	cnf := config.GetConfig()
	urlLink := fmt.Sprintf("%s/ipfs/%s", cnf.RetrieveFile, cid)

	req, err := http.NewRequest(
		"GET",
		urlLink,
		nil,
	)

	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	//req.Header.Add("accept", "application/json")
	//req.Header.Add("content-type", "application/json")
	//req.Header.Add("authorization", fmt.Sprintf("Bearer %s", apikey))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//fmt.Println("body", string(body))
	return body, nil
}
