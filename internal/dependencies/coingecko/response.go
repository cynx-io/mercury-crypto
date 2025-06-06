package coingecko

import "mercury/api/proto/crypto"

type SearchResponse struct {
	Coins []CoinSearchData `json:"coins"`
}

type CoinSearchData struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	ApiSymbol     string `json:"api_symbol"`
	Symbol        string `json:"symbol"`
	MarketCapRank int    `json:"market_cap_rank"`
	Thumb         string `json:"thumb"`
	Large         string `json:"large"`
}

func (c CoinSearchData) Response() *crypto.CoinSearchData {
	return &crypto.CoinSearchData{
		Id:     c.Id,
		Name:   c.Name,
		Symbol: c.Symbol,
		Thumb:  c.Thumb,
		Large:  c.Large,
	}
}

type CoinResponse struct {
}
