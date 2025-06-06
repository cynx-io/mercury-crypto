package app

import (
	"mercury/internal/module/cryptomodule"
)

type Services struct {
	CryptoService *cryptomodule.CryptoService
}

func NewServices(repos *Repos, dependencies *Dependencies) *Services {

	return &Services{
		CryptoService: cryptomodule.NewCryptoService(dependencies.CoinGecko),
	}
}
