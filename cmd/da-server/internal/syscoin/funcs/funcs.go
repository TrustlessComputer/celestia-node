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

func sysRequest(syscoinMethod string, params string) ([]byte, error) {
	cnf := config.GetConfig()
	client := &http.Client{}

	if params == "" {
		params = "[]"
	} else {
		params = fmt.Sprintf(`["%s"]`, params)
	}

	id := fmt.Sprintf("bmv-%d", time.Now().UTC().UnixNano())
	requestData := fmt.Sprintf(`{"jsonrpc": "1.0", "id": "%s", "method": "%s", "params": %s}`, id, syscoinMethod, params)
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

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		d := &SysRespErr{}
		err = json.Unmarshal(body, d)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		return nil, errors.WithStack(errors.New(fmt.Sprintf("%d, %s", d.Error.Code, d.Error.Message)))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return body, nil
}

func UploadData(data []byte) (string, error) {
	var err error
	dataBlobInHex := hex.EncodeToString(data)
	body, err := sysRequest(METHOD_CREATE_EVM_BLOD, dataBlobInHex)
	if err != nil {
		return "", errors.WithStack(err)
	}

	respData := SyscoinRPCResp{}
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return respData.result, nil
}

func GetData(hash string) ([]byte, error) {
	//TODO - implement me
	return nil, nil
}

// test
func Getblockchaininfo() ([]byte, error) {
	resp, err := sysRequest(METHOD_GET_BLOCK_INFO, "")
	if err != nil {
		return nil, err
	}
	return resp, nil
}
