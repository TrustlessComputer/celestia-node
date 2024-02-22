package funcs

type FileCoinResp struct {
	Name string `json:"Name"`
	Hash string `json:"Hash"`
	Size string `json:"Size"`
}

type DealParams struct {
	NumCopies       int      `json:"num_copies"`
	RepairThreshold int      `json:"repair_threshold"`
	RenewThreshold  int      `json:"renew_threshold"`
	Miner           []string `json:"miner"`
	Network         string   `json:"network"`
	AddMockData     int      `json:"add_mock_data"`
}
