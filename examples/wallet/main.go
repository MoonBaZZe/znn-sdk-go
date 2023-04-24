package main

import (
	"znn-sdk-go/wallet"
	"znn-sdk-go/zenon"
)

func main() {
	kf, err := wallet.NewKeyStore()
	if err != nil {
		zenon.WalletLogger.Error("Error while creating key file", "wallet", err)
	} else {
		if err := wallet.WriteKeyFile(kf, "keyfile-sdk", "123456"); err != nil {
			zenon.WalletLogger.Error("Error while trying to create key file", "wallet", err)
		}
		if kf, err = wallet.ReadKeyFile("keyfile-sdk", "123456", ""); err != nil {
			zenon.WalletLogger.Error("Error while trying to read key file", "wallet", err)
		}
	}
}
