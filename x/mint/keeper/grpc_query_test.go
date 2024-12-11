package keeper_test

import (
	gocontext "context"
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/stretchr/testify/require"
	"github.com/trstlabs/intento/x/mint/keeper"
	"github.com/trstlabs/intento/x/mint/types"
)

func TestGRPCParams(t *testing.T) {
	app, ctx := createTestApp(true)
	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, keeper.NewQueryServer(app.MintKeeper))
	queryClient := types.NewQueryClient(queryHelper)

	params, err := queryClient.Params(gocontext.Background(), &types.QueryParamsRequest{})

	require.NoError(t, err)
	mintParams, _ := app.MintKeeper.GetParams(ctx)
	require.Equal(t, params.Params, mintParams)

	annualProvisions, err := queryClient.AnnualProvisions(gocontext.Background(), &types.QueryAnnualProvisionsRequest{})
	require.NoError(t, err)
	minter, _ := app.MintKeeper.GetMinter(ctx)
	require.Equal(t, annualProvisions.AnnualProvisions, minter.AnnualProvisions)
}
