package coingeckomodel

import "mercury/api/proto/crypto"

func (c CoinSearchData) Response() *crypto.CoinSearchData {
	return &crypto.CoinSearchData{
		Id:     c.Id,
		Name:   c.Name,
		Symbol: c.Symbol,
		Thumb:  c.Thumb,
		Large:  c.Large,
	}
}

func (c GetCoinResponse) Website() string {
	if len(c.Links.Homepage) > 0 {
		return c.Links.Homepage[0]
	}
	return ""
}

func (c GetCoinResponse) Github() string {
	if len(c.Links.ReposURL.Github) > 0 {
		return c.Links.ReposURL.Github[0]
	}
	return ""
}

func (c GetCoinResponse) Twitter() string {
	return c.Links.TwitterScreenName
}

func (c GetCoinResponse) Reddit() string {
	return c.Links.SubredditURL
}
