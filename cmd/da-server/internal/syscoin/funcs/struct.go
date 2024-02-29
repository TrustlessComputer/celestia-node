package funcs

const (
	METHOD_CREATE_EVM_BLOD string = "syscoincreatenevmblob"
	METHOD_GET_EVM_BLOD    string = "getnevmblobdata"
	METHOD_GET_BLOCK_INFO  string = "getblockchaininfo"
)

type SyscoinRPCResp struct {
	Result interface{}     `json:"result"`
	Error  *SyscoinErrResp `json:"error"`
	Id     string          `json:"id"`
}

type SyscoinUploadResp struct {
	Versionhash string `json:"versionhash"`
}

type SyscoinGetUploadedResp struct {
	Versionhash string `json:"versionhash"`
	Mpt         int    `json:"mpt"`
	Datasize    int    `json:"datasize"`
	Data        string `json:"data"`
}

type SyscoinErrResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
