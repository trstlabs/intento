package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/danieljdd/tpp/x/tpp/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	// this line is used by starport scaffolding # 1

	cmd.AddCommand(CmdCreateEstimator())
	cmd.AddCommand(CmdUpdateEstimator())
	cmd.AddCommand(CmdDeleteEstimator())

	cmd.AddCommand(CmdCreateBuyer())
	cmd.AddCommand(CmdUpdateBuyer())
	cmd.AddCommand(CmdDeleteBuyer())

	cmd.AddCommand(CmdCreateItem())
	cmd.AddCommand(CmdUpdateItem())
	cmd.AddCommand(CmdDeleteItem())

	return cmd
}
