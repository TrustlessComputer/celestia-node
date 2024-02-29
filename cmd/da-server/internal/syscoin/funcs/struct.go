package funcs

const (
	METHOD_CREATE_EVM_BLOD string = "syscoincreatenevmblob"
	METHOD_GET_BLOCK_INFO  string = "getblockchaininfo"
)

type SyscoinRPCResp struct {
	jsonrpc string
	result  string
}

type SysRespErr struct {
	Result interface{} `json:"result"`
	Error  struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	Id string `json:"id"`
}
