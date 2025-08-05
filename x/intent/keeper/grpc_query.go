package keeper

import (
	"context"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"

	"github.com/trstlabs/intento/x/intent/types"
)

var _ types.QueryServer = &QueryServer{}

// QueryServer implements the module gRPC query service.
type QueryServer struct {
	keeper Keeper
}

// NewQueryServer creates a new gRPC query server.
func NewQueryServer(keeper Keeper) *QueryServer {
	return &QueryServer{
		keeper: keeper,
	}
}

// Params returns params of the alloc module.
func (q QueryServer) Params(goCtx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params, err := q.keeper.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, err
}

// InterchainAccountFromAddress implements the Query/InterchainAccountFromAddress gRPC method
func (q QueryServer) InterchainAccountFromAddress(goCtx context.Context, req *types.QueryInterchainAccountFromAddressRequest) (*types.QueryInterchainAccountFromAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	portID, err := icatypes.NewControllerPortID(req.Owner)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "could not find account: %s", err)
	}

	ica, found := q.keeper.icaControllerKeeper.GetInterchainAccountAddress(ctx, req.ConnectionId, portID)
	if !found {
		return nil, status.Errorf(codes.NotFound, "no account found for connectionID: %s and portID: %s",
			req.ConnectionId, portID)
	}

	return types.NewQueryInterchainAccountResponse(ica), nil
}

// Flow implements the Query/FlowgRPC method
func (q QueryServer) Flow(c context.Context, req *types.QueryFlowRequest) (*types.QueryFlowResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	id, err := strconv.ParseUint(req.Id, 10, 64)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	flowInfo, err := q.keeper.TryGetFlowInfo(ctx, id)
	if err != nil {
		return nil, err
	}
	// for msg := range flowInfo.Msgs{
	// 	makeReadableMsgData(&flowInfo, msg)
	// }

	return &types.QueryFlowResponse{
		FlowInfo: flowInfo,
	}, nil
}

func (q QueryServer) FlowHistory(ctx context.Context, req *types.QueryFlowHistoryRequest) (*types.QueryFlowHistoryResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	// Convert id from the request
	id, err := strconv.ParseUint(req.Id, 10, 64)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request")
	}

	// Assuming FlowHistoryEntry items are stored with keys prefixed by id

	store := runtime.KVStoreAdapter(q.keeper.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, types.GetFlowHistoryKey(id))
	// Paginate over the prefixed store
	var historyEntries []types.FlowHistoryEntry
	pageRes, err := query.Paginate(prefixStore, req.Pagination, func(key []byte, value []byte) error {
		var historyEntry types.FlowHistoryEntry
		if err := q.keeper.cdc.Unmarshal(value, &historyEntry); err != nil {
			return err
		}
		historyEntries = append(historyEntries, historyEntry)
		return nil
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "paginate: %v", err)
	}

	// Return paginated results
	return &types.QueryFlowHistoryResponse{
		History:    historyEntries,
		Pagination: pageRes,
	}, nil
}

// Flows implements the Query/Flows gRPC method
func (q QueryServer) Flows(c context.Context, req *types.QueryFlowsRequest) (*types.QueryFlowsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	flows := make([]types.FlowInfo, 0)
	store := runtime.KVStoreAdapter(q.keeper.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, types.FlowKeyPrefix)

	pageRes, err := query.FilteredPaginate(prefixStore, req.Pagination, func(_ []byte, value []byte, accumulate bool) (bool, error) {
		if accumulate {
			var c types.FlowInfo
			q.keeper.cdc.MustUnmarshal(value, &c)
			flows = append(flows, c)

		}
		return true, nil
	})

	if err != nil {
		return nil, err
	}

	return &types.QueryFlowsResponse{
		FlowInfos:  flows,
		Pagination: pageRes,
	}, nil
}

// FlowsForOwner implements the Query/FlowsForOwner gRPC method
func (q QueryServer) FlowsForOwner(c context.Context, req *types.QueryFlowsForOwnerRequest) (*types.QueryFlowsForOwnerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	flows := make([]types.FlowInfo, 0)

	ownerAddress, err := sdk.AccAddressFromBech32(req.Owner)
	if err != nil {
		return nil, err
	}
	store := runtime.KVStoreAdapter(q.keeper.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, types.GetFlowsByOwnerPrefix(ownerAddress))

	pageRes, err := query.FilteredPaginate(prefixStore, req.Pagination, func(key []byte, _ []byte, accumulate bool) (bool, error) {
		if accumulate {
			flowID := types.GetIDFromBytes(key)
			flowInfo := q.keeper.GetFlowInfo(ctx, flowID)

			flows = append(flows, flowInfo)

		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}

	return &types.QueryFlowsForOwnerResponse{
		FlowInfos:  flows,
		Pagination: pageRes,
	}, nil
}

// TrustlessExecutionAgent implements the Query/TrustlessExecutionAgent gRPC method
func (q QueryServer) TrustlessExecutionAgent(c context.Context, req *types.QueryTrustlessExecutionAgentRequest) (*types.QueryTrustlessExecutionAgentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	hosted, err := q.keeper.TryGetTrustlessExecutionAgent(ctx, req.AgentAddress)
	if err != nil {
		return nil, err
	}

	return &types.QueryTrustlessExecutionAgentResponse{
		TrustlessExecutionAgent: hosted,
	}, nil
}

// TrustlessExecutionAgents implements the Query/TrustlessExecutionAgents gRPC method
func (q QueryServer) TrustlessExecutionAgents(c context.Context, req *types.QueryTrustlessExecutionAgentsRequest) (*types.QueryTrustlessExecutionAgentsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	trustlessExecutionAgents := make([]types.TrustlessExecutionAgent, 0)
	store := runtime.KVStoreAdapter(q.keeper.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, types.TrustlessExecutionAgentKeyPrefix)

	pageRes, err := query.FilteredPaginate(prefixStore, req.Pagination, func(_ []byte, value []byte, accumulate bool) (bool, error) {
		if accumulate {
			var c types.TrustlessExecutionAgent
			q.keeper.cdc.MustUnmarshal(value, &c)
			trustlessExecutionAgents = append(trustlessExecutionAgents, c)

		}
		return true, nil
	})

	if err != nil {
		return nil, err
	}

	return &types.QueryTrustlessExecutionAgentsResponse{
		TrustlessExecutionAgents: trustlessExecutionAgents,
		Pagination:               pageRes,
	}, nil
}

// TrustlessExecutionAgentsByFeeAdmin implements the Query/TrustlessExTrustlessExecutionAgentsByFeeAdminecutionAgentsByAdmin gRPC method
func (q QueryServer) TrustlessExecutionAgentsByFeeAdmin(c context.Context, req *types.QueryTrustlessExecutionAgentsByFeeAdminRequest) (*types.QueryTrustlessExecutionAgentsByFeeAdminResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	trustlessExecutionAgents := make([]types.TrustlessExecutionAgent, 0)

	admin, err := sdk.AccAddressFromBech32(req.FeeAdmin)
	if err != nil {
		return nil, err
	}
	store := runtime.KVStoreAdapter(q.keeper.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, types.GetTrustlessExecutionAgentsByAdminPrefix(admin))
	pageRes, err := query.FilteredPaginate(prefixStore, req.Pagination, func(key []byte, _ []byte, accumulate bool) (bool, error) {
		if accumulate {
			trustlessExecutionAgentAddress := string(key)
			flowInfo := q.keeper.GetTrustlessExecutionAgent(ctx, trustlessExecutionAgentAddress)

			trustlessExecutionAgents = append(trustlessExecutionAgents, flowInfo)

		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}

	return &types.QueryTrustlessExecutionAgentsByFeeAdminResponse{
		TrustlessExecutionAgents: trustlessExecutionAgents,
		Pagination:               pageRes,
	}, nil
}
