package keeper

import (
	"fmt"

	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	cosmwasm "github.com/trstlabs/trst/go-cosmwasm/types"
)

func TestEncryptedAttributesFromInitWithoutSubmessageWithoutReply(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, events, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"add_attributes":{}}`, true, defaultGasForTests)
	require.Empty(t, err)

	require.Equal(t,
		[]ContractEvent{
			{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
			},
		},
		events,
	)

}

func TestEncryptedAttributesFromInitWithSubmessageWithoutReply(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, events, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"add_attributes_with_submessage":{"id":0}}`, true, defaultGasForTests)
	require.Empty(t, err)
	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[0])
	require.ElementsMatch(t, ContractEvent{
		{Key: "attr3", Value: []byte("🍉"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr4", Value: []byte("🥝"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[1])

	/*require.ElementsMatch(t,
		[]ContractEvent{
			{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
			},
			{

				{Key: "attr3", Value: []byte("🍉"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr4", Value: []byte("🥝"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
			},
		},
		events,
	)*/

}

func TestEncryptedAttributesFromInitWithSubmessageWithReply(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	_, _, contractAddress, events, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"add_attributes_with_submessage":{"id":2200}}`, true, defaultGasForTests)
	require.Empty(t, err)

	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[0])
	require.ElementsMatch(t, ContractEvent{
		{Key: "attr3", Value: []byte("🍉"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr4", Value: []byte("🥝"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[1])
	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr5", Value: []byte("🤯"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr6", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[2])
}

func TestEncryptedAttributesFromExecuteWithoutSubmessageWithoutReply(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)
	_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"add_attributes":{}}`, true, defaultGasForTests, 0)
	require.Empty(t, err)

	require.Equal(t,
		[]ContractEvent{
			{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
			},
		},
		events,
	)

}

func TestEncryptedAttributesFromExecuteWithSubmessageWithoutReply(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)
	_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"add_attributes_with_submessage":{"id":0}}`, true, defaultGasForTests, 0)
	require.Empty(t, err)

	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[0])
	require.ElementsMatch(t, ContractEvent{
		{Key: "attr3", Value: []byte("🍉"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr4", Value: []byte("🥝"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[1])

	/*
		require.ElementsMatch(t,
			[]ContractEvent{
				{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
				},
				{

					{Key: "attr3", Value: []byte("🍉"), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr4", Value: []byte("🥝"), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				},
			},
			events,
		)*/

}

func TestEncryptedAttributesFromExecuteWithSubmessageWithReply(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)
	_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"add_attributes_with_submessage":{"id":2200}}`, true, defaultGasForTests, 0)
	require.Empty(t, err)

	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[0])
	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr3", Value: []byte("🍉"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr4", Value: []byte("🥝"), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[1])
	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr5", Value: []byte("🤯"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr6", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[2])

}

func TestPlaintextFromInitWithoutSubmessageWithoutReply(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, events, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"add_plaintext_attributes":{}}`, true, defaultGasForTests, true)
	require.Empty(t, err)

	require.Equal(t,
		[]ContractEvent{
			{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
			},
		},
		events,
	)

}

func TestPlaintextAttributesFromInitWithSubmessageWithoutReply(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, events, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"add_plaintext_attributes_with_submessage":{"id":0}}`, true, defaultGasForTests, true)
	require.Empty(t, err)

	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[0])
	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr3", Value: []byte("🍉"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr4", Value: []byte("🥝"), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[1])
	/*
		require.ElementsMatch(t,
			[]ContractEvent{
				{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
				},
				{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr3", Value: []byte("🍉"), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr4", Value: []byte("🥝"), AccAddr: "", Encrypted: false, PubDb: false},
				},
			},
			events,
		)*/

}

func TestPlaintextAttributesFromInitWithSubmessageWithReply(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	_, _, contractAddress, events, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"add_plaintext_attributes_with_submessage":{"id":2300}}`, true, defaultGasForTests, true)
	require.Empty(t, err)

	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[0])
	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr3", Value: []byte("🍉"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr4", Value: []byte("🥝"), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[1])
	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr5", Value: []byte("🤯"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr6", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[2])
}

func TestPlaintextAttributesFromExecuteWithoutSubmessageWithoutReply(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)
	_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"add_plaintext_attributes":{}}`, true, defaultGasForTests, 0, true)
	require.Empty(t, err)

	require.Equal(t,
		[]ContractEvent{
			{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
			},
		},
		events,
	)

}

func TestPlaintextAttributesFromExecuteWithSubmessageWithoutReply(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)
	_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"add_plaintext_attributes_with_submessage":{"id":0}}`, true, defaultGasForTests, 0, true)
	require.Empty(t, err)

	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[0])
	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr3", Value: []byte("🍉"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr4", Value: []byte("🥝"), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[1])

}

func TestPlaintextAttributesFromExecuteWithSubmessageWithReply(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)
	_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"add_plaintext_attributes_with_submessage":{"id":2300}}`, true, defaultGasForTests, 0, true)
	require.Empty(t, err)

	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[0])
	require.ElementsMatch(t, ContractEvent{
		{Key: "attr3", Value: []byte("🍉"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr4", Value: []byte("🥝"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[1])
	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr5", Value: []byte("🤯"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr6", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[2])

}

func TestEncryptedEventsFromInitWithoutSubmessageWithoutReply(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	nonce, ctx, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"add_events":{}}`, true, defaultGasForTests)

	require.Empty(t, err)

	events := ctx.EventManager().Events()

	hadCyber1 := false
	hadCyber2 := false
	for _, e := range events {
		if e.Type == "wasm-cyber1" {
			require.False(t, hadCyber1)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.Equal(t,
				[]cosmwasm.Attribute{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
				},
				attrs,
			)

			hadCyber1 = true
		}

		if e.Type == "wasm-cyber2" {
			require.False(t, hadCyber2)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.Equal(t,
				[]cosmwasm.Attribute{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr3", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr4", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
				},
				attrs,
			)

			hadCyber2 = true
		}
	}

	require.True(t, hadCyber1)
	require.True(t, hadCyber2)
}

func TestEncryptedEventsFromInitWithSubmessageWithoutReply(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	nonce, ctx, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"add_events_with_submessage":{"id":0}}`, true, defaultGasForTests)
	require.Empty(t, err)

	events := ctx.EventManager().Events()

	hadCyber1 := false
	hadCyber2 := false
	hadCyber3 := false
	hadCyber4 := false

	for _, e := range events {
		if e.Type == "wasm-cyber1" {
			require.False(t, hadCyber1)
			attrs, _ := parseAndDecryptAttributes(e.Attributes, nonce)
			//require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)

			hadCyber1 = true
		}

		if e.Type == "wasm-cyber2" {
			require.False(t, hadCyber2)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr3", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr4", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)

			hadCyber2 = true
		}

		if e.Type == "wasm-cyber3" {
			require.False(t, hadCyber3)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "attr1", Value: []byte("🤯"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("🤟"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)

			hadCyber3 = true
		}

		if e.Type == "wasm-cyber4" {
			require.False(t, hadCyber4)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr3", Value: []byte("😅"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr4", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)

			hadCyber4 = true
		}
	}

	require.True(t, hadCyber1)
	require.True(t, hadCyber2)
	require.True(t, hadCyber3)
	require.True(t, hadCyber4)
}

func TestEncryptedEventsFromInitWithSubmessageWithReply(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	nonce, ctx, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"add_events_with_submessage":{"id":2400}}`, true, defaultGasForTests)
	require.Empty(t, err)

	events := ctx.EventManager().Events()

	hadCyber1 := false
	hadCyber2 := false
	hadCyber3 := false
	hadCyber4 := false
	hadCyber5 := false
	hadCyber6 := false

	for _, e := range events {
		if e.Type == "wasm-cyber1" {
			require.False(t, hadCyber1)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)

			hadCyber1 = true
		}

		if e.Type == "wasm-cyber2" {
			require.False(t, hadCyber2)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr3", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr4", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)

			hadCyber2 = true
		}

		if e.Type == "wasm-cyber3" {
			require.False(t, hadCyber3)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("🤯"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("🤟"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)
			hadCyber3 = true
		}

		if e.Type == "wasm-cyber4" {
			require.False(t, hadCyber4)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr3", Value: []byte("😅"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr4", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)

			hadCyber4 = true
		}

		if e.Type == "wasm-cyber5" {
			require.False(t, hadCyber5)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("😗"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("😋"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)
			hadCyber5 = true
		}

		if e.Type == "wasm-cyber6" {
			require.False(t, hadCyber6)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr3", Value: []byte("😉"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr4", Value: []byte("😊"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)
			hadCyber6 = true
		}
	}

	require.True(t, hadCyber1)
	require.True(t, hadCyber2)
	require.True(t, hadCyber3)
	require.True(t, hadCyber4)
	require.True(t, hadCyber5)
	require.True(t, hadCyber6)
}

func TestEncryptedEventsFromExecuteWithoutSubmessageWithoutReply(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)
	nonce, ctx, _, _, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"add_events":{}}`, true, defaultGasForTests, 0)
	require.Empty(t, err)

	events := ctx.EventManager().Events()

	hadCyber1 := false
	hadCyber2 := false

	for _, e := range events {
		if e.Type == "wasm-cyber1" {
			require.False(t, hadCyber1)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.Equal(t,
				[]cosmwasm.Attribute{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
				},
				attrs,
			)

			hadCyber1 = true
		}

		if e.Type == "wasm-cyber2" {
			require.False(t, hadCyber2)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr3", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr4", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)

			hadCyber2 = true
		}
	}

	require.True(t, hadCyber1)
	require.True(t, hadCyber2)
}

func TestEncryptedEventsFromExecuteWithSubmessageWithoutReply(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)
	nonce, ctx, _, _, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"add_events_with_submessage":{"id":0}}`, true, defaultGasForTests, 0)
	require.Empty(t, err)

	events := ctx.EventManager().Events()

	hadCyber1 := false
	hadCyber2 := false
	hadCyber3 := false
	hadCyber4 := false

	for _, e := range events {
		if e.Type == "wasm-cyber1" {
			require.False(t, hadCyber1)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)

			hadCyber1 = true
		}

		if e.Type == "wasm-cyber2" {
			require.False(t, hadCyber2)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr3", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr4", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)

			hadCyber2 = true
		}

		if e.Type == "wasm-cyber3" {
			require.False(t, hadCyber3)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("🤯"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("🤟"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)
			hadCyber3 = true
		}

		if e.Type == "wasm-cyber4" {
			require.False(t, hadCyber4)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr3", Value: []byte("😅"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr4", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)
			hadCyber4 = true
		}
	}

	require.True(t, hadCyber1)
	require.True(t, hadCyber2)
	require.True(t, hadCyber3)
	require.True(t, hadCyber4)
}

func TestEncryptedEventsFromExecuteWithSubmessageWithReply(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)
	nonce, ctx, _, _, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"add_events_with_submessage":{"id":2400}}`, true, defaultGasForTests, 0)
	require.Empty(t, err)

	events := ctx.EventManager().Events()

	hadCyber1 := false
	hadCyber2 := false
	hadCyber3 := false
	hadCyber4 := false
	hadCyber5 := false
	hadCyber6 := false

	for _, e := range events {
		if e.Type == "wasm-cyber1" {
			require.False(t, hadCyber1)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.Equal(t,
				[]cosmwasm.Attribute{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
				},
				attrs,
			)

			hadCyber1 = true
		}

		if e.Type == "wasm-cyber2" {
			require.False(t, hadCyber2)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr3", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr4", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)
			hadCyber2 = true
		}

		if e.Type == "wasm-cyber3" {
			require.False(t, hadCyber3)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("🤯"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("🤟"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)
			hadCyber3 = true
		}

		if e.Type == "wasm-cyber4" {
			require.False(t, hadCyber4)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr3", Value: []byte("😅"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr4", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)
			hadCyber4 = true
		}

		if e.Type == "wasm-cyber5" {
			require.False(t, hadCyber5)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("😗"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("😋"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)
			hadCyber5 = true
		}

		if e.Type == "wasm-cyber6" {
			require.False(t, hadCyber6)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr3", Value: []byte("😉"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr4", Value: []byte("😊"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)
			hadCyber6 = true
		}
	}

	require.True(t, hadCyber1)
	require.True(t, hadCyber2)
	require.True(t, hadCyber3)
	require.True(t, hadCyber4)
	require.True(t, hadCyber5)
	require.True(t, hadCyber6)
}

func TestMixedLogsFromInitWithoutSubmessageWithoutReply(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	nonce, ctx, contractAddress, logs, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"add_mixed_attributes_and_events":{}}`, true, defaultGasForTests, true)

	require.Empty(t, err)

	events := ctx.EventManager().Events()

	hadCyber1 := false
	for _, e := range events {
		if e.Type == "wasm-cyber1" {
			require.False(t, hadCyber1)
			attrs, _ := parseAndDecryptAttributes(e.Attributes, nonce)
			//require.NotEmpty(t, hasPlaintext)
			require.Equal(t,
				[]cosmwasm.Attribute{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
				},
				attrs,
			)

			hadCyber1 = true
		}
	}

	require.True(t, hadCyber1)

	require.ElementsMatch(t,
		ContractEvent{

			{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
			{Key: "attr3", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
			{Key: "attr4", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
		},
		logs[0],
	)
}

func TestMixedAttributesAndEventsFromInitWithSubmessageWithoutReply(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	nonce, ctx, contractAddress, logs, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"add_mixed_attributes_and_events_with_submessage":{"id":0}}`, true, defaultGasForTests)
	require.Empty(t, err)

	events := ctx.EventManager().Events()

	hadCyber1 := false
	hadCyber2 := false

	for _, e := range events {
		if e.Type == "wasm-cyber1" {
			require.False(t, hadCyber1)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)

			hadCyber1 = true
		}

		if e.Type == "wasm-cyber2" {
			require.False(t, hadCyber2)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr5", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr6", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)
			hadCyber2 = true
		}
	}

	require.True(t, hadCyber1)
	require.True(t, hadCyber2)

	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr3", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr4", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
	}, logs[0])

}

func TestMixedAttributesAndEventsFromInitWithSubmessageWithReply(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	nonce, ctx, contractAddress, logs, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"add_mixed_attributes_and_events_with_submessage":{"id":2500}}`, true, defaultGasForTests)
	require.Empty(t, err)

	events := ctx.EventManager().Events()

	hadCyber1 := false
	hadCyber2 := false
	hadCyber3 := false

	for _, e := range events {
		if e.Type == "wasm-cyber1" {
			require.False(t, hadCyber1)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)
			hadCyber1 = true
		}

		if e.Type == "wasm-cyber2" {
			require.False(t, hadCyber2)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr5", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr6", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)
			hadCyber2 = true
		}

		if e.Type == "wasm-cyber3" {
			require.False(t, hadCyber3)
			attrs, _ := parseAndDecryptAttributes(e.Attributes, nonce)
			//require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr9", Value: []byte("🤯"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr10", Value: []byte("🤟"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)
			hadCyber3 = true
		}
	}

	require.True(t, hadCyber1)
	require.True(t, hadCyber2)
	require.True(t, hadCyber3)

	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr3", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr4", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
	}, logs[0])
	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
	}, logs[1])
	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr11", Value: []byte("😉"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr12", Value: []byte("😊"), AccAddr: "", Encrypted: false, PubDb: false},
	}, logs[4])

}

func TestMixedAttributesAndEventsFromExecuteWithoutSubmessageWithoutReply(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)
	nonce, ctx, _, logs, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"add_mixed_attributes_and_events":{}}`, true, defaultGasForTests, 0, true)
	require.Empty(t, err)

	events := ctx.EventManager().Events()

	hadCyber1 := false

	for _, e := range events {
		if e.Type == "wasm-cyber1" {
			require.False(t, hadCyber1)
			attrs, _ := parseAndDecryptAttributes(e.Attributes, nonce)
			//require.NotEmpty(t, hasPlaintext)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)

			hadCyber1 = true
		}
	}

	require.True(t, hadCyber1)

	require.ElementsMatch(t,
		ContractEvent{

			{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
			{Key: "attr3", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
			{Key: "attr4", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
		},
		logs[0],
	)
}

func TestMixedAttributesAndEventsFromExecuteWithSubmessageWithoutReply(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)
	nonce, ctx, _, logs, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"add_mixed_attributes_and_events_with_submessage":{"id":0}}`, true, defaultGasForTests, 0)
	require.Empty(t, err)

	events := ctx.EventManager().Events()

	hadCyber1 := false
	hadCyber2 := false

	for _, e := range events {
		if e.Type == "wasm-cyber1" {
			require.False(t, hadCyber1)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)

			hadCyber1 = true
		}

		if e.Type == "wasm-cyber2" {
			require.False(t, hadCyber2)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr5", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr6", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)

			hadCyber2 = true
		}
	}

	require.True(t, hadCyber1)
	require.True(t, hadCyber2)

	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr3", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr4", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
	}, logs[0])
	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr7", Value: []byte("😅"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr8", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
	}, logs[2])

}

func TestMixedAttributesAndEventsFromExecuteWithSubmessageWithReply(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)
	nonce, ctx, _, logs, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"add_mixed_attributes_and_events_with_submessage":{"id":2500}}`, true, defaultGasForTests, 0)
	require.Empty(t, err)

	events := ctx.EventManager().Events()

	hadCyber1 := false
	hadCyber2 := false
	hadCyber3 := false

	for _, e := range events {
		fmt.Printf("ev %+v \n", e)
		fmt.Printf("Attributes %+v \n", e.Attributes)
		fmt.Printf("Attribute2 %s \n", e.Attributes[1].Key)
		if e.Type == "wasm-cyber1" {
			require.False(t, hadCyber1)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)
			hadCyber1 = true
		}

		if e.Type == "wasm-cyber2" {
			require.False(t, hadCyber2)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr5", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr6", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)

			hadCyber2 = true
		}

		if e.Type == "wasm-cyber3" {
			require.False(t, hadCyber3)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr9", Value: []byte("🤯"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr10", Value: []byte("🤟"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)

			hadCyber3 = true
		}
	}

	require.True(t, hadCyber1)
	require.True(t, hadCyber2)
	require.True(t, hadCyber3)

	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr3", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr4", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
	}, logs[0])
	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr7", Value: []byte("😅"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr8", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
	}, logs[2])
	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr11", Value: []byte("😉"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr12", Value: []byte("😊"), AccAddr: "", Encrypted: false, PubDb: false},
	}, logs[4])
}
