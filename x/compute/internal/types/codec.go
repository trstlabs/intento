package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// "github.com/cosmos/cosmos-sdk/x/supply/exported"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// RegisterCodec registers the account types and interface
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(MsgStoreCode{}, "wasm/MsgStoreCode", nil) //TODO Remove for mainnet, code is to be uploaded through governance
	cdc.RegisterConcrete(&StoreCodeProposal{}, "wasm/StoreCodeProposal", nil)
	cdc.RegisterConcrete(MsgInstantiateContract{}, "wasm/MsgInstantiateContract", nil)
	cdc.RegisterConcrete(MsgExecuteContract{}, "wasm/MsgExecuteContract", nil)

	cdc.RegisterConcrete(MsgDiscardAutoMsg{}, "wasm/MsgDiscardAutoMsg", nil)

	cdc.RegisterConcrete(&InstantiateContractProposal{}, "wasm/InstantiateContractProposal", nil)
	cdc.RegisterConcrete(&ExecuteContractProposal{}, "wasm/ExecuteContractProposal", nil)

}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgStoreCode{}, //TODO Remove for mainnet, code to be uploaded through governance
		&MsgInstantiateContract{},
		&MsgExecuteContract{},
		&MsgDiscardAutoMsg{},
	)
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&StoreCodeProposal{},
		&InstantiateContractProposal{},
		&ExecuteContractProposal{},
	)
}

// ModuleCdc generic sealed codec to be used throughout module
var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/wasm module codec.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
