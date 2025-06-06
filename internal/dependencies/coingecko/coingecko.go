package coingecko

import (
	"mercury/internal/dependencies"
	"mercury/internal/helper"
)

type Client struct {
	config *dependencies.CoinGeckoConfig
}

func NewCoinGecko(config *dependencies.CoinGeckoConfig) *Client {
	return &Client{
		config: config,
	}
}

func (c *Client) headers() map[string]string {
	return map[string]string{
		"x-cg-demo-api-key": c.config.ApiKey,
	}
}

func (c *Client) Search(query string) (SearchResponse, error) {
	url := c.config.BaseUrl + "/api/v3/search?query=" + query
	return helper.HttpRequest[SearchResponse](url, nil, c.headers(), helper.GET)
}

func (c *Client) GetCoin(coinId string) (CoinResponse, error) {
	url := c.config.BaseUrl + "/api/v3/coins/" + coinId
	return helper.HttpRequest[CoinResponse](url, nil, c.headers(), helper.GET)
}
