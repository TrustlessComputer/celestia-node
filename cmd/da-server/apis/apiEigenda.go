package apis

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
)

const (
	NAMESPACE_2   = "tceigenda"
	EIGENDA_NODE  = "disperser-goerli.eigenda.xyz:443"
	EIGENDA_PROTO = "./api/proto/disperser/disperser.proto"

	WORKING_DIR = "/root/data/eigenda"
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
	cmd := exec.Command("/root/go/bin/grpcurl", "-proto", EIGENDA_PROTO, "-d", dataCmd, EIGENDA_NODE, "disperser.Disperser/DisperseBlob")

	// set working dir:
	cmd.Dir = WORKING_DIR

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

	fmt.Println("Result:", result.Result)
	fmt.Println("RequestID:", result.RequestId)

	// get ID:
	if result.Result == "PROCESSING" {

		for i := 0; i < 20; i++ {

			fmt.Println("try times: ", i+1)

			dataCmd := fmt.Sprintf(`{"request_id": "%s"}`, result.RequestId)

			fmt.Println("dataCmd INFO", dataCmd)

			// get INFO:
			cmd := exec.Command("/root/go/bin/grpcurl", "-proto", EIGENDA_PROTO, "-d", dataCmd, EIGENDA_NODE, "disperser.Disperser/GetBlobStatus")

			// set working dir:
			cmd.Dir = WORKING_DIR

			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println("Error CombinedOutput: ", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			fmt.Println("output get detail: ", string(output))

			var resultDetail EigendaDataResp
			if err := json.Unmarshal(output, &resultDetail); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			fmt.Println("status get detail: ", resultDetail.Status)

			if resultDetail.Status == "CONFIRMED" {
				decodedBytes, err := base64.StdEncoding.DecodeString(resultDetail.Info.BlobVerificationProof.QuorumIndexes)
				if err != nil {
					fmt.Println("StdEncoding.DecodeString", err)
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				height := int(decodedBytes[0])                                                            // the quorumIndexes
				commitmentBase64 := resultDetail.Info.BlobVerificationProof.BatchMetadata.BatchHeaderHash // base64
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
				return
			}
			time.Sleep(time.Second * 30)

		}

	} else {
		http.Error(w, "NOT PROCESSING", http.StatusBadRequest)
		return
	}

	http.Error(w, "timed out", http.StatusBadRequest)
	return

}

func ApiGetEigenda(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	indexStr := vars["index"]
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	headerHashB64, err := hex.DecodeString(vars["headerHash"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	headerHashString := base64.StdEncoding.EncodeToString(headerHashB64)

	dataCmd := fmt.Sprintf(`{"batch_header_hash": "%s" ,"blob_index": "%d"}`, headerHashString, index)

	fmt.Println("dataCmd ApiGetEigenda: ", dataCmd)

	err = os.Chdir(WORKING_DIR)
	if err != nil {
		fmt.Println("Chdir(workingDir):", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get data:
	cmd := exec.Command("/root/go/bin/grpcurl", "-proto", EIGENDA_PROTO, "-d", dataCmd, EIGENDA_NODE, "disperser.Disperser/RetrieveBlob")
	fmt.Println(cmd.String())
	// set working dir:
	cmd.Dir = WORKING_DIR

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

func ApiGetEigendaBK(w http.ResponseWriter, r *http.Request) {

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

	err = os.Chdir(WORKING_DIR)
	if err != nil {
		fmt.Println("Chdir(workingDir) err:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get data:
	cmd := exec.Command("grpcurl", "-proto", EIGENDA_PROTO, "-d", dataCmd, EIGENDA_NODE, "disperser.Disperser/RetrieveBlob")

	// set working dir:
	cmd.Dir = WORKING_DIR

	// output, err := cmd.CombinedOutput()
	// if err != nil {
	// 	fmt.Println("Error CombinedOutput: ", err)
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("StderrPipe stderr:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = cmd.Start()
	if err != nil {
		fmt.Println("cmd.Start() err:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	outputErr := make([]byte, 0, 512)
	for {
		buf := make([]byte, 512)
		n, err := stderr.Read(buf)
		if n > 0 {
			outputErr = append(outputErr, buf[:n]...)
		}
		if err != nil {
			break
		}
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Println("cmd.Wait() err:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Result: " + string(outputErr))

	var result EigendaDataResp
	if err := json.Unmarshal(outputErr, &result); err != nil {
		fmt.Println("Unmarshal err:", err)
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
		fmt.Println("Write err:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	return
}
