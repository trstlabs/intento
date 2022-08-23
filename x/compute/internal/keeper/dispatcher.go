package keeper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
	wasmTypes "github.com/trstlabs/trst/go-cosmwasm/types"
	"github.com/trstlabs/trst/x/compute/internal/types"
)

// Messenger is an extension point for custom wasmd message handling

type Messenger interface {
	// DispatchMsg encodes the wasmVM message and dispatches it.
	DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmTypes.CosmosMsg) (events []sdk.Event, data [][]byte, err error)
}

// Replyer is a subset of keeper that can handle replies to submessages
type Replyer interface {
	reply(ctx sdk.Context, contractAddress sdk.AccAddress, reply wasmTypes.Reply, ogTx []byte, ogSigInfo wasmTypes.VerificationInfo) ([]byte, error)
}

// MessageDispatcher coordinates message sending and submessage reply/ state commits
type MessageDispatcher struct {
	messenger Messenger
	keeper    Replyer
}

// NewMessageDispatcher constructor
func NewMessageDispatcher(messenger Messenger, keeper Replyer) *MessageDispatcher {
	return &MessageDispatcher{messenger: messenger, keeper: keeper}
}

func filterEvents(events []sdk.Event) []sdk.Event {
	// pre-allocate space for efficiency
	res := make([]sdk.Event, 0, len(events))
	for _, ev := range events {
		if ev.Type != "message" {
			res = append(res, ev)
		}
	}
	return res
}

func sdkAttributesToWasmVMAttributes(attrs []abci.EventAttribute) []wasmTypes.Attribute {
	res := make([]wasmTypes.Attribute, len(attrs))
	for i, attr := range attrs {
		res[i] = wasmTypes.Attribute{
			Key:       string(attr.Key),
			Value:     attr.Value,
			Encrypted: false,
			PubDb:     false,
			AccAddr:   "",
		}
	}
	return res
}

func sdkEventsToWasmVMEvents(events []sdk.Event) []wasmTypes.Event {
	res := make([]wasmTypes.Event, len(events))
	for i, ev := range events {
		res[i] = wasmTypes.Event{
			Type:       ev.Type,
			Attributes: sdkAttributesToWasmVMAttributes(ev.Attributes),
		}
	}
	return res
}

// dispatchMsgWithGasLimit sends a message with gas limit applied
func (d MessageDispatcher) dispatchMsgWithGasLimit(ctx sdk.Context, contractAddr sdk.AccAddress, ibcPort string, msg wasmTypes.CosmosMsg, gasLimit uint64) (events []sdk.Event, data [][]byte, err error) {
	limitedMeter := sdk.NewGasMeter(gasLimit)
	subCtx := ctx.WithGasMeter(limitedMeter)
	fmt.Printf("Dispatch for %s \n", contractAddr.String())
	// catch out of gas panic and just charge the entire gas limit
	defer func() {
		if r := recover(); r != nil {
			// if it's not an OutOfGas error, raise it again
			if _, ok := r.(sdk.ErrorOutOfGas); !ok {
				// log it to get the original stack trace somewhere (as panic(r) keeps message but stacktrace to here
				ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName)).Info("SubMsg rethrowing panic: %#v", r)
				panic(r)
			}
			ctx.GasMeter().ConsumeGas(gasLimit, "SubMsg OutOfGas panic")
			err = sdkerrors.Wrap(sdkerrors.ErrOutOfGas, "SubMsg hit gas limit")
		}
	}()
	events, data, err = d.messenger.DispatchMsg(subCtx, contractAddr, ibcPort, msg)

	// make sure we charge the parent what was spent
	spent := subCtx.GasMeter().GasConsumed()
	ctx.GasMeter().ConsumeGas(spent, "From gas-limited SubMsg")

	return events, data, err
}

type InvalidRequest struct {
	Err     string `json:"error"`
	Request []byte `json:"request"`
}

func (e InvalidRequest) Error() string {
	return fmt.Sprintf("invalid request: %s - original request: %s", e.Err, string(e.Request))
}

type InvalidResponse struct {
	Err      string `json:"error"`
	Response []byte `json:"response"`
}

func (e InvalidResponse) Error() string {
	return fmt.Sprintf("invalid response: %s - original response: %s", e.Err, string(e.Response))
}

type NoSuchContract struct {
	Addr string `json:"addr,omitempty"`
}

func (e NoSuchContract) Error() string {
	return fmt.Sprintf("no such contract: %s", e.Addr)
}

type Unknown struct{}

func (e Unknown) Error() string {
	return "unknown system error"
}

type UnsupportedRequest struct {
	Kind string `json:"kind,omitempty"`
}

func (e UnsupportedRequest) Error() string {
	return fmt.Sprintf("unsupported request: %s", e.Kind)
}

// Reply is encrypted only when it is a contract reply
func isReplyEncrypted(msg wasmTypes.CosmosMsg, reply wasmTypes.Reply) bool {
	return (msg.Wasm != nil)
}

// Issue #759 - we don't return error string for worries of non-determinism
func redactError(err error) error {
	// Do not redact encrypted wasm contract errors
	if strings.HasPrefix(err.Error(), "encrypted:") {
		// remove encrypted sign
		e := strings.ReplaceAll(err.Error(), "encrypted: ", "")
		e = strings.ReplaceAll(e, ": execute contract failed", "")
		e = strings.ReplaceAll(e, ": instantiate contract failed", "")
		return fmt.Errorf("%s", e)
	}

	// Do not redact system errors
	// SystemErrors must be created in x/wasm and we can ensure determinism
	if wasmTypes.ToSystemError(err) != nil {
		return err
	}

	// FIXME: do we want to hardcode some constant string mappings here as well?
	// Or better document them? (SDK error string may change on a patch release to fix wording)
	// sdk/11 is out of gas
	// sdk/5 is insufficient funds (on bank send)
	// (we can theoretically redact less in the future, but this is a first step to safety)
	codespace, code, _ := sdkerrors.ABCIInfo(err, false)
	return fmt.Errorf("codespace: %s, code: %d", codespace, code)
}

// DispatchSubmessages builds a sandbox to execute these messages and returns the execution result to the contract
// that dispatched them, both on success as well as failure
func (d MessageDispatcher) DispatchSubmessages(ctx sdk.Context, contractAddr sdk.AccAddress, ibcPort string, msgs []wasmTypes.SubMsg, ogTx []byte, ogSigInfo wasmTypes.VerificationInfo) ([]byte, error) {
	var rsp []byte
	for _, msg := range msgs {
		// Check replyOn validity
		switch msg.ReplyOn {
		case wasmTypes.ReplySuccess, wasmTypes.ReplyError, wasmTypes.ReplyAlways, wasmTypes.ReplyNever:
		default:
			return nil, sdkerrors.Wrap(types.ErrInvalid, "replyOn value")
		}
		//fmt.Printf("SubMsg for %s \n", contractAddr.String())
		fmt.Printf("SubMsg %+v\n", msg)
		// first, we build a sub-context which we can use inside the submessages
		subCtx, commit := ctx.CacheContext()
		em := sdk.NewEventManager()
		subCtx = subCtx.WithEventManager(em)

		// check how much gas left locally, optionally wrap the gas meter
		gasRemaining := ctx.GasMeter().Limit() - ctx.GasMeter().GasConsumed()
		limitGas := msg.GasLimit != nil && (*msg.GasLimit < gasRemaining)

		var err error
		var events []sdk.Event
		var data [][]byte
		if limitGas {
			fmt.Printf("Dispatch msg with limit gas for %s \n", contractAddr.String())
			events, data, err = d.dispatchMsgWithGasLimit(subCtx, contractAddr, ibcPort, msg.Msg, *msg.GasLimit)
		} else {
			fmt.Printf("Dispatch msg with no limit gas for %s \n", contractAddr.String())
			events, data, err = d.messenger.DispatchMsg(subCtx, contractAddr, ibcPort, msg.Msg)
		}
		//ctx.EventManager().EmitEvents(events)

		// if it succeeds, commit state changes from submessage, and pass on events to Event Manager
		var filteredEvents []sdk.Event
		if err == nil {
			commit()
			filteredEvents = filterEvents(append(em.Events(), events...))
			ctx.EventManager().EmitEvents(filteredEvents)

			if msg.Msg.Wasm == nil {
				filteredEvents = []sdk.Event{}
			} else {
				for _, e := range filteredEvents {
					attributes := e.Attributes
					sort.SliceStable(attributes, func(i, j int) bool {
						return bytes.Compare(attributes[i].Key, attributes[j].Key) < 0
					})
				}
			}
		} // on failure, revert state from sandbox, and ignore events (just skip doing the above)

		// we only callback if requested. Short-circuit here the cases we don't want to
		if (msg.ReplyOn == wasmTypes.ReplySuccess || msg.ReplyOn == wasmTypes.ReplyNever) && err != nil {
			return nil, err
		}

		if msg.ReplyOn == wasmTypes.ReplyNever || (msg.ReplyOn == wasmTypes.ReplyError && err == nil) {
			continue
		}

		// If we are here it means that ReplySuccess and success OR ReplyError and there were errors OR ReplyAlways.
		// Basically, handle replying to the contract
		// We need to create a SubMsgResult and pass it into the calling contract
		var result wasmTypes.SubMsgResult
		if err == nil {
			//fmt.Printf("Reply data0 %v \n", data[0])
			// just take the first one for now if there are multiple sub-sdk messages
			// and safely return nothing if no data
			var responseData []byte
			if len(data) > 0 {
				responseData = data[0]
			}

			result = wasmTypes.SubMsgResult{
				Ok: &wasmTypes.SubMsgResponse{
					Events: sdkEventsToWasmVMEvents(filteredEvents), //wasmTypes.Events{}, //
					Data:   responseData,
				},
			}
		} else {
			// Issue #759 - we don't return error string for worries of non-determinism
			ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName)).Info("Redacting submessage error", "cause", err)
			result = wasmTypes.SubMsgResult{
				Err: redactError(err).Error(),
			}
		}

		msg_id := []byte(fmt.Sprint(msg.ID))
		// now handle the reply, we use the parent context, and abort on error
		reply := wasmTypes.Reply{
			ID:     msg_id,
			Result: result,
		}
		// we can ignore any result returned as there is nothing to do with the data
		// and the events are already in the ctx.EventManager()

		// In order to specify that the reply isn't signed by the enclave we use "SIGN_MODE_UNSPECIFIED"
		// The SGX will notice that the value is SIGN_MODE_UNSPECIFIED and will treat the message as plaintext.
		replySigInfo := wasmTypes.VerificationInfo{
			Bytes:     []byte{},
			ModeInfo:  []byte{},
			PublicKey: []byte{},
			Signature: []byte{},
			SignMode:  "SIGN_MODE_UNSPECIFIED",
		}

		if isReplyEncrypted(msg.Msg, reply) {
			var dataWithInternalReplyInfo wasmTypes.DataWithInternalReplyInfo

			if reply.Result.Ok != nil {
				//fmt.Printf("Reply data raw %v \n", reply.Result.Ok.Data)
				err = json.Unmarshal(reply.Result.Ok.Data, &dataWithInternalReplyInfo)
				if err != nil {
					return nil, fmt.Errorf("cannot serialize DataWithInternalReplyInfo into json : %w", err)
				}

				reply.Result.Ok.Data = dataWithInternalReplyInfo.Data
				//fmt.Printf("Reply Ok %v \n", reply.Result.Ok)

			} else {
				err = json.Unmarshal(data[0], &dataWithInternalReplyInfo)
				if err != nil {
					return nil, fmt.Errorf("cannot serialize DataWithInternalReplyInfo into json : %w", err)
				}
			}

			if len(dataWithInternalReplyInfo.InternalMsgId) == 0 || len(dataWithInternalReplyInfo.InternalReplyEnclaveSig) == 0 {
				return nil, fmt.Errorf("when sending a reply both InternalReplyEnclaveSig and InternalMsgId are expected to be initialized")
			}
			replySigInfo = ogSigInfo
			reply.ID = dataWithInternalReplyInfo.InternalMsgId
			replySigInfo.CallbackSignature = dataWithInternalReplyInfo.InternalReplyEnclaveSig

		}
		rspData, err := d.keeper.reply(ctx, contractAddr, reply, ogTx, replySigInfo)

		switch {
		case err != nil:
			fmt.Printf("Got err for SubMsg for %v \n", err.Error())
			return nil, err
		case rspData != nil:
			rsp = rspData
		}
	}

	return rsp, nil
}
