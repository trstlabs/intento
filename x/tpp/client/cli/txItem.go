package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

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
		Use:   "create-item [title] [description] [shippingcost] [localpickup] [estimationcounthash] [tags] [condition] [shippingregion]",
		Short: "Creates a new item",
		Args:  cobra.ExactArgs(8),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsTitle := string(args[0])
			argsDescription := string(args[1])
			argsShippingcost, _ := strconv.ParseInt(args[2],10,64)
			argsLocalpickup := true
			if args[3] == "0"{
				argsLocalpickup = false
			}
			
			//argsEstimationcounthash := string(args[4])

			estimationcheck, ok := sdk.NewIntFromString(args[4])
			if ok != true {
				return sdkerrors.Wrap(types.ErrArgumentMissingOrNonUInteger, "not a number or lower than zero")
			}

			var estimationcount = fmt.Sprint(estimationcheck)
			var estimationcountHash = sha256.Sum256([]byte(estimationcount))
			var estimationcountHashString = hex.EncodeToString(estimationcountHash[:])


			argsTags := string(args[5])

			
		

			argsCondition, _ := strconv.ParseInt(args[6],10,64)

		
		
			argsShippingregion := string(args[7])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateItem(clientCtx.GetFromAddress().String(), string(argsTitle), string(argsDescription), int64(argsShippingcost), bool(argsLocalpickup), string(estimationcountHashString), string(argsTags), int64(argsCondition), string(argsShippingregion))
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
		Use:   "update-item [id] [title] [description] [shippingcost] [condition] [shippingregion]",
		Short: "Update a item",
		Args:  cobra.ExactArgs(7),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			argsTitle := string(args[1])
			argsDescription := string(args[2])
	
			argsShippingcost, _ := strconv.ParseInt(args[3],10,64)
			argsLocalpickup := true
			if args[4] == "0"{
				argsLocalpickup = false
			}
		

			argsCondition, _ := strconv.ParseInt(args[5],10,64)


			argsShippingregion := string(args[6])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateItem(clientCtx.GetFromAddress().String(), id, string(argsTitle), string(argsDescription), int64(argsShippingcost), bool(argsLocalpickup), int64(argsCondition), string(argsShippingregion))
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
