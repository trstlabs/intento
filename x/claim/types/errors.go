package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/claim module sentinel errors
var (
	ErrIncorrectModuleAccountBalance = errorsmod.Register(ModuleName, 1100,
		"claim module account balance != sum of all claim record InitialClaimableAmounts")
)
