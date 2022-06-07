package api

import (
	"context"
	"github.com/zenon-network/go-zenon/common/types"
	"github.com/zenon-network/go-zenon/rpc/api/subscribe"
	"github.com/zenon-network/go-zenon/rpc/server"
	rpc "github.com/zenon-network/go-zenon/rpc/server"
)

type SubscriberApi struct {
	client *server.Client
}

func NewSubscriberApi(client *server.Client) *SubscriberApi {
	return &SubscriberApi{
		client: client,
	}
}

func (sa *SubscriberApi) ToMomentums() (*rpc.ClientSubscription, chan []subscribe.Momentum, error) {
	ctx := context.Background()
	ch := make(chan []subscribe.Momentum)
	subscription, err := sa.client.Subscribe(ctx, "ledger", ch, "momentums")
	if err != nil {
		return nil, nil, err
	}
	return subscription, ch, err
}

func (sa *SubscriberApi) ToAllAccountBlocks() (*rpc.ClientSubscription, chan []subscribe.AccountBlock, error) {
	ctx := context.Background()
	ch := make(chan []subscribe.AccountBlock)
	subscription, err := sa.client.Subscribe(ctx, "ledger", ch, "allAccountBlocks")
	if err != nil {
		return nil, nil, err
	}
	return subscription, ch, err
}

func (sa *SubscriberApi) ToAccountBlocksByAddress(address types.Address) (*rpc.ClientSubscription, chan []subscribe.AccountBlock, error) {
	ctx := context.Background()
	ch := make(chan []subscribe.AccountBlock)
	subscription, err := sa.client.Subscribe(ctx, "ledger", ch, "accountBlocksByAddress", address.String())
	if err != nil {
		return nil, nil, err
	}
	return subscription, ch, err
}

func (sa *SubscriberApi) ToUnreceivedAccountBlocksByAddress(address types.Address) (*rpc.ClientSubscription, chan []subscribe.AccountBlock, error) {
	ctx := context.Background()
	ch := make(chan []subscribe.AccountBlock)
	subscription, err := sa.client.Subscribe(ctx, "ledger", ch, "unreceivedAccountBlocksByAddress", address.String())
	if err != nil {
		return nil, nil, err
	}
	return subscription, ch, err
}
