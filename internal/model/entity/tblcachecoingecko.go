package entity

import (
	"encoding/json"
	"mercury/internal/repository/externalapi/coingecko/coingeckomodel"
)

type TblCacheCoinGecko struct {
	EssentialEntity
	TokenId  string `gorm:"size:255;not null;uniqueIndex" json:"token_id"`
	Endpoint string `gorm:"size:255;not null;index" json:"endpoint"`
	Response string `gorm:"type:longtext" json:"response"`
}

func (data TblCacheCoinGecko) ToGetCoinResponse() (coingeckomodel.GetCoinResponse, error) {

	var resp coingeckomodel.GetCoinResponse

	if err := json.Unmarshal([]byte(data.Response), &resp); err != nil {
		return resp, err
	}

	return resp, nil
}
