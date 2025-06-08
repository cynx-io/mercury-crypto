package gopluslabs

import (
	pb "mercury/api/proto/crypto"
	"mercury/internal/pkg/logger"
	"strconv"
	"strings"
)

func (c GetTokenSecurityResponse) Data() *TokenSecurityData {
	for _, data := range c.Result {
		return &data
	}
	return nil
}

func (c GetTokenSecurityResponse) Top10HoldersPercentage() float64 {
	if len(c.Data().Holders) == 0 {
		return 0
	}

	totalPercentage := 0.0

	for i, holder := range c.Data().Holders {

		if i >= 10 {
			break
		}

		percent, err := strconv.ParseFloat(holder.Percent, 64)
		if err != nil {
			continue
		}

		totalPercentage += percent
	}

	totalPercentage *= 100.0 // Convert to percentage

	logger.Debug("total percentage: ", totalPercentage)
	return totalPercentage
}

func (c GetTokenSecurityResponse) Top10LpHoldersPercentage() float64 {
	if len(c.Data().LpHolders) == 0 {
		return 0
	}

	totalPercentage := 0.0

	for i, holder := range c.Data().LpHolders {

		if i >= 10 {
			break
		}

		percent, err := strconv.ParseFloat(holder.Percent, 64)
		if err != nil {
			continue
		}

		totalPercentage += percent
	}

	totalPercentage *= 100.0 // Convert to percentage

	logger.Debug("total percentage: ", totalPercentage)
	return totalPercentage
}

func (c GetTokenSecurityResponse) HoldersResponse() []*pb.Holder {

	holders := make([]*pb.Holder, 0, len(c.Data().Holders))

	for _, holder := range c.Data().Holders {

		percent, err := strconv.ParseFloat(holder.Percent, 64)
		if err != nil {
			continue
		}

		balance, err := strconv.ParseFloat(holder.Balance, 64)
		if err != nil {
			logger.Error("Error parsing holder balance: ", err)
			continue
		}

		// Check if the holder is locked
		isLocked := false
		if holder.IsLocked == 1 {
			isLocked = true
		}

		holders = append(holders, &pb.Holder{
			Address:    holder.Address,
			Percentage: percent,
			Balance:    balance,
			IsLocked:   isLocked,
		})
	}

	return holders
}

func (c GetTokenSecurityResponse) LpHoldersResponse() []*pb.Holder {

	holders := make([]*pb.Holder, 0, len(c.Data().Holders))

	for _, holder := range c.Data().LpHolders {

		percent, err := strconv.ParseFloat(holder.Percent, 64)
		if err != nil {
			continue
		}

		balance, err := strconv.ParseFloat(holder.Balance, 64)
		if err != nil {
			logger.Error("Error parsing holder balance: ", err)
			continue
		}

		// Check if the holder is locked
		isLocked := false
		if holder.IsLocked == 1 {
			isLocked = true
		}

		holders = append(holders, &pb.Holder{
			Address:    holder.Address,
			Percentage: percent,
			Balance:    balance,
			IsLocked:   isLocked,
		})
	}

	return holders
}

func (c GetTokenSecurityResponse) IsHoneypot() bool {
	isHoneypot, err := strconv.ParseBool(c.Data().IsHoneypot)
	if err != nil {
		return false
	}
	return isHoneypot
}

func (c GetTokenSecurityResponse) IsOpenSource() bool {
	isOpenSource, err := strconv.ParseBool(c.Data().IsOpenSource)
	if err != nil {
		return false
	}
	return isOpenSource
}

func (c GetTokenSecurityResponse) OwnershipRenounced() bool {
	return strings.ToLower(c.Data().OwnerAddress) == "0x0000000000000000000000000000000000000000"
}

func (c GetTokenSecurityResponse) IsProxy() bool {
	isProxy, err := strconv.ParseBool(c.Data().IsProxy)
	if err != nil {
		return false
	}
	return isProxy
}

func (c GetTokenSecurityResponse) IsInDex() bool {
	isInDex, err := strconv.ParseBool(c.Data().IsInDex)
	if err != nil {
		return false
	}
	return isInDex
}

func (c GetTokenSecurityResponse) IsAntiWhale() bool {
	isAntiWhale, err := strconv.ParseBool(c.Data().IsAntiWhale)
	if err != nil {
		return false
	}
	return isAntiWhale
}

func (c GetTokenSecurityResponse) IsGasAbuser() bool {
	isGasAbuser, err := strconv.ParseBool(c.Data().CannotBuy)
	if err != nil {
		return false
	}
	return isGasAbuser
}

func (c GetTokenSecurityResponse) IsBlacklisted() bool {
	isBlacklisted, err := strconv.ParseBool(c.Data().IsBlacklisted)
	if err != nil {
		return false
	}
	return isBlacklisted
}

func (c GetTokenSecurityResponse) IsWhitelisted() bool {
	isWhitelisted, err := strconv.ParseBool(c.Data().IsWhitelisted)
	if err != nil {
		return false
	}
	return isWhitelisted
}

func (c GetTokenSecurityResponse) IsTransferPausable() bool {
	isTransferPausable, err := strconv.ParseBool(c.Data().TransferPausable)
	if err != nil {
		return false
	}
	return isTransferPausable
}
