package keeper_test

// Note: this is for dockernet

import (
	"fmt"

	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
)

func (s *KeeperTestSuite) TestIBCDenom() {
	chainId := "osmosis-test-5"
	denom := "ibc/DE6792CF9E521F6AD6E9A4BDF6225C9571A3B74ACC0A529F92BC5122A39D2E58"
	for i := 0; i < 4; i++ {
		sourcePrefix := transfertypes.GetDenomPrefix("transfer", fmt.Sprintf("channel-%d", i))
		prefixedDenom := sourcePrefix + denom

		fmt.Printf("IBC_%s_CHANNEL_%d_DENOM='%s'\n", chainId, i, transfertypes.ParseDenomTrace(prefixedDenom).IBCDenom())
	}
}
