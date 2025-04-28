package keeper_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"cosmossdk.io/math"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/trstlabs/intento/x/intent/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	ibctesting "github.com/cosmos/ibc-go/v8/testing"

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

func (suite *KeeperTestSuite) TestOnRecvTransferPacketWithSubmitFlow() {
	suite.SetupTest()

	addr := suite.HostChain.SenderAccount.GetAddress().String()
	addrTo := suite.TestAccs[0].String()
	msg := fmt.Sprintf(`{
		"@type":"/cosmos.bank.v1beta1.MsgSend",
		"amount": [{
			"amount": "70",
			"denom": "stake"
		}],
		"from_address": "%s",
		"to_address": "%s"
	}`, derivePlaceholderSender(ibctesting.FirstChannelID, addr).String(), addrTo)

	ackBytes := suite.receiveTransferPacket(addr, fmt.Sprintf(`{"flow": {"owner": "%s","label": "my flow", "msgs": [%s], "duration": "500s", "interval": "60s", "start_at": "0"} }`, addr, msg))

	var ack map[string]string // This can't be unmarshalled to Acknowledgement because it's fetched from the events
	err := json.Unmarshal(ackBytes, &ack)
	suite.Require().NoError(err)
	suite.Require().NotContains(ack, "error")

	flow := GetICAApp(suite.IntentoChain).IntentKeeper.GetFlowInfo(suite.IntentoChain.GetContext(), 1)

	suite.Require().Equal(flow.Label, "my flow")
	suite.Require().Equal(flow.ICAConfig.PortID, "")
	suite.Require().Equal(flow.Interval, time.Second*60)

	var txMsgAny codectypes.Any
	cdc := codec.NewProtoCodec(GetICAApp(suite.IntentoChain).InterfaceRegistry())

	err = cdc.UnmarshalJSON([]byte(msg), &txMsgAny)
	suite.Require().NoError(err)
	suite.True(flow.Msgs[0].Equal(txMsgAny))
}

func (suite *KeeperTestSuite) TestOnRecvTransferPacketWithUpdateFlow() {
	suite.SetupTest()

	addr := suite.HostChain.SenderAccount.GetAddress().String()
	addrTo := suite.TestAccs[0].String()
	msg := fmt.Sprintf(`{
		"@type":"/cosmos.bank.v1beta1.MsgSend",
		"amount": [{
			"amount": "70",
			"denom": "stake"
		}],
		"from_address": "%s",
		"to_address": "%s"
	}`, derivePlaceholderSender(ibctesting.FirstChannelID, addr).String(), addrTo)

	ackBytes := suite.receiveTransferPacket(addr, fmt.Sprintf(`{"flow": {"owner": "%s","label": "my flowwwwww", "msgs": [%s], "duration": "500s", "interval": "60s", "start_at": "0"} }`, addr, msg))

	var ack map[string]string // This can't be unmarshalled to Acknowledgement because it's fetched from the events
	err := json.Unmarshal(ackBytes, &ack)
	suite.Require().NoError(err)
	suite.Require().NotContains(ack, "error")
	flow := GetICAApp(suite.IntentoChain).IntentKeeper.GetFlowInfo(suite.IntentoChain.GetContext(), 1)
	suite.Require().Equal(flow.Owner, derivePlaceholderSender(ibctesting.FirstChannelID, addr).String())
	ackBytes = suite.receiveTransferPacketWithSequence(addr, fmt.Sprintf(`{"flow": {"owner": "%s","id": "1","label": "my flow", "msgs": [%s], "duration": "500s", "interval": "60s", "start_at": "0"} }`, derivePlaceholderSender(ibctesting.FirstChannelID, addr).String(), msg), 1)

	err = json.Unmarshal(ackBytes, &ack)
	suite.Require().NoError(err)
	suite.Require().NotContains(ack, "error")

	flow = GetICAApp(suite.IntentoChain).IntentKeeper.GetFlowInfo(suite.IntentoChain.GetContext(), 1)

	suite.Require().Equal(flow.Label, "my flow")
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

	addr := suite.HostChain.SenderAccount.GetAddress()
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

	path := NewICAPath(suite.IntentoChain, suite.HostChain)
	suite.Coordinator.SetupConnections(path)

	err := GetICAApp(suite.IntentoChain).BankKeeper.SendCoins(suite.IntentoChain.GetContext(), suite.IntentoChain.SenderAccount.GetAddress(), suite.HostChain.SenderAccount.GetAddress(), sdk.NewCoins(sdk.NewInt64Coin("stake", 100000)))
	suite.Require().NoError(err)

	//fix in test as we do not have the channels over the same connectionID in testing
	path.EndpointA.ConnectionID = "connection-1"
	ackBytes := suite.receiveTransferPacket(addr.String(), fmt.Sprintf(`{"flow": {"owner": "%s","label": "my flow", "cid":"%s", "host_cid":"%s","msgs": [%s, %s], "duration": "500s", "interval": "60s", "start_at": "0", "fallback": "true" } }`, addr.String(), path.EndpointA.ConnectionID, path.EndpointB.ConnectionID, msg, msg))

	var ack map[string]string // This can't be unmarshalled to Acknowledgement because it's fetched from the events
	err = json.Unmarshal(ackBytes, &ack)
	suite.Require().NoError(err)
	suite.Require().NotContains(ack, "error")

	flow := GetICAApp(suite.IntentoChain).IntentKeeper.GetFlowInfo(suite.IntentoChain.GetContext(), 1)

	suite.Require().Equal(flow.Owner, addr.String())
	suite.Require().Equal(flow.Label, "my flow")
	suite.Require().Equal(flow.Configuration.FallbackToOwnerBalance, true)
	suite.Require().Equal(flow.ICAConfig.PortID, "icacontroller-"+addr.String())
	suite.Require().Equal(flow.ICAConfig.ConnectionID, path.EndpointA.ConnectionID)

	suite.Require().Equal(flow.Interval, time.Second*60)

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
	GetICAApp(suite.IntentoChain).IntentKeeper.SetParams(suite.IntentoChain.GetContext(), params)
	types.Denom = "stake"

	addr := suite.HostChain.SenderAccount.GetAddress()
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
	err = GetICAApp(suite.IntentoChain).BankKeeper.SendCoins(suite.IntentoChain.GetContext(), suite.IntentoChain.SenderAccount.GetAddress(), suite.HostChain.SenderAccount.GetAddress(), sdk.NewCoins(sdk.NewInt64Coin("stake", 100000)))
	suite.Require().NoError(err)

	path.EndpointA.ConnectionID = "connection-1"
	ackBytes := suite.receiveTransferPacket(addr.String(), fmt.Sprintf(`{"flow": {"owner": "%s","label": "my flow", "cid":"%s","host_cid":"%s","msgs": [%s, %s], "duration": "120s", "interval": "60s", "start_at": "0", "fallback":"true" }}`, addr.String(), path.EndpointA.ConnectionID, path.EndpointB.ConnectionID, msg, msg))
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

	//fix as we do not have the channels over the same connectionID in testing
	flow.ICAConfig.ConnectionID = "connection-0"
	flowKeeper.HandleFlow(suite.IntentoChain.GetContext(), flowKeeper.Logger(suite.IntentoChain.GetContext()), flow, suite.IntentoChain.GetContext().BlockTime(), nil)

	flow = flowKeeper.GetFlowInfo(suite.IntentoChain.GetContext(), 1)
	flowHistory, _ := flowKeeper.GetFlowHistory(suite.IntentoChain.GetContext(), flow.ID)
	suite.Require().NotNil(flowHistory)
	suite.Require().Empty(flowHistory[0].Errors)
	suite.Require().Equal(flow.Owner, addr.String())
	suite.Require().Equal(flow.Label, "my flow")
	suite.Require().Equal(flow.ICAConfig.PortID, "icacontroller-"+addr.String())

	unpackedMsgs = flow.GetTxMsgs(unpacker)
	suite.Require().False(strings.Contains(unpackedMsgs[0].String(), types.ParseICAValue))
	suite.Require().Equal(flow.Interval, time.Second*60)
}

func (suite *KeeperTestSuite) TestOnRecvTransferPacketFlowWithConditionsAndDerivedSender() {
	suite.SetupTest()

	addr := suite.HostChain.SenderAccount.GetAddress().String()
	addrTo := suite.TestAccs[0].String()
	msg := fmt.Sprintf(`{
		"@type":"/cosmos.bank.v1beta1.MsgSend",
		"amount": [{
			"amount": "70",
			"denom": "stake"
		}],
		"from_address": "%s",
		"to_address": "%s"
	}`, derivePlaceholderSender(ibctesting.FirstChannelID, addr), addrTo)

	ackBytes := suite.receiveTransferPacket(addr, fmt.Sprintf(`{"flow": {"label": "my flow", "msgs": [%s], "duration": "500s", "interval": "60s", "start_at": "0","conditions":{"stop_on_failure_of": [12345], "feedback_loops": [{"response_index":0,"response_key": "Amount.[0]", "msgs_index":1, "msg_key":"Amount","value_type": "sdk.Coin"}], "comparisons": [{"response_index":0,"response_key": "Amount.[0]", "operand":"1'$HOST_DENOM'", "operator":4,"value_type": "sdk.Coin"}]}}}`, msg))

	var ack map[string]string // This can't be unmarshalled to Acknowledgement because it's fetched from the events
	err := json.Unmarshal(ackBytes, &ack)
	suite.Require().NoError(err)
	suite.Require().NotContains(ack, "error")

	flow := GetICAApp(suite.IntentoChain).IntentKeeper.GetFlowInfo(suite.IntentoChain.GetContext(), 1)

	suite.Require().Equal(flow.Owner, derivePlaceholderSender(ibctesting.FirstChannelID, addr).String())
	suite.Require().Equal(flow.Label, "my flow")
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

func (suite *KeeperTestSuite) TestOnRecvTransferPacketSubmitTxAndAddressParsingMsgExecuteContract() {
	suite.SetupTest()
	params := types.DefaultParams()
	params.GasFeeCoins = sdk.NewCoins(sdk.NewCoin("stake", math.OneInt()))
	GetICAApp(suite.IntentoChain).IntentKeeper.SetParams(suite.IntentoChain.GetContext(), params)
	types.Denom = "stake"

	addr := suite.HostChain.SenderAccount.GetAddress()
	// 1) Build inner contract msg and Base64‑encode it
	contractMsg := map[string]interface{}{
		"transfer": map[string]string{
			"recipient": "ICA_ADDR",
			"amount":    "1000",
		},
	}
	contractMsgJSON, _ := json.Marshal(contractMsg)
	contractMsgBase64 := base64.StdEncoding.EncodeToString(contractMsgJSON)

	// 2) Embed raw Base64 in the JSON for MsgExecuteContract
	msg := fmt.Sprintf(`{
        "@type":   "/cosmwasm.wasm.v1.MsgExecuteContract",
        "sender":   "ICA_ADDR",
        "contract": "wasm1deadbeefdeadbeefdeadbeefdeadbeefdead00",
        "msg":      "%s",
        "funds":   [{"denom":"stake","amount":"1"}]
    }`, contractMsgBase64)

	// Set up IBC/ICA path, fund accounts, receive the packet...
	path := NewICAPath(suite.IntentoChain, suite.HostChain)
	suite.Coordinator.SetupConnections(path)
	suite.Require().NoError(suite.SetupICAPath(path, addr.String()))
	suite.Require().NoError(
		GetICAApp(suite.IntentoChain).BankKeeper.SendCoins(
			suite.IntentoChain.GetContext(),
			suite.IntentoChain.SenderAccount.GetAddress(),
			suite.HostChain.SenderAccount.GetAddress(),
			sdk.NewCoins(sdk.NewInt64Coin("stake", 10_000_000)),
		),
	)

	path.EndpointA.ConnectionID = "connection-1"
	ackBytes := suite.receiveTransferPacket(
		addr.String(),
		fmt.Sprintf(
			`{"flow":{"owner":"%s","label":"my flow","cid":"%s","host_cid":"%s","msgs":[%s,%s],"duration":"120s","interval":"60s","start_at":"0","fallback":"true"}}`,
			addr.String(), path.EndpointA.ConnectionID, path.EndpointB.ConnectionID, msg, msg,
		),
	)
	var ack map[string]string
	suite.Require().NoError(json.Unmarshal(ackBytes, &ack))
	suite.Require().NotContains(ack, "error")

	// Pull out the msgs before substitution
	flowKeeper := GetICAApp(suite.IntentoChain).IntentKeeper
	flow := flowKeeper.GetFlowInfo(suite.IntentoChain.GetContext(), 1)
	unpacked := flow.GetTxMsgs(suite.IntentoChain.Codec)
	suite.Require().NotEmpty(unpacked)

	// === BEFORE substitution ===
	exec, ok := unpacked[0].(*wasmtypes.MsgExecuteContract)
	suite.Require().True(ok)

	// 1) JSON‑unmarshal msgExec.Msg into a string
	var rawB64 string
	suite.Require().NoError(json.Unmarshal(exec.Msg, &rawB64))

	// 2) Base64‑decode into raw JSON bytes
	decoded, err := base64.StdEncoding.DecodeString(rawB64)
	suite.Require().NoError(err)

	// 3) JSON‑unmarshal into a map
	var before map[string]interface{}
	suite.Require().NoError(json.Unmarshal(decoded, &before))

	beforeJSON, _ := json.MarshalIndent(before, "", "  ")
	suite.Require().True(
		strings.Contains(string(beforeJSON), types.ParseICAValue),
		"expected placeholder before substitution; got:\n%s", beforeJSON,
	)

	// === DO the substitution ===
	suite.IntentoChain.CurrentHeader.Time = suite.IntentoChain.CurrentHeader.Time.Add(time.Minute)
	flow.ICAConfig.ConnectionID = "connection-0"
	flowKeeper.HandleFlow(
		suite.IntentoChain.GetContext(),
		flowKeeper.Logger(suite.IntentoChain.GetContext()),
		flow,
		suite.IntentoChain.GetContext().BlockTime(),
		nil,
	)

	// === AFTER substitution ===
	flow = flowKeeper.GetFlowInfo(suite.IntentoChain.GetContext(), 1)
	unpacked = flow.GetTxMsgs(suite.IntentoChain.Codec)
	suite.Require().NotEmpty(unpacked)

	exec, ok = unpacked[0].(*wasmtypes.MsgExecuteContract)
	suite.Require().True(ok)

	// repeat the three‐step decode
	rawB64 = ""
	suite.Require().NoError(json.Unmarshal(exec.Msg, &rawB64))
	decoded, err = base64.StdEncoding.DecodeString(rawB64)
	suite.Require().NoError(err)

	var after map[string]interface{}
	suite.Require().NoError(json.Unmarshal(decoded, &after))

	afterJSON, _ := json.MarshalIndent(after, "", "  ")
	fmt.Println("✅ Decoded contract msg:\n", string(afterJSON))
	suite.Require().False(
		strings.Contains(string(afterJSON), types.ParseICAValue),
		"unexpected placeholder after substitution; got:\n%s", afterJSON,
	)
}

func derivePlaceholderSender(channel, originalSender string) sdk.AccAddress {
	senderStr := fmt.Sprintf("%s/%s", channel, originalSender)
	senderHash32 := address.Hash(types.SenderPrefix, []byte(senderStr))
	sender := sdk.AccAddress(senderHash32[:])
	return sender
}
