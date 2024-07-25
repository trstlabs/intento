package keeper_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/trstlabs/intento/x/intent/types"
	intenttypes "github.com/trstlabs/intento/x/intent/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"

	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
)

func (suite *KeeperTestSuite) TestOnRecvTransferPacket() {
	var (
		trace    transfertypes.DenomTrace
		amount   sdkmath.Int
		receiver string
	)

	suite.SetupTest()

	path := NewTransferPath(suite.chainA, suite.chainB)
	suite.coordinator.Setup(path)
	receiver = suite.chainB.SenderAccount.GetAddress().String() // must be explicitly changed

	amount = sdk.NewInt(100) // must be explicitly changed in malleate
	seq := uint64(1)

	trace = transfertypes.ParseDenomTrace(sdk.DefaultBondDenom)

	// send coin from chainA to chainB
	transferMsg := transfertypes.NewMsgTransfer(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sdk.NewCoin(trace.IBCDenom(), amount), suite.chainA.SenderAccount.GetAddress().String(), receiver, clienttypes.NewHeight(1, 110), 0, "")
	_, err := suite.chainA.SendMsgs(transferMsg)
	suite.Require().NoError(err) // message committed

	data := transfertypes.NewFungibleTokenPacketData(trace.GetFullDenomPath(), amount.String(), suite.chainA.SenderAccount.GetAddress().String(), receiver, "")
	packet := channeltypes.NewPacket(data.GetBytes(), seq, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, clienttypes.NewHeight(1, 100), 0)

	ack := suite.chainB.GetIntoApp().TransferStack.OnRecvPacket(suite.chainB.GetContext(), packet, suite.chainA.SenderAccount.GetAddress())

	suite.Require().True(ack.Success())

}

func (suite *KeeperTestSuite) TestOnRecvTransferPacketWithAction() {
	suite.SetupTest()

	params := intenttypes.DefaultParams()
	params.GasFeeCoins = sdk.NewCoins(sdk.NewCoin("stake", sdk.OneInt()))
	params.ActionFlexFeeMul = 1
	suite.chainA.GetIntoApp().IntentKeeper.SetParams(suite.chainA.GetContext(), params)

	addr := suite.chainA.SenderAccount.GetAddress()
	msg := `{
		"@type":"/cosmos.bank.v1beta1.MsgSend",
		"amount": [{
			"amount": "70",
			"denom": "stake"
		}],
		"from_address": "into12gxmzpucje8aflw2vz45rv8x4nyaaj3rp8vjh03dulehkdl5fu6s93ewkp",
		"to_address": "into1ykql5ktedxkpjszj5trzu8f5dxajvgv95nuwjx"
	}`

	ackBytes := suite.receiveTransferPacket(addr.String(), fmt.Sprintf(`{"action": {"owner": "%s","label": "my_trigger", "msgs": [%s], "duration": "500s", "interval": "60s", "start_at": "0"} }`, addr, msg))

	var ack map[string]string // This can't be unmarshalled to Acknowledgement because it's fetched from the events
	err := json.Unmarshal(ackBytes, &ack)
	suite.Require().NoError(err)
	suite.Require().NotContains(ack, "error")

	action := suite.chainA.GetIntoApp().IntentKeeper.GetActionInfo(suite.chainA.GetContext(), 1)

	suite.Require().Equal(action.Owner, addr.String())
	suite.Require().Equal(action.Label, "my_trigger")
	suite.Require().Equal(action.ICAConfig.PortID, "")
	suite.Require().Equal(action.Interval, time.Second*60)

	var txMsgAny codectypes.Any
	cdc := codec.NewProtoCodec(suite.chainA.GetIntoApp().InterfaceRegistry())

	err = cdc.UnmarshalJSON([]byte(msg), &txMsgAny)
	suite.Require().NoError(err)
	suite.True(action.Msgs[0].Equal(txMsgAny))
}

func (suite *KeeperTestSuite) TestOnRecvTransferPacketAndMultippleActions() {
	suite.SetupTest()

	params := intenttypes.DefaultParams()
	params.GasFeeCoins = sdk.NewCoins(sdk.NewCoin("stake", sdk.OneInt()))
	params.ActionFlexFeeMul = 1
	suite.chainA.GetIntoApp().IntentKeeper.SetParams(suite.chainA.GetContext(), params)

	addr := suite.chainA.SenderAccount.GetAddress()
	msg := `{
		"@type":"/cosmos.bank.v1beta1.MsgSend",
		"amount": [{
			"amount": "70",
			"denom": "stake"
		}],
		"from_address": "into12gxmzpucje8aflw2vz45rv8x4nyaaj3rp8vjh03dulehkdl5fu6s93ewkp",
		"to_address": "into1ykql5ktedxkpjszj5trzu8f5dxajvgv95nuwjx"
	}`

	path := NewICAPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupConnections(path)
	err := SetupICAPath(path, addr.String())
	suite.Require().NoError(err)

	//chainB sends packet to chainA. connectionID to execute on chainB is on chainAs config
	ackBytes := suite.receiveTransferPacket(addr.String(), fmt.Sprintf(`{"action": {"owner": "%s","label": "my_trigger", "cid":"%s", "host_cid":"%s","msgs": [%s, %s], "duration": "500s", "interval": "60s", "start_at": "0", "fallback": "true" } }`, addr.String(), path.EndpointA.ConnectionID, path.EndpointB.ConnectionID, msg, msg))

	var ack map[string]string // This can't be unmarshalled to Acknowledgement because it's fetched from the events
	err = json.Unmarshal(ackBytes, &ack)
	suite.Require().NoError(err)
	suite.Require().NotContains(ack, "error")

	action := suite.chainA.GetIntoApp().IntentKeeper.GetActionInfo(suite.chainA.GetContext(), 1)

	suite.Require().Equal(action.Owner, addr.String())
	suite.Require().Equal(action.Label, "my_trigger")
	suite.Require().Equal(action.ICAConfig.PortID, "icacontroller-"+addr.String())
	suite.Require().Equal(action.ICAConfig.ConnectionID, path.EndpointA.ConnectionID)

	suite.Require().Equal(action.Interval, time.Second*60)

	_, found := suite.chainA.GetIntoApp().ICAControllerKeeper.GetInterchainAccountAddress(suite.chainA.GetContext(), action.ICAConfig.ConnectionID, action.ICAConfig.PortID)
	suite.Require().True(found)

	var txMsgAny codectypes.Any
	cdc := codec.NewProtoCodec(suite.chainA.GetIntoApp().InterfaceRegistry())

	err = cdc.UnmarshalJSON([]byte(msg), &txMsgAny)
	suite.Require().NoError(err)
	suite.True(action.Msgs[0].Equal(txMsgAny))
}

func (suite *KeeperTestSuite) TestOnRecvTransferPacketSubmitTxAndAddressParsing() {
	suite.SetupTest()

	params := intenttypes.DefaultParams()
	params.GasFeeCoins = sdk.NewCoins(sdk.NewCoin("stake", sdk.OneInt()))
	params.ActionFlexFeeMul = 1
	suite.chainA.GetIntoApp().IntentKeeper.SetParams(suite.chainA.GetContext(), params)

	addr := suite.chainA.SenderAccount.GetAddress()
	msg := `{
		"@type":"/cosmos.bank.v1beta1.MsgSend",
		"amount": [{
			"amount": "70",
			"denom": "stake"
		}],
		"from_address": "ICA_ADDR",
		"to_address": "into1ykql5ktedxkpjszj5trzu8f5dxajvgv95nuwjx"
	}`

	path := NewICAPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupConnections(path)
	err := SetupICAPath(path, addr.String())
	suite.Require().NoError(err)

	ackBytes := suite.receiveTransferPacket(addr.String(), fmt.Sprintf(`{"action": {"owner": "%s","label": "my trigger", "cid":"%s","host_cid":"%s","msgs": [%s, %s], "duration": "120s", "interval": "60s", "start_at": "0", "fallback": "true" }}`, addr.String(), path.EndpointA.ConnectionID, path.EndpointB.ConnectionID, msg, msg))
	var ack map[string]string // This can't be unmarshalled to Acknowledgement because it's fetched from the events
	err = json.Unmarshal(ackBytes, &ack)
	suite.Require().NoError(err)
	suite.Require().NotContains(ack, "error")

	actionKeeper := suite.chainA.GetIntoApp().IntentKeeper
	action := actionKeeper.GetActionInfo(suite.chainA.GetContext(), 1)
	unpacker := suite.chainA.Codec
	unpackedMsgs := action.GetTxMsgs(unpacker)
	suite.Require().True(strings.Contains(unpackedMsgs[0].String(), types.ParseICAValue))

	suite.chainA.CurrentHeader.Time = suite.chainA.CurrentHeader.Time.Add(time.Minute)
	FakeBeginBlocker(suite.chainA.GetContext(), actionKeeper, sdk.ConsAddress(suite.chainA.Vals.Proposer.Address))

	action = actionKeeper.GetActionInfo(suite.chainA.GetContext(), 1)
	actionHistory, _ := actionKeeper.GetActionHistory(suite.chainA.GetContext(), action.ID)
	suite.Require().NotNil(actionHistory.History)
	suite.Require().Empty(actionHistory.History[0].Errors)
	suite.Require().Equal(action.Owner, addr.String())
	suite.Require().Equal(action.Label, "my trigger")
	suite.Require().Equal(action.ICAConfig.PortID, "icacontroller-"+addr.String())
	suite.Require().Equal(action.ICAConfig.ConnectionID, path.EndpointA.ConnectionID)

	unpackedMsgs = action.GetTxMsgs(unpacker)
	suite.Require().False(strings.Contains(unpackedMsgs[0].String(), types.ParseICAValue))
	suite.Require().Equal(action.Interval, time.Second*60)
}

func (suite *KeeperTestSuite) TestOnRecvTransferPacketSubmitTxWithSentDenomInParams() {
	suite.SetupTest()

	addr := suite.chainA.SenderAccount.GetAddress()
	msg := `{
		"@type":"/cosmos.bank.v1beta1.MsgSend",
		"amount": [{
			"amount": "70",
			"denom": "stake"
		}],
		"from_address": "ICA_ADDR",
		"to_address": "into1ykql5ktedxkpjszj5trzu8f5dxajvgv95nuwjx"
	}`

	path := NewICAPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupConnections(path)
	err := SetupICAPath(path, addr.String())
	suite.Require().NoError(err)

	ackBytes := suite.receiveTransferPacket(addr.String(), fmt.Sprintf(`{"action": {"owner": "%s","label": "my trigger", "cid":"%s","host_cid":"%s","msgs": [%s, %s], "duration": "120s", "interval": "60s", "start_at": "0", "fallback": "true" }}`, addr.String(), path.EndpointA.ConnectionID, path.EndpointB.ConnectionID, msg, msg))
	var ack map[string]string // This can't be unmarshalled to Acknowledgement because it's fetched from the events
	err = json.Unmarshal(ackBytes, &ack)
	suite.Require().NoError(err)
	suite.Require().NotContains(ack, "error")

	actionKeeper := suite.chainA.GetIntoApp().IntentKeeper
	action := actionKeeper.GetActionInfo(suite.chainA.GetContext(), 1)
	feeAddr, _ := sdk.AccAddressFromBech32(action.FeeAddress)
	bDenom := suite.chainA.GetIntoApp().BankKeeper.GetAllBalances(suite.chainA.GetContext(), feeAddr)[0].Denom
	params := intenttypes.DefaultParams()
	params.GasFeeCoins = sdk.NewCoins(sdk.NewCoin(bDenom, sdk.NewInt(2)), sdk.NewCoin("stake", sdk.OneInt()))
	params.ActionFlexFeeMul = 1
	suite.chainA.GetIntoApp().IntentKeeper.SetParams(suite.chainA.GetContext(), params)

	unpacker := suite.chainA.Codec
	unpackedMsgs := action.GetTxMsgs(unpacker)
	suite.Require().True(strings.Contains(unpackedMsgs[0].String(), types.ParseICAValue))

	suite.chainA.CurrentHeader.Time = suite.chainA.CurrentHeader.Time.Add(time.Minute)
	FakeBeginBlocker(suite.chainA.GetContext(), actionKeeper, sdk.ConsAddress(suite.chainA.Vals.Proposer.Address))

	action = actionKeeper.GetActionInfo(suite.chainA.GetContext(), 1)
	actionHistory, _ := actionKeeper.GetActionHistory(suite.chainA.GetContext(), action.ID)
	suite.Require().NotNil(actionHistory.History)
	suite.Require().Empty(actionHistory.History[0].Errors)
}
