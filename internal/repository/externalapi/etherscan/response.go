package etherscan

type GetContractCreationResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  []struct {
		ContractAddress  string `json:"contractAddress"`
		ContractCreator  string `json:"contractCreator"`
		TxHash           string `json:"txHash"`
		BlockNumber      string `json:"blockNumber"`
		Timestamp        string `json:"timestamp"`
		ContractFactory  string `json:"contractFactory"`
		CreationBytecode string `json:"creationBytecode"`
	} `json:"result"`
}
