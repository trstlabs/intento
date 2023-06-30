package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/trstlabs/trst/x/item/types"
)

func CmdListProfiles() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-profiles",
		Short: "list all profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllProfilesRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.AllProfiles(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdGetProfile() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-profile [id]",
		Short: "gets an estimator profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			Owner := args[0]

			params := &types.QueryGetProfileRequest{
				Owner: Owner,
			}

			res, err := queryClient.Profile(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
