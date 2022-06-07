package main

import (
	"znn-sdk-go/zenon"
)

func main() {

	// Initialize zenon client with keyFile
	z, err := zenon.NewZenon("keyfile-sdk")

	if err != nil {
		zenon.CommonLogger.Error("Error while creating Zenon SDK instance", "error", err)
		return
	}

	if err := z.Start("123456", "ws://127.0.0.1:35998", 0); err != nil {
		zenon.CommonLogger.Error("Error while trying to connect to node", "error", err)
		return
	}

	if err := z.Stop(); err != nil {
		zenon.CommonLogger.Error("Error while stopping Zenon SDK instance", "error", err)
	}
}
