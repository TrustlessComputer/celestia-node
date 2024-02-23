package funcs

type FileCoinResp struct {
	Name string `json:"Name"`
	Hash string `json:"Hash"`
	Size string `json:"Size"`
}

type FileCoinReq struct {
	Cid      string `json:"cid"`
	FileName string `json:"fileName"`
}

type DealParams struct {
	NumCopies       int      `json:"num_copies,omitempty"`
	RepairThreshold int      `json:"repair_threshold,omitempty"`
	RenewThreshold  int      `json:"renew_threshold,omitempty"`
	Miner           []string `json:"miner,omitempty"`
	Network         string   `json:"network,omitempty"`
	AddMockData     int      `json:"add_mock_data,omitempty"`
}
