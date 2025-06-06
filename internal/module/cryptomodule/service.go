package cryptomodule

import (
	"context"
	"mercury/internal/dependencies/coingecko"
	"mercury/internal/model/response/responsecode"

	pb "mercury/api/proto/crypto"
)

type CryptoService struct {
	CoinGecko *coingecko.Client
}

func NewCryptoService(coinGecko *coingecko.Client) *CryptoService {
	return &CryptoService{
		CoinGecko: coinGecko,
	}
}

func (s *CryptoService) SearchCoin(ctx context.Context, query string) (*pb.SearchCoinResponse, error) {
	coins, err := s.CoinGecko.Search(query)
	if err != nil {
		return &pb.SearchCoinResponse{
			Base: &pb.BaseResponse{
				Code: responsecode.CodeCoinGeckoError.String(),
				Desc: "Database error while searching coin",
			},
		}, err
	}

	coinsResponse := make([]*pb.CoinSearchData, 0, len(coins.Coins))
	for _, coin := range coins.Coins {
		coinsResponse = append(coinsResponse, coin.Response())
	}

	return &pb.SearchCoinResponse{
		Base:  &pb.BaseResponse{Code: responsecode.CodeSuccess.String(), Desc: "Success"},
		Coins: coinsResponse,
	}, nil
}
