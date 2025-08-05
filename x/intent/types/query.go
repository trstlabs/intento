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

// NewQueryTrustlessExecutionAgentRequest creates and returns a new QueryTrustlessExecutionAgentsRequest
func NewQueryTrustlessExecutionAgentRequest(agentAddress string) *QueryTrustlessExecutionAgentRequest {
	return &QueryTrustlessExecutionAgentRequest{AgentAddress: agentAddress}
}

// NewQueryTrustlessExecutionAgentResponse creates and returns a new QueryTrustlessExecutionAgentsResponse
func NewQueryTrustlessExecutionAgentResponse(trustlessExecutionAgent TrustlessExecutionAgent) *QueryTrustlessExecutionAgentResponse {
	return &QueryTrustlessExecutionAgentResponse{
		TrustlessExecutionAgent: trustlessExecutionAgent,
	}
}

// NewQueryTrustlessExecutionAgentRequest creates and returns a new QueryTrustlessExecutionAgentsRequest
func NewQueryTrustlessExecutionAgentsRequest(pagination *query.PageRequest) *QueryTrustlessExecutionAgentsRequest {
	return &QueryTrustlessExecutionAgentsRequest{
		Pagination: pagination,
	}
}

// NewQueryTrustlessExecutionAgentResponse creates and returns a new QueryTrustlessExecutionAgentsResponse
func NewQueryTrustlessExecutionAgentsResponse(trustlessExecutionAgents []TrustlessExecutionAgent) *QueryTrustlessExecutionAgentsResponse {
	return &QueryTrustlessExecutionAgentsResponse{
		TrustlessExecutionAgents: trustlessExecutionAgents,
	}
}
