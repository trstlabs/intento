package keeper

import (
	"errors"
	"testing"
	"time"

	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v3/testing"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/trstlabs/trst/x/auto-ibc-tx/types"
)

func TestQueryAutoTxsByOwnerList(t *testing.T) {
	ctx, keepers := CreateTestInput(t, false)

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 1000000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin("denom", 500))
	//coordinator, path := NewICA(t)
	creator, _ := CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, deposit)
	var allExpecedAutoTxs []types.AutoTxInfo
	portID, err := icatypes.NewControllerPortID(creator.String())
	require.NoError(t, err)
	//anyAddr, _ := CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, deposit)

	/*

		require.NoError(t, err)

		coordinator.SetupConnections(path)

		keepers.AutoIbcTxKeeper.RegisterInterchainAccount(ctx, ibctesting.FirstConnectionID, creator.String())
		require.NoError(t, err)

		channel := channeltypes.NewChannel(
			channeltypes.OPEN,
			channeltypes.ORDERED,
			channeltypes.NewCounterparty(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID),
			[]string{path.EndpointA.ConnectionID},
			path.EndpointA.ChannelConfig.Version,
		)

		keepers.IbcKeeper.ChannelKeeper.SetChannel(ctx, portID, ibctesting.FirstChannelID, channel)

		keepers.ICAControllerKeeper.SetActiveChannelID(ctx, ibctesting.FirstConnectionID, portID, ibctesting.FirstChannelID)*/

	// create 10 auto-txs
	for i := 0; i < 10; i++ {
		autoTx, err := CreateFakeAutoTx(keepers.AutoIbcTxKeeper, ctx, creator, portID, []byte("fake_ica_msg"), ibctesting.FirstConnectionID, time.Minute, time.Hour, ctx.BlockTime(), topUp)
		require.NoError(t, err)
		allExpecedAutoTxs = append(allExpecedAutoTxs, autoTx)

	}

	specs := map[string]struct {
		srcQuery       *types.QueryAutoTxsForOwnerRequest
		expAutoTxInfos []types.AutoTxInfo
		expErr         error
	}{
		"query all": {
			srcQuery: &types.QueryAutoTxsForOwnerRequest{
				Owner: creator.String(),
			},
			expAutoTxInfos: allExpecedAutoTxs,
			expErr:         nil,
		},
		"with pagination offset": {
			srcQuery: &types.QueryAutoTxsForOwnerRequest{
				Owner: creator.String(),
				Pagination: &query.PageRequest{
					Offset: 1,
				},
			},
			expAutoTxInfos: allExpecedAutoTxs[1:],
			expErr:         nil,
		},
		"with pagination limit": {
			srcQuery: &types.QueryAutoTxsForOwnerRequest{
				Owner: creator.String(),
				Pagination: &query.PageRequest{
					Limit: 1,
				},
			},
			expAutoTxInfos: allExpecedAutoTxs[0:1],
			expErr:         nil,
		},
		"nil creator": {
			srcQuery: &types.QueryAutoTxsForOwnerRequest{
				Pagination: &query.PageRequest{},
			},
			expAutoTxInfos: allExpecedAutoTxs,
			expErr:         errors.New("empty address string is not allowed"),
		},
		"nil req": {
			srcQuery:       nil,
			expAutoTxInfos: allExpecedAutoTxs,
			expErr:         status.Error(codes.InvalidArgument, "empty request"),
		},
	}

	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			got, err := keepers.AutoIbcTxKeeper.AutoTxsForOwner(sdk.WrapSDKContext(ctx), spec.srcQuery)

			if spec.expErr != nil {
				require.Equal(t, spec.expErr, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, spec.expAutoTxInfos, got.AutoTxInfos)
		})
	}
}

func NewICA(t *testing.T) (*ibctesting.Coordinator, *ibctesting.Path) {
	coordinator := ibctesting.NewCoordinator(t, 2)
	chainA := coordinator.GetChain(ibctesting.GetChainID(1))
	chainB := coordinator.GetChain(ibctesting.GetChainID(2))
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = icatypes.PortID
	path.EndpointB.ChannelConfig.PortID = icatypes.PortID
	path.EndpointA.ChannelConfig.Order = channeltypes.ORDERED
	path.EndpointB.ChannelConfig.Order = channeltypes.ORDERED
	// path.EndpointA.ChannelConfig.Version = TestVersion
	// path.EndpointB.ChannelConfig.Version = TestVersion

	return coordinator, path
}

func CreateFakeAutoTx(k Keeper, ctx sdk.Context, owner sdk.AccAddress, portID string, data []byte, connectionId string, duration time.Duration, interval time.Duration, startAt time.Time, feeFunds sdk.Coins) (types.AutoTxInfo, error) {

	txID := k.autoIncrementID(ctx, types.KeyLastTxID)
	autoTxAddress, err := k.createFeeAccount(ctx, txID, owner, feeFunds)
	if err != nil {
		return types.AutoTxInfo{}, err
	}
	endTime, execTime, interval := k.calculateAndInsertQueue(ctx, startAt, duration, txID, interval)
	autoTx := types.AutoTxInfo{
		TxID:      txID,
		Address:   autoTxAddress,
		Owner:     owner,
		Data:      data,
		Interval:  interval,
		Duration:  duration,
		StartTime: startAt,
		ExecTime:  execTime,
		EndTime:   endTime,
		PortID:    portID,
	}

	k.SetAutoTxInfo(ctx, &autoTx)
	k.addToAutoTxOwnerIndex(ctx, owner /* startAt, */, txID)
	return autoTx, nil
}
