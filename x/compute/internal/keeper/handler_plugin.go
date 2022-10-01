package keeper

import (
	"encoding/json"
	"errors"
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	channelkeeper "github.com/cosmos/ibc-go/v3/modules/core/04-channel/keeper"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"

	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	wasmTypes "github.com/trstlabs/trst/go-cosmwasm/types"
	"github.com/trstlabs/trst/x/compute/internal/types"
)

// MessageHandlerChain defines a chain of handlers that are called one by one until it can be handled.
type MessageHandlerChain struct {
	handlers []Messenger
}

type MessageHandler struct {
	// router is used to route StargateMsg and any other msg except for MsgExecuteContract & MsgInstantiateContrat.
	router MsgServiceRouter
	// legacyRouter is used to route MsgExecuteContract & MsgInstantiateContrat.
	// the reason is those msgs use the data field internally for reply, which is
	// truncated if the msg erred
	legacyRouter sdk.Router
	encoders     MessageEncoders
}

func NewSDKMessageHandler(router MsgServiceRouter, legacyRouter sdk.Router, encoders MessageEncoders) MessageHandler {
	return MessageHandler{
		router:       router,
		legacyRouter: legacyRouter,
		encoders:     encoders,
	}
}

// IBCRawPacketHandler handels IBC.SendPacket messages which are published to an IBC channel.
type IBCRawPacketHandler struct {
	channelKeeper    channelkeeper.Keeper
	capabilityKeeper capabilitykeeper.ScopedKeeper
}

func NewIBCRawPacketHandler(chk channelkeeper.Keeper, cak capabilitykeeper.ScopedKeeper) IBCRawPacketHandler {
	return IBCRawPacketHandler{channelKeeper: chk, capabilityKeeper: cak}
}

func NewMessageHandlerChain(first Messenger, others ...Messenger) *MessageHandlerChain {
	r := &MessageHandlerChain{handlers: append([]Messenger{first}, others...)}
	for i := range r.handlers {
		if r.handlers[i] == nil {
			panic(fmt.Sprintf("handler must not be nil at position : %d", i))
		}
	}
	return r
}

func NewMessageHandler(
	msgRouter MsgServiceRouter,
	router sdk.Router,
	customEncoders *MessageEncoders,
	channelKeeper channelkeeper.Keeper,
	capabilityKeeper capabilitykeeper.ScopedKeeper,
	portSource types.ICS20TransferPortSource,
	unpacker codectypes.AnyUnpacker) Messenger {
	encoders := DefaultEncoders(portSource, unpacker).Merge(customEncoders)
	return NewMessageHandlerChain(
		NewSDKMessageHandler(msgRouter, router, encoders),
		NewIBCRawPacketHandler(channelKeeper, capabilityKeeper),
	)
}

// DispatchMsg dispatch message and calls chained handlers one after another in
// order to find the right one to process given message. If a handler cannot
// process given message (returns ErrUnknownMsg), its result is ignored and the
// next handler is executed.
func (m MessageHandlerChain) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmTypes.CosmosMsg) ([]sdk.Event, [][]byte, error) {
	for _, h := range m.handlers {
		events, data, err := h.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
		switch {
		case err == nil:
			return events, data, nil
		case errors.Is(err, types.ErrUnknownMsg):
			continue
		default:
			return events, data, err
		}
	}
	return nil, nil, sdkerrors.Wrap(types.ErrUnknownMsg, "no handler found")
}

// DispatchMsg publishes a raw IBC packet onto the channel.
func (h IBCRawPacketHandler) DispatchMsg(ctx sdk.Context, _ sdk.AccAddress, contractIBCPortID string, msg wasmTypes.CosmosMsg) (events []sdk.Event, data [][]byte, err error) {
	if msg.IBC == nil || msg.IBC.SendPacket == nil {
		return nil, nil, types.ErrUnknownMsg
	}
	if contractIBCPortID == "" {
		return nil, nil, sdkerrors.Wrapf(types.ErrUnsupportedForContract, "ibc not supported")
	}
	contractIBCChannelID := msg.IBC.SendPacket.ChannelID
	if contractIBCChannelID == "" {
		return nil, nil, sdkerrors.Wrapf(types.ErrEmpty, "ibc channel")
	}

	sequence, found := h.channelKeeper.GetNextSequenceSend(ctx, contractIBCPortID, contractIBCChannelID)
	if !found {
		return nil, nil, sdkerrors.Wrapf(channeltypes.ErrSequenceSendNotFound,
			"source port: %s, source channel: %s", contractIBCPortID, contractIBCChannelID,
		)
	}

	channelInfo, ok := h.channelKeeper.GetChannel(ctx, contractIBCPortID, contractIBCChannelID)
	if !ok {
		return nil, nil, sdkerrors.Wrap(channeltypes.ErrInvalidChannel, "not found")
	}
	channelCap, ok := h.capabilityKeeper.GetCapability(ctx, host.ChannelCapabilityPath(contractIBCPortID, contractIBCChannelID))
	if !ok {
		return nil, nil, sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}
	packet := channeltypes.NewPacket(
		msg.IBC.SendPacket.Data,
		sequence,
		contractIBCPortID,
		contractIBCChannelID,
		channelInfo.Counterparty.PortId,
		channelInfo.Counterparty.ChannelId,
		convertWasmIBCTimeoutHeightToCosmosHeight(msg.IBC.SendPacket.Timeout.Block),
		msg.IBC.SendPacket.Timeout.Timestamp,
	)
	return nil, nil, h.channelKeeper.SendPacket(ctx, channelCap, packet)
}

type BankEncoder func(sender sdk.AccAddress, msg *wasmTypes.BankMsg) ([]sdk.Msg, error)
type CustomEncoder func(sender sdk.AccAddress, msg json.RawMessage) ([]sdk.Msg, error)
type DistributionEncoder func(sender sdk.AccAddress, msg *wasmTypes.DistributionMsg) ([]sdk.Msg, error)
type GovEncoder func(sender sdk.AccAddress, msg *wasmTypes.GovMsg) ([]sdk.Msg, error)
type IBCEncoder func(ctx sdk.Context, sender sdk.AccAddress, contractIBCPortID string, msg *wasmTypes.IBCMsg) ([]sdk.Msg, error)
type StakingEncoder func(sender sdk.AccAddress, msg *wasmTypes.StakingMsg) ([]sdk.Msg, error)
type StargateEncoder func(sender sdk.AccAddress, msg *wasmTypes.StargateMsg) ([]sdk.Msg, error)
type WasmEncoder func(sender sdk.AccAddress, msg *wasmTypes.WasmMsg) ([]sdk.Msg, error)

type MessageEncoders struct {
	Bank         BankEncoder
	Custom       CustomEncoder
	Distribution DistributionEncoder
	Gov          GovEncoder
	IBC          IBCEncoder
	Staking      StakingEncoder
	Stargate     StargateEncoder
	Wasm         WasmEncoder
}

func DefaultEncoders(portSource types.ICS20TransferPortSource, unpacker codectypes.AnyUnpacker) MessageEncoders {
	return MessageEncoders{
		Bank:         EncodeBankMsg,
		Custom:       NoCustomMsg,
		Distribution: EncodeDistributionMsg,
		Gov:          EncodeGovMsg,
		IBC:          EncodeIBCMsg(portSource),
		Staking:      EncodeStakingMsg,
		Stargate:     EncodeStargateMsg(unpacker),
		Wasm:         EncodeWasmMsg,
	}
}

func (e MessageEncoders) Merge(o *MessageEncoders) MessageEncoders {
	if o == nil {
		return e
	}
	if o.Bank != nil {
		e.Bank = o.Bank
	}
	if o.Custom != nil {
		e.Custom = o.Custom
	}
	if o.Staking != nil {
		e.Staking = o.Staking
	}
	if o.Wasm != nil {
		e.Wasm = o.Wasm
	}
	if o.Gov != nil {
		e.Gov = o.Gov
	}
	return e
}

func (e MessageEncoders) Encode(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmTypes.CosmosMsg) ([]sdk.Msg, error) {
	switch {
	case msg.Bank != nil:
		return e.Bank(contractAddr, msg.Bank)
	case msg.Custom != nil:
		return e.Custom(contractAddr, msg.Custom)
	case msg.Distribution != nil:
		return e.Distribution(contractAddr, msg.Distribution)
	case msg.Gov != nil:
		return e.Gov(contractAddr, msg.Gov)
	case msg.IBC != nil:
		return e.IBC(ctx, contractAddr, contractIBCPortID, msg.IBC)
	case msg.Staking != nil:
		return e.Staking(contractAddr, msg.Staking)
	case msg.Stargate != nil:
		return e.Stargate(contractAddr, msg.Stargate)
	case msg.Wasm != nil:
		return e.Wasm(contractAddr, msg.Wasm)
	}

	return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown variant of Wasm")
}

var VoteOptionMap = map[wasmTypes.VoteOption]string{
	wasmTypes.Yes:        "VOTE_OPTION_YES",
	wasmTypes.Abstain:    "VOTE_OPTION_ABSTAIN",
	wasmTypes.No:         "VOTE_OPTION_NO",
	wasmTypes.NoWithVeto: "VOTE_OPTION_NO_WITH_VETO",
}

func EncodeGovMsg(sender sdk.AccAddress, msg *wasmTypes.GovMsg) ([]sdk.Msg, error) {
	if msg.Vote == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown variant of Gov")
	}

	opt, exists := VoteOptionMap[msg.Vote.Vote]
	if !exists {
		// if it's not found, let the `VoteOptionFromString` below fail
		opt = ""
	}

	option, err := govtypes.VoteOptionFromString(opt)
	if err != nil {
		return nil, err
	}

	sdkMsg := govtypes.NewMsgVote(sender, msg.Vote.ProposalId, option)
	return []sdk.Msg{sdkMsg}, nil
}

func EncodeIBCMsg(portSource types.ICS20TransferPortSource) func(ctx sdk.Context, sender sdk.AccAddress, contractIBCPortID string, msg *wasmTypes.IBCMsg) ([]sdk.Msg, error) {
	return func(ctx sdk.Context, sender sdk.AccAddress, contractIBCPortID string, msg *wasmTypes.IBCMsg) ([]sdk.Msg, error) {
		switch {
		case msg.CloseChannel != nil:
			return []sdk.Msg{&channeltypes.MsgChannelCloseInit{
				PortId:    PortIDForContract(sender),
				ChannelId: msg.CloseChannel.ChannelID,
				Signer:    sender.String(),
			}}, nil
		case msg.Transfer != nil:
			//amount, err := convertWasmCoinToSdkCoin(msg.Transfer.Amount)
			//if err != nil {
			//	return nil, sdkerrors.Wrap(err, "amount")
			//}
			msg := &ibctransfertypes.MsgTransfer{
				SourcePort:       portSource.GetPort(ctx),
				SourceChannel:    msg.Transfer.ChannelID,
				Token:            msg.Transfer.Amount,
				Sender:           sender.String(),
				Receiver:         msg.Transfer.ToAddress,
				TimeoutHeight:    convertWasmIBCTimeoutHeightToCosmosHeight(msg.Transfer.Timeout.Block),
				TimeoutTimestamp: msg.Transfer.Timeout.Timestamp,
			}
			return []sdk.Msg{msg}, nil
		default:
			return nil, sdkerrors.Wrap(types.ErrUnknownMsg, "Unknown variant of IBC")
		}
	}
}

func EncodeBankMsg(sender sdk.AccAddress, msg *wasmTypes.BankMsg) ([]sdk.Msg, error) {
	if msg.Send == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown variant of Bank")
	}
	if len(msg.Send.Amount) == 0 {
		return nil, nil
	}
	// validate that the addresses are valid
	_, stderr := sdk.AccAddressFromBech32(msg.Send.ToAddress)
	if stderr != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Send.ToAddress)
	}

	//toSend, err := convertWasmCoinsToSdkCoins(msg.Send.Amount)
	//if err != nil {
	//	return nil, err
	//}
	sdkMsg := banktypes.MsgSend{
		FromAddress: sender.String(),
		ToAddress:   msg.Send.ToAddress,
		Amount:      msg.Send.Amount,
	}
	return []sdk.Msg{&sdkMsg}, nil
}

func NoCustomMsg(sender sdk.AccAddress, msg json.RawMessage) ([]sdk.Msg, error) {
	return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Custom variant not supported")
}

func EncodeDistributionMsg(sender sdk.AccAddress, msg *wasmTypes.DistributionMsg) ([]sdk.Msg, error) {
	switch {
	case msg.SetWithdrawAddress != nil:
		setMsg := distrtypes.MsgSetWithdrawAddress{
			DelegatorAddress: sender.String(),
			WithdrawAddress:  msg.SetWithdrawAddress.Address,
		}
		return []sdk.Msg{&setMsg}, nil
	case msg.WithdrawDelegatorReward != nil:
		withdrawMsg := distrtypes.MsgWithdrawDelegatorReward{
			DelegatorAddress: sender.String(),
			ValidatorAddress: msg.WithdrawDelegatorReward.Validator,
		}
		return []sdk.Msg{&withdrawMsg}, nil
	default:
		return nil, sdkerrors.Wrap(types.ErrUnknownMsg, "unknown variant of Distribution")
	}
}

func EncodeStakingMsg(sender sdk.AccAddress, msg *wasmTypes.StakingMsg) ([]sdk.Msg, error) {
	var err error
	switch {
	case msg.Delegate != nil:
		// Check that the address belongs to a validator.
		validator, err := sdk.ValAddressFromBech32(msg.Delegate.Validator)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Delegate.Validator)
		}
		coin, err := convertWasmCoinToSdkCoin(msg.Delegate.Amount)
		if err != nil {
			return nil, err
		}
		//sdkMsg := stakingtypes.MsgDelegate{
		//	DelegatorAddress: sender.String(),
		//	ValidatorAddress: msg.Delegate.Validator,
		//	Amount:           coin,
		//}
		sdkMsg := stakingtypes.NewMsgDelegate(sender, validator, coin)
		return []sdk.Msg{sdkMsg}, nil

	case msg.Redelegate != nil:
		// Check that the addresses belong to validators.
		_, err = sdk.ValAddressFromBech32(msg.Redelegate.SrcValidator)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Redelegate.SrcValidator)
		}
		_, err = sdk.ValAddressFromBech32(msg.Redelegate.DstValidator)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Redelegate.DstValidator)
		}
		coin, err := convertWasmCoinToSdkCoin(msg.Redelegate.Amount)
		if err != nil {
			return nil, err
		}
		sdkMsg := stakingtypes.MsgBeginRedelegate{
			DelegatorAddress:    sender.String(),
			ValidatorSrcAddress: msg.Redelegate.SrcValidator,
			ValidatorDstAddress: msg.Redelegate.DstValidator,
			Amount:              coin,
		}
		return []sdk.Msg{&sdkMsg}, nil
	case msg.Undelegate != nil:
		// Check that the address belongs to a validator.
		_, err = sdk.ValAddressFromBech32(msg.Undelegate.Validator)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Undelegate.Validator)
		}
		coin, err := convertWasmCoinToSdkCoin(msg.Undelegate.Amount)
		if err != nil {
			return nil, err
		}
		sdkMsg := stakingtypes.MsgUndelegate{
			DelegatorAddress: sender.String(),
			ValidatorAddress: msg.Undelegate.Validator,
			Amount:           coin,
		}
		return []sdk.Msg{&sdkMsg}, nil
	case msg.Withdraw != nil:
		senderAddr := sender.String()
		rcpt := senderAddr
		if len(msg.Withdraw.Recipient) != 0 {
			// Check that the address belongs to a real account.
			_, err = sdk.AccAddressFromBech32(msg.Withdraw.Recipient)
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Withdraw.Recipient)
			}
			rcpt = msg.Withdraw.Recipient
		}
		// Check that the address belongs to a validator.
		_, err = sdk.ValAddressFromBech32(msg.Withdraw.Validator)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Withdraw.Validator)
		}
		setMsg := distrtypes.MsgSetWithdrawAddress{
			DelegatorAddress: senderAddr,
			WithdrawAddress:  rcpt,
		}
		withdrawMsg := distrtypes.MsgWithdrawDelegatorReward{
			DelegatorAddress: senderAddr,
			ValidatorAddress: msg.Withdraw.Validator,
		}
		return []sdk.Msg{&setMsg, &withdrawMsg}, nil
	default:
		return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown variant of Staking")
	}
}

func EncodeStargateMsg(unpacker codectypes.AnyUnpacker) StargateEncoder {
	return func(sender sdk.AccAddress, msg *wasmTypes.StargateMsg) ([]sdk.Msg, error) {
		any := codectypes.Any{
			TypeUrl: msg.TypeURL,
			Value:   msg.Value,
		}
		var sdkMsg sdk.Msg
		if err := unpacker.UnpackAny(&any, &sdkMsg); err != nil {
			return nil, sdkerrors.Wrap(types.ErrInvalidMsg, fmt.Sprintf("Cannot unpack proto message with type URL: %s", msg.TypeURL))
		}
		if err := codectypes.UnpackInterfaces(sdkMsg, unpacker); err != nil {
			return nil, sdkerrors.Wrap(types.ErrInvalidMsg, fmt.Sprintf("UnpackInterfaces inside msg: %s", err))
		}
		return []sdk.Msg{sdkMsg}, nil
	}
}

func EncodeWasmMsg(sender sdk.AccAddress, msg *wasmTypes.WasmMsg) ([]sdk.Msg, error) {
	switch {
	case msg.Execute != nil:
		contractAddr, err := sdk.AccAddressFromBech32(msg.Execute.ContractAddr)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Execute.ContractAddr)
		}
		coins, err := convertWasmCoinsToSdkCoins(msg.Execute.Funds)
		if err != nil {
			return nil, err
		}

		sdkMsg := types.MsgExecuteContract{
			Sender:      sender.String(),
			Contract:    contractAddr.String(),
			CodeHash:    msg.Execute.CodeHash,
			Msg:         msg.Execute.Msg,
			Funds:       coins,
			CallbackSig: msg.Execute.CallbackSignature,
		}
		return []sdk.Msg{&sdkMsg}, nil
	case msg.Instantiate != nil:
		coins, err := convertWasmCoinsToSdkCoins(msg.Instantiate.Funds)
		if err != nil {
			return nil, err
		}

		sdkMsg := types.MsgInstantiateContract{
			Sender:          sender.String(),
			CodeID:          msg.Instantiate.CodeID,
			ContractId:      msg.Instantiate.ContractID,
			CodeHash:        msg.Instantiate.CodeHash,
			Msg:             msg.Instantiate.Msg,
			AutoMsg:         msg.Instantiate.AutoMsg,
			Duration:        msg.Instantiate.Duration,
			Interval:        msg.Instantiate.Interval,
			StartDurationAt: msg.Instantiate.StartDurationAt,
			Funds:           coins,
			CallbackSig:     msg.Instantiate.CallbackSignature,
		}
		return []sdk.Msg{&sdkMsg}, nil
	default:
		return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown variant of Wasm")
	}
}

func (h MessageHandler) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmTypes.CosmosMsg) ([]sdk.Event, [][]byte, error) {
	sdkMsgs, err := h.encoders.Encode(ctx, contractAddr, contractIBCPortID, msg)
	if err != nil {
		return nil, nil, err
	}

	var (
		events []sdk.Event
		data   [][]byte
	)
	for _, sdkMsg := range sdkMsgs {
		sdkEvents, sdkData, err := h.handleSdkMessage(ctx, contractAddr, sdkMsg)

		if err != nil {
			data = append(data, sdkData)
			return nil, data, err
		}
		// append data
		data = append(data, sdkData)
		events = append(events, sdkEvents...)
	}

	return events, data, nil

	//return nil, nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown variant of CosmosMsgVersion")
}

func (h MessageHandler) handleSdkMessage(ctx sdk.Context, contractAddr sdk.Address, msg sdk.Msg) (sdk.Events, []byte, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, nil, err
	}

	// make sure this account can send it
	for _, acct := range msg.GetSigners() {
		if !acct.Equals(contractAddr) {
			return nil, nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "contract doesn't have permission")
		}
	}

	var res *sdk.Result
	var err error
	if legacyMsg, ok := msg.(legacytx.LegacyMsg); ok {
		msgRoute := legacyMsg.Route()
		handler := h.legacyRouter.Route(ctx, msgRoute)
		if handler == nil {
			return nil, nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized message route: %s", msgRoute)
		}

		res, err = handler(ctx, msg)
		if err != nil {
			var errData []byte
			errData = nil
			if res != nil {
				errData = make([]byte, len(res.Data))
				copy(errData, res.Data)
			}

			return nil, errData, err
		}
	} else {
		return nil, nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized legacy message route: %s", sdk.MsgTypeURL(msg))

		// todo: grpc routing
		//handler := k.serviceRouter.Handler(msg)
		//if handler == nil {
		//	return nil, nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, sdk.MsgTypeURL(msg))
		//}
		//res, err := handler(ctx, msg)
		//if err != nil {
		//	return nil, nil, err
		//}
	}

	var events []sdk.Event
	data := make([]byte, len(res.Data))
	copy(data, res.Data)
	//
	// convert Tendermint.Events to sdk.Event
	sdkEvents := make(sdk.Events, len(res.Events))
	for i := range res.Events {
		sdkEvents[i] = sdk.Event(res.Events[i])
	}

	// append message action attribute
	events = append(events, sdkEvents...)

	return events, data, nil
}

// convertWasmIBCTimeoutHeightToCosmosHeight converts a wasm type ibc timeout height to ibc module type height
func convertWasmIBCTimeoutHeightToCosmosHeight(ibcTimeoutBlock *wasmTypes.IBCTimeoutBlock) ibcclienttypes.Height {
	if ibcTimeoutBlock == nil {
		return ibcclienttypes.NewHeight(0, 0)
	}
	return ibcclienttypes.NewHeight(ibcTimeoutBlock.Revision, ibcTimeoutBlock.Height)
}

func convertWasmCoinsToSdkCoins(coins []wasmTypes.Coin) (sdk.Coins, error) {
	var toSend sdk.Coins
	for _, coin := range coins {
		c, err := convertWasmCoinToSdkCoin(coin)
		if err != nil {
			return nil, err
		}
		toSend = append(toSend, c)
	}
	return toSend, nil
}

func convertWasmCoinToSdkCoin(coin wasmTypes.Coin) (sdk.Coin, error) {
	amount, ok := sdk.NewIntFromString(coin.Amount)
	if !ok {
		return sdk.Coin{}, sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, coin.Amount+coin.Denom)
	}
	return sdk.Coin{
		Denom:  coin.Denom,
		Amount: amount,
	}, nil
}
