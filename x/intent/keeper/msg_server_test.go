package keeper_test

import (
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"
	"github.com/trstlabs/intento/x/intent/keeper"
	"github.com/trstlabs/intento/x/intent/types"
)

func (suite *KeeperTestSuite) TestRegisterInterchainAccount() {
	var (
		owner string
		path  *ibctesting.Path
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"success", func() {}, true,
		},
		{
			"port is already bound",
			func() {
				GetICAApp(suite.chainA.TestChain).GetIBCKeeper().PortKeeper.BindPort(suite.chainA.GetContext(), TestPortID)
			},
			false,
		},
		{
			"fails to generate port-id",
			func() {
				owner = ""
			},
			false,
		},
		{
			"MsgChanOpenInit fails - channel is already active",
			func() {
				portID, err := icatypes.NewControllerPortID(owner)
				suite.Require().NoError(err)

				channel := channeltypes.NewChannel(
					channeltypes.OPEN,
					channeltypes.ORDERED,
					channeltypes.NewCounterparty(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID),
					[]string{path.EndpointA.ConnectionID},
					path.EndpointA.ChannelConfig.Version,
				)
				GetICAApp(suite.chainA.TestChain).IBCKeeper.ChannelKeeper.SetChannel(suite.chainA.GetContext(), portID, ibctesting.FirstChannelID, channel)

				GetICAApp(suite.chainA.TestChain).ICAControllerKeeper.SetActiveChannelID(suite.chainA.GetContext(), ibctesting.FirstConnectionID, portID, ibctesting.FirstChannelID)
			},
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest()

			owner = TestOwnerAddress // must be explicitly changed

			path = NewICAPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupConnections(path)

			tc.malleate() // malleate mutates test data

			msgSrv := keeper.NewMsgServerImpl(GetActionKeeper(suite.chainA))
			msg := types.NewMsgRegisterAccount(owner, path.EndpointA.ConnectionID, path.EndpointA.ChannelConfig.Version)

			res, err := msgSrv.RegisterAccount(sdk.WrapSDKContext(suite.chainA.GetContext()), msg)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestSubmitTx() {
	var (
		path                      *ibctesting.Path
		registerInterchainAccount bool
		owner                     string
		connectionId              string
		sdkMsg                    sdk.Msg
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"success", func() {
				registerInterchainAccount = true
				owner = TestOwnerAddress
				connectionId = path.EndpointA.ConnectionID
			}, true,
		},
		{
			"failure - owner address is empty", func() {
				registerInterchainAccount = true
				owner = ""
				connectionId = path.EndpointA.ConnectionID

			}, false,
		},
		{
			"failure - active channel does not exist for connection ID", func() {
				registerInterchainAccount = true
				owner = TestOwnerAddress
				connectionId = "connection-100"

			}, false,
		},
		{
			"failure - active channel does not exist for port ID", func() {
				registerInterchainAccount = true
				owner = "cosmos153lf4zntqt33a4v0sm5cytrxyqn78q7kz8j8x5"
				connectionId = path.EndpointA.ConnectionID

			}, false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest()

			icaAppA := GetICAApp(suite.chainA.TestChain)
			icaAppB := GetICAApp(suite.chainB.TestChain)

			path = NewICAPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupConnections(path)

			tc.malleate() // malleate mutates test data

			if registerInterchainAccount {
				err := SetupICAPath(path, TestOwnerAddress)
				suite.Require().NoError(err)

				portID, err := icatypes.NewControllerPortID(TestOwnerAddress)
				suite.Require().NoError(err)

				// Get the address of the interchain account stored in state during handshake step
				interchainAccountAddr, found := GetICAApp(suite.chainA.TestChain).ICAControllerKeeper.GetInterchainAccountAddress(suite.chainA.GetContext(), path.EndpointA.ConnectionID, portID)
				suite.Require().True(found)

				icaAddr, err := sdk.AccAddressFromBech32(interchainAccountAddr)
				suite.Require().NoError(err)

				// Check if account is created
				interchainAccount := icaAppB.AccountKeeper.GetAccount(suite.chainB.GetContext(), icaAddr)
				suite.Require().Equal(interchainAccount.GetAddress().String(), interchainAccountAddr)

				// Create bank transfer message to execute on the host
				sdkMsg = &banktypes.MsgSend{
					FromAddress: interchainAccountAddr,
					ToAddress:   suite.chainB.SenderAccount.GetAddress().String(),
					Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
				}
			}

			msgSrv := keeper.NewMsgServerImpl(GetActionKeeperFromApp(icaAppA))
			msg, err := types.NewMsgSubmitTx(owner, sdkMsg, connectionId)
			suite.Require().NoError(err)

			res, err := msgSrv.SubmitTx(sdk.WrapSDKContext(suite.chainA.GetContext()), msg)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestSubmitAction() {
	var (
		path                      *ibctesting.Path
		registerInterchainAccount bool
		noOwner                   bool
		connectionId              string
		hostConnectionId          string
		sdkMsg                    sdk.Msg
		parseIcaAddress           bool
		startAtBeforeBlockHeight  bool
		transferMsg               bool
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"success - IBC ICA action", func() {
				registerInterchainAccount = true
				connectionId = path.EndpointA.ConnectionID
				hostConnectionId = path.EndpointB.ConnectionID
			}, true,
		},
		{
			"success - local action", func() {
				registerInterchainAccount = false
				connectionId = ""
				hostConnectionId = ""
				sdkMsg = &banktypes.MsgSend{
					FromAddress: suite.chainA.SenderAccount.GetAddress().String(),
					ToAddress:   TestOwnerAddress,
					Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
				}
			}, true,
		},
		{
			"success - IBC transfer", func() {
				registerInterchainAccount = false
				connectionId = ""
				hostConnectionId = ""
				transferMsg = true
			}, true,
		},
		{
			"success - parse ICA address", func() {
				registerInterchainAccount = true
				connectionId = path.EndpointA.ConnectionID
				hostConnectionId = path.EndpointB.ConnectionID
				parseIcaAddress = true
				transferMsg = false
			}, true,
		},
		{
			"failure - start before block height", func() {
				registerInterchainAccount = true
				noOwner = false
				connectionId = path.EndpointA.ConnectionID
				hostConnectionId = path.EndpointB.ConnectionID
				startAtBeforeBlockHeight = true
			}, false,
		},
		{
			"failure - owner address is empty", func() {
				registerInterchainAccount = false
				noOwner = true
				connectionId = path.EndpointA.ConnectionID
				hostConnectionId = path.EndpointB.ConnectionID
				parseIcaAddress = false
			}, false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest()
			var owner string

			icaAppA := GetICAApp(suite.chainA.TestChain)
			icaAppB := GetICAApp(suite.chainB.TestChain)
			path = NewICAPath(suite.chainA, suite.chainB)

			suite.coordinator.SetupConnections(path)

			tc.malleate() // malleate mutates test data

			if transferMsg {
				path = NewTransferPath(suite.chainA, suite.chainB)
				suite.coordinator.SetupConnections(path)
				sdkMsg = transfertypes.NewMsgTransfer(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100)), suite.chainA.SenderAccount.GetAddress().String(), suite.chainB.SenderAccount.GetAddress().String(), suite.chainB.GetTimeoutHeight(), 0, "")
			}

			if noOwner {
				owner = ""
			} else {
				owner = suite.chainA.SenderAccount.GetAddress().String()
			}
			if registerInterchainAccount {
				err := SetupICAPath(path, owner)
				suite.Require().NoError(err)

				portID, err := icatypes.NewControllerPortID(owner)
				suite.Require().NoError(err)

				// Get the address of the interchain account stored in state during handshake step
				interchainAccountAddr, found := GetICAApp(suite.chainA.TestChain).ICAControllerKeeper.GetInterchainAccountAddress(suite.chainA.GetContext(), path.EndpointA.ConnectionID, portID)
				suite.Require().True(found)

				icaAddr, err := sdk.AccAddressFromBech32(interchainAccountAddr)
				suite.Require().NoError(err)

				// Check if account is created
				interchainAccount := icaAppB.AccountKeeper.GetAccount(suite.chainB.GetContext(), icaAddr)
				suite.Require().Equal(interchainAccount.GetAddress().String(), interchainAccountAddr)
				if parseIcaAddress {
					interchainAccountAddr = types.ParseICAValue
				}
				// Create bank transfer message to execute on the host
				sdkMsg = &banktypes.MsgSend{
					FromAddress: interchainAccountAddr,
					ToAddress:   suite.chainB.SenderAccount.GetAddress().String(),
					Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
				}
			}

			label := "label"
			duration := time.Second * 200
			durationTimeText := duration.String()
			interval := time.Second * 100
			intervalTimeText := interval.String()
			startAt := uint64(0)
			if startAtBeforeBlockHeight {
				startAt = uint64(suite.chainA.GetContext().BlockTime().Unix() - 60*60)
			}
			msg, err := types.NewMsgSubmitAction(owner, label, []sdk.Msg{sdkMsg}, connectionId, hostConnectionId, durationTimeText, intervalTimeText, startAt, sdk.Coins{}, &types.ExecutionConfiguration{FallbackToOwnerBalance: true})
			suite.Require().NoError(err)
			wrappedCtx := sdk.WrapSDKContext(suite.chainA.GetContext())

			msgSrv := keeper.NewMsgServerImpl(GetActionKeeperFromApp(icaAppA))
			res, err := msgSrv.SubmitAction(wrappedCtx, msg)

			if !tc.expPass {
				suite.Require().Error(err)
				suite.Require().Nil(res)
				return
			}

			suite.Require().NoError(err)
			suite.Require().NotNil(res)

			if parseIcaAddress {
				err := sdkMsg.ValidateBasic()
				suite.Require().Contains(err.Error(), "bech32")
				err = msg.ValidateBasic()
				suite.Require().NoError(err)
			}
			actionKeeper := icaAppA.IntentKeeper

			suite.chainA.CurrentHeader.Time = suite.chainA.CurrentHeader.Time.Add(interval)
			types.Denom = "stake"
			FakeBeginBlocker(suite.chainA.GetContext(), actionKeeper, sdk.ConsAddress(suite.chainA.Vals.Proposer.Address))
			suite.chainA.NextBlock()
			action := actionKeeper.GetActionInfo(suite.chainA.GetContext(), 1)
			actionHistory, err := actionKeeper.GetActionHistory(suite.chainA.GetContext(), 1)

			suite.Require().NoError(err)
			suite.Require().NotEqual(action, types.ActionInfo{})
			suite.Require().Equal(action.Owner, owner)
			suite.Require().Equal(action.Label, label)

			//ibc

			if action.ICAConfig.PortID != "" {
				suite.Require().Equal(action.ICAConfig.PortID, "icacontroller-"+owner)
			}
			if !transferMsg {
				suite.Require().Empty(actionHistory.History[0].Errors)

			}

			if parseIcaAddress {
				var txMsg sdk.Msg
				err = icaAppA.AppCodec().UnpackAny(action.Msgs[0], &txMsg)
				suite.Require().NoError(err)
				index := strings.Index(txMsg.String(), types.ParseICAValue)
				suite.Require().Equal(-1, index)

			}

		})
	}
}

func (suite *KeeperTestSuite) TestSubmitActionAuthZ() {
	var (
		path                      *ibctesting.Path
		registerInterchainAccount bool
		connectionId              string
		hostConnectionId          string
		sdkMsg                    sdk.Msg
		parseIcaAddress           bool
		startAtBeforeBlockHeight  bool
	)
	sdkMsg = &banktypes.MsgSend{
		FromAddress: TestOwnerAddress,
		ToAddress:   suite.chainA.SenderAccount.GetAddress().String(),
		Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
	}
	anyMsg, _ := types.PackTxMsgAnys([]sdk.Msg{sdkMsg})

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{{

		"fail - MsgExec action", func() {
			registerInterchainAccount = true
			connectionId = ""
			hostConnectionId = ""
			sdkMsg = &authztypes.MsgExec{Grantee: TestOwnerAddress,
				Msgs: anyMsg,
			}
		}, false,
	},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest()
			var owner string

			icaAppA := GetICAApp(suite.chainA.TestChain)
			icaAppB := GetICAApp(suite.chainB.TestChain)
			path = NewICAPath(suite.chainA, suite.chainB)

			suite.coordinator.SetupConnections(path)

			tc.malleate() // malleate mutates test data

			owner = suite.chainA.SenderAccount.GetAddress().String()

			if registerInterchainAccount {
				err := SetupICAPath(path, owner)
				suite.Require().NoError(err)

				portID, err := icatypes.NewControllerPortID(owner)
				suite.Require().NoError(err)

				// Get the address of the interchain account stored in state during handshake step
				interchainAccountAddr, found := GetICAApp(suite.chainA.TestChain).ICAControllerKeeper.GetInterchainAccountAddress(suite.chainA.GetContext(), path.EndpointA.ConnectionID, portID)
				suite.Require().True(found)

				icaAddr, err := sdk.AccAddressFromBech32(interchainAccountAddr)
				suite.Require().NoError(err)

				// Check if account is created
				interchainAccount := icaAppB.AccountKeeper.GetAccount(suite.chainB.GetContext(), icaAddr)
				suite.Require().Equal(interchainAccount.GetAddress().String(), interchainAccountAddr)
				if parseIcaAddress {
					interchainAccountAddr = types.ParseICAValue
				}

			}

			label := "label"
			duration := time.Second * 200
			durationTimeText := duration.String()
			interval := time.Second * 100
			intervalTimeText := interval.String()
			startAt := uint64(0)
			if startAtBeforeBlockHeight {
				startAt = uint64(suite.chainA.GetContext().BlockTime().Unix() - 60*60)
			}
			msg, err := types.NewMsgSubmitAction(owner, label, []sdk.Msg{sdkMsg}, connectionId, hostConnectionId, durationTimeText, intervalTimeText, startAt, sdk.Coins{}, &types.ExecutionConfiguration{FallbackToOwnerBalance: true})
			suite.Require().NoError(err)
			wrappedCtx := sdk.WrapSDKContext(suite.chainA.GetContext())

			msgSrv := keeper.NewMsgServerImpl(GetActionKeeperFromApp(icaAppA))
			res, err := msgSrv.SubmitAction(wrappedCtx, msg)

			if !tc.expPass {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), "exec msg signer error")
				suite.Require().Nil(res)
				return
			}
		})
	}
}

func (suite *KeeperTestSuite) TestUpdateAction() {
	var (
		path                      *ibctesting.Path
		registerInterchainAccount bool
		owner                     string
		connectionId              string
		sdkMsg                    sdk.Msg
		newEndTime                uint64
		newStartAt                uint64
		newInterval               string
		addFakeExecHistory        bool
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"success - update action", func() {
				registerInterchainAccount = true
				owner = TestOwnerAddress
				connectionId = path.EndpointA.ConnectionID
				newStartAt = uint64(time.Now().Unix())
				newEndTime = uint64(time.Now().Add(time.Hour).Unix())
				newInterval = "8m20s"
			}, true,
		},
		{
			"success - update local action", func() {
				registerInterchainAccount = false
				owner = TestOwnerAddress
				connectionId = ""
			}, true,
		},
		{
			"failure - interval shorter than min duration", func() {
				registerInterchainAccount = true
				owner = TestOwnerAddress
				connectionId = path.EndpointA.ConnectionID
				newInterval = "1s"
			}, false,
		},
		{
			"failure - start time can not be changed after execution entry", func() {
				registerInterchainAccount = true
				owner = TestOwnerAddress
				connectionId = path.EndpointA.ConnectionID
				newStartAt = uint64(time.Date(2020, time.April, 22, 12, 0, 0, 0, time.UTC).Unix())
				newInterval = "8m20s"
				addFakeExecHistory = true
			}, false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest()

			icaAppA := GetICAApp(suite.chainA.TestChain)
			icaAppB := GetICAApp(suite.chainB.TestChain)

			path = NewICAPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupConnections(path)

			tc.malleate() // malleate mutates test data

			if registerInterchainAccount {
				err := SetupICAPath(path, TestOwnerAddress)
				suite.Require().NoError(err)

				portID, err := icatypes.NewControllerPortID(TestOwnerAddress)
				suite.Require().NoError(err)

				// Get the address of the interchain account stored in state during handshake step
				interchainAccountAddr, found := GetICAApp(suite.chainA.TestChain).ICAControllerKeeper.GetInterchainAccountAddress(suite.chainA.GetContext(), path.EndpointA.ConnectionID, portID)
				suite.Require().True(found)

				icaAddr, err := sdk.AccAddressFromBech32(interchainAccountAddr)
				suite.Require().NoError(err)

				// Check if account is created
				interchainAccount := icaAppB.AccountKeeper.GetAccount(suite.chainB.GetContext(), icaAddr)
				suite.Require().Equal(interchainAccount.GetAddress().String(), interchainAccountAddr)

				// Create bank transfer message to execute on the host
				sdkMsg = &banktypes.MsgSend{
					FromAddress: interchainAccountAddr,
					ToAddress:   suite.chainB.SenderAccount.GetAddress().String(),
					Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
				}
			}

			msg, err := types.NewMsgSubmitAction(owner, "label", []sdk.Msg{sdkMsg}, connectionId, "", "200s", "100s", uint64(suite.chainA.GetContext().BlockTime().Add(time.Hour).Unix()), sdk.Coins{}, &types.ExecutionConfiguration{SaveMsgResponses: false})
			suite.Require().NoError(err)
			wrappedCtx := sdk.WrapSDKContext(suite.chainA.GetContext())
			msgSrv := keeper.NewMsgServerImpl(GetActionKeeperFromApp(icaAppA))
			_, err = msgSrv.SubmitAction(wrappedCtx, msg)
			suite.Require().NoError(err)

			if addFakeExecHistory {
				action := icaAppA.IntentKeeper.GetActionInfo(sdk.UnwrapSDKContext(wrappedCtx), 1)
				fakeEntry := types.ActionHistoryEntry{ScheduledExecTime: action.ExecTime, ActualExecTime: action.ExecTime}

				icaAppA.IntentKeeper.SetActionHistoryEntry(sdk.UnwrapSDKContext(wrappedCtx), action.ID, &fakeEntry)
				suite.chainA.NextBlock()
				actionHistory := icaAppA.IntentKeeper.MustGetActionHistory(sdk.UnwrapSDKContext(wrappedCtx), 1)
				suite.Require().NotZero(actionHistory.History[0].ActualExecTime)
			}
			updateMsg, err := types.NewMsgUpdateAction(owner, 1, "new_label", []sdk.Msg{sdkMsg}, connectionId, newEndTime, newInterval, newStartAt, sdk.Coins{}, &types.ExecutionConfiguration{SaveMsgResponses: false})
			suite.Require().NoError(err)
			suite.chainA.Coordinator.IncrementTime()
			suite.Require().NotEqual(suite.chainA.GetContext(), wrappedCtx)
			wrappedCtx = sdk.WrapSDKContext(suite.chainA.GetContext())

			if addFakeExecHistory {
				actionHistory := icaAppA.IntentKeeper.MustGetActionHistory(sdk.UnwrapSDKContext(wrappedCtx), 1)
				suite.Require().NotZero(actionHistory.History[0].ActualExecTime)
				actionHistoryEntry, err := icaAppA.IntentKeeper.GetLatestActionHistoryEntry(sdk.UnwrapSDKContext(wrappedCtx), 1)
				suite.Require().NoError(err)
				suite.Require().NotNil(actionHistoryEntry)
			}

			res, err := msgSrv.UpdateAction(wrappedCtx, updateMsg)

			if !tc.expPass {
				suite.Require().Error(err)

			} else {
				suite.chainA.NextBlock()
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				action := icaAppA.IntentKeeper.GetActionInfo(sdk.UnwrapSDKContext(wrappedCtx), 1)
				suite.Require().Equal(action.Label, updateMsg.Label)
				suite.Require().Equal(action.StartTime.Unix(), int64(updateMsg.StartAt))
				suite.Require().Equal(action.EndTime.Unix(), int64(updateMsg.EndTime))
				suite.Require().Equal((action.Interval.String()), updateMsg.Interval)

			}

		})
	}
}
