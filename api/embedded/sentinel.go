package embedded

import (
	"math/big"

	"github.com/zenon-network/go-zenon/chain/nom"
	"github.com/zenon-network/go-zenon/common"
	"github.com/zenon-network/go-zenon/common/types"
	"github.com/zenon-network/go-zenon/rpc/api/embedded"
	"github.com/zenon-network/go-zenon/rpc/server"
	"github.com/zenon-network/go-zenon/vm/constants"
	"github.com/zenon-network/go-zenon/vm/embedded/definition"
)

type SentinelApi struct {
	client *server.Client
}

func NewSentinelApi(client *server.Client) *SentinelApi {
	return &SentinelApi{
		client: client,
	}
}

func (sa *SentinelApi) GetByOwner(address types.Address) (*embedded.SentinelInfo, error) {
	ans := new(embedded.SentinelInfo)
	if err := sa.client.Call(ans, "embedded.sentinel.getByOwner", address.String()); err != nil {
		return nil, err
	}
	return ans, nil
}

func (sa *SentinelApi) GetAllActive(pageIndex, pageSize uint32) (*embedded.SentinelInfoList, error) {
	ans := new(embedded.SentinelInfoList)
	if err := sa.client.Call(ans, "embedded.sentinel.getAllActive", pageIndex, pageSize); err != nil {
		return nil, err
	}
	return ans, nil
}

func (sa *SentinelApi) GetDepositedQsr(address types.Address) (*big.Int, error) {
	var ans string
	if err := sa.client.Call(&ans, "embedded.sentinel.getDepositedQsr", address); err != nil {
		return nil, err
	}
	return common.StringToBigInt(ans), nil
}

func (sa *SentinelApi) GetUncollectedReward(address types.Address) (*definition.RewardDeposit, error) {
	ans := new(definition.RewardDeposit)
	if err := sa.client.Call(ans, "embedded.sentinel.getUncollectedReward", address); err != nil {
		return nil, err
	}
	return ans, nil
}

func (sa *SentinelApi) GetFrontierRewardByPage(address types.Address, pageIndex, pageSize uint32) (*embedded.RewardHistoryList, error) {
	ans := new(embedded.RewardHistoryList)
	if err := sa.client.Call(ans, "embedded.sentinel.getFrontierRewardByPage", address, pageIndex, pageSize); err != nil {
		return nil, err
	}
	return ans, nil
}

// Contract calls

func (sa *SentinelApi) Register() *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.SentinelContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        constants.SentinelZnnRegisterAmount,
		Data:          definition.ABISentinel.PackMethodPanic(definition.RegisterSentinelMethodName),
	}
}

func (sa *SentinelApi) Revoke() *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.SentinelContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        common.Big0,
		Data:          definition.ABISentinel.PackMethodPanic(definition.RevokeSentinelMethodName),
	}
}

func (sa *SentinelApi) DepositQsr(amount *big.Int) *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.SentinelContract,
		TokenStandard: types.QsrTokenStandard,
		Amount:        amount,
		Data:          definition.ABISentinel.PackMethodPanic(definition.DepositQsrMethodName),
	}
}

func (sa *SentinelApi) WithdrawQsr() *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.SentinelContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        common.Big0,
		Data:          definition.ABISentinel.PackMethodPanic(definition.WithdrawQsrMethodName),
	}
}

func (sa *SentinelApi) CollectReward() *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.SentinelContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        common.Big0,
		Data:          definition.ABISentinel.PackMethodPanic(definition.CollectRewardMethodName),
	}
}
