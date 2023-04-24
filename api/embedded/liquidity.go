package embedded

import (
	"github.com/zenon-network/go-zenon/chain/nom"
	"github.com/zenon-network/go-zenon/common"
	"github.com/zenon-network/go-zenon/common/types"
	"github.com/zenon-network/go-zenon/rpc/api/embedded"
	"github.com/zenon-network/go-zenon/rpc/server"
	"github.com/zenon-network/go-zenon/vm/embedded/definition"
	"math/big"
)

type LiquidityApi struct {
	client *server.Client
}

func NewLiquidityApi(client *server.Client) *LiquidityApi {
	return &LiquidityApi{
		client: client,
	}
}

func (sa *LiquidityApi) GetUncollectedReward(address types.Address) (*definition.RewardDeposit, error) {
	ans := new(definition.RewardDeposit)
	if err := sa.client.Call(ans, "embedded.liquidity.getUncollectedReward", address.String()); err != nil {
		return nil, err
	}
	return ans, nil
}

func (sa *LiquidityApi) GetFrontierRewardByPage(address types.Address, pageIndex, pageSize uint32) (*embedded.RewardHistoryList, error) {
	ans := new(embedded.RewardHistoryList)
	if err := sa.client.Call(ans, "embedded.liquidity.getFrontierRewardByPage", address.String(), pageIndex, pageSize); err != nil {
		return nil, err
	}
	return ans, nil
}

func (sa *LiquidityApi) GetLiquidityInfo() (*definition.LiquidityInfo, error) {
	ans := new(definition.LiquidityInfo)
	if err := sa.client.Call(ans, "embedded.liquidity.getLiquidityInfo"); err != nil {
		return nil, err
	}
	return ans, nil
}

func (sa *LiquidityApi) GetSecurityInfo() (*definition.SecurityInfoVariable, error) {
	ans := new(definition.SecurityInfoVariable)
	if err := sa.client.Call(ans, "embedded.liquidity.getSecurityInfo"); err != nil {
		return nil, err
	}
	return ans, nil
}

func (sa *LiquidityApi) GetLiquidityStakeEntriesByAddress(address types.Address, pageIndex, pageSize uint32) (*embedded.LiquidityStakeList, error) {
	ans := new(embedded.LiquidityStakeList)
	if err := sa.client.Call(ans, "embedded.liquidity.getLiquidityStakeEntriesByAddress", address.String(), pageIndex, pageSize); err != nil {
		return nil, err
	}
	return ans, nil
}

func (sa *LiquidityApi) GetTimeChallengesInfo() (*embedded.TimeChallengesList, error) {
	ans := new(embedded.TimeChallengesList)
	if err := sa.client.Call(ans, "embedded.liquidity.getTimeChallengesInfo"); err != nil {
		return nil, err
	}
	return ans, nil
}

// Contract calls

func (sa *LiquidityApi) SetTokenTupleMethod(tokenStandards []string, znnPercentages []uint32, qsrPercentages []uint32, minAmounts []*big.Int) *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.LiquidityContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        common.Big0,
		Data: definition.ABILiquidity.PackMethodPanic(definition.SetTokenTupleMethodName,
			tokenStandards, znnPercentages, qsrPercentages, minAmounts),
	}
}

func (sa *LiquidityApi) LiquidityStake(durationInSec int64, amount *big.Int, zts types.ZenonTokenStandard) *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.LiquidityContract,
		TokenStandard: zts,
		Amount:        amount,
		Data:          definition.ABILiquidity.PackMethodPanic(definition.LiquidityStakeMethodName, durationInSec),
	}
}

func (sa *LiquidityApi) SetIsHalted(value bool) *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.LiquidityContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        common.Big0,
		Data: definition.ABILiquidity.PackMethodPanic(
			definition.SetIsHaltedMethodName,
			value,
		),
	}
}

func (sa *LiquidityApi) CollectReward() *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.LiquidityContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        common.Big0,
		Data:          definition.ABILiquidity.PackMethodPanic(definition.CollectRewardMethodName),
	}
}

func (sa *LiquidityApi) CancelLiquidity(id types.Hash) *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.LiquidityContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        common.Big0,
		Data: definition.ABILiquidity.PackMethodPanic(definition.CancelLiquidityStakeMethodName,
			id),
	}
}

func (sa *LiquidityApi) UnlockLiquidityStakeEntries(zts types.ZenonTokenStandard) *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.LiquidityContract,
		TokenStandard: zts,
		Amount:        common.Big0,
		Data:          definition.ABILiquidity.PackMethodPanic(definition.UnlockLiquidityStakeEntriesMethodName),
	}
}

func (sa *LiquidityApi) SetAdditionalReward(znnReward *big.Int, qsrAmount *big.Int) *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.LiquidityContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        common.Big0,
		Data: definition.ABILiquidity.PackMethodPanic(definition.SetAdditionalRewardMethodName,
			znnReward, qsrAmount),
	}
}

func (sa *LiquidityApi) NominateGuardians(guardians []types.Address) *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.LiquidityContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        common.Big0,
		Data: definition.ABILiquidity.PackMethodPanic(
			definition.NominateGuardiansMethodName,
			guardians,
		),
	}
}
