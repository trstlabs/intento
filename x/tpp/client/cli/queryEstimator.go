package cli

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/danieljdd/trst/x/trst/types"
	"github.com/spf13/cobra"
)

func CmdListEstimator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-estimator",
		Short: "list all estimator",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllEstimatorRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.EstimatorAll(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowEstimator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-estimator [id]",
		Short: "shows a estimator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			Itemid, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			params := &types.QueryGetEstimatorRequest{
				Itemid: Itemid,
			}

			res, err := queryClient.Estimator(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
