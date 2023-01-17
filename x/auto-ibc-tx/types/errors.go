package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrIBCAccountAlreadyExist = sdkerrors.Register(ModuleName, 2, "interchain account already registered")
	ErrIBCAccountNotExist     = sdkerrors.Register(ModuleName, 3, "interchain account not exist")
	ErrAccountExists          = sdkerrors.Register(ModuleName, 6, "fee account already exists")
	ErrDuplicate              = sdkerrors.Register(ModuleName, 14, "duplicate")
	ErrInvalid                = sdkerrors.Register(ModuleName, 1, "custom error message")
	ErrEmpty                  = sdkerrors.Register(ModuleName, 11, "empty")
)
