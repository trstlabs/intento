package keeper

import (

	//"log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/danieljdd/tpp/x/compute/internal/types"
)

// ContractPayout pays the creator of the contract
func (k Keeper) ContractPayout(ctx sdk.Context, contractAddress sdk.AccAddress) error {
	store := ctx.KVStore(k.storeKey)
	contractBz := store.Get(types.GetContractAddressKey(contractAddress))
	if contractBz == nil {
		return sdkerrors.Wrap(types.ErrNotFound, "contract")
	}
	var contract types.ContractInfo
	k.cdc.MustUnmarshalBinaryBare(contractBz, &contract)

	//payout contract coins to the creator
	balance := k.bankKeeper.GetAllBalances(ctx, contractAddress)
	if !balance.Empty() {
		k.bankKeeper.SendCoins(ctx, contractAddress, contract.Creator, balance)
	}
	return nil
}

// CallLastMsg executes a final message before end-blocker deletion
func (k Keeper) CallLastMsg(ctx sdk.Context, contractAddress sdk.AccAddress) (err error) {

	//get codeid first
	info, err := k.GetContractInfo(ctx, contractAddress)
	if err != nil {
		return err
	}

	/*

		signerSig := []byte{}
		signBytes := []byte{}
		//var err error

			signerSig, signBytes, err = k.GetSignerInfo(ctx, contractAddress)
			if err != nil {
				return err
			}

		verificationInfo := types.NewVerificationInfo(signBytes, signerSig, nil)

		codeInfo, prefixStore, err := k.contractInstance(ctx, contractAddress)
		if err != nil {
			return nil
		}

		store := ctx.KVStore(k.storeKey)


		contractKey := store.Get(types.GetContractEnclaveKey(contractAddress))
		if contractKey == nil {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "contract key not found")
		}
		//fmt.Printf("Contract Execute: Got contract Key for contract %s: %s\n", contractAddress, base64.StdEncoding.EncodeToString(contractKey))
		params := types.NewEnv(ctx, contractAddress, sdk.NewCoins(sdk.NewCoin("tpp", sdk.ZeroInt())), contractAddress, contractKey)
		//fmt.Printf("Contract Execute: key from params %s \n", params.Key)

		// prepare querier
		querier := QueryHandler{
			Ctx:     ctx,
			Plugins: k.queryPlugins,
		}

		gas := gasForContract(ctx)
		//	fmt.Printf("Execute message before wasm is %s \n", base64.StdEncoding.EncodeToString(msg))
		res, gasUsed, execErr := k.wasmer.Execute(codeInfo.CodeHash, params, codeInfo.LastMsg, prefixStore, cosmwasmAPI, querier, gasMeter(ctx), gas, verificationInfo)
		consumeGas(ctx, gasUsed)

		if execErr != nil {
			return nil, sdkerrors.Wrap(types.ErrExecuteFailed, execErr.Error())
		}




	*/
	if info.LastMsg != nil {
		res, err := k.Execute(ctx, contractAddress, contractAddress, info.LastMsg, sdk.NewCoins(sdk.NewCoin("tpp", sdk.ZeroInt())), nil)
		if err != nil {
			return err
		}
		k.SetContractResult(ctx, contractAddress, res)
	}
	return nil
}

/*
func (k Keeper) getConsensusIoPubKey(ctx sdk.Context) ([]byte, error) {
	var masterIoKey reg.MasterCertificate

	//route := fmt.Sprintf("custom/%s/%s", types.RegisterQuerierRoute, types.QueryMasterCertificate)

	res, _, err := ctx.Context().Query("/tpp.x.registration.v1beta1.Query/MasterKey")
	if err != nil {
		//	res, _, err = ctx.CLIContext.Query("/tpp.x.registration.v1beta1.Query/MasterKey")
		if err != nil {
			return nil, err
		}
	}

	err = encoding.GetCodec(proto.Name).Unmarshal(res, &response)
	if err != nil {
		return nil, err
	}

	ioPubkey, err := ra.VerifyRaCert(response.MasterKey.Bytes)
	if err != nil {
		return nil, err
	}

	return ioPubkey, nil
}

func (k Keeper) getTxEncryptionKey(txSenderPrivKey []byte, nonce []byte) ([]byte, error) {

	consensusIoPubKeyBytes, err := k.getConsensusIoPubKey()
	if err != nil {
		fmt.Println("Failed to get IO key. Make sure the CLI and the node you are targeting are operating in the same SGX mode")
		return nil, err
	}

	txEncryptionIkm, err := curve25519.X25519(txSenderPrivKey, consensusIoPubKeyBytes)
	if err != nil {
		return nil, err
	}
	kdfFunc := hkdf.New(sha256.New, append(txEncryptionIkm[:], nonce...), hkdfSalt, []byte{})

	txEncryptionKey := make([]byte, 32)
	if _, err := io.ReadFull(kdfFunc, txEncryptionKey); err != nil {
		return nil, err
	}

	return txEncryptionKey, nil
}

// Encrypt encrypts
func (k Keeper) Encrypt(plaintext []byte) ([]byte, error) {
	txSenderPrivKey, txSenderPubKey, err := k.GetTxSenderKeyPair()
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, 32)
	_, err = rand.Read(nonce)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	txEncryptionKey, err := k.getTxEncryptionKey(txSenderPrivKey, nonce)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return encryptData(txEncryptionKey, txSenderPubKey, plaintext, nonce)
}

func encryptData(aesEncryptionKey []byte, txSenderPubKey []byte, plaintext []byte, nonce []byte) ([]byte, error) {
	cipher, err := miscreant.NewAESCMACSIV(aesEncryptionKey)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	ciphertext, err := cipher.Seal(nil, plaintext, []byte{})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// ciphertext = nonce(32) || wallet_pubkey(32) || ciphertext
	ciphertext = append(nonce, append(txSenderPubKey, ciphertext...)...)

	return ciphertext, nil
}
*/
