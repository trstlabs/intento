package msg_registry

import (
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	proto "github.com/cosmos/gogoproto/proto"
	cosmosevm "github.com/trstlabs/intento/x/intent/msg_registry/cosmos/evm/v1"
	elysamm "github.com/trstlabs/intento/x/intent/msg_registry/elys/amm"
	elyscommitment "github.com/trstlabs/intento/x/intent/msg_registry/elys/commitment"
	elysestaking "github.com/trstlabs/intento/x/intent/msg_registry/elys/estaking"
	elysleveragelp "github.com/trstlabs/intento/x/intent/msg_registry/elys/leveragelp"
	elysmasterchef "github.com/trstlabs/intento/x/intent/msg_registry/elys/masterchef"
	elysperpetual "github.com/trstlabs/intento/x/intent/msg_registry/elys/perpetual"
	elysstablestake "github.com/trstlabs/intento/x/intent/msg_registry/elys/stablestake"
	elystradeshield "github.com/trstlabs/intento/x/intent/msg_registry/elys/tradeshield"
	osmosisgammv1beta1 "github.com/trstlabs/intento/x/intent/msg_registry/osmosis/gamm/v1beta1"
	types "github.com/trstlabs/intento/x/intent/types"
)

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {

	registry.RegisterInterface(
		"cosmos.evm.vm.v1.TxData",
		(*cosmosevm.TxData)(nil),
		&cosmosevm.DynamicFeeTx{},
		&cosmosevm.AccessListTx{},
		&cosmosevm.LegacyTx{},
	)

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

		&elystradeshield.MsgCreateSpotOrder{},
		&elystradeshield.MsgUpdateSpotOrder{},
		&elystradeshield.MsgCancelSpotOrder{},
		&elystradeshield.MsgCancelSpotOrders{},
		&elystradeshield.MsgCreatePerpetualOpenOrder{},
		&elystradeshield.MsgCreatePerpetualCloseOrder{},
		&elystradeshield.MsgUpdatePerpetualOrder{},
		&elystradeshield.MsgCancelPerpetualOrder{},
		&elystradeshield.MsgCancelPerpetualOrders{},
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
		&elysmasterchef.MsgUpdatePoolMultipliers{},
		&elysmasterchef.MsgClaimRewards{},
		&elysmasterchef.MsgTogglePoolEdenRewards{},

		&elysstablestake.MsgBond{},
		&elysstablestake.MsgUnbond{},
		&elysstablestake.MsgAddPool{},
		&elysstablestake.MsgUpdatePool{},

		&elyscommitment.MsgCommitClaimedRewards{},
		&elyscommitment.MsgUncommitTokens{},
		&elyscommitment.MsgVest{},
		&elyscommitment.MsgCancelVest{},
		&elyscommitment.MsgClaimVesting{},
		&elyscommitment.MsgUpdateEnableVestNow{},
		&elyscommitment.MsgVestNow{},
		&elyscommitment.MsgVestLiquid{},
		&elyscommitment.MsgStake{},
		&elyscommitment.MsgUnstake{},
		&elyscommitment.MsgClaimKol{},
		&elyscommitment.MsgClaimRewardProgram{},

		// EVM
		&cosmosevm.MsgEthereumTx{},
	)

	// For comparison and feedback Loop logic
	registry.RegisterImplementations(
		(*proto.Message)(nil),
		&cosmosevm.MsgEthereumTxResponse{},
		&elysestaking.MsgWithdrawElysStakingRewardsResponse{},
		&elysestaking.MsgWithdrawAllRewardsResponse{},
		&elysestaking.MsgWithdrawRewardResponse{},
		&elysamm.MsgSwapExactAmountInResponse{},
		&elysamm.MsgSwapExactAmountOutResponse{},
		&elysamm.MsgFeedMultipleExternalLiquidityResponse{},
		&elysamm.MsgFeedMultipleExternalLiquidityResponse{},
		&elysamm.MsgSwapByDenomResponse{},
		&elyscommitment.MsgStakeResponse{},
		&osmosisgammv1beta1.MsgSwapExactAmountInResponse{},
		&osmosisgammv1beta1.MsgSwapExactAmountOutResponse{},
		&osmosisgammv1beta1.MsgJoinPoolResponse{},
		&osmosisgammv1beta1.MsgExitPoolResponse{},
	)

}

type RawContractMessage []byte

var MsgRegistry = map[string]struct {
	NewResponse func() proto.Message
	RewardType  int
}{
	// Osmosis Gamm
	sdk.MsgTypeURL(&osmosisgammv1beta1.MsgExitPool{}):                {func() proto.Message { return &osmosisgammv1beta1.MsgExitPoolResponse{} }, types.KeyFlowIncentiveForOsmoTx},
	sdk.MsgTypeURL(&osmosisgammv1beta1.MsgExitSwapExternAmountOut{}): {func() proto.Message { return &osmosisgammv1beta1.MsgExitSwapExternAmountOutResponse{} }, types.KeyFlowIncentiveForOsmoTx},
	sdk.MsgTypeURL(&osmosisgammv1beta1.MsgExitSwapShareAmountIn{}):   {func() proto.Message { return &osmosisgammv1beta1.MsgExitSwapShareAmountInResponse{} }, types.KeyFlowIncentiveForOsmoTx},
	sdk.MsgTypeURL(&osmosisgammv1beta1.MsgJoinPool{}):                {func() proto.Message { return &osmosisgammv1beta1.MsgJoinPoolResponse{} }, types.KeyFlowIncentiveForOsmoTx},
	sdk.MsgTypeURL(&osmosisgammv1beta1.MsgJoinSwapExternAmountIn{}):  {func() proto.Message { return &osmosisgammv1beta1.MsgJoinSwapExternAmountInResponse{} }, types.KeyFlowIncentiveForOsmoTx},
	sdk.MsgTypeURL(&osmosisgammv1beta1.MsgJoinSwapShareAmountOut{}):  {func() proto.Message { return &osmosisgammv1beta1.MsgJoinSwapShareAmountOutResponse{} }, types.KeyFlowIncentiveForOsmoTx},
	sdk.MsgTypeURL(&osmosisgammv1beta1.MsgSwapExactAmountIn{}):       {func() proto.Message { return &osmosisgammv1beta1.MsgSwapExactAmountInResponse{} }, types.KeyFlowIncentiveForOsmoTx},
	sdk.MsgTypeURL(&osmosisgammv1beta1.MsgSwapExactAmountOut{}):      {func() proto.Message { return &osmosisgammv1beta1.MsgSwapExactAmountOutResponse{} }, types.KeyFlowIncentiveForOsmoTx},

	// Elys AMM
	sdk.MsgTypeURL(&elysamm.MsgCreatePool{}):                    {func() proto.Message { return &elysamm.MsgCreatePoolResponse{} }, -1},
	sdk.MsgTypeURL(&elysamm.MsgJoinPool{}):                      {func() proto.Message { return &elysamm.MsgJoinPoolResponse{} }, -1},
	sdk.MsgTypeURL(&elysamm.MsgExitPool{}):                      {func() proto.Message { return &elysamm.MsgExitPoolResponse{} }, -1},
	sdk.MsgTypeURL(&elysamm.MsgUpFrontSwapExactAmountIn{}):      {func() proto.Message { return &elysamm.MsgUpFrontSwapExactAmountInResponse{} }, -1},
	sdk.MsgTypeURL(&elysamm.MsgSwapExactAmountIn{}):             {func() proto.Message { return &elysamm.MsgSwapExactAmountInResponse{} }, -1},
	sdk.MsgTypeURL(&elysamm.MsgSwapExactAmountOut{}):            {func() proto.Message { return &elysamm.MsgSwapExactAmountOutResponse{} }, -1},
	sdk.MsgTypeURL(&elysamm.MsgSwapByDenom{}):                   {func() proto.Message { return &elysamm.MsgSwapByDenomResponse{} }, -1},
	sdk.MsgTypeURL(&elysamm.MsgFeedMultipleExternalLiquidity{}): {func() proto.Message { return &elysamm.MsgFeedMultipleExternalLiquidityResponse{} }, -1},

	// Elys TradeShield
	sdk.MsgTypeURL(&elystradeshield.MsgCreateSpotOrder{}):           {func() proto.Message { return &elystradeshield.MsgCreateSpotOrderResponse{} }, -1},
	sdk.MsgTypeURL(&elystradeshield.MsgUpdateSpotOrder{}):           {func() proto.Message { return &elystradeshield.MsgUpdateSpotOrderResponse{} }, -1},
	sdk.MsgTypeURL(&elystradeshield.MsgCancelSpotOrder{}):           {func() proto.Message { return &elystradeshield.MsgCancelSpotOrderResponse{} }, -1},
	sdk.MsgTypeURL(&elystradeshield.MsgCancelSpotOrders{}):          {func() proto.Message { return &elystradeshield.MsgCancelSpotOrdersResponse{} }, -1},
	sdk.MsgTypeURL(&elystradeshield.MsgCreatePerpetualOpenOrder{}):  {func() proto.Message { return &elystradeshield.MsgCreatePerpetualOpenOrderResponse{} }, -1},
	sdk.MsgTypeURL(&elystradeshield.MsgCreatePerpetualCloseOrder{}): {func() proto.Message { return &elystradeshield.MsgCreatePerpetualCloseOrderResponse{} }, -1},
	sdk.MsgTypeURL(&elystradeshield.MsgUpdatePerpetualOrder{}):      {func() proto.Message { return &elystradeshield.MsgUpdatePerpetualOrderResponse{} }, -1},
	sdk.MsgTypeURL(&elystradeshield.MsgCancelPerpetualOrder{}):      {func() proto.Message { return &elystradeshield.MsgCancelPerpetualOrderResponse{} }, -1},
	sdk.MsgTypeURL(&elystradeshield.MsgCancelPerpetualOrders{}):     {func() proto.Message { return &elystradeshield.MsgCancelPerpetualOrdersResponse{} }, -1},
	sdk.MsgTypeURL(&elystradeshield.MsgExecuteOrders{}):             {func() proto.Message { return &elystradeshield.MsgExecuteOrdersResponse{} }, -1},

	// Elys Perpetual
	sdk.MsgTypeURL(&elysperpetual.MsgOpen{}):                     {func() proto.Message { return &elysperpetual.MsgOpenResponse{} }, -1},
	sdk.MsgTypeURL(&elysperpetual.MsgClose{}):                    {func() proto.Message { return &elysperpetual.MsgCloseResponse{} }, -1},
	sdk.MsgTypeURL(&elysperpetual.MsgUpdateParams{}):             {func() proto.Message { return &elysperpetual.MsgUpdateParamsResponse{} }, -1},
	sdk.MsgTypeURL(&elysperpetual.MsgWhitelist{}):                {func() proto.Message { return &elysperpetual.MsgWhitelistResponse{} }, -1},
	sdk.MsgTypeURL(&elysperpetual.MsgDewhitelist{}):              {func() proto.Message { return &elysperpetual.MsgDewhitelistResponse{} }, -1},
	sdk.MsgTypeURL(&elysperpetual.MsgUpdateStopLoss{}):           {func() proto.Message { return &elysperpetual.MsgUpdateStopLossResponse{} }, -1},
	sdk.MsgTypeURL(&elysperpetual.MsgClosePositions{}):           {func() proto.Message { return &elysperpetual.MsgClosePositionsResponse{} }, -1},
	sdk.MsgTypeURL(&elysperpetual.MsgUpdateTakeProfitPrice{}):    {func() proto.Message { return &elysperpetual.MsgUpdateTakeProfitPriceResponse{} }, -1},
	sdk.MsgTypeURL(&elysperpetual.MsgUpdateMaxLeverageForPool{}): {func() proto.Message { return &elysperpetual.MsgUpdateMaxLeverageForPoolResponse{} }, -1},
	sdk.MsgTypeURL(&elysperpetual.MsgUpdateEnabledPools{}):       {func() proto.Message { return &elysperpetual.MsgUpdateEnabledPoolsResponse{} }, -1},

	// Elys eStaking
	sdk.MsgTypeURL(&elysestaking.MsgWithdrawReward{}):             {func() proto.Message { return &elysestaking.MsgWithdrawRewardResponse{} }, -1},
	sdk.MsgTypeURL(&elysestaking.MsgWithdrawElysStakingRewards{}): {func() proto.Message { return &elysestaking.MsgWithdrawElysStakingRewardsResponse{} }, -1},
	sdk.MsgTypeURL(&elysestaking.MsgWithdrawAllRewards{}):         {func() proto.Message { return &elysestaking.MsgWithdrawAllRewardsResponse{} }, -1},

	// Elys Leverage LP
	sdk.MsgTypeURL(&elysleveragelp.MsgOpen{}):               {func() proto.Message { return &elysleveragelp.MsgOpenResponse{} }, -1},
	sdk.MsgTypeURL(&elysleveragelp.MsgClose{}):              {func() proto.Message { return &elysleveragelp.MsgCloseResponse{} }, -1},
	sdk.MsgTypeURL(&elysleveragelp.MsgClaimRewards{}):       {func() proto.Message { return &elysleveragelp.MsgClaimRewardsResponse{} }, -1},
	sdk.MsgTypeURL(&elysleveragelp.MsgUpdateParams{}):       {func() proto.Message { return &elysleveragelp.MsgUpdateParamsResponse{} }, -1},
	sdk.MsgTypeURL(&elysleveragelp.MsgAddPool{}):            {func() proto.Message { return &elysleveragelp.MsgAddPoolResponse{} }, -1},
	sdk.MsgTypeURL(&elysleveragelp.MsgRemovePool{}):         {func() proto.Message { return &elysleveragelp.MsgRemovePoolResponse{} }, -1},
	sdk.MsgTypeURL(&elysleveragelp.MsgWhitelist{}):          {func() proto.Message { return &elysleveragelp.MsgWhitelistResponse{} }, -1},
	sdk.MsgTypeURL(&elysleveragelp.MsgDewhitelist{}):        {func() proto.Message { return &elysleveragelp.MsgDewhitelistResponse{} }, -1},
	sdk.MsgTypeURL(&elysleveragelp.MsgUpdateStopLoss{}):     {func() proto.Message { return &elysleveragelp.MsgUpdateStopLossResponse{} }, -1},
	sdk.MsgTypeURL(&elysleveragelp.MsgClosePositions{}):     {func() proto.Message { return &elysleveragelp.MsgClosePositionsResponse{} }, -1},
	sdk.MsgTypeURL(&elysleveragelp.MsgUpdatePool{}):         {func() proto.Message { return &elysleveragelp.MsgUpdatePoolResponse{} }, -1},
	sdk.MsgTypeURL(&elysleveragelp.MsgUpdateEnabledPools{}): {func() proto.Message { return &elysleveragelp.MsgUpdateEnabledPoolsResponse{} }, -1},

	// Elys MasterChef
	sdk.MsgTypeURL(&elysmasterchef.MsgAddExternalRewardDenom{}): {func() proto.Message { return &elysmasterchef.MsgAddExternalRewardDenomResponse{} }, -1},
	sdk.MsgTypeURL(&elysmasterchef.MsgAddExternalIncentive{}):   {func() proto.Message { return &elysmasterchef.MsgAddExternalIncentiveResponse{} }, -1},
	sdk.MsgTypeURL(&elysmasterchef.MsgUpdatePoolMultipliers{}):  {func() proto.Message { return &elysmasterchef.MsgUpdatePoolMultipliersResponse{} }, -1},
	sdk.MsgTypeURL(&elysmasterchef.MsgClaimRewards{}):           {func() proto.Message { return &elysmasterchef.MsgClaimRewardsResponse{} }, -1},
	sdk.MsgTypeURL(&elysmasterchef.MsgTogglePoolEdenRewards{}):  {func() proto.Message { return &elysmasterchef.MsgTogglePoolEdenRewardsResponse{} }, -1},

	// Elys StableStake
	sdk.MsgTypeURL(&elysstablestake.MsgBond{}):       {func() proto.Message { return &elysstablestake.MsgBondResponse{} }, -1},
	sdk.MsgTypeURL(&elysstablestake.MsgUnbond{}):     {func() proto.Message { return &elysstablestake.MsgUnbondResponse{} }, -1},
	sdk.MsgTypeURL(&elysstablestake.MsgAddPool{}):    {func() proto.Message { return &elysstablestake.MsgAddPoolResponse{} }, -1},
	sdk.MsgTypeURL(&elysstablestake.MsgUpdatePool{}): {func() proto.Message { return &elysstablestake.MsgUpdatePoolResponse{} }, -1},

	// Elys Commitment
	sdk.MsgTypeURL(&elyscommitment.MsgCommitClaimedRewards{}): {func() proto.Message { return &elyscommitment.MsgCommitClaimedRewardsResponse{} }, -1},
	sdk.MsgTypeURL(&elyscommitment.MsgUncommitTokens{}):       {func() proto.Message { return &elyscommitment.MsgUncommitTokensResponse{} }, -1},
	sdk.MsgTypeURL(&elyscommitment.MsgVest{}):                 {func() proto.Message { return &elyscommitment.MsgVestResponse{} }, -1},
	sdk.MsgTypeURL(&elyscommitment.MsgCancelVest{}):           {func() proto.Message { return &elyscommitment.MsgCancelVestResponse{} }, -1},
	sdk.MsgTypeURL(&elyscommitment.MsgClaimVesting{}):         {func() proto.Message { return &elyscommitment.MsgClaimVestingResponse{} }, -1},
	sdk.MsgTypeURL(&elyscommitment.MsgUpdateEnableVestNow{}):  {func() proto.Message { return &elyscommitment.MsgUpdateEnableVestNowResponse{} }, -1},
	sdk.MsgTypeURL(&elyscommitment.MsgVestNow{}):              {func() proto.Message { return &elyscommitment.MsgVestNowResponse{} }, -1},
	sdk.MsgTypeURL(&elyscommitment.MsgVestLiquid{}):           {func() proto.Message { return &elyscommitment.MsgVestLiquidResponse{} }, -1},
	sdk.MsgTypeURL(&elyscommitment.MsgStake{}):                {func() proto.Message { return &elyscommitment.MsgStakeResponse{} }, -1},
	sdk.MsgTypeURL(&elyscommitment.MsgUnstake{}):              {func() proto.Message { return &elyscommitment.MsgUnstakeResponse{} }, -1},
	sdk.MsgTypeURL(&elyscommitment.MsgClaimKol{}):             {func() proto.Message { return &elyscommitment.MsgClaimKolResponse{} }, -1},
	sdk.MsgTypeURL(&elyscommitment.MsgClaimRewardProgram{}):   {func() proto.Message { return &elyscommitment.MsgClaimRewardProgramResponse{} }, -1},

	// EVM
	sdk.MsgTypeURL(&cosmosevm.MsgEthereumTx{}): {func() proto.Message { return &cosmosevm.MsgEthereumTxResponse{} }, -1},
}
