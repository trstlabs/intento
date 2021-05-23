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

func CmdCreateEstimation() *cobra.Command {
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
			argsItemID, err := strconv.ParseUint(args[4], 10, 64)
			if err != nil {
				return err
			}
			msg := types.NewMsgCreateEstimation(clientCtx.GetFromAddress().String(), int64(argsEstimation), uint64(argsItemID), int64(argsDeposit), bool(interested), string(argsComment))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdUpdateLike() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-estimator [itemid] [interested]",
		Short: "Update a estimator",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			argsItemID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			interested := false
			if args[1] == "1" {
				interested = true
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateLike(clientCtx.GetFromAddress().String(), uint64(argsItemID), bool(interested))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdDeleteEstimation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-estimator [id] [estimation] [estimatorestimationhash] [itemid] [deposit] [interested] [comment] [flag]",
		Short: "Delete a estimator by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteEstimation(clientCtx.GetFromAddress().String(), uint64(id))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdFlagItem() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-flag [itemid]",
		Short: "create a new flag for item",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			itemid, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgFlagItem(clientCtx.GetFromAddress().String(), uint64(itemid))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
