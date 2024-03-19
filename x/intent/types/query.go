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

// NewQueryActionsForOwnerRequest creates and returns a new QueryActionsForOwnerFromAddressRequest
func NewQueryActionsForOwnerRequest(owner string, pagination *query.PageRequest) *QueryActionsForOwnerRequest {
	return &QueryActionsForOwnerRequest{
		Owner:      owner,
		Pagination: pagination,
	}
}

// NewQueryActionsForOwnerResponse creates and returns a new QueryActionsForOwnerFromAddressResponse
func NewQueryActionsForOwnerResponse(autoTxInfos []ActionInfo) *QueryActionsForOwnerResponse {
	return &QueryActionsForOwnerResponse{
		ActionInfos: autoTxInfos,
	}
}

// NewQueryActionsForOwnerRequest creates and returns a new QueryActionsForOwnerFromAddressRequest
func NewQueryActionsRequest(pagination *query.PageRequest) *QueryActionsRequest {
	return &QueryActionsRequest{
		Pagination: pagination,
	}
}

// NewQueryActionsForOwnerRequest creates and returns a new QueryActionsForOwnerFromAddressRequest
func NewQueryActionRequest(id string) *QueryActionRequest {
	return &QueryActionRequest{Id: id}
}

// NewQueryActionsForOwnerResponse creates and returns a new QueryActionsForOwnerFromAddressResponse
func NewQueryActionsResponse(autoTxInfos []ActionInfo) *QueryActionsResponse {
	return &QueryActionsResponse{
		ActionInfos: autoTxInfos,
	}
}
