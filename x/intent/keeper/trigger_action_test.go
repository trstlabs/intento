package keeper

import (
	"encoding/json"
	"testing"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/stretchr/testify/require"
)

func TestHandleContractMsgSalt(t *testing.T) {
	// Original message as it starts in the user request
	msgJson := `{
      "send": {
        "channel_id": 2,
        "timeout_height": "0",
        "timeout_timestamp": "1769385047460000000",
        "salt": "0xcfc1d1a8cceceb17cd53831d50ffd12daa0f70b2bb6265b3cd3a60501208d7dc",
        "instruction": "0x00"
      }
    }`

	execMsg := &wasmtypes.MsgExecuteContract{
		Msg: []byte(msgJson),
	}

	err := handleContractMsgSalt(execMsg)
	require.NoError(t, err)

	// We expect the result to still be a valid JSON object, not a string wrapping base64
	var result map[string]interface{}
	err = json.Unmarshal(execMsg.Msg, &result)
	require.NoError(t, err, "Msg should be unmarshalable as a JSON object")

	// Verify salt was incremented
	sendMap, ok := result["send"].(map[string]interface{})
	require.True(t, ok)
	salt, ok := sendMap["salt"].(string)
	require.True(t, ok)
	require.Equal(t, "0xcfc1d1a8cceceb17cd53831d50ffd12daa0f70b2bb6265b3cd3a60501208d7dd", salt)
}
