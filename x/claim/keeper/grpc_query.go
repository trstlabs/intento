package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trstlabs/intento/x/claim/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

// Params returns p of the claim module.
func (q QueryServer) ModuleAccountBalance(
	c context.Context, _ *types.QueryModuleAccountBalanceRequest) (*types.QueryModuleAccountBalanceResponse, error) {

	ctx := sdk.UnwrapSDKContext(c)
	moduleAccBal := sdk.NewCoins(q.keeper.GetModuleAccountBalance(ctx))

	return &types.QueryModuleAccountBalanceResponse{ModuleAccountBalance: moduleAccBal}, nil
}

// Params returns params of the claim module.
func (q QueryServer) Params(goCtx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params, err := q.keeper.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, err
}

// ClaimRecord returns claimrecord per account
func (q QueryServer) ClaimRecord(
	goCtx context.Context,
	req *types.QueryClaimRecordRequest,
) (*types.QueryClaimRecordResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	claimRecord, err := q.keeper.GetClaimRecord(ctx, addr)
	return &types.QueryClaimRecordResponse{ClaimRecord: claimRecord}, err
}

// Activities returns activities
func (q QueryServer) ClaimableForAction(
	goCtx context.Context,
	req *types.QueryClaimableForActionRequest,
) (*types.QueryClaimableForActionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}
	p, err := q.keeper.GetParams(ctx)
	if err != nil {
		return nil, err
	}
	total, err := q.keeper.GetTotalClaimableAmountPerAction(ctx, addr)

	return &types.QueryClaimableForActionResponse{
		Total: sdk.NewCoin(p.ClaimDenom, total),
	}, err
}

// Activities returns activities
func (q QueryServer) TotalClaimable(
	goCtx context.Context,
	req *types.QueryTotalClaimableRequest,
) (*types.QueryTotalClaimableResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	total, err := q.keeper.GetTotalClaimableForAddr(ctx, addr)

	return &types.QueryTotalClaimableResponse{
		Total: total,
	}, err
}
