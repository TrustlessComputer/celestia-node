package apis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type PinataMetadata struct {
	Name string `json:"name"`
}

type Response struct {
	IpfsHash  string `json:"IpfsHash"`
	PinSize   int    `json:"PinSize"`
	Timestamp string `json:"Timestamp"`
}

// {
// 	"IpfsHash": "QmVLwvmGehsrNEvhcCnnsw5RQNseohgEkFNN1848zNzdng",
// 	"PinSize": 32942,
// 	"Timestamp": "2023-06-16T17:24:37.998Z"
//   }

type CheckHashExist struct {
	Count int `json:"count"`
	Rows  []struct {
		ID           string    `json:"id"`
		IpfsPinHash  string    `json:"ipfs_pin_hash"`
		Size         int       `json:"size"`
		UserID       string    `json:"user_id"`
		DatePinned   time.Time `json:"date_pinned"`
		DateUnpinned any       `json:"date_unpinned"`
		Metadata     struct {
			Name      string `json:"name"`
			Keyvalues any    `json:"keyvalues"`
		} `json:"metadata"`
		Regions []struct {
			RegionID                string `json:"regionId"`
			CurrentReplicationCount int    `json:"currentReplicationCount"`
			DesiredReplicationCount int    `json:"desiredReplicationCount"`
		} `json:"regions"`
		MimeType      string `json:"mime_type"`
		NumberOfFiles int    `json:"number_of_files"`
	} `json:"rows"`
}

func CheckDataExist(apikey, hash string) (bool, error) {
	urlLink := fmt.Sprintf("https://api.pinata.cloud/data/pinList?hashContains=%s", hash)

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

	if len(respBody.Rows) == 0 {
		return false, nil
	}

	return true, nil
}

func UploadData(apikey, fileName string, data []byte) (string, error) {
	urlLink := "https://api.pinata.cloud/pinning/pinFileToIPFS"

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

	return respBody.IpfsHash, nil
}
