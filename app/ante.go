package app

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func GateCreateValidatorAnteHandler(disableGatekeeping bool) sdk.AnteDecorator {
	return gateCreateValidatorDecorator{disableGatekeeping: disableGatekeeping}
}

type gateCreateValidatorDecorator struct {
	disableGatekeeping bool
}

func (d gateCreateValidatorDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	if d.disableGatekeeping {
		return next(ctx, tx, simulate)
	}

	msgs := tx.GetMsgs()
	for _, m := range msgs {
		switch m.(type) {
		case *stakingtypes.MsgCreateValidator:
			return ctx, errorsmod.Wrap(
				sdkerrors.ErrUnauthorized,
				"MsgCreateValidator is gated by governance; direct usage is disabled",
			)
		}
	}
	return next(ctx, tx, simulate)
}
