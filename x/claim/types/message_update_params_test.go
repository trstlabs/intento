package types_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/trstlabs/intento/x/claim/types"
)

func TestMsgUpdateParams_ValidateBasic(t *testing.T) {
	validParams := types.Params{
		ClaimDenom:             "uinto",
		DurationUntilDecay:     time.Hour,
		DurationOfDecay:        time.Hour * 2,
		DurationVestingPeriods: []time.Duration{time.Hour, time.Hour},
	}
	acc, _ := sdk.AccAddressFromHexUnsafe("91e17c2a0c4d8c1b0d3d7d0a7f5c8a1c0b1a09")
	msg := types.NewMsgUpdateParams(acc.String(), validParams)
	require.NoError(t, msg.ValidateBasic())

	msgBad := types.NewMsgUpdateParams("badaddress", validParams)
	require.Error(t, msgBad.ValidateBasic())

	invalidParams := validParams
	invalidParams.ClaimDenom = ""
	msg2 := types.NewMsgUpdateParams(acc.String(), invalidParams)
	require.Error(t, msg2.ValidateBasic())
}
