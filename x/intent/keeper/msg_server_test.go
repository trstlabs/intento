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
		connectionId string
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
				connectionId = path.EndpointA.ConnectionID
			}, true,
		},
		{
			"failure - active channel does not exist for connection ID", func() {
				owner = TestOwnerAddress
				connectionId = "connection-100"

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
			msg, err := types.NewMsgSubmitTx(owner, sdkMsg, connectionId)
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
		connectionId              string
		hostConnectionId          string
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
				connectionId = path.EndpointA.ConnectionID
				hostConnectionId = path.EndpointB.ConnectionID
			}, true,
		},
		{
			"success - local flow", func() {
				registerInterchainAccount = false
				connectionId = ""
				hostConnectionId = ""
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
				connectionId = ""
				hostConnectionId = ""
				transferMsg = true
			}, true,
		},
		{
			"success - ICQ transfer", func() {
				registerInterchainAccount = false
				connectionId = ""
				hostConnectionId = ""
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
				connectionId = path.EndpointA.ConnectionID
				hostConnectionId = path.EndpointB.ConnectionID
				parseIcaAddress = true
				transferMsg = false
				conditions = types.ExecutionConditions{}
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

			msg, err := types.NewMsgSubmitFlow(owner, label, []sdk.Msg{sdkMsg}, connectionId, hostConnectionId, durationTimeText, intervalTimeText, startAt, sdk.Coins{}, "", sdk.Coin{}, &types.ExecutionConfiguration{FallbackToOwnerBalance: true}, &conditions)

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
			flow := flowKeeper.GetFlowInfo(ctx, 1)

			if len(flow.Conditions.FeedbackLoops) != 0 && flow.Conditions.FeedbackLoops[0].ICQConfig != nil {
				flowKeeper.SubmitInterchainQueries(ctx, flow, flowKeeper.Logger(ctx))
			}
			flowKeeper.HandleFlow(ctx, flowKeeper.Logger(ctx), flow, ctx.BlockTime(), nil)
			suite.IntentoChain.NextBlock()
			flow = flowKeeper.GetFlowInfo(ctx, 1)
			flowHistory, err := flowKeeper.GetFlowHistory(ctx, 1)

			suite.Require().NoError(err)
			suite.Require().NotEqual(flow, types.FlowInfo{})
			suite.Require().Equal(flow.Owner, owner)
			suite.Require().Equal(flow.Label, label)

			//ibc
			if flow.ICAConfig.PortID != "" {
				suite.Require().Equal(flow.ICAConfig.PortID, "icacontroller-"+owner)
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

				// flowNew := flowKeeper.GetFlowInfo(ctx, 1)
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
			msg, err := types.NewMsgSubmitFlow(owner, label, []sdk.Msg{sdkMsg}, "", "", durationTimeText, intervalTimeText, startAt, sdk.Coins{}, "", sdk.Coin{}, &types.ExecutionConfiguration{FallbackToOwnerBalance: true}, nil)
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
			connectionId := path.EndpointA.ConnectionID
			hostConnectionId := path.EndpointB.ConnectionID
			msgRegisterAndSubmit, err := types.NewMsgRegisterAccountAndSubmitFlow(owner, label, []sdk.Msg{sdkMsg}, connectionId, hostConnectionId, durationTimeText, intervalTimeText, startAt, sdk.Coins{}, &types.ExecutionConfiguration{FallbackToOwnerBalance: true}, "")
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

func (suite *KeeperTestSuite) TestCreateHostedAccount() {
	var (
		path                     *ibctesting.Path
		connectionId             string
		hostConnectionId         string
		sdkMsg                   sdk.Msg
		startAtBeforeBlockHeight bool
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{{

		"success - Create hosted account flow", func() {
			connectionId = path.EndpointA.ConnectionID
			hostConnectionId = path.EndpointB.ConnectionID
			sdkMsg = &banktypes.MsgSend{
				FromAddress: TestOwnerAddress,
				ToAddress:   suite.TestAccs[0].String(),
				Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100))),
			}
		}, true,
	},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest()
			path = NewICAPath(suite.IntentoChain, suite.HostChain)
			suite.Coordinator.SetupConnections(path)

			creator := suite.IntentoChain.SenderAccount.GetAddress().String()
			err := suite.SetupICAPath(path, creator)
			suite.Require().NoError(err)
			// Check if account is created
			portID, err := icatypes.NewControllerPortID(creator)
			suite.Require().NoError(err)
			interchainAccountAddr, found := GetICAApp(suite.IntentoChain).ICAControllerKeeper.GetInterchainAccountAddress(suite.IntentoChain.GetContext(), path.EndpointA.ConnectionID, portID)
			suite.Require().True(found)
			sdkMsg = &banktypes.MsgSend{
				FromAddress: interchainAccountAddr,
				ToAddress:   suite.HostChain.SenderAccount.GetAddress().String(),
				Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100))),
			}

			tc.malleate() // malleate mutates test data

			msgHosted := types.NewMsgCreateHostedAccount(creator, path.EndpointA.ConnectionID, path.EndpointA.ChannelConfig.Version, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.OneInt())))

			msgSrv := keeper.NewMsgServerImpl(GetICAApp(suite.IntentoChain).IntentKeeper)
			resHosted, err := msgSrv.CreateHostedAccount(suite.IntentoChain.GetContext(), msgHosted)
			suite.Require().Nil(err)
			suite.Require().NotNil(resHosted.Address)
			suite.Require().NotEqual(resHosted.Address, creator)

			label := "label"
			duration := time.Second * 200
			durationTimeText := duration.String()
			interval := time.Second * 100
			intervalTimeText := interval.String()
			startAt := uint64(0)
			if startAtBeforeBlockHeight {
				startAt = uint64(suite.IntentoChain.GetContext().BlockTime().Unix() - 60*60)
			}
			msg, err := types.NewMsgSubmitFlow(creator, label, []sdk.Msg{sdkMsg}, connectionId, hostConnectionId, durationTimeText, intervalTimeText, startAt, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.OneInt())), resHosted.Address, sdk.NewCoin(sdk.DefaultBondDenom, math.OneInt()), &types.ExecutionConfiguration{FallbackToOwnerBalance: true}, nil)
			suite.Require().NoError(err)

			msgSrv = keeper.NewMsgServerImpl(GetICAApp(suite.IntentoChain).IntentKeeper)
			res, err := msgSrv.SubmitFlow(suite.IntentoChain.GetContext(), msg)
			if !tc.expPass {
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
			"success - update flow", func() {
				registerInterchainAccount = true
				owner = TestOwnerAddress
				connectionId = path.EndpointA.ConnectionID
				newStartAt = uint64(time.Now().Unix())
				newEndTime = uint64(time.Now().Add(time.Hour).Unix())
				newInterval = "8m20s"
			}, true,
		},
		{
			"success - update local flow", func() {
				registerInterchainAccount = false
				owner = TestOwnerAddress
				connectionId = ""
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

			msg, err := types.NewMsgSubmitFlow(owner, "label", []sdk.Msg{sdkMsg}, connectionId, "", "200s", "100s", uint64(suite.IntentoChain.GetContext().BlockTime().Add(time.Hour).Unix()), sdk.Coins{}, "", sdk.Coin{}, &types.ExecutionConfiguration{SaveResponses: false}, nil)
			suite.Require().NoError(err)

			msgSrv := keeper.NewMsgServerImpl(GetICAApp(suite.IntentoChain).IntentKeeper)
			_, err = msgSrv.SubmitFlow(suite.IntentoChain.GetContext(), msg)
			suite.Require().NoError(err)

			if addFakeExecHistory {
				flow := icaAppA.IntentKeeper.GetFlowInfo(suite.IntentoChain.GetContext(), 1)
				fakeEntry := types.FlowHistoryEntry{ScheduledExecTime: flow.ExecTime, ActualExecTime: flow.ExecTime}

				icaAppA.IntentKeeper.SetFlowHistoryEntry(suite.IntentoChain.GetContext(), flow.ID, &fakeEntry)
				suite.IntentoChain.NextBlock()
				flowHistory := icaAppA.IntentKeeper.MustGetFlowHistory(suite.IntentoChain.GetContext(), 1)
				suite.Require().NotZero(flowHistory[0].ActualExecTime)
			}
			updateMsg, err := types.NewMsgUpdateFlow(owner, 1, "new_label", []sdk.Msg{sdkMsg}, connectionId, newEndTime, newInterval, newStartAt, sdk.Coins{}, "", sdk.Coin{}, &types.ExecutionConfiguration{SaveResponses: false}, nil)
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
				flow := icaAppA.IntentKeeper.GetFlowInfo(suite.IntentoChain.GetContext(), 1)
				suite.Require().Equal(flow.Label, updateMsg.Label)
				suite.Require().Equal(flow.StartTime.Unix(), int64(updateMsg.StartAt))
				suite.Require().Equal(flow.EndTime.Unix(), int64(updateMsg.EndTime))
				suite.Require().Equal((flow.Interval.String()), updateMsg.Interval)

			}

		})
	}
}

func (suite *KeeperTestSuite) TestUpdateHostedAccount() {
	var (
		path            *ibctesting.Path
		newConnectionId string
		newAdmin        string
		newFeeAmount    uint64
		newDenom        string
		newVersion      string
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"success - update hosted", func() {
				newConnectionId = "connection-123"
				newAdmin = "cosmos1wdplq6qjh2xruc7qqagma9ya665q6qhcwju3ng"
				newFeeAmount = 32434554
				newDenom = "utrst"
				newVersion = "v123"

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

			msgHosted := types.NewMsgCreateHostedAccount(admin, path.EndpointA.ConnectionID, path.EndpointA.ChannelConfig.Version, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.OneInt())))

			msgSrv := keeper.NewMsgServerImpl(GetICAApp(suite.IntentoChain).IntentKeeper)
			resHosted, err := msgSrv.CreateHostedAccount(suite.IntentoChain.GetContext(), msgHosted)
			suite.Require().Nil(err)
			suite.Require().NotNil(resHosted.Address)

			hosted := GetICAApp(suite.IntentoChain).IntentKeeper.GetHostedAccount(suite.IntentoChain.GetContext(), resHosted.Address)

			msg := types.NewMsgUpdateHostedAccount(admin, resHosted.Address, newConnectionId, newVersion, newAdmin, sdk.NewCoins(sdk.NewCoin(newDenom, math.NewIntFromUint64(newFeeAmount))))
			suite.Require().NoError(err)

			msgSrv = keeper.NewMsgServerImpl(GetICAApp(suite.IntentoChain).IntentKeeper)
			res, err := msgSrv.UpdateHostedAccount(suite.IntentoChain.GetContext(), msg)
			suite.IntentoChain.NextBlock()
			hostedNew := GetICAApp(suite.IntentoChain).IntentKeeper.GetHostedAccount(suite.IntentoChain.GetContext(), resHosted.Address)
			if !tc.expPass {
				suite.Require().Error(err)
				suite.Require().Nil(res)
				return
			}

			suite.Require().Equal(hosted.HostedAddress, hostedNew.HostedAddress)
			suite.Require().NotEqual(hosted.HostFeeConfig.Admin, hostedNew.HostFeeConfig.Admin)
			suite.Require().NotEqual(hosted.HostFeeConfig.FeeCoinsSuported.Denoms(), hostedNew.HostFeeConfig.FeeCoinsSuported.Denoms())
			suite.Require().NotEqual(hosted.HostFeeConfig.FeeCoinsSuported[0].Amount, hostedNew.HostFeeConfig.FeeCoinsSuported[0].Amount)
			suite.Require().NotEqual(hosted.ICAConfig.ConnectionID, hostedNew.ICAConfig.ConnectionID)

		})
	}
}
