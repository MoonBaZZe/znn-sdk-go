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

type PillarApi struct {
	client *server.Client
}

func NewPillarApi(client *server.Client) *PillarApi {
	return &PillarApi{
		client: client,
	}
}

func (pa *PillarApi) GetDepositedQsr(address types.Address) (*big.Int, error) {
	var ans string
	if err := pa.client.Call(ans, "embedded.pillar.getDepositedQsr", address.String()); err != nil {
		return nil, err
	}
	return common.StringToBigInt(ans), nil
}
func (pa *PillarApi) GetQsrRegistrationCost() (*big.Int, error) {
	var ans string
	if err := pa.client.Call(ans, "embedded.pillar.getQsrRegistrationCost"); err != nil {
		return nil, err
	}
	return common.StringToBigInt(ans), nil
}

func (pa *PillarApi) GetUncollectedReward(address types.Address) (*definition.RewardDeposit, error) {
	ans := new(definition.RewardDeposit)
	if err := pa.client.Call(ans, "embedded.pillar.getUncollectedReward", address.String()); err != nil {
		return nil, err
	}
	return ans, nil
}

func (pa *PillarApi) GetFrontierRewardByPage(address types.Address, pageIndex, pageSize uint32) (*embedded.RewardHistoryList, error) {
	ans := new(embedded.RewardHistoryList)
	if err := pa.client.Call(ans, "embedded.pillar.getFrontierRewardByPage", address.String(), pageIndex, pageSize); err != nil {
		return nil, err
	}
	return ans, nil
}

func (pa *PillarApi) GetAll(pageIndex, pageSize uint32) (*embedded.PillarInfoList, error) {
	ans := new(embedded.PillarInfoList)
	if err := pa.client.Call(ans, "embedded.pillar.getAll", pageIndex, pageSize); err != nil {
		return nil, err
	}
	return ans, nil
}

func (pa *PillarApi) GetByOwner(address types.Address) ([]*embedded.PillarInfo, error) {
	var ans []*embedded.PillarInfo
	if err := pa.client.Call(&ans, "embedded.pillar.getByOwner", address.String()); err != nil {
		return nil, err
	}
	return ans, nil
}

func (pa *PillarApi) GetByName(name string) (*embedded.PillarInfo, error) {
	ans := new(embedded.PillarInfo)
	if err := pa.client.Call(ans, "embedded.pillar.getByName", name); err != nil {
		return nil, err
	}
	return ans, nil
}

func (pa *PillarApi) CheckNameAvailability(name string) (*bool, error) {
	ans := new(bool)
	if err := pa.client.Call(ans, "embedded.pillar.checkNameAvailability", name); err != nil {
		return nil, err
	}
	return ans, nil
}

func (pa *PillarApi) GetDelegatedPillar(address types.Address) (*embedded.GetDelegatedPillarResponse, error) {
	ans := new(embedded.GetDelegatedPillarResponse)
	if err := pa.client.Call(ans, "embedded.pillar.getDelegatedPillar", address); err != nil {
		return nil, err
	}
	return ans, nil
}

func (pa *PillarApi) GetPillarEpochHistory(pillarName string, pageIndex, pageSize uint32) (*embedded.PillarEpochHistoryList, error) {
	ans := new(embedded.PillarEpochHistoryList)
	if err := pa.client.Call(ans, "embedded.pillar.getPillarEpochHistory", pillarName, pageIndex, pageSize); err != nil {
		return nil, err
	}
	return ans, nil
}

func (pa *PillarApi) GetPillarsHistoryByEpoch(epoch uint64, pageIndex, pageSize uint32) (*embedded.PillarEpochHistoryList, error) {
	ans := new(embedded.PillarEpochHistoryList)
	if err := pa.client.Call(ans, "embedded.pillar.getPillarsHistoryByEpoch", epoch, pageIndex, pageSize); err != nil {
		return nil, err
	}
	return ans, nil
}

// Contract calls

func (pa *PillarApi) Register(name string, producerAddress, rewardAddress types.Address, blockProducingPercentage, delegationPercentage uint8) *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.PillarContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        constants.PillarStakeAmount,
		Data: definition.ABIPillars.PackMethodPanic(
			definition.RegisterMethodName,
			name,
			producerAddress,
			rewardAddress,
			blockProducingPercentage,
			delegationPercentage,
		),
	}
}

func (pa *PillarApi) UpdatePillar(name string, producerAddress, rewardAddress types.Address, blockProducingPercentage, delegationPercentage uint8) *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.PillarContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        common.Big0,
		Data: definition.ABIPillars.PackMethodPanic(
			definition.UpdatePillarMethodName,
			name,
			producerAddress,
			rewardAddress,
			blockProducingPercentage,
			delegationPercentage,
		),
	}
}

func (pa *PillarApi) Revoke() *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.PillarContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        common.Big0,
		Data:          definition.ABIPillars.PackMethodPanic(definition.RevokeMethodName),
	}
}

func (pa *PillarApi) RegisterLegacy(name string, producerAddress, rewardAddress types.Address, blockProducingPercentage, delegationPercentage uint8, publicKey, signature string) *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.PillarContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        constants.PillarStakeAmount,
		Data: definition.ABIPillars.PackMethodPanic(
			definition.LegacyRegisterMethodName,
			name,
			producerAddress,
			rewardAddress,
			blockProducingPercentage,
			delegationPercentage,
			publicKey,
			signature,
		),
	}
}

func (pa *PillarApi) Delegate(name string) *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.PillarContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        common.Big0,
		Data:          definition.ABIPillars.PackMethodPanic(definition.DelegateMethodName, name),
	}
}

func (pa *PillarApi) Undelegate() *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.PillarContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        common.Big0,
		Data:          definition.ABIPillars.PackMethodPanic(definition.UndelegateMethodName),
	}
}

func (pa *PillarApi) DepositQsr(amount *big.Int) *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.PillarContract,
		TokenStandard: types.QsrTokenStandard,
		Amount:        amount,
		Data:          definition.ABIPillars.PackMethodPanic(definition.DepositQsrMethodName),
	}
}

func (pa *PillarApi) WithdrawQsr() *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.PillarContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        common.Big0,
		Data:          definition.ABIPillars.PackMethodPanic(definition.WithdrawQsrMethodName),
	}
}

func (pa *PillarApi) CollectReward() *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.PillarContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        common.Big0,
		Data:          definition.ABIPillars.PackMethodPanic(definition.CollectRewardMethodName),
	}
}
