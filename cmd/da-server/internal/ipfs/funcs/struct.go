package funcs

import "time"

type PinataMetadata struct {
	Name string `json:"name"`
}

type Response struct {
	IpfsHash  string         `json:"IpfsHash"`
	PinSize   int            `json:"PinSize"`
	Timestamp string         `json:"Timestamp"`
	Err       *ResponseError `json:"error"`
}

type ResponseError struct {
	Reason  string `json:"reason"`
	Details string `json:"details"`
}

type CheckHashExist struct {
	Count int `json:"count"`
	Rows  []struct {
		ID           string    `json:"id"`
		IpfsPinHash  string    `json:"ipfs_pin_hash"`
		Size         int       `json:"size"`
		UserID       string    `json:"user_id"`
		DatePinned   time.Time `json:"date_pinned"`
		DateUnpinned any       `json:"date_unpinned"`
		Metadata     struct {
			Name      string `json:"name"`
			Keyvalues any    `json:"keyvalues"`
		} `json:"metadata"`
		Regions []struct {
			RegionID                string `json:"regionId"`
			CurrentReplicationCount int    `json:"currentReplicationCount"`
			DesiredReplicationCount int    `json:"desiredReplicationCount"`
		} `json:"regions"`
		MimeType      string `json:"mime_type"`
		NumberOfFiles int    `json:"number_of_files"`
	} `json:"rows"`
	Err *ResponseError `json:"error"`
}
