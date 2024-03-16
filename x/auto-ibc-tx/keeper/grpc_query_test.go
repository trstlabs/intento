package keeper

import (
	"errors"
	"testing"
	"time"

	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/trstlabs/intento/x/auto-ibc-tx/types"
)

func TestQueryAutoTxsByOwnerList(t *testing.T) {
	ctx, keepers, _ := CreateTestInput(t, false)
	autoTxKeeper := keepers.AutoIbcTxKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 1000000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin("denom", 500))

	creator, _ := CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, deposit)
	var expectedAutoTxs []types.AutoTxInfo
	portID, err := icatypes.NewControllerPortID(creator.String())
	require.NoError(t, err)

	// create 10 auto-txs
	for i := 0; i < 10; i++ {
		autoTx, err := CreateFakeAutoTx(autoTxKeeper, ctx, creator, portID, ibctesting.FirstConnectionID, time.Minute, time.Hour, ctx.BlockTime(), topUp)
		require.NoError(t, err)

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
			got, err := autoTxKeeper.AutoTxsForOwner(sdk.WrapSDKContext(ctx), spec.srcQuery)

			if spec.expErr != nil {
				require.Equal(t, spec.expErr, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
			for i, expectedAutoTx := range spec.expAutoTxInfos {
				assert.Equal(t, expectedAutoTx.GetTxMsgs(autoTxKeeper.cdc), got.AutoTxInfos[i].GetTxMsgs(autoTxKeeper.cdc))
				assert.Equal(t, expectedAutoTx.ICAConfig.PortID, got.AutoTxInfos[i].ICAConfig.PortID)
				assert.Equal(t, expectedAutoTx.Owner, got.AutoTxInfos[i].Owner)
				assert.Equal(t, expectedAutoTx.ICAConfig.ConnectionID, got.AutoTxInfos[i].ICAConfig.ConnectionID)
				assert.Equal(t, expectedAutoTx.Interval, got.AutoTxInfos[i].Interval)
				assert.Equal(t, expectedAutoTx.EndTime, got.AutoTxInfos[i].EndTime)
				assert.Equal(t, expectedAutoTx.Configuration, got.AutoTxInfos[i].Configuration)
			}
		})
	}
}

func TestQueryAutoTxHistory(t *testing.T) {
	ctx, keepers, _ := CreateTestInput(t, false)
	autoTxKeeper := keepers.AutoIbcTxKeeper

	autoTxHistory, err := CreateFakeAutoTxHistory(autoTxKeeper, ctx, ctx.BlockTime())
	require.NoError(t, err)

	ID := "1"
	got, err := autoTxKeeper.AutoTxHistory(sdk.WrapSDKContext(ctx), &types.QueryAutoTxHistoryRequest{Id: ID})
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, got.History[0].ScheduledExecTime, autoTxHistory.History[0].ScheduledExecTime)
	require.Equal(t, got.History[0].ActualExecTime, autoTxHistory.History[0].ActualExecTime)

}

func TestQueryAutoTxsList(t *testing.T) {
	ctx, keepers, _ := CreateTestInput(t, false)
	autoTxKeeper := keepers.AutoIbcTxKeeper
	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 1000000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin("denom", 500))

	creator, _ := CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, deposit)
	var expectedAutoTxs []types.AutoTxInfo
	portID, err := icatypes.NewControllerPortID(creator.String())
	require.NoError(t, err)

	// create 10 auto-txs
	for i := 0; i < 10; i++ {
		autoTx, err := CreateFakeAutoTx(autoTxKeeper, ctx, creator, portID, ibctesting.FirstConnectionID, time.Minute, time.Hour, ctx.BlockTime(), topUp)
		require.NoError(t, err)

		expectedAutoTxs = append(expectedAutoTxs, autoTx)
	}

	got, err := autoTxKeeper.AutoTxs(sdk.WrapSDKContext(ctx), &types.QueryAutoTxsRequest{})

	require.NoError(t, err)
	require.NotNil(t, got)
	for i, expectedAutoTx := range expectedAutoTxs {

		assert.Equal(t, expectedAutoTx.GetTxMsgs(autoTxKeeper.cdc), got.AutoTxInfos[i].GetTxMsgs(autoTxKeeper.cdc))
		assert.Equal(t, expectedAutoTx.ICAConfig.PortID, got.AutoTxInfos[i].ICAConfig.PortID)
		assert.Equal(t, expectedAutoTx.Owner, got.AutoTxInfos[i].Owner)
		assert.Equal(t, expectedAutoTx.ICAConfig.ConnectionID, got.AutoTxInfos[i].ICAConfig.ConnectionID)
		assert.Equal(t, expectedAutoTx.Interval, got.AutoTxInfos[i].Interval)
		assert.Equal(t, expectedAutoTx.EndTime, got.AutoTxInfos[i].EndTime)
		assert.Equal(t, expectedAutoTx.Configuration, got.AutoTxInfos[i].Configuration)
		assert.Equal(t, expectedAutoTx.UpdateHistory, got.AutoTxInfos[i].UpdateHistory)
	}
}

func TestQueryAutoTxsListWithAuthZMsg(t *testing.T) {
	ctx, keepers, _ := CreateTestInput(t, false)
	autoTxKeeper := keepers.AutoIbcTxKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 1000000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin("denom", 500))

	creator, _ := CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, deposit)

	portID, err := icatypes.NewControllerPortID(creator.String())
	require.NoError(t, err)

	expectedAutoTx, err := CreateFakeAuthZAutoTx(autoTxKeeper, ctx, creator, portID, ibctesting.FirstConnectionID, time.Minute, time.Hour, ctx.BlockTime(), topUp)
	require.NoError(t, err)
	got, err := autoTxKeeper.AutoTxs(sdk.WrapSDKContext(ctx), &types.QueryAutoTxsRequest{})

	require.NoError(t, err)
	require.NotNil(t, got)

	var txMsg sdk.Msg
	autoTxKeeper.cdc.UnpackAny(expectedAutoTx.Msgs[0], &txMsg)

	var gotMsg sdk.Msg
	autoTxKeeper.cdc.UnpackAny(got.AutoTxInfos[0].Msgs[0], &gotMsg)

	assert.Equal(t, expectedAutoTx.Msgs, got.AutoTxInfos[0].Msgs)
	//	assert.Equal(t, txMsg, gotMsg)
	assert.Equal(t, expectedAutoTx.ICAConfig.PortID, got.AutoTxInfos[0].ICAConfig.PortID)
	assert.Equal(t, expectedAutoTx.Owner, got.AutoTxInfos[0].Owner)
	assert.Equal(t, expectedAutoTx.ICAConfig.ConnectionID, got.AutoTxInfos[0].ICAConfig.ConnectionID)
	assert.Equal(t, expectedAutoTx.Interval, got.AutoTxInfos[0].Interval)
	assert.Equal(t, expectedAutoTx.EndTime, got.AutoTxInfos[0].EndTime)
	assert.Equal(t, expectedAutoTx.Configuration, got.AutoTxInfos[0].Configuration)
	assert.Equal(t, expectedAutoTx.UpdateHistory, got.AutoTxInfos[0].UpdateHistory)

}

func TestQueryParams(t *testing.T) {
	ctx, keepers, _ := CreateTestInput(t, false)
	autoTxKeeper := keepers.AutoIbcTxKeeper

	resp, err := autoTxKeeper.Params(sdk.WrapSDKContext(ctx), &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, resp.Params, types.DefaultParams())
}

func NewICA(t *testing.T) (*ibctesting.Coordinator, *ibctesting.Path) {
	coordinator := ibctesting.NewCoordinator(t, 2)
	chainA := coordinator.GetChain(ibctesting.GetChainID(1))
	chainB := coordinator.GetChain(ibctesting.GetChainID(2))
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = icatypes.HostPortID
	path.EndpointB.ChannelConfig.PortID = icatypes.HostPortID
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

		StartTime:     startAt,
		ExecTime:      execTime,
		EndTime:       endTime,
		ICAConfig:     &types.ICAConfig{PortID: portID},
		Configuration: &types.ExecutionConfiguration{SaveMsgResponses: true},
	}

	k.SetAutoTxInfo(ctx, &autoTx)
	k.addToAutoTxOwnerIndex(ctx, owner, txID)

	var newAutoTx types.AutoTxInfo
	autoTxBz := k.cdc.MustMarshal(&autoTx)
	k.cdc.MustUnmarshal(autoTxBz, &newAutoTx)
	return newAutoTx, nil
}

func CreateFakeAutoTxHistory(k Keeper, ctx sdk.Context, startAt time.Time) (types.AutoTxHistory, error) {

	entry := types.AutoTxHistoryEntry{
		ScheduledExecTime: startAt.Add(time.Minute),
		ActualExecTime:    startAt.Add(time.Minute).Add(time.Microsecond),
		Errors:            []string{"text"},
		Executed:          true,
	}

	autoTxHistory := types.AutoTxHistory{
		History: []types.AutoTxHistoryEntry{entry},
	}

	k.SetAutoTxHistory(ctx, 1, &autoTxHistory)

	return autoTxHistory, nil
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
		ICAConfig:     &types.ICAConfig{PortID: portID},
	}
	k.SetAutoTxInfo(ctx, &autoTx)
	k.addToAutoTxOwnerIndex(ctx, owner, txID)
	autoTxBz := k.cdc.MustMarshal(&autoTx)
	var newAutoTx types.AutoTxInfo
	k.cdc.MustUnmarshal(autoTxBz, &newAutoTx)
	return newAutoTx, nil
}
