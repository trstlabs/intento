package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/tpp module sentinel errors
var (
	ErrSample = sdkerrors.Register(ModuleName, 1100, "sample error")
	ErrInvalid                      = sdkerrors.Register(ModuleName, 1, "custom error message")
	ErrArgumentMissingOrNonUInteger = sdkerrors.Register(ModuleName, 338, "argument is missing or is not an unsigned integer")
)
