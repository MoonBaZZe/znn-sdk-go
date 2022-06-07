package wallet

import "errors"

var (
	ErrWalletManagerStopped = errors.New("wallet manager has not started")
)
