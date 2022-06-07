package embedded

import (
	"math/big"

	"github.com/zenon-network/go-zenon/chain/nom"
	"github.com/zenon-network/go-zenon/common"
	"github.com/zenon-network/go-zenon/common/types"
	"github.com/zenon-network/go-zenon/rpc/api"
	"github.com/zenon-network/go-zenon/rpc/api/embedded"
	"github.com/zenon-network/go-zenon/rpc/server"
	"github.com/zenon-network/go-zenon/vm/constants"
	"github.com/zenon-network/go-zenon/vm/embedded/definition"
)

type TokenApi struct {
	client *server.Client
}

func NewTokenApi(client *server.Client) *TokenApi {
	return &TokenApi{
		client: client,
	}
}

func (ta *TokenApi) GetAll(pageIndex, pageSize uint32) (*embedded.TokenList, error) {
	ans := new(embedded.TokenList)
	if err := ta.client.Call(ans, "embedded.token.getAll", pageIndex, pageSize); err != nil {
		return nil, err
	}
	return ans, nil
}

func (ta *TokenApi) GetByOwner(address types.Address, pageIndex, pageSize uint32) (*embedded.TokenList, error) {
	ans := new(embedded.TokenList)
	if err := ta.client.Call(ans, "embedded.token.getByOwner", address, pageIndex, pageSize); err != nil {
		return nil, err
	}
	return ans, nil
}

func (ta *TokenApi) GetByZts(zts types.ZenonTokenStandard) (*api.Token, error) {
	ans := new(api.Token)
	if err := ta.client.Call(ans, "embedded.token.getByZts", zts.String()); err != nil {
		return nil, err
	}
	return ans, nil
}

// Contract calls

func (ta *TokenApi) IssueToken(tokenName, tokenSymbol, tokenDomain string, totalSupply, maxSupply *big.Int, decimals uint8, isMintable, isBurnable, isUtility bool) *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.TokenContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        constants.TokenIssueAmount,
		Data: definition.ABIToken.PackMethodPanic(
			definition.IssueMethodName,
			tokenName,
			tokenSymbol,
			tokenDomain,
			totalSupply,
			maxSupply,
			decimals,
			isMintable,
			isBurnable,
			isUtility,
		),
	}
}

func (ta *TokenApi) Mint(tokenStandard types.ZenonTokenStandard, amount *big.Int, receiver types.Address) *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.TokenContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        common.Big0,
		Data: definition.ABIToken.PackMethodPanic(
			definition.MintMethodName,
			tokenStandard,
			amount,
			receiver,
		),
	}
}

func (ta *TokenApi) Burn(tokenStandard types.ZenonTokenStandard, amount *big.Int) *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.TokenContract,
		TokenStandard: tokenStandard,
		Amount:        amount,
		Data:          definition.ABIToken.PackMethodPanic(definition.BurnMethodName),
	}
}

func (ta *TokenApi) UpdateToken(tokenStandard types.ZenonTokenStandard, owner types.Address, isMintable, isBurnable bool) *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.TokenContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        common.Big0,
		Data: definition.ABIToken.PackMethodPanic(
			definition.UpdateTokenMethodName,
			tokenStandard,
			owner,
			isMintable,
			isBurnable,
		),
	}
}
