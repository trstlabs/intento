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

// NewQueryFlowsForOwnerRequest creates and returns a new QueryFlowsForOwnerFromAddressRequest
func NewQueryFlowsForOwnerRequest(owner string, pagination *query.PageRequest) *QueryFlowsForOwnerRequest {
	return &QueryFlowsForOwnerRequest{
		Owner:      owner,
		Pagination: pagination,
	}
}

// NewQueryFlowsForOwnerResponse creates and returns a new QueryFlowsForOwnerFromAddressResponse
func NewQueryFlowsForOwnerResponse(flowInfos []FlowInfo) *QueryFlowsForOwnerResponse {
	return &QueryFlowsForOwnerResponse{
		FlowInfos: flowInfos,
	}
}

// NewQueryFlowsForOwnerRequest creates and returns a new QueryFlowsForOwnerFromAddressRequest
func NewQueryFlowsRequest(pagination *query.PageRequest) *QueryFlowsRequest {
	return &QueryFlowsRequest{
		Pagination: pagination,
	}
}

// NewQueryFlowsForOwnerResponse creates and returns a new QueryFlowsForOwnerFromAddressResponse
func NewQueryFlowsResponse(flowInfos []FlowInfo) *QueryFlowsResponse {
	return &QueryFlowsResponse{
		FlowInfos: flowInfos,
	}
}

// NewQueryFlowRequest creates and returns a new QueryFlowRequest
func NewQueryFlowRequest(id string) *QueryFlowRequest {
	return &QueryFlowRequest{Id: id}
}

// NewQueryFlowHistoryRequest creates and returns a new QueryFlowHistoryRequest
func NewQueryFlowHistoryRequest(id string) *QueryFlowHistoryRequest {
	return &QueryFlowHistoryRequest{Id: id}
}

// NewQueryTrustlessAgentRequest creates and returns a new QueryTrustlessAgentsRequest
func NewQueryTrustlessAgentRequest(agentAddress string) *QueryTrustlessAgentRequest {
	return &QueryTrustlessAgentRequest{AgentAddress: agentAddress}
}

// NewQueryTrustlessAgentResponse creates and returns a new QueryTrustlessAgentsResponse
func NewQueryTrustlessAgentResponse(trustlessAgent TrustlessAgent) *QueryTrustlessAgentResponse {
	return &QueryTrustlessAgentResponse{
		TrustlessAgent: trustlessAgent,
	}
}

// NewQueryTrustlessAgentRequest creates and returns a new QueryTrustlessAgentsRequest
func NewQueryTrustlessAgentsRequest(pagination *query.PageRequest) *QueryTrustlessAgentsRequest {
	return &QueryTrustlessAgentsRequest{
		Pagination: pagination,
	}
}

// NewQueryTrustlessAgentResponse creates and returns a new QueryTrustlessAgentsResponse
func NewQueryTrustlessAgentsResponse(trustlessAgents []TrustlessAgent) *QueryTrustlessAgentsResponse {
	return &QueryTrustlessAgentsResponse{
		TrustlessAgents: trustlessAgents,
	}
}
