package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/danieljdd/tpp/x/tpp/types"
	"github.com/spf13/cobra"
)

func CmdListItem() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-item",
		Short: "list all item",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllItemRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.ItemAll(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}


func CmdListInactiveItems() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-inactive-items",
		Short: "list all inactive items",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllInactiveItemsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.InactiveItemsAll(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}


func CmdSellerItems() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "seller-items [seller]",
		Short: "list all seller items",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QuerySellerItemsRequest{
				Seller: args[0],
			}

			res, err := queryClient.SellerItems(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}


func CmdShowItem() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-item [id]",
		Short: "shows a item",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetItemRequest{
				Id: args[0],
			}

			res, err := queryClient.Item(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
