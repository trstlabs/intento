package keeper_test

import (
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"
	"github.com/trstlabs/trst/x/auto-ibc-tx/keeper"
	"github.com/trstlabs/trst/x/auto-ibc-tx/types"
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

			msgSrv := keeper.NewMsgServerImpl(GetAutoTxKeeper(suite.chainA))
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

			msgSrv := keeper.NewMsgServerImpl(GetAutoTxKeeperFromApp(icaAppA))
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

func (suite *KeeperTestSuite) TestSubmitAutoTx() {
	var (
		path                      *ibctesting.Path
		registerInterchainAccount bool
		noOwner                   bool
		connectionId              string
		sdkMsg                    sdk.Msg
		parseIcaAddress           bool
		startAtBeforeBlockHeight  bool
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"success - IBC autoTx", func() {
				registerInterchainAccount = true
				connectionId = path.EndpointA.ConnectionID
			}, true,
		},
		{
			"success - local autoTx", func() {
				registerInterchainAccount = false
				connectionId = ""
				sdkMsg = &banktypes.MsgSend{
					FromAddress: suite.chainA.SenderAccount.GetAddress().String(),
					ToAddress:   TestOwnerAddress,
					Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
				}
			}, true,
		},
		{
			"success - parse ICA address", func() {
				registerInterchainAccount = true
				connectionId = path.EndpointA.ConnectionID
				parseIcaAddress = true
			}, true,
		},
		{
			"failure - start_at before block_height", func() {
				registerInterchainAccount = true
				noOwner = false
				connectionId = path.EndpointA.ConnectionID
				startAtBeforeBlockHeight = true
			}, false,
		},
		{
			"failure - owner address is empty", func() {
				registerInterchainAccount = false
				noOwner = true
				connectionId = path.EndpointA.ConnectionID
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
			msg, err := types.NewMsgSubmitAutoTx(owner, label, []sdk.Msg{sdkMsg}, connectionId, durationTimeText, intervalTimeText, startAt, sdk.Coins{}, &types.ExecutionConfiguration{SaveMsgResponses: false})
			suite.Require().NoError(err)
			wrappedCtx := sdk.WrapSDKContext(suite.chainA.GetContext())

			msgSrv := keeper.NewMsgServerImpl(GetAutoTxKeeperFromApp(icaAppA))
			res, err := msgSrv.SubmitAutoTx(wrappedCtx, msg)

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
			autoTxKeeper := icaAppA.AutoIbcTxKeeper
			autoTx := autoTxKeeper.GetAutoTxInfo(suite.chainA.GetContext(), 1)

			suite.chainA.CurrentHeader.Time = suite.chainA.CurrentHeader.Time.Add(interval)
			FakeBeginBlocker(suite.chainA.GetContext(), autoTxKeeper, sdk.ConsAddress(suite.chainA.Vals.Proposer.Address))
			suite.chainA.NextBlock()

			autoTx = autoTxKeeper.GetAutoTxInfo(suite.chainA.GetContext(), 1)

			suite.Require().NotEqual(autoTx, types.AutoTxInfo{})
			suite.Require().Equal(autoTx.Owner, owner)
			suite.Require().Equal(autoTx.Label, label)
			suite.Require().Contains(autoTx.Msgs[0].TypeUrl, "MsgSend")

			//ibc autotx
			if connectionId != "" {
				suite.Require().Equal(autoTx.PortID, "icacontroller-"+owner)
				suite.Require().Equal(autoTx.ConnectionID, path.EndpointA.ConnectionID)
			}

			if parseIcaAddress {
				var txMsg sdk.Msg
				err = icaAppA.AppCodec().UnpackAny(autoTx.Msgs[0], &txMsg)
				suite.Require().NoError(err)
				index := strings.Index(txMsg.String(), types.ParseICAValue)
				suite.Require().Equal(-1, index)

			}

		})
	}
}

func (suite *KeeperTestSuite) TestUpdateAutoTx() {
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
			"success - update autoTx", func() {
				registerInterchainAccount = true
				owner = TestOwnerAddress
				connectionId = path.EndpointA.ConnectionID
				newStartAt = uint64(time.Now().Unix())
				newEndTime = uint64(time.Now().Add(time.Hour).Unix())
				newInterval = "8m20s"
			}, true,
		},
		{
			"success - update local autoTx", func() {
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

			//sendTo, _ := CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
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

			msg, err := types.NewMsgSubmitAutoTx(owner, "label", []sdk.Msg{sdkMsg}, connectionId, "200s", "100s", uint64(suite.chainA.GetContext().BlockTime().Add(time.Hour).Unix()), sdk.Coins{}, &types.ExecutionConfiguration{SaveMsgResponses: false})
			suite.Require().NoError(err)
			wrappedCtx := sdk.WrapSDKContext(suite.chainA.GetContext())
			msgSrv := keeper.NewMsgServerImpl(GetAutoTxKeeperFromApp(icaAppA))
			_, err = msgSrv.SubmitAutoTx(wrappedCtx, msg)
			suite.Require().NoError(err)

			if addFakeExecHistory {
				autoTx := icaAppA.AutoIbcTxKeeper.GetAutoTxInfo(sdk.UnwrapSDKContext(wrappedCtx), 1)
				suite.Require().NotNil(autoTx)
				fakeEntry := types.AutoTxHistoryEntry{ScheduledExecTime: autoTx.ExecTime}
				autoTx.AutoTxHistory = []*types.AutoTxHistoryEntry{&fakeEntry}
				icaAppA.AutoIbcTxKeeper.SetAutoTxInfo(sdk.UnwrapSDKContext(wrappedCtx), &autoTx)
				suite.chainA.NextBlock()
			}
			updateMsg, err := types.NewMsgUpdateAutoTx(owner, 1, "new_label", []sdk.Msg{sdkMsg}, connectionId, newEndTime, newInterval, newStartAt, sdk.Coins{}, &types.ExecutionConfiguration{SaveMsgResponses: false})
			suite.Require().NoError(err)
			suite.chainA.Coordinator.IncrementTime()
			suite.Require().NotEqual(suite.chainA.GetContext(), wrappedCtx)
			wrappedCtx = sdk.WrapSDKContext(suite.chainA.GetContext())

			msgSrv = keeper.NewMsgServerImpl(GetAutoTxKeeperFromApp(icaAppA))
			res, err := msgSrv.UpdateAutoTx(wrappedCtx, updateMsg)

			if !tc.expPass {
				suite.Require().Nil(res)
				suite.Require().Error(err)

			} else {
				suite.chainA.NextBlock()
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				autoTx := icaAppA.AutoIbcTxKeeper.GetAutoTxInfo(sdk.UnwrapSDKContext(wrappedCtx), 1)
				suite.Require().Equal(autoTx.Label, updateMsg.Label)
				suite.Require().Equal(autoTx.StartTime.Unix(), int64(updateMsg.StartAt))
				suite.Require().Equal(autoTx.EndTime.Unix(), int64(updateMsg.EndTime))
				suite.Require().Equal((autoTx.Interval.String()), updateMsg.Interval)

			}

		})
	}
}
