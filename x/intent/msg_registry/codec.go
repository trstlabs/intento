package msg_registry

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	proto "github.com/cosmos/gogoproto/proto"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
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
	// Cosmos SDK bank and staking messages
	sdk.MsgTypeURL(&banktypes.MsgSend{}): {
		func() proto.Message { return &banktypes.MsgSendResponse{} },
		types.KeyFlowIncentiveForLowGas,
	},
	sdk.MsgTypeURL(&stakingtypes.MsgDelegate{}): {
		func() proto.Message { return &stakingtypes.MsgDelegateResponse{} },
		types.KeyFlowIncentiveForLowGas,
	},
	sdk.MsgTypeURL(&stakingtypes.MsgUndelegate{}): {
		func() proto.Message { return &stakingtypes.MsgUndelegateResponse{} },
		types.KeyFlowIncentiveForLowGas,
	},
	sdk.MsgTypeURL(&stakingtypes.MsgBeginRedelegate{}): {
		func() proto.Message { return &stakingtypes.MsgBeginRedelegateResponse{} },
		types.KeyFlowIncentiveForLowGas,
	},

	sdk.MsgTypeURL(&distributiontypes.MsgWithdrawDelegatorReward{}): {
		func() proto.Message { return &distributiontypes.MsgWithdrawDelegatorRewardResponse{} },
		types.KeyFlowIncentiveForLowGas,
	},

	sdk.MsgTypeURL(&authztypes.MsgExec{}): {
		func() proto.Message { return &authztypes.MsgExecResponse{} },
		types.KeyFlowIncentiveForAuthzExec,
	},

	sdk.MsgTypeURL(&wasmtypes.MsgExecuteContract{}): {
		func() proto.Message { return &wasmtypes.MsgExecuteContractResponse{} },
		types.KeyFlowIncentiveForHighGas,
	},

	sdk.MsgTypeURL(&wasmtypes.MsgInstantiateContract{}): {
		func() proto.Message { return &wasmtypes.MsgInstantiateContractResponse{} },
		types.KeyFlowIncentiveForHighGas,
	},

	sdk.MsgTypeURL(&ibctransfertypes.MsgTransfer{}): {
		func() proto.Message { return &ibctransfertypes.MsgTransferResponse{} },
		types.KeyFlowIncentiveForLowGas,
	},

	// Osmosis Gamm
	sdk.MsgTypeURL(&osmosisgammv1beta1.MsgExitPool{}):                {func() proto.Message { return &osmosisgammv1beta1.MsgExitPoolResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&osmosisgammv1beta1.MsgExitSwapExternAmountOut{}): {func() proto.Message { return &osmosisgammv1beta1.MsgExitSwapExternAmountOutResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&osmosisgammv1beta1.MsgExitSwapShareAmountIn{}):   {func() proto.Message { return &osmosisgammv1beta1.MsgExitSwapShareAmountInResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&osmosisgammv1beta1.MsgJoinPool{}):                {func() proto.Message { return &osmosisgammv1beta1.MsgJoinPoolResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&osmosisgammv1beta1.MsgJoinSwapExternAmountIn{}):  {func() proto.Message { return &osmosisgammv1beta1.MsgJoinSwapExternAmountInResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&osmosisgammv1beta1.MsgJoinSwapShareAmountOut{}):  {func() proto.Message { return &osmosisgammv1beta1.MsgJoinSwapShareAmountOutResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&osmosisgammv1beta1.MsgSwapExactAmountIn{}):       {func() proto.Message { return &osmosisgammv1beta1.MsgSwapExactAmountInResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&osmosisgammv1beta1.MsgSwapExactAmountOut{}):      {func() proto.Message { return &osmosisgammv1beta1.MsgSwapExactAmountOutResponse{} }, types.KeyFlowIncentiveForMediumGas},

	// Elys AMM
	sdk.MsgTypeURL(&elysamm.MsgCreatePool{}):                    {func() proto.Message { return &elysamm.MsgCreatePoolResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysamm.MsgJoinPool{}):                      {func() proto.Message { return &elysamm.MsgJoinPoolResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysamm.MsgExitPool{}):                      {func() proto.Message { return &elysamm.MsgExitPoolResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysamm.MsgUpFrontSwapExactAmountIn{}):      {func() proto.Message { return &elysamm.MsgUpFrontSwapExactAmountInResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysamm.MsgSwapExactAmountIn{}):             {func() proto.Message { return &elysamm.MsgSwapExactAmountInResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysamm.MsgSwapExactAmountOut{}):            {func() proto.Message { return &elysamm.MsgSwapExactAmountOutResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysamm.MsgSwapByDenom{}):                   {func() proto.Message { return &elysamm.MsgSwapByDenomResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysamm.MsgFeedMultipleExternalLiquidity{}): {func() proto.Message { return &elysamm.MsgFeedMultipleExternalLiquidityResponse{} }, types.KeyFlowIncentiveForMediumGas},

	// Elys TradeShield
	sdk.MsgTypeURL(&elystradeshield.MsgCreateSpotOrder{}):           {func() proto.Message { return &elystradeshield.MsgCreateSpotOrderResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elystradeshield.MsgUpdateSpotOrder{}):           {func() proto.Message { return &elystradeshield.MsgUpdateSpotOrderResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elystradeshield.MsgCancelSpotOrder{}):           {func() proto.Message { return &elystradeshield.MsgCancelSpotOrderResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elystradeshield.MsgCancelSpotOrders{}):          {func() proto.Message { return &elystradeshield.MsgCancelSpotOrdersResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elystradeshield.MsgCreatePerpetualOpenOrder{}):  {func() proto.Message { return &elystradeshield.MsgCreatePerpetualOpenOrderResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elystradeshield.MsgCreatePerpetualCloseOrder{}): {func() proto.Message { return &elystradeshield.MsgCreatePerpetualCloseOrderResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elystradeshield.MsgUpdatePerpetualOrder{}):      {func() proto.Message { return &elystradeshield.MsgUpdatePerpetualOrderResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elystradeshield.MsgCancelPerpetualOrder{}):      {func() proto.Message { return &elystradeshield.MsgCancelPerpetualOrderResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elystradeshield.MsgCancelPerpetualOrders{}):     {func() proto.Message { return &elystradeshield.MsgCancelPerpetualOrdersResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elystradeshield.MsgExecuteOrders{}):             {func() proto.Message { return &elystradeshield.MsgExecuteOrdersResponse{} }, types.KeyFlowIncentiveForMediumGas},

	// Elys Perpetual
	sdk.MsgTypeURL(&elysperpetual.MsgOpen{}):                     {func() proto.Message { return &elysperpetual.MsgOpenResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysperpetual.MsgClose{}):                    {func() proto.Message { return &elysperpetual.MsgCloseResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysperpetual.MsgUpdateParams{}):             {func() proto.Message { return &elysperpetual.MsgUpdateParamsResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysperpetual.MsgWhitelist{}):                {func() proto.Message { return &elysperpetual.MsgWhitelistResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysperpetual.MsgDewhitelist{}):              {func() proto.Message { return &elysperpetual.MsgDewhitelistResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysperpetual.MsgUpdateStopLoss{}):           {func() proto.Message { return &elysperpetual.MsgUpdateStopLossResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysperpetual.MsgClosePositions{}):           {func() proto.Message { return &elysperpetual.MsgClosePositionsResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysperpetual.MsgUpdateTakeProfitPrice{}):    {func() proto.Message { return &elysperpetual.MsgUpdateTakeProfitPriceResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysperpetual.MsgUpdateMaxLeverageForPool{}): {func() proto.Message { return &elysperpetual.MsgUpdateMaxLeverageForPoolResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysperpetual.MsgUpdateEnabledPools{}):       {func() proto.Message { return &elysperpetual.MsgUpdateEnabledPoolsResponse{} }, types.KeyFlowIncentiveForMediumGas},

	// Elys eStaking
	sdk.MsgTypeURL(&elysestaking.MsgWithdrawReward{}):             {func() proto.Message { return &elysestaking.MsgWithdrawRewardResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysestaking.MsgWithdrawElysStakingRewards{}): {func() proto.Message { return &elysestaking.MsgWithdrawElysStakingRewardsResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysestaking.MsgWithdrawAllRewards{}):         {func() proto.Message { return &elysestaking.MsgWithdrawAllRewardsResponse{} }, types.KeyFlowIncentiveForMediumGas},

	// Elys Leverage LP
	sdk.MsgTypeURL(&elysleveragelp.MsgOpen{}):               {func() proto.Message { return &elysleveragelp.MsgOpenResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysleveragelp.MsgClose{}):              {func() proto.Message { return &elysleveragelp.MsgCloseResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysleveragelp.MsgClaimRewards{}):       {func() proto.Message { return &elysleveragelp.MsgClaimRewardsResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysleveragelp.MsgUpdateParams{}):       {func() proto.Message { return &elysleveragelp.MsgUpdateParamsResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysleveragelp.MsgAddPool{}):            {func() proto.Message { return &elysleveragelp.MsgAddPoolResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysleveragelp.MsgRemovePool{}):         {func() proto.Message { return &elysleveragelp.MsgRemovePoolResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysleveragelp.MsgWhitelist{}):          {func() proto.Message { return &elysleveragelp.MsgWhitelistResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysleveragelp.MsgDewhitelist{}):        {func() proto.Message { return &elysleveragelp.MsgDewhitelistResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysleveragelp.MsgUpdateStopLoss{}):     {func() proto.Message { return &elysleveragelp.MsgUpdateStopLossResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysleveragelp.MsgClosePositions{}):     {func() proto.Message { return &elysleveragelp.MsgClosePositionsResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysleveragelp.MsgUpdatePool{}):         {func() proto.Message { return &elysleveragelp.MsgUpdatePoolResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysleveragelp.MsgUpdateEnabledPools{}): {func() proto.Message { return &elysleveragelp.MsgUpdateEnabledPoolsResponse{} }, types.KeyFlowIncentiveForMediumGas},

	// Elys MasterChef
	sdk.MsgTypeURL(&elysmasterchef.MsgAddExternalRewardDenom{}): {func() proto.Message { return &elysmasterchef.MsgAddExternalRewardDenomResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysmasterchef.MsgAddExternalIncentive{}):   {func() proto.Message { return &elysmasterchef.MsgAddExternalIncentiveResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysmasterchef.MsgUpdatePoolMultipliers{}):  {func() proto.Message { return &elysmasterchef.MsgUpdatePoolMultipliersResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysmasterchef.MsgClaimRewards{}):           {func() proto.Message { return &elysmasterchef.MsgClaimRewardsResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysmasterchef.MsgTogglePoolEdenRewards{}):  {func() proto.Message { return &elysmasterchef.MsgTogglePoolEdenRewardsResponse{} }, types.KeyFlowIncentiveForMediumGas},

	// Elys StableStake
	sdk.MsgTypeURL(&elysstablestake.MsgBond{}):       {func() proto.Message { return &elysstablestake.MsgBondResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysstablestake.MsgUnbond{}):     {func() proto.Message { return &elysstablestake.MsgUnbondResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysstablestake.MsgAddPool{}):    {func() proto.Message { return &elysstablestake.MsgAddPoolResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elysstablestake.MsgUpdatePool{}): {func() proto.Message { return &elysstablestake.MsgUpdatePoolResponse{} }, types.KeyFlowIncentiveForMediumGas},

	// Elys Commitment
	sdk.MsgTypeURL(&elyscommitment.MsgCommitClaimedRewards{}): {func() proto.Message { return &elyscommitment.MsgCommitClaimedRewardsResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elyscommitment.MsgUncommitTokens{}):       {func() proto.Message { return &elyscommitment.MsgUncommitTokensResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elyscommitment.MsgVest{}):                 {func() proto.Message { return &elyscommitment.MsgVestResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elyscommitment.MsgCancelVest{}):           {func() proto.Message { return &elyscommitment.MsgCancelVestResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elyscommitment.MsgClaimVesting{}):         {func() proto.Message { return &elyscommitment.MsgClaimVestingResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elyscommitment.MsgUpdateEnableVestNow{}):  {func() proto.Message { return &elyscommitment.MsgUpdateEnableVestNowResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elyscommitment.MsgVestNow{}):              {func() proto.Message { return &elyscommitment.MsgVestNowResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elyscommitment.MsgVestLiquid{}):           {func() proto.Message { return &elyscommitment.MsgVestLiquidResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elyscommitment.MsgStake{}):                {func() proto.Message { return &elyscommitment.MsgStakeResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elyscommitment.MsgUnstake{}):              {func() proto.Message { return &elyscommitment.MsgUnstakeResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elyscommitment.MsgClaimKol{}):             {func() proto.Message { return &elyscommitment.MsgClaimKolResponse{} }, types.KeyFlowIncentiveForMediumGas},
	sdk.MsgTypeURL(&elyscommitment.MsgClaimRewardProgram{}):   {func() proto.Message { return &elyscommitment.MsgClaimRewardProgramResponse{} }, types.KeyFlowIncentiveForMediumGas},

	// EVM
	sdk.MsgTypeURL(&cosmosevm.MsgEthereumTx{}): {func() proto.Message { return &cosmosevm.MsgEthereumTxResponse{} }, types.KeyFlowIncentiveForHighGas},
}
