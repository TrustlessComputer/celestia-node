package funcs

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/celestiaorg/celestia-node/cmd/da-server/internal/syscoin/config"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"time"
)

func credential() string {
	cnf := config.GetConfig()
	s := fmt.Sprintf("%s:%s", cnf.User, cnf.Password)
	sEnc := b64.StdEncoding.EncodeToString([]byte(s))
	return sEnc
}

func createReqHeader(req *http.Request) {
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", credential()))
}

func sysRequest(syscoinMethod string, params []interface{}) ([]byte, error) {
	cnf := config.GetConfig()
	client := &http.Client{}

	paramsStr := "[]"
	if len(params) != 0 {
		_b, err1 := json.Marshal(params)
		if err1 != nil {
			return nil, err1
		}
		paramsStr = string(_b)
	}
	id := fmt.Sprintf("bmv-%d", time.Now().UTC().UnixNano())
	requestData := fmt.Sprintf(`{"jsonrpc": "1.0", "id": "%s", "method": "%s", "params": %s}`, id, syscoinMethod, paramsStr)
	req, err := http.NewRequest(
		"POST",
		cnf.RpcURL,
		bytes.NewReader([]byte(requestData)),
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	createReqHeader(req)
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	d := &SyscoinRPCResp{}
	err = json.Unmarshal(body, d)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if d.Error != nil {
		return nil, errors.WithStack(errors.New(fmt.Sprintf("%d, %s", d.Error.Code, d.Error.Message)))
	}

	_b, err := json.Marshal(d.Result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return _b, nil
}

func UploadData(data []byte) (string, error) {
	var err error
	dataBlobInHex := hex.EncodeToString(data)
	body, err := sysRequest(METHOD_CREATE_EVM_BLOD, []interface{}{dataBlobInHex})
	if err != nil {
		return "", errors.WithStack(err)
	}

	respData := SyscoinUploadResp{}
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return respData.Versionhash, nil
}

func GetData(hash string) ([]byte, error) {
	var err error

	//https://docs.syscoin.org/docs/tech/poda#getnevmblobdata
	//example "params": ["hash_string", true]
	body, err := sysRequest(METHOD_GET_EVM_BLOD, []interface{}{hash, true})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	respData := SyscoinGetUploadedResp{}
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	_b, err := hex.DecodeString(respData.Data)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return _b, nil
}

// test
func Getblockchaininfo() ([]byte, error) {
	resp, err := sysRequest(METHOD_GET_BLOCK_INFO, []interface{}{})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
