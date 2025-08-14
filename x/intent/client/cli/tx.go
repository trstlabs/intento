package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/trstlabs/intento/x/intent/types"
)

// GetTxCmd creates and returns the intent tx command
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transaction subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		getRegisterAccountCmd(),
		getSubmitTxCmd(),
		getSubmitFlowCmd(),
		getRegisterAccountAndSubmitFlowCmd(),
		getUpdateFlowCmd(),
		getCreateTrustlessAgent(),
		getUpdateTrustlessAgentCmd(),
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
				HostConnectionId:       viper.GetString(flagHostConnectionID),
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
	cmd.Flags().AddFlagSet(fsIBC)
	_ = cmd.MarkFlagRequired(flagConnectionID)
	_ = cmd.MarkFlagRequired(flagHostConnectionID)

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func getSubmitTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "submit-ica-tx [path/to/sdk_msg.json]",
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

	cmd.Flags().AddFlagSet(fsIBC)
	_ = cmd.MarkFlagRequired(flagConnectionID)
	_ = cmd.MarkFlagRequired(flagHostConnectionID)

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func getSubmitFlowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "submit-flow [path/to/sdk_msg.json]",
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

			conditions := types.ExecutionConditions{}
			conditionsString := viper.GetString(flagConditions)
			if conditionsString != "" {
				if err := json.Unmarshal([]byte(conditionsString), &conditions); err != nil {
					return errors.Wrap(err, "error unmarshalling conditions")
				}
			}

			// Get execution configuration
			configuration := getExecutionConfiguration()
			funds := sdk.Coins{}
			amount := viper.GetString(flagFeeFunds)
			if amount != "" {
				funds, err = sdk.ParseCoinsNormalized(amount)
				if err != nil {
					return err
				}
			}

			trustlessAgentFeeLimit := sdk.Coins{}
			trustlessAgentFeeLimitString := viper.GetString(flagTrustlessAgentFeeLimit)
			if trustlessAgentFeeLimitString != "" {
				trustlessAgentFeeLimit, err = sdk.ParseCoinsNormalized(trustlessAgentFeeLimitString)
				if err != nil {
					return err
				}
			}

			msg, err := types.NewMsgSubmitFlow(clientCtx.GetFromAddress().String(), viper.GetString(flagLabel), txMsgs, viper.GetString(flagConnectionID), viper.GetString(flagDuration), viper.GetString(flagInterval), viper.GetUint64(flagStartAt), funds, viper.GetString(flagTrustlessAgent), trustlessAgentFeeLimit, configuration, &conditions)
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(fsFlow)
	cmd.Flags().String(flagDuration, "", "A custom duration for the flow e.g. 2h, 6000s, 72h3m0.5s")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func getRegisterAccountAndSubmitFlowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "register-ica-and-submit-flow [path/to/sdk_msg.json]",
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
			configuration := getExecutionConfiguration()
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
				HostConnectionId:       viper.GetString(flagHostConnectionID),
				Encoding:               icatypes.EncodingProtobuf,
				TxType:                 icatypes.TxTypeSDKMultiMsg,
			}))

			msg, err := types.NewMsgRegisterAccountAndSubmitFlow(clientCtx.GetFromAddress().String(), viper.GetString(flagLabel), txMsgs, viper.GetString(flagConnectionID), viper.GetString(flagHostConnectionID), viper.GetString(flagDuration), viper.GetString(flagInterval), viper.GetUint64(flagStartAt), funds, configuration, version)
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(fsFlow)
	cmd.Flags().String(flagDuration, "", "A custom duration for the flow e.g. 2h, 6000s, 72h3m0.5s")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func getUpdateFlowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "update-flow [id] [path/to/sdk_msg.json, optional] ",
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

			conditions := types.ExecutionConditions{}
			conditionsString := viper.GetString(flagConditions)
			if conditionsString != "" {
				if err := json.Unmarshal([]byte(conditionsString), &conditions); err != nil {
					return errors.Wrap(err, "error unmarshalling conditions")
				}
			}

			// Get execution configuration
			configuration := getExecutionConfiguration()
			funds := sdk.Coins{}
			amount := viper.GetString(flagFeeFunds)
			if amount != "" {
				funds, err = sdk.ParseCoinsNormalized(amount)
				if err != nil {
					return err
				}
			}
			trustlessAgentFeeLimit := sdk.Coins{}
			trustlessAgentFeeLimitString := viper.GetString(flagTrustlessAgentFeeLimit)
			if trustlessAgentFeeLimitString != "" {
				trustlessAgentFeeLimit, err = sdk.ParseCoinsNormalized(trustlessAgentFeeLimitString)
				if err != nil {
					return err
				}
			}
			msg, err := types.NewMsgUpdateFlow(clientCtx.GetFromAddress().String(), id, viper.GetString(flagLabel), txMsgs, viper.GetString(flagConnectionID), viper.GetUint64(flagEndTime), viper.GetString(flagInterval), viper.GetUint64(flagStartAt), funds, viper.GetString(flagTrustlessAgent), trustlessAgentFeeLimit, configuration, &conditions)
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(fsFlow)
	cmd.Flags().String(flagEndTime, "", "A custom end-time in UNIX time, optional")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func getExecutionConfiguration() *types.ExecutionConfiguration {

	updatingDisabled := viper.GetBool(flagUpdatingDisabled)
	SaveResponses := viper.GetBool(flagSaveResponses)
	walletFallback := viper.GetBool(flagWalletFallback)
	stopOnSuccess := viper.GetBool(flagStopOnSuccess)
	stopOnFailure := viper.GetBool(flagStopOnFailure)
	stopOnTimeout := viper.GetBool(flagStopOnTimeout)
	configuration := types.ExecutionConfiguration{
		UpdatingDisabled: updatingDisabled,
		SaveResponses:    SaveResponses,
		StopOnSuccess:    stopOnSuccess,
		StopOnFailure:    stopOnFailure,
		StopOnTimeout:    stopOnTimeout,
		WalletFallback:   walletFallback,
	}

	return &configuration
}

func getCreateTrustlessAgent() *cobra.Command {
	cmd := &cobra.Command{
		Use: "create-trustless-agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			feeCoinsSupported := sdk.Coins{} //e.g. 54utrst,56uinto,57ucosm
			amount := viper.GetString(flagFeeCoinsSupported)
			if amount != "" {
				feeCoinsSupported, err = sdk.ParseCoinsNormalized(amount)
				if err != nil {
					return err
				}
			}

			// TestVersion defines a reusable interchainaccounts version string for testing purposes
			version := string(icatypes.ModuleCdc.MustMarshalJSON(&icatypes.Metadata{
				Version:                icatypes.Version,
				ControllerConnectionId: viper.GetString(flagConnectionID),
				HostConnectionId:       viper.GetString(flagHostConnectionID),
				Encoding:               icatypes.EncodingProtobuf,
				TxType:                 icatypes.TxTypeSDKMultiMsg,
			}))

			msg := types.NewMsgCreateTrustlessAgent(
				clientCtx.GetFromAddress().String(),
				viper.GetString(flagConnectionID),
				version,
				feeCoinsSupported,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(fsIBC)
	cmd.Flags().String(flagFeeCoinsSupported, "", "Coins supported as fees for hosted")
	// _ = cmd.MarkFlagRequired(flagFeeCoinsSupported)
	_ = cmd.MarkFlagRequired(flagConnectionID)
	_ = cmd.MarkFlagRequired(flagHostConnectionID)

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func getUpdateTrustlessAgentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "update-trustless-agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			feeCoinsSupported := sdk.Coins{} //e.g. 54utrst,56uinto,57ucosm
			amount := viper.GetString(flagFeeCoinsSupported)
			if amount != "" {
				feeCoinsSupported, err = sdk.ParseCoinsNormalized(amount)
				if err != nil {
					return err
				}
			}

			msg := types.NewMsgUpdateTrustlessAgent(
				clientCtx.GetFromAddress().String(),
				viper.GetString(flagTrustlessAgent),
				viper.GetString(flagNewAdmin),
				feeCoinsSupported,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	// cmd.Flags().String(flagFeeCoinsSupported, "", "Coins supported as fees for hosted, optional")
	cmd.Flags().String(flagNewAdmin, "", "A new admin, optional")
	cmd.Flags().String(flagTrustlessAgent, "", "A trustless agent to execute actions on a host")
	_ = cmd.MarkFlagRequired(flagTrustlessAgent)

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
