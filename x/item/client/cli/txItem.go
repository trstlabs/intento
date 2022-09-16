package cli

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"strings"

	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/trstlabs/trst/x/item/types"

	wasmUtils "github.com/trstlabs/trst/x/compute/client/utils"
)

const (
	flagPhotos          = "photos"
	flagTokenURI        = "token_uri"
	flagShippingCost    = "shipping_cost"
	flagLocation        = "location"
	flagShippingRegion  = "shipping_region"
	flagDepositAmount   = "deposit_amount"
	flagEstimationCount = "estimation_count"
	flagCondition       = "condition"
)

func CmdCreateItem() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-item [title] [description] [tags] --deposit_amount number --estimation_count number",
		Short: "Creates a new Trustless Item",
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			argsTitle := string(args[0])
			argsDescription := string(args[1])
			argsTags := strings.Split(args[2], ",")

			flags := cmd.Flags()
			argsLocation, _ := flags.GetString(flagLocation)
			argsShippingCost, _ := flags.GetInt64(flagShippingCost)
			argsCondition, _ := flags.GetInt64(flagCondition)
			argsEstimationCount, err := flags.GetInt64(flagEstimationCount)
			if err != nil {
				return err
			}
			argsDepositAmount, err := flags.GetInt64(flagDepositAmount)
			if err != nil {
				return err
			}

			argsShippingRegion, _ := flags.GetStringSlice(flagShippingRegion)

			argsPhotos, err := flags.GetStringSlice(flagPhotos)

			if err != nil {
				return err
			}
			argsTokenURI, _ := flags.GetString(flagTokenURI)

			wasmCtx := wasmUtils.WASMContext{CLIContext: cliCtx}

			count := map[string]string{"estimation_count": strconv.FormatInt(argsEstimationCount, 10), "deposit_required": strconv.FormatInt(argsDepositAmount, 10)}

			msg := types.ContractMsg{}
			msg.Msg, err = json.Marshal(count)

			if err != nil {
				return err
			}
			var codeId uint64 = 1
			if argsLocation == "" && len(argsShippingRegion) == 0 {
				codeId = 2
			}
			queryClient := types.NewQueryClient(cliCtx)
			params := &types.QueryCodeHashRequest{
				Codeid: codeId,
			}
			res, err := queryClient.CodeHash(context.Background(), params)
			if err != nil {
				return err
			}

			var encryptedMsg []byte

			msg.CodeHash = []byte(hex.EncodeToString(res.Codehash))
			encryptedMsg, err = wasmCtx.Encrypt(msg.Serialize())
			if err != nil {
				return err
			}

			autoMsg := types.ContractMsg{}
			auto := types.ParseAuto{}
			autoMsg.Msg, err = json.Marshal(auto)
			if err != nil {
				return err
			}

			autoMsg.CodeHash = msg.CodeHash
			autoMsgEncrypted, err := wasmCtx.Encrypt(autoMsg.Serialize())
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txMsg := types.NewMsgCreateItem(clientCtx.GetFromAddress().String(), string(argsTitle), string(argsDescription), argsShippingCost, string(argsLocation), int64(argsEstimationCount), []string(argsTags), int64(argsCondition), []string(argsShippingRegion), int64(argsDepositAmount), encryptedMsg, autoMsgEncrypted, []string(argsPhotos), argsTokenURI)
			if err := txMsg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), txMsg)
		},
	}
	//Estimation related
	cmd.Flags().Int64(flagDepositAmount, 0, "deposit amount for the estimators to estimate the item (higher = more accurate, lower = faster)'  [number] ")
	cmd.Flags().Int64(flagEstimationCount, 0, "estimation count of estimators to estimate the item (higher = more accurate, lower = faster)'  [number] ")

	//Property related
	cmd.Flags().String(flagTokenURI, "", "token_uri of the item [string] (optional)")
	cmd.Flags().StringSlice(flagPhotos, []string{}, "photos of the item, max 9 [string array] (optional)")
	cmd.Flags().Int64(flagCondition, 0, "condition of the item if applicable (optional)")
	//Transfer related
	cmd.Flags().StringSlice(flagShippingRegion, []string{}, "shipping regions  of the item e.g. 'UK,NL,DE' [string array] (optional)")
	cmd.Flags().Int64(flagShippingCost, 0, "shipping_cost of the item [string] (optional)")
	cmd.Flags().String(flagLocation, "", "location location of the item [string] (optional)")

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdDeleteItem() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-item [item id] ",
		Short: "Delete an item by item id",
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

			revealMsg := types.ContractMsg{}
			reveal := types.ParseReveal{}

			//msg.Msg = []byte("{\"estimationcount\": \"3\"}")
			revealMsg.Msg, err = json.Marshal(reveal)
			//fmt.Printf("json message: %X\n", estimation)
			if err != nil {
				return err
			}

			//quite a long way to get a single value, however we can't directy access the keeper
			queryClient := types.NewQueryClient(cliCtx)
			estimationOnly, _ := cmd.Flags().GetBool(flagEstimationOnly)
			var codeId uint64 = 1
			if estimationOnly {
				codeId = 2
			}
			params := &types.QueryCodeHashRequest{
				Codeid: codeId,
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
			//	fmt.Printf("sending msg: %X\n", revealMsg.Msg)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().Bool(flagEstimationOnly, false, "for an item that is estimation-only (has no shipping or location set")

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdItemTransferable() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "item-transferable [itemid]",
		Short: "set item transferability",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			wasmCtx := wasmUtils.WASMContext{CLIContext: cliCtx}

			itemID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			//msgTransfer := map[string]string{"transferable": ""}
			//msg.Msg = []byte("{\"estimationcount\": \"3\"}")
			msg := types.ContractMsg{}
			init := types.ParseTransferable{}
			msg.Msg, err = json.Marshal(init)
			//fmt.Printf("json message: %X\n", estimation)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)
			params := &types.QueryCodeHashRequest{
				Codeid: 1,
			}
			res, err := queryClient.CodeHash(context.Background(), params)
			if err != nil {
				return err
			}

			var encryptedMsg []byte

			msg.CodeHash = []byte(hex.EncodeToString(res.Codehash))
			encryptedMsg, err = wasmCtx.Encrypt(msg.Serialize())
			if err != nil {
				return err
			}

			txMsg := types.NewMsgItemTransferable(clientCtx.GetFromAddress().String(), encryptedMsg, uint64(itemID))
			if err := txMsg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), txMsg)
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
		Use:   "item-resell [itemid] [shipping_cost] [discount] [location] [shipping_region] [note] ",
		Short: "Resell an item",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsItemID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			argsShippingCost, _ := strconv.ParseInt(args[1], 10, 64)

			argsDiscount, _ := strconv.ParseInt(args[2], 10, 64)

			argsLocation := string(args[3])

			argsShippingRegion := strings.Split(args[4], ",")

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			note := args[5]

			msg := types.NewMsgItemResell(clientCtx.GetFromAddress().String(), uint64(argsItemID), int64(argsShippingCost), int64(argsDiscount), string(argsLocation), []string(argsShippingRegion), string(note))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdTokenizeItem() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tokenize-item [item id] ",
		Short: "Tokenize an item by item id",
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

			msg := types.NewMsgTokenizeItem(clientCtx.GetFromAddress().String(), uint64(id))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdUnTokenizeItem() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "un-tokenize-item [item id] ",
		Short: "Un-Tokenize an item by item id",
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

			msg := types.NewMsgUnTokenizeItem(clientCtx.GetFromAddress().String(), uint64(id))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
