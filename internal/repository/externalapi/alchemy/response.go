package alchemy

type GetTokenCodeResponse struct {
	JsonRpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  string `json:"result"`
}
