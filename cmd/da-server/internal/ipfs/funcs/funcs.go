package funcs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/celestiaorg/celestia-node/cmd/da-server/internal/ipfs/config"
	"github.com/cockroachdb/errors"
	"io"
	"mime/multipart"
	"net/http"
)

func UploadData(fileName string, data []byte) (string, error) {
	conf := config.GetConfig()
	apikey := conf.JWT
	urlLink := fmt.Sprintf("%s/pinning/pinFileToIPFS", conf.API)

	// const pinataMetadata = JSON.stringify({
	// 	name: 'File name',
	//   });
	//   formData.append('pinataMetadata', pinataMetadata);

	//   const pinataOptions = JSON.stringify({
	// 	cidVersion: 0,
	//   })
	//   formData.append('pinataOptions', pinataOptions);

	metadata := PinataMetadata{
		Name: fileName,
	}
	metaBytes, err := json.Marshal(metadata)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fw, err := w.CreateFormFile("file", fileName)
	if err != nil {
		return "", err
	}

	if _, err = fw.Write(data); err != nil {
		return "", err
	}
	fw, err = w.CreateFormField("pinataMetadata")
	if err != nil {
		return "", err
	}
	if _, err = fw.Write(metaBytes); err != nil {
		return "", err
	}

	fw, err = w.CreateFormField("pinataOptions")
	if err != nil {
		return "", err
	}
	if _, err = fw.Write([]byte(`{"cidVersion": 0}`)); err != nil {
		return "", err
	}

	w.Close()

	req, err := http.NewRequest(
		"POST",
		urlLink,
		&b,
	)

	if err != nil {
		return "", errors.WithStack(err)
	}

	client := &http.Client{}
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", w.FormDataContentType())
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", apikey))

	// 'Content-Type': `multipart/form-data; boundary=${formData._boundary}`,

	resp, err := client.Do(req)
	if err != nil {
		return "", errors.WithStack(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.WithStack(err)
	}

	fmt.Println("body", string(body))
	var respBody Response
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		return "", errors.WithStack(err)
	}

	if respBody.Err != nil {
		return "", errors.WithStack(errors.New(fmt.Sprintf("[%s] - %s", respBody.Err.Reason, respBody.Err.Details)))
	}

	return respBody.IpfsHash, nil
}

func CheckDataExist(hash string) (bool, error) {
	conf := config.GetConfig()
	apikey := conf.JWT
	urlLink := fmt.Sprintf("%s/data/pinList?hashContains=%s", conf.API, hash)

	req, err := http.NewRequest(
		"GET",
		urlLink,
		nil,
	)

	if err != nil {
		return false, errors.WithStack(err)
	}

	client := &http.Client{}
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", apikey))

	resp, err := client.Do(req)
	if err != nil {
		return false, errors.WithStack(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, errors.WithStack(err)
	}

	fmt.Println("body", string(body))

	var respBody CheckHashExist
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		return false, errors.WithStack(err)
	}

	if respBody.Err != nil {
		return false, errors.WithStack(errors.New(fmt.Sprintf("[%s] - %s", respBody.Err.Reason, respBody.Err.Details)))
	}

	if len(respBody.Rows) == 0 {
		return false, nil
	}

	return true, nil
}

func GetData(hash string) ([]byte, error) {
	conf := config.GetConfig()
	//apikey := conf.JWT
	urlLink := fmt.Sprintf("%s/%s", conf.GetAPI, hash)

	req, err := http.NewRequest(
		"GET",
		urlLink,
		nil,
	)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	client := &http.Client{}
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	//req.Header.Add("authorization", fmt.Sprintf("Bearer %s", apikey))

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	fmt.Println("body", string(body))
	return body, nil
}
