package embedded

import (
	"math/big"

	"github.com/zenon-network/go-zenon/chain/nom"
	"github.com/zenon-network/go-zenon/common"
	"github.com/zenon-network/go-zenon/common/types"
	"github.com/zenon-network/go-zenon/rpc/api/embedded"
	"github.com/zenon-network/go-zenon/rpc/server"
	"github.com/zenon-network/go-zenon/vm/embedded/definition"
)

type StakeApi struct {
	client *server.Client
}

func NewStakeApi(client *server.Client) *StakeApi {
	return &StakeApi{
		client: client,
	}
}

func (sa *StakeApi) GetUncollectedReward(address types.Address) (*definition.RewardDeposit, error) {
	ans := new(definition.RewardDeposit)
	if err := sa.client.Call(ans, "embedded.stake.getUncollectedReward", address.String()); err != nil {
		return nil, err
	}
	return ans, nil
}

func (sa *StakeApi) GetFrontierRewardByPage(address types.Address, pageIndex, pageSize uint32) (*embedded.RewardHistoryList, error) {
	ans := new(embedded.RewardHistoryList)
	if err := sa.client.Call(ans, "embedded.stake.getFrontierRewardByPage", address.String(), pageIndex, pageSize); err != nil {
		return nil, err
	}
	return ans, nil
}

func (sa *StakeApi) GetEntriesByAddress(address types.Address, pageIndex, pageSize uint32) (*embedded.StakeList, error) {
	ans := new(embedded.StakeList)
	if err := sa.client.Call(ans, "embedded.stake.getEntriesByAddress", address.String(), pageIndex, pageSize); err != nil {
		return nil, err
	}
	return ans, nil
}

// Contract calls

func (sa *StakeApi) Stake(durationInSec int64, amount *big.Int) *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.StakeContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        amount,
		Data:          definition.ABIStake.PackMethodPanic(definition.StakeMethodName, durationInSec),
	}
}

func (sa *StakeApi) Cancel(id types.Hash) *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.StakeContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        common.Big0,
		Data:          definition.ABIStake.PackMethodPanic(definition.CancelStakeMethodName, id),
	}
}

func (sa *StakeApi) CollectReward() *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.StakeContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        common.Big0,
		Data:          definition.ABIStake.PackMethodPanic(definition.CollectRewardMethodName),
	}
}
