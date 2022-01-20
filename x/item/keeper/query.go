package keeper

import (
	// this line is used by starport scaffolding # 1
	"github.com/trstlabs/trst/x/item/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		var (
			res []byte
			err error
		)

		switch path[0] {
		// this line is used by starport scaffolding # 2
		case types.QueryGetProfile:
			return getProfile(ctx, path[1], k, legacyQuerierCdc)

		case types.QueryListProfiles:
			return listProfiles(ctx, k, legacyQuerierCdc)

			/*	case types.QueryGetBuyer:
					return getBuyer(ctx, path[1], k, legacyQuerierCdc)

				case types.QueryListBuyer:
					return listBuyer(ctx, k, legacyQuerierCdc)*/

		case types.QueryGetItem:
			return getItem(ctx, path[1], k, legacyQuerierCdc)

		case types.QueryListItem:
			return listItem(ctx, k, legacyQuerierCdc)
		case types.QueryListListedItems:
			return listListedItems(ctx, k, legacyQuerierCdc)
		case types.QuerySellerItems:
			return sellerItems(ctx, path[1], k, legacyQuerierCdc)
		case types.QueryBuyerItems:
			return buyerItems(ctx, path[1], k, legacyQuerierCdc)

		default:
			err = sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
		}

		return res, err
	}
}
