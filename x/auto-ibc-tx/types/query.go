package types

import query "github.com/cosmos/cosmos-sdk/types/query"

// NewQueryInterchainAccountRequest creates and returns a new QueryInterchainAccountFromAddressRequest
func NewQueryInterchainAccountRequest(owner, connectionID string) *QueryInterchainAccountFromAddressRequest {
	return &QueryInterchainAccountFromAddressRequest{
		Owner:        owner,
		ConnectionId: connectionID,
	}
}

// NewQueryInterchainAccountResponse creates and returns a new QueryInterchainAccountFromAddressResponse
func NewQueryInterchainAccountResponse(interchainAccAddr string) *QueryInterchainAccountFromAddressResponse {
	return &QueryInterchainAccountFromAddressResponse{
		InterchainAccountAddress: interchainAccAddr,
	}
}

// NewQueryAutoTxsForOwnerRequest creates and returns a new QueryAutoTxsForOwnerFromAddressRequest
func NewQueryAutoTxsForOwnerRequest(owner string, pagination *query.PageRequest) *QueryAutoTxsForOwnerRequest {
	return &QueryAutoTxsForOwnerRequest{
		Owner:      owner,
		Pagination: pagination,
	}
}

// NewQueryAutoTxsForOwnerResponse creates and returns a new QueryAutoTxsForOwnerFromAddressResponse
func NewQueryAutoTxsForOwnerResponse(autoTxInfos []AutoTxInfo) *QueryAutoTxsForOwnerResponse {
	return &QueryAutoTxsForOwnerResponse{
		AutoTxInfos: autoTxInfos,
	}
}
