package appparams

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
)

// CustomAppConfig defines the configuration for the Nois app.
type CustomAppConfig struct {
	serverconfig.Config
	Wasm    wasmtypes.NodeConfig `mapstructure:"wasm" json:"wasm"`
	Intento IntentoConfig        `mapstructure:"intento" json:"intento"`
}

type IntentoConfig struct {
	PoA PoAConfig `mapstructure:"poa" json:"poa"`
}

type PoAConfig struct {
	DisableGatekeeping bool `mapstructure:"disable_gatekeeping" json:"disable_gatekeeping"`
}

func CustomconfigTemplate(config wasmtypes.NodeConfig) string {
	return serverconfig.DefaultConfigTemplate + wasmtypes.ConfigTemplate(config) + `
###############################################################################
###                        Custom Intento Configuration                     ###
###############################################################################

[intento.poa]
# disable_gatekeeping allows direct usage of MsgCreateValidator.
# If false, MsgCreateValidator is gated by governance.
disable_gatekeeping = false
`
}

func DefaultConfig() (string, interface{}) {
	serverConfig := serverconfig.DefaultConfig()
	serverConfig.MinGasPrices = "0uinto"

	wasmConfig := wasmtypes.DefaultNodeConfig()
	simulationLimit := uint64(50_000_000)

	wasmConfig.SimulationGasLimit = &simulationLimit
	wasmConfig.SmartQueryGasLimit = 25_000_000
	wasmConfig.MemoryCacheSize = 512
	wasmConfig.ContractDebugMode = false

	customConfig := CustomAppConfig{
		Config: *serverConfig,
		Wasm:   wasmConfig,
		Intento: IntentoConfig{
			PoA: PoAConfig{
				DisableGatekeeping: false,
			},
		},
	}

	return CustomconfigTemplate(wasmConfig), customConfig
}
