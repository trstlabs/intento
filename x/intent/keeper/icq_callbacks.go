package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	icqtypes "github.com/trstlabs/intento/x/interchainquery/types"
)

const (
	ICQCallbackID_Action = "action"
)

// ICQCallbacks wrapper struct for keeper
type ICQCallback func(Keeper, sdk.Context, []byte, icqtypes.Query) error

type ICQCallbacks struct {
	k         Keeper
	callbacks map[string]ICQCallback
}

var _ icqtypes.QueryCallbacks = ICQCallbacks{}

func (k Keeper) ICQCallbackHandler() ICQCallbacks {
	return ICQCallbacks{k, make(map[string]ICQCallback)}
}

func (c ICQCallbacks) CallICQCallback(ctx sdk.Context, id string, args []byte, query icqtypes.Query) error {
	return c.callbacks[id](c.k, ctx, args, query)
}

func (c ICQCallbacks) HasICQCallback(id string) bool {
	_, found := c.callbacks[id]
	return found
}

func (c ICQCallbacks) AddICQCallback(id string, fn interface{}) icqtypes.QueryCallbacks {
	c.callbacks[id] = fn.(ICQCallback)
	return c
}

func (c ICQCallbacks) RegisterICQCallbacks() icqtypes.QueryCallbacks {
	return c.
		AddICQCallback(ICQCallbackID_Action, ICQCallback(HandleActionCallback))

}

func HandleActionCallback(k Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	actionID, err := strconv.ParseUint(query.Id, 10, 64)
	if err != nil {
		return err
	}
	action, err := k.TryGetActionInfo(ctx, actionID)
	if err != nil {
		return err
	}

	k.HandleAction(ctx, k.Logger(ctx), action, ctx.BlockTime(), args)
	return nil
}
