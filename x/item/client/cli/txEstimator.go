package cli

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	//"crypto/sha256"
	//"encoding/hex"
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	wasmUtils "github.com/trstlabs/trst/x/compute/client/utils"
	"github.com/trstlabs/trst/x/item/types"
)

const (
	flagEstimationOnly = "estimation-only"
	flagInterested     = "interested"
)

func CmdCreateEstimation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-estimation [estimation] [deposit amount] [comment] [itemid]",
		Short: "Creates a new estimation",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			//	argsMsg, _ := strconv.ParseInt(args[0], 10, 64)

			wasmCtx := wasmUtils.WASMContext{CLIContext: clientCtx}
			estimateMsg := types.TrustlessMsg{}

			estimation := map[string]string{"estimation_amount": args[0], "comment": args[2]}

			message := map[string]interface{}{"create_estimation": estimation}

			estimateMsg.Msg, err = json.Marshal(message)
			if err != nil {
				return err
			}

			//fmt.Printf("json message: %X\n", estimateMsg.Msg)
			//estimateMsg.Msg = []byte("{ amount:" + args[0] + ",comment:" + args[1] + "}")
			queryClient := types.NewQueryClient(clientCtx)
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
			//estimateMsg.CodeHash = res.Codehash

			estimateMsg.CodeHash = []byte(hex.EncodeToString(res.Codehash))
			//fmt.Printf("Got estimate .CodeHash hash: %X\n", estimateMsg.CodeHash)
			encryptedMsg, err = wasmCtx.Encrypt(estimateMsg.Serialize())
			if err != nil {
				return err
			}

			argsDeposit, _ := strconv.ParseInt(args[1], 10, 64)

			interested, _ := cmd.Flags().GetBool(flagInterested)

			//argsComment := string(args[3])
			argsItemID, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return err
			}
			msg := types.NewMsgCreateEstimation(clientCtx.GetFromAddress().String(), encryptedMsg, uint64(argsItemID), int64(argsDeposit), bool(interested))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			//fmt.Printf("sending msg: %X\n", estimateMsg.Msg)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().Bool(flagEstimationOnly, false, "for an item that is estimation-only (has no shipping or location set")
	cmd.Flags().Bool(flagInterested, false, "for an item that is estimation-only (has no shipping or location set")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdUpdateLike() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-estimator [itemid] [like]",
		Short: "Update a like",
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
		Use:   "delete-estimation [id] ",
		Short: "Delete a estimation by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			wasmCtx := wasmUtils.WASMContext{CLIContext: cliCtx}

			deletegMsg := types.TrustlessMsg{}
			delete := types.ParseDelete{}
			deletegMsg.Msg, err = json.Marshal(delete)
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
			deletegMsg.CodeHash = []byte(hex.EncodeToString(res.Codehash))
			encryptedMsg, err = wasmCtx.Encrypt(deletegMsg.Serialize())
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteEstimation(cliCtx.GetFromAddress().String(), uint64(id), encryptedMsg)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().Bool(flagEstimationOnly, false, "for an item that is estimation-only (has no shipping or location set")

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdFlagItem() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-flag [itemid]",
		Short: "create a new flag for item",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			itemid, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			wasmCtx := wasmUtils.WASMContext{CLIContext: cliCtx}

			flagMsg := types.TrustlessMsg{}
			flag := types.ParseFlag{}

			//initMsg.Msg = []byte("{\"estimationcount\": \"3\"}")
			flagMsg.Msg, err = json.Marshal(flag)
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
			flagMsg.CodeHash = []byte(hex.EncodeToString(res.Codehash))
			encryptedMsg, err = wasmCtx.Encrypt(flagMsg.Serialize())
			if err != nil {
				return err
			}

			msg := types.NewMsgFlagItem(cliCtx.GetFromAddress().String(), uint64(itemid), encryptedMsg)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			fmt.Printf("json message: %X\n", flag)
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().Bool(flagEstimationOnly, false, "for an item that is estimation-only (has no shipping or location set")

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
