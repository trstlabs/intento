package keeper

import (
	"github.com/trstlabs/trst/x/item/types"
)

var _ types.QueryServer = Keeper{}
