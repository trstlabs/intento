package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/danieljdd/tpp/x/tpp/types"
	"github.com/spf13/cobra"
	//"cosmos/base/v1beta1/coin.proto"
)

func CmdCreateBuyer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-buyer [itemid] [deposit amount]",
		Short: "Creates a new buyer",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsItemid := string(args[0])

			argsDeposit, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return err
			}
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateBuyer(clientCtx.GetFromAddress().String(), string(argsItemid), int64(argsDeposit))
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
		Use:   "update-buyer[itemid] [deposit]",
		Short: "Update a buyer",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			argsItemid := string(args[0])
		
			argsDeposit, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return err
			}
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateBuyer(clientCtx.GetFromAddress().String(), string(argsItemid), int64(argsDeposit))
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
		Use:   "delete-buyer [id]",
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
		Use:   "item-transfer [itemID]",
		Short: "Set a new buyer",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			
			argsItemID := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgItemTransfer(clientCtx.GetFromAddress().String(), string(argsItemID))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdItemRating() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "item-rating [Rating] [Note] [itemID]",
		Short: "Set a new buyer",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
	

			argsRating, _ := strconv.ParseInt(args[0], 10, 64)


			argsNote := string(args[1])
			argsItemID := string(args[2])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgItemRating(clientCtx.GetFromAddress().String(), string(argsItemID), int64(argsRating), string(argsNote))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
