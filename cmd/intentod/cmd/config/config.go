package config

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	appparams "github.com/trstlabs/intento/app/params"
)

const (
	// Bech32Prefix defines the Bech32 prefix used for EthAccounts
	Bech32Prefix = "into"

	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	Bech32PrefixAccAddr = Bech32Prefix
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key
	Bech32PrefixAccPub = Bech32Prefix + sdk.PrefixPublic
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
	Bech32PrefixValAddr = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
	Bech32PrefixValPub = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
	Bech32PrefixConsAddr = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
	Bech32PrefixConsPub = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus + sdk.PrefixPublic
)

const (
	HumanReadableCoinUnit = "INTO"
	BaseCoinUnit          = "uinto"
	DefaultBondDenom      = BaseCoinUnit
)

func SetupConfig() {
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount(appparams.Bech32PrefixAccAddr, appparams.Bech32PrefixAccPub)
	cfg.SetBech32PrefixForValidator(appparams.Bech32PrefixValAddr, appparams.Bech32PrefixValPub)
	cfg.SetBech32PrefixForConsensusNode(appparams.Bech32PrefixConsAddr, appparams.Bech32PrefixConsPub)
	cfg.Seal()
}

// SetBip44CoinType sets the global coin type to be used in hierarchical deterministic wallets.
func SetBip44CoinType(config *sdk.Config) {
	config.SetCoinType(sdk.CoinType)
	config.SetPurpose(sdk.Purpose) // Shared
}

// RegisterDenoms registers the base and display denominations to the SDK.
func RegisterDenoms() {
	if err := sdk.RegisterDenom(HumanReadableCoinUnit, math.LegacyOneDec()); err != nil {
		panic(err)
	}

	if err := sdk.RegisterDenom(BaseCoinUnit, math.LegacyNewDecWithPrec(1, 6)); err != nil {
		panic(err)
	}
}

// SetAddressPrefixes builds the Config with Bech32 addressPrefix and publKeyPrefix for accounts, validators, and consensus nodes and verifies that addresses have correct format.

func SetAddressPrefixes(config *sdk.Config) {
	// This is copied from the cosmos sdk v0.43.0-beta1
	// source: https://github.com/cosmos/cosmos-sdk/blob/v0.43.0-beta1/types/address.go#L141
	config.SetAddressVerifier(func(bytes []byte) error {
		if len(bytes) == 0 {
			return errorsmod.Wrap(sdkerrors.ErrUnknownAddress, "addresses cannot be empty")
		}

		if len(bytes) > address.MaxAddrLen {
			return errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "address max length is %d, got %d", address.MaxAddrLen, len(bytes))
		}

		// TODO: Do we want to allow addresses of lengths other than 20 and 32 bytes?
		if len(bytes) != 20 && len(bytes) != 32 {
			return errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "address length must be 20 or 32 bytes, got %d", len(bytes))
		}

		return nil
	})
}
