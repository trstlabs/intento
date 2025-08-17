package keeper

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/math"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v8/testing"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/trstlabs/intento/x/intent/types"
)

func TestQueryFlowsByOwnerList(t *testing.T) {
	ctx, keepers, _ := CreateTestInput(t, false)

	qs := NewQueryServer(keepers.IntentKeeper)

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 1000000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin("denom", 500))

	creator, _ := CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, deposit)
	var expectedFlows []types.Flow
	portID, err := icatypes.NewControllerPortID(creator.String())
	require.NoError(t, err)

	// create 10 flows
	for i := 0; i < 10; i++ {
		flow, err := CreateFakeFlow(keepers.IntentKeeper, ctx, creator, portID, ibctesting.FirstConnectionID, time.Minute, time.Hour, ctx.BlockTime(), topUp)
		require.NoError(t, err)

		expectedFlows = append(expectedFlows, flow)
	}

	specs := map[string]struct {
		srcQuery *types.QueryFlowsForOwnerRequest
		expflows []types.Flow
		expErr   error
	}{
		"query all": {
			srcQuery: &types.QueryFlowsForOwnerRequest{
				Owner: creator.String(),
			},
			expflows: expectedFlows,
			expErr:   nil,
		},
		"with pagination offset": {
			srcQuery: &types.QueryFlowsForOwnerRequest{
				Owner: creator.String(),
				Pagination: &query.PageRequest{
					Offset: 1,
				},
			},
			expflows: expectedFlows[1:],
			expErr:   nil,
		},
		"with pagination limit": {
			srcQuery: &types.QueryFlowsForOwnerRequest{
				Owner: creator.String(),
				Pagination: &query.PageRequest{
					Limit: 1,
				},
			},
			expflows: expectedFlows[0:1],
			expErr:   nil,
		},
		"nil creator": {
			srcQuery: &types.QueryFlowsForOwnerRequest{
				Pagination: &query.PageRequest{},
			},
			expflows: expectedFlows,
			expErr:   errors.New("empty address string is not allowed"),
		},
		"nil req": {
			srcQuery: nil,
			expflows: expectedFlows,
			expErr:   status.Error(codes.InvalidArgument, "empty request"),
		},
	}

	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			got, err := qs.FlowsForOwner(ctx, spec.srcQuery)

			if spec.expErr != nil {
				require.Equal(t, spec.expErr, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
			for i, expectedFlow := range spec.expflows {
				assert.Equal(t, expectedFlow.GetTxMsgs(keepers.IntentKeeper.cdc), got.Flows[i].GetTxMsgs(keepers.IntentKeeper.cdc))
				assert.Equal(t, expectedFlow.SelfHostedICA.PortID, got.Flows[i].SelfHostedICA.PortID)
				assert.Equal(t, expectedFlow.Owner, got.Flows[i].Owner)
				assert.Equal(t, expectedFlow.SelfHostedICA.ConnectionID, got.Flows[i].SelfHostedICA.ConnectionID)
				assert.Equal(t, expectedFlow.Interval, got.Flows[i].Interval)
				assert.Equal(t, expectedFlow.EndTime, got.Flows[i].EndTime)
				assert.Equal(t, expectedFlow.Configuration, got.Flows[i].Configuration)
			}
		})
	}
}

func TestQueryFlowHistory(t *testing.T) {
	ctx, keepers, _ := CreateTestInput(t, false)

	qs := NewQueryServer(keepers.IntentKeeper)
	flowHistory, err := CreateFakeFlowHistory(keepers.IntentKeeper, ctx, ctx.BlockTime())
	require.NoError(t, err)

	ID := "1"
	got, err := qs.FlowHistory(ctx, &types.QueryFlowHistoryRequest{Id: ID})
	require.NoError(t, err)
	require.NotNil(t, got)

	require.Equal(t, got.History[0].ScheduledExecTime, flowHistory.History[0].ScheduledExecTime)
	require.Equal(t, got.History[0].ActualExecTime, flowHistory.History[0].ActualExecTime)

}

func TestQueryFlowHistoryLimit(t *testing.T) {
	ctx, keepers, _ := CreateTestInput(t, false)

	qs := NewQueryServer(keepers.IntentKeeper)

	flowHistory, err := CreateFakeFlowHistory(keepers.IntentKeeper, ctx, ctx.BlockTime())
	require.NoError(t, err)

	ID := "1"
	got, err := qs.FlowHistory(ctx, &types.QueryFlowHistoryRequest{Id: ID, Pagination: &query.PageRequest{Limit: 3}})
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, got.History[0].ScheduledExecTime, flowHistory.History[0].ScheduledExecTime)
	require.Equal(t, got.History[0].ActualExecTime, flowHistory.History[0].ActualExecTime)

}

func TestQueryFlowsList(t *testing.T) {
	ctx, keepers, _ := CreateTestInput(t, false)

	qs := NewQueryServer(keepers.IntentKeeper)
	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 1000000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin("denom", 500))

	creator, _ := CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, deposit)
	var expectedFlows []types.Flow
	portID, err := icatypes.NewControllerPortID(creator.String())
	require.NoError(t, err)

	// create 10 flows
	for i := 0; i < 10; i++ {
		flow, err := CreateFakeFlow(keepers.IntentKeeper, ctx, creator, portID, ibctesting.FirstConnectionID, time.Minute, time.Hour, ctx.BlockTime(), topUp)
		require.NoError(t, err)

		expectedFlows = append(expectedFlows, flow)
	}

	got, err := qs.Flows(ctx, &types.QueryFlowsRequest{})

	require.NoError(t, err)
	require.NotNil(t, got)
	for i, expectedFlow := range expectedFlows {

		assert.Equal(t, expectedFlow.GetTxMsgs(keepers.IntentKeeper.cdc), got.Flows[i].GetTxMsgs(keepers.IntentKeeper.cdc))
		assert.Equal(t, expectedFlow.SelfHostedICA.PortID, got.Flows[i].SelfHostedICA.PortID)
		assert.Equal(t, expectedFlow.Owner, got.Flows[i].Owner)
		assert.Equal(t, expectedFlow.SelfHostedICA.ConnectionID, got.Flows[i].SelfHostedICA.ConnectionID)
		assert.Equal(t, expectedFlow.Interval, got.Flows[i].Interval)
		assert.Equal(t, expectedFlow.EndTime, got.Flows[i].EndTime)
		assert.Equal(t, expectedFlow.Configuration, got.Flows[i].Configuration)
		assert.Equal(t, expectedFlow.UpdateHistory, got.Flows[i].UpdateHistory)
	}
}

func TestQueryFlowsListWithAuthZMsg(t *testing.T) {
	ctx, keepers, _ := CreateTestInput(t, false)

	qs := NewQueryServer(keepers.IntentKeeper)
	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 1000000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin("denom", 500))

	creator, _ := CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, deposit)

	portID, err := icatypes.NewControllerPortID(creator.String())
	require.NoError(t, err)

	expectedFlow, err := CreateFakeAuthZFlow(keepers.IntentKeeper, ctx, creator, portID, ibctesting.FirstConnectionID, time.Minute, time.Hour, ctx.BlockTime(), topUp)
	require.NoError(t, err)
	got, err := qs.Flows(ctx, &types.QueryFlowsRequest{})

	require.NoError(t, err)
	require.NotNil(t, got)

	var txMsg sdk.Msg
	_ = keepers.IntentKeeper.cdc.UnpackAny(expectedFlow.Msgs[0], &txMsg)

	var gotMsg sdk.Msg
	_ = keepers.IntentKeeper.cdc.UnpackAny(got.Flows[0].Msgs[0], &gotMsg)

	assert.Equal(t, expectedFlow.Msgs, got.Flows[0].Msgs)
	//	assert.Equal(t, txMsg, gotMsg)
	assert.Equal(t, expectedFlow.SelfHostedICA.PortID, got.Flows[0].SelfHostedICA.PortID)
	assert.Equal(t, expectedFlow.Owner, got.Flows[0].Owner)
	assert.Equal(t, expectedFlow.SelfHostedICA.ConnectionID, got.Flows[0].SelfHostedICA.ConnectionID)
	assert.Equal(t, expectedFlow.Interval, got.Flows[0].Interval)
	assert.Equal(t, expectedFlow.EndTime, got.Flows[0].EndTime)
	assert.Equal(t, expectedFlow.Configuration, got.Flows[0].Configuration)
	assert.Equal(t, expectedFlow.UpdateHistory, got.Flows[0].UpdateHistory)

}

func TestQueryParams(t *testing.T) {
	ctx, keepers, _ := CreateTestInput(t, false)

	qs := NewQueryServer(keepers.IntentKeeper)

	resp, err := qs.Params(ctx, &types.QueryParamsRequest{})
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

func TestQueryTrustlessAgentsByFeeAdmin(t *testing.T) {
	ctx, keepers, _ := CreateTestInput(t, false)

	qs := NewQueryServer(keepers.IntentKeeper)

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 1000000))

	creator, _ := CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, deposit)
	var expectedTrustlessAgents []types.TrustlessAgent
	portID, err := icatypes.NewControllerPortID(creator.String())
	require.NoError(t, err)

	// create 10
	for i := 0; i < 10; i++ {
		hostAcc, err := CreateFakeTrustlessAgent(keepers.IntentKeeper, ctx, creator.String(), portID, ibctesting.FirstConnectionID+(fmt.Sprint(i)), ibctesting.FirstConnectionID)
		require.NoError(t, err)

		expectedTrustlessAgents = append(expectedTrustlessAgents, hostAcc)
	}

	specs := map[string]struct {
		srcQuery           *types.QueryTrustlessAgentsByFeeAdminRequest
		expTrustlessAgents []types.TrustlessAgent
		expErr             error
	}{
		"query all": {
			srcQuery: &types.QueryTrustlessAgentsByFeeAdminRequest{
				FeeAdmin: creator.String(),
			},
			expTrustlessAgents: expectedTrustlessAgents,
			expErr:             nil,
		},
		"with pagination offset": {
			srcQuery: &types.QueryTrustlessAgentsByFeeAdminRequest{
				FeeAdmin: creator.String(),
				Pagination: &query.PageRequest{
					Offset: 1,
				},
			},
			expTrustlessAgents: expectedTrustlessAgents[1:],
			expErr:             nil,
		},
		"with pagination limit": {
			srcQuery: &types.QueryTrustlessAgentsByFeeAdminRequest{
				FeeAdmin: creator.String(),
				Pagination: &query.PageRequest{
					Limit: 1,
				},
			},
			expTrustlessAgents: expectedTrustlessAgents[0:1],
			expErr:             nil,
		},
		"nil admin": {
			srcQuery: &types.QueryTrustlessAgentsByFeeAdminRequest{
				Pagination: &query.PageRequest{},
			},
			expTrustlessAgents: expectedTrustlessAgents,
			expErr:             errors.New("empty address string is not allowed"),
		},
		"nil req": {
			srcQuery:           nil,
			expTrustlessAgents: expectedTrustlessAgents,
			expErr:             status.Error(codes.InvalidArgument, "empty request"),
		},
	}

	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			got, err := qs.TrustlessAgentsByFeeAdmin(ctx, spec.srcQuery)
			if spec.expErr != nil {
				require.Equal(t, spec.expErr, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
			for i, trustlessAgent := range spec.expTrustlessAgents {
				assert.Equal(t, trustlessAgent.FeeConfig.FeeAdmin, got.TrustlessAgents[i].FeeConfig.FeeAdmin)

			}
		})
	}
}

func CreateFakeFlow(k Keeper, ctx sdk.Context, owner sdk.AccAddress, portID, connectionId string, duration time.Duration, interval time.Duration, startAt time.Time, feeFunds sdk.Coins) (types.Flow, error) {

	id := k.autoIncrementID(ctx, types.KeyLastID)
	flowAddress, err := k.createFeeAccount(ctx, id, owner, feeFunds)
	if err != nil {
		return types.Flow{}, err
	}
	fakeMsg := banktypes.NewMsgSend(owner, flowAddress, feeFunds)

	anys, err := types.PackTxMsgAnys([]sdk.Msg{fakeMsg})
	if err != nil {
		return types.Flow{}, err
	}

	// fakeData, _ := icatypes.SerializeCosmosTx(k.cdc, []sdk.Msg{fakeMsg})
	endTime, execTime := k.calculateTimeAndInsertQueue(ctx, startAt, duration, id, interval)
	flow := types.Flow{
		ID:         id,
		FeeAddress: flowAddress.String(),
		Owner:      owner.String(),
		// Data:       fakeData,
		Msgs:     anys,
		Interval: interval,

		StartTime:     startAt,
		ExecTime:      execTime,
		EndTime:       endTime,
		SelfHostedICA: &types.ICAConfig{PortID: portID},
		Configuration: &types.ExecutionConfiguration{SaveResponses: true},
	}

	k.SetFlow(ctx, &flow)
	k.addToFlowOwnerIndex(ctx, owner, id)

	var newFlow types.Flow
	flowBz := k.cdc.MustMarshal(&flow)
	k.cdc.MustUnmarshal(flowBz, &newFlow)
	return newFlow, nil
}

func CreateFakeFlowHistory(k Keeper, ctx sdk.Context, startAt time.Time) (types.FlowHistory, error) {

	// Create an empty FlowHistory with a pre-allocated slice for efficiency
	flowHistory := types.FlowHistory{
		History: make([]types.FlowHistoryEntry, 0, 10), // Pre-allocate space for 10 entries
	}

	// Loop to create and append 10 entries
	for i := 0; i < 10; i++ {
		entry := types.FlowHistoryEntry{
			ScheduledExecTime: startAt.Add(time.Duration(i) * time.Minute),
			ActualExecTime:    startAt.Add(time.Duration(i) * time.Minute).Add(time.Microsecond),
			Errors:            []string{"error text"}, // Example error text
			Executed:          true,
		}

		k.SetFlowHistoryEntry(ctx, 1, &entry)
		flowHistory.History = append(flowHistory.History, entry)
	}

	return flowHistory, nil
}

func CreateFakeAuthZFlow(k Keeper, ctx sdk.Context, owner sdk.AccAddress, portID, connectionId string, duration time.Duration, interval time.Duration, startAt time.Time, feeFunds sdk.Coins) (types.Flow, error) {

	id := k.autoIncrementID(ctx, types.KeyLastID)
	flowAddress, err := k.createFeeAccount(ctx, id, owner, feeFunds)
	if err != nil {
		return types.Flow{}, err
	}
	fakeMsg := banktypes.NewMsgSend(owner, flowAddress, feeFunds)
	anys, err := types.PackTxMsgAnys([]sdk.Msg{fakeMsg})
	if err != nil {
		return types.Flow{}, err
	}
	fakeAuthZMsg := authztypes.MsgExec{Grantee: "ICA_ADDR", Msgs: anys}

	//fakeAuthZMsg := feegranttypes.Se{Grantee: "ICA_ADDR", Msgs: anys}
	anys, err = types.PackTxMsgAnys([]sdk.Msg{&fakeAuthZMsg})
	if err != nil {
		return types.Flow{}, err
	}

	// fakeData, _ := icatypes.SerializeCosmosTx(k.cdc, []sdk.Msg{fakeMsg})
	endTime, execTime := k.calculateTimeAndInsertQueue(ctx, startAt, duration, id, interval)
	flow := types.Flow{
		ID:         id,
		FeeAddress: flowAddress.String(),
		Owner:      owner.String(),
		// Data:       fakeData,
		Msgs:          anys,
		Interval:      interval,
		UpdateHistory: nil,
		StartTime:     startAt,
		ExecTime:      execTime,
		EndTime:       endTime,
		SelfHostedICA: &types.ICAConfig{PortID: portID},
	}
	k.SetFlow(ctx, &flow)
	k.addToFlowOwnerIndex(ctx, owner, id)
	flowBz := k.cdc.MustMarshal(&flow)
	var newFlow types.Flow
	k.cdc.MustUnmarshal(flowBz, &newFlow)
	return newFlow, nil
}

func CreateFakeTrustlessAgent(k Keeper, ctx sdk.Context, creator, portID, connectionId, hostConnectionId string) (types.TrustlessAgent, error) {
	agentAddress, err := DeriveAgentAddress(creator, connectionId)
	if err != nil {
		return types.TrustlessAgent{}, err
	}

	creatorAddr, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return types.TrustlessAgent{}, err
	}
	trustlessAgent := types.TrustlessAgent{AgentAddress: agentAddress.String(), FeeConfig: &types.TrustlessAgentFeeConfig{FeeAdmin: creator, FeeCoinsSupported: sdk.NewCoins(sdk.NewCoin(types.Denom, math.NewInt(1)))}, ICAConfig: &types.ICAConfig{ConnectionID: connectionId, PortID: portID}}
	//store hosted config by address on hosted key prefix
	k.SetTrustlessAgent(ctx, &trustlessAgent)
	k.addToTrustlessAgentAdminIndex(ctx, creatorAddr, agentAddress.String())

	return trustlessAgent, nil
}
