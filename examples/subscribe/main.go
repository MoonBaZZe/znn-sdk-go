package main

import (
	"fmt"
	"github.com/MoonBaZZe/znn-sdk-go/zenon"

	"github.com/zenon-network/go-zenon/rpc/api/subscribe"
)

func main() {
	z, err := zenon.NewZenon("keyfile-sdk")

	if err != nil {
		zenon.CommonLogger.Error("Error while creating Zenon SDK instance", "error", err)
		return
	}

	if err := z.Start("123456", "ws://127.0.0.1:35998", 0); err != nil {
		zenon.CommonLogger.Error("Error while trying to connect to node", "error", err)
		return
	}

	var ch chan []subscribe.Momentum
	if _, ch, err = z.Client.SubscriberApi.ToMomentums(); err != nil {
		zenon.CommonLogger.Error("Error while trying to call subscribe method", "error", err)
	}

	for i := 0; i < 5; i++ {
		select {
		case momentum := <-ch:
			fmt.Println("Height: ", momentum[0].Height, "; Hash: ", momentum[0].Hash)
		}
	}

	if err := z.Stop(); err != nil {
		zenon.CommonLogger.Error("Error while stopping Zenon SDK instance", "Error", err)
	}
}
