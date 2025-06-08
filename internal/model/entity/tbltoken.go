package entity

type TblToken struct {
	EssentialEntity
	// CoinGecko metadata
	TokenId     string `gorm:"size:255;not null;uniqueIndex" json:"token_id"`
	Name        string `gorm:"size:255;not null;index" json:"name"`
	Symbol      string `gorm:"size:255;not null;index" json:"symbol"`
	Thumb       string `gorm:"size:255;not null" json:"thumb"`
	Large       string `gorm:"size:255;not null" json:"large"`
	Description string `gorm:"size:255;type:text" json:"description"`

	// Token/Contract info (Ethereum only)
	EthAddress string `gorm:"size:255;uniqueIndex" json:"eth_address"` // Contract address on Ethereum
}
