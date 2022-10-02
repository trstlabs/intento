package keeper

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"testing"

	crypto "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	wasmTypes "github.com/trstlabs/trst/go-cosmwasm/types"
	"github.com/trstlabs/trst/x/compute/internal/types"
)

const defaultGasForIbcTests = 600_000

func ibcChannelConnectHelper(
	t *testing.T, keeper Keeper, ctx sdk.Context,
	contractAddr sdk.AccAddress, creatorPrivKey crypto.PrivKey,
	gas uint64, shouldSendOpenAck bool, channel wasmTypes.IBCChannel,
) (sdk.Context, []ContractEvent, wasmTypes.StdError) {
	// create new ctx with the same storage and a gas limit
	// this is to reset the event manager, so we won't get
	// events from past calls
	gasMeter := &WasmCounterGasMeter{0, sdk.NewGasMeter(gas)}
	ctx = sdk.NewContext(
		ctx.MultiStore(),
		ctx.BlockHeader(),
		ctx.IsCheckTx(),
		log.NewNopLogger(),
	).WithGasMeter(gasMeter)

	var ibcChannelConnectMsg wasmTypes.IBCChannelConnectMsg
	if shouldSendOpenAck {
		ibcChannelConnectMsg = wasmTypes.IBCChannelConnectMsg{
			OpenAck: &wasmTypes.IBCOpenAck{
				Channel:             channel,
				CounterpartyVersion: "",
			},
			OpenConfirm: nil,
		}
	} else {
		ibcChannelConnectMsg = wasmTypes.IBCChannelConnectMsg{
			OpenAck: nil,
			OpenConfirm: &wasmTypes.IBCOpenConfirm{
				Channel: channel,
			},
		}
	}

	err := keeper.OnConnectChannel(ctx, contractAddr, ibcChannelConnectMsg)

	require.NotZero(t, gasMeter.GetWasmCounter(), err)

	if err != nil {
		return ctx, nil, wasmTypes.StdError{GenericErr: &wasmTypes.GenericErr{Msg: err.Error()}}
	}

	// wasmEvents comes from all the callbacks as well
	wasmEvents := tryDecryptWasmEvents(ctx, []byte{}, true)

	return ctx, wasmEvents, wasmTypes.StdError{}
}

func ibcChannelOpenHelper(
	t *testing.T, keeper Keeper, ctx sdk.Context,
	contractAddr sdk.AccAddress, creatorPrivKey crypto.PrivKey,
	gas uint64, shouldSendOpenTry bool, channel wasmTypes.IBCChannel,
) (string, wasmTypes.StdError) {
	// create new ctx with the same storage and a gas limit
	// this is to reset the event manager, so we won't get
	// events from past calls
	gasMeter := &WasmCounterGasMeter{0, sdk.NewGasMeter(gas)}
	ctx = sdk.NewContext(
		ctx.MultiStore(),
		ctx.BlockHeader(),
		ctx.IsCheckTx(),
		log.NewNopLogger(),
	).WithGasMeter(gasMeter)

	var ibcChannelOpenMsg wasmTypes.IBCChannelOpenMsg
	if shouldSendOpenTry {
		ibcChannelOpenMsg = wasmTypes.IBCChannelOpenMsg{
			OpenTry: &wasmTypes.IBCOpenTry{
				Channel:             channel,
				CounterpartyVersion: "",
			},
			OpenInit: nil,
		}
	} else {
		ibcChannelOpenMsg = wasmTypes.IBCChannelOpenMsg{
			OpenTry: nil,
			OpenInit: &wasmTypes.IBCOpenInit{
				Channel: channel,
			},
		}
	}

	res, err := keeper.OnOpenChannel(ctx, contractAddr, ibcChannelOpenMsg)
	fmt.Printf("ibc help err %+v", err)
	require.NotZero(t, gasMeter.GetWasmCounter(), err)

	if err != nil {
		return "", wasmTypes.StdError{GenericErr: &wasmTypes.GenericErr{Msg: err.Error()}}
	}

	return res, wasmTypes.StdError{}
}

func ibcChannelCloseHelper(
	t *testing.T, keeper Keeper, ctx sdk.Context,
	contractAddr sdk.AccAddress, creatorPrivKey crypto.PrivKey,
	gas uint64, shouldSendCloseConfirn bool, channel wasmTypes.IBCChannel,
) (sdk.Context, []ContractEvent, wasmTypes.StdError) {
	// create new ctx with the same storage and a gas limit
	// this is to reset the event manager, so we won't get
	// events from past calls
	gasMeter := &WasmCounterGasMeter{0, sdk.NewGasMeter(gas)}
	ctx = sdk.NewContext(
		ctx.MultiStore(),
		ctx.BlockHeader(),
		ctx.IsCheckTx(),
		log.NewNopLogger(),
	).WithGasMeter(gasMeter)

	var ibcChannelCloseMsg wasmTypes.IBCChannelCloseMsg
	if shouldSendCloseConfirn {
		ibcChannelCloseMsg = wasmTypes.IBCChannelCloseMsg{
			CloseConfirm: &wasmTypes.IBCCloseConfirm{
				Channel: channel,
			},
			CloseInit: nil,
		}
	} else {
		ibcChannelCloseMsg = wasmTypes.IBCChannelCloseMsg{
			CloseConfirm: nil,
			CloseInit: &wasmTypes.IBCCloseInit{
				Channel: channel,
			},
		}
	}

	err := keeper.OnCloseChannel(ctx, contractAddr, ibcChannelCloseMsg)

	require.NotZero(t, gasMeter.GetWasmCounter(), err)

	if err != nil {
		return ctx, nil, wasmTypes.StdError{GenericErr: &wasmTypes.GenericErr{Msg: err.Error()}}
	}

	// wasmEvents comes from all the callbacks as well
	wasmEvents := tryDecryptWasmEvents(ctx, []byte{}, true)

	return ctx, wasmEvents, wasmTypes.StdError{}
}

func createIBCEndpoint(port string, channel string) wasmTypes.IBCEndpoint {
	return wasmTypes.IBCEndpoint{
		PortID:    port,
		ChannelID: channel,
	}
}

func createIBCTimeout(timeout uint64) wasmTypes.IBCTimeout {
	return wasmTypes.IBCTimeout{
		Block:     nil,
		Timestamp: timeout,
	}
}

func createIBCPacket(src wasmTypes.IBCEndpoint, dest wasmTypes.IBCEndpoint, sequence uint64, timeout wasmTypes.IBCTimeout, data []byte) wasmTypes.IBCPacket {
	return wasmTypes.IBCPacket{
		Data:     data,
		Src:      src,
		Dest:     dest,
		Sequence: sequence,
		Timeout:  timeout,
	}
}

func ibcPacketReceiveHelper(
	t *testing.T, keeper Keeper, ctx sdk.Context,
	contractAddr sdk.AccAddress, creatorPrivKey crypto.PrivKey,
	shouldEncryptMsg bool, gas uint64, packet wasmTypes.IBCPacket,
) (sdk.Context, []byte, []ContractEvent, []byte, wasmTypes.StdError) {
	var nonce []byte
	internalPacket := packet

	if shouldEncryptMsg {
		contractHash, err := keeper.GetContractHash(ctx, contractAddr)
		require.NoError(t, err)
		hashStr := hex.EncodeToString(contractHash)

		msg := types.ContractMsg{
			CodeHash: []byte(hashStr),
			Msg:      packet.Data,
		}

		dataBz, err := wasmCtx.Encrypt(msg.Serialize())
		require.NoError(t, err)
		nonce = dataBz[0:32]
		internalPacket.Data = dataBz
	}

	// create new ctx with the same storage and a gas limit
	// this is to reset the event manager, so we won't get
	// events from past calls
	gasMeter := &WasmCounterGasMeter{0, sdk.NewGasMeter(gas)}
	ctx = sdk.NewContext(
		ctx.MultiStore(),
		ctx.BlockHeader(),
		ctx.IsCheckTx(),
		log.NewNopLogger(),
	).WithGasMeter(gasMeter)

	ibcPacketReceiveMsg := wasmTypes.IBCPacketReceiveMsg{
		Packet:  internalPacket,
		Relayer: "relayer",
	}

	res, err := keeper.OnRecvPacket(ctx, contractAddr, ibcPacketReceiveMsg)

	require.NotZero(t, gasMeter.GetWasmCounter(), err)

	// wasmEvents comes from all the callbacks as well
	wasmEvents := tryDecryptWasmEvents(ctx, nonce, !shouldEncryptMsg)

	if err != nil {
		if shouldEncryptMsg {
			return ctx, nil, nil, nil, extractInnerError(t, err, nonce, true)
		}

		return ctx, nil, nil, nil, wasmTypes.StdError{GenericErr: &wasmTypes.GenericErr{Msg: err.Error()}}
	}

	data := res
	if shouldEncryptMsg {
		data = getDecryptedData(t, res, nonce)
	}

	return ctx, nonce, wasmEvents, data, wasmTypes.StdError{}
}

func ibcPacketAckHelper(
	t *testing.T, keeper Keeper, ctx sdk.Context,
	contractAddr sdk.AccAddress, creatorPrivKey crypto.PrivKey, gas uint64, originalPacket wasmTypes.IBCPacket, ack []byte,
) (sdk.Context, []ContractEvent, wasmTypes.StdError) {
	// create new ctx with the same storage and a gas limit
	// this is to reset the event manager, so we won't get
	// events from past calls
	gasMeter := &WasmCounterGasMeter{0, sdk.NewGasMeter(gas)}
	ctx = sdk.NewContext(
		ctx.MultiStore(),
		ctx.BlockHeader(),
		ctx.IsCheckTx(),
		log.NewNopLogger(),
	).WithGasMeter(gasMeter)

	ibcPacketAckMsg := wasmTypes.IBCPacketAckMsg{
		Acknowledgement: wasmTypes.IBCAcknowledgement{
			Data: ack,
		},
		OriginalPacket: originalPacket,
		Relayer:        "relayer",
	}

	err := keeper.OnAckPacket(ctx, contractAddr, ibcPacketAckMsg)

	require.NotZero(t, gasMeter.GetWasmCounter(), err)

	if err != nil {
		return ctx, nil, wasmTypes.StdError{GenericErr: &wasmTypes.GenericErr{Msg: err.Error()}}
	}

	// wasmEvents comes from all the callbacks as well
	wasmEvents := tryDecryptWasmEvents(ctx, []byte{}, true)

	return ctx, wasmEvents, wasmTypes.StdError{}
}

func ibcPacketTimeoutHelper(
	t *testing.T, keeper Keeper, ctx sdk.Context,
	contractAddr sdk.AccAddress, creatorPrivKey crypto.PrivKey, gas uint64, originalPacket wasmTypes.IBCPacket,
) (sdk.Context, []ContractEvent, wasmTypes.StdError) {
	// create new ctx with the same storage and a gas limit
	// this is to reset the event manager, so we won't get
	// events from past calls
	gasMeter := &WasmCounterGasMeter{0, sdk.NewGasMeter(gas)}
	ctx = sdk.NewContext(
		ctx.MultiStore(),
		ctx.BlockHeader(),
		ctx.IsCheckTx(),
		log.NewNopLogger(),
	).WithGasMeter(gasMeter)

	ibcPacketTimeoutMsg := wasmTypes.IBCPacketTimeoutMsg{
		Packet:  originalPacket,
		Relayer: "relayer",
	}

	err := keeper.OnTimeoutPacket(ctx, contractAddr, ibcPacketTimeoutMsg)

	require.NotZero(t, gasMeter.GetWasmCounter(), err)

	if err != nil {
		return ctx, nil, wasmTypes.StdError{GenericErr: &wasmTypes.GenericErr{Msg: err.Error()}}
	}

	// wasmEvents comes from all the callbacks as well
	wasmEvents := tryDecryptWasmEvents(ctx, []byte{}, true)

	return ctx, wasmEvents, wasmTypes.StdError{}
}

func TestIBCChannelOpen(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, TestContractPaths[ibcContract], sdk.NewCoins())

	_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"init":{}}`, true, defaultGasForTests)
	require.Empty(t, err)

	ibcChannel := wasmTypes.IBCChannel{
		Endpoint:             createIBCEndpoint(PortIDForContract(contractAddress), "channel.0"),
		CounterpartyEndpoint: createIBCEndpoint(PortIDForContract(contractAddress), "channel.1"),
		Order:                wasmTypes.Unordered,
		Version:              "1",
		ConnectionID:         "1",
	}
	version, err := ibcChannelOpenHelper(t, keeper, ctx, contractAddress, privKeyA, defaultGasForTests, false, ibcChannel)
	require.Empty(t, err)
	require.Equal(t, version, "ibc-v1")

	queryRes, err := queryHelper(t, keeper, ctx, contractAddress, `{"q":{}}`, true, math.MaxUint64)
	require.Empty(t, err)

	require.Equal(t, "1", queryRes)
}

func TestIBCChannelOpenTry(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, TestContractPaths[ibcContract], sdk.NewCoins())

	_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"init":{}}`, true, defaultGasForTests)
	require.Empty(t, err)

	ibcChannel := wasmTypes.IBCChannel{
		Endpoint:             createIBCEndpoint(PortIDForContract(contractAddress), "channel.0"),
		CounterpartyEndpoint: createIBCEndpoint(PortIDForContract(contractAddress), "channel.1"),
		Order:                wasmTypes.Unordered,
		Version:              "1",
		ConnectionID:         "1",
	}

	version, err := ibcChannelOpenHelper(t, keeper, ctx, contractAddress, privKeyA, defaultGasForTests, true, ibcChannel)
	require.Empty(t, err)
	require.Equal(t, version, "ibc-v1")

	queryRes, err := queryHelper(t, keeper, ctx, contractAddress, `{"q":{}}`, true, math.MaxUint64)
	require.Empty(t, err)

	require.Equal(t, "2", queryRes)
}

func TestIBCChannelConnect(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, TestContractPaths[ibcContract], sdk.NewCoins())

	_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"init":{}}`, true, defaultGasForTests)
	require.Empty(t, err)

	for _, test := range []struct {
		description   string
		connectionID  string
		output        string
		isSuccess     bool
		hasAttributes bool
		hasEvents     bool
	}{
		{
			description:   "Default",
			connectionID:  "0",
			output:        "4",
			isSuccess:     true,
			hasAttributes: false,
			hasEvents:     false,
		},
		{
			description:   "SubmessageNoReply",
			connectionID:  "1",
			output:        "10",
			isSuccess:     true,
			hasAttributes: false,
			hasEvents:     false,
		},
		{
			description:   "SubmessageWithReply",
			connectionID:  "2",
			output:        "17",
			isSuccess:     true,
			hasAttributes: false,
			hasEvents:     false,
		},
		{
			description:   "Attributes",
			connectionID:  "3",
			output:        "7",
			isSuccess:     true,
			hasAttributes: true,
			hasEvents:     false,
		},
		{
			description:   "Events",
			connectionID:  "4",
			output:        "8",
			isSuccess:     true,
			hasAttributes: false,
			hasEvents:     true,
		},
		{
			description:   "Error",
			connectionID:  "5",
			output:        "",
			isSuccess:     false,
			hasAttributes: false,
			hasEvents:     false,
		},
	} {
		t.Run(test.description, func(t *testing.T) {
			ibcChannel := wasmTypes.IBCChannel{
				Endpoint:             createIBCEndpoint(PortIDForContract(contractAddress), "channel.0"),
				CounterpartyEndpoint: createIBCEndpoint(PortIDForContract(contractAddress), "channel.1"),
				Order:                wasmTypes.Unordered,
				Version:              "1",
				ConnectionID:         test.connectionID,
			}

			ctx, events, err := ibcChannelConnectHelper(t, keeper, ctx, contractAddress, privKeyA, defaultGasForIbcTests, false, ibcChannel)

			if !test.isSuccess {
				require.Contains(t, fmt.Sprintf("%+v", err), "Intentional")
			} else {
				require.Empty(t, err)
				if test.hasAttributes {
					require.Equal(t,
						[]ContractEvent{
							{
								{Key: "contract_address", Value: []byte(contractAddress.String())},
								{Key: "attr1", Value: []byte("😗")},
							},
						},
						events,
					)
				}

				if test.hasEvents {
					hadCyber1 := false
					evts := ctx.EventManager().Events()
					for _, e := range evts {
						if e.Type == "wasm-cyber1" {
							require.False(t, hadCyber1)
							attrs, err := parseAndDecryptAttributes(e.Attributes, []byte{})
							require.Empty(t, err)

							require.Equal(t,
								[]wasmTypes.Attribute{
									{Key: "contract_address", Value: []byte(contractAddress.String())},
									{Key: "attr1", Value: []byte("🤯")},
								},
								attrs,
							)

							hadCyber1 = true
						}
					}

					require.True(t, hadCyber1)
				}

				queryRes, err := queryHelper(t, keeper, ctx, contractAddress, `{"q":{}}`, true, math.MaxUint64)

				require.Empty(t, err)

				require.Equal(t, test.output, queryRes)
			}
		})
	}
}

func TestIBCChannelConnectOpenAck(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, TestContractPaths[ibcContract], sdk.NewCoins())

	_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"init":{}}`, true, defaultGasForTests)
	require.Empty(t, err)

	ibcChannel := wasmTypes.IBCChannel{
		Endpoint:             createIBCEndpoint(PortIDForContract(contractAddress), "channel.0"),
		CounterpartyEndpoint: createIBCEndpoint(PortIDForContract(contractAddress), "channel.1"),
		Order:                wasmTypes.Unordered,
		Version:              "1",
		ConnectionID:         "1",
	}

	ctx, _, err = ibcChannelConnectHelper(t, keeper, ctx, contractAddress, privKeyA, defaultGasForTests, true, ibcChannel)
	require.Empty(t, err)

	queryRes, err := queryHelper(t, keeper, ctx, contractAddress, `{"q":{}}`, true, math.MaxUint64)
	require.Empty(t, err)

	require.Equal(t, "3", queryRes)
}

func TestIBCChannelClose(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, TestContractPaths[ibcContract], sdk.NewCoins())

	_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"init":{}}`, true, defaultGasForIbcTests)
	require.Empty(t, err)

	for _, test := range []struct {
		description   string
		connectionID  string
		output        string
		isSuccess     bool
		hasAttributes bool
		hasEvents     bool
	}{
		{
			description:   "Default",
			connectionID:  "0",
			output:        "6",
			isSuccess:     true,
			hasAttributes: false,
			hasEvents:     false,
		},
		{
			description:   "SubmessageNoReply",
			connectionID:  "1",
			output:        "12",
			isSuccess:     true,
			hasAttributes: false,
			hasEvents:     false,
		},
		{
			description:   "SubmessageWithReply",
			connectionID:  "2",
			output:        "19",
			isSuccess:     true,
			hasAttributes: false,
			hasEvents:     false,
		},
		{
			description:   "Attributes",
			connectionID:  "3",
			output:        "9",
			isSuccess:     true,
			hasAttributes: true,
			hasEvents:     false,
		},
		{
			description:   "Events",
			connectionID:  "4",
			output:        "10",
			isSuccess:     true,
			hasAttributes: false,
			hasEvents:     true,
		},
		{
			description:   "Error",
			connectionID:  "5",
			output:        "",
			isSuccess:     false,
			hasAttributes: false,
			hasEvents:     false,
		},
	} {
		t.Run(test.description, func(t *testing.T) {
			ibcChannel := wasmTypes.IBCChannel{
				Endpoint:             createIBCEndpoint(PortIDForContract(contractAddress), "channel.0"),
				CounterpartyEndpoint: createIBCEndpoint(PortIDForContract(contractAddress), "channel.1"),
				Order:                wasmTypes.Unordered,
				Version:              "1",
				ConnectionID:         test.connectionID,
			}

			ctx, events, err := ibcChannelCloseHelper(t, keeper, ctx, contractAddress, privKeyA, defaultGasForIbcTests, true, ibcChannel)

			if !test.isSuccess {
				require.Contains(t, fmt.Sprintf("%+v", err), "Intentional")
			} else {
				require.Empty(t, err)
				if test.hasAttributes {
					require.Equal(t,
						[]ContractEvent{
							{
								{Key: "contract_address", Value: []byte(contractAddress.String())},
								{Key: "attr1", Value: []byte("😗")},
							},
						},
						events,
					)
				}

				if test.hasEvents {
					hadCyber1 := false
					evts := ctx.EventManager().Events()
					for _, e := range evts {
						if e.Type == "wasm-cyber1" {
							require.False(t, hadCyber1)
							attrs, err := parseAndDecryptAttributes(e.Attributes, []byte{})
							require.Empty(t, err)

							require.Equal(t,
								[]wasmTypes.Attribute{
									{Key: "contract_address", Value: []byte(contractAddress.String())},
									{Key: "attr1", Value: []byte("🤯")},
								},
								attrs,
							)

							hadCyber1 = true
						}
					}

					require.True(t, hadCyber1)
				}

				queryRes, err := queryHelper(t, keeper, ctx, contractAddress, `{"q":{}}`, true, math.MaxUint64)

				require.Empty(t, err)

				require.Equal(t, test.output, queryRes)
			}
		})
	}
}

func TestIBCChannelCloseInit(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, TestContractPaths[ibcContract], sdk.NewCoins())

	_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"init":{}}`, true, defaultGasForTests)
	require.Empty(t, err)

	ibcChannel := wasmTypes.IBCChannel{
		Endpoint:             createIBCEndpoint(PortIDForContract(contractAddress), "channel.0"),
		CounterpartyEndpoint: createIBCEndpoint(PortIDForContract(contractAddress), "channel.1"),
		Order:                wasmTypes.Unordered,
		Version:              "1",
		ConnectionID:         "1",
	}

	ctx, _, err = ibcChannelCloseHelper(t, keeper, ctx, contractAddress, privKeyA, defaultGasForTests, false, ibcChannel)
	require.Empty(t, err)

	queryRes, err := queryHelper(t, keeper, ctx, contractAddress, `{"q":{}}`, true, math.MaxUint64)
	require.Empty(t, err)

	require.Equal(t, "5", queryRes)
}

func TestIBCPacketReceive(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, TestContractPaths[ibcContract], sdk.NewCoins())

	_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"init":{}}`, true, defaultGasForTests)
	require.Empty(t, err)
	for _, isEncrypted := range []bool{false, true} {
		for _, test := range []struct {
			description   string
			sequence      uint64
			output        string
			isSuccess     bool
			hasAttributes bool
			hasEvents     bool
		}{
			{
				description:   "Default",
				sequence:      0,
				output:        "7",
				isSuccess:     true,
				hasAttributes: false,
				hasEvents:     false,
			},
			/*	{
					description:   "SubmessageNoReply",
					sequence:      1,
					output:        "13",
					isSuccess:     true,
					hasAttributes: false,
					hasEvents:     false,
				},
				{
					description:   "SubmessageWithReply",
					sequence:      2,
					output:        "20",
					isSuccess:     true,
					hasAttributes: false,
					hasEvents:     false,
				},
				{
					description:   "Attributes",
					sequence:      3,
					output:        "10",
					isSuccess:     true,
					hasAttributes: true,
					hasEvents:     false,
				},
				{
					description:   "Events",
					sequence:      4,
					output:        "11",
					isSuccess:     true,
					hasAttributes: false,
					hasEvents:     true,
				},
				{
					description:   "Error",
					sequence:      5,
					output:        "",
					isSuccess:     false,
					hasAttributes: false,
					hasEvents:     false,
				},
				{
					description:   "SubmessageWithReplyThatCallsToSubmessage",
					sequence:      6,
					output:        "35",
					isSuccess:     true,
					hasAttributes: false,
					hasEvents:     false,
				},*/
		} {
			t.Run(fmt.Sprintf("%s-Encryption:%t", test.description, isEncrypted), func(t *testing.T) {
				ibcPacket := createIBCPacket(createIBCEndpoint(PortIDForContract(contractAddress), "channel.1"),
					createIBCEndpoint(PortIDForContract(contractAddress), "channel.0"),
					test.sequence,
					createIBCTimeout(math.MaxUint64),
					[]byte{},
				)
				ctx, nonce, events, data, err := ibcPacketReceiveHelper(t, keeper, ctx, contractAddress, privKeyA, isEncrypted, defaultGasForIbcTests, ibcPacket)

				if !test.isSuccess {
					require.Contains(t, fmt.Sprintf("%+v", err), "Intentional")
				} else {
					require.Empty(t, err)
					require.Equal(t, "\"out\"", string(data))

					if test.hasAttributes {
						require.Equal(t,
							[]ContractEvent{
								{
									{Key: "contract_address", Value: []byte(contractAddress.String())},
									{Key: "attr1", Value: []byte("😗")},
								},
							},
							events,
						)
					}

					if test.hasEvents {
						hadCyber1 := false
						evts := ctx.EventManager().Events()
						for _, e := range evts {
							if e.Type == "wasm-cyber1" {
								require.False(t, hadCyber1)
								attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
								require.Empty(t, err)

								require.Equal(t,
									[]wasmTypes.Attribute{
										{Key: "contract_address", Value: []byte(contractAddress.String())},
										{Key: "attr1", Value: []byte("🤯")},
									},
									attrs,
								)

								hadCyber1 = true
							}
						}

						require.True(t, hadCyber1)
					}

					queryRes, err := queryHelper(t, keeper, ctx, contractAddress, `{"q":{}}`, true, math.MaxUint64)

					require.Empty(t, err)
					require.Equal(t, test.output, queryRes)
				}
			})
		}
	}
}

type ContractInfo struct {
	Address string `json:"address"`
	Hash    string `json:"hash"`
}

func TestIBCPacketAck(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, TestContractPaths[ibcContract], sdk.NewCoins())

	_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"init":{}}`, true, defaultGasForIbcTests)
	require.Empty(t, err)

	for _, test := range []struct {
		description   string
		sequence      uint64
		output        string
		isSuccess     bool
		hasAttributes bool
		hasEvents     bool
	}{
		{
			description:   "Default",
			sequence:      0,
			output:        "8",
			isSuccess:     true,
			hasAttributes: false,
			hasEvents:     false,
		},
		{
			description:   "SubmessageNoReply",
			sequence:      1,
			output:        "14",
			isSuccess:     true,
			hasAttributes: false,
			hasEvents:     false,
		},
		{
			description:   "SubmessageWithReply",
			sequence:      2,
			output:        "21",
			isSuccess:     true,
			hasAttributes: false,
			hasEvents:     false,
		},
		{
			description:   "Attributes",
			sequence:      3,
			output:        "11",
			isSuccess:     true,
			hasAttributes: true,
			hasEvents:     false,
		},
		{
			description:   "Events",
			sequence:      4,
			output:        "12",
			isSuccess:     true,
			hasAttributes: false,
			hasEvents:     true,
		},
		{
			description:   "Error",
			sequence:      5,
			output:        "",
			isSuccess:     false,
			hasAttributes: false,
			hasEvents:     false,
		},
	} {
		t.Run(test.description, func(t *testing.T) {
			ibcPacket := createIBCPacket(createIBCEndpoint(PortIDForContract(contractAddress), "channel.1"),
				createIBCEndpoint(PortIDForContract(contractAddress), "channel.0"),
				test.sequence,
				createIBCTimeout(math.MaxUint64),
				[]byte{},
			)
			ack := make([]byte, 8)
			binary.LittleEndian.PutUint64(ack, uint64(test.sequence))

			ctx, events, err := ibcPacketAckHelper(t, keeper, ctx, contractAddress, privKeyA, defaultGasForIbcTests, ibcPacket, ack)

			if !test.isSuccess {
				require.Contains(t, fmt.Sprintf("%+v", err), "Intentional")
			} else {
				require.Empty(t, err)
				if test.hasAttributes {
					require.Equal(t,
						[]ContractEvent{
							{
								{Key: "contract_address", Value: []byte(contractAddress.String())},
								{Key: "attr1", Value: []byte("😗")},
							},
						},
						events,
					)
				}

				if test.hasEvents {
					hadCyber1 := false
					evts := ctx.EventManager().Events()
					for _, e := range evts {
						if e.Type == "wasm-cyber1" {
							require.False(t, hadCyber1)
							attrs, err := parseAndDecryptAttributes(e.Attributes, []byte{})
							require.Empty(t, err)

							require.Equal(t,
								[]wasmTypes.Attribute{
									{Key: "contract_address", Value: []byte(contractAddress.String())},
									{Key: "attr1", Value: []byte("🤯")},
								},
								attrs,
							)

							hadCyber1 = true
						}
					}

					require.True(t, hadCyber1)
				}

				queryRes, err := queryHelper(t, keeper, ctx, contractAddress, `{"q":{}}`, true, math.MaxUint64)

				require.Empty(t, err)

				require.Equal(t, test.output, queryRes)
			}
		})
	}
}

func TestIBCPacketTimeout(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, TestContractPaths[ibcContract], sdk.NewCoins())

	_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"init":{}}`, true, defaultGasForIbcTests)
	require.Empty(t, err)

	for _, test := range []struct {
		description   string
		sequence      uint64
		output        string
		isSuccess     bool
		hasAttributes bool
		hasEvents     bool
	}{
		{
			description:   "Default",
			sequence:      0,
			output:        "9",
			isSuccess:     true,
			hasAttributes: false,
			hasEvents:     false,
		},
		{
			description:   "SubmessageNoReply",
			sequence:      1,
			output:        "15",
			isSuccess:     true,
			hasAttributes: false,
			hasEvents:     false,
		},
		{
			description:   "SubmessageWithReply",
			sequence:      2,
			output:        "22",
			isSuccess:     true,
			hasAttributes: false,
			hasEvents:     false,
		},
		{
			description:   "Attributes",
			sequence:      3,
			output:        "12",
			isSuccess:     true,
			hasAttributes: true,
			hasEvents:     false,
		},
		{
			description:   "Events",
			sequence:      4,
			output:        "13",
			isSuccess:     true,
			hasAttributes: false,
			hasEvents:     true,
		},
		{
			description:   "Error",
			sequence:      5,
			output:        "",
			isSuccess:     false,
			hasAttributes: false,
			hasEvents:     false,
		},
	} {
		t.Run(test.description, func(t *testing.T) {
			ibcPacket := createIBCPacket(createIBCEndpoint(PortIDForContract(contractAddress), "channel.1"),
				createIBCEndpoint(PortIDForContract(contractAddress), "channel.0"),
				test.sequence,
				createIBCTimeout(math.MaxUint64),
				[]byte{},
			)

			ctx, events, err := ibcPacketTimeoutHelper(t, keeper, ctx, contractAddress, privKeyA, defaultGasForIbcTests, ibcPacket)

			if !test.isSuccess {
				require.Contains(t, fmt.Sprintf("%+v", err), "Intentional")
			} else {
				require.Empty(t, err)
				if test.hasAttributes {
					require.Equal(t,
						[]ContractEvent{
							{
								{Key: "contract_address", Value: []byte(contractAddress.String())},
								{Key: "attr1", Value: []byte("😗")},
							},
						},
						events,
					)
				}

				if test.hasEvents {
					hadCyber1 := false
					evts := ctx.EventManager().Events()
					for _, e := range evts {
						if e.Type == "wasm-cyber1" {
							require.False(t, hadCyber1)
							attrs, err := parseAndDecryptAttributes(e.Attributes, []byte{})
							require.Empty(t, err)

							require.Equal(t,
								[]wasmTypes.Attribute{
									{Key: "contract_address", Value: []byte(contractAddress.String())},
									{Key: "attr1", Value: []byte("🤯")},
								},
								attrs,
							)

							hadCyber1 = true
						}
					}

					require.True(t, hadCyber1)
				}

				queryRes, err := queryHelper(t, keeper, ctx, contractAddress, `{"q":{}}`, true, math.MaxUint64)

				require.Empty(t, err)

				require.Equal(t, test.output, queryRes)
			}
		})
	}
}

// getDecryptedData decrypts the output of the first function to be called
// Only returns the data, logs and messages from the first function call
func getDecryptedData(t *testing.T, data []byte, nonce []byte) []byte {
	if len(data) == 0 {
		return data
	}

	dataPlaintextBase64, err := wasmCtx.Decrypt(data, nonce)
	require.NoError(t, err)

	dataPlaintext, err := base64.StdEncoding.DecodeString(string(dataPlaintextBase64))
	require.NoError(t, err)

	return dataPlaintext
}
