package gopluslabs

import (
	pb "mercury/api/proto/crypto"
	"strconv"
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

	return totalPercentage
}

func (c GetTokenSecurityResponse) HoldersResponse() []*pb.Holder {

	holders := make([]*pb.Holder, 0, len(c.Data().Holders))

	for _, holder := range c.Data().Holders {

		percent, err := strconv.ParseFloat(holder.Percent, 64)
		if err != nil {
			continue
		}

		holders = append(holders, &pb.Holder{
			Address:    holder.Address,
			Percentage: percent,
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
