package keeper_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	//"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"

	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"
	"github.com/stretchr/testify/suite"
	icaapp "github.com/trstlabs/trst/app"
	autoIbcTxKeeper "github.com/trstlabs/trst/x/auto-ibc-tx/keeper"
	trstibctesting "github.com/trstlabs/trst/x/auto-ibc-tx/keeper/tests"
)

var (
	// TestAccAddress defines a resuable bech32 address for testing purposes
	// TODO: update crypto.AddressHash() when sdk uses address.Module()
	//TestAccAddress = icatypes.GenerateAddress(sdk.AccAddress(crypto.AddressHash([]byte(icatypes.ModuleName))), ibctesting.FirstConnectionID, TestPortID)
	// TestOwnerAddress defines a reusable bech32 address for testing purposes
	TestOwnerAddress = "cosmos17dtl0mjt3t77kpuhg2edqzjpszulwhgzuj9ljs"
	// TestPortID defines a resuable port identifier for testing purposes
	TestPortID, _ = icatypes.NewControllerPortID(TestOwnerAddress)
	// TestVersion defines a resuable interchainaccounts version string for testing purposes
	TestVersion = string(icatypes.ModuleCdc.MustMarshalJSON(&icatypes.Metadata{
		Version:                icatypes.Version,
		ControllerConnectionId: ibctesting.FirstConnectionID,
		HostConnectionId:       ibctesting.FirstConnectionID,
		Encoding:               icatypes.EncodingProtobuf,
		TxType:                 icatypes.TxTypeSDKMultiMsg,
	}))
)

// KeeperTestSuite is a testing suite to test keeper functions
type KeeperTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	// testing chains used for convenience and readability
	chainA *trstibctesting.TestChain
	chainB *trstibctesting.TestChain
}

func GetICAApp(chain *ibctesting.TestChain) *icaapp.TrstApp {
	app, ok := chain.App.(*icaapp.TrstApp)
	if !ok {
		panic("not ica app")
	}

	return app
}

func GetICAKeeper(chain *trstibctesting.TestChain) autoIbcTxKeeper.Keeper {
	app, ok := chain.App.(*icaapp.TrstApp)
	if !ok {
		panic("not ica app")
	}

	return app.AutoIbcTxKeeper
}

func GetICAKeeper2(app *icaapp.TrstApp) autoIbcTxKeeper.Keeper {

	return app.AutoIbcTxKeeper
}

// TestKeeperTestSuite runs all the tests within this package.
func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

// SetupTest creates a coordinator with 2 test chains.
func (suite *KeeperTestSuite) SetupTest() {

	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)
	ibctesting.DefaultTestingAppInit = trstibctesting.SetupTestingApp
	suite.chainA = &trstibctesting.TestChain{TestChain: suite.coordinator.GetChain(ibctesting.GetChainID(1))}
	suite.chainB = &trstibctesting.TestChain{TestChain: suite.coordinator.GetChain(ibctesting.GetChainID(2))}

}

func NewICAPath(chainA, chainB *trstibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA.TestChain, chainB.TestChain)
	path.EndpointA.ChannelConfig.PortID = icatypes.HostPortID
	path.EndpointB.ChannelConfig.PortID = icatypes.HostPortID
	path.EndpointA.ChannelConfig.Order = channeltypes.ORDERED
	path.EndpointB.ChannelConfig.Order = channeltypes.ORDERED
	path.EndpointA.ChannelConfig.Version = TestVersion
	path.EndpointB.ChannelConfig.Version = TestVersion

	return path
}

// ToDo: Move this to osmosistesting to avoid repetition
func NewTransferPath(chainA, chainB *trstibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA.TestChain, chainB.TestChain)
	path.EndpointA.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointB.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointA.ChannelConfig.Version = transfertypes.Version
	path.EndpointB.ChannelConfig.Version = transfertypes.Version

	return path
}

// SetupICAPath invokes the InterchainAccounts entrypoint and subsequent channel handshake handlers
func SetupICAPath(path *ibctesting.Path, owner string) error {
	if err := RegisterInterchainAccount(path.EndpointA, owner); err != nil {
		return err
	}
	if err := path.EndpointB.ChanOpenTry(); err != nil {
		return err
	}

	if err := path.EndpointA.ChanOpenAck(); err != nil {
		return err
	}

	if err := path.EndpointB.ChanOpenConfirm(); err != nil {
		return err
	}

	return nil
}

// RegisterInterchainAccount is a helper function for starting the channel handshake
func RegisterInterchainAccount(endpoint *ibctesting.Endpoint, owner string) error {
	portID, err := icatypes.NewControllerPortID(owner)
	if err != nil {
		return err
	}

	channelSequence := endpoint.Chain.App.GetIBCKeeper().ChannelKeeper.GetNextChannelSequence(endpoint.Chain.GetContext())

	if err := GetICAApp(endpoint.Chain).ICAControllerKeeper.RegisterInterchainAccount(endpoint.Chain.GetContext(), endpoint.ConnectionID, owner, TestVersion); err != nil {
		return err
	}

	// commit state changes for proof verification
	endpoint.Chain.NextBlock()

	// update port/channel ids
	endpoint.ChannelID = channeltypes.FormatChannelIdentifier(channelSequence)
	endpoint.ChannelConfig.PortID = portID

	return nil
}

func (suite *KeeperTestSuite) TestOnRecvTransferPacketWorks() {
	var (
		trace    transfertypes.DenomTrace
		amount   sdk.Int
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

	ack := suite.chainB.GetTrstApp().TransferStack.OnRecvPacket(suite.chainB.GetContext(), packet, suite.chainA.SenderAccount.GetAddress())

	suite.Require().True(ack.Success())

}

func (suite *KeeperTestSuite) TestOnRecvTransferPacketWithAutoTxWorks() {
	suite.SetupTest() // reset

	addr := suite.chainA.SenderAccount.GetAddress()
	msg := `{
		"@type":"/cosmos.bank.v1beta1.MsgSend",
		"amount": [{
			"amount": "70",
			"denom": "stake"
		}],
		"from_address": "trust12gxmzpucje8aflw2vz45rv8x4nyaaj3rp8vjh03dulehkdl5fu6s93ewkp",
		"to_address": "trust1ykql5ktedxkpjszj5trzu8f5dxajvgv95nuwjx"
	}`

	ackBytes := suite.receivePacket(addr.String(), fmt.Sprintf(`{"auto_tx": {"owner": "%s","label": "my_trigger", "msgs": [%s], "duration": "500s", "interval": "60s", "start_at": "0"} }`, addr, msg))
	// ackStr := string(ackBytes)
	// fmt.Println(ackStr)
	var ack map[string]string // This can't be unmarshalled to Acknowledgement because it's fetched from the events
	err := json.Unmarshal(ackBytes, &ack)
	suite.Require().NoError(err)
	suite.Require().NotContains(ack, "error")

	autoTx := suite.chainA.GetTrstApp().AutoIbcTxKeeper.GetAutoTxInfo(suite.chainA.GetContext(), 1)

	suite.Require().Equal(autoTx.Owner, addr.String())
	suite.Require().Equal(autoTx.Label, "my_trigger")
	suite.Require().Equal(autoTx.PortID, "")
	suite.Require().Equal(autoTx.Interval, time.Second*60)

	var txMsgAny codectypes.Any
	cdc := codec.NewProtoCodec(suite.chainA.GetTrstApp().InterfaceRegistry())

	err = cdc.UnmarshalJSON([]byte(msg), &txMsgAny)
	suite.Require().NoError(err)
	suite.True(autoTx.Msgs[0].Equal(txMsgAny))
}

func (suite *KeeperTestSuite) TestOnRecvTransferPacketAndMultippleAutoTxsWorks() {
	suite.SetupTest() // reset

	addr := suite.chainA.SenderAccount.GetAddress()
	msg := `{
		"@type":"/cosmos.bank.v1beta1.MsgSend",
		"amount": [{
			"amount": "70",
			"denom": "stake"
		}],
		"from_address": "trust12gxmzpucje8aflw2vz45rv8x4nyaaj3rp8vjh03dulehkdl5fu6s93ewkp",
		"to_address": "trust1ykql5ktedxkpjszj5trzu8f5dxajvgv95nuwjx"
	}`

	path := NewICAPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupConnections(path)
	err := SetupICAPath(path, addr.String())
	suite.Require().NoError(err)

	//chainB sends packet to chainA. connectionID to execute on chainB is on chainAs config
	ackBytes := suite.receivePacket(addr.String(), fmt.Sprintf(`{"auto_tx": {"owner": "%s","label": "my_trigger", "connection_id":"%s", "msgs": [%s, %s], "duration": "500s", "interval": "60s", "start_at": "0"} }`, addr.String(), path.EndpointA.ConnectionID, msg, msg))
	// ackStr := string(ackBytes)
	// fmt.Println(ackStr)
	var ack map[string]string // This can't be unmarshalled to Acknowledgement because it's fetched from the events
	err = json.Unmarshal(ackBytes, &ack)
	suite.Require().NoError(err)
	suite.Require().NotContains(ack, "error")

	autoTx := suite.chainA.GetTrstApp().AutoIbcTxKeeper.GetAutoTxInfo(suite.chainA.GetContext(), 1)

	suite.Require().Equal(autoTx.Owner, addr.String())
	suite.Require().Equal(autoTx.Label, "my_trigger")
	suite.Require().Equal(autoTx.PortID, "icacontroller-"+addr.String())
	suite.Require().Equal(autoTx.ConnectionID, path.EndpointA.ConnectionID)

	suite.Require().Equal(autoTx.Interval, time.Second*60)

	_, found := suite.chainA.GetTrstApp().ICAControllerKeeper.GetInterchainAccountAddress(suite.chainA.GetContext(), autoTx.ConnectionID, autoTx.PortID)
	suite.Require().True(found)

	var txMsgAny codectypes.Any
	cdc := codec.NewProtoCodec(suite.chainA.GetTrstApp().InterfaceRegistry())

	err = cdc.UnmarshalJSON([]byte(msg), &txMsgAny)
	suite.Require().NoError(err)
	suite.True(autoTx.Msgs[0].Equal(txMsgAny))
}

func (suite *KeeperTestSuite) TestOnRecvTransferPacketWithRegistrationAndMultippleAutoTxsWorks() {
	suite.SetupTest() // reset

	addr := suite.chainA.SenderAccount.GetAddress()
	msg := `{
		"@type":"/cosmos.bank.v1beta1.MsgSend",
		"amount": [{
			"amount": "70",
			"denom": "stake"
		}],
		"from_address": "trust12gxmzpucje8aflw2vz45rv8x4nyaaj3rp8vjh03dulehkdl5fu6s93ewkp",
		"to_address": "trust1ykql5ktedxkpjszj5trzu8f5dxajvgv95nuwjx"
	}`

	path := NewICAPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupConnections(path)

	ackBytes := suite.receivePacket(addr.String(), fmt.Sprintf(`{"auto_tx": {"owner": "%s","label": "my_trigger", "connection_id":"%s", "msgs": [%s, %s], "duration": "500s", "interval": "60s", "start_at": "0", "register_ica": "true"} }`, addr.String(), path.EndpointA.ConnectionID, msg, msg))
	// ackStr := string(ackBytes)
	// fmt.Println(ackStr)

	var ack map[string]string // This can't be unmarshalled to Acknowledgement because it's fetched from the events
	err := json.Unmarshal(ackBytes, &ack)
	suite.Require().NoError(err)
	suite.Require().NotContains(ack, "error")

	autoTx := suite.chainA.GetTrstApp().AutoIbcTxKeeper.GetAutoTxInfo(suite.chainA.GetContext(), 1)

	suite.Require().Equal(autoTx.Owner, addr.String())
	suite.Require().Equal(autoTx.Label, "my_trigger")
	suite.Require().Equal(autoTx.PortID, "icacontroller-"+addr.String())
	suite.Require().Equal(autoTx.ConnectionID, path.EndpointA.ConnectionID)

	suite.Require().Equal(autoTx.Interval, time.Second*60)
	/*
		// Update both clients
		err = path.EndpointB.UpdateClient()
		suite.Require().NoError(err)
		err = path.EndpointA.UpdateClient()
		suite.Require().NoError(err)

		suite.chainA.NextBlock()
		suite.chainB.NextBlock() */

	// _, found := suite.chainA.GetTrstApp().ICAControllerKeeper.GetInterchainAccountAddress(suite.chainA.GetContext(), autoTx.ConnectionID, autoTx.PortID)
	// suite.Require().True(found)

	var txMsgAny codectypes.Any
	cdc := codec.NewProtoCodec(suite.chainA.GetTrstApp().InterfaceRegistry())

	err = cdc.UnmarshalJSON([]byte(msg), &txMsgAny)
	suite.Require().NoError(err)
	suite.True(autoTx.Msgs[0].Equal(txMsgAny))
}

func (suite *KeeperTestSuite) receivePacket(receiver, memo string) []byte {
	return suite.receivePacketWithSequence(receiver, memo, 0)
}

func (suite *KeeperTestSuite) receivePacketWithSequence(receiver, memo string, prevSequence uint64) []byte {
	fmt.Println(memo)
	path := NewTransferPath(suite.chainA, suite.chainB)

	suite.coordinator.Setup(path)
	channelCap := suite.chainB.GetChannelCapability(
		path.EndpointB.ChannelConfig.PortID,
		path.EndpointB.ChannelID)
	packet := suite.makeMockPacket(receiver, memo, prevSequence, path)

	_, err := suite.chainB.GetTrstApp().IBCKeeper.ChannelKeeper.SendPacket(
		suite.chainB.GetContext(), channelCap, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, clienttypes.ZeroHeight(), uint64(suite.chainB.GetContext().BlockTime().Add(time.Minute).UnixNano()), packet.Data)
	suite.Require().NoError(err, "IBC send failed. Expected success. %s", err)

	// Update both clients
	err = path.EndpointB.UpdateClient()
	suite.Require().NoError(err)
	err = path.EndpointA.UpdateClient()
	suite.Require().NoError(err)

	// recv in chain a
	res, err := path.EndpointA.RecvPacketWithResult(packet)
	suite.Require().NoError(err)
	// get the ack from the chain a's response
	ack, err := ibctesting.ParseAckFromEvents(res.GetEvents())
	suite.Require().NoError(err)

	// manually send the acknowledgement to chain b
	err = path.EndpointA.AcknowledgePacket(packet, ack)
	suite.Require().NoError(err)
	return ack
}

func (suite *KeeperTestSuite) makeMockPacket(receiver, memo string, prevSequence uint64, path *ibctesting.Path) channeltypes.Packet {
	packetData := transfertypes.FungibleTokenPacketData{
		Denom:    sdk.DefaultBondDenom,
		Amount:   "1",
		Sender:   suite.chainB.SenderAccount.GetAddress().String(),
		Receiver: receiver,
		Memo:     memo,
	}

	return channeltypes.NewPacket(
		packetData.GetBytes(),
		prevSequence+1,
		path.EndpointB.ChannelConfig.PortID,
		path.EndpointB.ChannelID,
		path.EndpointA.ChannelConfig.PortID,
		path.EndpointA.ChannelID,
		clienttypes.ZeroHeight(),
		uint64(suite.chainB.GetContext().BlockTime().Add(time.Minute).UnixNano()),
	)
}
