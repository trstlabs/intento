package compute

import (
	"fmt"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/trstlabs/trst/x/compute/internal/types"
)

// NewHandler returns a handler for "bank" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {

		case *MsgStoreCode:
			return handleStoreCode(ctx, k, msg)
		case *MsgInstantiateContract:
			return handleInstantiate(ctx, k, msg)
		case *MsgExecuteContract:
			return handleExecute(ctx, k, msg)
		case *types.MsgDiscardAutoMsg:
			return handleDiscardAutoMsg(ctx, k, msg)
			/*
				case MsgMigrateContract:
					return handleMigration(ctx, k, &msg)
				case MsgUpdateAdmin:
					return handleUpdateContractAdmin(ctx, k, &msg)
				case MsgClearAdmin:
					return handleClearContractAdmin(ctx, k, &msg)
			*/
		default:
			errMsg := fmt.Sprintf("unrecognized wasm message type: %T", msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// filteredMessageEvents returns the same events with all of type == EventTypeMessage removed.
// this is so only our top-level message event comes through
func filteredMessageEvents(manager *sdk.EventManager) []abci.Event {
	events := manager.ABCIEvents()
	res := make([]abci.Event, 0, len(events)+1)
	for _, e := range events {
		if e.Type != sdk.EventTypeMessage {
			res = append(res, e)
		}
	}
	return res
}

func handleStoreCode(ctx sdk.Context, k Keeper, msg *MsgStoreCode) (*sdk.Result, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	p := k.GetParams(ctx)
	duration, err := time.ParseDuration(msg.DefaultDuration)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}
	if duration > p.MaxContractDuration {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "contract duration must be shorter than maximum duration")
	}
	if duration != 0 && duration < p.MinContractDuration {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "contract duration must be longer than minimum duration")
	}
	interval, err := time.ParseDuration(msg.DefaultInterval)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}
	if interval != 0 && interval < p.MinContractInterval {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "contract interval must be longer than minimum interval")
	}

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	codeID, err := k.Create(ctx, sender, msg.WASMByteCode, msg.Source, msg.Builder, duration, interval, msg.Title, msg.Description)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
			sdk.NewAttribute(types.AttributeKeySigner, msg.Sender),
			sdk.NewAttribute(types.AttributeKeyCodeID, fmt.Sprintf("%d", codeID)),
		),
	})

	return &sdk.Result{
		Data:   []byte(fmt.Sprintf("%d", codeID)),
		Events: ctx.EventManager().ABCIEvents(),
	}, nil
}

func handleInstantiate(ctx sdk.Context, k Keeper, msg *MsgInstantiateContract) (*sdk.Result, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	var duration time.Duration = 0
	if msg.Duration != "" {
		duration, err = time.ParseDuration(msg.Duration)
		if err != nil {
			return nil, err
		}
	}
	var interval time.Duration = 0
	if msg.Interval != "" {
		interval, err = time.ParseDuration(msg.Interval)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
		}
	}
	var startTime time.Time = ctx.BlockHeader().Time
	if msg.StartDurationAt != 0 {
		startTime = time.Unix(int64(msg.StartDurationAt), 0)
		if err != nil {
			return nil, err
		}
	}
	p := k.GetParams(ctx)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}
	if interval != 0 && interval < p.MinContractInterval {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "contract interval must be longer than minimum interval")
	}
	if duration != 0 {
		if duration > p.MaxContractDuration {
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "contract duration must be shorter than maximum duration")
		}
		if duration < p.MinContractDuration {
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "contract duration must be longer than minimum duration")
		}
		if startTime.After(ctx.BlockHeader().Time.Add(duration)) {
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "start time must be before contract end time")
		}

	}

	owner, _ := sdk.AccAddressFromBech32(msg.Owner)
	contractAddr, data, err := k.Instantiate(ctx, msg.CodeID, sender, msg.Msg, msg.AutoMsg, msg.ContractId, msg.Funds, msg.CallbackSig, duration, interval, startTime, owner)
	if err != nil {

		return nil, err
	}

	events := filteredMessageEvents(ctx.EventManager())
	custom := sdk.Events{sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
		sdk.NewAttribute(types.AttributeKeySigner, msg.Sender),
		sdk.NewAttribute(types.AttributeKeyCodeID, fmt.Sprintf("%d", msg.CodeID)),
		sdk.NewAttribute(types.AttributeKeyContract, contractAddr.String()),
	)}
	events = append(events, custom.ToABCIEvents()...)

	return &sdk.Result{
		Data:   data,
		Events: events,
	}, nil
}
func handleExecute(ctx sdk.Context, k Keeper, msg *MsgExecuteContract) (*sdk.Result, error) {

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	contract, err := sdk.AccAddressFromBech32(msg.Contract)
	if err != nil {
		return nil, err
	}
	info := k.GetContractInfo(ctx, contract)
	if info == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "contract code does not exist")
	}
	/*	if info.CodeID < 2 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "cannot execute on internal contract code")
	}*/
	res, err := k.Execute(
		ctx,
		contract,
		sender,
		msg.Msg,
		msg.Funds,
		msg.CallbackSig,
	)

	if err != nil {
		return nil, err
	}

	events := filteredMessageEvents(ctx.EventManager())

	custom := sdk.Events{sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
		sdk.NewAttribute(types.AttributeKeySigner, msg.Sender),
		sdk.NewAttribute(types.AttributeKeyContract, msg.Contract),
	)}

	events = append(events, custom.ToABCIEvents()...)

	res.Events = events

	return res, nil
}

func handleDiscardAutoMsg(ctx sdk.Context, k Keeper, msg *types.MsgDiscardAutoMsg) (*sdk.Result, error) {

	info := k.GetContractInfo(ctx, msg.ContractAddress)
	if info == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "contract code does not exist")
	}

	err := k.DiscardAutoMsg(
		ctx,
		*info,
		msg.ContractAddress,
		msg.Sender,
	)

	if err != nil {
		return nil, err
	}

	events := filteredMessageEvents(ctx.EventManager())

	custom := sdk.Events{sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
		sdk.NewAttribute(types.AttributeKeySigner, msg.Sender.String()),
		sdk.NewAttribute(types.AttributeKeyContract, msg.ContractAddress.String()),
		sdk.NewAttribute("cancelled", msg.ContractAddress.String()),
	)}

	events = append(events, custom.ToABCIEvents()...)

	return &sdk.Result{
		Data:   msg.ContractAddress,
		Events: events,
	}, nil
}

/*
func handleDeleteContract(ctx sdk.Context, k Keeper, msg *MsgDeleteContract) (*sdk.Result, error) {
	res, err := k.DeleteContract(ctx, msg.Contract, msg.Sender, msg.CodeID, msg.DeleteContractMsg) // for MsgMigrateContract, there is only one signer which is msg.Sender (https://github.com/trstlabs/trst/blob/d7813792fa07b93a10f0885eaa4c5e0a0a698854/x/compute/internal/types/msg.go#L228-L230)
	if err != nil {
		return nil, err
	}

	events := filteredMessageEvents(ctx.EventManager())
	ourEvent := sdk.Events{sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
		sdk.NewAttribute(types.AttributeKeySigner, msg.Sender),
		sdk.NewAttribute(types.AttributeKeyContract, msg.Contract.String()),
	)}
	res.Events = append(events, ourEvent.ToABCIEvents()...)
	return res, nil
}
*/
/*
func handleMigration(ctx sdk.Context, k Keeper, msg *MsgMigrateContract) (*sdk.Result, error) {
	res, err := k.Migrate(ctx, msg.Contract, msg.Sender, msg.CodeID, msg.MigrateMsg) // for MsgMigrateContract, there is only one signer which is msg.Sender (https://github.com/trstlabs/trst/blob/d7813792fa07b93a10f0885eaa4c5e0a0a698854/x/compute/internal/types/msg.go#L228-L230)
	if err != nil {
		return nil, err
	}

	events := filteredMessageEvents(ctx.EventManager())
	ourEvent := sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
		sdk.NewAttribute(types.AttributeKeySigner, msg.Sender),
		sdk.NewAttribute(types.AttributeKeyContract, msg.Contract.String()),
	)
	res.Events = append(events, ourEvent)
	return res, nil
}

func handleUpdateContractAdmin(ctx sdk.Context, k Keeper, msg *MsgUpdateAdmin) (*sdk.Result, error) {
	if err := k.UpdateContractAdmin(ctx, msg.Contract, msg.Sender, msg.NewAdmin); err != nil {
		return nil, err
	}
	events := ctx.EventManager().Events()
	ourEvent := sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
		sdk.NewAttribute(types.AttributeKeySigner, msg.Sender),
		sdk.NewAttribute(types.AttributeKeyContract, msg.Contract.String()),
	)
	return &sdk.Result{
		Events: append(events, ourEvent),
	}, nil
}

func handleClearContractAdmin(ctx sdk.Context, k Keeper, msg *MsgClearAdmin) (*sdk.Result, error) {
	if err := k.ClearContractAdmin(ctx, msg.Contract, msg.Sender); err != nil {
		return nil, err
	}
	events := ctx.EventManager().Events()
	ourEvent := sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
		sdk.NewAttribute(types.AttributeKeySigner, msg.Sender),
		sdk.NewAttribute(types.AttributeKeyContract, msg.Contract.String()),
	)
	return &sdk.Result{
		Events: append(events, ourEvent),
	}, nil
}
*/
