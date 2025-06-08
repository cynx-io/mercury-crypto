package alchemy

import (
	"mercury/internal/dependencies"
	"mercury/internal/helper"
)

type Client struct {
	config *dependencies.AlchemyConfig
}

func NewClient(config *dependencies.AlchemyConfig) *Client {
	return &Client{
		config: config,
	}
}

func (c *Client) headers() map[string]string {
	return map[string]string{}
}

func (c *Client) GetTokenCode(tokenAddress string) (GetTokenCodeResponse, error) {
	url := c.config.BaseUrl + "/v2/" + c.config.ApiKey

	body := &GetTokenCodeRequest{
		JsonRpc: "2.0",
		Method:  "eth_getCode",
		Params:  []string{tokenAddress, "latest"},
		ID:      1,
	}

	response, err := helper.HttpRequest[GetTokenCodeResponse](url, body, c.headers(), helper.POST)
	if err != nil {
		return response, err
	}
	return response, nil
}
