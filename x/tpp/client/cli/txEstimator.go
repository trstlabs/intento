package cli

import (
	"github.com/spf13/cobra"
	//"crypto/sha256"
	//"encoding/hex"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/danieljdd/tpp/x/tpp/types"
)

func CmdCreateEstimator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-estimator [estimation] [deposit] [interested] [comment] [itemid]",
		Short: "Creates a new estimator",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			argsEstimation, _ := strconv.ParseInt(args[0], 10, 64)

			argsDeposit, _ := strconv.ParseInt(args[1], 10, 64)
			interested := false
			if args[2] == "1" {
				interested = true
			}
			argsComment := string(args[3])
			argsItemid := string(args[4])

			msg := types.NewMsgCreateEstimator(clientCtx.GetFromAddress().String(), int64(argsEstimation), string(argsItemid), int64(argsDeposit), bool(interested), string(argsComment))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdUpdateEstimator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-estimator [id] [estimation] [estimatorestimationhash] [itemid] [deposit] [interested] [comment] [flag]",
		Short: "Update a estimator",
		Args:  cobra.ExactArgs(8),
		RunE: func(cmd *cobra.Command, args []string) error {

			argsItemid := string(args[0])

			interested := false
			if args[1] == "1" {
				interested = true
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateEstimator(clientCtx.GetFromAddress().String(), string(argsItemid), bool(interested))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdDeleteEstimator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-estimator [id] [estimation] [estimatorestimationhash] [itemid] [deposit] [interested] [comment] [flag]",
		Short: "Delete a estimator by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteEstimator(clientCtx.GetFromAddress().String(), id)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdCreateFlag() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-flag [flag] [itemid] ",
		Short: "create a new flag for item",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			flag := false
			if args[0] == "1" {
				flag = true
			}

			itemid := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateFlag(clientCtx.GetFromAddress().String(), bool(flag), string(itemid))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
