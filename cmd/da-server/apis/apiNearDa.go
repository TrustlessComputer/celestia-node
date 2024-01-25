package apis

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	near "github.com/near/rollup-data-availability/gopkg/da-rpc"
	"net/http"
	"os"
)

const (
	// DA_KEY      = "ed25519:5rruwJXodZu6phNsApFcAm9LFxSy7nYpwnCB8vQDAvJKVgDZ424uGyXQiHQGTM3sbeBkvVXKbGxFiuswXaNRWKjv"
	// DA_CONTRACT = "8363d67bf00d1aea294285ad3fd1066378ac502d25995cf4a474326535a37201"
	// DA_ACCOUNT  = "8363d67bf00d1aea294285ad3fd1066378ac502d25995cf4a474326535a37201"

	NAMESPACE_3 = "tcnearda"
)

func GetNearDaConfig() (string, string, string) {

	_DA_CONTRACT := "8363d67bf00d1aea294285ad3fd1066378ac502d25995cf4a474326535a37201"
	_DA_ACCOUNT := "8363d67bf00d1aea294285ad3fd1066378ac502d25995cf4a474326535a37201"
	_DA_KEY := "ed25519:5rruwJXodZu6phNsApFcAm9LFxSy7nYpwnCB8vQDAvJKVgDZ424uGyXQiHQGTM3sbeBkvVXKbGxFiuswXaNRWKjv"

	env := os.Getenv("api_env")

	if env == "mainnet" {
		_DA_CONTRACT = "507cf5df56c8d98e6f5983599da44a1beacbd60f974336ed68b669769c164d44"
		_DA_ACCOUNT = "507cf5df56c8d98e6f5983599da44a1beacbd60f974336ed68b669769c164d44"
		_DA_KEY = os.Getenv("NEAR_DA_KEY")
	}

	return _DA_KEY, _DA_CONTRACT, _DA_ACCOUNT

}

func ApiTestNearDA(w http.ResponseWriter, r *http.Request) {

	DA_ACCOUNT, DA_CONTRACT, DA_KEY := GetNearDaConfig()

	config, err := near.NewConfig(DA_ACCOUNT, DA_CONTRACT, DA_KEY, 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("config", config)
	if config.Namespace.Id != 1 {
		err = errors.New("Expected namespace id to be equal")
	}
	if config.Namespace.Version != 0 {
		err = errors.New("Expected namespace version to be equal")
	}
	http.Error(w, err.Error(), http.StatusOK)

	return

}

func ConvertDataToHex(data []byte) (uint64, string, error) {
	frameRef := near.FrameRef{}
	err := frameRef.UnmarshalBinary(data)
	if err != nil {
		return 0, "", err
	}
	fmt.Println("frameRef.TxId", frameRef.TxId, "frameRef.TxCommitment", frameRef.TxCommitment)

	fmt.Println("frameRef.TxId.String()", string(frameRef.TxId), "frameRef.TxCommitment.String()", string(frameRef.TxCommitment))

	commitmentHex := hex.EncodeToString(frameRef.TxCommitment)

	height := binary.BigEndian.Uint64(frameRef.TxId)

	return height, commitmentHex, nil
}

func ApiStoreNearDA(w http.ResponseWriter, r *http.Request) {

	type RequestData struct {
		Data string `json:"data"`
	}
	data := RequestData{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	DA_ACCOUNT, DA_CONTRACT, DA_KEY := GetNearDaConfig()

	config, err := near.NewConfig(DA_ACCOUNT, DA_CONTRACT, DA_KEY, 1)
	if err != nil {
		fmt.Println("NewConfig err:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var candidateHex string = "0xfF00000000000000000000000000000000000000"

	// dataBase64 := "AILhYPddiq2lIgaPiWDGH4gAAAAABDR42trhz/DDZ0FpR5WUq4O43rK3H2vblkz91f83pJaXpSFMILzttVTkZv9mPk7WBWxtbz59VVqjuTKofvq2kPlf7Wvs7xZpNv8sqP8Za3hmmnxL6tTv6gcgBq5icFCNPpB650ShRNcC/o3XAr7teppjMe1Il2jnxNuXs2cTbaAh1EAls5DLpbtXHSvTTNpy+FyfckCo8QTjGtdJYbWsmmsZO74QbaA11MA0brPry47/V2dIU50xW+jtD9MF0p3RtzR0CxfWlxltnZrczMfJtqCVUcD0a7lmev2HM5vuHX8lwrhs14Wa9GOP73FOUazUCpQFGegKNfDy57dK4XM1PeduWnf/O7tQ1571Rlcfv3jqvYl5kswHqdeFRBvoDzXQz+ey+4bIytzV9VVPdBtnvuy/cv7JCX0Bj00eEsw1THmmzXyc7AuuXr/Wc2HTnSCD8AeFKhN8/v2ZE6xQdEKN/RtXiYK3zmYekIGRUAN/sWpV1QXvSdXkSxB9V+L0SlWv9Na6jQ4GEYcbVy2bd0ONaAOToQauNZjQvvdD0bUFiVyOE7TXc8j9O3xxzYY/e/39Ph9o0din18zHybHg0jeWIv7wDOHH1Rnvry5n8Veq0BMvSKzcs/v9YY9vLCsSQAbmQg28kWx0OqUvbTX/VR6D6zIn6h3/rJ+zbKLUkluCl26e/OI4h2gDy6EGXmvLnzU1/tyqeXM85z86e9rroH8m5yGWlDCX6TOWS15gOEq0gY1QAyPkK/tVgpb+2vkzICTlx6oXcv9afWxVl7YnCTNV/s97t72Zj5NzwWlzgRMXH3xL3TmpiTF5k8f6yCQOTjU+A72W0mdGwZIFn0AGdkMNPDUrRM7aMO/dxjmPd4eX/ooxUFn927nq6PX3C2W3LM2TMiTawKlQA9PnvzefkHzpfZqda1zr3NUxr1J2dzVUiE+cmaw86W8Z74FmPk6uBZ8mcaU6Xj4f/q9ypcXDmcubPuj+PNYzXeXr1VWfnJxCY9VABs6HGsgyIya+I0DbbWGZNNfPHwrPOYyOBy49O6X0qtnz59fvXXlOtIEroQZ2i2YxdVjXJib1Vbb3F7pqR03w+C/Rm+xbVjWjjNPpdBzRBm6GGqiQGDp/LitnpNYOv5A+J/myS2t/8Sxyzvh6uXLOHYbbAZLNfJzcC3oWbVO/8PlSOuMDl6SeGY4LtusdPusmek/j2WoFK5nEJ/YgA/dCDXzw8JLMTKeVvzJl7ygZv3AS7Uo3l70k+8A92DU2Nd7vKjPRBh6HGrhwikgcpyd7fUj732PB/lerb/OdO267eV2nem33pJqTpjua+Th5FvA9cFmTs7Y3s030xhu2Jna/IzFWb74cVZ3131d28vRzRnNABl6EGrjALHh5QYhC9YptPkmrn9pt/nKtOEzCpyBcNWG9Smz2dCeiDbx9ABAAAP//YyLw+gE="

	dataBase64 := data.Data

	// bytes := make([]byte, 64)
	// copy(bytes, []byte(dataBase64))
	decodedBytes, err := base64.StdEncoding.DecodeString(dataBase64)
	if err != nil {
		fmt.Println("Error decoding Base64:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("decodedBytes: ", decodedBytes)

	result, err := config.Submit(candidateHex, []byte(decodedBytes))

	if err != nil {
		fmt.Println("submit err:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resultHex := hex.EncodeToString(result)

	fmt.Println("resultHex ", resultHex)

	height, commitmentHex, err := ConvertDataToHex(result)
	if err != nil {
		fmt.Println("ConvertDataToHex err:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("height: ", height)
	fmt.Println("commitmentHex: ", commitmentHex)

	_, err = w.Write([]byte(fmt.Sprintf("/%s/%s", NAMESPACE_3, resultHex)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	return

}

func ApiGetNearDA(w http.ResponseWriter, r *http.Request) {

	// dataHex := "5dc15471df1cc3fe66c79fd183076b1c0a255ec89026e021cab1fb591b8641e5847927a9251f1063703f640348cba026b214eb28dbd85837409c41aedccc7238"

	vars := mux.Vars(r)

	dataHex := vars["dataHex"]

	// convert to []byte:
	resultByte, err := hex.DecodeString(dataHex)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	DA_ACCOUNT, DA_CONTRACT, DA_KEY := GetNearDaConfig()
	config, err := near.NewConfig(DA_ACCOUNT, DA_CONTRACT, DA_KEY, 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	blob, err := config.Get(resultByte, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// fmt.Println("blob byte: ", blob)
	// fmt.Println("blob string: ", string(blob))

	//resultHex := hex.EncodeToString(blob)
	//
	//fmt.Println("resultHex: ", resultHex)
	//
	//// convert string to []byte
	//decodedBytes, err := base64.StdEncoding.DecodeString(resultHex)
	//if err != nil {
	//	fmt.Println("Error decoding Base64:", err)
	//	http.Error(w, err.Error(), http.StatusBadRequest)
	//	return
	//}
	//
	//fmt.Println("decodedBytes: ", decodedBytes)

	_, err = w.Write(blob)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	return
}
