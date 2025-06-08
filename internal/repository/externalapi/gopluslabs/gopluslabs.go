package gopluslabs

import (
	"mercury/internal/dependencies"
	"mercury/internal/helper"
)

type Client struct {
	config *dependencies.GoPlusLabsConfig
}

func NewClient(config *dependencies.GoPlusLabsConfig) *Client {
	return &Client{
		config: config,
	}
}

func (c *Client) headers() map[string]string {
	return map[string]string{}
}

func (c *Client) GetTokenSecurity(tokenAddress string) (GetTokenSecurityResponse, error) {
	url := c.config.BaseUrl + "/api/v1/token_security/1?contract_addresses=" + tokenAddress
	return helper.HttpRequest[GetTokenSecurityResponse](url, nil, c.headers(), helper.GET)
}
