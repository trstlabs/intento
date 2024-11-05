package keeper

import (
	"errors"
	"fmt"
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
	intentKeeper := keepers.IntentKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 1000000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin("denom", 500))

	creator, _ := CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, deposit)
	var expectedActions []types.ActionInfo
	portID, err := icatypes.NewControllerPortID(creator.String())
	require.NoError(t, err)

	// create 10 actions
	for i := 0; i < 10; i++ {
		action, err := CreateFakeAction(intentKeeper, ctx, creator, portID, ibctesting.FirstConnectionID, time.Minute, time.Hour, ctx.BlockTime(), topUp)
		require.NoError(t, err)

		expectedActions = append(expectedActions, action)
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
			got, err := intentKeeper.ActionsForOwner(sdk.WrapSDKContext(ctx), spec.srcQuery)

			if spec.expErr != nil {
				require.Equal(t, spec.expErr, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
			for i, expectedAction := range spec.expActionInfos {
				assert.Equal(t, expectedAction.GetTxMsgs(intentKeeper.cdc), got.ActionInfos[i].GetTxMsgs(intentKeeper.cdc))
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
	intentKeeper := keepers.IntentKeeper

	actionHistory, err := CreateFakeActionHistory(intentKeeper, ctx, ctx.BlockTime())
	require.NoError(t, err)

	ID := "1"
	got, err := intentKeeper.ActionHistory(sdk.WrapSDKContext(ctx), &types.QueryActionHistoryRequest{Id: ID})
	require.NoError(t, err)
	require.NotNil(t, got)

	require.Equal(t, got.History[0].ScheduledExecTime, actionHistory.History[0].ScheduledExecTime)
	require.Equal(t, got.History[0].ActualExecTime, actionHistory.History[0].ActualExecTime)

}

func TestQueryActionHistoryLimit(t *testing.T) {
	ctx, keepers, _ := CreateTestInput(t, false)
	intentKeeper := keepers.IntentKeeper

	actionHistory, err := CreateFakeActionHistory(intentKeeper, ctx, ctx.BlockTime())
	require.NoError(t, err)

	ID := "1"
	got, err := intentKeeper.ActionHistory(sdk.WrapSDKContext(ctx), &types.QueryActionHistoryRequest{Id: ID, Pagination: &query.PageRequest{Limit: 3}})
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, got.History[0].ScheduledExecTime, actionHistory.History[0].ScheduledExecTime)
	require.Equal(t, got.History[0].ActualExecTime, actionHistory.History[0].ActualExecTime)

}

func TestQueryActionsList(t *testing.T) {
	ctx, keepers, _ := CreateTestInput(t, false)
	intentKeeper := keepers.IntentKeeper
	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 1000000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin("denom", 500))

	creator, _ := CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, deposit)
	var expectedActions []types.ActionInfo
	portID, err := icatypes.NewControllerPortID(creator.String())
	require.NoError(t, err)

	// create 10 actions
	for i := 0; i < 10; i++ {
		action, err := CreateFakeAction(intentKeeper, ctx, creator, portID, ibctesting.FirstConnectionID, time.Minute, time.Hour, ctx.BlockTime(), topUp)
		require.NoError(t, err)

		expectedActions = append(expectedActions, action)
	}

	got, err := intentKeeper.Actions(sdk.WrapSDKContext(ctx), &types.QueryActionsRequest{})

	require.NoError(t, err)
	require.NotNil(t, got)
	for i, expectedAction := range expectedActions {

		assert.Equal(t, expectedAction.GetTxMsgs(intentKeeper.cdc), got.ActionInfos[i].GetTxMsgs(intentKeeper.cdc))
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
	intentKeeper := keepers.IntentKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 1000000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin("denom", 500))

	creator, _ := CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, deposit)

	portID, err := icatypes.NewControllerPortID(creator.String())
	require.NoError(t, err)

	expectedAction, err := CreateFakeAuthZAction(intentKeeper, ctx, creator, portID, ibctesting.FirstConnectionID, time.Minute, time.Hour, ctx.BlockTime(), topUp)
	require.NoError(t, err)
	got, err := intentKeeper.Actions(sdk.WrapSDKContext(ctx), &types.QueryActionsRequest{})

	require.NoError(t, err)
	require.NotNil(t, got)

	var txMsg sdk.Msg
	_ = intentKeeper.cdc.UnpackAny(expectedAction.Msgs[0], &txMsg)

	var gotMsg sdk.Msg
	_ = intentKeeper.cdc.UnpackAny(got.ActionInfos[0].Msgs[0], &gotMsg)

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
	intentKeeper := keepers.IntentKeeper

	resp, err := intentKeeper.Params(sdk.WrapSDKContext(ctx), &types.QueryParamsRequest{})
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

func TestQueryHostedAccountsByAdmin(t *testing.T) {
	ctx, keepers, _ := CreateTestInput(t, false)
	intentKeeper := keepers.IntentKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 1000000))

	creator, _ := CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, deposit)
	var expectedHostedAccounts []types.HostedAccount
	portID, err := icatypes.NewControllerPortID(creator.String())
	require.NoError(t, err)

	// create 10
	for i := 0; i < 10; i++ {
		hostAcc, err := CreateFakeHostedAcc(intentKeeper, ctx, creator.String(), portID, ibctesting.FirstConnectionID+(fmt.Sprint(i)), ibctesting.FirstConnectionID)
		require.NoError(t, err)

		expectedHostedAccounts = append(expectedHostedAccounts, hostAcc)
	}

	specs := map[string]struct {
		srcQuery          *types.QueryHostedAccountsByAdminRequest
		expHostedAccounts []types.HostedAccount
		expErr            error
	}{
		"query all": {
			srcQuery: &types.QueryHostedAccountsByAdminRequest{
				Admin: creator.String(),
			},
			expHostedAccounts: expectedHostedAccounts,
			expErr:            nil,
		},
		"with pagination offset": {
			srcQuery: &types.QueryHostedAccountsByAdminRequest{
				Admin: creator.String(),
				Pagination: &query.PageRequest{
					Offset: 1,
				},
			},
			expHostedAccounts: expectedHostedAccounts[1:],
			expErr:            nil,
		},
		"with pagination limit": {
			srcQuery: &types.QueryHostedAccountsByAdminRequest{
				Admin: creator.String(),
				Pagination: &query.PageRequest{
					Limit: 1,
				},
			},
			expHostedAccounts: expectedHostedAccounts[0:1],
			expErr:            nil,
		},
		"nil admin": {
			srcQuery: &types.QueryHostedAccountsByAdminRequest{
				Pagination: &query.PageRequest{},
			},
			expHostedAccounts: expectedHostedAccounts,
			expErr:            errors.New("empty address string is not allowed"),
		},
		"nil req": {
			srcQuery:          nil,
			expHostedAccounts: expectedHostedAccounts,
			expErr:            status.Error(codes.InvalidArgument, "empty request"),
		},
	}

	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			got, err := intentKeeper.HostedAccountsByAdmin(sdk.WrapSDKContext(ctx), spec.srcQuery)
			//fmt.Println(got)
			if spec.expErr != nil {
				require.Equal(t, spec.expErr, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
			for i, hostedAccount := range spec.expHostedAccounts {
				assert.Equal(t, hostedAccount.HostFeeConfig.Admin, got.HostedAccounts[i].HostFeeConfig.Admin)

			}
		})
	}
}

func CreateFakeAction(k Keeper, ctx sdk.Context, owner sdk.AccAddress, portID, connectionId string, duration time.Duration, interval time.Duration, startAt time.Time, feeFunds sdk.Coins) (types.ActionInfo, error) {

	id := k.autoIncrementID(ctx, types.KeyLastID)
	actionAddress, err := k.createFeeAccount(ctx, id, owner, feeFunds)
	if err != nil {
		return types.ActionInfo{}, err
	}
	fakeMsg := banktypes.NewMsgSend(owner, actionAddress, feeFunds)

	anys, err := types.PackTxMsgAnys([]sdk.Msg{fakeMsg})
	if err != nil {
		return types.ActionInfo{}, err
	}

	// fakeData, _ := icatypes.SerializeCosmosTx(k.cdc, []sdk.Msg{fakeMsg})
	endTime, execTime := k.calculateTimeAndInsertQueue(ctx, startAt, duration, id, interval)
	action := types.ActionInfo{
		ID:         id,
		FeeAddress: actionAddress.String(),
		Owner:      owner.String(),
		// Data:       fakeData,
		Msgs:     anys,
		Interval: interval,

		StartTime:     startAt,
		ExecTime:      execTime,
		EndTime:       endTime,
		ICAConfig:     &types.ICAConfig{PortID: portID},
		Configuration: &types.ExecutionConfiguration{SaveResponses: true},
	}

	k.SetActionInfo(ctx, &action)
	k.addToActionOwnerIndex(ctx, owner, id)

	var newAction types.ActionInfo
	actionBz := k.cdc.MustMarshal(&action)
	k.cdc.MustUnmarshal(actionBz, &newAction)
	return newAction, nil
}

func CreateFakeActionHistory(k Keeper, ctx sdk.Context, startAt time.Time) (types.ActionHistory, error) {

	// Create an empty ActionHistory with a pre-allocated slice for efficiency
	actionHistory := types.ActionHistory{
		History: make([]types.ActionHistoryEntry, 0, 10), // Pre-allocate space for 10 entries
	}

	// Loop to create and append 10 entries
	for i := 0; i < 10; i++ {
		entry := types.ActionHistoryEntry{
			ScheduledExecTime: startAt.Add(time.Duration(i) * time.Minute),
			ActualExecTime:    startAt.Add(time.Duration(i) * time.Minute).Add(time.Microsecond),
			Errors:            []string{"error text"}, // Example error text
			Executed:          true,
		}

		k.SetActionHistoryEntry(ctx, 1, &entry)
		actionHistory.History = append(actionHistory.History, entry)
	}

	return actionHistory, nil
}

func CreateFakeAuthZAction(k Keeper, ctx sdk.Context, owner sdk.AccAddress, portID, connectionId string, duration time.Duration, interval time.Duration, startAt time.Time, feeFunds sdk.Coins) (types.ActionInfo, error) {

	id := k.autoIncrementID(ctx, types.KeyLastID)
	actionAddress, err := k.createFeeAccount(ctx, id, owner, feeFunds)
	if err != nil {
		return types.ActionInfo{}, err
	}
	fakeMsg := banktypes.NewMsgSend(owner, actionAddress, feeFunds)
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
	action := types.ActionInfo{
		ID:         id,
		FeeAddress: actionAddress.String(),
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
	k.SetActionInfo(ctx, &action)
	k.addToActionOwnerIndex(ctx, owner, id)
	actionBz := k.cdc.MustMarshal(&action)
	var newAction types.ActionInfo
	k.cdc.MustUnmarshal(actionBz, &newAction)
	return newAction, nil
}

func CreateFakeHostedAcc(k Keeper, ctx sdk.Context, creator, portID, connectionId, hostConnectionId string) (types.HostedAccount, error) {
	hostedAddress, err := DeriveHostedAddress(creator, connectionId)
	if err != nil {
		return types.HostedAccount{}, err
	}

	creatorAddr, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return types.HostedAccount{}, err
	}
	hostedAcc := types.HostedAccount{HostedAddress: hostedAddress.String(), HostFeeConfig: &types.HostFeeConfig{Admin: creator, FeeCoinsSuported: sdk.NewCoins(sdk.NewCoin(types.Denom, sdk.NewInt(1)))}, ICAConfig: &types.ICAConfig{ConnectionID: connectionId, HostConnectionID: hostConnectionId, PortID: portID}}
	//store hosted config by address on hosted key prefix
	k.SetHostedAccount(ctx, &hostedAcc)
	k.addToHostedAccountAdminIndex(ctx, creatorAddr, hostedAddress.String())

	return hostedAcc, nil
}
