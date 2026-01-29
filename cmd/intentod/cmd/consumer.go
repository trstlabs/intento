package cmd

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	ccvconsumertypes "github.com/cosmos/interchain-security/v6/x/ccv/consumer/types"
	"github.com/spf13/cobra"
)

type GenesisMutator interface {
	AlterConsumerModuleState(cmd *cobra.Command, callback func(state *GenesisData, appState map[string]json.RawMessage) error) error
}

type DefaultGenesisIO struct {
	DefaultGenesisReader
}

func NewDefaultGenesisIO() *DefaultGenesisIO {
	return &DefaultGenesisIO{DefaultGenesisReader: DefaultGenesisReader{}}
}

func (x DefaultGenesisIO) AlterConsumerModuleState(cmd *cobra.Command, callback func(state *GenesisData, appState map[string]json.RawMessage) error) error {
	g, err := x.ReadGenesis(cmd)
	if err != nil {
		return err
	}
	if err := callback(g, g.AppState); err != nil {
		return err
	}
	if err := g.ConsumerModuleState.Validate(); err != nil {
		return err
	}
	clientCtx := client.GetClientContextFromCmd(cmd)
	consumerGenStateBz, err := clientCtx.Codec.MarshalJSON(g.ConsumerModuleState)
	if err != nil {
		return errors.Wrap(err, "marshal consumer genesis state")
	}

	g.AppState[ccvconsumertypes.ModuleName] = consumerGenStateBz
	appStateJSON, err := json.Marshal(g.AppState)
	if err != nil {
		return errors.Wrap(err, "marshal application genesis state")
	}

	g.AppGenesis.AppState = appStateJSON
	return genutil.ExportGenesisFile(g.AppGenesis, g.GenesisFile)
}

type DefaultGenesisReader struct{}

func (d DefaultGenesisReader) ReadGenesis(cmd *cobra.Command) (*GenesisData, error) {
	clientCtx := client.GetClientContextFromCmd(cmd)
	serverCtx := server.GetServerContextFromCmd(cmd)
	config := serverCtx.Config
	config.SetRoot(clientCtx.HomeDir)

	genFile := config.GenesisFile()
	appState, appGenesis, err := genutiltypes.GenesisStateFromGenFile(genFile)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal genesis state: %w", err)
	}

	return NewGenesisData(
		genFile,
		appGenesis,
		appState,
		nil,
	), nil
}

type GenesisData struct {
	GenesisFile         string
	AppGenesis          *genutiltypes.AppGenesis
	AppState            map[string]json.RawMessage
	ConsumerModuleState *ccvconsumertypes.GenesisState
}

func NewGenesisData(genesisFile string, appGenesis *genutiltypes.AppGenesis, appState map[string]json.RawMessage, consumerModuleState *ccvconsumertypes.GenesisState) *GenesisData {
	return &GenesisData{GenesisFile: genesisFile, AppGenesis: appGenesis, AppState: appState, ConsumerModuleState: consumerModuleState}
}
