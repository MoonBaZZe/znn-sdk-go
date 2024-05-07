package embedded

import (
	"github.com/zenon-network/go-zenon/chain/nom"
	"github.com/zenon-network/go-zenon/common"
	"github.com/zenon-network/go-zenon/common/types"
	"github.com/zenon-network/go-zenon/rpc/api/embedded"
	"github.com/zenon-network/go-zenon/rpc/server"
	"github.com/zenon-network/go-zenon/vm/embedded/definition"
)

type MergeMiningApi struct {
	client *server.Client
}

func NewMergeMiningApi(client *server.Client) *MergeMiningApi {
	return &MergeMiningApi{
		client: client,
	}
}

func (mma *MergeMiningApi) GetMergeMiningInfo() (*definition.MergeMiningInfoVariable, error) {
	ans := new(definition.MergeMiningInfoVariable)
	if err := mma.client.Call(ans, "embedded.merge_mining.getMergeMiningInfo"); err != nil {
		return nil, err
	}
	return ans, nil
}

func (mma *MergeMiningApi) GetHeaderChainInfo() (*definition.HeaderChainInfoVariable, error) {
	ans := new(definition.HeaderChainInfoVariable)
	if err := mma.client.Call(ans, "embedded.merge_mining.getHeaderChainInfo"); err != nil {
		return nil, err
	}
	return ans, nil
}

func (mma *MergeMiningApi) GetShareChainInfo(id uint8) (*definition.ShareChainInfoVariable, error) {
	ans := new(definition.ShareChainInfoVariable)
	if err := mma.client.Call(ans, "embedded.merge_mining.getShareChainInfo", id); err != nil {
		return nil, err
	}
	return ans, nil
}

func (mma *MergeMiningApi) GetBlockHeader(hash types.Hash) (*definition.BlockHeaderVariable, error) {
	ans := new(definition.BlockHeaderVariable)
	if err := mma.client.Call(ans, "embedded.merge_mining.getShareChainInfo", hash); err != nil {
		return nil, err
	}
	return ans, nil
}

func (mma *MergeMiningApi) GetTimeChallengesInfo() (*embedded.TimeChallengesList, error) {
	ans := new(embedded.TimeChallengesList)
	if err := mma.client.Call(ans, "embedded.merge_mining.getTimeChallengesInfo"); err != nil {
		return nil, err
	}
	return ans, nil
}

func (mma *MergeMiningApi) GetSecurityInfo() (*definition.SecurityInfoVariable, error) {
	ans := new(definition.SecurityInfoVariable)
	if err := mma.client.Call(ans, "embedded.merge_mining.getSecurityInfo"); err != nil {
		return nil, err
	}
	return ans, nil
}

func (mma *MergeMiningApi) SetInitialBitcoinBlock(version int32, prevBlock, merkleRoot types.Hash, timestamp, bits, nonce uint32) *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.MergeMiningContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        common.Big0,
		Data: definition.ABIBridge.PackMethodPanic(
			definition.SetInitialBitcoinBlockHeaderMethodName,
			version,
			prevBlock,
			merkleRoot,
			timestamp,
			bits,
			nonce,
		),
	}
}

func (mma *MergeMiningApi) AddBitcoinBlockHeader(version int32, prevBlock, merkleRoot types.Hash, timestamp, bits, nonce uint32) *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.MergeMiningContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        common.Big0,
		Data: definition.ABIBridge.PackMethodPanic(
			definition.AddBitcoinBlockHeaderMethodName,
			version,
			prevBlock,
			merkleRoot,
			timestamp,
			bits,
			nonce,
		),
	}
}

func (mma *MergeMiningApi) AddShare(shareChainId uint8, witness bool, version int32, prevBlock, merkleRoot types.Hash, timestamp, bits, nonce uint32,
	proof, prooff, proofff, prooffff, additionalData types.Hash) *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.MergeMiningContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        common.Big0,
		Data: definition.ABIBridge.PackMethodPanic(
			definition.AddShareMethodName,
			shareChainId,
			witness,
			version,
			prevBlock,
			merkleRoot,
			timestamp,
			bits,
			nonce,
			proof, prooff, proofff, prooffff,
			additionalData,
		),
	}
}

func (mma *MergeMiningApi) SetShareChain(id uint8, bits, rewardMultiplier uint32) *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.MergeMiningContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        common.Big0,
		Data: definition.ABIBridge.PackMethodPanic(
			definition.SetShareChainMethodName,
			id,
			bits,
			rewardMultiplier,
		),
	}
}

func (mma *MergeMiningApi) Emergency() *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.MergeMiningContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        common.Big0,
		Data: definition.ABIBridge.PackMethodPanic(
			definition.EmergencyMethodName,
		),
	}
}

func (mma *MergeMiningApi) ChangeTssECDSAPubKey(pubKey, signature, newSignature string) *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.MergeMiningContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        common.Big0,
		Data: definition.ABIBridge.PackMethodPanic(
			definition.ChangeTssECDSAPubKeyMethodName,
			pubKey,
			signature,
			newSignature,
		),
	}
}

func (mma *MergeMiningApi) ChangeAdministrator(administrator types.Address) *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.MergeMiningContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        common.Big0,
		Data: definition.ABIBridge.PackMethodPanic(
			definition.ChangeAdministratorMethodName,
			administrator,
		),
	}
}

func (mma *MergeMiningApi) NominateGuardians(guardians []types.Address) *nom.AccountBlock {
	return &nom.AccountBlock{
		BlockType:     nom.BlockTypeUserSend,
		ToAddress:     types.MergeMiningContract,
		TokenStandard: types.ZnnTokenStandard,
		Amount:        common.Big0,
		Data: definition.ABIBridge.PackMethodPanic(
			definition.NominateGuardiansMethodName,
			guardians,
		),
	}
}
