package main

import (
	"fmt"
	"github.com/MoonBaZZe/znn-sdk-go/zenon"
)

func main() {

	// Initialize zenon client with keyFile
	z, err := zenon.NewZenon("")

	if err != nil {
		zenon.CommonLogger.Error("Error while creating Zenon SDK instance", "error", err)
		return
	}

	if err := z.Start("", "ws://127.0.0.1:35998", 0); err != nil {
		zenon.CommonLogger.Error("Error while trying to connect to node", "error", err)
		return
	}

	if momentumList, err := z.Client.LedgerApi.GetDetailedMomentumsByHeight(5, 5); err != nil {
		zenon.CommonLogger.Error("Error while trying to call RPC", "error", err)
	} else {
		for _, momentum := range momentumList.List {
			fmt.Println("Momentum height: ", momentum.Momentum.Height)
			fmt.Println("Momentum hash: ", momentum.Momentum.Hash.String())
		}
	}

	if pillarInfoList, err := z.Client.PillarApi.GetAll(0, 5); err != nil {
		zenon.CommonLogger.Error("Error while trying to call RPC", "error", err)
	} else {
		for _, pillar := range pillarInfoList.List {
			fmt.Println("Pillar: ", pillar)
		}
	}

	if networkInfo, err := z.Client.StatsApi.NetworkInfo(); err != nil {
		zenon.CommonLogger.Error("Error while trying to call RPC", "error", err)
	} else {
		fmt.Println("NetworkInfo node public key: ", networkInfo.Self.PublicKey)
	}

	if err := z.Stop(); err != nil {
		zenon.CommonLogger.Error("Error while stopping Zenon SDK instance", "Error", err)
	}
}
