package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/trstlabs/intento/x/intent/types"
)

// GetTxCmd creates and returns the intent tx command
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
		getSubmitActionCmd(),
		getRegisterAccountAndSubmitActionCmd(),
		getUpdateActionCmd(),
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

			// TestVersion defines a reusable interchainaccounts version string for testing purposes
			version := string(icatypes.ModuleCdc.MustMarshalJSON(&icatypes.Metadata{
				Version:                icatypes.Version,
				ControllerConnectionId: viper.GetString(flagConnectionID),
				HostConnectionId:       viper.GetString(flagCounterpartyConnectionID),
				Encoding:               icatypes.EncodingProtobuf,
				TxType:                 icatypes.TxTypeSDKMultiMsg,
			}))

			msg := types.NewMsgRegisterAccount(
				clientCtx.GetFromAddress().String(),
				viper.GetString(flagConnectionID),
				version,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(flagCounterpartyConnectionID, "", "Counterparty Connection ID, from host chain")
	cmd.Flags().String(flagConnectionID, "", "Connection ID, an IBC ID from this chain to the host chain")
	_ = cmd.MarkFlagRequired(flagConnectionID)

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
					return errors.Wrap(err, "error unmarshalling sdk msg")
				}
			}

			msg, err := types.NewMsgSubmitTx(clientCtx.GetFromAddress().String(), txMsg, viper.GetString(flagConnectionID))
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(flagConnectionID, "", "Connection ID, an IBC ID from this chain to the host chain, optional")

	_ = cmd.MarkFlagRequired(flagConnectionID)

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func getSubmitActionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "submit-action [path/to/sdk_msg.json]",
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			cdc := codec.NewProtoCodec(clientCtx.InterfaceRegistry)

			var txMsgs []sdk.Msg
			for _, arg := range args {
				var txMsg sdk.Msg
				if err := cdc.UnmarshalInterfaceJSON([]byte(arg), &txMsg); err != nil {
					// Check for file path if JSON input is not provided
					var msgContents []byte

					// Check if arg is a valid file path
					if _, err := os.Stat(arg); err == nil {
						// arg is a file path, read the file contents
						msgContents, err = os.ReadFile(arg)
						if err != nil {
							return errors.Wrap(err, "failed to read file")
						}
					} else {
						// arg is not a valid file path, assume it is a JSON string
						msgContents = []byte(arg)
					}

					// Parse the JSON content
					if err := cdc.UnmarshalInterfaceJSON(msgContents, &txMsg); err != nil {
						return errors.Wrap(err, "error unmarshalling sdk msg")
					}
					txMsgs = append(txMsgs, txMsg)
				}
			}

			// Get execution configuration
			configuration := getExecutionConfiguration(cmd)
			funds := sdk.Coins{}
			amount := viper.GetString(flagFeeFunds)
			if amount != "" {
				funds, err = sdk.ParseCoinsNormalized(amount)
				if err != nil {
					return err
				}
			}

			msg, err := types.NewMsgSubmitAction(clientCtx.GetFromAddress().String(), viper.GetString(flagLabel), txMsgs, viper.GetString(flagConnectionID), viper.GetString(flagCounterpartyConnectionID), viper.GetString(flagDuration), viper.GetString(flagInterval), viper.GetUint64(flagStartAt), funds, configuration)
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(fsAction)
	cmd.Flags().String(flagDuration, "", "A custom duration for the Action e.g. 2h, 6000s, 72h3m0.5s")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func getRegisterAccountAndSubmitActionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "register-ica-and-submit-action [path/to/sdk_msg.json]",
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			cdc := codec.NewProtoCodec(clientCtx.InterfaceRegistry)

			var txMsgs []sdk.Msg
			for _, arg := range args {
				var txMsg sdk.Msg
				if err := cdc.UnmarshalInterfaceJSON([]byte(arg), &txMsg); err != nil {
					// Check for file path if JSON input is not provided
					var msgContents []byte

					// Check if arg is a valid file path
					if _, err := os.Stat(arg); err == nil {
						// arg is a file path, read the file contents
						msgContents, err = os.ReadFile(arg)
						if err != nil {
							return errors.Wrap(err, "failed to read file")
						}
					} else {
						// arg is not a valid file path, assume it is a JSON string
						msgContents = []byte(arg)
					}

					// Parse the JSON content
					if err := cdc.UnmarshalInterfaceJSON(msgContents, &txMsg); err != nil {
						return errors.Wrap(err, "error unmarshalling sdk msg")
					}
					txMsgs = append(txMsgs, txMsg)
				}
			}

			// Get execution configuration
			configuration := getExecutionConfiguration(cmd)
			funds := sdk.Coins{}
			amount := viper.GetString(flagFeeFunds)
			if amount != "" {
				funds, err = sdk.ParseCoinsNormalized(amount)
				if err != nil {
					return err
				}
			}
			// TestVersion defines a reusable interchainaccounts version string for testing purposes
			version := string(icatypes.ModuleCdc.MustMarshalJSON(&icatypes.Metadata{
				Version:                icatypes.Version,
				ControllerConnectionId: viper.GetString(flagConnectionID),
				HostConnectionId:       viper.GetString(flagCounterpartyConnectionID),
				Encoding:               icatypes.EncodingProtobuf,
				TxType:                 icatypes.TxTypeSDKMultiMsg,
			}))

			msg, err := types.NewMsgRegisterAccountAndSubmitAction(clientCtx.GetFromAddress().String(), viper.GetString(flagLabel), txMsgs, viper.GetString(flagConnectionID), viper.GetString(flagDuration), viper.GetString(flagInterval), viper.GetUint64(flagStartAt), funds, configuration, version)
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(fsAction)
	cmd.Flags().String(flagDuration, "", "A custom duration for the Action e.g. 2h, 6000s, 72h3m0.5s")
	cmd.Flags().String(flagCounterpartyConnectionID, "", "Counterparty Connection ID, from host chain, optional")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func getUpdateActionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "update-action [id] [path/to/sdk_msg.json, optional] ",
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			cdc := codec.NewProtoCodec(clientCtx.InterfaceRegistry)

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			var txMsgs []sdk.Msg
			for _, arg := range args {
				if arg == args[0] {
					continue
				}
				var txMsg sdk.Msg
				if err := cdc.UnmarshalInterfaceJSON([]byte(arg), &txMsg); err != nil {
					// Check for file path if JSON input is not provided
					var msgContents []byte

					// Check if arg is a valid file path
					if _, err := os.Stat(arg); err == nil {
						// arg is a file path, read the file contents
						msgContents, err = os.ReadFile(arg)
						if err != nil {
							return errors.Wrap(err, "failed to read file")
						}
					} else {
						// arg is not a valid file path, assume it is a JSON string
						msgContents = []byte(arg)
					}

					// Parse the JSON content
					if err := cdc.UnmarshalInterfaceJSON(msgContents, &txMsg); err != nil {
						return errors.Wrap(err, "error unmarshalling sdk msg")
					}
					txMsgs = append(txMsgs, txMsg)
				}
			}
			// Get execution configuration
			configuration := getExecutionConfiguration(cmd)
			funds := sdk.Coins{}
			amount := viper.GetString(flagFeeFunds)
			if amount != "" {
				funds, err = sdk.ParseCoinsNormalized(amount)
				if err != nil {
					return err
				}
			}
			msg, err := types.NewMsgUpdateAction(clientCtx.GetFromAddress().String(), id, viper.GetString(flagLabel), txMsgs, viper.GetString(flagConnectionID), viper.GetUint64(flagEndTime), viper.GetString(flagInterval), viper.GetUint64(flagStartAt), funds, configuration /* viper.GetUint64(flagRetries) */ /* , viper.GetString(flagVersion) */)
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(fsAction)
	cmd.Flags().String(flagEndTime, "", "A custom end-time in UNIX time, optional")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func getExecutionConfiguration(cmd *cobra.Command) *types.ExecutionConfiguration {
	updatingDisabled := viper.GetBool(flagUpdatingDisabled)
	saveMsgResponses := viper.GetBool(flagSaveMsgResponses)
	fallbackToOwnerBalance := viper.GetBool(flagFallbackToOwnerBalance)
	reregisterICAAfterTimeout := viper.GetBool(flagReregisterICAAfterTimeout)
	stopOnSuccess := viper.GetBool(flagStopOnSuccess)
	stopOnFailure := viper.GetBool(flagStopOnFailure)

	configuration := &types.ExecutionConfiguration{
		UpdatingDisabled:          updatingDisabled,
		SaveMsgResponses:          saveMsgResponses,
		StopOnSuccess:             stopOnSuccess,
		StopOnFailure:             stopOnFailure,
		FallbackToOwnerBalance:    fallbackToOwnerBalance,
		ReregisterICAAfterTimeout: reregisterICAAfterTimeout,
	}

	return configuration
}

// func parseIntSlice(flagName string) ([]int64, error) {
// 	stringSlice := viper.GetStringSlice(flagName)
// 	var result []int64

// 	for _, id := range stringSlice {
// 		txId, err := strconv.ParseInt(id, 10, 64)
// 		if err != nil {
// 			return nil, errors.Wrap(err, "invalid id, must be a number")
// 		}
// 		result = append(result, txId)
// 	}

// 	return result, nil
// }
