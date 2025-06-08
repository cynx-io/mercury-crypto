package app

import (
	"mercury/internal/repository/database"
	"mercury/internal/repository/externalapi/alchemy"
	"mercury/internal/repository/externalapi/coingecko"
	"mercury/internal/repository/externalapi/etherscan"
	"mercury/internal/repository/externalapi/gopluslabs"
)

type Repos struct {
	TblToken *database.TblToken

	CoinGecko  *coingecko.Client
	GoPlusLabs *gopluslabs.Client
	Alchemy    *alchemy.Client
	Etherscan  *etherscan.Client
}

func NewRepos(dependencies *Dependencies) *Repos {

	tblCacheCoingecko := database.NewCacheCoinGeckoClient(dependencies.DatabaseClient.Db)

	return &Repos{
		TblToken: database.NewTblToken(dependencies.DatabaseClient.Db),

		CoinGecko:  coingecko.NewClient(&dependencies.Config.CoinGecko, tblCacheCoingecko),
		GoPlusLabs: gopluslabs.NewClient(&dependencies.Config.GoPlusLabs),
		Alchemy:    alchemy.NewClient(&dependencies.Config.Alchemy),
		Etherscan:  etherscan.NewClient(&dependencies.Config.Etherscan),
	}
}
