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

	"github.com/trstlabs/intento/x/intent/types"
)

func TestQueryActionsByOwnerList(t *testing.T) {
	ctx, keepers, _ := CreateTestInput(t, false)
	autoTxKeeper := keepers.IntentKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 1000000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin("denom", 500))

	creator, _ := CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, deposit)
	var expectedActions []types.ActionInfo
	portID, err := icatypes.NewControllerPortID(creator.String())
	require.NoError(t, err)

	// create 10 actions
	for i := 0; i < 10; i++ {
		autoTx, err := CreateFakeAction(autoTxKeeper, ctx, creator, portID, ibctesting.FirstConnectionID, time.Minute, time.Hour, ctx.BlockTime(), topUp)
		require.NoError(t, err)

		expectedActions = append(expectedActions, autoTx)
	}

	specs := map[string]struct {
		srcQuery       *types.QueryActionsForOwnerRequest
		expActionInfos []types.ActionInfo
		expErr         error
	}{
		"query all": {
			srcQuery: &types.QueryActionsForOwnerRequest{
				Owner: creator.String(),
			},
			expActionInfos: expectedActions,
			expErr:         nil,
		},
		"with pagination offset": {
			srcQuery: &types.QueryActionsForOwnerRequest{
				Owner: creator.String(),
				Pagination: &query.PageRequest{
					Offset: 1,
				},
			},
			expActionInfos: expectedActions[1:],
			expErr:         nil,
		},
		"with pagination limit": {
			srcQuery: &types.QueryActionsForOwnerRequest{
				Owner: creator.String(),
				Pagination: &query.PageRequest{
					Limit: 1,
				},
			},
			expActionInfos: expectedActions[0:1],
			expErr:         nil,
		},
		"nil creator": {
			srcQuery: &types.QueryActionsForOwnerRequest{
				Pagination: &query.PageRequest{},
			},
			expActionInfos: expectedActions,
			expErr:         errors.New("empty address string is not allowed"),
		},
		"nil req": {
			srcQuery:       nil,
			expActionInfos: expectedActions,
			expErr:         status.Error(codes.InvalidArgument, "empty request"),
		},
	}

	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			got, err := autoTxKeeper.ActionsForOwner(sdk.WrapSDKContext(ctx), spec.srcQuery)

			if spec.expErr != nil {
				require.Equal(t, spec.expErr, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
			for i, expectedAction := range spec.expActionInfos {
				assert.Equal(t, expectedAction.GetTxMsgs(autoTxKeeper.cdc), got.ActionInfos[i].GetTxMsgs(autoTxKeeper.cdc))
				assert.Equal(t, expectedAction.ICAConfig.PortID, got.ActionInfos[i].ICAConfig.PortID)
				assert.Equal(t, expectedAction.Owner, got.ActionInfos[i].Owner)
				assert.Equal(t, expectedAction.ICAConfig.ConnectionID, got.ActionInfos[i].ICAConfig.ConnectionID)
				assert.Equal(t, expectedAction.Interval, got.ActionInfos[i].Interval)
				assert.Equal(t, expectedAction.EndTime, got.ActionInfos[i].EndTime)
				assert.Equal(t, expectedAction.Configuration, got.ActionInfos[i].Configuration)
			}
		})
	}
}

func TestQueryActionHistory(t *testing.T) {
	ctx, keepers, _ := CreateTestInput(t, false)
	autoTxKeeper := keepers.IntentKeeper

	autoTxHistory, err := CreateFakeActionHistory(autoTxKeeper, ctx, ctx.BlockTime())
	require.NoError(t, err)

	ID := "1"
	got, err := autoTxKeeper.ActionHistory(sdk.WrapSDKContext(ctx), &types.QueryActionHistoryRequest{Id: ID})
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, got.History[0].ScheduledExecTime, autoTxHistory.History[0].ScheduledExecTime)
	require.Equal(t, got.History[0].ActualExecTime, autoTxHistory.History[0].ActualExecTime)

}

func TestQueryActionsList(t *testing.T) {
	ctx, keepers, _ := CreateTestInput(t, false)
	autoTxKeeper := keepers.IntentKeeper
	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 1000000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin("denom", 500))

	creator, _ := CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, deposit)
	var expectedActions []types.ActionInfo
	portID, err := icatypes.NewControllerPortID(creator.String())
	require.NoError(t, err)

	// create 10 actions
	for i := 0; i < 10; i++ {
		autoTx, err := CreateFakeAction(autoTxKeeper, ctx, creator, portID, ibctesting.FirstConnectionID, time.Minute, time.Hour, ctx.BlockTime(), topUp)
		require.NoError(t, err)

		expectedActions = append(expectedActions, autoTx)
	}

	got, err := autoTxKeeper.Actions(sdk.WrapSDKContext(ctx), &types.QueryActionsRequest{})

	require.NoError(t, err)
	require.NotNil(t, got)
	for i, expectedAction := range expectedActions {

		assert.Equal(t, expectedAction.GetTxMsgs(autoTxKeeper.cdc), got.ActionInfos[i].GetTxMsgs(autoTxKeeper.cdc))
		assert.Equal(t, expectedAction.ICAConfig.PortID, got.ActionInfos[i].ICAConfig.PortID)
		assert.Equal(t, expectedAction.Owner, got.ActionInfos[i].Owner)
		assert.Equal(t, expectedAction.ICAConfig.ConnectionID, got.ActionInfos[i].ICAConfig.ConnectionID)
		assert.Equal(t, expectedAction.Interval, got.ActionInfos[i].Interval)
		assert.Equal(t, expectedAction.EndTime, got.ActionInfos[i].EndTime)
		assert.Equal(t, expectedAction.Configuration, got.ActionInfos[i].Configuration)
		assert.Equal(t, expectedAction.UpdateHistory, got.ActionInfos[i].UpdateHistory)
	}
}

func TestQueryActionsListWithAuthZMsg(t *testing.T) {
	ctx, keepers, _ := CreateTestInput(t, false)
	autoTxKeeper := keepers.IntentKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 1000000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin("denom", 500))

	creator, _ := CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, deposit)

	portID, err := icatypes.NewControllerPortID(creator.String())
	require.NoError(t, err)

	expectedAction, err := CreateFakeAuthZAction(autoTxKeeper, ctx, creator, portID, ibctesting.FirstConnectionID, time.Minute, time.Hour, ctx.BlockTime(), topUp)
	require.NoError(t, err)
	got, err := autoTxKeeper.Actions(sdk.WrapSDKContext(ctx), &types.QueryActionsRequest{})

	require.NoError(t, err)
	require.NotNil(t, got)

	var txMsg sdk.Msg
	autoTxKeeper.cdc.UnpackAny(expectedAction.Msgs[0], &txMsg)

	var gotMsg sdk.Msg
	autoTxKeeper.cdc.UnpackAny(got.ActionInfos[0].Msgs[0], &gotMsg)

	assert.Equal(t, expectedAction.Msgs, got.ActionInfos[0].Msgs)
	//	assert.Equal(t, txMsg, gotMsg)
	assert.Equal(t, expectedAction.ICAConfig.PortID, got.ActionInfos[0].ICAConfig.PortID)
	assert.Equal(t, expectedAction.Owner, got.ActionInfos[0].Owner)
	assert.Equal(t, expectedAction.ICAConfig.ConnectionID, got.ActionInfos[0].ICAConfig.ConnectionID)
	assert.Equal(t, expectedAction.Interval, got.ActionInfos[0].Interval)
	assert.Equal(t, expectedAction.EndTime, got.ActionInfos[0].EndTime)
	assert.Equal(t, expectedAction.Configuration, got.ActionInfos[0].Configuration)
	assert.Equal(t, expectedAction.UpdateHistory, got.ActionInfos[0].UpdateHistory)

}

func TestQueryParams(t *testing.T) {
	ctx, keepers, _ := CreateTestInput(t, false)
	autoTxKeeper := keepers.IntentKeeper

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

func CreateFakeAction(k Keeper, ctx sdk.Context, owner sdk.AccAddress, portID, connectionId string, duration time.Duration, interval time.Duration, startAt time.Time, feeFunds sdk.Coins) (types.ActionInfo, error) {

	id := k.autoIncrementID(ctx, types.KeyLastID)
	autoTxAddress, err := k.createFeeAccount(ctx, id, owner, feeFunds)
	if err != nil {
		return types.ActionInfo{}, err
	}
	fakeMsg := banktypes.NewMsgSend(owner, autoTxAddress, feeFunds)

	anys, err := types.PackTxMsgAnys([]sdk.Msg{fakeMsg})
	if err != nil {
		return types.ActionInfo{}, err
	}

	// fakeData, _ := icatypes.SerializeCosmosTx(k.cdc, []sdk.Msg{fakeMsg})
	endTime, execTime := k.calculateTimeAndInsertQueue(ctx, startAt, duration, id, interval)
	autoTx := types.ActionInfo{
		ID:         id,
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

	k.SetActionInfo(ctx, &autoTx)
	k.addToActionOwnerIndex(ctx, owner, id)

	var newAction types.ActionInfo
	autoTxBz := k.cdc.MustMarshal(&autoTx)
	k.cdc.MustUnmarshal(autoTxBz, &newAction)
	return newAction, nil
}

func CreateFakeActionHistory(k Keeper, ctx sdk.Context, startAt time.Time) (types.ActionHistory, error) {

	entry := types.ActionHistoryEntry{
		ScheduledExecTime: startAt.Add(time.Minute),
		ActualExecTime:    startAt.Add(time.Minute).Add(time.Microsecond),
		Errors:            []string{"text"},
		Executed:          true,
	}

	autoTxHistory := types.ActionHistory{
		History: []types.ActionHistoryEntry{entry},
	}

	k.SetActionHistory(ctx, 1, &autoTxHistory)

	return autoTxHistory, nil
}

func CreateFakeAuthZAction(k Keeper, ctx sdk.Context, owner sdk.AccAddress, portID, connectionId string, duration time.Duration, interval time.Duration, startAt time.Time, feeFunds sdk.Coins) (types.ActionInfo, error) {

	id := k.autoIncrementID(ctx, types.KeyLastID)
	autoTxAddress, err := k.createFeeAccount(ctx, id, owner, feeFunds)
	if err != nil {
		return types.ActionInfo{}, err
	}
	fakeMsg := banktypes.NewMsgSend(owner, autoTxAddress, feeFunds)
	anys, err := types.PackTxMsgAnys([]sdk.Msg{fakeMsg})
	if err != nil {
		return types.ActionInfo{}, err
	}
	fakeAuthZMsg := authztypes.MsgExec{Grantee: "ICA_ADDR", Msgs: anys}

	//fakeAuthZMsg := feegranttypes.Se{Grantee: "ICA_ADDR", Msgs: anys}
	anys, err = types.PackTxMsgAnys([]sdk.Msg{&fakeAuthZMsg})
	if err != nil {
		return types.ActionInfo{}, err
	}

	// fakeData, _ := icatypes.SerializeCosmosTx(k.cdc, []sdk.Msg{fakeMsg})
	endTime, execTime := k.calculateTimeAndInsertQueue(ctx, startAt, duration, id, interval)
	autoTx := types.ActionInfo{
		ID:         id,
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
	k.SetActionInfo(ctx, &autoTx)
	k.addToActionOwnerIndex(ctx, owner, id)
	autoTxBz := k.cdc.MustMarshal(&autoTx)
	var newAction types.ActionInfo
	k.cdc.MustUnmarshal(autoTxBz, &newAction)
	return newAction, nil
}
