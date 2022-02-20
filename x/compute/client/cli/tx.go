package cli

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	//"io/ioutil"
	"os"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/tx"

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
	flagRunAs                  = "run-as"
	flagInstantiateByEverybody = "instantiate-everybody"
	flagInstantiateByAddress   = "instantiate-only-address"
	flagProposalType           = "type"
	flagIoMasterKey            = "enclave-key"
	flagCodeHash               = "code-hash"
	flagDuration               = "duration"

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
		StoreCodeCmd(),
		InstantiateContractCmd(),
		ExecuteContractCmd(),
		// Currently not supporting these commands
		//MigrateContractCmd(cdc),
		//UpdateContractAdminCmd(cdc),
		//ClearContractAdminCmd(cdc),
	)
	return txCmd
}

// StoreCodeCmd will upload code to be reused.
func StoreCodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "store [wasm file] [title] [description]  --source [source] ",
		Short: "Upload a wasm binary",
		Args:  cobra.ExactArgs(3),
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

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd.Flags().String(flagSource, "", "A valid URI reference to the contract's source code, optional")
	cmd.Flags().String(flagBuilder, "", "A valid docker tag for the build system, optional")
	//cmd.Flags().String(flagInstantiateByEverybody, "", "Everybody can instantiate a contract from the code, optional")
	//cmd.Flags().String(flagInstantiateByAddress, "", "Only this address can instantiate a contract instance from the code, optional")
	cmd.Flags().String(flagDuration, "", "A duration for the contract e.g. 2h, 6000s, 72h3m0.5s, optional")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func parseStoreCodeArgs(args []string, cliCtx client.Context, flags *flag.FlagSet) (types.MsgStoreCode, error) {
	wasm, err := os.ReadFile(args[0])
	if err != nil {
		return types.MsgStoreCode{}, err
	}

	argsTitle := string(args[1])
	argsDescription := string(args[2])

	// gzip the wasm file
	if wasmUtils.IsWasm(wasm) {
		wasm, err = wasmUtils.GzipIt(wasm)

		if err != nil {
			return types.MsgStoreCode{}, err
		}
	} else if !wasmUtils.IsGzip(wasm) {
		return types.MsgStoreCode{}, fmt.Errorf("invalid input file. Use wasm binary or gzip")
	}

	contractDuration, err := flags.GetString(flagDuration) //strconv.ParseInt(args[3], 10, 64)
	if err != nil {
		return types.MsgStoreCode{}, fmt.Errorf("contract duration: %s", err)
	}
	/*
	   var perm *types.AccessConfig
	   if onlyAddrStr := viper.GetString(flagInstantiateByAddress); onlyAddrStr != "" {
	       allowedAddr, err := sdk.AccAddressFromBech32(onlyAddrStr)
	       if err != nil {
	           return types.MsgStoreCode{}, sdkerrors.Wrap(err, flagInstantiateByAddress)
	       }
	       x := types.OnlyAddress.With(allowedAddr)
	       perm = &x
	   } else if everybody := viper.GetBool(flagInstantiateByEverybody); everybody {
	       perm = &types.AllowEverybody
	   }
	*/

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
		Sender:         cliCtx.GetFromAddress().String(),
		WASMByteCode:   wasm,
		Source:         source,
		Builder:        builder,
		ContractPeriod: contractDuration,
		Title:          argsTitle,
		Description:    argsDescription,
		// InstantiatePermission: perm,
	}
	return msg, nil
}

// InstantiateContractCmd will instantiate a contract from previously uploaded code.
func InstantiateContractCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instantiate [code id] [JSON args] --contract_id [unique contractId] " /* --admin [address,optional] */ + "--amount [coins] (optional)",
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
	var sendAutoMsg bool
	cmd.Flags().String(flagCodeHash, "", "For offline transactions, use this to specify the target contract's code hash")
	cmd.Flags().String(flagIoMasterKey, "", "For offline transactions, use this to specify the path to the "+
		"io-master-cert.der file, which you can get using the command `trstd q register trst-enclave-params` ")
	cmd.Flags().String(flagAmount, "", "Coins to send to the contract during instantiation")
	cmd.Flags().String(flagContractId, "", "A human-readable name for this contract in lists")
	cmd.Flags().BoolVar(&sendAutoMsg, "auto_msg", false, "(A auto message to send, before the contract ends (optional)")

	// cmd.Flags().String(flagAdmin, "", "Address of an admin")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func parseInstantiateArgs(args []string, cliCtx client.Context, initFlags *flag.FlagSet) (types.MsgInstantiateContract, error) {
	// get the id of the code to instantiate
	codeID, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return types.MsgInstantiateContract{}, err
	}

	amountStr, err := initFlags.GetString(flagAmount)
	if err != nil {
		return types.MsgInstantiateContract{}, fmt.Errorf("amount: %s", err)
	}

	amount, err := sdk.ParseCoinsNormalized(amountStr)
	if err != nil {
		return types.MsgInstantiateContract{}, err
	}

	contractId, err := initFlags.GetString(flagContractId)
	if contractId == "" {
		return types.MsgInstantiateContract{}, fmt.Errorf("contract id is required on all contracts")
	}
	sendAutoMsg, _ := initFlags.GetBool("auto_msg")

	if err != nil {
		return types.MsgInstantiateContract{}, err
	}

	wasmCtx := wasmUtils.WASMContext{CLIContext: cliCtx}
	initMsg := types.TrustlessMsg{}

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
		initMsg.CodeHash = []byte(codeHash)

		encryptedMsg, _ = wasmCtx.OfflineEncrypt(initMsg.Serialize(), ioKeyPath)
	} else {
		// if we aren't creating an offline transaction we can validate the chosen contractId
		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryContractAddress, contractId)
		res, _, _ := cliCtx.Query(route)
		if res != nil {
			return types.MsgInstantiateContract{}, fmt.Errorf("contractId already exists. You must choose a unique contractId for your contract instance")
		}

		initMsg.CodeHash, err = GetCodeHashByCodeId(cliCtx, args[0])
		if err != nil {
			return types.MsgInstantiateContract{}, err
		}

		//fmt.Printf("Got code hash: %X\n", initMsg.CodeHash)
		// todo: Add check that this is valid json and stuff
		initMsg.Msg = []byte(args[1])
		//initMsg.Msg = []byte("{\"estimationcount\": \"3\"}")
		//initMsg.Msg, err = json.Marshal(estimation)
		//fmt.Printf("json message: %X\n", estimation)

		//initMsg.Msg = []byte(initMsg.Msg)

		encryptedMsg, _ = wasmCtx.Encrypt(initMsg.Serialize())

	}
	var autoMsgEncrypted []byte
	autoMsgEncrypted = nil
	if sendAutoMsg {
		autoMsg := types.TrustlessMsg{}
		auto := types.ParseAuto{}
		autoMsg.Msg, err = json.Marshal(auto)
		if err != nil {
			return types.MsgInstantiateContract{}, err
		}

		autoMsg.CodeHash = initMsg.CodeHash
		autoMsgEncrypted, err = wasmCtx.Encrypt(autoMsg.Serialize())
		if err != nil {
			return types.MsgInstantiateContract{}, err
		}

	}

	// build and sign the transaction, then broadcast to Tendermint
	msg := types.MsgInstantiateContract{
		Sender:           cliCtx.GetFromAddress().String(),
		CallbackCodeHash: "",
		CodeID:           codeID,
		ContractId:       contractId,
		InitFunds:        amount,
		InitMsg:          encryptedMsg,
		AutoMsg:          autoMsgEncrypted,
	}
	return msg, nil
}

// ExecuteContractCmd will instantiate a contract from previously uploaded code.
func ExecuteContractCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execute [contract address] [json encoded send args] [codeId]",
		Short: "Execute a command on a wasm-based Trustless Contract",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			var contractAddr []byte
			var msg []byte
			var codeHash string
			var ioKeyPath string
			var codeId string

			genOnly, _ := cmd.Flags().GetBool(flags.FlagGenerateOnly)

			amountStr, _ := cmd.Flags().GetString(flagAmount)

			if len(args) == 2 {

				if genOnly {
					return fmt.Errorf("offline transactions must contain contract address")
				}

				contractId, err := cmd.Flags().GetString(flagContractId)
				if err != nil {
					return fmt.Errorf("error with contractId: %s", err)
				}
				if contractId == "" {
					return fmt.Errorf("contractId or bech32 contract address is required")
				}

				route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryContractAddress, contractId)
				res, _, err := cliCtx.Query(route)
				if err != nil {
					return err
				}

				contractAddr = res
				msg = []byte(args[0])
				codeId = args[1]
			} else {
				// get address to execute
				res, err := sdk.AccAddressFromBech32(args[0])
				if err != nil {
					return err
				}

				contractAddr = res
				msg = []byte(args[1])
				codeId = args[2]

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

			return ExecuteWithData(cmd, contractAddr, msg, amountStr, genOnly, ioKeyPath, codeHash, cliCtx, codeId)
		},
	}

	cmd.Flags().String(flagCodeHash, "", "For offline transactions, use this to specify the target contract's code hash")
	cmd.Flags().String(flagIoMasterKey, "", "For offline transactions, use this to specify the path to the "+
		"io-master-cert.der file, which you can get using the command `trstd q register trst-enclave-params` ")
	cmd.Flags().String(flagAmount, "", "Coins to send to the contract along with command")
	cmd.Flags().String(flagContractId, "", "A human-readable name for this contract in lists")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func ExecuteWithData(cmd *cobra.Command, contractAddress sdk.AccAddress, msg []byte, amount string, genOnly bool, ioMasterKeyPath string, codeHash string, cliCtx client.Context, codeId string) error {
	wasmCtx := wasmUtils.WASMContext{CLIContext: cliCtx}
	execMsg := types.TrustlessMsg{}

	execMsg.Msg = msg

	coins, err := sdk.ParseCoinsNormalized(amount)
	if err != nil {
		return err
	}
	//	fmt.Print("Executing msg...")
	var encryptedMsg []byte
	if genOnly {
		execMsg.CodeHash = []byte(codeHash)
		encryptedMsg, err = wasmCtx.OfflineEncrypt(execMsg.Serialize(), ioMasterKeyPath)
	} else {

		execMsg.CodeHash, err = GetCodeHashByCodeId(cliCtx, codeId)

		if err != nil {
			return err
		}

		encryptedMsg, err = wasmCtx.Encrypt(execMsg.Serialize())

	}
	if err != nil {
		return err
	}

	// build and sign the transaction, then broadcast to Tendermint
	msgExec := types.MsgExecuteContract{
		Sender:           cliCtx.GetFromAddress().String(),
		Contract:         contractAddress.String(),
		CallbackCodeHash: "",
		SentFunds:        coins,
		Msg:              encryptedMsg,
	}

	//	fmt.Printf("Execute message before is %s \n", string(encryptedMsg))

	return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), &msgExec)
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
		return nil, err
	}

	return []byte(hex.EncodeToString(codeResp.CodeHash)), nil
}

/*
func GetCodeHashByContractAddr(cliCtx client.Context, contractAddr sdk.AccAddress) ([]byte, error) {
	route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryContractHash, contractAddr.String())
	res, _, err := cliCtx.Query(route)
	if err != nil {
		return nil, err
	}

	return []byte(hex.EncodeToString(res)), nil
}
*/
