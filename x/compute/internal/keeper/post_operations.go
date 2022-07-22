package keeper

import (

	//"encoding/json"
	"encoding/json"
	"fmt"

	//"log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	wasmTypes "github.com/trstlabs/trst/go-cosmwasm/types"
	"github.com/trstlabs/trst/x/compute/internal/types"
)

// handleContractResponse processes the contract response data by emitting events and sending sub-/messages.
func (k *Keeper) handleContractResponse(
	ctx sdk.Context,
	contractAddr sdk.AccAddress,
	ibcPort string,
	resp wasmTypes.Response,
	msgs []wasmTypes.SubMsg,
	evts wasmTypes.Events,
	data []byte,
	// original TX in order to extract the first 64bytes of signing info
	ogTx []byte,
	// sigInfo of the initial message that triggered the original contract call
	// This is used mainly in replies in order to decrypt their data.
	ogSigInfo wasmTypes.VerificationInfo,
) ([]byte, error) {
	events := types.ContractLogsToSdkEvents(resp.Attributes, contractAddr)
	ctx.EventManager().EmitEvents(events)

	if len(evts) > 0 {
		customEvents, err := types.NewCustomEvents(evts, contractAddr)
		if err != nil {
			return nil, err
		}
		ctx.EventManager().EmitEvents(customEvents)
	}

	responseHandler := NewContractResponseHandler(NewMessageDispatcher(k.messenger, k))
	return responseHandler.Handle(ctx, contractAddr, ibcPort, msgs, data, ogTx, ogSigInfo)
}

type MsgDispatcher interface {
	DispatchSubmessages(ctx sdk.Context, contractAddr sdk.AccAddress, ibcPort string, msgs []wasmTypes.SubMsg, ogTx []byte, ogSigInfo wasmTypes.VerificationInfo) ([]byte, error)
}

// ContractResponseHandler default implementation that first dispatches submessage then normal messages.
// The Submessage execution may include an success/failure response handling by the contract that can overwrite the
// original
type ContractResponseHandler struct {
	md MsgDispatcher
}

func NewContractResponseHandler(md MsgDispatcher) *ContractResponseHandler {
	return &ContractResponseHandler{md: md}
}

// Handle processes the data returned by a contract invocation.
func (h ContractResponseHandler) Handle(ctx sdk.Context, contractAddr sdk.AccAddress, ibcPort string, messages []wasmTypes.SubMsg, origRspData []byte, ogTx []byte, ogSigInfo wasmTypes.VerificationInfo) ([]byte, error) {
	result := origRspData
	switch rsp, err := h.md.DispatchSubmessages(ctx, contractAddr, ibcPort, messages, ogTx, ogSigInfo); {
	case err != nil:
		return nil, sdkerrors.Wrap(err, "submessages")
	case rsp != nil:
		result = rsp
	}
	return result, nil
}

// reply is only called from keeper internal functions (dispatchSubmessages) after processing the submessage
func (k Keeper) reply(ctx sdk.Context, contractAddress sdk.AccAddress, reply wasmTypes.Reply, ogTx []byte, ogSigInfo wasmTypes.VerificationInfo, replyToContractHash []byte) ([]byte, error) {
	contractInfo, codeInfo, prefixStore, err := k.contractInstance(ctx, contractAddress)
	if err != nil {
		return nil, err
	}
	fmt.Printf("reply for %s \n", contractAddress.String())
	// always consider this pinned
	ctx.GasMeter().ConsumeGas(types.InstanceCost, "Loading Compute module: reply")

	store := ctx.KVStore(k.storeKey)
	contractKey := store.Get(types.GetContractEnclaveKey(contractAddress))

	env := types.NewEnv(ctx, contractAddress, sdk.Coins{}, contractAddress, contractKey)

	// prepare querier
	querier := QueryHandler{
		Ctx:     ctx,
		Plugins: k.queryPlugins,
	}

	// instantiate wasm contract
	gas := gasForContract(ctx)
	marshaledReply, error := json.Marshal(reply)
	//marshaledReply = append(replyToContractHash, marshaledReply...)
	marshaledReply = append(ogTx[0:64], marshaledReply...)

	if error != nil {
		return nil, error
	}

	res, gasUsed, execErr := k.wasmer.Execute(codeInfo.CodeHash, env, marshaledReply, prefixStore, cosmwasmAPI, querier, ctx.GasMeter(), gas, ogSigInfo, wasmTypes.HandleTypeReply)
	if execErr != nil {
		return nil, sdkerrors.Wrap(types.ErrReplyFailed, execErr.Error())
	}

	consumeGas(ctx, gasUsed)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeReply,
		sdk.NewAttribute(types.AttributeKeyContractAddr, contractAddress.String()),
	))

	data, err := k.handleContractResponse(ctx, contractAddress, contractInfo.IBCPortID, *res, res.Messages, res.Events, res.Data, ogTx, ogSigInfo)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrReplyFailed, err.Error())
	}

	return data, nil

}
