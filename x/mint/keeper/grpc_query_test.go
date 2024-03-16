package keeper_test

import (
	gocontext "context"
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/stretchr/testify/require"
	"github.com/trstlabs/intento/x/mint/types"
)

func TestGRPCParams(t *testing.T) {
	app, ctx := createTestApp(true)
	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.MintKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	params, err := queryClient.Params(gocontext.Background(), &types.QueryParamsRequest{})

	require.NoError(t, err)
	require.Equal(t, params.Params, app.MintKeeper.GetParams(ctx))

	annualProvisions, err := queryClient.AnnualProvisions(gocontext.Background(), &types.QueryAnnualProvisionsRequest{})
	require.NoError(t, err)
	require.Equal(t, annualProvisions.AnnualProvisions, app.MintKeeper.GetMinter(ctx).AnnualProvisions)
}
