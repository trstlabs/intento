package cli

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	//"io/ioutil"
	"os"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/tx"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	//"github.com/trstlabs/trst/x/compute/internal/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	wasmUtils "github.com/trstlabs/trst/x/compute/client/utils"
	"github.com/trstlabs/trst/x/compute/internal/types"
)

const (
	flagAmount                 = "amount"
	flagSource                 = "source"
	flagBuilder                = "builder"
	flagContractId             = "contract_id"
	flagCreator                = "creator"
	flagInstantiateByEverybody = "instantiate-everybody"
	flagInstantiateByAddress   = "instantiate-only-address"
	flagProposalType           = "type"
	flagIoMasterKey            = "enclave-key"
	flagCodeHash               = "code-hash"
	flagDuration               = "duration"
	flagInterval               = "interval"
	flagTitle                  = "contract-title"
	flagDescription            = "contract-description"
	flagAutoMsg                = "auto_msg"
	flagStartAt                = "start_at"

	// flagAdmin                  = "admin"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Compute transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		CmdStoreCode(),
		CmdInstantiateContract(),
		CmdExecuteContract(),
		CmdDiscardAutoMsg(),
		// Currently not supporting these commands
		//MigrateContractCmd(cdc),
		//UpdateContractAdminCmd(cdc),
		//ClearContractAdminCmd(cdc),
	)
	return txCmd
}

// CmdStoreCode will upload code to be reused.
func CmdStoreCode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "store [wasm file] --contract-title [text] --contract-description [text] --source [source] ",
		Short: "Upload a wasm binary",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg, err := parseStoreCodeArgs(args, clientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			fmt.Printf("CLI TX with duration: %s \n", msg.DefaultDuration)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd.Flags().String(flagSource, "", "A valid URI reference to the contract's source code, optional")
	cmd.Flags().String(flagBuilder, "", "A valid docker tag for the build system, optional")
	cmd.Flags().String(flagTitle, "", "Title of contract")
	cmd.Flags().String(flagDescription, "", "Description of contract")
	cmd.Flags().String(flagDuration, "", "A default duration for the contract e.g. 2h, 6000s, 72h3m0.5s, optional")
	cmd.Flags().String(flagInterval, "", "A default interval for the contract e.g. 2h, 6000s, 72h3m0.5s, optional")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func parseStoreCodeArgs(args []string, cliCtx client.Context, flags *flag.FlagSet) (types.MsgStoreCode, error) {
	wasm, err := os.ReadFile(args[0])
	if err != nil {
		return types.MsgStoreCode{}, err
	}

	argsTitle, err := flags.GetString(flagTitle)
	if err != nil {
		return types.MsgStoreCode{}, fmt.Errorf("title: %s", err)
	}
	argsDescription, err := flags.GetString(flagDescription)
	if err != nil {
		return types.MsgStoreCode{}, fmt.Errorf("description: %s", err)
	}

	// gzip the wasm file
	if wasmUtils.IsWasm(wasm) {
		wasm, err = wasmUtils.GzipIt(wasm)

		if err != nil {
			return types.MsgStoreCode{}, err
		}
	} else if !wasmUtils.IsGzip(wasm) {
		return types.MsgStoreCode{}, fmt.Errorf("invalid input file. Use wasm binary or gzip")
	}

	duration, err := flags.GetString(flagDuration) //strconv.ParseInt(args[3], 10, 64)
	if err != nil {
		return types.MsgStoreCode{}, fmt.Errorf("contract duration: %s", err)
	}
	interval, err := flags.GetString(flagInterval) //strconv.ParseInt(args[3], 10, 64)
	if err != nil {
		return types.MsgStoreCode{}, fmt.Errorf("contract interval: %s", err)
	}

	source, err := flags.GetString(flagSource)
	if err != nil {
		return types.MsgStoreCode{}, fmt.Errorf("source: %s", err)
	}
	builder, err := flags.GetString(flagBuilder)
	if err != nil {
		return types.MsgStoreCode{}, fmt.Errorf("builder: %s", err)
	}

	// build and sign the transaction, then broadcast to Tendermint
	msg := types.MsgStoreCode{
		Sender:          cliCtx.GetFromAddress().String(),
		WASMByteCode:    wasm,
		Source:          source,
		Builder:         builder,
		DefaultDuration: duration,
		DefaultInterval: interval,
		Title:           argsTitle,
		Description:     argsDescription,
	}
	return msg, nil
}

// CmdInstantiateContract will instantiate a contract from previously uploaded code.
func CmdInstantiateContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instantiate [code id] [JSON args] --contract_id [unique contract ID] " /* --admin [address,optional] */ + "--amount [coins] (optional)  --auto_msg [json args, optional] --duration [custom duration e.g. 400s/5h] (optional)  --interval [custom dration e.g. 400s/5h]  (optional) --start_at [UNIX time]",
		Short: "Instantiate a Trustless Contract",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg, err := parseInstantiateArgs(args, cliCtx, cmd.Flags())
			if err != nil {
				return err
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), &msg)

		},
	}
	//var sendAutoMsg bool
	cmd.Flags().String(flagCodeHash, "", "For offline transactions, use this to specify the target contract's code hash")
	cmd.Flags().String(flagIoMasterKey, "", "For offline transactions, use this to specify the path to the "+
		"io-master-cert.der file, which you can get using the command `trstd q registration node-enclave-params` ")
	cmd.Flags().String(flagAmount, "", "Coins to send to the contract during instantiation")
	cmd.Flags().String(flagContractId, "", "A human-readable name for this contract in lists")
	cmd.Flags().String(flagAutoMsg, "", "An automatic message to send, that the contract executes after a set duration (optional)")
	cmd.Flags().String(flagDuration, "", "A custom duration for the contract e.g. 2h, 6000s, 72h3m0.5s, optional")
	cmd.Flags().String(flagInterval, "", "A custom interval for the contract e.g. 2h, 6000s, 72h3m0.5s, optional")
	cmd.Flags().String(flagStartAt, "0", "A custom start time for the contract self-execution, in UNIX time")

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func parseInstantiateArgs(args []string, cliCtx client.Context, initFlags *flag.FlagSet) (types.MsgInstantiateContract, error) {
	// get the id of the code to instantiate
	codeID, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return types.MsgInstantiateContract{}, err
	}

	amount, err := initFlags.GetString(flagAmount)
	if err != nil {
		return types.MsgInstantiateContract{}, fmt.Errorf("amount: %s", err)
	}

	funds, err := sdk.ParseCoinsNormalized(amount)
	if err != nil {
		return types.MsgInstantiateContract{}, err
	}

	contractId, _ := initFlags.GetString(flagContractId)
	if contractId == "" {
		return types.MsgInstantiateContract{}, fmt.Errorf("contract id is required on all contracts")
	}

	autoMsgString, err := initFlags.GetString(flagAutoMsg)
	if err != nil {
		return types.MsgInstantiateContract{}, err
	}
	duration, err := initFlags.GetString(flagDuration) //strconv.ParseInt(args[3], 10, 64)
	if err != nil {
		return types.MsgInstantiateContract{}, fmt.Errorf("contract duration: %s", err)
	}
	interval, err := initFlags.GetString(flagInterval) //strconv.ParseInt(args[3], 10, 64)
	if err != nil {
		return types.MsgInstantiateContract{}, fmt.Errorf("contract interval: %s", err)
	}

	startAtStr, err := initFlags.GetString(flagStartAt) //strconv.ParseInt(args[3], 10, 64)
	if err != nil {
		return types.MsgInstantiateContract{}, fmt.Errorf("startAt string: %s", err)
	}

	startAt, err := strconv.ParseUint(startAtStr, 10, 64)
	if err != nil {
		return types.MsgInstantiateContract{}, fmt.Errorf("failed to parse start duration at: %s", err)
	}

	//sendAutoMsg, _ := initFlags.GetBool("auto_msg")

	if err != nil {
		return types.MsgInstantiateContract{}, err
	}

	wasmCtx := wasmUtils.WASMContext{CLIContext: cliCtx}
	msg := types.ContractMsg{}

	var encryptedMsg []byte
	genOnly, err := initFlags.GetBool(flags.FlagGenerateOnly)
	if err != nil && genOnly {
		// if we're creating an offline transaction we just need the path to the io master key
		ioKeyPath, err := initFlags.GetString(flagIoMasterKey)
		if err != nil {
			return types.MsgInstantiateContract{}, fmt.Errorf("ioKeyPath: %s", err)
		}
		if ioKeyPath == "" {
			return types.MsgInstantiateContract{}, fmt.Errorf("missing flag --%s. To create an offline transaction, you must specify path to the enclave key", flagIoMasterKey)
		}

		codeHash, err := initFlags.GetString(flagCodeHash)
		if err != nil {
			return types.MsgInstantiateContract{}, fmt.Errorf("codeHash: %s", err)
		}
		if codeHash == "" {
			return types.MsgInstantiateContract{}, fmt.Errorf("missing flag --%s. To create an offline transaction, you must set the target contract's code hash", flagCodeHash)
		}
		msg.CodeHash = []byte(codeHash)

		encryptedMsg, _ = wasmCtx.OfflineEncrypt(msg.Serialize(), ioKeyPath)
	} else {
		// if we aren't creating an offline transaction we can validate the chosen contractId
		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryContractAddress, contractId)
		res, _, _ := cliCtx.Query(route)
		if res != nil {
			return types.MsgInstantiateContract{}, fmt.Errorf("contractId already exists. You must choose a unique contractId for your contract instance")
		}

		msg.CodeHash, err = GetCodeHashByCodeId(cliCtx, args[0])
		if err != nil {
			return types.MsgInstantiateContract{}, err
		}

		msg.Msg = []byte(args[1])

		encryptedMsg, _ = wasmCtx.Encrypt(msg.Serialize())

	}
	var autoMsgEncrypted []byte
	autoMsgEncrypted = nil
	if autoMsgString != "" {
		autoMsg := types.ContractMsg{}

		autoMsg.Msg = []byte(autoMsgString)
		autoMsg.CodeHash = msg.CodeHash
		autoMsgEncrypted, err = wasmCtx.Encrypt(autoMsg.Serialize())
		if err != nil {
			return types.MsgInstantiateContract{}, err
		}
	}

	// build and sign the transaction, then broadcast to Tendermint
	msgInit := types.MsgInstantiateContract{
		Sender:          cliCtx.GetFromAddress().String(),
		CodeHash:        "",
		CodeID:          codeID,
		ContractId:      contractId,
		Funds:           funds,
		Msg:             encryptedMsg,
		AutoMsg:         autoMsgEncrypted,
		Duration:        duration,
		Interval:        interval,
		StartDurationAt: startAt,
	}
	return msgInit, nil
}

// CmdExecuteContract will execute a contract from previously instantiated code.
func CmdExecuteContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execute [contract address] [json encoded send args]",
		Short: "Execute a command on a wasm-based Trustless Contract",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			var contractAddr []byte
			var msg []byte
			var codeHash string
			var ioKeyPath string
			//var codeId string

			genOnly, _ := cmd.Flags().GetBool(flags.FlagGenerateOnly)

			amountStr, _ := cmd.Flags().GetString(flagAmount)

			if len(args) == 1 {

				if genOnly {
					return fmt.Errorf("offline transactions must contain contract address")
				}

				contractId, err := cmd.Flags().GetString(flagContractId)
				if err != nil {
					return fmt.Errorf("error with contractId: %s", err)
				}
				if contractId == "" {
					return fmt.Errorf("contract Id or bech32 contract address is required")
				}

				route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryContractAddress, contractId)
				res, _, err := cliCtx.Query(route)
				if err != nil {
					return err
				} else if res == nil {
					return fmt.Errorf("no contract address found")
				}
				contractAddr = res

				msg = []byte(args[0])

			} else {
				// get address to execute
				contractAddr, err = sdk.AccAddressFromBech32(args[0])
				if err != nil {
					return err
				}

				msg = []byte(args[1])

			}

			if genOnly {

				ioKeyPath, err = cmd.Flags().GetString(flagIoMasterKey)
				if err != nil {
					return fmt.Errorf("error with ioKeyPath: %s", err)
				}
				if ioKeyPath == "" {
					return fmt.Errorf("missing flag --%s. To create an offline transaction, you must specify path to the enclave key", flagIoMasterKey)
				}

				codeHash, err = cmd.Flags().GetString(flagCodeHash)
				if err != nil {
					return fmt.Errorf("error with codeHash: %s", err)
				}
				if codeHash == "" {
					return fmt.Errorf("missing flag --%s. To create an offline transaction, you must set the target contract's code hash", flagCodeHash)
				}
			}

			msgExec, err := parseExecuteArgs(cmd, contractAddr, msg, amountStr, genOnly, ioKeyPath, codeHash, cliCtx)
			if err != nil {
				return err
			}
			if err := msgExec.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), &msgExec)
		},
	}

	cmd.Flags().String(flagCodeHash, "", "For offline transactions, use this to specify the target contract's code hash")
	cmd.Flags().String(flagIoMasterKey, "", "For offline transactions, use this to specify the path to the "+
		"io-master-cert.der file, which you can get using the command `trstd q registration node-enclave-params` ")
	cmd.Flags().String(flagAmount, "", "Coins to send to the contract along with command")
	cmd.Flags().String(flagContractId, "", "A human-readable name for this contract in lists")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func parseExecuteArgs(cmd *cobra.Command, contractAddress sdk.AccAddress, msg []byte, amount string, genOnly bool, ioMasterKeyPath string, codeHash string, cliCtx client.Context) (types.MsgExecuteContract, error) {
	wasmCtx := wasmUtils.WASMContext{CLIContext: cliCtx}
	execMsg := types.ContractMsg{}

	execMsg.Msg = msg

	funds, err := sdk.ParseCoinsNormalized(amount)
	if err != nil {
		return types.MsgExecuteContract{}, err
	}
	//	fmt.Print("Executing msg...")
	var encryptedMsg []byte
	if genOnly {
		execMsg.CodeHash = []byte(codeHash)
		encryptedMsg, err = wasmCtx.OfflineEncrypt(execMsg.Serialize(), ioMasterKeyPath)
	} else {
		execMsg.CodeHash, err = GetCodeHashByContractAddr(cliCtx, contractAddress)
		if err != nil {
			return types.MsgExecuteContract{}, err
		}
		encryptedMsg, err = wasmCtx.Encrypt(execMsg.Serialize())
	}
	if err != nil {
		return types.MsgExecuteContract{}, err
	}

	// build and sign the transaction, then broadcast to Tendermint
	msgExec := types.MsgExecuteContract{
		Sender:   cliCtx.GetFromAddress().String(),
		Contract: contractAddress.String(),
		CodeHash: "",
		Funds:    funds,
		Msg:      encryptedMsg,
	}

	//	fmt.Printf("Execute message before is %s \n", string(encryptedMsg))
	return msgExec, nil

}

func GetCodeHashByCodeId(cliCtx client.Context, codeID string) ([]byte, error) {
	route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryGetCode, codeID)
	res, _, err := cliCtx.Query(route)
	if err != nil {
		return nil, err
	}

	var codeResp types.QueryCodeResponse

	err = json.Unmarshal(res, &codeResp)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrNotFound, "code not found")
	}

	return []byte(hex.EncodeToString(codeResp.CodeHash)), nil
}

func GetCodeHashByContractAddr(cliCtx client.Context, contractAddr sdk.AccAddress) ([]byte, error) {
	route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryContractHash, contractAddr.String())
	res, _, err := cliCtx.Query(route)
	if err != nil {
		return nil, err
	}

	return []byte(hex.EncodeToString(res)), nil
}

func CmdDiscardAutoMsg() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-auto-msg",
		Short: "Cancel the auto-message for an automated contract",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			// get address to execute
			contractAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgDiscardAutoMsg(
				clientCtx.GetFromAddress(),
				contractAddr,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
