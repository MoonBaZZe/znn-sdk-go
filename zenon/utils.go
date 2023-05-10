package zenon

import (
	"math/big"
	"znn-sdk-go/rpc_client"

	"github.com/inconshreveable/log15"
	"github.com/zenon-network/go-zenon/chain/nom"
	"github.com/zenon-network/go-zenon/common/crypto"
	"github.com/zenon-network/go-zenon/common/types"
	"github.com/zenon-network/go-zenon/pow"
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

func CheckAndSetFields(client *rpc_client.RpcClient, transaction *nom.AccountBlock, address types.Address, public []byte) error {
	transaction.Address = address
	transaction.PublicKey = public

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

func SetDifficulty(client *rpc_client.RpcClient, transaction *nom.AccountBlock) error {
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
	if response.RequiredDifficulty > 0 {
		transaction.FusedPlasma = response.AvailablePlasma
		transaction.Difficulty = response.RequiredDifficulty

		CommonLogger.Info("Generating Plasma, please wait\n")

		hashBytes := crypto.Hash(transaction.Address.Bytes(), transaction.PreviousHash.Bytes())
		powData, err := types.BytesToHash(hashBytes)
		if err != nil {
			return err
		}

		nonceFound := pow.GetPoWNonce(new(big.Int).SetUint64(response.RequiredDifficulty), powData)

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

func DebugAccountBlock(ab *nom.AccountBlock) {
	CommonLogger.Debug("AccountBlock Type: ", ab.BlockType, nil)
	CommonLogger.Debug("AccountBlock Hash: ", ab.Hash.String(), nil)
	CommonLogger.Debug("AccountBlock PrevHash: ", ab.PreviousHash.String(), nil)
	CommonLogger.Debug("AccountBlock Height: ", ab.Height, nil)
	CommonLogger.Debug("AccountBlock MomentumAcknowledged: ", ab.MomentumAcknowledged, nil)
	CommonLogger.Debug("AccountBlock Address: ", ab.Address.String(), nil)
	CommonLogger.Debug("AccountBlock ToAddress: ", ab.ToAddress.String(), nil)
	CommonLogger.Debug("AccountBlock Amount ", ab.Amount.Uint64(), nil)
	CommonLogger.Debug("AccountBlock FromBlockHash: ", ab.FromBlockHash.String(), nil)
	CommonLogger.Debug("AccountBlock Data: ", ab.Data, nil)
}
