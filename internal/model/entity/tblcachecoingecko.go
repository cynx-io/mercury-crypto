package entity

import (
	"encoding/json"
	"mercury/internal/repository/externalapi/coingecko/coingeckomodel"
)

type TblCacheCoinGecko struct {
	EssentialEntity
	TokenId  string `gorm:"size:255;not null;uniqueIndex;index:idx_endpoint_token_id" json:"token_id"`
	Endpoint string `gorm:"size:255;not null;index;index:idx_endpoint_token_id" json:"endpoint"`
	Response string `gorm:"type:longtext" json:"response"`
}

type PartialTblCacheCoinGeckoAddresses struct {
	TokenId    string  `gorm:"size:255;not null;uniqueIndex" json:"token_id"`
	EthAddress *string `gorm:"size:255" json:"eth_address,omitempty"`
}

func (data TblCacheCoinGecko) ToGetCoinResponse() (coingeckomodel.GetCoinResponse, error) {

	var resp coingeckomodel.GetCoinResponse

	if err := json.Unmarshal([]byte(data.Response), &resp); err != nil {
		return resp, err
	}

	return resp, nil
}
