package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/trstlabs/trst/x/compute/internal/types"
)

func ProposalStoreCodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wasm-store [wasm file] --title [proposal text] --description [proposal text] --contract-title [text] --contract-description [text] --run-as [address]",
		Short: "Submit a wasm binary proposal",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			src, err := parseStoreCodeArgs(args, clientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			runAs, err := cmd.Flags().GetString(flagRunAs)
			if err != nil {
				return fmt.Errorf("run-as: %s", err)
			}
			if len(runAs) == 0 {
				return errors.New("run-as address is required")
			}
			proposalTitle, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return fmt.Errorf("proposal title: %s", err)
			}
			proposalDescr, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return fmt.Errorf("proposal description: %s", err)
			}
			depositArg, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositArg)
			if err != nil {
				return err
			}

			content := types.StoreCodeProposal{
				Title:               proposalTitle,
				Description:         proposalDescr,
				RunAs:               runAs,
				WASMByteCode:        src.WASMByteCode,
				ContractTitle:       src.Title,
				ContractDescription: src.Description,
				ContractDuration:    src.ContractPeriod,
			}

			msg, err := govtypes.NewMsgSubmitProposal(&content, deposit, clientCtx.GetFromAddress())
			if err != nil {
				return err
			}
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(flagRunAs, "", "The address that is stored as code creator")
	//cmd.Flags().String(flagInstantiateByEverybody, "", "Everybody can instantiate a contract from the code, optional")
	//cmd.Flags().String(flagInstantiateByAddress, "", "Only this address can instantiate a contract instance from the code, optional")
	cmd.Flags().String(flagSource, "", "A valid URI reference to the contract's source code, optional")
	cmd.Flags().String(flagBuilder, "", "A valid docker tag for the build system, optional")

	cmd.Flags().String(flagDuration, "", "A max duration for the contract e.g. 2h, 6000s, 72h3m0.5s, optional")
	// proposal flags
	cmd.Flags().String(cli.FlagTitle, "", "Title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "Description of proposal")
	cmd.Flags().String(flagTitle, "", "Title of contract")
	cmd.Flags().String(flagDescription, "", "Description of contract")
	cmd.Flags().String(cli.FlagDeposit, "", "Deposit of proposal")
	cmd.Flags().String(cli.FlagProposal, "", "Proposal file path (if this path is given, other proposal flags are ignored)")
	// type values must match the "ProposalHandler" "routes" in cli
	cmd.Flags().String(flagProposalType, "", "Permission of proposal, types: store-code/instantiate/migrate/update-admin/clear-admin/text/parameter_change/software_upgrade")
	return cmd
}
