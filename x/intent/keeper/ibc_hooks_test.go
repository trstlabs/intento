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

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"

	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
)

func (suite *KeeperTestSuite) TestOnRecvTransferPacketWorks() {
	var (
		trace    transfertypes.DenomTrace
		amount   sdkmath.Int
		receiver string
	)

	suite.SetupTest() // reset

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

func (suite *KeeperTestSuite) TestOnRecvTransferPacketWithActionWorks() {
	suite.SetupTest() // reset

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
	// ackStr := string(ackBytes)
	// fmt.Println(ackStr)
	var ack map[string]string // This can't be unmarshalled to Acknowledgement because it's fetched from the events
	err := json.Unmarshal(ackBytes, &ack)
	suite.Require().NoError(err)
	suite.Require().NotContains(ack, "error")

	autoTx := suite.chainA.GetIntoApp().IntentKeeper.GetActionInfo(suite.chainA.GetContext(), 1)

	suite.Require().Equal(autoTx.Owner, addr.String())
	suite.Require().Equal(autoTx.Label, "my_trigger")
	suite.Require().Equal(autoTx.ICAConfig.PortID, "")
	suite.Require().Equal(autoTx.Interval, time.Second*60)

	var txMsgAny codectypes.Any
	cdc := codec.NewProtoCodec(suite.chainA.GetIntoApp().InterfaceRegistry())

	err = cdc.UnmarshalJSON([]byte(msg), &txMsgAny)
	suite.Require().NoError(err)
	suite.True(autoTx.Msgs[0].Equal(txMsgAny))
}

func (suite *KeeperTestSuite) TestOnRecvTransferPacketAndMultippleActionsWorks() {
	suite.SetupTest() // reset

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
	ackBytes := suite.receiveTransferPacket(addr.String(), fmt.Sprintf(`{"action": {"owner": "%s","label": "my_trigger", "cid":"%s", "host_cid":"%s","msgs": [%s, %s], "duration": "500s", "interval": "60s", "start_at": "0"} }`, addr.String(), path.EndpointA.ConnectionID, path.EndpointB.ConnectionID, msg, msg))

	var ack map[string]string // This can't be unmarshalled to Acknowledgement because it's fetched from the events
	err = json.Unmarshal(ackBytes, &ack)
	suite.Require().NoError(err)
	suite.Require().NotContains(ack, "error")

	autoTx := suite.chainA.GetIntoApp().IntentKeeper.GetActionInfo(suite.chainA.GetContext(), 1)

	suite.Require().Equal(autoTx.Owner, addr.String())
	suite.Require().Equal(autoTx.Label, "my_trigger")
	suite.Require().Equal(autoTx.ICAConfig.PortID, "icacontroller-"+addr.String())
	suite.Require().Equal(autoTx.ICAConfig.ConnectionID, path.EndpointA.ConnectionID)

	suite.Require().Equal(autoTx.Interval, time.Second*60)

	_, found := suite.chainA.GetIntoApp().ICAControllerKeeper.GetInterchainAccountAddress(suite.chainA.GetContext(), autoTx.ICAConfig.ConnectionID, autoTx.ICAConfig.PortID)
	suite.Require().True(found)

	var txMsgAny codectypes.Any
	cdc := codec.NewProtoCodec(suite.chainA.GetIntoApp().InterfaceRegistry())

	err = cdc.UnmarshalJSON([]byte(msg), &txMsgAny)
	suite.Require().NoError(err)
	suite.True(autoTx.Msgs[0].Equal(txMsgAny))
}

func (suite *KeeperTestSuite) TestOnRecvTransferPacketSubmitTxAndAddressParsingWorks() {
	suite.SetupTest() // reset

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

	ackBytes := suite.receiveTransferPacket(addr.String(), fmt.Sprintf(`{"action": {"owner": "%s","label": "my trigger", "cid":"%s","host_cid":"%s","msgs": [%s, %s], "duration": "120s", "interval": "60s", "start_at": "0"} }`, addr.String(), path.EndpointA.ConnectionID, path.EndpointB.ConnectionID, msg, msg))
	var ack map[string]string // This can't be unmarshalled to Acknowledgement because it's fetched from the events
	err = json.Unmarshal(ackBytes, &ack)
	suite.Require().NoError(err)
	suite.Require().NotContains(ack, "error")

	autoTxKeeper := suite.chainA.GetIntoApp().IntentKeeper
	autoTx := autoTxKeeper.GetActionInfo(suite.chainA.GetContext(), 1)
	feeAddr, _ := sdk.AccAddressFromBech32(autoTx.FeeAddress)
	types.Denom = suite.chainA.GetIntoApp().BankKeeper.GetAllBalances(suite.chainA.GetContext(), feeAddr)[0].Denom

	unpacker := suite.chainA.Codec
	unpackedMsgs := autoTx.GetTxMsgs(unpacker)
	suite.Require().True(strings.Contains(unpackedMsgs[0].String(), types.ParseICAValue))

	suite.chainA.CurrentHeader.Time = suite.chainA.CurrentHeader.Time.Add(time.Minute)
	FakeBeginBlocker(suite.chainA.GetContext(), autoTxKeeper, sdk.ConsAddress(suite.chainA.Vals.Proposer.Address))

	autoTx = autoTxKeeper.GetActionInfo(suite.chainA.GetContext(), 1)
	autoTxHistory, _ := autoTxKeeper.TryGetActionHistory(suite.chainA.GetContext(), autoTx.ID)
	suite.Require().NotNil(autoTxHistory.History)
	suite.Require().Empty(autoTxHistory.History[0].Errors)
	suite.Require().Equal(autoTx.Owner, addr.String())
	suite.Require().Equal(autoTx.Label, "my trigger")
	suite.Require().Equal(autoTx.ICAConfig.PortID, "icacontroller-"+addr.String())
	suite.Require().Equal(autoTx.ICAConfig.ConnectionID, path.EndpointA.ConnectionID)

	unpackedMsgs = autoTx.GetTxMsgs(unpacker)
	suite.Require().False(strings.Contains(unpackedMsgs[0].String(), types.ParseICAValue))
	suite.Require().Equal(autoTx.Interval, time.Second*60)
}
