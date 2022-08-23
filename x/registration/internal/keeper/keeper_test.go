package keeper

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	eng "github.com/trstlabs/trst/types"
	"github.com/trstlabs/trst/x/registration/internal/types"
	ra "github.com/trstlabs/trst/x/registration/remote_attestation"
)

func init() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(eng.Bech32PrefixAccAddr, eng.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(eng.Bech32PrefixValAddr, eng.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(eng.Bech32PrefixConsAddr, eng.Bech32PrefixConsPub)
	config.Seal()
}

func TestNewKeeper(t *testing.T) {
	tempDir := t.TempDir()
	_, regKeeper := CreateTestInput(t, false, tempDir, true)
	require.NotNil(t, regKeeper)
}

func TestNewKeeper_Node(t *testing.T) {
	tempDir := t.TempDir()

	seedPath := filepath.Join(tempDir, types.NodeCfgFolder, types.NodeSeedConfig)

	err := os.MkdirAll(filepath.Join(tempDir, types.NodeCfgFolder), 0o700)
	require.NoError(t, err)

	err = ioutil.WriteFile(seedPath, CreateTestSeedConfig(t), 0o700)
	require.NoError(t, err)

	_, regKeeper := CreateTestInput(t, false, tempDir, false)
	require.NotNil(t, regKeeper)
}

func TestKeeper_RegisterationStore(t *testing.T) {
	tempDir := t.TempDir()
	ctx, regKeeper := CreateTestInput(t, false, tempDir, true)

	cert, err := ioutil.ReadFile("../../testdata/attestation_cert_sw")
	require.NoError(t, err)

	regInfo := types.RegistrationNodeInfo{
		Certificate:   cert,
		EncryptedSeed: []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
	}

	regKeeper.SetRegistrationInfo(ctx, regInfo)

	publicKey, err := ra.VerifyRaCert(regInfo.Certificate)
	if err != nil {
		return
	}

	regInfo2 := regKeeper.getRegistrationInfo(ctx, publicKey)

	require.Equal(t, regInfo, *regInfo2)
}

func TestKeeper_RegisterNode(t *testing.T) {
	tempDir := t.TempDir()
	ctx, regKeeper := CreateTestInput(t, false, tempDir, true)

	cert, err := ioutil.ReadFile("../../testdata/attestation_cert_sw")
	require.NoError(t, err)

	regInfo := types.RegistrationNodeInfo{
		Certificate:   cert,
		EncryptedSeed: []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
	}

	regKeeper.SetRegistrationInfo(ctx, regInfo)

	_, err = regKeeper.RegisterNode(ctx, cert)
	require.NoError(t, err)
}
