package keeper

import (
	"errors"
	"fmt"
	"testing"
	"time"

	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	icatypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v4/testing"
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
	var expectedAutoTxs []types.AutoTxInfo
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
		autoTx, err := CreateFakeAutoTx(keepers.AutoIbcTxKeeper, ctx, creator, portID, ibctesting.FirstConnectionID, time.Minute, time.Hour, ctx.BlockTime(), topUp)
		require.NoError(t, err)
		// msg, err := icatypes.DeserializeCosmosTx(keepers.AutoIbcTxKeeper.cdc, autoTx.Data)
		// require.NoError(t, err)
		// makeReadableMsgData(&autoTx, msg)
		expectedAutoTxs = append(expectedAutoTxs, autoTx)
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
			expAutoTxInfos: expectedAutoTxs,
			expErr:         nil,
		},
		"with pagination offset": {
			srcQuery: &types.QueryAutoTxsForOwnerRequest{
				Owner: creator.String(),
				Pagination: &query.PageRequest{
					Offset: 1,
				},
			},
			expAutoTxInfos: expectedAutoTxs[1:],
			expErr:         nil,
		},
		"with pagination limit": {
			srcQuery: &types.QueryAutoTxsForOwnerRequest{
				Owner: creator.String(),
				Pagination: &query.PageRequest{
					Limit: 1,
				},
			},
			expAutoTxInfos: expectedAutoTxs[0:1],
			expErr:         nil,
		},
		"nil creator": {
			srcQuery: &types.QueryAutoTxsForOwnerRequest{
				Pagination: &query.PageRequest{},
			},
			expAutoTxInfos: expectedAutoTxs,
			expErr:         errors.New("empty address string is not allowed"),
		},
		"nil req": {
			srcQuery:       nil,
			expAutoTxInfos: expectedAutoTxs,
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
			for i, expectedAutoTx := range spec.expAutoTxInfos {
				assert.Equal(t, expectedAutoTx.GetTxMsgs(), got.AutoTxInfos[i].GetTxMsgs())
				assert.Equal(t, expectedAutoTx.AutoTxHistory, got.AutoTxInfos[i].AutoTxHistory)
				assert.Equal(t, expectedAutoTx.PortID, got.AutoTxInfos[i].PortID)
				assert.Equal(t, expectedAutoTx.Owner, got.AutoTxInfos[i].Owner)
				assert.Equal(t, expectedAutoTx.ConnectionID, got.AutoTxInfos[i].ConnectionID)
				assert.Equal(t, expectedAutoTx.Interval, got.AutoTxInfos[i].Interval)
				assert.Equal(t, expectedAutoTx.EndTime, got.AutoTxInfos[i].EndTime)
				assert.Equal(t, expectedAutoTx.DependsOnTxIds, got.AutoTxInfos[i].DependsOnTxIds)
			}
		})
	}
}

func TestQueryAutoTxsList(t *testing.T) {
	ctx, keepers := CreateTestInput(t, false)

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 1000000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin("denom", 500))

	creator, _ := CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, deposit)
	var expectedAutoTxs []types.AutoTxInfo
	portID, err := icatypes.NewControllerPortID(creator.String())
	require.NoError(t, err)

	// create 10 auto-txs
	for i := 0; i < 10; i++ {
		autoTx, err := CreateFakeAutoTx(keepers.AutoIbcTxKeeper, ctx, creator, portID, ibctesting.FirstConnectionID, time.Minute, time.Hour, ctx.BlockTime(), topUp)
		require.NoError(t, err)

		expectedAutoTxs = append(expectedAutoTxs, autoTx)
	}

	got, err := keepers.AutoIbcTxKeeper.AutoTxs(sdk.WrapSDKContext(ctx), &types.QueryAutoTxsRequest{})

	require.NoError(t, err)
	require.NotNil(t, got)
	for i, expectedAutoTx := range expectedAutoTxs {

		assert.Equal(t, expectedAutoTx.GetTxMsgs(), got.AutoTxInfos[i].GetTxMsgs())
		assert.Equal(t, expectedAutoTx.AutoTxHistory, got.AutoTxInfos[i].AutoTxHistory)
		assert.Equal(t, expectedAutoTx.PortID, got.AutoTxInfos[i].PortID)
		assert.Equal(t, expectedAutoTx.Owner, got.AutoTxInfos[i].Owner)
		assert.Equal(t, expectedAutoTx.ConnectionID, got.AutoTxInfos[i].ConnectionID)
		assert.Equal(t, expectedAutoTx.Interval, got.AutoTxInfos[i].Interval)
		assert.Equal(t, expectedAutoTx.EndTime, got.AutoTxInfos[i].EndTime)
		assert.Equal(t, expectedAutoTx.DependsOnTxIds, got.AutoTxInfos[i].DependsOnTxIds)
		assert.Equal(t, expectedAutoTx.UpdateHistory, got.AutoTxInfos[i].UpdateHistory)
	}
}

func TestQueryAutoTxsListWithAuthZMsg(t *testing.T) {
	ctx, keepers := CreateTestInput(t, false)

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 1000000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin("denom", 500))

	creator, _ := CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, deposit)

	portID, err := icatypes.NewControllerPortID(creator.String())
	require.NoError(t, err)

	expectedAutoTx, err := CreateFakeAuthZAutoTx(keepers.AutoIbcTxKeeper, ctx, creator, portID, ibctesting.FirstConnectionID, time.Minute, time.Hour, ctx.BlockTime(), topUp)
	require.NoError(t, err)
	fmt.Printf("%v\n", len(expectedAutoTx.Msgs))
	got, err := keepers.AutoIbcTxKeeper.AutoTxs(sdk.WrapSDKContext(ctx), &types.QueryAutoTxsRequest{})

	require.NoError(t, err)
	require.NotNil(t, got)

	var txMsg sdk.Msg
	keepers.AutoIbcTxKeeper.cdc.UnpackAny(expectedAutoTx.Msgs[0], &txMsg)

	var gotMsg sdk.Msg
	keepers.AutoIbcTxKeeper.cdc.UnpackAny(got.AutoTxInfos[0].Msgs[0], &gotMsg)

	assert.Equal(t, expectedAutoTx.Msgs, got.AutoTxInfos[0].Msgs)
	//	assert.Equal(t, txMsg, gotMsg)
	assert.Equal(t, expectedAutoTx.AutoTxHistory, got.AutoTxInfos[0].AutoTxHistory)
	assert.Equal(t, expectedAutoTx.PortID, got.AutoTxInfos[0].PortID)
	assert.Equal(t, expectedAutoTx.Owner, got.AutoTxInfos[0].Owner)
	assert.Equal(t, expectedAutoTx.ConnectionID, got.AutoTxInfos[0].ConnectionID)
	assert.Equal(t, expectedAutoTx.Interval, got.AutoTxInfos[0].Interval)
	assert.Equal(t, expectedAutoTx.EndTime, got.AutoTxInfos[0].EndTime)
	assert.Equal(t, expectedAutoTx.DependsOnTxIds, got.AutoTxInfos[0].DependsOnTxIds)
	assert.Equal(t, expectedAutoTx.UpdateHistory, got.AutoTxInfos[0].UpdateHistory)

}

func TestQueryParams(t *testing.T) {
	ctx, keepers := CreateTestInput(t, false)
	resp, err := keepers.AutoIbcTxKeeper.Params(sdk.WrapSDKContext(ctx), &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, resp.Params, types.DefaultParams())
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

func CreateFakeAutoTx(k Keeper, ctx sdk.Context, owner sdk.AccAddress, portID, connectionId string, duration time.Duration, interval time.Duration, startAt time.Time, feeFunds sdk.Coins) (types.AutoTxInfo, error) {

	txID := k.autoIncrementID(ctx, types.KeyLastTxID)
	autoTxAddress, err := k.createFeeAccount(ctx, txID, owner, feeFunds)
	if err != nil {
		return types.AutoTxInfo{}, err
	}
	fakeMsg := banktypes.NewMsgSend(owner, autoTxAddress, feeFunds)

	anys, err := types.PackTxMsgAnys([]sdk.Msg{fakeMsg})
	if err != nil {
		return types.AutoTxInfo{}, err
	}

	// fakeData, _ := icatypes.SerializeCosmosTx(k.cdc, []sdk.Msg{fakeMsg})
	endTime, execTime := k.calculateTimeAndInsertQueue(ctx, startAt, duration, txID, interval)
	autoTx := types.AutoTxInfo{
		TxID:       txID,
		FeeAddress: autoTxAddress.String(),
		Owner:      owner.String(),
		// Data:       fakeData,
		Msgs:     anys,
		Interval: interval,

		StartTime: startAt,
		ExecTime:  execTime,
		EndTime:   endTime,
		PortID:    portID,
	}

	k.SetAutoTxInfo(ctx, &autoTx)
	k.addToAutoTxOwnerIndex(ctx, owner, txID)

	var newAutoTx types.AutoTxInfo
	autoTxBz := k.cdc.MustMarshal(&autoTx)
	k.cdc.MustUnmarshal(autoTxBz, &newAutoTx)
	return newAutoTx, nil
}

func CreateFakeAuthZAutoTx(k Keeper, ctx sdk.Context, owner sdk.AccAddress, portID, connectionId string, duration time.Duration, interval time.Duration, startAt time.Time, feeFunds sdk.Coins) (types.AutoTxInfo, error) {

	txID := k.autoIncrementID(ctx, types.KeyLastTxID)
	autoTxAddress, err := k.createFeeAccount(ctx, txID, owner, feeFunds)
	if err != nil {
		return types.AutoTxInfo{}, err
	}
	fakeMsg := banktypes.NewMsgSend(owner, autoTxAddress, feeFunds)
	anys, err := types.PackTxMsgAnys([]sdk.Msg{fakeMsg})
	if err != nil {
		return types.AutoTxInfo{}, err
	}
	fakeAuthZMsg := authztypes.MsgExec{Grantee: "ICA_ADDR", Msgs: anys}

	//fakeAuthZMsg := feegranttypes.Se{Grantee: "ICA_ADDR", Msgs: anys}
	anys, err = types.PackTxMsgAnys([]sdk.Msg{&fakeAuthZMsg})
	if err != nil {
		return types.AutoTxInfo{}, err
	}

	// fakeData, _ := icatypes.SerializeCosmosTx(k.cdc, []sdk.Msg{fakeMsg})
	endTime, execTime := k.calculateTimeAndInsertQueue(ctx, startAt, duration, txID, interval)
	autoTx := types.AutoTxInfo{
		TxID:       txID,
		FeeAddress: autoTxAddress.String(),
		Owner:      owner.String(),
		// Data:       fakeData,
		Msgs:          anys,
		Interval:      interval,
		UpdateHistory: nil,
		StartTime:     startAt,
		ExecTime:      execTime,
		EndTime:       endTime,
		PortID:        portID,
	}
	k.SetAutoTxInfo(ctx, &autoTx)
	k.addToAutoTxOwnerIndex(ctx, owner, txID)
	autoTxBz := k.cdc.MustMarshal(&autoTx)
	var newAutoTx types.AutoTxInfo
	k.cdc.MustUnmarshal(autoTxBz, &newAutoTx)
	return newAutoTx, nil
}
