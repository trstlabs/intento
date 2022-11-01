package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

// Aliases for internal types
const (
	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	Bech32PrefixAccAddr = "trust"
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key
	Bech32PrefixAccPub = "trustpub"
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
	Bech32PrefixValAddr = "trustvaloper"
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
	Bech32PrefixValPub = "trustvaloperpub"
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
	Bech32PrefixConsAddr = "trustvalcons"
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
	Bech32PrefixConsPub = "trustvalconspub"
)

// AddressVerifier secret address verifier
var AddressVerifier = func(bytes []byte) error {
	// 20 bytes = module accounts, base accounts, secret contracts
	// 32 bytes = ICA accounts
	if len(bytes) != 20 && len(bytes) != 32 {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "address length must be 20 or 32 bytes, got %d", len(bytes))
	}

	return nil
}
