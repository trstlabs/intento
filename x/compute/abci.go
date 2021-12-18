package compute

import (
	//"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/trstlabs/trst/x/compute/internal/keeper"
	"github.com/trstlabs/trst/x/compute/internal/types"
)

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {

	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	logger := k.Logger(ctx)
	//fmt.Printf("ABCI ENDBLOCK COMPUTE")
	// delete inactive contracts from store and its deposits
	k.IterateContractQueue(ctx, ctx.BlockHeader().Time, func(contract types.ContractInfoWithAddress) bool {

		logger.Info(
			"contract was expired",
			"contract", contract.Address.String(),
		)

		//err := k.CallAutoMsg(ctx, contract.Address)
		if contract.ContractInfo.AutoMsg != nil {
			/*	codeHash := k.GetCodeHash(ctx, contract.ContractInfo.CodeID)
				logger.Info(
					"contract codeHash %s \n", codeHash,
				)*/
			logger.Info(
				"calback sig %s \n", len(contract.ContractInfo.CallbackSig),
			)
			res, err := k.Execute(ctx, contract.Address, contract.Address, contract.ContractInfo.AutoMsg, sdk.NewCoins(sdk.NewCoin("utrst", sdk.ZeroInt())), contract.ContractInfo.CallbackSig)
			logger.Info(
				"contract err %s \n", err.Error(),
			)
			/*var msg wasmTypes.ExecuteMsg
			logger.Info(
				"message msg", msg,
			)
			addr := contract.Address.String()
			logger.Info(
				"message contract addr", addr,
			)
			msg.ContractAddr = addr
			logger.Info(
				"message contract addr", msg.ContractAddr,
			)
			var toSend wasmTypes.CosmosMsg

			msg.Msg = contract.ContractInfo.AutoMsg
			logger.Info(
				"message message",
			)
			var wasmMsg wasmTypes.WasmMsg
			logger.Info(
				"message types",
			)

			wasmMsg.Execute = &msg
			logger.Info(
				"message toSend",
			)
			toSend.Wasm = &wasmMsg
			logger.Info(
				"message complete for contract", toSend.Wasm.Execute.ContractAddr,
			)
			res, _, err := k.Dispatch(ctx, contract.Address, toSend)*/
			if err != nil {
				logger.Info(
					"Error lastMsg, creator payout", contract.Address,
				)

				err = k.ContractPayout(ctx, contract.Address)
				logger.Info(
					"Error lastMsg, err", err.Error(),
				)
			}
			logger.Info(
				"result log", res.Log[0],
			)

			k.SetContractResult(ctx, contract.Address, res)
		}

		k.RemoveFromContractQueue(ctx, contract.Address.String(), contract.ContractInfo.EndTime)
		_ = k.Delete(ctx, contract.Address)
		logger.Info(
			"Deleted contract", contract.Address,
		)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeContractExpired,
				sdk.NewAttribute(types.AttributeKeyContractAddr, contract.Address.String()),
			),
		)
		return false

		/*if err != nil {
			//fmt.Printf("contract.ContractInfo.CodeID")
			//fmt.Printf(".AttributeKeyContractAddr  is:  %s ", addr.String())
			return false
		}
		*/

	})

	return []abci.ValidatorUpdate{}
}
