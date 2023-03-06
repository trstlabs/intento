package keeper_test

import (
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	//banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	icatypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v4/testing"
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
				GetICAApp(suite.chainA.TestChain).AppKeepers.IbcKeeper.ChannelKeeper.SetChannel(suite.chainA.GetContext(), portID, ibctesting.FirstChannelID, channel)

				GetICAApp(suite.chainA.TestChain).AppKeepers.ICAControllerKeeper.SetActiveChannelID(suite.chainA.GetContext(), ibctesting.FirstConnectionID, portID, ibctesting.FirstChannelID)
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

			msgSrv := keeper.NewMsgServerImpl(GetICAKeeper(suite.chainA))
			msg := types.NewMsgRegisterAccount(owner, path.EndpointA.ConnectionID, path.EndpointA.ChannelConfig.Version)

			res, err := msgSrv.RegisterAccount(sdk.WrapSDKContext(suite.chainA.GetContext()), msg)

			// resp, err := GetICAApp(suite.chainA).AppKeepers.AutoIBCTXKeeper.InterchainAccountFromAddress(sdk.WrapSDKContext(suite.chainA.GetContext()), &types.QueryInterchainAccountFromAddressRequest{Owner: owner, ConnectionId: path.EndpointA.ConnectionID})
			// fmt.Printf("%v", resp)

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
		icaMsg                    sdk.Msg
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
				interchainAccountAddr, found := GetICAApp(suite.chainA.TestChain).AppKeepers.ICAControllerKeeper.GetInterchainAccountAddress(suite.chainA.GetContext(), path.EndpointA.ConnectionID, portID)
				suite.Require().True(found)

				icaAddr, err := sdk.AccAddressFromBech32(interchainAccountAddr)
				suite.Require().NoError(err)

				// Check if account is created
				interchainAccount := icaAppB.AppKeepers.AccountKeeper.GetAccount(suite.chainB.GetContext(), icaAddr)
				suite.Require().Equal(interchainAccount.GetAddress().String(), interchainAccountAddr)

				// Create bank transfer message to execute on the host
				icaMsg = &banktypes.MsgSend{
					FromAddress: interchainAccountAddr,
					ToAddress:   suite.chainB.SenderAccount.GetAddress().String(),
					Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
				}
			}

			// ownerAddr, err := sdk.AccAddressFromBech32(owner)
			// suite.Require().NoError(err)

			msgSrv := keeper.NewMsgServerImpl(GetICAKeeper2(icaAppA))
			msg, err := types.NewMsgSubmitTx(owner, icaMsg, connectionId)
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
		owner                     string
		connectionId              string
		icaMsg                    sdk.Msg
		parseIcaAddress           bool
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"success - IBC autoTx", func() {
				registerInterchainAccount = true
				owner = TestOwnerAddress
				connectionId = path.EndpointA.ConnectionID
			}, true,
		},
		{
			"success - local autoTx", func() {
				registerInterchainAccount = true
				owner = TestOwnerAddress
				connectionId = ""
			}, true,
		},
		{
			"success - parse ICA address", func() {
				registerInterchainAccount = true
				owner = TestOwnerAddress
				connectionId = path.EndpointA.ConnectionID
				parseIcaAddress = true
			}, true,
		},
		{
			"failure - owner address is empty", func() {
				registerInterchainAccount = true
				owner = ""
				connectionId = path.EndpointA.ConnectionID
				parseIcaAddress = false
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
				interchainAccountAddr, found := GetICAApp(suite.chainA.TestChain).AppKeepers.ICAControllerKeeper.GetInterchainAccountAddress(suite.chainA.GetContext(), path.EndpointA.ConnectionID, portID)
				suite.Require().True(found)

				icaAddr, err := sdk.AccAddressFromBech32(interchainAccountAddr)
				suite.Require().NoError(err)

				// Check if account is created
				interchainAccount := icaAppB.AppKeepers.AccountKeeper.GetAccount(suite.chainB.GetContext(), icaAddr)
				suite.Require().Equal(interchainAccount.GetAddress().String(), interchainAccountAddr)
				if parseIcaAddress {
					interchainAccountAddr = "ICA_ADDR"
				}
				// Create bank transfer message to execute on the host
				icaMsg = &banktypes.MsgSend{
					FromAddress: interchainAccountAddr,
					ToAddress:   suite.chainB.SenderAccount.GetAddress().String(),
					Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
				}
			}

			GetICAApp(suite.chainA.TestChain).AppKeepers.AutoIBCTXKeeper.SetParams(suite.chainA.GetContext(), types.Params{
				AutoTxFundsCommission:      2,
				AutoTxConstantFee:          1_000_000,
				AutoTxFlexFeeMul:           100,
				RecurringAutoTxConstantFee: 1_000_000,
				MaxAutoTxDuration:          time.Hour * 24 * 366 * 10,
				MinAutoTxDuration:          time.Second * 60,
				MinAutoTxInterval:          time.Second * 20,
			})

			msg, err := types.NewMsgSubmitAutoTx(owner, "label", []sdk.Msg{icaMsg}, connectionId, "200s", "100s", 0, []uint64{})
			suite.Require().NoError(err)
			wrappedCtx := sdk.WrapSDKContext(suite.chainA.GetContext())
			msgSrv := keeper.NewMsgServerImpl(GetICAKeeper2(icaAppA))
			res, err := msgSrv.SubmitAutoTx(wrappedCtx, msg)

			if parseIcaAddress {
				icaApp := GetICAApp(suite.chainA.TestChain)
				err := icaMsg.ValidateBasic()
				suite.Require().Contains(err.Error(), "bech32")
				err = msg.ValidateBasic()
				suite.Require().NoError(err)
				suite.chainA.NextBlock()
				UnwrappedCtx := sdk.UnwrapSDKContext(wrappedCtx)
				autoTx := icaApp.AppKeepers.AutoIBCTXKeeper.GetAutoTxInfo(UnwrappedCtx, 1)
				suite.Require().NotEqual(autoTx, types.AutoTxInfo{})

				err = icaApp.AppKeepers.AutoIBCTXKeeper.SendAutoTx(UnwrappedCtx, autoTx)
				suite.Require().NoError(err)
				icaAddrToParse := "ICA_ADDR"
				autoTx = icaApp.AppKeepers.AutoIBCTXKeeper.GetAutoTxInfo(UnwrappedCtx, 1)
				//txMsgs := autoTx.GetTxMsgs()

				var txMsg sdk.Msg
				err = icaApp.AppCodec().UnpackAny(autoTx.Msgs[0], &txMsg)
				suite.Require().NoError(err)
				msgJSON, err := icaApp.AppCodec().MarshalInterfaceJSON(txMsg)
				suite.Require().NoError(err)

				msgJSONString := string(msgJSON)
				index := strings.Index(msgJSONString, icaAddrToParse)
				suite.Require().Equal(-1, index)

			}

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
