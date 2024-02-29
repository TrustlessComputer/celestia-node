package funcs

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
