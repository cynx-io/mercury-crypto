package helper

import (
	"strings"
	"sync"
)

const (
	codeMint              = "40c10f19" // mint(address,uint256)
	codeOwner             = "8da5cb5b" // owner()
	codePause             = "8456cb59" // pause()
	codeUnpause           = "3f4ba83a" // unpause()
	codeSetFeePercent     = "36568abe" // setFeePercent(uint256)
	codeTransferOwnership = "f2fde38b" // transferOwnership(address)
	codeRenounceOwnership = "715018a6" // renounceOwnership()
)

type EthTokenFunctions struct {
	Mint              bool `json:"mint"`
	Owner             bool `json:"owner"`
	Pause             bool `json:"pause"`
	Unpause           bool `json:"unpause"`
	SetFeePercent     bool `json:"setFeePercent"`
	TransferOwnership bool `json:"transferOwnership"`
	RenounceOwnership bool `json:"renounceOwnership"`
}

func GetTokenCodeFunctions(code string) EthTokenFunctions {
	code = strings.ToLower(code)

	var wg sync.WaitGroup
	result := &EthTokenFunctions{}

	search := func(selector string, set func()) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if strings.Contains(code, selector) {
				set()
			}
		}()
	}

	search(codeMint, func() { result.Mint = true })
	search(codeOwner, func() { result.Owner = true })
	search(codePause, func() { result.Pause = true })
	search(codeUnpause, func() { result.Unpause = true })
	search(codeSetFeePercent, func() { result.SetFeePercent = true })
	search(codeTransferOwnership, func() { result.TransferOwnership = true })
	search(codeRenounceOwnership, func() { result.RenounceOwnership = true })

	wg.Wait()
	return *result
}
