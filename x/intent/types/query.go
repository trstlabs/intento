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
func NewQueryActionsForOwnerResponse(actionInfos []ActionInfo) *QueryActionsForOwnerResponse {
	return &QueryActionsForOwnerResponse{
		ActionInfos: actionInfos,
	}
}

// NewQueryActionsForOwnerRequest creates and returns a new QueryActionsForOwnerFromAddressRequest
func NewQueryActionsRequest(pagination *query.PageRequest) *QueryActionsRequest {
	return &QueryActionsRequest{
		Pagination: pagination,
	}
}

// NewQueryActionsForOwnerResponse creates and returns a new QueryActionsForOwnerFromAddressResponse
func NewQueryActionsResponse(actionInfos []ActionInfo) *QueryActionsResponse {
	return &QueryActionsResponse{
		ActionInfos: actionInfos,
	}
}

// NewQueryActionRequest creates and returns a new QueryActionRequest
func NewQueryActionRequest(id string) *QueryActionRequest {
	return &QueryActionRequest{Id: id}
}

// NewQueryActionHistoryRequest creates and returns a new QueryActionHistoryRequest
func NewQueryActionHistoryRequest(id string) *QueryActionHistoryRequest {
	return &QueryActionHistoryRequest{Id: id}
}

// NewQueryHostedAccountRequest creates and returns a new QueryHostedAccountsRequest
func NewQueryHostedAccountRequest(address string) *QueryHostedAccountRequest {
	return &QueryHostedAccountRequest{Address: address}
}

// NewQueryHostedAccountResponse creates and returns a new QueryHostedAccountsResponse
func NewQueryHostedAccountResponse(hostedAccount HostedAccount) *QueryHostedAccountResponse {
	return &QueryHostedAccountResponse{
		HostedAccount: hostedAccount,
	}
}

// NewQueryHostedAccountRequest creates and returns a new QueryHostedAccountsRequest
func NewQueryHostedAccountsRequest(pagination *query.PageRequest) *QueryHostedAccountsRequest {
	return &QueryHostedAccountsRequest{
		Pagination: pagination,
	}
}

// NewQueryHostedAccountResponse creates and returns a new QueryHostedAccountsResponse
func NewQueryHostedAccountsResponse(hostedAccounts []HostedAccount) *QueryHostedAccountsResponse {
	return &QueryHostedAccountsResponse{
		HostedAccounts: hostedAccounts,
	}
}
