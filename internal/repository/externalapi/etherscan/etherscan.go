package etherscan

import (
	"mercury/internal/dependencies"
	"mercury/internal/helper"
)

type Client struct {
	config *dependencies.EtherscanConfig
}

func NewClient(config *dependencies.EtherscanConfig) *Client {
	return &Client{
		config: config,
	}
}

func (c *Client) headers() map[string]string {
	return map[string]string{}
}

func (c *Client) GetContractCreation(tokenAddress string) (GetContractCreationResponse, error) {
	url := c.config.BaseUrl + "/api?module=contract&action=getcontractcreation&contractaddresses=" + tokenAddress + "&apikey=" + c.config.ApiKey

	response, err := helper.HttpRequest[GetContractCreationResponse](url, nil, c.headers(), helper.GET)
	if err != nil {
		return response, err
	}
	return response, nil
}
