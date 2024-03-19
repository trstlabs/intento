package keeper_test

import (
	"testing"
	"time"

	//"github.com/cosmos/cosmos-sdk/simapp"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"

	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"
	"github.com/stretchr/testify/suite"
	icaapp "github.com/trstlabs/intento/app"
	keeper "github.com/trstlabs/intento/x/intent/keeper"
	intentoibctesting "github.com/trstlabs/intento/x/intent/keeper/tests"
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
	chainA *intentoibctesting.TestChain
	chainB *intentoibctesting.TestChain
}

func GetICAApp(chain *ibctesting.TestChain) *icaapp.IntoApp {
	app, ok := chain.App.(*icaapp.IntoApp)
	if !ok {
		panic("not ica app")
	}

	return app
}

func GetActionKeeper(chain *intentoibctesting.TestChain) keeper.Keeper {
	app, ok := chain.App.(*icaapp.IntoApp)
	if !ok {
		panic("not ica app")
	}

	return app.IntentKeeper
}

func GetActionKeeperFromApp(app *icaapp.IntoApp) keeper.Keeper {

	return app.IntentKeeper
}

// TestKeeperTestSuite runs all the tests within this package.
func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

// SetupTest creates a coordinator with 2 test chains.
func (suite *KeeperTestSuite) SetupTest() {

	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)
	ibctesting.DefaultTestingAppInit = intentoibctesting.SetupTestingApp
	suite.chainA = &intentoibctesting.TestChain{TestChain: suite.coordinator.GetChain(ibctesting.GetChainID(1))}
	suite.chainB = &intentoibctesting.TestChain{TestChain: suite.coordinator.GetChain(ibctesting.GetChainID(2))}

}

func NewICAPath(chainA, chainB *intentoibctesting.TestChain) *ibctesting.Path {
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
func NewTransferPath(chainA, chainB *intentoibctesting.TestChain) *ibctesting.Path {
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

func (suite *KeeperTestSuite) receiveTransferPacket(receiver, memo string) []byte {
	return suite.receiveTransferPacketWithSequence(receiver, memo, 0)
}

func (suite *KeeperTestSuite) receiveTransferPacketWithSequence(receiver, memo string, prevSequence uint64) []byte {
	// fmt.Println(memo)
	path := NewTransferPath(suite.chainA, suite.chainB)

	suite.coordinator.Setup(path)
	channelCap := suite.chainB.GetChannelCapability(
		path.EndpointB.ChannelConfig.PortID,
		path.EndpointB.ChannelID)
	packet := suite.makeMockPacket(receiver, memo, prevSequence, path)

	_, err := suite.chainB.GetIntoApp().IBCKeeper.ChannelKeeper.SendPacket(
		suite.chainB.GetContext(), channelCap, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, clienttypes.ZeroHeight(), uint64(suite.chainB.GetContext().BlockTime().Add(time.Minute).UnixNano()), packet.Data)
	suite.Require().NoError(err, "IBC send failed. Expected success. %s", err)

	// Update both clients
	// recv in chain a
	// get the ack from the chain a's response
	// manually send the acknowledgement to chain b
	ack := receivePacket(path, suite, packet)
	return ack
}

func receivePacket(path *ibctesting.Path, suite *KeeperTestSuite, packet channeltypes.Packet) []byte {
	err := path.EndpointB.UpdateClient()
	suite.Require().NoError(err)
	err = path.EndpointA.UpdateClient()
	suite.Require().NoError(err)

	res, err := path.EndpointA.RecvPacketWithResult(packet)
	suite.Require().NoError(err)

	ack, err := ibctesting.ParseAckFromEvents(res.GetEvents())
	suite.Require().NoError(err)

	err = path.EndpointA.AcknowledgePacket(packet, ack)
	suite.Require().NoError(err)
	return ack
}

func (suite *KeeperTestSuite) makeMockPacket(receiver, memo string, prevSequence uint64, path *ibctesting.Path) channeltypes.Packet {

	packetData := transfertypes.FungibleTokenPacketData{
		Denom:    sdk.DefaultBondDenom,
		Amount:   "1000000",
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
