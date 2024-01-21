package apis

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	near "github.com/near/rollup-data-availability/gopkg/da-rpc"
	"net/http"
	// "time"
)

const (
	DA_KEY      = "ed25519:5rruwJXodZu6phNsApFcAm9LFxSy7nYpwnCB8vQDAvJKVgDZ424uGyXQiHQGTM3sbeBkvVXKbGxFiuswXaNRWKjv"
	DA_CONTRACT = "8363d67bf00d1aea294285ad3fd1066378ac502d25995cf4a474326535a37201"
	DA_ACCOUNT  = "8363d67bf00d1aea294285ad3fd1066378ac502d25995cf4a474326535a37201"
)

func ApiTestNearDA(w http.ResponseWriter, r *http.Request) {

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

	config, err := near.NewConfig(DA_ACCOUNT, DA_CONTRACT, DA_KEY, 1)
	if err != nil {
		fmt.Println("NewConfig err:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var candidateHex string = "0xfF00000000000000000000000000000000000000"

	// var data []byte = []byte("elvis")

	dataBase64 := "AILhYPddiq2lIgaPiWDGH4gAAAAABDR42trhz/DDZ0FpR5WUq4O43rK3H2vblkz91f83pJaXpSFMILzttVTkZv9mPk7WBWxtbz59VVqjuTKofvq2kPlf7Wvs7xZpNv8sqP8Za3hmmnxL6tTv6gcgBq5icFCNPpB650ShRNcC/o3XAr7teppjMe1Il2jnxNuXs2cTbaAh1EAls5DLpbtXHSvTTNpy+FyfckCo8QTjGtdJYbWsmmsZO74QbaA11MA0brPry47/V2dIU50xW+jtD9MF0p3RtzR0CxfWlxltnZrczMfJtqCVUcD0a7lmev2HM5vuHX8lwrhs14Wa9GOP73FOUazUCpQFGegKNfDy57dK4XM1PeduWnf/O7tQ1571Rlcfv3jqvYl5kswHqdeFRBvoDzXQz+ey+4bIytzV9VVPdBtnvuy/cv7JCX0Bj00eEsw1THmmzXyc7AuuXr/Wc2HTnSCD8AeFKhN8/v2ZE6xQdEKN/RtXiYK3zmYekIGRUAN/sWpV1QXvSdXkSxB9V+L0SlWv9Na6jQ4GEYcbVy2bd0ONaAOToQauNZjQvvdD0bUFiVyOE7TXc8j9O3xxzYY/e/39Ph9o0din18zHybHg0jeWIv7wDOHH1Rnvry5n8Veq0BMvSKzcs/v9YY9vLCsSQAbmQg28kWx0OqUvbTX/VR6D6zIn6h3/rJ+zbKLUkluCl26e/OI4h2gDy6EGXmvLnzU1/tyqeXM85z86e9rroH8m5yGWlDCX6TOWS15gOEq0gY1QAyPkK/tVgpb+2vkzICTlx6oXcv9afWxVl7YnCTNV/s97t72Zj5NzwWlzgRMXH3xL3TmpiTF5k8f6yCQOTjU+A72W0mdGwZIFn0AGdkMNPDUrRM7aMO/dxjmPd4eX/ooxUFn927nq6PX3C2W3LM2TMiTawKlQA9PnvzefkHzpfZqda1zr3NUxr1J2dzVUiE+cmaw86W8Z74FmPk6uBZ8mcaU6Xj4f/q9ypcXDmcubPuj+PNYzXeXr1VWfnJxCY9VABs6HGsgyIya+I0DbbWGZNNfPHwrPOYyOBy49O6X0qtnz59fvXXlOtIEroQZ2i2YxdVjXJib1Vbb3F7pqR03w+C/Rm+xbVjWjjNPpdBzRBm6GGqiQGDp/LitnpNYOv5A+J/myS2t/8Sxyzvh6uXLOHYbbAZLNfJzcC3oWbVO/8PlSOuMDl6SeGY4LtusdPusmek/j2WoFK5nEJ/YgA/dCDXzw8JLMTKeVvzJl7ygZv3AS7Uo3l70k+8A92DU2Nd7vKjPRBh6HGrhwikgcpyd7fUj732PB/lerb/OdO267eV2nem33pJqTpjua+Th5FvA9cFmTs7Y3s030xhu2Jna/IzFWb74cVZ3131d28vRzRnNABl6EGrjALHh5QYhC9YptPkmrn9pt/nKtOEzCpyBcNWG9Smz2dCeiDbx9ABAAAP//YyLw+gE="

	bytes := make([]byte, 64)
	copy(bytes, []byte(dataBase64))

	result, err := config.Submit(candidateHex, bytes)

	if err != nil {
		fmt.Println("submit err:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	height, commitmentHex, err := ConvertDataToHex(result)
	if err != nil {
		fmt.Println("ConvertDataToHex err:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("height: ", height)
	fmt.Println("commitmentHex: ", commitmentHex)

	// if string(frameRef.TxId) != "11111111111111111111111111111111" {
	// 	err = errors.New("Expected id to be equal")
	// 	fmt.Println(frameRef.TxId, " err:", err)
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return

	// }
	// if string(frameRef.TxCommitment) != "22222222222222222222222222222222" {
	// 	err = errors.New("Expected commitment to be equal")
	// 	fmt.Println(frameRef.TxCommitment, " err:", err)
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// time.Sleep(30 * time.Second)

	return

}

func ApiGetNearDA(w http.ResponseWriter, r *http.Request) {

	commitmentHex := "ed8e75db33506660bbbb1e7c98b9e9708b02587314b4b7a171304b90fadc49dc"
	height := "5814586574713041586"

	commitmentHashB64, err := hex.DecodeString(commitmentHex)
	if err != nil {
		fmt.Println("DecodeString err:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	commitmentStr := base64.StdEncoding.EncodeToString(commitmentHashB64)

	id := make([]byte, 32)
	copy(id, []byte(height))

	commitment := make([]byte, 32)
	copy(commitment, []byte(commitmentStr))
	frameRef := near.FrameRef{
		TxId:         id,
		TxCommitment: commitment,
	}
	binary, err := frameRef.MarshalBinary()
	println("binary, id, commitment", binary, id, commitment)
	if err != nil {
		fmt.Println("MarshalBinary err:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	config, err := near.NewConfig(DA_ACCOUNT, DA_CONTRACT, DA_KEY, 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	blob, err := config.Get(binary, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("blob byte: ", blob)
	fmt.Println("blob string: ", string(blob))

	dataHex := hex.EncodeToString(blob)

	fmt.Println("dataHex: ", dataHex)

	return
}
