package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/trstlabs/intento/x/interchainquery/types"
)

var _ types.QueryServiceServer = &QueryServer{}

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

// Queries all queries that have been requested but have not received a response
func (q QueryServer) PendingQueries(c context.Context, req *types.QueryPendingQueriesRequest) (*types.QueryPendingQueriesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	pendingQueries := []types.Query{}
	for _, query := range q.keeper.AllQueries(ctx) {
		if query.RequestSent {
			pendingQueries = append(pendingQueries, query)
		}
	}

	return &types.QueryPendingQueriesResponse{PendingQueries: pendingQueries}, nil
}
