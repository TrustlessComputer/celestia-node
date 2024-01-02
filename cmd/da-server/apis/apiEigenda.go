package apis

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os/exec"
	"strconv"
)

const (
	NAMESPACE_2   = "tceigenda"
	EIGENDA_NODE  = "disperser-goerli.eigenda.xyz:443"
	EIGENDA_PROTO = "./api/proto/disperser/disperser.proto"
)

func ApiStoreEigenda(w http.ResponseWriter, r *http.Request) {

	type RequestData struct {
		Data string `json:"data"`
	}
	data := RequestData{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dataCmd := fmt.Sprintf(`{"data": "%s", "security_params": [{"quorum_id": 0, "adversary_threshold": 25, "quorum_threshold": 50}]}`, data.Data)

	fmt.Println("dataCmd", dataCmd)

	// todo: update some filed to configs:
	cmd := exec.Command("grpcurl", "-proto", EIGENDA_PROTO, "-d", dataCmd, EIGENDA_NODE, "disperser.Disperser/DisperseBlob")

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error CombinedOutput: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var result EigendaDataResp
	if err := json.Unmarshal(output, &result); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Status:", result.Status)
	fmt.Println("RequestID:", result.RequestId)

	// get ID:
	if result.Status == "PROCESSING" {

		dataCmd := fmt.Sprintf(`{"request_id": "%s"}`, result.RequestId)

		fmt.Println("dataCmd", dataCmd)

		// get INFO:
		cmd := exec.Command("grpcurl", "-proto", EIGENDA_PROTO, "-d", dataCmd, EIGENDA_NODE, "disperser.Disperser/GetBlobStatus")

		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error CombinedOutput: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var result EigendaDataResp
		if err := json.Unmarshal(output, &result); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		decodedBytes, err := base64.StdEncoding.DecodeString(result.Info.BlobVerificationProof.QuorumIndexes)
		if err != nil {
			fmt.Println("StdEncoding.DecodeString", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		height := int(decodedBytes[0])                                                      // the quorumIndexes
		commitmentBase64 := result.Info.BlobVerificationProof.BatchMetadata.BatchHeaderHash // base64
		// convert to hex:
		decodedBytes2, err := base64.StdEncoding.DecodeString(commitmentBase64)
		if err != nil {
			fmt.Println("StdEncoding.DecodeString:", err)
			return
		}
		commitmentHex := hex.EncodeToString(decodedBytes2)

		_, err = w.Write([]byte(fmt.Sprintf("/%s/%d/%s", NAMESPACE_2, height, commitmentHex)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	return

}

func ApiGetEigenda(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	heightStr := vars["height"]
	height, err := strconv.Atoi(heightStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	commitmentHex, err := hex.DecodeString(vars["commitment"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	base64String := base64.StdEncoding.EncodeToString(commitmentHex)

	dataCmd := fmt.Sprintf(`{"batch_header_hash": "%s"}, {"blob_index": "%d"}`, base64String, height)

	fmt.Println("dataCmd: ", dataCmd)

	// get data:
	cmd := exec.Command("grpcurl", "-proto", EIGENDA_PROTO, "-d", dataCmd, EIGENDA_NODE, "disperser.Disperser/RetrieveBlob")

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error CombinedOutput: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var result EigendaDataResp
	if err := json.Unmarshal(output, &result); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// convert string to []byte
	decodedBytes, err := base64.StdEncoding.DecodeString(result.Data)
	if err != nil {
		fmt.Println("Error decoding Base64:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = w.Write(decodedBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	return
}
