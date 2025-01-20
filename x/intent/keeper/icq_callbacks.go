package keeper

import (
	"fmt"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trstlabs/intento/x/intent/types"
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
	k.Logger(ctx).Debug("Action Callback", "queryId", query.Id, "callbackArgs", args)
	// Parsing the actionID
	parts := strings.Split(query.Id, ":")
	if len(parts) < 3 {
		return fmt.Errorf("invalid query ID format: %s", query.Id)
	}

	// Parse actionID
	actionID, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse action ID: %w", err)
	}

	action, err := k.TryGetActionInfo(ctx, actionID)
	if err != nil {
		return err
	}

	// Parse the index
	index, err := strconv.Atoi(parts[2])
	if err != nil {
		return fmt.Errorf("failed to parse index: %w", err)
	}

	// Save to action ICQConfig
	// Check prefix to determine the type (feedback loop or condition)
	prefix := parts[0]
	switch prefix {
	case string(types.ActionFeedbackLoopQueryKeyPrefix):
		action.Conditions.FeedbackLoops[index].ICQConfig.Response = args
	case string(types.ActionComparisonQueryKeyPrefix): // Example prefix for conditions
		action.Conditions.Comparisons[index].ICQConfig.Response = args
	default:
		return fmt.Errorf("unknown prefix in query ID: %s", prefix)
	}

	k.SetActionInfo(ctx, &action)

	// Only if all responses are present (and thus it is the last one), handle the action
	for _, comparison := range action.Conditions.Comparisons {
		if comparison.ICQConfig != nil {
			if comparison.ICQConfig.Response == nil {
				return nil
			}
		}
	}

	for _, feedbackLoop := range action.Conditions.FeedbackLoops {
		if feedbackLoop.ICQConfig != nil {
			if feedbackLoop.ICQConfig.Response == nil {
				return nil
			}
		}
	}

	k.HandleAction(ctx, k.Logger(ctx), action, ctx.BlockTime(), args)
	return nil
}
