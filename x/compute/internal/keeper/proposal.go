package keeper

import (
	"encoding/hex"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/trstlabs/trst/x/compute/internal/types"
)

// NewWasmProposalHandler creates a new governance Handler for wasm proposals
func NewWasmProposalHandler(k Keeper, enabledProposalTypes []types.ProposalType) govtypes.Handler {
	enabledTypes := make(map[string]struct{}, len(enabledProposalTypes))
	for i := range enabledProposalTypes {
		enabledTypes[string(enabledProposalTypes[i])] = struct{}{}
	}
	return func(ctx sdk.Context, content govtypes.Content) error {
		if content == nil {
			return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "content must not be empty")
		}
		if _, ok := enabledTypes[content.ProposalType()]; !ok {
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unsupported wasm proposal content type: %q", content.ProposalType())
		}
		switch c := content.(type) {
		case *types.StoreCodeProposal:
			return handleStoreCodeProposal(k, ctx, *c)
		case *types.InstantiateContractProposal:
			return handleInstantiateProposal(ctx, k, *c)
		case *types.ExecuteContractProposal:
			return handleExecuteProposal(ctx, k, *c)
		case *types.UpdateCodeProposal:
			return handleUpdateCodeProposal(ctx, k, *c)
		case *types.UpdateContractProposal:
			return handleUpdateContractProposal(ctx, k, *c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized wasm proposal content type: %T", c)
		}
	}
}

func handleStoreCodeProposal(k Keeper, ctx sdk.Context, p types.StoreCodeProposal) error {
	if err := p.ValidateBasic(); err != nil {
		return err
	}

	creatorAddr, err := sdk.AccAddressFromBech32(p.Creator)
	if err != nil {
		return sdkerrors.Wrap(err, "run as address")
	}
	duration, err := time.ParseDuration(p.DefaultDuration)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}
	interval, err := time.ParseDuration(p.DefaultInterval)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}
	fmt.Printf("p %s \n", p.ContractTitle)
	fmt.Printf("p %s \n", p.ContractDescription)
	fmt.Printf("p %s \n", duration)
	fmt.Printf("p %s \n", creatorAddr)
	_, err = k.Create(ctx, creatorAddr, p.WASMByteCode, "", "", duration, interval, p.ContractTitle, p.ContractDescription)
	if err != nil {
		return sdkerrors.Wrap(err, "err creating code")
	}
	return nil
}

func handleInstantiateProposal(ctx sdk.Context, k Keeper, p types.InstantiateContractProposal) error {
	if err := p.ValidateBasic(); err != nil {
		return err
	}
	var encryptedAutoMsg []byte = nil
	sig, encryptedMsg, err := k.CreateCommunityPoolCallbackSig(ctx, p.Msg, p.CodeID, p.Funds)
	if err != nil {
		return err
	}
	if p.AutoMsg != nil {
		_, encryptedAutoMsg, err = k.CreateCommunityPoolCallbackSig(ctx, p.AutoMsg, p.CodeID, nil)
		if err != nil {
			return err
		}
	}
	duration, err := time.ParseDuration(p.Duration)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}
	interval, err := time.ParseDuration(p.Interval)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}
	var startTime time.Time = ctx.BlockHeader().Time
	if p.StartDurationAt != 0 {
		startTime = time.Unix(int64(p.StartDurationAt), 0)
		if err != nil {
			return err
		}
	}
	addr, _, err := k.Instantiate(ctx, p.CodeID, k.accountKeeper.GetModuleAddress("distribution"), encryptedMsg, encryptedAutoMsg, p.ContractId, p.Funds, sig, duration, interval, startTime)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeGovContractResult,
		sdk.NewAttribute(types.AttributeKeyAddress, hex.EncodeToString(addr)),
	))
	return nil
}

func handleExecuteProposal(ctx sdk.Context, k Keeper, p types.ExecuteContractProposal) error {
	if err := p.ValidateBasic(); err != nil {
		return err
	}

	contractAddr, err := sdk.AccAddressFromBech32(p.Contract)
	if err != nil {
		return sdkerrors.Wrap(err, "contract")
	}
	info := k.GetContractInfo(ctx, sdk.AccAddress(p.Contract))

	fmt.Printf("contract info %s", info)
	sig, encryptedMsg, err := k.CreateCommunityPoolCallbackSig(ctx, p.Msg, info.CodeID, p.Funds)
	if err != nil {
		return err
	}
	fmt.Printf("sig %s", sig)
	res, err := k.Execute(ctx, contractAddr, k.accountKeeper.GetModuleAddress("distribution"), encryptedMsg, p.Funds, sig)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeGovContractResult,
		sdk.NewAttribute(types.AttributeKeyResultDataHex, hex.EncodeToString(res.Data)),
	))
	return nil
}

func handleUpdateCodeProposal(ctx sdk.Context, k Keeper, p types.UpdateCodeProposal) error {
	if err := p.ValidateBasic(); err != nil {
		return err
	}

	err := k.SetCodeInfo(ctx, p.CodeId, p.Creator, p.DefaultDuration, p.DefaultInterval)
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeGovContractResult,
		sdk.NewAttribute(types.AttributeKeyCodeID, fmt.Sprintf("%d", p.CodeId)),
	),
	)
	return nil
}

func handleUpdateContractProposal(ctx sdk.Context, k Keeper, p types.UpdateContractProposal) error {
	if err := p.ValidateBasic(); err != nil {
		return err
	}

	err := k.UpdateContractInfo(ctx, p.Contract, p.Owner, p.StartTime, p.EndTime, p.Interval)
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeGovContractResult,
		sdk.NewAttribute(types.AttributeKeyAddress, p.Contract),
	),
	)
	return nil
}
