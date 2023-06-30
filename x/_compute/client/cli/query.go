package cli

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/gogo/protobuf/proto"

	//"io/ioutil"
	"os"
	"strconv"

	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	flag "github.com/spf13/pflag"
	cosmwasmTypes "github.com/trstlabs/trst/go-cosmwasm/types"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"

	wasmUtils "github.com/trstlabs/trst/x/compute/client/utils"

	"github.com/trstlabs/trst/x/compute/internal/types"
)

func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the compute module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	queryCmd.AddCommand(
		GetCmdListCodes(),
		GetCmdListContractByCode(),
		GetCmdQueryCode(),
		GetCmdGetContractInfo(),
		GetCmdGetContractState(),
		GetCmdQuery(),
		GetQueryDecryptTxCmd(),
		GetCmdQueryContractID(),
		GetCmdCodeHashByContract(),
		CmdDecryptText(),
		// GetCmdGetContractHistory(cdc),
	)
	return queryCmd
}

// GetCmdListCodes lists all wasm code uploaded
func GetCmdListCodes() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-codes",
		Short: "List all wasm bytecode on the chain",
		Long:  "List all wasm bytecode on the chain",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryListCode)
			res, _, err := clientCtx.Query(route)
			if err != nil {
				return err
			}
			fmt.Println(string(res))
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryContractID checks if a contract-id is in use
func GetCmdQueryContractID() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contract-id [contract-id]",
		Short: "Check if a contract-id is in use",
		Long:  "Check if a contract-id is in use",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryContractAddress, args[0])
			res, _, err := clientCtx.Query(route)
			if err != nil {
				if err == sdkErrors.ErrUnknownAddress {
					fmt.Printf("ContractId is available and not in use\n")
					return nil
				}
				return fmt.Errorf("error querying: %s", err)
			}

			addr := sdk.AccAddress{}
			if addr == nil {
				fmt.Printf("no contract address found for this contract id: %s", args[0])
				return nil
			}
			err = addr.Unmarshal(res)
			if err != nil {
				return fmt.Errorf("error unwrapping address: %s", err)
			}
			fmt.Printf("ContractId exists")
			fmt.Printf("ContractId is in use by contract address: %s\n", addr.String())
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdCodeHashByContract gets the code hash given a contract address
func GetCmdCodeHashByContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "code-hash [address]",
		Short: "Return the code hash of a contract",
		Long:  "Return the code hash of a contract",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryContractHash, args[0])
			res, _, err := clientCtx.Query(route)
			if err != nil {
				return fmt.Errorf("error querying code hash: %s", err)
			}

			codeHash := hex.EncodeToString(res)
			fmt.Printf("0x%s", codeHash)
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdListContractByCode lists all wasm code uploaded for given code id
func GetCmdListContractByCode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-contracts-by-code [code_id]",
		Short: "List all wasm bytecode on the chain for a given code id",
		Long:  "List all wasm bytecode on the chain for a given code id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			codeID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s/%d", types.QuerierRoute, types.QueryListContractByCode, codeID)
			res, _, err := clientCtx.Query(route)
			if err != nil {
				return err
			}
			fmt.Println(string(res))
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryCode returns the bytecode for a given contract
func GetCmdQueryCode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "code [code_id] [output filename]",
		Short: "Downloads wasm bytecode for a given code id",
		Long:  "Downloads wasm bytecode for a given code id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			codeID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s/%d", types.QuerierRoute, types.QueryGetCode, codeID)
			res, _, err := clientCtx.Query(route)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				return fmt.Errorf("contract not found")
			}
			var code types.QueryCodeResponse
			err = json.Unmarshal(res, &code)
			if err != nil {
				return err
			}

			if len(code.Data) == 0 {
				return fmt.Errorf("contract not found")
			}

			fmt.Printf("Downloading wasm code to %s\n", args[1])
			fmt.Printf("This code has a contract duration of %s\n hours", code.DefaultDuration)
			return os.WriteFile(args[1], code.Data, 0644)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdGetContractInfo gets details about a given contract
func GetCmdGetContractInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contract [bech32_address]",
		Short: "Prints out metadata of a contract given its address",
		Long:  "Prints out metadata of a contract given its address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryGetContract, addr.String())
			res, _, err := clientCtx.Query(route)
			if err != nil {
				return err
			}
			fmt.Println(string(res))
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdGetContractState gets public state details about a given contract
func GetCmdGetContractState() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "state [bech32_address]",
		Short: "Prints out the last available state of a contract given its address",
		Long:  "Prints out the last available state of a contract given its address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryGetContractPublicState, addr.String())
			res, _, err := clientCtx.Query(route)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdDecryptText() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decrypt [encrypted_data]",
		Short: "Attempt to decrypt an encrypted blob",
		Long: "Attempt to decrypt a base-64 encoded encrypted message. This is intended to be used if manual decrypt" +
			"is required for data that is unavailable to be decrypted using the 'query compute tx' command",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			encodedInput := args[0]

			dataCipherBz, err := base64.StdEncoding.DecodeString(encodedInput)
			if err != nil {
				return fmt.Errorf("error while trying to decode the encrypted output data from base64: %w", err)
			}

			nonce := dataCipherBz[0:32]
			originalTxSenderPubkey := dataCipherBz[32:64]

			wasmCtx := wasmUtils.WASMContext{CLIContext: clientCtx}
			_, myPubkey, _ := wasmCtx.GetTxSenderKeyPair()

			if !bytes.Equal(originalTxSenderPubkey, myPubkey) {
				return fmt.Errorf("cannot decrypt, not original tx sender")
			}

			dataPlaintextB64Bz, err := wasmCtx.Decrypt(dataCipherBz[64:], nonce)
			if err != nil {
				return fmt.Errorf("error while trying to decrypt the output data: %w", err)
			}

			fmt.Printf("Decrypted data: %s", dataPlaintextB64Bz)
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// QueryDecryptTxCmd the default command for a tx query + IO decryption if I'm the tx sender.
// Coppied from https://github.com/cosmos/cosmos-sdk/blob/v0.38.4/x/auth/client/cli/query.go#L157-L184 and added IO decryption (Could not wrap it because it prints directly to stdout)
func GetQueryDecryptTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tx [hash]",
		Short: "Query for a transaction by hash in a committed block, decrypt input and outputs if I'm the tx sender",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			result, err := authtx.QueryTx(clientCtx, args[0])
			if err != nil {
				return err
			}

			if result.Empty() {
				return fmt.Errorf("no transaction found with hash %s", args[0])
			}

			var answer types.DecryptedAnswer
			var encryptedInput []byte
			var dataOutputHexB64 string

			txInputs := result.GetTx().GetMsgs()

			if len(txInputs) != 1 {
				return fmt.Errorf("can only decrypt txs with 1 input. Got %d", len(txInputs))
			}
			txInput, ok := txInputs[0].(*types.MsgExecuteContract)
			if !ok {
				txInput2, ok := txInputs[0].(*types.MsgInstantiateContract)
				if !ok {
					txInput3, ok := txInputs[0].(*types.MsgStoreCode)
					if ok {
						txInput3.WASMByteCode = nil
						return clientCtx.PrintProto(txInput3)
					} else {
						return fmt.Errorf("TX is not a compute transaction")
					}
				} else {
					encryptedInput = txInput2.Msg
					dataOutputHexB64 = result.Data
					answer.Type = "instantiate"
				}
			} else {
				encryptedInput = txInput.Msg
				dataOutputHexB64 = result.Data
				answer.Type = "execute"
			}

			// decrypt input
			if len(encryptedInput) < 64 {
				return fmt.Errorf("input must be > 64 bytes. Got %d", len(encryptedInput))
			}

			nonce := encryptedInput[0:32]
			originalTxSenderPubkey := encryptedInput[32:64]

			wasmCtx := wasmUtils.WASMContext{CLIContext: clientCtx}
			_, myPubkey, err := wasmCtx.GetTxSenderKeyPair()
			if err != nil {
				return fmt.Errorf("error in GetTxSenderKeyPair: %w", err)
			}

			if !bytes.Equal(originalTxSenderPubkey, myPubkey) {
				return fmt.Errorf("cannot decrypt, not original tx sender")
			}

			ciphertextInput := encryptedInput[64:]
			var plaintextInput []byte
			if len(ciphertextInput) > 0 {
				plaintextInput, err = wasmCtx.Decrypt(ciphertextInput, nonce)
				if err != nil {
					return fmt.Errorf("error while trying to decrypt the tx input: %w", err)
				}
			}

			answer.Input = string(plaintextInput)

			// decrypt data
			if answer.Type == "execute" {
				dataOutputAsProtobuf, err := hex.DecodeString(dataOutputHexB64)
				if err != nil {
					return fmt.Errorf("error while trying to decode the encrypted output data from hex string: %w", err)
				}

				var txData sdk.MsgData
				proto.Unmarshal(dataOutputAsProtobuf, &txData)

				dataOutputCipherBz, err := base64.StdEncoding.DecodeString(string(txData.Data))
				if err != nil {
					return fmt.Errorf("error while trying to decode the encrypted output data from base64 '%v': %w", string(txData.Data), err)
				}

				dataPlaintextB64Bz, err := wasmCtx.Decrypt(dataOutputCipherBz, nonce)
				if err != nil {
					return fmt.Errorf("error while trying to decrypt the output data: %w", err)
				}
				dataPlaintextB64 := string(dataPlaintextB64Bz)
				answer.OutputData = dataPlaintextB64

				dataPlaintext, err := base64.StdEncoding.DecodeString(dataPlaintextB64)
				if err != nil {
					return fmt.Errorf("error while trying to decode the decrypted output data from base64 '%v': %w", dataPlaintextB64, err)
				}

				answer.OutputDataAsString = string(dataPlaintext)
			}

			// decrypt logs
			answer.OutputLogs = []sdk.StringEvent{}
			for _, l := range result.Logs {
				for _, e := range l.Events {
					if e.Type == "wasm" {
						for i, a := range e.Attributes {
							if a.Key != "contract_address" {
								// key
								if a.Key != "" {
									// Try to decrypt the log key. If it doesn't look encrypted, leave it as-is
									keyCiphertext, err := base64.StdEncoding.DecodeString(a.Key)
									if err == nil {
										keyPlaintext, err := wasmCtx.Decrypt(keyCiphertext, nonce)
										if err == nil {
											a.Key = string(keyPlaintext)
										}
									}
								}

								// value
								if a.Value != "" {
									// Try to decrypt the log value. If it doesn't look encrypted, leave it as-is
									valueCiphertext, err := base64.StdEncoding.DecodeString(a.Value)
									if err == nil {
										valuePlaintext, err := wasmCtx.Decrypt(valueCiphertext, nonce)
										if err == nil {
											a.Value = string(valuePlaintext)
										}
									}
								}

								e.Attributes[i] = a
							}
						}
						answer.OutputLogs = append(answer.OutputLogs, e)
					}
				}
			}

			if types.IsEncryptedErrorCode(result.Code) && types.ContainsEncryptedString(result.RawLog) {
				stdErr, err := wasmCtx.DecryptError(result.RawLog, answer.Type, nonce)
				if err != nil {
					return err
				}

				answer.OutputError = stdErr
			} else if types.ContainsEnclaveError(result.RawLog) {
				answer.PlaintextError = result.RawLog
			}

			return clientCtx.PrintObjectLegacy(&answer)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetCmdQuery() *cobra.Command {
	decoder := newArgDecoder(asciiDecodeString)

	cmd := &cobra.Command{
		Use:   "query [contract address] [query]", // TODO add --from wallet
		Short: "Calls contract with given address with query data and prints the returned result",
		Long:  "Calls contract with given address with query data and prints the returned result",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			var msg string
			var contractAddr string

			if len(args) == 1 {

				contractId, err := cmd.Flags().GetString(flagContractId)
				if err != nil {
					return fmt.Errorf("Trustless Contract ID or bech32 contract address is required")
				}

				route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryContractAddress, contractId)
				res, _, err := clientCtx.Query(route)
				if err != nil {
					return err
				}

				contractAddr = string(res)
				msg = args[0]
			} else {
				// get the id of the code to instantiate
				contractAddr = args[0]
				msg = args[1]
			}

			queryData, err := decoder.DecodeString(msg)
			if err != nil {
				return fmt.Errorf("decode query: %s", err)
			}

			if !json.Valid(queryData) {
				return errors.New("query data must be json")
			}

			return QueryWithData(contractAddr, queryData, clientCtx)
		},
	}
	decoder.RegisterFlags(cmd.PersistentFlags(), "query argument")
	cmd.Flags().String(flagContractId, "", "A human-readable name for this contract in lists")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func QueryWithData(contractAddress string, queryData []byte, cliCtx client.Context) error {
	addr, err := sdk.AccAddressFromBech32(contractAddress)
	if err != nil {
		return err
	}

	route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryGetContractPrivateState, addr.String())

	wasmCtx := wasmUtils.WASMContext{CLIContext: cliCtx}

	/*codeHash, err := GetCodeHashByContractAddr(cliCtx, addr)
	if err != nil {
		return fmt.Errorf("contract not found: %s", addr)
	}*/

	/*	codeHash, err := GetCodeHashByCodeId(cliCtx, codeid)
		if err != nil {
			return fmt.Errorf("code id not found: %s", codeHash)
		}*/

	routeHash := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryContractHash, addr.String())
	hash, _, err := cliCtx.Query(routeHash)
	if err != nil {
		return fmt.Errorf("error querying code hash: %s", err)
	}

	//codeHash := hex.EncodeToString(res)
	msg := types.ContractMsg{
		CodeHash: []byte(hex.EncodeToString(hash)),
		Msg:      queryData,
	}

	queryData, err = wasmCtx.Encrypt(msg.Serialize())
	if err != nil {
		return fmt.Errorf("error encrypting contract data: %s", err)
	}
	nonce := queryData[:32]

	res, _, err := cliCtx.QueryWithData(route, queryData)

	//res, err := types.QueryContractState(sdk.Context,contractAddress, queryData, types)

	if err != nil {
		if types.ErrContainsQueryError(err) {
			errorPlainBz, err := wasmCtx.DecryptError(err.Error(), "query", nonce)
			if err != nil {
				return err
			}
			var stdErr cosmwasmTypes.StdError
			err = json.Unmarshal(errorPlainBz, &stdErr)
			if err != nil {
				return fmt.Errorf("error while trying to parse the error as json: '%s': %w", string(errorPlainBz), err)
			}
			return fmt.Errorf("query result: %v", stdErr.Error())
		}

		return fmt.Errorf("error querying contract data: %s", err)
	}

	var resDecrypted []byte
	if len(res) > 0 {
		resDecrypted, err = wasmCtx.Decrypt(res, nonce)
		if err != nil {
			return fmt.Errorf("error decrypting contract data: %s", err)
		}
	}

	decodedResp, err := base64.StdEncoding.DecodeString(string(resDecrypted))
	if err != nil {
		return fmt.Errorf("error decoding contract data: %s", err)
	}

	fmt.Println(string(decodedResp))
	return nil
}

type argumentDecoder struct {
	// dec is the default decoder
	dec                func(string) ([]byte, error)
	asciiF, hexF, b64F bool
}

func newArgDecoder(def func(string) ([]byte, error)) *argumentDecoder {
	return &argumentDecoder{dec: def}
}

func (a *argumentDecoder) RegisterFlags(f *flag.FlagSet, argName string) {
	f.BoolVar(&a.asciiF, "ascii", false, "ascii encoded "+argName)
	f.BoolVar(&a.hexF, "hex", false, "hex encoded  "+argName)
	f.BoolVar(&a.b64F, "b64", false, "base64 encoded "+argName)
}

func (a *argumentDecoder) DecodeString(s string) ([]byte, error) {
	found := -1
	for i, v := range []*bool{&a.asciiF, &a.hexF, &a.b64F} {
		if !*v {
			continue
		}
		if found != -1 {
			return nil, errors.New("multiple decoding flags used")
		}
		found = i
	}
	switch found {
	case 0:
		return asciiDecodeString(s)
	case 1:
		return hex.DecodeString(s)
	case 2:
		return base64.StdEncoding.DecodeString(s)
	default:
		return a.dec(s)
	}
}

func asciiDecodeString(s string) ([]byte, error) {
	return []byte(s), nil
}
