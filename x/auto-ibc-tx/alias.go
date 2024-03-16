package autoibctx

import (
	"github.com/trstlabs/intento/x/auto-ibc-tx/keeper"
	"github.com/trstlabs/intento/x/auto-ibc-tx/types"
)

const (
	ModuleName = types.ModuleName
	StoreKey   = types.StoreKey
)

var (
	ExportGenesis = keeper.ExportGenesis
)
