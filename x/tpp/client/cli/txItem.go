package cli

import (
	//"crypto/sha256"
	//"encoding/hex"
	//"fmt"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	//sdk "github.com/cosmos/cosmos-sdk/types"

	//sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/danieljdd/tpp/x/tpp/types"

	//"cosmos/base/v1beta1/coin.proto"
	wasmUtils "github.com/danieljdd/tpp/x/compute/client/utils"
)

func CmdCreateItem() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-item [title] [description] [shippingcost] [localpickup] [estimationcount] [tags] [condition] [shippingregion] [depositamount]",
		Short: "Creates a new item",
		Args:  cobra.ExactArgs(9),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			argsTitle := string(args[0])
			argsDescription := string(args[1])
			argsShippingcost, _ := strconv.ParseInt(args[2], 10, 64)
			argsLocalpickup := string(args[3])

			//argsEstimationcounthash := string(args[4])

			//estimationcheck, ok := sdk.NewIntFromString(args[4])
			//if ok != true {
			//	return sdkerrors.Wrap(types.ErrArgumentMissingOrNonUInteger, "not a number or lower than zero")
			//}

			argsTags := strings.Split(args[5], ",")

			argsEstimationcount, _ := strconv.ParseInt(args[4], 10, 64)
			wasmCtx := wasmUtils.WASMContext{CLIContext: cliCtx}

			initMsg := types.SecretMsg{}

			count := map[string]string{"estimationcount": args[4]}
			//initMsg.Msg = []byte("{\"estimationcount\": \"3\"}")
			initMsg.Msg, err = json.Marshal(count)
			//fmt.Printf("json message: %X\n", estimation)
			if err != nil {
				return err
			}
			//initMsg.Msg = []byte(initMsg.Msg)
			//fmt.Printf("message: %X\n", initMsg.Msg)
			//quite a long way to get a single value, however we can't directy access the keeper
			queryClient := types.NewQueryClient(cliCtx)
			params := &types.QueryCodeHashRequest{
				Codeid: 1,
			}
			res, err := queryClient.CodeHash(context.Background(), params)
			if err != nil {
				return err
			}

			fmt.Printf("Got code hash: %X\n", res.Codehash)
			var encryptedMsg []byte

			initMsg.CodeHash = []byte(hex.EncodeToString(res.Codehash))
			fmt.Printf("Got initMsg.CodeHash hash: %X\n", initMsg.CodeHash)
			encryptedMsg, err = wasmCtx.Encrypt(initMsg.Serialize())
			if err != nil {
				return err
			}
			fmt.Printf("encryptedMsg: %X\n", encryptedMsg)
			argsCondition, _ := strconv.ParseInt(args[6], 10, 64)

			argsShippingregion := strings.Split(args[7], ",")

			argsDepositAmount, _ := strconv.ParseInt(args[8], 10, 64)

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateItem(clientCtx.GetFromAddress().String(), string(argsTitle), string(argsDescription), int64(argsShippingcost), string(argsLocalpickup), int64(argsEstimationcount), []string(argsTags), int64(argsCondition), []string(argsShippingregion), int64(argsDepositAmount), encryptedMsg)
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
		Use:   "delete-item [item id] ",
		Short: "Delete a item by item id",
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

			msg := types.NewMsgDeleteItem(clientCtx.GetFromAddress().String(), uint64(id))
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
		Use:   "reveal-estimation [item ID]",
		Short: "reveal a new estimation by item ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			argsItemID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			wasmCtx := wasmUtils.WASMContext{CLIContext: cliCtx}

			revealMsg := types.SecretMsg{}
			reveal := types.ParseReveal{}

			//initMsg.Msg = []byte("{\"estimationcount\": \"3\"}")
			revealMsg.Msg, err = json.Marshal(reveal)
			//fmt.Printf("json message: %X\n", estimation)
			if err != nil {
				return err
			}

			//quite a long way to get a single value, however we can't directy access the keeper
			queryClient := types.NewQueryClient(cliCtx)
			params := &types.QueryCodeHashRequest{
				Codeid: 1,
			}
			res, err := queryClient.CodeHash(context.Background(), params)
			if err != nil {
				return err
			}

			var encryptedMsg []byte
			revealMsg.CodeHash = []byte(hex.EncodeToString(res.Codehash))
			encryptedMsg, err = wasmCtx.Encrypt(revealMsg.Serialize())
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgRevealEstimation(clientCtx.GetFromAddress().String(), uint64(argsItemID), encryptedMsg)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			fmt.Printf("sending msg: %X\n", revealMsg.Msg)
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

			itemID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgItemTransferable(clientCtx.GetFromAddress().String(), bool(transferBool), uint64(itemID))
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
			argsItemID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgItemShipping(clientCtx.GetFromAddress().String(), bool(shippingtrackingBool), uint64(argsItemID))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdItemResell() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "item-resell [itemid] [shippingcost] [discount] [localpickup] [shippingregion] [note] ",
		Short: "Resell an item",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsItemID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			argsShippingcost, _ := strconv.ParseInt(args[1], 10, 64)

			argsDiscount, _ := strconv.ParseInt(args[2], 10, 64)

			argsLocalpickup := string(args[3])

			argsShippingregion := strings.Split(args[4], ",")

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			note := args[5]

			msg := types.NewMsgItemResell(clientCtx.GetFromAddress().String(), uint64(argsItemID), int64(argsShippingcost), int64(argsDiscount), string(argsLocalpickup), []string(argsShippingregion), string(note))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
