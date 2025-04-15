package msg_registry

import (
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	elysamm "github.com/trstlabs/intento/x/intent/msg_registry/elys/amm"
	elysestaking "github.com/trstlabs/intento/x/intent/msg_registry/elys/estaking"
	elysleveragelp "github.com/trstlabs/intento/x/intent/msg_registry/elys/leveragelp"
	elysmasterchef "github.com/trstlabs/intento/x/intent/msg_registry/elys/masterchef"
	elysperpetual "github.com/trstlabs/intento/x/intent/msg_registry/elys/perpetual"
	elysstablestake "github.com/trstlabs/intento/x/intent/msg_registry/elys/stablestake"
	elystradeshield "github.com/trstlabs/intento/x/intent/msg_registry/elys/tradeshield"
	osmosisgammv1beta1 "github.com/trstlabs/intento/x/intent/msg_registry/osmosis/gamm/v1beta1"
)

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {

	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		//osmosis
		&osmosisgammv1beta1.MsgExitPool{},
		&osmosisgammv1beta1.MsgExitSwapExternAmountOut{},
		&osmosisgammv1beta1.MsgExitSwapShareAmountIn{},
		&osmosisgammv1beta1.MsgJoinPool{},
		&osmosisgammv1beta1.MsgJoinSwapExternAmountIn{},
		&osmosisgammv1beta1.MsgJoinSwapShareAmountOut{},
		&osmosisgammv1beta1.MsgSwapExactAmountIn{},
		&osmosisgammv1beta1.MsgSwapExactAmountOut{},
		//elys
		&elysamm.MsgCreatePool{},
		&elysamm.MsgJoinPool{},
		&elysamm.MsgExitPool{},
		&elysamm.MsgUpFrontSwapExactAmountIn{},
		&elysamm.MsgSwapExactAmountIn{},
		&elysamm.MsgSwapExactAmountOut{},
		&elysamm.MsgSwapByDenom{},
		&elysamm.MsgFeedMultipleExternalLiquidity{},
		&elysamm.MsgUpdatePoolParams{},
		&elysamm.MsgUpdateParams{},

		&elystradeshield.MsgCreateSpotOrder{},
		&elystradeshield.MsgUpdateSpotOrder{},
		&elystradeshield.MsgCancelSpotOrder{},
		&elystradeshield.MsgCancelSpotOrders{},
		&elystradeshield.MsgCreatePerpetualOpenOrder{},
		&elystradeshield.MsgCreatePerpetualCloseOrder{},
		&elystradeshield.MsgUpdatePerpetualOrder{},
		&elystradeshield.MsgCancelPerpetualOrder{},
		&elystradeshield.MsgCancelPerpetualOrders{},
		&elystradeshield.MsgUpdateParams{},
		&elystradeshield.MsgExecuteOrders{},

		&elysperpetual.MsgOpen{},
		&elysperpetual.MsgClose{},
		&elysperpetual.MsgUpdateParams{},
		&elysperpetual.MsgWhitelist{},
		&elysperpetual.MsgDewhitelist{},
		&elysperpetual.MsgUpdateStopLoss{},
		&elysperpetual.MsgClosePositions{},
		&elysperpetual.MsgUpdateTakeProfitPrice{},
		&elysperpetual.MsgUpdateMaxLeverageForPool{},
		&elysperpetual.MsgUpdateEnabledPools{},

		&elysestaking.MsgUpdateParams{},
		&elysestaking.MsgWithdrawReward{},
		&elysestaking.MsgWithdrawElysStakingRewards{},
		&elysestaking.MsgWithdrawAllRewards{},

		&elysleveragelp.MsgOpen{},
		&elysleveragelp.MsgClose{},
		&elysleveragelp.MsgClaimRewards{},
		&elysleveragelp.MsgUpdateParams{},
		&elysleveragelp.MsgAddPool{},
		&elysleveragelp.MsgRemovePool{},
		&elysleveragelp.MsgWhitelist{},
		&elysleveragelp.MsgDewhitelist{},
		&elysleveragelp.MsgUpdateStopLoss{},
		&elysleveragelp.MsgClosePositions{},
		&elysleveragelp.MsgUpdatePool{},
		&elysleveragelp.MsgUpdateEnabledPools{},

		&elysmasterchef.MsgAddExternalRewardDenom{},
		&elysmasterchef.MsgAddExternalIncentive{},
		&elysmasterchef.MsgUpdateParams{},
		&elysmasterchef.MsgUpdatePoolMultipliers{},
		&elysmasterchef.MsgClaimRewards{},
		&elysmasterchef.MsgTogglePoolEdenRewards{},

		&elysstablestake.MsgBond{},
		&elysstablestake.MsgUnbond{},
		&elysstablestake.MsgUpdateParams{},
		&elysstablestake.MsgAddPool{},
		&elysstablestake.MsgUpdatePool{},
	)
}

type RawContractMessage []byte
