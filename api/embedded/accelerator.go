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

type AcceleratorApi struct {
	client *server.Client
}

func NewAcceleratorApi(client *server.Client) *AcceleratorApi {
	return &AcceleratorApi{
		client: client,
	}
}

func (aa *AcceleratorApi) GetAll(pageIndex, pageSize uint32) (*embedded.ProjectList, error) {
	ans := new(embedded.ProjectList)
	if err := aa.client.Call(ans, "embedded.accelerator.GetAll", pageIndex, pageSize); err != nil {
		return nil, err
	}
	return ans, nil
}

func (aa *AcceleratorApi) GetProjectById(id types.Hash) (*embedded.Project, error) {
	ans := new(embedded.Project)
	if err := aa.client.Call(ans, "embedded.accelerator.getProjectById", id.String()); err != nil {
		return nil, err
	}
	return ans, nil
}

func (aa *AcceleratorApi) GetPhaseById(id types.Hash) (*embedded.Phase, error) {
	ans := new(embedded.Phase)
	if err := aa.client.Call(ans, "embedded.accelerator.getPhaseById", id.String()); err != nil {
		return nil, err
	}
	return ans, nil
}

func (aa *AcceleratorApi) GetVoteBreakdown(id types.Hash) (*definition.VoteBreakdown, error) {
	ans := new(definition.VoteBreakdown)
	if err := aa.client.Call(ans, "embedded.accelerator.getVoteBreakdown", id.String()); err != nil {
		return nil, err
	}
	return ans, nil
}

func (aa *AcceleratorApi) GetPillarVotes(name string, hashes []types.Hash) ([]*definition.PillarVote, error) {
	var ans []*definition.PillarVote
	if err := aa.client.Call(&ans, "embedded.accelerator.getPillarVotes", name, hashes); err != nil {
		return nil, err
	}
	return ans, nil
}

func (aa *AcceleratorApi) CreateProject(name, description, url string, znnFundsNeeded, qsrFundsNeeded *big.Int) *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.AcceleratorContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        constants.ProjectCreationAmount,
		Data: definition.ABIAccelerator.PackMethodPanic(
			definition.CreateProjectMethodName,
			name,
			description,
			url,
			znnFundsNeeded,
			qsrFundsNeeded,
		),
	}
}

func (aa *AcceleratorApi) AddPhase(id types.Hash, name, description, url string, znnFundsNeeded, qsrFundsNeeded *big.Int) *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.AcceleratorContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        common.Big0,
		Data: definition.ABIAccelerator.PackMethodPanic(
			definition.AddPhaseMethodName,
			id,
			name,
			description,
			url,
			znnFundsNeeded,
			qsrFundsNeeded,
		),
	}
}

func (aa *AcceleratorApi) UpdatePhase(id types.Hash, name, description, url string, znnFundsNeeded, qsrFundsNeeded *big.Int) *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.AcceleratorContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        common.Big0,
		Data: definition.ABIAccelerator.PackMethodPanic(
			definition.UpdateMethodName,
			id,
			name,
			description,
			url,
			znnFundsNeeded,
			qsrFundsNeeded,
		),
	}
}

func (aa *AcceleratorApi) Donate(amount *big.Int, tokenStandard types.ZenonTokenStandard) *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.AcceleratorContract,
		TokenStandard: tokenStandard,
		Amount:        amount,
		Data:          definition.ABIAccelerator.PackMethodPanic(definition.DonateMethodName),
	}
}

func (aa *AcceleratorApi) VoteByName(id types.Hash, pillarName string, vote uint8) *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.AcceleratorContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        common.Big0,
		Data: definition.ABIAccelerator.PackMethodPanic(
			definition.VoteByNameMethodName,
			id,
			pillarName,
			vote,
		),
	}
}

func (aa *AcceleratorApi) VoteByProducerAddress(id types.Hash, vote uint8) *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.AcceleratorContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        common.Big0,
		Data: definition.ABIAccelerator.PackMethodPanic(
			definition.VoteByProdAddressMethodName,
			id,
			vote,
		),
	}
}
