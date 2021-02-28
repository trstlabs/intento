package cli

import (
	"github.com/spf13/cobra"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/danieljdd/tpp/x/tpp/types"
	//"cosmos/base/v1beta1/coin.proto"
)

func CmdCreateBuyer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-buyer [itemid] [deposit]",
		Short: "Creates a new buyer",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsItemid := string(args[0])
			
			argsDeposit, _ := sdk.ParseCoinNormalized(args[2])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateBuyer(clientCtx.GetFromAddress().String(), string(argsItemid), sdk.Coin(argsDeposit))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdUpdateBuyer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-buyer [id] [itemid] [transferable] [deposit]",
		Short: "Update a buyer",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {

			argsItemid := string(args[1])
			argsTransferable := false
			if args[2] == "1" {
				argsTransferable = true
			}
			argsDeposit,_ := sdk.ParseCoinNormalized(args[3])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateBuyer(clientCtx.GetFromAddress().String(), string(argsItemid), bool(argsTransferable), sdk.Coin(argsDeposit))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdDeleteBuyer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-buyer [id] [itemid] [transferable] [deposit]",
		Short: "Delete a buyer by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteBuyer(clientCtx.GetFromAddress().String(), id)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdItemTransfer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "item-transfer [transterbool] [itemID]",
		Short: "Set a new buyer",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsTransferbool := false
			if args[0] == "1" {
				argsTransferbool = true
			}
			argsItemID := string(args[1])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgItemTransfer(clientCtx.GetFromAddress().String(), string(argsItemID), bool(argsTransferbool))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}