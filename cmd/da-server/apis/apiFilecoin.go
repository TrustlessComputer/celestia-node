package apis

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

const (
	proxySchemeFilecoin = "http"
	proxyHostFilecoin   = "127.0.0.1:22260"
)

func ApiStoreFileCoin(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	url := fmt.Sprintf("%s://%s%s", proxySchemeFilecoin, proxyHostFilecoin, r.RequestURI)
	proxyReq, err := http.NewRequest(r.Method, url, bytes.NewReader(body))
	proxyReq.Header = r.Header
	httpClient := http.Client{}
	resp, err := httpClient.Do(proxyReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	return

}

func ApiGetFileCoin(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	url := fmt.Sprintf("%s://%s%s", proxySchemeFilecoin, proxyHostFilecoin, r.RequestURI)
	proxyReq, err := http.NewRequest(r.Method, url, bytes.NewReader(body))
	proxyReq.Header = r.Header
	httpClient := http.Client{}
	resp, err := httpClient.Do(proxyReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	return
}
