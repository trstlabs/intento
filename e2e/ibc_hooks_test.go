package interchaintest

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	math "cosmossdk.io/math"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/icza/dyno"
	interchaintest "github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/strangelove-ventures/interchaintest/v8/testreporter"
	"github.com/strangelove-ventures/interchaintest/v8/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestIBCHooksE2E(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	t.Parallel()

	ctx := context.Background()

	// Modify the the timeout_commit in the config.toml node files
	// to reduce the block commit times. This speeds up the tests
	// by about 35%
	configFileOverrides := make(map[string]any)
	configTomlOverrides := make(testutil.Toml)
	consensus := make(testutil.Toml)
	consensus["timeout_commit"] = "1s"
	configTomlOverrides["consensus"] = consensus
	configFileOverrides["config/config.toml"] = configTomlOverrides
	genesisKVMods := []cosmos.GenesisKV{
		cosmos.NewGenesisKV("app_state.feemarket.params.enabled", false),
		cosmos.NewGenesisKV("app_state.feemarket.params.min_base_gas_price", "0.001000000000000000"),
		cosmos.NewGenesisKV("app_state.feemarket.state.base_gas_price", "0.001000000000000000"),
	}
	// Chain Factory
	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{Name: "gaia", Version: "v22.1.0", ChainConfig: ibc.ChainConfig{
			GasAdjustment:       1.5,
			GasPrices:           "0.01uatom",
			ModifyGenesis:       cosmos.ModifyGenesis(genesisKVMods), //SetupGaiaGenesis([]string{"*"}),
			ConfigFileOverrides: configFileOverrides,
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
			GasAdjustment:  1.5,
			TrustingPeriod: "508h",
			NoHostMount:    false},
		},
	})

	// logger.go:146: 2025-02-22T14:32:29.624+0100 INFO    Exec    {"validator": true, "i": 0, "chain_id": "gaia-1", "test": "TestIBCHooksE2E", "image": "ghcr.io/strangelove-ventures/heighliner/gaia:v22.1.0", "test_name": "TestIBCHooksE2E", "command": "gaiad tx gov submit-legacy-proposal consumer-addition /var/cosmos-chain/gaia-1/proposal_zpen.json --gas auto --gas-prices 0.01uatom --gas-adjustment 1.5 --from proposer --keyring-backend test --output json -y --chain-id gaia-1 --home /var/cosmos-chain/gaia-1 --node tcp://gaia-1-val-0-TestIBCHooksE2E:26657", "hostname": "TestIBCHooksE2E-hsclyn", "container": "TestIBCHooksE2E-hsclyn"}
	// ibc_hooks_test.go:164:
	//             Error Trace:    /root/intento/e2e/ibc_hooks_test.go:164
	//             Error:          Received unexpected error:
	//                             failed to start chains: failed to start provider chain gaia-1: failed to submit consumer addition proposal: exit code 1:  Error: failed to parse proposal: proposal type is required
	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)
	gaia, intento := chains[0], chains[1]

	// Relayer Factory
	client, network := interchaintest.DockerSetup(t)
	r := interchaintest.NewBuiltinRelayerFactory(ibc.CosmosRly, zaptest.NewLogger(t)).Build(
		t, client, network)

	// Prep Interchain
	const gaiaIntentoICSPath = "gaia-intento-ics-path"
	const ibcPath = "gaia-intento-demo"
	ic := interchaintest.NewInterchain().
		AddChain(gaia).
		AddChain(intento).
		AddRelayer(r, "relayer").
		AddProviderConsumerLink(interchaintest.ProviderConsumerLink{Provider: gaia, Consumer: intento, Relayer: r, Path: gaiaIntentoICSPath}).
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
	users2 := interchaintest.GetAndFundTestUsers(t, ctx, "default", fundAmount, gaia, intento)
	gaiaUser := users[0]
	intentoUser := users[1]
	intentoUser2 := users2[1]

	gaiaUserBalInitial, err := gaia.GetBalance(ctx, gaiaUser.FormattedAddress(), gaia.Config().Denom)
	require.NoError(t, err)
	require.True(t, gaiaUserBalInitial.Equal(fundAmount))

	// Get Channel ID
	gaiaChannelInfo, err := r.GetChannels(ctx, eRep, gaia.Config().ChainID)
	require.NoError(t, err)
	gaiaChannelID := gaiaChannelInfo[0].ChannelID

	intentoChannelInfo, err := r.GetChannels(ctx, eRep, intento.Config().ChainID)
	require.NoError(t, err)
	intentoChannelID := intentoChannelInfo[0].ChannelID

	height, err := intento.Height(ctx)
	require.NoError(t, err)

	// Send Transaction
	amountToSend := math.NewInt(1_000_000)
	dstAddress := intentoUser.FormattedAddress()
	msgSend := fmt.Sprintf(`{
		"@type":"/cosmos.bank.v1beta1.MsgSend",
		"amount": [{
			"amount": "70",
			"denom": "uatom"
		}],
		"from_address": "%s",
		"to_address": "%s"
	}`, intentoUser, intentoUser2)

	memo := fmt.Sprintf(`{"flow": {"owner": "%s","label": "my_trigger", "msgs": [%s], "duration": "500s", "interval": "60s", "start_at": "0"} }`, intentoUser, msgSend)
	// memoBytes, err := json.Marshal(memo)
	// require.NoError(t, err)

	transfer := ibc.WalletAmount{
		Address: dstAddress,
		Denom:   gaia.Config().Denom,
		Amount:  amountToSend,
	}
	tx, err := gaia.SendIBCTransfer(ctx, gaiaChannelID, gaiaUser.KeyName(), transfer, ibc.TransferOptions{Memo: memo})
	require.NoError(t, err)
	require.NoError(t, tx.Validate())

	// ==== IMPLEMENT IBC TRANSFER WITH FLOW MEMO ====

	// ==== VERIFY FLOW EXECUTION ON Intento ====
	time.Sleep(5 * time.Second) // Wait for IBC packet delivery

	// Query Intent Keeper on Intento for Flow Info
	query := []string{"q", "intent", "flow", " 1"}
	stdOut, stErr, err := chains[1].Exec(ctx, query, nil)
	require.NoError(t, err)
	require.Nil(t, stErr)
	require.NotNil(t, stdOut)
	require.Contains(t, string(stdOut), "1")

	t.Logf("Flow info on Intento: %s", string(stdOut))

	srcDenomTrace := ibctransfertypes.ParseDenomTrace(ibctransfertypes.GetPrefixedDenom("transfer", intentoChannelID, gaia.Config().Denom))
	dstIbcDenom := srcDenomTrace.IBCDenom()

	// Test destination wallet has increased funds
	intentoUserBalNew, err := intento.GetBalance(ctx, intentoUser.FormattedAddress(), dstIbcDenom)
	require.NoError(t, err)
	require.True(t, intentoUserBalNew.Equal(amountToSend))

	intentoUser2BalNew, err := intento.GetBalance(ctx, intentoUser.FormattedAddress(), dstIbcDenom)
	require.NoError(t, err)
	require.True(t, intentoUser2BalNew.Sub(fundAmount).Equal(math.NewInt(70)))

	chain := intento.(*cosmos.CosmosChain)
	reg := chain.Config().EncodingConfig.InterfaceRegistry
	msgUpdateClient, err := cosmos.PollForMessage[*clienttypes.MsgUpdateClient](ctx, chain, reg, height, height+10, nil)
	require.NoError(t, err)

	require.Equal(t, "07-tendermint-0", msgUpdateClient.ClientId)
	require.NotEmpty(t, msgUpdateClient.Signer)
	t.Log("IBC-Hooks E2E test passed!")
}

// Sets custom fields for the Gaia genesis file that interchaintest isn't aware of by default.
//
// allowed_messages - explicitly allowed messages to be accepted by the the interchainaccounts section
func SetupGaiaGenesis(allowed_messages []string) func(ibc.ChainConfig, []byte) ([]byte, error) {
	return func(chainConfig ibc.ChainConfig, genbz []byte) ([]byte, error) {
		//g := make(map[string]interface{})
		g := []cosmos.GenesisKV{
			cosmos.NewGenesisKV("app_state.feemarket.params.enabled", false),
			cosmos.NewGenesisKV("app_state.feemarket.params.min_base_gas_price", "0.001000000000000000"),
			cosmos.NewGenesisKV("app_state.feemarket.state.base_gas_price", "0.001000000000000000"),
		}
		if err := json.Unmarshal(genbz, &g); err != nil {
			return nil, fmt.Errorf("failed to unmarshal genesis file: %w", err)
		}

		if err := dyno.Set(g, allowed_messages, "app_state", "interchainaccounts", "host_genesis_state", "params", "allow_messages"); err != nil {
			return nil, fmt.Errorf("failed to set allow_messages for interchainaccount host in genesis json: %w", err)
		}

		out, err := json.Marshal(g)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal genesis bytes to json: %w", err)
		}
		return out, nil
	}
}
