package cli

import (
	//"crypto/sha256"
	//"encoding/hex"
	//"fmt"
	"strings"
	//sdk "github.com/cosmos/cosmos-sdk/types"
	//sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"strconv"

	"github.com/spf13/cobra"

	//sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/danieljdd/tpp/x/tpp/types"
	//"cosmos/base/v1beta1/coin.proto"
)

func CmdCreateItem() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-item [title] [description] [shippingcost] [localpickup] [estimationcount] [tags] [condition] [shippingregion]",
		Short: "Creates a new item",
		Args:  cobra.ExactArgs(9),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsTitle := string(args[0])
			argsDescription := string(args[1])
			argsShippingcost, _ := strconv.ParseInt(args[2], 10, 64)
			argsLocalpickup := true
			if args[3] == "0" {
				argsLocalpickup = false
			}

			//argsEstimationcounthash := string(args[4])

			//estimationcheck, ok := sdk.NewIntFromString(args[4])
			//if ok != true {
			//	return sdkerrors.Wrap(types.ErrArgumentMissingOrNonUInteger, "not a number or lower than zero")
			//}

			argsTags := strings.Split(args[5], ",")

			argsEstimationcount, _ := strconv.ParseInt(args[4], 10, 64)

			argsCondition, _ := strconv.ParseInt(args[6], 10, 64)

			argsShippingregion := strings.Split(args[7], ",")

			argsDepositAmount, _ := strconv.ParseInt(args[8], 10, 64)

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateItem(clientCtx.GetFromAddress().String(), string(argsTitle), string(argsDescription), int64(argsShippingcost), bool(argsLocalpickup), int64(argsEstimationcount), []string(argsTags), int64(argsCondition), []string(argsShippingregion), int64(argsDepositAmount))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdUpdateItem() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-item [id]  [shippingcost] [localpickup] [shippingregion]",
		Short: "Update a item",
		Args:  cobra.ExactArgs(7),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			argsShippingcost, _ := strconv.ParseInt(args[1], 10, 64)
			argsLocalpickup := true
			if args[2] == "0" {
				argsLocalpickup = false
			}

			argsShippingregion := strings.Split(args[3], ",")

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateItem(clientCtx.GetFromAddress().String(), id, int64(argsShippingcost), bool(argsLocalpickup), []string(argsShippingregion))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdDeleteItem() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-item [id] [title] [description] [shippingcost] [localpickup] [estimationcounthash] [bestestimator] [lowestestimator] [highestestimator] [estimationprice] [estimatorlist] [estimatorestimationhashlist] [transferable] [buyer] [tracking] [status] [comments] [tags] [flags] [condition] [shippingregion]",
		Short: "Delete a item by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteItem(clientCtx.GetFromAddress().String(), id)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdRevealEstimation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reveal-estimation [itemID]",
		Short: "reveal a new estimation by itemID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			itemID := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgRevealEstimation(clientCtx.GetFromAddress().String(), string(itemID))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdItemTransferable() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "item-transferable [yes/no] [itemid]",
		Short: "set item transferability",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			transferBool := true
			if args[0] == "0" {
				transferBool = false
			}

			itemID := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgItemTransferable(clientCtx.GetFromAddress().String(), bool(transferBool), string(itemID))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdItemShipping() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "item-shipping [yes/no] [itemid]",
		Short: "set item transferability",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			shippingtrackingBool := false
			if args[0] == "1" {
				shippingtrackingBool = true
			}
			itemID := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgItemShipping(clientCtx.GetFromAddress().String(), bool(shippingtrackingBool), string(itemID))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
