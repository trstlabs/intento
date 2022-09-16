package cli

import (
	"fmt"
	"strconv"

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
			creator, err := cmd.Flags().GetString(flagCreator)
			if err != nil {
				return fmt.Errorf("run-as: %s", err)
			}
			if len(creator) == 0 {
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
				Creator:             creator,
				WASMByteCode:        src.WASMByteCode,
				ContractTitle:       src.Title,
				ContractDescription: src.Description,
				DefaultDuration:     src.DefaultDuration,
				DefaultInterval:     src.DefaultInterval,
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

	cmd.Flags().String(flagCreator, "", "The address that is stored as code creator")
	//cmd.Flags().String(flagInstantiateByEverybody, "", "Everybody can instantiate a contract from the code, optional")
	//cmd.Flags().String(flagInstantiateByAddress, "", "Only this address can instantiate a contract instance from the code, optional")
	cmd.Flags().String(flagSource, "", "A valid URI reference to the contract's source code, optional")
	cmd.Flags().String(flagBuilder, "", "A valid docker tag for the build system, optional")
	// proposal flags
	cmd.Flags().String(flagDuration, "", "A default duration for the contract e.g. 2h, 6000s, 72h3m0.5s, optional")
	cmd.Flags().String(flagInterval, "", "A default interval for the contract e.g. 2h, 6000s, 72h3m0.5s, optional")
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

func ProposalInstantiateContractCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instantiate-contract [code_id_int64] [json_encoded_init_args] --contract_id [text] --title [proposal text] --description [proposal text] --amount [coins,optional] --deposit [coins,optional] --auto_msg [json args, optional]  --duration [custom duration e.g. 400s/5h] (optional)  --interval [custom dration e.g. 400s/5h]  (optional) --start_at [UNIX time]",
		Short: "Submit an instantiate wasm contract proposal (run by community)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			/*runAs, err := cmd.Flags().GetString(flagRunAs)
			if err != nil {
				return fmt.Errorf("run-as: %s", err)
			}
			if len(runAs) == 0 {
				return errors.New("run-as address is required")
			}*/
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
			codeID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			contractId, _ := cmd.Flags().GetString(flagContractId)
			if contractId == "" {
				return fmt.Errorf("contract id is required on all contracts")
			}
			autoMsg, err := cmd.Flags().GetString(flagAutoMsg)
			if err != nil {
				return err
			}
			amountStr, err := cmd.Flags().GetString(flagAmount)
			if err != nil {
				return fmt.Errorf("amount: %s", err)
			}

			funds, err := sdk.ParseCoinsNormalized(amountStr)
			if err != nil {
				return err
			}
			duration, err := cmd.Flags().GetString(flagDuration) //strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return fmt.Errorf("contract duration: %s", err)
			}
			interval, err := cmd.Flags().GetString(flagInterval) //strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return fmt.Errorf("contract interval: %s", err)
			}

			startAtStr, err := cmd.Flags().GetString(flagStartAt) //strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return fmt.Errorf("startAt string: %s", err)
			}

			startAt, err := strconv.ParseUint(startAtStr, 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse start duration at: %s", err)
			}
			content := types.InstantiateContractProposal{
				Title:       proposalTitle,
				Description: proposalDescr,
				//RunAs:       runAs,
				//Proposer: clientCtx.GetFromAddress().String(),
				//Admin:       src.Admin,
				CodeID:          codeID,
				ContractId:      contractId,
				Msg:             []byte(args[1]),
				AutoMsg:         []byte(autoMsg),
				Funds:           funds,
				Duration:        duration,
				Interval:        interval,
				StartDurationAt: startAt,
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
	cmd.Flags().String(flagAmount, "", "Coins to send to the contract during instantiation")
	cmd.Flags().String(flagContractId, "", "A human-readable name for this contract in lists")
	//cmd.Flags().String(flagAdmin, "", "Address of an admin")
	//cmd.Flags().String(flagRunAs, "", "The address that runs the contract. It is the creator of the contract and passed to the contract as sender on proposal execution")
	cmd.Flags().String(flagAutoMsg, "", "An automatic message to send, that the contract executes after a set duration (optional)")

	// proposal flags
	cmd.Flags().String(cli.FlagTitle, "", "Title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "Description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "Deposit of proposal")
	cmd.Flags().String(cli.FlagProposal, "", "Proposal file path (if this path is given, other proposal flags are ignored)")
	// type values must match the "ProposalHandler" "routes" in cli
	cmd.Flags().String(flagProposalType, "", "Permission of proposal, types: store-code/instantiate/migrate/update-admin/clear-admin/text/parameter_change/software_upgrade")
	cmd.Flags().String(flagDuration, "", "A custom duration for the contract e.g. 2h, 6000s, 72h3m0.5s, optional")
	cmd.Flags().String(flagInterval, "", "A custom interval for the contract e.g. 2h, 6000s, 72h3m0.5s, optional")
	cmd.Flags().String(flagStartAt, "0", "A custom start time for the contract self-execution, in UNIX time")
	return cmd
}

func ProposalExecuteContractCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execute-contract [contract_addr_bech32] [json_encoded_migration_args] --title [text] --description [text]  --amount [coins,optional]",
		Short: "Submit a execute wasm contract proposal (run by community)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			//	contract := args[0]
			//	execMsg := []byte(args[1])
			amountStr, err := cmd.Flags().GetString(flagAmount)
			if err != nil {
				return fmt.Errorf("amount: %s", err)
			}
			// get address to execute
			contractAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			src, err := parseExecuteArgs(cmd, contractAddr, []byte(args[1]), amountStr, false, "", "", clientCtx)
			if err != nil {
				return err
			}
			/*
				funds, err := sdk.ParseCoinsNormalized(amountStr)
				if err != nil {
					return fmt.Errorf("amount: %s", err)
				}

					/*runAs, err := cmd.Flags().GetString(flagRunAs)
					if err != nil {
						return fmt.Errorf("run-as: %s", err)
					}

					if len(runAs) == 0 {
						return errors.New("run-as address is required")
					}*/
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

			content := types.ExecuteContractProposal{
				Title:       proposalTitle,
				Description: proposalDescr,
				Contract:    src.Contract,
				Msg:         src.Msg,
				//RunAs:       runAs,
				//Proposer:  clientCtx.GetFromAddress().String(),
				Funds: src.Funds,
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
	//cmd.Flags().String(flagRunAs, "", "The address that is passed as sender to the contract on proposal execution")
	cmd.Flags().String(flagAmount, "", "Coins to send to the contract during instantiation")
	cmd.Flags().String(flagContractId, "", "A human-readable name for this contract in lists")

	// proposal flags
	cmd.Flags().String(cli.FlagTitle, "", "Title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "Description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "Deposit of proposal")
	cmd.Flags().String(cli.FlagProposal, "", "Proposal file path (if this path is given, other proposal flags are ignored)")
	// type values must match the "ProposalHandler" "routes" in cli
	cmd.Flags().String(flagProposalType, "", "Permission of proposal, types: store-code/instantiate/migrate/update-admin/clear-admin/text/parameter_change/software_upgrade")
	return cmd
}
