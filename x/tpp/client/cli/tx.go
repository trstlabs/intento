package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/danieljdd/trst/x/trst/types"
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

	cmd.AddCommand(CmdCreateEstimation())
	cmd.AddCommand(CmdUpdateLike())
	cmd.AddCommand(CmdDeleteEstimation())
	cmd.AddCommand(CmdFlagItem())

	cmd.AddCommand(CmdPrepayment())

	cmd.AddCommand(CmdWithdrawal())
	cmd.AddCommand(CmdItemTransfer())
	cmd.AddCommand(CmdItemRating())

	cmd.AddCommand(CmdCreateItem())

	cmd.AddCommand(CmdDeleteItem())
	cmd.AddCommand(CmdRevealEstimation())
	cmd.AddCommand(CmdItemTransferable())
	cmd.AddCommand(CmdItemShipping())
	cmd.AddCommand(CmdItemResell())
	cmd.AddCommand(CmdTokenizeItem())
	cmd.AddCommand(CmdUnTokenizeItem())

	return cmd
}
