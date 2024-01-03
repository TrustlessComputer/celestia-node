package apis

type BlobHeader struct {
	Commitment       string `json:"commitment"`
	DataLength       int    `json:"dataLength"`
	BlobQuorumParams []struct {
		AdversaryThresholdPercentage int `json:"adversaryThresholdPercentage"`
		QuorumThresholdPercentage    int `json:"quorumThresholdPercentage"`
		ChunkLength                  int `json:"chunkLength"`
	} `json:"blobQuorumParams"`
}

type BatchHeader struct {
	BatchRoot               string `json:"batchRoot"`
	QuorumNumbers           string `json:"quorumNumbers"`
	QuorumSignedPercentages string `json:"quorumSignedPercentages"`
	ReferenceBlockNumber    int    `json:"referenceBlockNumber"`
}

type BatchMetadata struct {
	BatchHeader             BatchHeader `json:"batchHeader"`
	SignatoryRecordHash     string      `json:"signatoryRecordHash"`
	Fee                     string      `json:"fee"`
	ConfirmationBlockNumber int         `json:"confirmationBlockNumber"`
	BatchHeaderHash         string      `json:"batchHeaderHash"`
}

type BlobVerificationProof struct {
	BatchId        int           `json:"batchId"`
	BatchMetadata  BatchMetadata `json:"batchMetadata"`
	InclusionProof string        `json:"inclusionProof"`
	QuorumIndexes  string        `json:"quorumIndexes"`
}

type Info struct {
	BlobHeader            BlobHeader            `json:"blobHeader"`
	BlobVerificationProof BlobVerificationProof `json:"blobVerificationProof"`
}

type EigendaDataResp struct {
	Result    string `json:"result"`
	Status    string `json:"status"`
	Info      Info   `json:"info"`
	RequestId string `json:"requestId"`
	Data      string `json:"data"`
}
