package wallet

import (
	"crypto/rand"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/tyler-smith/go-bip39"
	"github.com/zenon-network/go-zenon/node"
	zwallet "github.com/zenon-network/go-zenon/wallet"
)

func WriteKeyFile(ks *zwallet.KeyStore, name, password string) error {
	kf, err := ks.Encrypt(password)
	if err != nil {
		return err
	}
	if name == "" {
		name = kf.BaseAddress.String()
	}
	kf.Path = filepath.Join(node.DefaultDataDir(), "wallet", name)
	keyFileJson, err := json.MarshalIndent(kf, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(kf.Path, keyFileJson, 0644)
}

func ReadKeyFile(name, password string) (*zwallet.KeyStore, error) {
	path := filepath.Join(node.DefaultDataDir(), "wallet", name)
	kf, err := zwallet.ReadKeyFile(path)
	if err != nil {
		return nil, err
	}
	ks, err := kf.Decrypt(password)
	if err != nil {
		return nil, err
	}

	return ks, nil
}

func NewKeyStoreFromEntropy(entropy []byte) (*zwallet.KeyStore, error) {
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, err
	}

	ks := &zwallet.KeyStore{
		Entropy:  entropy,
		Seed:     bip39.NewSeed(mnemonic, ""),
		Mnemonic: mnemonic,
	}

	// setup base address
	if _, kp, err := ks.DeriveForIndexPath(0); err == nil {
		ks.BaseAddress = kp.Address
	} else {
		return nil, err
	}

	return ks, nil
}

func NewKeyStore() (*zwallet.KeyStore, error) {
	entropy := make([]byte, 32)
	if _, err := rand.Read(entropy); err != nil {
		return nil, err
	}

	return NewKeyStoreFromEntropy(entropy)
}
