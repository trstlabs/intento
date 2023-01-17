package keeper_test

import (
	"encoding/json"
	"testing"

	//"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/trstlabs/trst/x/compute"

	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v3/testing"
	simapp "github.com/cosmos/ibc-go/v3/testing/simapp"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	icaapp "github.com/trstlabs/trst/app"
	autoibctxkeeper "github.com/trstlabs/trst/x/auto-ibc-tx/keeper"
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
	}))
)

func init() {
	ibctesting.DefaultTestingAppInit = SetupICATestingApp
}

func SetupICATestingApp() (ibctesting.TestingApp, map[string]json.RawMessage) {
	db := dbm.NewMemDB()
	// encCdc := icaapp.MakeEncodingConfig()
	app := icaapp.NewTrstApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, icaapp.DefaultNodeHome, 5, true, simapp.EmptyAppOptions{}, compute.DefaultWasmConfig(), icaapp.GetEnabledProposals())
	// TODO: figure out if it's ok that w MakeEncodingConfig inside of our Genesis.go. It would be a different instance than the one used in app
	return app, icaapp.NewDefaultGenesisState(app.AppCodec())
}

// KeeperTestSuite is a testing suite to test keeper functions
type KeeperTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	// testing chains used for convenience and readability
	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain
}

/*
	 func GetICAApp(chain *ibctesting.TestChain) *simapp.SimApp {
		app, ok := chain.App.(*simapp.SimApp)
		if !ok {
			panic("not ica app")
		}

		return app
	}
*/

func GetICAApp(chain *ibctesting.TestChain) *icaapp.TrstApp {
	app, ok := chain.App.(*icaapp.TrstApp)
	if !ok {
		panic("not ica app")
	}

	return app
}

func GetICAKeeper(chain *ibctesting.TestChain) autoibctxkeeper.Keeper {
	app, ok := chain.App.(*icaapp.TrstApp)
	if !ok {
		panic("not ica app")
	}

	return *app.AppKeepers.ICAAuthKeeper
}

func GetICAKeeper2(app *icaapp.TrstApp) autoibctxkeeper.Keeper {

	return *app.AppKeepers.ICAAuthKeeper
}

// TestKeeperTestSuite runs all the tests within this package.
func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

// SetupTest creates a coordinator with 2 test chains.
func (suite *KeeperTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(2))
}

func NewICAPath(chainA, chainB *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = icatypes.PortID
	path.EndpointB.ChannelConfig.PortID = icatypes.PortID
	path.EndpointA.ChannelConfig.Order = channeltypes.ORDERED
	path.EndpointB.ChannelConfig.Order = channeltypes.ORDERED
	path.EndpointA.ChannelConfig.Version = TestVersion
	path.EndpointB.ChannelConfig.Version = TestVersion

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

	if err := GetICAApp(endpoint.Chain).AppKeepers.ICAControllerKeeper.RegisterInterchainAccount(endpoint.Chain.GetContext(), endpoint.ConnectionID, owner); err != nil {
		return err
	}

	// commit state changes for proof verification
	endpoint.Chain.NextBlock()

	// update port/channel ids
	endpoint.ChannelID = channeltypes.FormatChannelIdentifier(channelSequence)
	endpoint.ChannelConfig.PortID = portID

	return nil
}
