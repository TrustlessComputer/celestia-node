package apis

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/JackalLabs/jackalgo/handlers/file_io_handler"
	"github.com/JackalLabs/jackalgo/handlers/file_upload_handler"
	"github.com/JackalLabs/jackalgo/handlers/storage_handler"
	"github.com/JackalLabs/jackalgo/handlers/wallet_handler"
	"github.com/gorilla/mux"
	"jackalda/config"
	"net/http"
	"time"
)

const (
	NAMESPACE_JACKAL = "tcjackal"
	FOLDER_NAME      = "tcjackal"
	FOLDER_JACKAL    = "s/" + FOLDER_NAME
)

func ApiStoreJackal(w http.ResponseWriter, r *http.Request) {
	type RequestData struct {
		Data string `json:"data"`
	}
	data := RequestData{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Use the data
	decodedBytes, err := base64.StdEncoding.DecodeString(data.Data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	configInfo := config.GetConfig()

	wallet, err := wallet_handler.NewWalletHandler(configInfo.Seed, configInfo.RPC, configInfo.ChainId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	wallet = wallet.WithGas("500000")
	s := storage_handler.NewStorageHandler(wallet)
	_, err = s.BuyStorage(wallet.GetAddress(), 720, 2)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := map[string]interface{}{
			"error": err.Error(),
		}
		bytes, _ := json.Marshal(&resp)
		w.Write(bytes)
		return
	}

	fileIO, err := file_io_handler.NewFileIoHandler(wallet)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = fileIO.GenerateInitialDirs([]string{FOLDER_NAME})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	folder, err := fileIO.DownloadFolder(FOLDER_JACKAL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fileName := fmt.Sprintf("%v", time.Now().UnixMicro())
	file, err := file_upload_handler.TrackVirtualFile(decodedBytes, fileName, FOLDER_JACKAL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	failed, fids, cids, err := fileIO.StaggeredUploadFiles([]*file_upload_handler.FileUploadHandler{file}, folder, true, true)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = failed
	_ = cids
	_ = fids

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(fmt.Sprintf("/%s/%s", NAMESPACE_JACKAL, fileName)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	return
}

func ApiGetJackal(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileName := vars["fileName"]
	if fileName == "" {
		http.Error(w, "fileName is required", http.StatusBadRequest)
		return

	}
	configInfo := config.GetConfig()
	wallet, err := wallet_handler.NewWalletHandler(configInfo.Seed, configInfo.RPC, configInfo.ChainId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fileIO, err := file_io_handler.NewFileIoHandler(wallet.WithGas("500000"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	folder, err := fileIO.DownloadFolder(FOLDER_JACKAL)

	children := folder.GetChildFiles()
	fmt.Println(children)

	bytes, err := fileIO.DownloadRawFile(fmt.Sprintf("%s/%s", FOLDER_JACKAL, fileName))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}

	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	return
}
