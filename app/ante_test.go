package app_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/trstlabs/intento/app"
)

func TestGateCreateValidatorAnteHandler(t *testing.T) {
	testCases := []struct {
		name           string
		disable        bool
		msg            sdk.Msg
		expectErr      bool
		expectedErrStr string
	}{
		{
			name:      "Allowed: Non-CreateValidator Msg",
			disable:   false,
			msg:       testdata.NewTestMsg(),
			expectErr: false,
		},
		{
			name:           "Blocked: CreateValidator Msg when gatekeeping enabled",
			disable:        false,
			msg:            &stakingtypes.MsgCreateValidator{},
			expectErr:      true,
			expectedErrStr: "MsgCreateValidator is gated by governance",
		},
		{
			name:      "Allowed: CreateValidator Msg when gatekeeping disabled",
			disable:   true,
			msg:       &stakingtypes.MsgCreateValidator{},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Mock AnteHandler that always succeeds
			next := func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
				return ctx, nil
			}

			handler := app.GateCreateValidatorAnteHandler(tc.disable)

			// Create a dummy context
			ctx := sdk.Context{}

			// Create a mock tx
			// txBuilder := client.TxBuilder(nil) // Unused
			// Using a simpler approach: create a mock Tx that implements sdk.Tx
			tx := MockTx{Msgs: []sdk.Msg{tc.msg}}

			_, err := handler.AnteHandle(ctx, tx, false, next)

			if tc.expectErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErrStr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// MockTx implements sdk.Tx
type MockTx struct {
	Msgs []sdk.Msg
}

func (tx MockTx) GetMsgs() []sdk.Msg {
	return tx.Msgs
}

func (tx MockTx) GetMsgsV2() ([]protoreflect.ProtoMessage, error) {
	return nil, nil
}

func (tx MockTx) ValidateBasic() error {
	return nil
}
