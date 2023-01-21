package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/trstlabs/trst/x/auto-ibc-tx/types"
)

// GetTxCmd creates and returns the auto-ibc-tx tx command
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		getRegisterAccountCmd(),
		getSubmitTxCmd(),
		getSubmitAutoTxCmd(),
		getRegisterAccountAndSubmitAutoTxCmd(),
	)

	return cmd
}

func getRegisterAccountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "register",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgRegisterAccount(
				clientCtx.GetFromAddress().String(),
				viper.GetString(FlagConnectionID),
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(fsConnectionPair)
	_ = cmd.MarkFlagRequired(FlagConnectionID)

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func getSubmitTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "submit-tx [path/to/sdk_msg.json]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			cdc := codec.NewProtoCodec(clientCtx.InterfaceRegistry)

			var txMsg sdk.Msg
			if err := cdc.UnmarshalInterfaceJSON([]byte(args[0]), &txMsg); err != nil {

				// check for file path if JSON input is not provided
				contents, err := os.ReadFile(args[0])
				if err != nil {
					return errors.Wrap(err, "neither JSON input nor path to .json file for sdk msg were provided")
				}

				if err := cdc.UnmarshalInterfaceJSON(contents, &txMsg); err != nil {
					return errors.Wrap(err, "error unmarshalling sdk msg file")
				}
			}

			msg, err := types.NewMsgSubmitTx(clientCtx.GetFromAddress().String(), txMsg, viper.GetString(FlagConnectionID))
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(fsConnectionPair)

	_ = cmd.MarkFlagRequired(FlagConnectionID)

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func getSubmitAutoTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "submit-auto-tx [path/to/sdk_msg.json]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			cdc := codec.NewProtoCodec(clientCtx.InterfaceRegistry)

			var txMsg sdk.Msg
			if err := cdc.UnmarshalInterfaceJSON([]byte(args[0]), &txMsg); err != nil {

				// check for file path if JSON input is not provided
				contents, err := os.ReadFile(args[0])
				if err != nil {
					return errors.Wrap(err, "neither JSON input nor path to .json file for sdk msg were provided")
				}

				if err := cdc.UnmarshalInterfaceJSON(contents, &txMsg); err != nil {
					return errors.Wrap(err, "error unmarshalling sdk msg file")
				}
			}
			var txIds []uint64
			dependsOnString := viper.GetStringSlice(flagDependsOn)
			for _, id := range dependsOnString {
				txId, err := strconv.ParseUint(id, 10, 64)
				if err != nil {
					return errors.Wrap(err, "invalid id, must be a number")
				}
				txIds = append(txIds, txId)
			}
			msg, err := types.NewMsgSubmitAutoTx(clientCtx.GetFromAddress().String(), txMsg, viper.GetString(FlagConnectionID), viper.GetString(flagDuration), viper.GetString(flagInterval), viper.GetUint64(flagStartAt), txIds, viper.GetUint64(flagRetries))
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(fsConnectionPair)
	cmd.Flags().String(flagDuration, "", "A custom duration for the contract e.g. 2h, 6000s, 72h3m0.5s, optional")
	cmd.Flags().String(flagInterval, "", "A custom interval for the contract e.g. 2h, 6000s, 72h3m0.5s, optional")
	cmd.Flags().Uint64(flagStartAt, 0, "A custom start time for the contract self-execution, in UNIX time")

	cmd.Flags().StringArray(flagDependsOn, []string{}, "array of auto-tx-ids this auto-tx depends on e.g. 5, 6")
	cmd.Flags().Uint64(flagRetries, 0, "Maximum amount of retries to make the tx succeed, optional")

	_ = cmd.MarkFlagRequired(FlagConnectionID)
	_ = cmd.MarkFlagRequired(flagDuration)

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func getRegisterAccountAndSubmitAutoTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "register-ica-and-submit-auto-tx [path/to/sdk_msg.json]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			cdc := codec.NewProtoCodec(clientCtx.InterfaceRegistry)

			var txMsg sdk.Msg
			if err := cdc.UnmarshalInterfaceJSON([]byte(args[0]), &txMsg); err != nil {

				// check for file path if JSON input is not provided
				contents, err := os.ReadFile(args[0])
				if err != nil {
					return errors.Wrap(err, "neither JSON input nor path to .json file for sdk msg were provided")
				}

				if err := cdc.UnmarshalInterfaceJSON(contents, &txMsg); err != nil {
					return errors.Wrap(err, "error unmarshalling sdk msg file")
				}
			}
			var txIds []uint64
			dependsOnString := viper.GetStringSlice(flagDependsOn)
			for _, id := range dependsOnString {
				txId, err := strconv.ParseUint(id, 10, 64)
				if err != nil {
					return errors.Wrap(err, "invalid id, must be a number")
				}
				txIds = append(txIds, txId)
			}
			msg, err := types.NewMsgRegisterAccountAndSubmitAutoTx(clientCtx.GetFromAddress().String(), txMsg, viper.GetString(FlagConnectionID), viper.GetString(flagDuration), viper.GetString(flagInterval), viper.GetUint64(flagStartAt), txIds, viper.GetUint64(flagRetries))
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(fsConnectionPair)
	cmd.Flags().String(flagDuration, "", "A custom duration for the contract e.g. 2h, 6000s, 72h3m0.5s, optional")
	cmd.Flags().String(flagInterval, "", "A custom interval for the contract e.g. 2h, 6000s, 72h3m0.5s, optional")
	cmd.Flags().String(flagStartAt, "0", "A custom start time for the contract self-execution, in UNIX time")

	cmd.Flags().StringArray(flagDependsOn, []string{}, "array of auto-tx-ids this auto-tx depends on e.g. 5, 6")
	cmd.Flags().Uint64(flagRetries, 0, "Maximum amount of retries to make the tx succeed, optional")

	_ = cmd.MarkFlagRequired(FlagConnectionID)
	_ = cmd.MarkFlagRequired(flagDuration)

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
