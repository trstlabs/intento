package keeper

import (
	"context"
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
func NewQuerier(keeper Keeper) grpcQuerier {
	return grpcQuerier{keeper: keeper}
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
		KeyPairs: rsp,
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

func (q grpcQuerier) SmartContractState(c context.Context, req *types.QuerySmartContractStateRequest) (*types.QuerySmartContractStateResponse, error) {
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(c).WithGasMeter(sdk.NewGasMeter(q.keeper.queryGasLimit))
	rsp, err := q.keeper.QuerySmart(ctx, addr, req.QueryData, false)
	switch {
	case err != nil:
		return nil, err
	case rsp == nil:
		return nil, types.ErrNotFound
	}
	return &types.QuerySmartContractStateResponse{Data: rsp}, nil

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

func (q grpcQuerier) AddressByLabel(c context.Context, req *types.QueryContractAddressByContractIdRequest) (*types.QueryContractAddressByContractIdResponse, error) {
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
	//var res sdk.Result
	pS := keeper.GetContractPublicState(ctx, addr)

	return pS, nil
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
	res := keeper.GetCodeInfo(ctx, codeID)
	if res == nil {
		// nil, nil leads to 404 in rest handler
		return nil, nil
	}
	info := types.CodeInfoResponse{
		CodeID:           codeID,
		Creator:          res.Creator,
		CodeHash:         res.CodeHash,
		Source:           res.Source,
		Builder:          res.Builder,
		ContractDuration: res.Duration,
		Title:            res.Title,
		Description:      res.Description,
		Instances:        res.Instances,
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
			CodeID:           i,
			Creator:          res.Creator,
			CodeHash:         res.CodeHash,
			Source:           res.Source,
			Builder:          res.Builder,
			ContractDuration: res.Duration,
			Title:            res.Title,
			Description:      res.Description,
			Instances:        res.Instances,
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
	return keeper.GetCodeInfo(ctx, res.CodeID).CodeHash, nil
}
