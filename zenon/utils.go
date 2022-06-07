package zenon

import (
	"znn-sdk-go/rpc_client"

	"github.com/inconshreveable/log15"
	"github.com/zenon-network/go-zenon/chain/nom"
	"github.com/zenon-network/go-zenon/common"
	"github.com/zenon-network/go-zenon/common/crypto"
	"github.com/zenon-network/go-zenon/common/types"
	"github.com/zenon-network/go-zenon/pow"
	"github.com/zenon-network/go-zenon/rpc/api"
	"github.com/zenon-network/go-zenon/rpc/api/embedded"
	"github.com/zenon-network/go-zenon/wallet"
)

var (
	RPCLogger    = log15.New("module", "go-sdk-rpc")
	WalletLogger = log15.New("module", "go-sdk-wallet")
	CommonLogger = log15.New("module", "go-sdk-common")
)

func autofillTransactionParameters(client *rpc_client.RpcClient, transaction *nom.AccountBlock) error {
	transaction.Version = transactionVersion
	transaction.ChainIdentifier = netId

	frontierAccountBlock, err := client.LedgerApi.GetFrontierAccountBlock(transaction.Address)

	if err != nil {
		return err
	}
	height := uint64(1)
	previousHash := types.ZeroHash
	if frontierAccountBlock != nil {
		height = frontierAccountBlock.Height + 1
		previousHash = frontierAccountBlock.Hash
	}
	transaction.Height = height
	transaction.PreviousHash = previousHash

	frontierMomentum, err := client.LedgerApi.GetFrontierMomentum()
	if err != nil {
		return err
	}
	momentumAcknowledged := types.HashHeight{
		Hash:   frontierMomentum.Hash,
		Height: frontierMomentum.Height,
	}
	transaction.MomentumAcknowledged = momentumAcknowledged

	return nil
}

func checkAndSetFields(client *rpc_client.RpcClient, transaction *nom.AccountBlock, currentKeyPair wallet.KeyPair) error {
	transaction.Address = currentKeyPair.Address
	transaction.PublicKey = currentKeyPair.Public

	if err := autofillTransactionParameters(client, transaction); err != nil {
		return err
	}

	if transaction.BlockType == nom.BlockTypeUserReceive {
		if transaction.FromBlockHash == types.ZeroHash {
			return ErrZeroFromHash
		}
		sendBlock, err := client.LedgerApi.GetAccountBlockByHash(transaction.FromBlockHash)
		if err != nil {
			return err
		} else if sendBlock == nil {
			return ErrNonExistentAccountBlock
		}

		if sendBlock.ToAddress != transaction.Address {
			return ErrDifferentReceiver
		}

		if len(transaction.Data) > 0 {
			return ErrContainsData
		}
	}

	nonce, err := transaction.Nonce.MarshalText()
	if err != nil {
		return err
	}
	if transaction.Difficulty > 0 && len(nonce) == 0 {
		return ErrNoNonce
	}

	return nil
}

func setDifficulty(client *rpc_client.RpcClient, transaction *nom.AccountBlock) error {
	powParam := embedded.GetRequiredParam{
		SelfAddr:  transaction.Address,
		BlockType: transaction.BlockType,
		ToAddr:    &transaction.ToAddress,
		Data:      transaction.Data,
	}
	response, err := client.PlasmaApi.GetRequiredPoWForAccountBlock(powParam)
	if err != nil {
		return err
	}
	if response.RequiredDifficulty.Cmp(common.Big0) > 0 {
		transaction.FusedPlasma = response.AvailablePlasma
		transaction.Difficulty = response.RequiredDifficulty.Uint64()

		CommonLogger.Info("Generating Plasma, please wait\n")

		hashBytes := crypto.Hash(transaction.Address.Bytes(), transaction.PreviousHash.Bytes())
		powData, err := types.BytesToHash(hashBytes)
		if err != nil {
			return err
		}

		nonceFound := pow.GetPoWNonce(response.RequiredDifficulty, powData)

		transaction.Nonce = nom.DeSerializeNonce(nonceFound)
	} else {
		transaction.FusedPlasma = response.BasePlasma
		transaction.Difficulty = uint64(0)
		if err := transaction.Nonce.UnmarshalText([]byte("0000000000000000")); err != nil {
			return err
		}
	}

	return nil
}

func setHashAndSignature(transaction *nom.AccountBlock, currentKeyPair wallet.KeyPair) {
	transaction.Hash = transaction.ComputeHash()
	transaction.Signature = currentKeyPair.Sign(transaction.Hash.Bytes())
}

func DebugAccountBlock(ab *api.AccountBlock) {
	CommonLogger.Debug("AccountBlock Type: %d\n", ab.BlockType)
	CommonLogger.Debug("AccountBlock Hash: %s\n", ab.Hash.String())
	CommonLogger.Debug("AccountBlock PrevHash: %s\n", ab.PreviousHash.String())
	CommonLogger.Debug("AccountBlock Height: %d\n", ab.Height)
	CommonLogger.Debug("AccountBlock MomentumAcknowledged: %v\n", ab.MomentumAcknowledged)
	CommonLogger.Debug("AccountBlock Address: %s\n", ab.Address.String())
	CommonLogger.Debug("AccountBlock ToAddress: %s\n", ab.ToAddress.String())
	CommonLogger.Debug("AccountBlock Amount %d\n", ab.Amount.Uint64())
	CommonLogger.Debug("AccountBlock FromBlockHash: %s\n", ab.FromBlockHash.String())
	CommonLogger.Debug("AccountBlock Data: %v\n", ab.Data)
}
