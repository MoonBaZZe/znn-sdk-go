package wallet

import (
	"path/filepath"

	"github.com/zenon-network/go-zenon/node"
	zwallet "github.com/zenon-network/go-zenon/wallet"
)

var (
	DefaultWalletDir = filepath.Join(node.DefaultDataDir(), node.DefaultWalletDir)
)

const (
	DefaultMaxIndex = zwallet.DefaultMaxIndex
)
