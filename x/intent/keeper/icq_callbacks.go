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
	ICQCallbackID_Flow = "flow"
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
		AddICQCallback(ICQCallbackID_Flow, ICQCallback(HandleFlowCallback))

}

func HandleFlowCallback(k Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	k.Logger(ctx).Debug("Flow Callback", "queryId", query.Id, "callbackArgs", args)
	// Parsing the flowID
	parts := strings.Split(query.Id, ":")
	if len(parts) < 3 {
		return fmt.Errorf("invalid query ID format: %s", query.Id)
	}

	// Parse flowID
	flowID, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse flow ID: %w", err)
	}

	flow, err := k.TryGetFlowInfo(ctx, flowID)
	if err != nil {
		return err
	}

	// Parse the index
	index, err := strconv.Atoi(parts[2])
	if err != nil {
		return fmt.Errorf("failed to parse index: %w", err)
	}

	// Save to flow ICQConfig
	// Check prefix to determine the type (feedback loop or condition)
	prefix := parts[0]
	switch prefix {
	case string(types.FlowFeedbackLoopQueryKeyPrefix):
		flow.Conditions.FeedbackLoops[index].ICQConfig.Response = args
	case string(types.FlowComparisonQueryKeyPrefix): // Example prefix for conditions
		flow.Conditions.Comparisons[index].ICQConfig.Response = args
	default:
		return fmt.Errorf("unknown prefix in query ID: %s", prefix)
	}

	k.SetFlowInfo(ctx, &flow)

	// Only if all responses are present (and thus it is the last one), handle the flow
	for _, comparison := range flow.Conditions.Comparisons {
		if comparison.ICQConfig != nil {
			if comparison.ICQConfig.Response == nil {
				return nil
			}
		}
	}

	for _, feedbackLoop := range flow.Conditions.FeedbackLoops {
		if feedbackLoop.ICQConfig != nil {
			if feedbackLoop.ICQConfig.Response == nil {
				return nil
			}
		}
	}

	k.HandleFlow(ctx, k.Logger(ctx), flow, ctx.BlockTime(), args)
	return nil
}
