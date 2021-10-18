package cli

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/danieljdd/trst/x/trst/types"
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

func CmdListListedItems() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-listed-items",
		Short: "list all listed items",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllListedItemsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.ListedItemsAll(context.Background(), params)
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
			if len(args[0]) > 0 {
				params := &types.QuerySellerItemsRequest{
					Seller: args[0],
				}

				res, err := queryClient.SellerItems(context.Background(), params)
				if err != nil {
					return err
				}

				return clientCtx.PrintProto(res)
			} else {
				return sdkerrors.Wrap(types.ErrArgumentMissingOrNonUInteger, "address missing")
			}
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
			Id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			params := &types.QueryGetItemRequest{
				Id: Id,
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
