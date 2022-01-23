package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/trstlabs/trst/x/compute/internal/types"

	// authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	// "github.com/trstlabs/trst/x/compute/internal/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// InitGenesis sets supply information for genesis.
//
// CONTRACT: all types of accounts must have been already initialized/created
// InitGenesis initializes the trst module state
func (k Keeper) InitGenesis(ctx sdk.Context, gs types.GenesisState) error {
	// NOTE: since the compute module is a module account, the auth module should
	// take care of importing the amount into the account except for the
	// genesis block

	if k.GetComputeModuleBalance(ctx).IsZero() {
		err := k.InitializeComputeModule(ctx, sdk.NewCoin("utrst", sdk.ZeroInt()))
		if err != nil {
			panic(err)
		}
	}
	var maxCodeID uint64
	for i, code := range gs.Codes {
		err := k.importCode(ctx, code.CodeID, code.CodeInfo, code.CodeBytes)
		if err != nil {
			return sdkerrors.Wrapf(err, "code %d with id: %d", i, code.CodeID)
		}
		if code.CodeID > maxCodeID {
			maxCodeID = code.CodeID
		}
	}

	var maxContractID int
	for i, contract := range gs.Contracts {
		err := k.importContract(ctx, contract.ContractAddress, &contract.ContractInfo, contract.ContractState)
		if err != nil {
			return sdkerrors.Wrapf(err, "contract number %d", i)
		}
		maxContractID = i + 1 // not ideal but max(contractID) is not persisted otherwise
	}

	for i, seq := range gs.Sequences {
		err := k.importAutoIncrementID(ctx, seq.IDKey, seq.Value)
		if err != nil {
			return sdkerrors.Wrapf(err, "sequence number %d", i)
		}
	}

	// sanity check seq values
	if k.peekAutoIncrementID(ctx, types.KeyLastCodeID) <= maxCodeID {
		return sdkerrors.Wrapf(types.ErrInvalid, "seq %s must be greater %d ", string(types.KeyLastCodeID), maxCodeID)
	}
	if k.peekAutoIncrementID(ctx, types.KeyLastInstanceID) <= uint64(maxContractID) {
		return sdkerrors.Wrapf(types.ErrInvalid, "seq %s must be greater %d ", string(types.KeyLastInstanceID), maxContractID)
	}
	//fmt.Print("setting paams...")
	k.SetParams(ctx, types.DefaultParams())
	//keeper.setParams(ctx, data.Params)

	return nil
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) *types.GenesisState {
	//var genState types.GenesisState
	genState := *types.DefaultGenesis()

	//genState.Params = keeper.GetParams(ctx)
	genState.Params = keeper.GetParams(ctx)

	keeper.IterateCodeInfos(ctx, func(codeID uint64, info types.CodeInfo) bool {
		bytecode, err := keeper.GetByteCode(ctx, codeID)
		if err != nil {
			panic(err)
		}
		genState.Codes = append(genState.Codes, types.Code{
			CodeID:    codeID,
			CodeInfo:  info,
			CodeBytes: bytecode,
		})
		return false
	})

	keeper.IterateContractInfo(ctx, func(addr sdk.AccAddress, contract types.ContractInfo) bool {
		contractStateIterator := keeper.GetContractState(ctx, addr)
		var state []types.Model
		for ; contractStateIterator.Valid(); contractStateIterator.Next() {
			m := types.Model{
				Key:   contractStateIterator.Key(),
				Value: contractStateIterator.Value(),
			}
			state = append(state, m)
		}
		// redact contract info
		contract.Created = nil

		genState.Contracts = append(genState.Contracts, types.Contract{
			ContractAddress: addr,
			ContractInfo:    contract,
			ContractState:   state,
		})

		return false
	})

	for _, k := range [][]byte{types.KeyLastCodeID, types.KeyLastInstanceID} {
		genState.Sequences = append(genState.Sequences, types.Sequence{
			IDKey: k,
			Value: keeper.peekAutoIncrementID(ctx, k),
		})
	}

	return &genState
}

// GetComputeModuleAccount returns the module account.
func (k Keeper) GetComputeModuleAccount(ctx sdk.Context) (ModuleName authtypes.ModuleAccountI) {
	return k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
}

// GetComputeModuleBalance returns the module account balance
func (k Keeper) GetComputeModuleBalance(ctx sdk.Context) sdk.Coin {
	return k.bankKeeper.GetBalance(ctx, k.GetComputeModuleAccount(ctx).GetAddress(), "utrst")
}

// InitializeComputeModule sets up the module account from genesis
func (k Keeper) InitializeComputeModule(ctx sdk.Context, funds sdk.Coin) error {
	return k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(funds))
}
