package keeper_test

import (
	"strings"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v8/testing"

	"github.com/trstlabs/intento/x/intent/keeper"
	"github.com/trstlabs/intento/x/intent/types"
	icqkeeper "github.com/trstlabs/intento/x/interchainquery/keeper"
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
				GetICAApp(suite.IntentoChain).GetIBCKeeper().PortKeeper.BindPort(suite.IntentoChain.GetContext(), TestPortID)
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
				GetICAApp(suite.IntentoChain).IBCKeeper.ChannelKeeper.SetChannel(suite.IntentoChain.GetContext(), portID, ibctesting.FirstChannelID, channel)

				GetICAApp(suite.IntentoChain).ICAControllerKeeper.SetActiveChannelID(suite.IntentoChain.GetContext(), ibctesting.FirstConnectionID, portID, ibctesting.FirstChannelID)
			},
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest()

			owner = TestOwnerAddress // must be explicitly changed

			path = NewICAPath(suite.IntentoChain, suite.HostChain)
			suite.Coordinator.SetupConnections(path)

			tc.malleate() // malleate mutates test data

			msgSrv := keeper.NewMsgServerImpl(GetICAApp(suite.IntentoChain).IntentKeeper)
			msg := types.NewMsgRegisterAccount(owner, path.EndpointA.ConnectionID, path.EndpointA.ChannelConfig.Version)

			res, err := msgSrv.RegisterAccount(suite.IntentoChain.GetContext(), msg)

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
		path         *ibctesting.Path
		owner        string
		connectionID string
		sdkMsg       sdk.Msg
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"success", func() {
				owner = TestOwnerAddress
				connectionID = path.EndpointA.ConnectionID
			}, true,
		},
		{
			"failure - active channel does not exist for connection ID", func() {
				owner = TestOwnerAddress
				connectionID = "connection-100"

			}, false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.SetupTest()
		path = NewICAPath(suite.IntentoChain, suite.HostChain)
		suite.Coordinator.SetupConnections(path)
		owner = suite.IntentoChain.SenderAccount.GetAddress().String()

		tc.malleate() // malleate mutates test data
		suite.Run(tc.name, func() {
			err := suite.SetupICAPath(path, owner)
			suite.Require().NoError(err)
			portID, err := icatypes.NewControllerPortID(owner)
			suite.Require().NoError(err)

			interchainAccountAddr, found := GetICAApp(suite.IntentoChain).ICAControllerKeeper.GetInterchainAccountAddress(suite.IntentoChain.GetContext(), path.EndpointA.ConnectionID, portID)
			suite.Require().True(found)
			sdkMsg = &banktypes.MsgSend{
				FromAddress: interchainAccountAddr,
				ToAddress:   suite.HostChain.SenderAccount.GetAddress().String(),
				Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100))),
			}

			msgSrv := keeper.NewMsgServerImpl(GetICAApp(suite.IntentoChain).IntentKeeper)
			msg, err := types.NewMsgSubmitTx(owner, sdkMsg, connectionID)
			suite.Require().NoError(err)

			res, err := msgSrv.SubmitTx(suite.IntentoChain.GetContext(), msg)

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

func (suite *KeeperTestSuite) TestSubmitFlow() {

	types.Denom = sdk.DefaultBondDenom
	var (
		path                      *ibctesting.Path
		registerInterchainAccount bool
		noOwner                   bool
		connectionID              string
		sdkMsg                    sdk.Msg
		parseIcaAddress           bool
		startAtBeforeBlockHeight  bool
		transferMsg               bool
		conditions                types.ExecutionConditions
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"success - IBC ICA flow", func() {
				registerInterchainAccount = true
				connectionID = path.EndpointA.ConnectionID
			}, true,
		},
		{
			"success - local flow", func() {
				registerInterchainAccount = false
				connectionID = ""
				sdkMsg = &banktypes.MsgSend{
					FromAddress: suite.IntentoChain.SenderAccount.GetAddress().String(),
					ToAddress:   TestOwnerAddress,
					Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100))),
				}
			}, true,
		},
		{
			"success - IBC transfer", func() {
				registerInterchainAccount = false
				connectionID = ""
				transferMsg = true
			}, true,
		},
		{
			"success - ICQ transfer", func() {
				registerInterchainAccount = false
				connectionID = ""
				transferMsg = false
				sdkMsg = &banktypes.MsgSend{
					FromAddress: suite.IntentoChain.SenderAccount.GetAddress().String(),
					ToAddress:   TestOwnerAddress,
					Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100))),
				}
				conditions = types.ExecutionConditions{FeedbackLoops: []*types.FeedbackLoop{{MsgsIndex: 0, MsgKey: "amount", ResponseIndex: 0, ResponseKey: "", ValueType: "sdk.Coin", ICQConfig: &types.ICQConfig{ChainId: suite.HostChain.ChainID, QueryType: "store/bank/key", QueryKey: "fake_key_owner", TimeoutPolicy: 1, TimeoutDuration: time.Second * 30}}}}
			}, true,
		},
		{
			"success - parse ICA address", func() {
				registerInterchainAccount = true
				connectionID = path.EndpointA.ConnectionID
				parseIcaAddress = true
				transferMsg = false
				conditions = types.ExecutionConditions{}
			}, true,
		},
		{
			"failure - start before block height", func() {
				registerInterchainAccount = true
				noOwner = false
				connectionID = path.EndpointA.ConnectionID
				startAtBeforeBlockHeight = true
			}, false,
		},
		{
			"failure - owner address is empty", func() {
				registerInterchainAccount = false
				noOwner = true
				connectionID = path.EndpointA.ConnectionID
				parseIcaAddress = false
			}, false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest()
			owner := suite.IntentoChain.SenderAccount.GetAddress().String()
			path = NewICAPath(suite.IntentoChain, suite.HostChain)
			suite.Coordinator.SetupConnections(path)

			params := types.DefaultParams()
			params.GasFeeCoins = sdk.NewCoins(sdk.NewCoin("stake", math.OneInt()))
			params.FlowFlexFeeMul = 1

			icaAppA := GetICAApp(suite.IntentoChain)
			icaAppA.IntentKeeper.SetParams(suite.IntentoChain.GetContext(), params)

			tc.malleate() // malleate mutates test data

			ctx := suite.IntentoChain.GetContext()
			if transferMsg {
				path = NewTransferPath(suite.IntentoChain, suite.HostChain)
				suite.Coordinator.SetupConnections(path)
				sdkMsg = transfertypes.NewMsgTransfer(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100)), owner, suite.HostChain.SenderAccount.GetAddress().String(), suite.HostChain.GetTimeoutHeight(), 0, "")
			}

			if conditions.FeedbackLoops != nil && conditions.FeedbackLoops[0].ICQConfig != nil {
				path = NewTransferPath(suite.IntentoChain, suite.HostChain)
				suite.Coordinator.SetupConnections(path)
				conditions.FeedbackLoops[0].ICQConfig.ConnectionId = path.EndpointA.ConnectionID
			}

			if noOwner {
				owner = ""
			}
			if registerInterchainAccount {

				err := suite.SetupICAPath(path, owner)
				suite.Require().NoError(err)
				// Check if account is created
				portID, err := icatypes.NewControllerPortID(owner)
				suite.Require().NoError(err)
				interchainAccountAddr, found := GetICAApp(suite.IntentoChain).ICAControllerKeeper.GetInterchainAccountAddress(ctx, path.EndpointA.ConnectionID, portID)
				suite.Require().True(found)
				sdkMsg = &banktypes.MsgSend{
					FromAddress: interchainAccountAddr,
					ToAddress:   suite.HostChain.SenderAccount.GetAddress().String(),
					Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100))),
				}
			}

			label := "label"
			duration := time.Second * 200
			durationTimeText := duration.String()
			interval := time.Second * 100
			intervalTimeText := interval.String()
			startAt := uint64(0)
			if startAtBeforeBlockHeight {
				startAt = uint64(ctx.BlockTime().Unix() - 60*60)
			}

			msg, err := types.NewMsgSubmitFlow(owner, label, []sdk.Msg{sdkMsg}, connectionID, durationTimeText, intervalTimeText, startAt, sdk.Coins{}, "", sdk.Coins{}, &types.ExecutionConfiguration{WalletFallback: true}, &conditions)

			suite.Require().NoError(err)

			msgSrv := keeper.NewMsgServerImpl(GetICAApp(suite.IntentoChain).IntentKeeper)
			res, err := msgSrv.SubmitFlow(suite.IntentoChain.GetContext(), msg)

			if !tc.expPass {
				suite.Require().Error(err)
				suite.Require().Nil(res)
				return
			}

			suite.Require().NoError(err)
			suite.Require().NotNil(res)

			flowKeeper := icaAppA.IntentKeeper

			suite.IntentoChain.CurrentHeader.Time = suite.IntentoChain.CurrentHeader.Time.Add(interval)
			flow := flowKeeper.GetFlow(ctx, 1)

			if len(flow.Conditions.FeedbackLoops) != 0 && flow.Conditions.FeedbackLoops[0].ICQConfig != nil {
				flowKeeper.SubmitInterchainQueries(ctx, flow, flowKeeper.Logger(ctx))
			}
			flowKeeper.HandleFlow(ctx, flowKeeper.Logger(ctx), flow, ctx.BlockTime())
			suite.IntentoChain.NextBlock()
			flow = flowKeeper.GetFlow(ctx, 1)
			flowHistory, err := flowKeeper.GetFlowHistory(ctx, 1)

			suite.Require().NoError(err)
			suite.Require().NotEqual(flow, types.Flow{})
			suite.Require().Equal(flow.Owner, owner)
			suite.Require().Equal(flow.Label, label)

			//ibc
			if flow.SelfHostedICA.PortID != "" {
				suite.Require().Equal(flow.SelfHostedICA.PortID, "icacontroller-"+owner)
			}
			if !transferMsg {
				//icq should return this error
				if flowHistory[0].Errors != nil {
					// fmt.Printf("ICQ MSG %v", flow.Msgs[0].GetCachedValue())
					suite.Require().Contains(flowHistory[0].Errors[0], "Error submitting ICQ")
				}

			}
			if len(flow.Conditions.FeedbackLoops) != 0 && msg.Conditions.FeedbackLoops[0].ICQConfig != nil {
				msgQueryResp, _ := suite.SetupMsgSubmitQueryResponse(*msg.Conditions.FeedbackLoops[0].ICQConfig, flow.ID)

				msgSrvICQ := icqkeeper.NewMsgServerImpl(GetICAApp(suite.IntentoChain).InterchainQueryKeeper)
				_, err := msgSrvICQ.SubmitQueryResponse(ctx, &msgQueryResp)

				//as we cannot fully test this, the code should run
				suite.Require().NoError(err)

				// flowNew := flowKeeper.GetFlow(ctx, 1)
				// suite.Require().Equal(flow.Msgs, flowNew.Msgs)
				// suite.Require().Nil(flowNew.ValidateBasic())

			}
			if parseIcaAddress {
				var txMsg sdk.Msg
				err = icaAppA.AppCodec().UnpackAny(flow.Msgs[0], &txMsg)
				suite.Require().NoError(err)
				index := strings.Index(txMsg.String(), types.ParseICAValue)
				suite.Require().Equal(-1, index)

			}

		})
	}
}

func (suite *KeeperTestSuite) TestSubmitFlowSigner() {
	icaAddrString := "into1wdplq6qjh2xruc7qqagma9ya665q6qhcpse4k6"
	testCases := []struct {
		name      string
		malleate  func()
		from      string
		expPass   bool
		withAuthZ bool
	}{
		{
			"success - MsgExec message with ICA sender", func() {}, suite.IntentoChain.SenderAccount.GetAddress().String(), true, true,
		},
		{
			"fail - MsgExec message with other from ", func() {}, TestOwnerAddress, false, true,
		},
		{
			"fail - MsgExec message with ICA with other from & prefix", func() {}, "into1g6qdx6kdhpf000afvvpte7hp0vnpzaputyyrem", false, true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {

			owner := suite.IntentoChain.SenderAccount.GetAddress().String()
			suite.SetupTest()

			tc.malleate() // malleate mutates test data

			var sdkMsg sdk.Msg
			sdkMsg = &banktypes.MsgSend{
				FromAddress: tc.from,
				ToAddress:   TestOwnerAddress,
				Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100))),
			}
			if tc.withAuthZ {
				anyMsg, _ := types.PackTxMsgAnys([]sdk.Msg{sdkMsg})
				sdkMsgAuthZ := &authztypes.MsgExec{Grantee: icaAddrString,
					Msgs: anyMsg,
				}
				sdkMsg = sdkMsgAuthZ
			}
			label := "label"
			duration := time.Second * 200
			durationTimeText := duration.String()
			interval := time.Second * 100
			intervalTimeText := interval.String()
			startAt := uint64(0)
			GetICAApp(suite.IntentoChain).ICAControllerKeeper.SetInterchainAccountAddress(suite.IntentoChain.GetContext(), "", "", icaAddrString)
			msg, err := types.NewMsgSubmitFlow(owner, label, []sdk.Msg{sdkMsg}, "", durationTimeText, intervalTimeText, startAt, sdk.Coins{}, "", sdk.Coins{}, &types.ExecutionConfiguration{WalletFallback: true}, nil)
			suite.Require().NoError(err)
			err = msg.ValidateBasic()
			suite.Require().NoError(err)

			msgSrv := keeper.NewMsgServerImpl(GetICAApp(suite.IntentoChain).IntentKeeper)
			res, err := msgSrv.SubmitFlow(suite.IntentoChain.GetContext(), msg)

			if !tc.expPass {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), "message signer is not message sender")
				suite.Require().Nil(res)
				return
			}
			suite.Require().NoError(err)

			//test the same for
			path := NewICAPath(suite.IntentoChain, suite.HostChain)
			suite.Coordinator.SetupConnections(path)
			connectionID := path.EndpointA.ConnectionID
			hostConnectionID := path.EndpointB.ConnectionID
			msgRegisterAndSubmit, err := types.NewMsgRegisterAccountAndSubmitFlow(owner, label, []sdk.Msg{sdkMsg}, connectionID, hostConnectionID, durationTimeText, intervalTimeText, startAt, sdk.Coins{}, &types.ExecutionConfiguration{WalletFallback: true}, "")
			suite.Require().NoError(err)
			err = msg.ValidateBasic()
			suite.Require().NoError(err)

			msgSrv = keeper.NewMsgServerImpl(GetICAApp(suite.IntentoChain).IntentKeeper)
			resRegisterAndSubmit, err := msgSrv.RegisterAccountAndSubmitFlow(suite.IntentoChain.GetContext(), msgRegisterAndSubmit)
			if !tc.expPass {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), "message signer is not message sender")
				suite.Require().Nil(resRegisterAndSubmit)
				return
			}
		})
	}
}

func (suite *KeeperTestSuite) TestCreateTrustlessAgent() {
	var (
		connectionID string
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"success - Create trustless agent flow",
			func() {
				path := NewICAPath(suite.IntentoChain, suite.HostChain)
				suite.Coordinator.SetupConnections(path)
				connectionID = path.EndpointA.ConnectionID
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.malleate()

			// Create a new trustless agent
			msgHosted := &types.MsgCreateTrustlessAgent{
				Creator:      suite.TestAccs[0].String(),
				ConnectionID: connectionID,
				Version:      TestVersion,
			}

			msgSrv := keeper.NewMsgServerImpl(GetICAApp(suite.IntentoChain).IntentKeeper)
			res, err := msgSrv.CreateTrustlessAgent(suite.IntentoChain.GetContext(), msgHosted)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().NotEmpty(res.Address)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestUpdateFlow() {
	var (
		path                      *ibctesting.Path
		registerInterchainAccount bool
		owner                     string
		connectionID              string
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
			"success - update flow", func() {
				registerInterchainAccount = true
				owner = TestOwnerAddress
				connectionID = path.EndpointA.ConnectionID
				newStartAt = uint64(time.Now().Unix())
				newEndTime = uint64(time.Now().Add(time.Hour).Unix())
				newInterval = "8m20s"
			}, true,
		},
		{
			"success - update local flow", func() {
				registerInterchainAccount = false
				owner = TestOwnerAddress
				connectionID = ""
				sdkMsg = &banktypes.MsgSend{
					FromAddress: TestOwnerAddress,
					ToAddress:   suite.HostChain.SenderAccount.GetAddress().String(),
					Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100))),
				}
			}, true,
		},
		{
			"failure - interval shorter than min duration", func() {
				registerInterchainAccount = true
				owner = TestOwnerAddress
				connectionID = path.EndpointA.ConnectionID
				newInterval = "1s"
			}, false,
		},
		{
			"failure - start time can not be changed after execution entry", func() {
				registerInterchainAccount = true
				owner = TestOwnerAddress
				connectionID = path.EndpointA.ConnectionID
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
			path = NewICAPath(suite.IntentoChain, suite.HostChain)
			suite.Coordinator.SetupConnections(path)

			icaAppA := GetICAApp(suite.IntentoChain)

			tc.malleate() // malleate mutates test data

			if registerInterchainAccount {

				err := suite.SetupICAPath(path, owner)
				suite.Require().NoError(err)
				// Check if account is created
				portID, err := icatypes.NewControllerPortID(owner)
				suite.Require().NoError(err)
				interchainAccountAddr, found := GetICAApp(suite.IntentoChain).ICAControllerKeeper.GetInterchainAccountAddress(suite.IntentoChain.GetContext(), path.EndpointA.ConnectionID, portID)
				suite.Require().True(found)
				sdkMsg = &banktypes.MsgSend{
					FromAddress: interchainAccountAddr,
					ToAddress:   suite.HostChain.SenderAccount.GetAddress().String(),
					Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100))),
				}
			}

			msg, err := types.NewMsgSubmitFlow(owner, "label", []sdk.Msg{sdkMsg}, connectionID, "200s", "100s", uint64(suite.IntentoChain.GetContext().BlockTime().Add(time.Hour).Unix()), sdk.Coins{}, "", sdk.Coins{}, &types.ExecutionConfiguration{SaveResponses: false}, nil)
			suite.Require().NoError(err)

			msgSrv := keeper.NewMsgServerImpl(GetICAApp(suite.IntentoChain).IntentKeeper)
			_, err = msgSrv.SubmitFlow(suite.IntentoChain.GetContext(), msg)
			suite.Require().NoError(err)

			if addFakeExecHistory {
				flow := icaAppA.IntentKeeper.GetFlow(suite.IntentoChain.GetContext(), 1)
				fakeEntry := types.FlowHistoryEntry{ScheduledExecTime: flow.ExecTime, ActualExecTime: flow.ExecTime}

				icaAppA.IntentKeeper.SetFlowHistoryEntry(suite.IntentoChain.GetContext(), flow.ID, &fakeEntry)
				suite.IntentoChain.NextBlock()
				flowHistory := icaAppA.IntentKeeper.MustGetFlowHistory(suite.IntentoChain.GetContext(), 1)
				suite.Require().NotZero(flowHistory[0].ActualExecTime)
			}
			updateMsg, err := types.NewMsgUpdateFlow(owner, 1, "new_label", []sdk.Msg{sdkMsg}, connectionID, newEndTime, newInterval, newStartAt, sdk.Coins{}, "", sdk.Coins{}, &types.ExecutionConfiguration{SaveResponses: false}, nil)
			suite.Require().NoError(err)
			suite.IntentoChain.Coordinator.IncrementTime()

			if addFakeExecHistory {
				flowHistory := icaAppA.IntentKeeper.MustGetFlowHistory(suite.IntentoChain.GetContext(), 1)
				suite.Require().NotZero(flowHistory[0].ActualExecTime)
				flowHistoryEntry, err := icaAppA.IntentKeeper.GetLatestFlowHistoryEntry(suite.IntentoChain.GetContext(), 1)
				suite.Require().NoError(err)
				suite.Require().NotNil(flowHistoryEntry)
			}

			res, err := msgSrv.UpdateFlow(suite.IntentoChain.GetContext(), updateMsg)

			if !tc.expPass {
				suite.Require().Error(err)

			} else {
				suite.IntentoChain.NextBlock()
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				flow := icaAppA.IntentKeeper.GetFlow(suite.IntentoChain.GetContext(), 1)
				suite.Require().Equal(flow.Label, updateMsg.Label)
				suite.Require().Equal(flow.StartTime.Unix(), int64(updateMsg.StartAt))
				suite.Require().Equal(flow.EndTime.Unix(), int64(updateMsg.EndTime))
				suite.Require().Equal((flow.Interval.String()), updateMsg.Interval)

			}

		})
	}
}

func (suite *KeeperTestSuite) TestUpdateTrustlessAgent() {
	var (
		path         *ibctesting.Path
		newAdmin     string
		newFeeAmount uint64
		newDenom     string
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"success - update hosted", func() {
				newAdmin = "cosmos1wdplq6qjh2xruc7qqagma9ya665q6qhcwju3ng"
				newFeeAmount = 32434554
				newDenom = "utrst"

			}, true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest()
			path = NewICAPath(suite.IntentoChain, suite.HostChain)
			suite.Coordinator.SetupConnections(path)

			tc.malleate() // malleate mutates test data

			admin := suite.IntentoChain.SenderAccount.GetAddress().String()

			msgHosted := types.NewMsgCreateTrustlessAgent(admin, path.EndpointA.ConnectionID, path.EndpointA.ChannelConfig.Version, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.OneInt())))

			msgSrv := keeper.NewMsgServerImpl(GetICAApp(suite.IntentoChain).IntentKeeper)
			resHosted, err := msgSrv.CreateTrustlessAgent(suite.IntentoChain.GetContext(), msgHosted)
			suite.Require().Nil(err)
			suite.Require().NotNil(resHosted.Address)

			hosted := GetICAApp(suite.IntentoChain).IntentKeeper.GetTrustlessAgent(suite.IntentoChain.GetContext(), resHosted.Address)

			msg := types.NewMsgUpdateTrustlessAgent(admin, resHosted.Address, newAdmin, sdk.NewCoins(sdk.NewCoin(newDenom, math.NewIntFromUint64(newFeeAmount))))
			suite.Require().NoError(err)

			msgSrv = keeper.NewMsgServerImpl(GetICAApp(suite.IntentoChain).IntentKeeper)
			res, err := msgSrv.UpdateTrustlessAgentFeeConfig(suite.IntentoChain.GetContext(), msg)
			suite.IntentoChain.NextBlock()
			hostedNew := GetICAApp(suite.IntentoChain).IntentKeeper.GetTrustlessAgent(suite.IntentoChain.GetContext(), resHosted.Address)
			if !tc.expPass {
				suite.Require().Error(err)
				suite.Require().Nil(res)
				return
			}

			suite.Require().Equal(hosted.AgentAddress, hostedNew.AgentAddress)
			suite.Require().NotEqual(hosted.FeeConfig.FeeAdmin, hostedNew.FeeConfig.FeeAdmin)
			suite.Require().NotEqual(hosted.FeeConfig.FeeCoinsSupported.Denoms(), hostedNew.FeeConfig.FeeCoinsSupported.Denoms())
			suite.Require().NotEqual(hosted.FeeConfig.FeeCoinsSupported[0].Amount, hostedNew.FeeConfig.FeeCoinsSupported[0].Amount)

		})
	}
}

func (suite *KeeperTestSuite) TestUpdateParams() {
	// Get the current params
	params, err := GetICAApp(suite.IntentoChain).IntentKeeper.GetParams(suite.IntentoChain.GetContext())
	suite.Require().NoError(err)

	// Create a new param value that's different from the current one
	newParams := params
	// Update a field in the params if needed
	// For example: newParams.SomeField = newValue

	testCases := []struct {
		name      string
		expPass   bool
		authority string
		params    types.Params
	}{
		{
			"valid authority and params",
			true,
			GetICAApp(suite.IntentoChain).IntentKeeper.GetAuthority(),
			newParams,
		},
		{
			"invalid authority",
			false,
			"invalid_authority",
			newParams,
		},
		{
			"empty authority",
			false,
			"",
			newParams,
		},
		{
			"invalid params",
			false,
			GetICAApp(suite.IntentoChain).IntentKeeper.GetAuthority(),
			types.Params{}, // Add invalid params here if needed
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			msg := types.MsgUpdateParams{
				Authority: tc.authority,
				Params:    tc.params,
			}

			msgSrv := keeper.NewMsgServerImpl(GetICAApp(suite.IntentoChain).IntentKeeper)
			txResponse, err := msgSrv.UpdateParams(suite.IntentoChain.GetContext(), &msg)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(txResponse)

				// Verify that the params were updated
				updatedParams, err := GetICAApp(suite.IntentoChain).IntentKeeper.GetParams(suite.IntentoChain.GetContext())
				suite.Require().NoError(err)
				suite.Require().Equal(tc.params, updatedParams)

				// Restore the original params for the next test case
				GetICAApp(suite.IntentoChain).IntentKeeper.SetParams(suite.IntentoChain.GetContext(), params)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
