package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/danieljdd/trst/x/registration/internal/types"
	ra "github.com/danieljdd/trst/x/registration/remote_attestation"
)

type grpcQuerier struct {
	keeper Keeper
}

// todo: this needs proper tests and doc
func NewQuerier(keeper Keeper) grpcQuerier {
	return grpcQuerier{keeper: keeper}
}

func (q grpcQuerier) MasterKey(c context.Context, req *types.QueryMasterKeyRequest) (*types.QueryMasterKeyResponse, error) {
	/*if req == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "empty request")
	}*/
	//var response types.QueryMasterKeyResponse
	//	var response *types.MasterCertificate

	rsp, err := queryMasterKey(sdk.UnwrapSDKContext(c), q.keeper)
	switch {
	case err != nil:
		return nil, err
	case rsp == nil:
		return nil, types.ErrNotFound
	}

	/*	err = encoding.GetCodec(proto.Name).Unmarshal(rsp.Bytes, &response.Bytes)
		if err != nil {
			return nil, err
		}
	*/
	/*err = encoding.GetCodec(proto.Name).Unmarshal(rsp.Bytes, &response)
	if err != nil {
		return nil, err
	}
	fmt.Printf("response %+v\n", &response)*/
	//	fmt.Printf("master key %+v\n", rsp)

	ioPubkey, err := ra.VerifyRaCert(rsp.Bytes)
	if err != nil {
		return nil, err
	}
	//	fmt.Printf("response %+v\n", ioPubkey)
	rsp.Bytes = ioPubkey
	return &types.QueryMasterKeyResponse{MasterKey: rsp}, nil
}

/*
func (q grpcQuerier) MasterKeyBytes(c context.Context, req *types.QueryMasterKeyBytesRequest) (*types.QueryMasterKeyBytesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	var response types.QueryMasterKeyResponse

	rsp, err := queryMasterKey(sdk.UnwrapSDKContext(c), q.keeper)
	switch {
	case err != nil:
		return nil, err
	case rsp == nil:
		return nil, types.ErrNotFound
	}

	result := &types.QueryMasterKeyResponse{MasterKey: rsp}

	err = encoding.GetCodec(proto.Name).Unmarshal(rsp, &response)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	//keeper.GetMasterCertificate(ctx, types.MasterIoKeyId)
	ioPubKey := q.keeper.GetMasterIoPubKeyArray(ctx)

	fmt.Printf("query result: %X\n", result)
	fmt.Printf("query ioPubKey bytes: %X\n", &ioPubKey)
	//store := ctx.KVStore(k.storeKey)

	//codeHash := store.Get([]byte(req.Codeid))
	if ioPubKey == nil {
		return nil, status.Error(codes.InvalidArgument, "no io PubKey")
	}

	err := encoding.GetCodec(proto.Name).Unmarshal(ioPubKey, &response.Bytes)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "marshalling failed")
	}

	ioPubkey, err := ra.VerifyRaCert(response.Bytes)
	if err != nil {
		return nil, err
	}

	fmt.Printf("query ioPubKey response: %X\n", &response)
	fmt.Printf("query ioPubKey response bytes: %X\n", response)
	//itemStore := prefix.NewStore(store, types.KeyPrefix(types.ItemKey))

	return &types.QueryMasterKeyBytesResponse{Masterkey: ioPubkey}, nil
}
*/ /*

func (q grpcQuerier) MasterKeyBytes(c context.Context, req *types.QueryMasterKeyBytesRequest) (*types.QueryMasterKeyBytesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	var response *types.MasterCertificate
	ctx := sdk.UnwrapSDKContext(c)
	//keeper.GetMasterCertificate(ctx, types.MasterIoKeyId)
	ioPubKey := q.keeper.GetMasterIoPubKeyArray(ctx)

	fmt.Printf("query ioPubKey: %X\n", ioPubKey)
	fmt.Printf("query ioPubKey bytes: %X\n", &ioPubKey)
	//store := ctx.KVStore(k.storeKey)

	//codeHash := store.Get([]byte(req.Codeid))
	if ioPubKey == nil {
		return nil, status.Error(codes.InvalidArgument, "no io PubKey")
	}

	err := encoding.GetCodec(proto.Name).Unmarshal(ioPubKey, &response.Bytes)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "marshalling failed")
	}

	ioPubkey, err := ra.VerifyRaCert(response.Bytes)
	if err != nil {
		return nil, err
	}

	fmt.Printf("query ioPubKey response: %X\n", &response)
	fmt.Printf("query ioPubKey response bytes: %X\n", response)
	//itemStore := prefix.NewStore(store, types.KeyPrefix(types.ItemKey))

	return &types.QueryMasterKeyBytesResponse{Masterkey: ioPubkey}, nil
}*/

func (q grpcQuerier) EncryptedSeed(c context.Context, req *types.QueryEncryptedSeedRequest) (*types.QueryEncryptedSeedResponse, error) {
	if req.PubKey == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "public key")
	}
	rsp, err := queryEncryptedSeed(sdk.UnwrapSDKContext(c), req.PubKey, q.keeper)
	switch {
	case err != nil:
		return nil, err
	case rsp == nil:
		return nil, types.ErrNotFound
	}
	return &types.QueryEncryptedSeedResponse{EncryptedSeed: rsp}, nil
}

func queryMasterKey(ctx sdk.Context, keeper Keeper) (*types.MasterCertificate, error) {
	ioKey := keeper.GetMasterCertificate(ctx, types.MasterIoKeyId)
	//nodeKey := keeper.GetMasterCertificate(ctx, types.MasterNodeKeyId)
	if ioKey == nil { //|| nodeKey == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownAddress, "Chain has not been initialized yet")
	}
	//resp := types.GenesisState{
	//	Registration:              nil,
	//	NodeExchMasterCertificate: nodeKey,
	//	IoMasterCertificate:       ioKey,
	//}

	//asBytes, err := json.Marshal(ioKey)
	//if err != nil {
	//	return nil, err
	//}

	return ioKey, nil
}

func queryEncryptedSeed(ctx sdk.Context, pubkeyBytes []byte, keeper Keeper) ([]byte, error) {
	seed := keeper.getRegistrationInfo(ctx, pubkeyBytes)
	if seed == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownAddress, "Node has not been authenticated yet")
	}

	return seed.EncryptedSeed, nil
}
