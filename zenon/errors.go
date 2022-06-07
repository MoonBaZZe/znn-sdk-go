package zenon

import "errors"

var (
	ErrZeroFromHash            = errors.New("cannot receive a block that comes from a zero hash send block")
	ErrNonExistentAccountBlock = errors.New("account-block does not exist")
	ErrDifferentReceiver       = errors.New("cannot receive a block with a different receiver")
	ErrContainsData            = errors.New("cannot receive a block that contains data")
	ErrNoNonce                 = errors.New("cannot have difficulty and no nonce")
)
