package zenon

import (
	"errors"
	"os"
	"os/signal"
	"syscall"
	"znn-sdk-go/rpc_client"
	"znn-sdk-go/wallet"

	"github.com/zenon-network/go-zenon/chain/nom"
	"github.com/zenon-network/go-zenon/common/types"
	zwallet "github.com/zenon-network/go-zenon/wallet"
)

type Zenon struct {
	walletManager      *zwallet.Manager
	keyStore           *zwallet.KeyStore
	keyPair            *zwallet.KeyPair
	defaultKeyFilePath string

	Client *rpc_client.RpcClient
	// Channel wait for termination notifications
	stopCh chan os.Signal
}

func NewZenon(keyFilePath string) (*Zenon, error) {
	config := &zwallet.Config{
		WalletDir:      wallet.DefaultWalletDir,
		MaxSearchIndex: wallet.DefaultMaxIndex,
	}

	newWalletManager := zwallet.New(config)

	z := &Zenon{
		walletManager:      newWalletManager,
		defaultKeyFilePath: keyFilePath,
		stopCh:             make(chan os.Signal, 1),
	}
	signal.Notify(z.stopCh, syscall.SIGINT, syscall.SIGTERM)

	return z, nil
}

func (z *Zenon) Start(password, url string, index uint32) error {
	if err := z.walletManager.Start(); err != nil {
		return err
	}

	if z.defaultKeyFilePath != "" {
		keystore, err := z.walletManager.GetKeyFileAndDecrypt(z.defaultKeyFilePath, password)
		if err != nil {
			return err
		}
		z.keyStore = keystore
		_, keyPair, err := keystore.DeriveForIndexPath(index)
		if err != nil {
			return err
		}
		z.keyPair = keyPair
	}

	go func() {
		<-z.stopCh
		if err := z.Stop(); err != nil {
			CommonLogger.Error("Error while stopping Zenon SDK instance", "error", err)
		}
		os.Exit(1)
	}()

	var err error
	z.Client, err = rpc_client.NewRpcClient(url)
	if err != nil {
		return err
	}

	return nil
}

func (z *Zenon) Address() types.Address {
	if z.keyPair == nil {
		panic("keyPair not init")
	}
	return z.keyPair.Address
}

func (z *Zenon) Send(transaction *nom.AccountBlock) error {
	if z.keyPair == nil {
		return errors.New("keyPair is not initialized")
	}
	if err := checkAndSetFields(z.Client, transaction, *z.keyPair); err != nil {
		return err
	}
	if err := setDifficulty(z.Client, transaction); err != nil {
		return err
	}
	setHashAndSignature(transaction, *z.keyPair)

	if err := z.Client.LedgerApi.PublishRawTransaction(transaction); err != nil {
		return err
	}
	RPCLogger.Info("Successfully sent transaction\n")
	return nil
}

func (z *Zenon) Stop() error {
	if z.walletManager == nil {
		return wallet.ErrWalletManagerStopped
	}
	z.walletManager.Stop()
	if z.Client != nil {
		z.Client.Stop()
	}

	return nil
}
