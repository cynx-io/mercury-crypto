package cryptomodule

import (
	"context"
	"mercury/internal/dependencies/coingecko"
	"mercury/internal/dependencies/gopluslabs"
	"mercury/internal/model/response/responsecode"
	"mercury/internal/pkg/logger"

	pb "mercury/api/proto/crypto"
)

type CryptoService struct {
	CoinGecko  *coingecko.Client
	GoPlusLabs *gopluslabs.Client
}

func NewCryptoService(coinGecko *coingecko.Client, goPlusLabs *gopluslabs.Client) *CryptoService {
	return &CryptoService{
		CoinGecko:  coinGecko,
		GoPlusLabs: goPlusLabs,
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

func (s *CryptoService) GetCoinRisk(ctx context.Context, coinId string) (*pb.GetCoinRiskResponse, error) {
	cgCoin, err := s.CoinGecko.GetCoin(coinId)
	if err != nil {
		return &pb.GetCoinRiskResponse{
			Base: &pb.BaseResponse{
				Code: responsecode.CodeCoinGeckoError.String(),
				Desc: "Error while getting coingecko data",
			},
		}, err
	}

	tokenAddress := cgCoin.Platforms.Ethereum
	if tokenAddress == "" {
		return &pb.GetCoinRiskResponse{
			Base: &pb.BaseResponse{
				Code: responsecode.CodeNoEthereumAddress.String(),
				Desc: "Token address not found for the given coin",
			},
		}, err
	}
	logger.Debug("Token Address: ", tokenAddress)

	gplSecurity, err := s.GoPlusLabs.GetTokenSecurity(tokenAddress)
	if err != nil {
		return &pb.GetCoinRiskResponse{
			Base: &pb.BaseResponse{
				Code: responsecode.CodeGoPlusLabsError.String(),
				Desc: "Error while getting GoPlusLabs data",
			},
		}, err
	}

	logger.Debug("Preparing response for coin risk")
	return &pb.GetCoinRiskResponse{
		Base: &pb.BaseResponse{Code: responsecode.CodeSuccess.String(), Desc: "Success"},
		TokenInfo: &pb.TokenInfo{
			Name:            cgCoin.Name,
			Symbol:          cgCoin.Symbol,
			ContractAddress: cgCoin.ContractAddress,
			LogoUrl:         cgCoin.Image.Large,
			Chain:           cgCoin.AssetPlatformID,
			TotalSupply:     cgCoin.MarketData.TotalSupply,
			Decimals:        0,
		},
		RiskFlags: &pb.RiskFlags{
			HasMintFunction:    &pb.RiskFlag{},
			HasPauseFunction:   &pb.RiskFlag{},
			OwnerCanChangeFees: &pb.RiskFlag{},
			OwnershipRenounced: &pb.RiskFlag{},
			VerifiedContract:   &pb.RiskFlag{},
			IsHoneypot: &pb.RiskFlag{
				Name:        "Is Honeypot",
				Value:       gplSecurity.IsHoneypot(),
				Description: "",
			},
			LiquidityLocked:     &pb.RiskFlag{},
			RecentDeployment:    &pb.RiskFlag{},
			SuspiciousTransfers: &pb.RiskFlag{},
			Notes:               "",
		},
		HolderInfo: &pb.HolderInfo{ // TODO Issue here
			Top_10HolderPercentage: gplSecurity.Top10HoldersPercentage() * 100,
			Holders:                gplSecurity.HoldersResponse(),
		},
		SocialInfo: &pb.SocialInfo{
			Website:     cgCoin.Website(),
			Github:      cgCoin.Github(),
			Twitter:     cgCoin.Twitter(),
			Reddit:      cgCoin.Reddit(),
			Description: cgCoin.Description.En,
		},
	}, nil
}
