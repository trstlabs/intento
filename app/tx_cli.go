package app

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/spf13/cobra"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
)

const (
	FlagTitle       = "title"
	FlagDescription = "description"
	FlagDeposit     = "deposit"
)

// NewCmdSubmitValidatorAddProposal implements the command to submit a ValidatorAddProposal
func NewCmdSubmitValidatorAddProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-add [valoper-addr] [pubkey-json] [moniker]",
		Args:  cobra.ExactArgs(3),
		Short: "Submit a proposal to add a new validator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a proposal to add a new validator.

Example:
$ %s tx validator-add \
  intentvaloper1... \
  '{"@type":"/cosmos.crypto.ed25519.PubKey","key":"..."}' \
  "my-validator" \
  --title="Add Validator" --description="Adding my-validator" --deposit=1000000stake
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			valoper := args[0]
			pubKeyStr := args[1]
			moniker := args[2]

			var pk cryptotypes.PubKey
			if err := clientCtx.Codec.UnmarshalInterfaceJSON([]byte(pubKeyStr), &pk); err != nil {
				return err
			}

			pubKey, err := codectypes.NewAnyWithValue(pk)
			if err != nil {
				return err
			}

			// Validate Title and Description
			title, err := cmd.Flags().GetString(FlagTitle)
			if err != nil {
				return err
			}
			description, err := cmd.Flags().GetString(FlagDescription)
			if err != nil {
				return err
			}
			depositStr, err := cmd.Flags().GetString(FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			content := NewValidatorAddProposal(title, description, valoper, *pubKey, moniker)
			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, clientCtx.GetFromAddress())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagTitle, "", "title of proposal")
	cmd.Flags().String(FlagDescription, "", "description of proposal")
	cmd.Flags().String(FlagDeposit, "", "deposit of proposal")

	return cmd
}

// NewCmdSubmitValidatorRemoveProposal implements the command to submit a ValidatorRemoveProposal
func NewCmdSubmitValidatorRemoveProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-remove [valoper-addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a proposal to remove a validator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a proposal to remove a validator.

Example:
$ %s tx validator-remove \
  intentvaloper1... \
  --title="Remove Validator" --description="Removing validator" --deposit=1000000stake
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			valoper := args[0]

			// Validate Title and Description
			title, err := cmd.Flags().GetString(FlagTitle)
			if err != nil {
				return err
			}
			description, err := cmd.Flags().GetString(FlagDescription)
			if err != nil {
				return err
			}
			depositStr, err := cmd.Flags().GetString(FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			content := NewValidatorRemoveProposal(title, description, valoper)
			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, clientCtx.GetFromAddress())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagTitle, "", "title of proposal")
	cmd.Flags().String(FlagDescription, "", "description of proposal")
	cmd.Flags().String(FlagDeposit, "", "deposit of proposal")

	return cmd
}
