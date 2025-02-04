package keeper_test

// Note: this is for dockernet

import (
	"fmt"

	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
)

func (s *KeeperTestSuite) TestIBCDenom() {
	chainId := "provider"
	denom := "uatom"
	for i := 0; i < 4; i++ {
		sourcePrefix := transfertypes.GetDenomPrefix("transfer", fmt.Sprintf("channel-%d", i))
		prefixedDenom := sourcePrefix + denom

		fmt.Printf("IBC_%s_CHANNEL_%d_DENOM='%s'\n", chainId, i, transfertypes.ParseDenomTrace(prefixedDenom).IBCDenom())
	}
}

// func (s *KeeperTestSuite) TestGetQueryKey() {
// 	denom := "uosmo"
// 	feeAddress := "osmo1wdplq6qjh2xruc7qqagma9ya665q6qhcxf0p96"
// 	//feeAddressBz, _ := sdk.AccAddressFromBech32(feeAddress)
// 	_, feeAddressBz, _ := bech32.DecodeAndConvert(feeAddress)
// 	// Generate the prefix and print its length and content
// 	prefix := banktypes.CreateAccountBalancesPrefix(feeAddressBz)
// 	fmt.Printf("Prefix length: %d\n", len(prefix))
// 	fmt.Printf("Prefix content: %x\n", prefix) // print in hex format

// 	// Append the denom bytes to the prefix
// 	queryData := append(prefix, []byte(denom)...)
// 	fmt.Printf("QueryData length: %d\n", len(queryData))

// 	// Encode queryData as base64
// 	base64QueryData := base64.StdEncoding.EncodeToString(queryData)
// 	fmt.Printf("Base64 QueryData: %s\n", base64QueryData)
// }
