package apis

import (
	"encoding/json"
	"net/http"
)

type RequestData struct {
	Data string `json:"data"`
}

func DecodeReqBody(r *http.Request) (RequestData, error) {
	data := RequestData{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}
