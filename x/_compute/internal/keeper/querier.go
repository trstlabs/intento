package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/golang/protobuf/ptypes/empty"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/trstlabs/trst/x/compute/internal/types"
)

type grpcQuerier struct {
	keeper Keeper
}

// todo: this needs proper tests and doc
func NewGrpcQuerier(keeper Keeper) grpcQuerier {
	return grpcQuerier{keeper: keeper}
}

// Params returns params of the mint module.
func (q grpcQuerier) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := q.keeper.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

func (q grpcQuerier) ContractInfo(c context.Context, req *types.QueryContractInfoRequest) (*types.QueryContractInfoResponse, error) {
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}
	rsp, err := queryContractInfo(sdk.UnwrapSDKContext(c), addr, q.keeper)
	switch {
	case err != nil:
		return nil, err
	case rsp == nil:
		return nil, types.ErrNotFound
	}
	return &types.QueryContractInfoResponse{
		ContractInfo: rsp.ContractInfo,
	}, nil
}

func (q grpcQuerier) ContractPublicState(c context.Context, req *types.QueryContractPublicStateRequest) (*types.QueryContractPublicStateResponse, error) {
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}
	rsp, err := queryContractPublicState(sdk.UnwrapSDKContext(c), addr, q.keeper)
	switch {
	case err != nil:
		return nil, err
	case rsp == nil:
		return nil, types.ErrNotFound
	}
	return &types.QueryContractPublicStateResponse{
		PublicContractState: rsp,
	}, nil
}

func (q grpcQuerier) ContractPublicStateForAccount(c context.Context, req *types.QueryContractPublicStateForAccountRequest) (*types.QueryContractPublicStateForAccountResponse, error) {
	//contract addr to query
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}
	//user account to query
	acc, err := sdk.AccAddressFromBech32(req.Account)
	if err != nil {
		return nil, err
	}
	rsp, err := queryContractPublicStateForAccount(sdk.UnwrapSDKContext(c), addr, acc, q.keeper)
	switch {
	case err != nil:
		return nil, err
	case rsp == nil:
		return nil, types.ErrNotFound
	}
	return &types.QueryContractPublicStateForAccountResponse{
		PublicContractState: rsp,
	}, nil
}

func (q grpcQuerier) ContractPublicStateByKey(c context.Context, req *types.QueryContractPublicStateByKeyRequest) (*types.QueryContractPublicStateByKeyResponse, error) {
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}
	rsp, err := queryContractPublicStateByKey(sdk.UnwrapSDKContext(c), addr, req.Key, q.keeper)
	switch {
	case err != nil:
		return nil, err
	case rsp == types.KeyPair{}:
		fmt.Printf("not found or empty keypair \n")
		return nil, types.ErrNotFound
	}
	return &types.QueryContractPublicStateByKeyResponse{
		KeyPair: &rsp,
	}, nil
}

func (q grpcQuerier) ContractsByCode(c context.Context, req *types.QueryContractsByCodeRequest) (*types.QueryContractsByCodeResponse, error) {
	if req.CodeId == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "code id")
	}
	rsp, err := queryContractListByCode(sdk.UnwrapSDKContext(c), req.CodeId, q.keeper)
	switch {
	case err != nil:
		return nil, err
	case rsp == nil:
		return nil, types.ErrNotFound
	}
	return &types.QueryContractsByCodeResponse{
		ContractInfos: rsp,
	}, nil
}

func (q grpcQuerier) ContractPrivateState(c context.Context, req *types.QueryContractPrivateStateRequest) (*types.QueryContractPrivateStateResponse, error) {
	addr, err := sdk.AccAddressFromBech32(req.ContractAddress)
	if err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(c).WithGasMeter(sdk.NewGasMeter(q.keeper.queryGasLimit))
	rsp, err := q.keeper.QueryPrivate(ctx, addr, req.QueryData, false)
	switch {
	case err != nil:
		return nil, err
	case rsp == nil:
		return nil, types.ErrNotFound
	}
	return &types.QueryContractPrivateStateResponse{Data: rsp}, nil

}

func (q grpcQuerier) Code(c context.Context, req *types.QueryCodeRequest) (*types.QueryCodeResponse, error) {
	if req.CodeId == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "code id")
	}
	rsp, err := QueryCode(sdk.UnwrapSDKContext(c), req.CodeId, q.keeper)
	switch {
	case err != nil:
		return nil, err
	case rsp == nil:
		return nil, types.ErrNotFound
	}
	return &types.QueryCodeResponse{
		CodeInfoResponse: rsp.CodeInfoResponse,
		Data:             rsp.Data,
	}, nil
}

func (q grpcQuerier) Codes(c context.Context, _ *empty.Empty) (*types.QueryCodesResponse, error) {
	rsp, err := queryCodeList(sdk.UnwrapSDKContext(c), q.keeper)
	switch {
	case err != nil:
		return nil, err
	case rsp == nil:
		return nil, types.ErrNotFound
	}
	return &types.QueryCodesResponse{CodeInfos: rsp}, nil
}

func (q grpcQuerier) AddressByContractId(c context.Context, req *types.QueryContractAddressByContractIdRequest) (*types.QueryContractAddressByContractIdResponse, error) {
	ctx := sdk.UnwrapSDKContext(c).WithGasMeter(sdk.NewGasMeter(q.keeper.queryGasLimit))
	rsp, err := queryContractAddress(ctx, req.ContractId, q.keeper)
	switch {
	case err != nil:
		return nil, err
	case rsp == nil:
		return nil, types.ErrNotFound
	}
	return &types.QueryContractAddressByContractIdResponse{Address: rsp}, nil

}

func (q grpcQuerier) CodeHash(c context.Context, req *types.QueryCodeHashRequest) (*types.QueryCodeHashResponse, error) {
	if req == nil {
		return nil, types.ErrNotFound
	}

	ctx := sdk.UnwrapSDKContext(c)

	rsp, err := queryCodeHash(ctx, req.CodeId, q.keeper)
	switch {
	case err != nil:
		return nil, err
	case rsp == nil:
		return nil, types.ErrNotFound
	}

	if rsp == nil {
		return nil, types.ErrNotFound
	}

	return &types.QueryCodeHashResponse{CodeHash: rsp, CodeHashString: hex.EncodeToString(rsp)}, nil
}

func (q grpcQuerier) ContractKey(c context.Context, req *types.QueryContractKeyRequest) (*types.QueryContractKeyResponse, error) {
	if err := sdk.VerifyAddressFormat(req.Address); err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(c).WithGasMeter(sdk.NewGasMeter(q.keeper.queryGasLimit))
	rsp, err := queryContractKey(ctx, req.Address, q.keeper)
	switch {
	case err != nil:
		return nil, err
	case rsp == nil:
		return nil, types.ErrNotFound
	}
	return &types.QueryContractKeyResponse{Key: rsp}, nil

}

func (q grpcQuerier) ContractHash(c context.Context, req *types.QueryContractHashRequest) (*types.QueryContractHashResponse, error) {
	if err := sdk.VerifyAddressFormat(req.Address); err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(c).WithGasMeter(sdk.NewGasMeter(q.keeper.queryGasLimit))
	rsp, err := queryContractHash(ctx, req.Address, q.keeper)
	switch {
	case err != nil:
		return nil, err
	case rsp == nil:
		return nil, types.ErrNotFound
	}
	return &types.QueryContractHashResponse{CodeHash: rsp}, nil

}

func queryContractInfo(ctx sdk.Context, addr sdk.AccAddress, keeper Keeper) (*types.ContractInfoWithAddress, error) {
	info := keeper.GetContractInfo(ctx, addr)
	if info == nil {
		return nil, nil
	}
	// redact the Created field (just used for sorting, not part of public API)
	info.Created = nil
	return &types.ContractInfoWithAddress{
		Address:      addr,
		ContractInfo: info,
	}, nil
}

func queryContractPublicState(ctx sdk.Context, addr sdk.AccAddress, keeper Keeper) ([]*types.KeyPair, error) {

	pS := keeper.GetContractPublicState(ctx, addr)

	return pS, nil
}
func queryContractPublicStateForAccount(ctx sdk.Context, addr sdk.AccAddress, acc sdk.AccAddress, keeper Keeper) ([]*types.KeyPair, error) {
	pSA := keeper.GetContractPublicStateForAccount(ctx, addr, acc)

	return pSA, nil
}
func queryContractPublicStateByKey(ctx sdk.Context, addr sdk.AccAddress, key string, keeper Keeper) (types.KeyPair, error) {
	//var res sdk.Result
	kv := keeper.GetContractPublicStateByKey(ctx, addr, []byte(key))

	return kv, nil
}

func queryContractListByCode(ctx sdk.Context, codeID uint64, keeper Keeper) ([]types.ContractInfoWithAddress, error) {
	var contracts []types.ContractInfoWithAddress
	keeper.IterateContractInfo(ctx, func(addr sdk.AccAddress, info types.ContractInfo) bool {
		if info.CodeID == codeID {
			// and add the address
			infoWithAddress := types.ContractInfoWithAddress{
				Address:      addr,
				ContractInfo: &info,
			}
			contracts = append(contracts, infoWithAddress)
		}
		return false
	})

	// now we sort them by AbsoluteTxPosition
	sort.Slice(contracts, func(i, j int) bool {
		return contracts[i].ContractInfo.Created.LessThan(contracts[j].ContractInfo.Created)
	})
	// and remove that info for the final json (yes, the json:"-" tag doesn't work)
	for i := range contracts {
		contracts[i].Created = nil
	}

	return contracts, nil
}

func QueryCode(ctx sdk.Context, codeID uint64, keeper Keeper) (*types.QueryCodeResponse, error) {
	if codeID == 0 {
		return nil, nil
	}
	res, err := keeper.GetCodeInfo(ctx, codeID)

	if err != nil {
		// nil, nil leads to 404 in rest handler
		return nil, nil
	}
	info := types.CodeInfoResponse{
		CodeID:          codeID,
		Creator:         res.Creator,
		CodeHash:        res.CodeHash,
		Source:          res.Source,
		Builder:         res.Builder,
		DefaultDuration: res.DefaultDuration,
		Title:           res.Title,
		Description:     res.Description,
		Instances:       res.Instances,
	}

	code, err := keeper.GetByteCode(ctx, codeID)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "loading wasm code")
	}

	return &types.QueryCodeResponse{CodeInfoResponse: &info, Data: code}, nil
}

func queryCodeList(ctx sdk.Context, keeper Keeper) ([]types.CodeInfoResponse, error) {
	var info []types.CodeInfoResponse
	keeper.IterateCodeInfos(ctx, func(i uint64, res types.CodeInfo) bool {
		info = append(info, types.CodeInfoResponse{
			CodeID:          i,
			Creator:         res.Creator,
			CodeHash:        res.CodeHash,
			Source:          res.Source,
			Builder:         res.Builder,
			DefaultDuration: res.DefaultDuration,
			Title:           res.Title,
			Description:     res.Description,
			Instances:       res.Instances,
		})
		return false
	})
	return info, nil
}

/*
func queryContractHistory(ctx sdk.Context, bech string, keeper Keeper) ([]byte, error) {
	contractAddr, err := sdk.AccAddressFromBech32(bech)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	entries := keeper.GetContractHistory(ctx, contractAddr)
	if entries == nil {
		// nil, nil leads to 404 in rest handler
		return nil, nil
	}
	// redact response
	for i := range entries {
		entries[i].Updated = nil
	}

	bz, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}
*/

func queryContractAddress(ctx sdk.Context, contractId string, keeper Keeper) (sdk.AccAddress, error) {
	res := keeper.GetContractAddress(ctx, contractId)
	if res == nil {
		return nil, nil
	}

	return res, nil
}

func queryCodeHash(ctx sdk.Context, codeID uint64, keeper Keeper) ([]byte, error) {
	res := keeper.GetCodeHash(ctx, codeID)
	if res == nil {
		return nil, nil
	}

	return res, nil
}

func queryContractKey(ctx sdk.Context, address sdk.AccAddress, keeper Keeper) ([]byte, error) {
	res := keeper.GetContractKey(ctx, address)
	if res == nil {
		return nil, nil
	}

	return res, nil
}

func queryContractHash(ctx sdk.Context, address sdk.AccAddress, keeper Keeper) ([]byte, error) {
	res := keeper.GetContractInfo(ctx, address)
	if res == nil {
		return nil, nil
	}
	info, err := keeper.GetCodeInfo(ctx, res.CodeID)
	if err != nil {
		return nil, nil
	}
	return info.CodeHash, nil
}
