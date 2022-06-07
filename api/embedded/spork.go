package embedded

import (
	"github.com/zenon-network/go-zenon/rpc/api/embedded"
	"github.com/zenon-network/go-zenon/rpc/server"
)

type SporkApi struct {
	client *server.Client
}

func NewSporkApi(client *server.Client) *SporkApi {
	return &SporkApi{
		client: client,
	}
}

func (sa *SporkApi) GetAll(pageIndex, pageSize uint32) (*embedded.SporkList, error) {
	ans := new(embedded.SporkList)
	if err := sa.client.Call(ans, "embedded.spork.getAll", pageIndex, pageSize); err != nil {
		return nil, err
	}
	return ans, nil
}
