package gopluslabs

type GetTokenSecurityResponse struct {
	Code    int                          `json:"code"`
	Message string                       `json:"message"`
	Result  map[string]TokenSecurityData `json:"result"`
}

type TokenSecurityData struct {
	BuyTax         string `json:"buy_tax"`
	CannotBuy      string `json:"cannot_buy"`
	CannotSellAll  string `json:"cannot_sell_all"`
	CreatorAddress string `json:"creator_address"`
	CreatorBalance string `json:"creator_balance"`
	CreatorPercent string `json:"creator_percent"`
	Dex            []struct {
		LiquidityType string `json:"liquidity_type"`
		Name          string `json:"name"`
		Liquidity     string `json:"liquidity"`
		Pair          string `json:"pair"`
		PoolManager   string `json:"poolManager,omitempty"`
	} `json:"dex"`
	HolderCount string `json:"holder_count"`
	Holders     []struct {
		Address    string `json:"address"`
		Tag        string `json:"tag"`
		IsContract int    `json:"is_contract"`
		Balance    string `json:"balance"`
		Percent    string `json:"percent"`
		IsLocked   int    `json:"is_locked"`
	} `json:"holders"`
	HoneypotWithSameCreator string `json:"honeypot_with_same_creator"`
	IsHoneypot              string `json:"is_honeypot"`
	IsInCex                 struct {
		Listed  string   `json:"listed"`
		CexList []string `json:"cex_list"`
	} `json:"is_in_cex"`
	IsInDex       string `json:"is_in_dex"`
	IsOpenSource  string `json:"is_open_source"`
	IsProxy       string `json:"is_proxy"`
	LpHolderCount string `json:"lp_holder_count"`
	LpHolders     []struct {
		Address    string `json:"address"`
		Tag        string `json:"tag"`
		Value      string `json:"value"`
		IsContract int    `json:"is_contract"`
		Balance    string `json:"balance"`
		Percent    string `json:"percent"`
		NFTList    []struct {
			Value         string `json:"value"`
			NFTID         string `json:"NFT_id"`
			Amount        string `json:"amount"`
			InEffect      string `json:"in_effect"`
			NFTPercentage string `json:"NFT_percentage"`
		} `json:"NFT_list"`
		IsLocked int `json:"is_locked"`
	} `json:"lp_holders"`
	LpTotalSupply string `json:"lp_total_supply"`
	OwnerAddress  string `json:"owner_address"`
	SellTax       string `json:"sell_tax"`
	TokenName     string `json:"token_name"`
	TokenSymbol   string `json:"token_symbol"`
	TotalSupply   string `json:"total_supply"`
	TransferTax   string `json:"transfer_tax"`
	TrustList     string `json:"trust_list"`
}
