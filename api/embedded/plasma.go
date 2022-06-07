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

type PlasmaApi struct {
	client *server.Client
}

func NewPlasmaApi(client *server.Client) *PlasmaApi {
	return &PlasmaApi{
		client: client,
	}
}

func (pa *PlasmaApi) Get(address types.Address) (*embedded.PlasmaInfo, error) {
	ans := new(embedded.PlasmaInfo)
	if err := pa.client.Call(ans, "embedded.plasma.get", address.String()); err != nil {
		return nil, err
	}
	return ans, nil
}

func (pa *PlasmaApi) GetEntriesByAddress(address types.Address, pageIndex, pageSize uint32) (*embedded.FusionEntryList, error) {
	ans := new(embedded.FusionEntryList)
	if err := pa.client.Call(ans, "embedded.plasma.getEntriesByAddress", address.String(), pageIndex, pageSize); err != nil {
		return nil, err
	}
	return ans, nil
}

func (pa *PlasmaApi) GetRequiredPoWForAccountBlock(param embedded.GetRequiredParam) (*embedded.GetRequiredResult, error) {
	ans := new(embedded.GetRequiredResult)
	if err := pa.client.Call(ans, "embedded.plasma.getRequiredPoWForAccountBlock", param); err != nil {
		return nil, err
	}
	return ans, nil
}

// Contract calls

func (pa *PlasmaApi) Fuse(address types.Address, amount *big.Int) *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.PlasmaContract,
		TokenStandard: types.QsrTokenStandard,
		Amount:        amount,
		Data:          definition.ABIPlasma.PackMethodPanic(definition.FuseMethodName, address),
	}
}

func (pa *PlasmaApi) Cancel(id types.Hash) *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.PlasmaContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        common.Big0,
		Data:          definition.ABIPlasma.PackMethodPanic(definition.CancelFuseMethodName, id),
	}
}
