package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmwasm "github.com/danieljdd/trst/go-cosmwasm"
	"github.com/danieljdd/trst/x/compute/internal/types"
)

var (
	CostHumanize  = 5 * types.GasMultiplier
	CostCanonical = 4 * types.GasMultiplier
)

func humanAddress(canon []byte) (string, uint64, error) {
	/* AddrLen not declared by package Types

	if len(canon) != sdk.AddrLen {
		return "", CostHumanize, fmt.Errorf("Expected %d byte address", sdk.AddrLen)
	}*/
	return sdk.AccAddress(canon).String(), CostHumanize, nil
}

func canonicalAddress(human string) ([]byte, uint64, error) {
	bz, err := sdk.AccAddressFromBech32(human)
	return bz, CostCanonical, err
}

var cosmwasmAPI = cosmwasm.GoAPI{
	HumanAddress:     humanAddress,
	CanonicalAddress: canonicalAddress,
}
