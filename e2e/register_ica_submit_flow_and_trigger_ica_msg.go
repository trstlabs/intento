package interchaintest

import (
	"context"
	"fmt"
	"testing"
	"time"

	math "cosmossdk.io/math"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	interchaintest "github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/strangelove-ventures/interchaintest/v8/testreporter"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestRegisterICASubmitFlowTriggerICAMsg(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	t.Parallel()

	ctx := context.Background()

	// Chain Factory
	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{Name: "gaia", Version: "v22.1.0", ChainConfig: ibc.ChainConfig{
			GasPrices: "0.0uatom",
		}},
		{ChainConfig: ibc.ChainConfig{
			Type:    "cosmos",
			Name:    "intento",
			ChainID: "intento-1",
			Images: []ibc.DockerImage{
				{
					Repository: "intento", // FOR LOCAL IMAGE USE: Docker Image Name
					Version:    "local",   // FOR LOCAL IMAGE USE: Docker Image Tag
					UIDGID:     "1025:1025",
				},
			},
			Bin:            "intentod",
			Bech32Prefix:   "into",
			Denom:          "uinto",
			GasPrices:      "0.0uinto",
			GasAdjustment:  1.3,
			TrustingPeriod: "508h",
			NoHostMount:    false},
		},
	})

	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)
	gaia, intento := chains[0], chains[1]

	// Relayer Factory
	client, network := interchaintest.DockerSetup(t)
	r := interchaintest.NewBuiltinRelayerFactory(ibc.CosmosRly, zaptest.NewLogger(t)).Build(
		t, client, network)

	// Prep Interchain
	const ibcPath = "gaia-intento-demo"
	ic := interchaintest.NewInterchain().
		AddChain(gaia).
		AddChain(intento).
		AddRelayer(r, "relayer").
		AddLink(interchaintest.InterchainLink{
			Chain1:  gaia,
			Chain2:  intento,
			Relayer: r,
			Path:    ibcPath,
		})

	// Log location
	f, err := interchaintest.CreateLogFile(fmt.Sprintf("%d.json", time.Now().Unix()))
	require.NoError(t, err)
	// Reporter/logs
	rep := testreporter.NewReporter(f)
	eRep := rep.RelayerExecReporter(t)

	// Build interchain
	require.NoError(t, ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:         t.Name(),
		Client:           client,
		NetworkID:        network,
		SkipPathCreation: false,
	},
	),
	)

	// Create and Fund User Wallets
	fundAmount := math.NewInt(10_000_000)
	users := interchaintest.GetAndFundTestUsers(t, ctx, "default", fundAmount, gaia, intento)
	gaiaUser := users[0]
	//intentoUser := users[1]
	gaiaUser2 := users[2]

	gaiaUserBalInitial, err := gaia.GetBalance(ctx, gaiaUser.FormattedAddress(), gaia.Config().Denom)
	require.NoError(t, err)
	require.True(t, gaiaUserBalInitial.Equal(fundAmount))

	// // Get Channel ID
	// gaiaChannelInfo, err := r.GetChannels(ctx, eRep, gaia.Config().ChainID)
	// require.NoError(t, err)
	// gaiaChannelID := gaiaChannelInfo[0].ChannelID
	// intentoChannelInfo, err := r.GetChannels(ctx, eRep, intento.Config().ChainID)
	// require.NoError(t, err)
	// intentoChannelID := intentoChannelInfo[0].ChannelID

	err = r.CreateConnections(ctx, rep.RelayerExecReporter(t), ibcPath)
	require.NoError(t, err)
	height, err := intento.Height(ctx)
	require.NoError(t, err)

	msgSend := fmt.Sprintf(`{
		"@type":"/cosmos.bank.v1beta1.MsgSend",
		"amount": [{
			"amount": "70",
			"denom": "uatom"
		}],
		"from_address": "%s",
		"to_address": "%s"
	}`, gaiaUser, gaiaUser2)

	// Submit Flow
	cmd := []string{"tx", "intent", "register-ica-and-submit-flow", msgSend, "--duration", "60s", "--connection-id", "0" /* or 1 */, "--host-connection-id", "0" /* or 1 */}
	stdOut, stErr, err := chains[1].Exec(ctx, cmd, nil)
	require.NoError(t, err)
	require.Nil(t, stErr)
	require.NotNil(t, stdOut)

	// Send Transaction

	// ==== VERIFY FLOW EXECUTION ON GAIA ====
	time.Sleep(5 * time.Second)  // Wait for IBC packet delivery
	time.Sleep(20 * time.Second) // Wait for Registration ICA
	time.Sleep(60 * time.Second) // Wait for Action Trigger
	time.Sleep(5 * time.Second)  // Wait for IBC packet delivery

	// Query  Flow Info
	query := []string{"q", "intent", "flow", " 1"}
	stdOut, stErr, err = chains[1].Exec(ctx, query, nil)
	require.NoError(t, err)
	require.Nil(t, stErr)
	require.NotNil(t, stdOut)
	require.Contains(t, string(stdOut), "1")

	t.Logf("Flow info on Intento: %s", string(stdOut))

	// Test destination wallet has increased funds
	gaiaUser2BalNew, err := gaia.GetBalance(ctx, gaiaUser2.FormattedAddress(), "uatom")
	require.NoError(t, err)
	require.True(t, gaiaUser2BalNew.Sub(fundAmount).Equal(math.NewInt(70)))

	chain := intento.(*cosmos.CosmosChain)
	reg := chain.Config().EncodingConfig.InterfaceRegistry
	msgUpdateClient, err := cosmos.PollForMessage[*clienttypes.MsgUpdateClient](ctx, chain, reg, height, height+10, nil)
	require.NoError(t, err)

	require.Equal(t, "07-tendermint-0", msgUpdateClient.ClientId)
	require.NotEmpty(t, msgUpdateClient.Signer)
	t.Log("IBC-Hooks E2E test passed!")
}
