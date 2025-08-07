package types_test

import (
	"testing"
	"time"

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
	msg := types.NewMsgUpdateParams("cosmos1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqe36h6r4", validParams)
	require.NoError(t, msg.ValidateBasic())

	msgBad := types.NewMsgUpdateParams("badaddress", validParams)
	require.Error(t, msgBad.ValidateBasic())

	invalidParams := validParams
	invalidParams.ClaimDenom = ""
	msg2 := types.NewMsgUpdateParams("cosmos1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqe36h6r4", invalidParams)
	require.Error(t, msg2.ValidateBasic())
}
