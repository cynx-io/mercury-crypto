package cryptomodule

import (
	"context"
	"mercury/internal/helper"
	"mercury/internal/model/response/responsecode"
	"mercury/internal/pkg/logger"
	"mercury/internal/repository/database"
	"mercury/internal/repository/externalapi/alchemy"
	"mercury/internal/repository/externalapi/coingecko"
	"mercury/internal/repository/externalapi/coingecko/coingeckomodel"
	"mercury/internal/repository/externalapi/etherscan"
	"mercury/internal/repository/externalapi/gopluslabs"
	"strconv"
	"time"

	pb "mercury/api/proto/crypto"
)

type CryptoService struct {
	tblToken *database.TblToken `gorm:"primaryKey;column:token"`

	CoinGecko  *coingecko.Client
	GoPlusLabs *gopluslabs.Client
	Alchemy    *alchemy.Client
	Etherscan  *etherscan.Client
}

func NewCryptoService(
	tblToken *database.TblToken,

	coinGecko *coingecko.Client, goPlusLabs *gopluslabs.Client, alchemyClient *alchemy.Client, etherscanClient *etherscan.Client,
) *CryptoService {
	return &CryptoService{
		tblToken: tblToken,

		CoinGecko:  coinGecko,
		GoPlusLabs: goPlusLabs,
		Alchemy:    alchemyClient,
		Etherscan:  etherscanClient,
	}
}

func (s *CryptoService) SearchCoin(ctx context.Context, query string) (*pb.SearchCoinResponse, error) {

	coins, err := s.CoinGecko.Search(ctx, query)
	if err != nil {
		return &pb.SearchCoinResponse{
			Base: &pb.BaseResponse{
				Code: responsecode.CodeCoinGeckoError.String(),
				Desc: "Database error while searching coin",
			},
		}, err
	}

	tokenIds := make([]string, 0, len(coins.Coins))
	for _, coin := range coins.Coins {
		tokenIds = append(tokenIds, coin.Id)
	}

	// Bulk fetch from CoinGecko (with cache support)
	coinResponses, err := s.CoinGecko.GetCoins(ctx, tokenIds)
	if err != nil {
		return &pb.SearchCoinResponse{
			Base: &pb.BaseResponse{
				Code: responsecode.CodeCoinGeckoError.String(),
				Desc: "Error while bulk getting coin",
			},
		}, err
	}

	// Map tokenId -> GetCoinResponse
	coinRespMap := make(map[string]coingeckomodel.GetCoinResponse)
	for _, resp := range coinResponses {
		coinRespMap[resp.ID] = resp
	}

	coinsResponse := make([]*pb.CoinSearchData, 0, len(coins.Coins))
	for _, coin := range coins.Coins {
		resp, ok := coinRespMap[coin.Id]
		if !ok {
			logger.Debug("Coin not found in response map: ", coin.Id)
			continue
		}

		if len(resp.Platforms.Ethereum) == 0 {
			logger.Debug("Skipping coin without Ethereum address: ", coin.Id)
			continue
		}

		// Populate your protobuf response object
		coinsResponse = append(coinsResponse, coin.Response())
	}

	return &pb.SearchCoinResponse{
		Base:  &pb.BaseResponse{Code: responsecode.CodeSuccess.String(), Desc: "Success"},
		Coins: coinsResponse,
	}, nil
}

func (s *CryptoService) GetCoinRisk(ctx context.Context, coinId string) (*pb.GetCoinRiskResponse, error) {
	cgCoin, err := s.CoinGecko.GetCoin(ctx, coinId, false)
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

	tokenCode, err := s.Alchemy.GetTokenCode(tokenAddress)
	if err != nil {
		logger.Error("Error while getting token code: ", err)
		return &pb.GetCoinRiskResponse{
			Base: &pb.BaseResponse{
				Code: responsecode.CodeAlchemyError.String(),
				Desc: "Error while getting Alchemy data",
			},
		}, err
	}

	functions := helper.GetTokenCodeFunctions(tokenCode.Result)

	contractCreation, err := s.Etherscan.GetContractCreation(tokenAddress)
	if err != nil {
		logger.Error("Error while getting contract creation: ", err)
		return &pb.GetCoinRiskResponse{
			Base: &pb.BaseResponse{
				Code: responsecode.CodeEtherscanError.String(),
				Desc: "Error while getting Etherscan data",
			},
		}, err
	}

	logger.Debug("Getting lp holders data")
	lpHolders, top10Percentage, lockedBalance, totalSupply, liquidityLockedPercentage, err := gplSecurity.GetLpHoldersData()
	if err != nil {
		logger.Error("Error getting LP holders data: ", err)
		return &pb.GetCoinRiskResponse{
			Base: &pb.BaseResponse{
				Code: responsecode.CodeGoPlusLabsError.String(),
				Desc: "Error while getting LP holders data",
			},
		}, err
	}

	liquidityLockedRisk := &pb.RiskFlag{
		Value:    false,
		Severity: 0,
		Reason:   "",
	}

	switch {
	case liquidityLockedPercentage >= 95:
		liquidityLockedRisk.Severity = SeverityNone
		liquidityLockedRisk.Reason = "Over 95% of liquidity is locked, which is excellent."
		break
	case liquidityLockedPercentage >= 80:
		liquidityLockedRisk.Severity = SeverityMinimal
		liquidityLockedRisk.Reason = "Over 80% of liquidity is locked. This is considered very secure."
		break
	case liquidityLockedPercentage >= 60:
		liquidityLockedRisk.Severity = SeverityLow
		liquidityLockedRisk.Reason = "Between 60%–79% of liquidity is locked. Moderate risk."
		break
	case liquidityLockedPercentage >= 40:
		liquidityLockedRisk.Severity = SeverityMedium
		liquidityLockedRisk.Reason = "Less than 60% of liquidity is locked. Security may be weak."
		break
	case liquidityLockedPercentage >= 20:
		liquidityLockedRisk.Severity = SeverityHigh
		liquidityLockedRisk.Reason = "Less than 40% of liquidity is locked. High risk of rug pull."
		break
	case liquidityLockedPercentage >= 5:
		liquidityLockedRisk.Severity = SeverityExtreme
		liquidityLockedRisk.Reason = "Only 5%–19% of liquidity is locked. Very high risk."
		break
	default:
		liquidityLockedRisk.Severity = 10
		liquidityLockedRisk.Reason = "Less than 5% of liquidity is locked. Extremely risky token."
		break
	}

	// If liquidity locked is less than 60%, we flag it as a risk
	if liquidityLockedPercentage < 60 {
		liquidityLockedRisk.Value = true
	}

	isRecentDeployment := false
	isRecentDeploymentReason := "Contract creation timestamp not available"
	if len(contractCreation.Result) > 0 {
		timestamp, err := strconv.Atoi(contractCreation.Result[0].Timestamp)
		if err != nil {
			logger.Error("Error parsing timestamp: ", err)
			return &pb.GetCoinRiskResponse{
				Base: &pb.BaseResponse{
					Code: responsecode.CodeEtherscanError.String(),
					Desc: "Error parsing contract creation timestamp",
				},
			}, err
		}

		timeDiff := time.Now().Unix() - int64(timestamp)
		isRecentDeployment = timeDiff < 30*24*60*60
		isRecentDeploymentReason = "Deployed " + helper.FormatTimeDiffToYMD(timeDiff) + " ago"
	}

	riskScore := 0
	maxRiskScore := 0

	logger.Debug("Preparing Risk Flags")
	riskFlags := &pb.RiskFlags{
		Functions: &pb.RiskFlagFunctions{
			Mint: &pb.RiskFlag{
				Value:    functions.Mint,
				Severity: boolToSeverity(functions.Mint, SeverityExtreme, &riskScore, &maxRiskScore),
				Reason:   boolToReason(functions.Mint, "Token has mint function, which can inflate supply arbitrarily."),
			},
			Owner: &pb.RiskFlag{
				Value:    functions.Owner,
				Severity: boolToSeverity(functions.Owner, SeverityHigh, &riskScore, &maxRiskScore),
				Reason:   boolToReason(functions.Owner, "Owner control detected; can change token behavior."),
			},
			Pause: &pb.RiskFlag{
				Value:    functions.Pause,
				Severity: boolToSeverity(functions.Pause, SeverityMedium, &riskScore, &maxRiskScore),
				Reason:   boolToReason(functions.Pause, "Token has pause function, may halt transfers."),
			},
			Unpause: &pb.RiskFlag{
				Value:    functions.Unpause,
				Severity: SeverityNone,
				Reason:   "",
			},
			SetFeePercent: &pb.RiskFlag{
				Value:    functions.SetFeePercent,
				Severity: boolToSeverity(functions.SetFeePercent, SeverityMedium, &riskScore, &maxRiskScore),
				Reason:   boolToReason(functions.SetFeePercent, "Token can dynamically change fee percentage."),
			},
			TransferOwnership: &pb.RiskFlag{
				Value:    functions.TransferOwnership,
				Severity: boolToSeverity(functions.TransferOwnership, SeverityLow, &riskScore, &maxRiskScore),
				Reason:   boolToReason(functions.TransferOwnership, "Ownership transfer is possible."),
			},
			RenounceOwnership: &pb.RiskFlag{
				Value:    functions.RenounceOwnership,
				Severity: boolToSeverity(functions.RenounceOwnership, SeverityMinimal, &riskScore, &maxRiskScore),
				Reason:   boolToReason(functions.RenounceOwnership, "Ownership renouncement function present."),
			},
			DisableTransfer: &pb.RiskFlag{
				Value:    gplSecurity.IsTransferPausable(),
				Severity: boolToSeverity(gplSecurity.IsTransferPausable(), SeverityHigh, &riskScore, &maxRiskScore),
				Reason:   boolToReason(gplSecurity.IsTransferPausable(), "Transfers can be disabled by contract."),
			},
			Blacklist: &pb.RiskFlag{
				Value:    gplSecurity.IsBlacklisted(),
				Severity: boolToSeverity(gplSecurity.IsBlacklisted(), SeverityExtreme, &riskScore, &maxRiskScore),
				Reason:   boolToReason(gplSecurity.IsBlacklisted(), "Blacklisting logic found; can prevent some addresses from transferring."),
			},
			Whitelist: &pb.RiskFlag{
				Value:    gplSecurity.IsWhitelisted(),
				Severity: boolToSeverity(gplSecurity.IsWhitelisted(), SeverityMedium, &riskScore, &maxRiskScore),
				Reason:   boolToReason(gplSecurity.IsWhitelisted(), "Whitelisting logic found; only specific addresses may be allowed to transact."),
			},
		},
		OwnershipRenounced: &pb.RiskFlag{
			Value:    gplSecurity.OwnershipRenounced(),
			Severity: boolToSeverity(!gplSecurity.OwnershipRenounced(), SeverityExtreme, &riskScore, &maxRiskScore),
			Reason:   boolToReason(!gplSecurity.OwnershipRenounced(), "Ownership not renounced; contract still controlled."),
		},
		VerifiedContract: &pb.RiskFlag{
			Value:    gplSecurity.IsOpenSource(),
			Severity: boolToSeverity(!gplSecurity.IsOpenSource(), SeverityExtreme, &riskScore, &maxRiskScore),
			Reason:   boolToReason(!gplSecurity.IsOpenSource(), "Contract not verified or open-source."),
		},
		IsHoneypot: &pb.RiskFlag{
			Value:    gplSecurity.IsHoneypot(),
			Severity: boolToSeverity(gplSecurity.IsHoneypot(), SeverityExtreme, &riskScore, &maxRiskScore),
			Reason:   boolToReason(gplSecurity.IsHoneypot(), "Potential honeypot detected: buyers may be unable to sell."),
		},
		LiquidityLocked: liquidityLockedRisk,
		RecentDeployment: &pb.RiskFlag{
			Value:    isRecentDeployment,
			Severity: boolToSeverity(isRecentDeployment, SeverityMedium, &riskScore, &maxRiskScore),
			Reason:   isRecentDeploymentReason,
		},
		IsProxy: &pb.RiskFlag{
			Value:    gplSecurity.IsProxy(),
			Severity: boolToSeverity(gplSecurity.IsProxy(), SeverityMedium, &riskScore, &maxRiskScore),
			Reason:   boolToReason(gplSecurity.IsProxy(), "Proxy contract; logic can be upgraded."),
		},
		IsOpenSource: &pb.RiskFlag{
			Value:    gplSecurity.IsOpenSource(),
			Severity: boolToSeverity(!gplSecurity.IsOpenSource(), SeverityMedium, &riskScore, &maxRiskScore),
			Reason:   boolToReason(!gplSecurity.IsOpenSource(), "Code not publicly verifiable."),
		},
		IsAntiWhale: &pb.RiskFlag{
			Value:    gplSecurity.IsAntiWhale(),
			Severity: boolToSeverity(gplSecurity.IsAntiWhale(), SeverityMinimal, &riskScore, &maxRiskScore),
			Reason:   boolToReason(gplSecurity.IsAntiWhale(), "Anti-whale mechanism limits large transfers."),
		},
		IsGasAbuser: &pb.RiskFlag{
			Value:    gplSecurity.IsGasAbuser(),
			Severity: boolToSeverity(gplSecurity.IsGasAbuser(), SeverityLow, &riskScore, &maxRiskScore),
			Reason:   boolToReason(gplSecurity.IsGasAbuser(), "Contract is inefficient or abuses gas."),
		},
		Notes: "",
	}

	riskScoreResponse := &pb.RiskScore{
		Score:       0,
		Description: "",
	}

	logger.Debug("Calculating risk score: ", riskScore, " out of ", maxRiskScore)
	riskScoreResponse.Score = uint32(riskScore * 100 / maxRiskScore)
	if riskScoreResponse.Score > 80 {
		riskScoreResponse.Description = "High risk: Proceed with caution."
	} else if riskScoreResponse.Score > 50 {
		riskScoreResponse.Description = "Moderate risk: Review carefully."
	} else {
		riskScoreResponse.Description = "Low risk: No immediate concerns."
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
		RiskFlags: riskFlags,
		RiskScore: riskScoreResponse,
		HolderInfo: &pb.HolderInfo{
			TopTenHolderPercentage:   gplSecurity.Top10HoldersPercentage(),
			Holders:                  gplSecurity.HoldersResponse(),
			TopTenLpHolderPercentage: top10Percentage,
			LpHolders:                lpHolders,
			LpLockedPercentage:       liquidityLockedPercentage,
			LpLockedBalance:          lockedBalance,
			LpTotalBalance:           totalSupply,
		},
		SocialInfo: &pb.SocialInfo{
			Website:     cgCoin.Website(),
			Github:      cgCoin.Github(),
			Twitter:     cgCoin.Twitter(),
			Reddit:      cgCoin.Reddit(),
			Description: cgCoin.Description.En,
		},
		MarketInfo: &pb.MarketInfo{
			MarketCap:                 cgCoin.MarketData.MarketCap.Usd,
			Volume_24H:                cgCoin.MarketData.TotalVolume.Usd,
			Price:                     cgCoin.MarketData.CurrentPrice.Usd,
			PriceChange_24H:           cgCoin.MarketData.PriceChange24H,
			PriceChangePercentage_24H: cgCoin.MarketData.PriceChangePercentage24H,
			CirculatingSupply:         cgCoin.MarketData.CirculatingSupply,
			TotalSupply:               cgCoin.MarketData.TotalSupply,
			MaxSupply:                 cgCoin.MarketData.MaxSupply,
		},
	}, nil
}
