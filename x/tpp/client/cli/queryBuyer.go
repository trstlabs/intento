package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/danieljdd/tpp/x/tpp/types"
	"github.com/spf13/cobra"
)

func CmdListBuyer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-buyer",
		Short: "list all buyer",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllBuyerRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.BuyerAll(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowBuyer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-buyer [id]",
		Short: "shows a buyer",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetBuyerRequest{
				Itemid: args[0],
			}

			res, err := queryClient.Buyer(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
