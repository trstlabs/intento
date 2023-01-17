package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

func (s Sequence) ValidateBasic() error {
	if len(s.IDKey) == 0 {
		return sdkerrors.Wrap(ErrEmpty, "id key")
	}
	return nil
}

func (s GenesisState) ValidateBasic() error {
	//fmt.Print("VALIDATING BASIC")
	if err := s.Params.Validate(); err != nil {
		return sdkerrors.Wrap(err, "params")
	}

	for i := range s.AutoTxInfos {
		if err := s.AutoTxInfos[i].ValidateBasic(); err != nil {
			return sdkerrors.Wrapf(err, "auto_tx: %d", i)
		}

	}
	for i := range s.Sequences {
		if err := s.Sequences[i].ValidateBasic(); err != nil {
			return sdkerrors.Wrapf(err, "sequence: %d", i)
		}
	}
	return nil
}

func (a AutoTxInfo) ValidateBasic() error {
	if a.TxID == 0 {
		return sdkerrors.Wrap(ErrEmpty, "code id")
	}

	return nil
}

// ValidateGenesis performs basic validation of supply genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(gs GenesisState) error {
	return gs.ValidateBasic()
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	//fmt.Println("Validating")
	gs.ValidateBasic()
	err := gs.Params.Validate()
	if err != nil {
		return err
	}

	return nil
}
