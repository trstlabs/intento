package keeper_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/trstlabs/intento/x/intent/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"

	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
)

func (suite *KeeperTestSuite) TestOnRecvTransferPacket() {
	var (
		trace    transfertypes.DenomTrace
		amount   math.Int
		receiver string
	)

	suite.SetupTest()

	path := NewTransferPath(suite.IntentoChain, suite.HostChain)
	suite.Coordinator.Setup(path)
	receiver = suite.HostChain.SenderAccount.GetAddress().String() // must be explicitly changed

	amount = math.NewInt(100) // must be explicitly changed in malleate
	seq := uint64(1)

	trace = transfertypes.ParseDenomTrace(sdk.DefaultBondDenom)

	// send coin from IntentoChain to HostChain
	transferMsg := transfertypes.NewMsgTransfer(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sdk.NewCoin(trace.IBCDenom(), amount), suite.IntentoChain.SenderAccount.GetAddress().String(), receiver, clienttypes.NewHeight(1, 110), 0, "")
	_, err := suite.IntentoChain.SendMsgs(transferMsg)
	suite.Require().NoError(err) // message committed

	data := transfertypes.NewFungibleTokenPacketData(trace.GetFullDenomPath(), amount.String(), suite.IntentoChain.SenderAccount.GetAddress().String(), receiver, "")
	packet := channeltypes.NewPacket(data.GetBytes(), seq, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, clienttypes.NewHeight(1, 100), 0)

	//a little hack as this check would be on HostChain OnRecvPacket
	ack := GetICAApp(suite.IntentoChain).TransferStack.OnRecvPacket(suite.IntentoChain.GetContext(), packet, suite.IntentoChain.SenderAccount.GetAddress())

	suite.Require().True(ack.Success())

}

func (suite *KeeperTestSuite) TestOnRecvTransferPacketWithFlow() {
	suite.SetupTest()

	params := types.DefaultParams()
	params.GasFeeCoins = sdk.NewCoins(sdk.NewCoin("stake", math.OneInt()))
	params.FlowFlexFeeMul = 1
	GetICAApp(suite.IntentoChain).IntentKeeper.SetParams(suite.IntentoChain.GetContext(), params)

	addr := suite.IntentoChain.SenderAccount.GetAddress().String()
	addrTo := suite.TestAccs[0].String()
	msg := fmt.Sprintf(`{
		"@type":"/cosmos.bank.v1beta1.MsgSend",
		"amount": [{
			"amount": "70",
			"denom": "stake"
		}],
		"from_address": "%s",
		"to_address": "%s"
	}`, addr, addrTo)

	ackBytes := suite.receiveTransferPacket(addr, fmt.Sprintf(`{"flow": {"owner": "%s","label": "my_trigger", "msgs": [%s], "duration": "500s", "interval": "60s", "start_at": "0"} }`, addr, msg))

	var ack map[string]string // This can't be unmarshalled to Acknowledgement because it's fetched from the events
	err := json.Unmarshal(ackBytes, &ack)
	suite.Require().NoError(err)
	suite.Require().NotContains(ack, "error")

	flow := GetICAApp(suite.IntentoChain).IntentKeeper.GetFlowInfo(suite.IntentoChain.GetContext(), 1)

	suite.Require().Equal(flow.Owner, addr)
	suite.Require().Equal(flow.Label, "my_trigger")
	suite.Require().Equal(flow.ICAConfig.PortID, "")
	suite.Require().Equal(flow.Interval, time.Second*60)

	var txMsgAny codectypes.Any
	cdc := codec.NewProtoCodec(GetICAApp(suite.IntentoChain).InterfaceRegistry())

	err = cdc.UnmarshalJSON([]byte(msg), &txMsgAny)
	suite.Require().NoError(err)
	suite.True(flow.Msgs[0].Equal(txMsgAny))
}

func (suite *KeeperTestSuite) TestOnRecvTransferPacketAndMultipleMsgs() {
	suite.SetupTest()

	params := types.DefaultParams()
	params.GasFeeCoins = sdk.NewCoins(sdk.NewCoin("stake", math.OneInt()))
	params.FlowFlexFeeMul = 1
	GetICAApp(suite.IntentoChain).IntentKeeper.SetParams(suite.IntentoChain.GetContext(), params)

	addr := suite.IntentoChain.SenderAccount.GetAddress()
	msg := `{
		"@type":"/cosmos.bank.v1beta1.MsgSend",
		"amount": [{
			"amount": "70",
			"denom": "stake"
		}],
		"from_address": "into12gxmzpucje8aflw2vz45rv8x4nyaaj3rp8vjh03dulehkdl5fu6s93ewkp",
		"to_address": "into1ykql5ktedxkpjszj5trzu8f5dxajvgv95nuwjx"
	}`

	path := NewICAPath(suite.IntentoChain, suite.HostChain)
	suite.Coordinator.SetupConnections(path)
	err := suite.SetupICAPath(path, addr.String())
	suite.Require().NoError(err)

	//HostChain sends packet to IntentoChain. connectionID to execute on HostChain is on IntentoChains config
	ackBytes := suite.receiveTransferPacket(addr.String(), fmt.Sprintf(`{"flow": {"owner": "%s","label": "my_trigger", "cid":"%s", "host_cid":"%s","msgs": [%s, %s], "duration": "500s", "interval": "60s", "start_at": "0", "fallback": "true" } }`, addr.String(), path.EndpointA.ConnectionID, path.EndpointB.ConnectionID, msg, msg))

	var ack map[string]string // This can't be unmarshalled to Acknowledgement because it's fetched from the events
	err = json.Unmarshal(ackBytes, &ack)
	suite.Require().NoError(err)
	suite.Require().NotContains(ack, "error")

	flow := GetICAApp(suite.IntentoChain).IntentKeeper.GetFlowInfo(suite.IntentoChain.GetContext(), 1)

	suite.Require().Equal(flow.Owner, addr.String())
	suite.Require().Equal(flow.Label, "my_trigger")
	suite.Require().Equal(flow.Configuration.FallbackToOwnerBalance, true)
	suite.Require().Equal(flow.ICAConfig.PortID, "icacontroller-"+addr.String())
	suite.Require().Equal(flow.ICAConfig.ConnectionID, path.EndpointA.ConnectionID)

	suite.Require().Equal(flow.Interval, time.Second*60)

	_, found := GetICAApp(suite.IntentoChain).ICAControllerKeeper.GetInterchainAccountAddress(suite.IntentoChain.GetContext(), flow.ICAConfig.ConnectionID, flow.ICAConfig.PortID)
	suite.Require().True(found)

	var txMsgAny codectypes.Any
	cdc := codec.NewProtoCodec(GetICAApp(suite.IntentoChain).InterfaceRegistry())

	err = cdc.UnmarshalJSON([]byte(msg), &txMsgAny)
	suite.Require().NoError(err)
	suite.True(flow.Msgs[0].Equal(txMsgAny))
}

func (suite *KeeperTestSuite) TestOnRecvTransferPacketSubmitTxAndAddressParsing() {
	suite.SetupTest()

	params := types.DefaultParams()
	params.GasFeeCoins = sdk.NewCoins(sdk.NewCoin("stake", math.OneInt()))
	params.FlowFlexFeeMul = 1
	GetICAApp(suite.IntentoChain).IntentKeeper.SetParams(suite.IntentoChain.GetContext(), params)

	addr := suite.IntentoChain.SenderAccount.GetAddress()
	msg := `{
		"@type":"/cosmos.bank.v1beta1.MsgSend",
		"amount": [{
			"amount": "70",
			"denom": "stake"
		}],
		"from_address": "ICA_ADDR",
		"to_address": "into1ykql5ktedxkpjszj5trzu8f5dxajvgv95nuwjx"
	}`

	path := NewICAPath(suite.IntentoChain, suite.HostChain)
	suite.Coordinator.SetupConnections(path)
	err := suite.SetupICAPath(path, addr.String())
	suite.Require().NoError(err)

	ackBytes := suite.receiveTransferPacket(addr.String(), fmt.Sprintf(`{"flow": {"owner": "%s","label": "my trigger", "cid":"%s","host_cid":"%s","msgs": [%s, %s], "duration": "120s", "interval": "60s", "start_at": "0", "fallback":"true" }}`, addr.String(), path.EndpointA.ConnectionID, path.EndpointB.ConnectionID, msg, msg))
	var ack map[string]string // This can't be unmarshalled to Acknowledgement because it's fetched from the events
	err = json.Unmarshal(ackBytes, &ack)
	suite.Require().NoError(err)
	suite.Require().NotContains(ack, "error")

	flowKeeper := GetICAApp(suite.IntentoChain).IntentKeeper
	flow := flowKeeper.GetFlowInfo(suite.IntentoChain.GetContext(), 1)
	unpacker := suite.IntentoChain.Codec
	unpackedMsgs := flow.GetTxMsgs(unpacker)
	suite.Require().True(strings.Contains(unpackedMsgs[0].String(), types.ParseICAValue))

	suite.IntentoChain.CurrentHeader.Time = suite.IntentoChain.CurrentHeader.Time.Add(time.Minute)
	flowKeeper.HandleFlow(suite.IntentoChain.GetContext(), flowKeeper.Logger(suite.IntentoChain.GetContext()), flow, suite.IntentoChain.GetContext().BlockTime(), nil)

	flow = flowKeeper.GetFlowInfo(suite.IntentoChain.GetContext(), 1)
	flowHistory, _ := flowKeeper.GetFlowHistory(suite.IntentoChain.GetContext(), flow.ID)
	suite.Require().NotNil(flowHistory)
	suite.Require().Empty(flowHistory[0].Errors)
	suite.Require().Equal(flow.Owner, addr.String())
	suite.Require().Equal(flow.Label, "my trigger")
	suite.Require().Equal(flow.ICAConfig.PortID, "icacontroller-"+addr.String())
	suite.Require().Equal(flow.ICAConfig.ConnectionID, path.EndpointA.ConnectionID)

	unpackedMsgs = flow.GetTxMsgs(unpacker)
	suite.Require().False(strings.Contains(unpackedMsgs[0].String(), types.ParseICAValue))
	suite.Require().Equal(flow.Interval, time.Second*60)
}

func (suite *KeeperTestSuite) TestOnRecvTransferPacketSubmitTxWithSentDenomInParams() {
	suite.SetupTest()

	addr := suite.IntentoChain.SenderAccount.GetAddress()
	msg := `{
		"@type":"/cosmos.bank.v1beta1.MsgSend",
		"amount": [{
			"amount": "70",
			"denom": "stake"
		}],
		"from_address": "ICA_ADDR",
		"to_address": "into1ykql5ktedxkpjszj5trzu8f5dxajvgv95nuwjx"
	}`

	path := NewICAPath(suite.IntentoChain, suite.HostChain)
	suite.Coordinator.SetupConnections(path)
	err := suite.SetupICAPath(path, addr.String())
	suite.Require().NoError(err)

	ackBytes := suite.receiveTransferPacket(addr.String(), fmt.Sprintf(`{"flow": {"owner": "%s","label": "my trigger", "cid":"%s","host_cid":"%s","msgs": [%s, %s], "duration": "120s", "interval": "60s", "start_at": "0", "fallback": "true" }}`, addr.String(), path.EndpointA.ConnectionID, path.EndpointB.ConnectionID, msg, msg))
	var ack map[string]string // This can't be unmarshalled to Acknowledgement because it's fetched from the events
	err = json.Unmarshal(ackBytes, &ack)
	suite.Require().NoError(err)
	suite.Require().NotContains(ack, "error")

	flowKeeper := GetICAApp(suite.IntentoChain).IntentKeeper
	flow := flowKeeper.GetFlowInfo(suite.IntentoChain.GetContext(), 1)
	feeAddr, _ := sdk.AccAddressFromBech32(flow.FeeAddress)
	bDenom := GetICAApp(suite.IntentoChain).BankKeeper.GetAllBalances(suite.IntentoChain.GetContext(), feeAddr)[0].Denom
	params := types.DefaultParams()
	params.GasFeeCoins = sdk.NewCoins(sdk.NewCoin(bDenom, math.NewInt(2)), sdk.NewCoin("stake", math.OneInt()))
	params.FlowFlexFeeMul = 1
	GetICAApp(suite.IntentoChain).IntentKeeper.SetParams(suite.IntentoChain.GetContext(), params)

	unpacker := suite.IntentoChain.Codec
	unpackedMsgs := flow.GetTxMsgs(unpacker)
	suite.Require().True(strings.Contains(unpackedMsgs[0].String(), types.ParseICAValue))

	suite.IntentoChain.CurrentHeader.Time = suite.IntentoChain.CurrentHeader.Time.Add(time.Minute)
	flowKeeper.HandleFlow(suite.IntentoChain.GetContext(), flowKeeper.Logger(suite.IntentoChain.GetContext()), flow, suite.IntentoChain.GetContext().BlockTime(), nil)

	flow = flowKeeper.GetFlowInfo(suite.IntentoChain.GetContext(), 1)
	flowHistory, _ := flowKeeper.GetFlowHistory(suite.IntentoChain.GetContext(), flow.ID)
	suite.Require().NotNil(flowHistory)
	suite.Require().Empty(flowHistory[0].Errors)
}

func (suite *KeeperTestSuite) TestOnRecvTransferPacketFlowWithConditions() {
	suite.SetupTest()

	params := types.DefaultParams()
	params.GasFeeCoins = sdk.NewCoins(sdk.NewCoin("stake", math.OneInt()))
	params.FlowFlexFeeMul = 1
	GetICAApp(suite.IntentoChain).IntentKeeper.SetParams(suite.IntentoChain.GetContext(), params)

	addr := suite.IntentoChain.SenderAccount.GetAddress().String()
	addrTo := suite.TestAccs[0].String()
	msg := fmt.Sprintf(`{
		"@type":"/cosmos.bank.v1beta1.MsgSend",
		"amount": [{
			"amount": "70",
			"denom": "stake"
		}],
		"from_address": "%s",
		"to_address": "%s"
	}`, addr, addrTo)

	ackBytes := suite.receiveTransferPacket(addr, fmt.Sprintf(`{"flow": {"owner": "%s","label": "my_trigger", "msgs": [%s], "duration": "500s", "interval": "60s", "start_at": "0","conditions":{"stop_on_failure_of": [12345], "feedback_loops": [{"response_index":0,"response_key": "Amount.[0]", "msgs_index":1, "msg_key":"Amount","value_type": "sdk.Coin"}], "comparisons": [{"response_index":0,"response_key": "Amount.[0]", "operand":"1'$HOST_DENOM'", "operator":4,"value_type": "sdk.Coin"}]}}}`, addr, msg))

	var ack map[string]string // This can't be unmarshalled to Acknowledgement because it's fetched from the events
	err := json.Unmarshal(ackBytes, &ack)
	suite.Require().NoError(err)
	suite.Require().NotContains(ack, "error")

	flow := GetICAApp(suite.IntentoChain).IntentKeeper.GetFlowInfo(suite.IntentoChain.GetContext(), 1)

	suite.Require().Equal(flow.Owner, addr)
	suite.Require().Equal(flow.Label, "my_trigger")
	suite.Require().Equal(flow.ICAConfig.PortID, "")
	suite.Require().Equal(flow.Interval, time.Second*60)
	suite.Require().Equal(flow.Conditions.StopOnFailureOf[0], uint64(12345))
	suite.Require().Equal(flow.Conditions.FeedbackLoops[0].MsgsIndex, uint32(1))
	suite.Require().Equal(flow.Conditions.Comparisons[0].Operator, types.ComparisonOperator(4))

	var txMsgAny codectypes.Any
	cdc := codec.NewProtoCodec(GetICAApp(suite.IntentoChain).InterfaceRegistry())

	err = cdc.UnmarshalJSON([]byte(msg), &txMsgAny)
	suite.Require().NoError(err)
	suite.True(flow.Msgs[0].Equal(txMsgAny))
}
