package keeper_test

import (
	"context"
	"time"

	"github.com/cometbft/cometbft/proto/tendermint/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	"github.com/spf13/cast"
	_ "github.com/stretchr/testify/suite"
	"github.com/trstlabs/intento/app/apptesting"

	//intenttypes "github.com/trstlabs/intento/x/intent/types"
	"github.com/trstlabs/intento/x/interchainquery/types"
)

type MsgSubmitQueryResponseTestCase struct {
	validMsg types.MsgSubmitQueryResponse
	goCtx    context.Context
	query    types.Query
}

// Given a connection ID, returns the light client height
func GetLightClientHeight(ibcKeeper ibckeeper.Keeper, ctx sdk.Context, connectionID string) (height uint64) {
	connection, found := ibcKeeper.ConnectionKeeper.GetConnection(ctx, connectionID)
	if !found {
		return 0
	}

	clientState, found := ibcKeeper.ClientKeeper.GetClientState(ctx, connection.ClientId)
	if !found {
		return 0
	}

	latestHeight, err := cast.ToUint64E(clientState.GetLatestHeight().GetRevisionHeight())
	if err != nil {
		return 0
	}
	return latestHeight
}

func (s *KeeperTestSuite) SetupMsgSubmitQueryResponse() MsgSubmitQueryResponseTestCase {
	// set up IBC
	s.CreateTransferChannel(apptesting.HostChainId)

	// define the query
	goCtx := sdk.WrapSDKContext(s.Ctx)
	h := GetLightClientHeight(*s.App.IBCKeeper, s.Ctx, s.TransferPath.EndpointA.ConnectionID)

	height := int64(h - 1) // start at the (LC height) - 1  height, which is the height the query executes at!
	result := []byte("result-example")
	proofOps := crypto.ProofOps{}
	fromAddress := s.TestAccs[0].String()
	expectedId := "9792c1d779a3846a8de7ae82f31a74d308b279a521fa9e0d5c4f08917117bf3e"

	_, addr, _ := bech32.DecodeAndConvert(s.TestAccs[0].String())
	data := banktypes.CreateAccountBalancesPrefix(addr)

	timeoutDuration := time.Minute
	query := types.Query{
		Id:               expectedId,
		CallbackId:       "withdrawalbalance",
		CallbackModule:   "intent",
		ChainId:          apptesting.HostChainId,
		ConnectionId:     s.TransferPath.EndpointA.ConnectionID,
		QueryType:        "store/bank", // intentionally leave off key to skip proof
		RequestData:      append(data, []byte(apptesting.HostChainId)...),
		TimeoutDuration:  timeoutDuration,
		TimeoutTimestamp: uint64(s.Ctx.BlockTime().Add(timeoutDuration).UnixNano()),
	}

	return MsgSubmitQueryResponseTestCase{
		validMsg: types.MsgSubmitQueryResponse{
			ChainId:     apptesting.HostChainId,
			QueryId:     expectedId,
			Result:      result,
			ProofOps:    &proofOps,
			Height:      height,
			FromAddress: fromAddress,
		},
		goCtx: goCtx,
		query: query,
	}
}

func (s *KeeperTestSuite) TestMsgSubmitQueryResponse_WrongProof() {
	tc := s.SetupMsgSubmitQueryResponse()

	tc.query.QueryType = types.BANK_STORE_QUERY_WITH_PROOF

	s.App.InterchainQueryKeeper.SetQuery(s.Ctx, tc.query)

	resp, err := s.GetMsgServer().SubmitQueryResponse(tc.goCtx, &tc.validMsg)
	s.Require().ErrorContains(err, "Unable to verify membership proof: proof cannot be empty")
	s.Require().Nil(resp)
}

func (s *KeeperTestSuite) TestMsgSubmitQueryResponse_UnknownId() {
	tc := s.SetupMsgSubmitQueryResponse()

	tc.query.Id = tc.query.Id + "INVALID_SUFFIX" // create an invalid query id
	s.App.InterchainQueryKeeper.SetQuery(s.Ctx, tc.query)

	resp, err := s.GetMsgServer().SubmitQueryResponse(tc.goCtx, &tc.validMsg)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Equal(&types.MsgSubmitQueryResponseResponse{}, resp)

	// check that the query is STILL in the store, as it should NOT be deleted because the query was not found
	_, found := s.App.InterchainQueryKeeper.GetQuery(s.Ctx, tc.query.Id)
	s.Require().True(found)
}

func (s *KeeperTestSuite) TestMsgSubmitQueryResponse_ProofStale() {
	tc := s.SetupMsgSubmitQueryResponse()

	// Set the submission time in the future
	tc.query.QueryType = types.BANK_STORE_QUERY_WITH_PROOF
	tc.query.SubmissionHeight = 100
	s.App.InterchainQueryKeeper.SetQuery(s.Ctx, tc.query)

	// Attempt to submit the response, it should fail because the response is stale
	_, err := s.GetMsgServer().SubmitQueryResponse(tc.goCtx, &tc.validMsg)
	s.Require().ErrorContains(err, "Query proof height (16) is older than the submission height (100)")
}

func (s *KeeperTestSuite) TestMsgSubmitQueryResponse_Timeout_RejectQuery() {
	tc := s.SetupMsgSubmitQueryResponse()

	// set timeout to be expired and set the policy to reject
	tc.query.TimeoutTimestamp = uint64(1)
	tc.query.TimeoutPolicy = types.TimeoutPolicy_REJECT_QUERY_RESPONSE
	s.App.InterchainQueryKeeper.SetQuery(s.Ctx, tc.query)

	_, err := s.GetMsgServer().SubmitQueryResponse(tc.goCtx, &tc.validMsg)
	s.Require().NoError(err)

	// check that the original query was deleted
	_, found := s.App.InterchainQueryKeeper.GetQuery(s.Ctx, tc.query.Id)
	s.Require().False(found, "original query should be removed")
}

func (s *KeeperTestSuite) TestMsgSubmitQueryResponse_Timeout_RetryQuery() {
	tc := s.SetupMsgSubmitQueryResponse()

	// set timeout to be expired and set the policy to retry
	tc.query.TimeoutTimestamp = uint64(1)
	tc.query.TimeoutPolicy = types.TimeoutPolicy_RETRY_QUERY_REQUEST
	s.App.InterchainQueryKeeper.SetQuery(s.Ctx, tc.query)

	_, err := s.GetMsgServer().SubmitQueryResponse(tc.goCtx, &tc.validMsg)
	s.Require().NoError(err)

	// check that the query original query was deleted,
	// but that a new one was created for the retry
	_, found := s.App.InterchainQueryKeeper.GetQuery(s.Ctx, tc.query.Id)
	s.Require().False(found, "original query should be removed")

	queries := s.App.InterchainQueryKeeper.AllQueries(s.Ctx)
	s.Require().Len(queries, 1, "there should be one new query")

	// Confirm original query attributes have not changed
	actualQuery := queries[0]
	s.Require().NotEqual(tc.query.Id, actualQuery.Id, "query ID")
	s.Require().Equal(tc.query.QueryType, actualQuery.QueryType, "query type")
	s.Require().Equal(tc.query.ConnectionId, actualQuery.ConnectionId, "query connection ID")
	s.Require().Equal(tc.query.CallbackModule, actualQuery.CallbackModule, "query callback module")
	s.Require().Equal(tc.query.CallbackData, actualQuery.CallbackData, "cquery allback data")
	s.Require().Equal(tc.query.TimeoutPolicy, actualQuery.TimeoutPolicy, "query timeout policy")
	s.Require().Equal(tc.query.TimeoutDuration, actualQuery.TimeoutDuration, "query timeout duration")

	// Confirm timeout was reset
	expectedTimeoutTimestamp := uint64(s.Ctx.BlockTime().Add(tc.query.TimeoutDuration).UnixNano())
	s.Require().Equal(expectedTimeoutTimestamp, actualQuery.TimeoutTimestamp, "timeout timestamp")
	s.Require().Equal(false, actualQuery.RequestSent, "request sent")
}

func (s *KeeperTestSuite) TestMsgSubmitQueryResponse_Timeout_ExecuteCallback() {
	tc := s.SetupMsgSubmitQueryResponse()

	// set timeout to be expired and set the policy to retry
	tc.query.TimeoutTimestamp = uint64(1)
	tc.query.TimeoutPolicy = types.TimeoutPolicy_EXECUTE_QUERY_CALLBACK
	s.App.InterchainQueryKeeper.SetQuery(s.Ctx, tc.query)

	// Rather than testing by executing the callback in its entirety,
	// check by invoking without the required mocked state and catching
	// the error that's thrown at the start of the callback
	_, err := s.GetMsgServer().SubmitQueryResponse(tc.goCtx, &tc.validMsg)
	s.Require().ErrorContains(err, "unable to determine balance from query response")
}

func (s *KeeperTestSuite) TestMsgSubmitQueryResponse_FindAndInvokeCallback() {
	tc := s.SetupMsgSubmitQueryResponse()

	s.App.InterchainQueryKeeper.SetQuery(s.Ctx, tc.query)

	// The withdrawal balance test is already covered in it's respective module
	// For this test, we just want to check that the callback function is invoked
	// To do this, we can just ignore the appropriate withdrawal balance callback
	// mocked state, and catch the expected error that happens at the beginning of
	// the callback
	_, err := s.GetMsgServer().SubmitQueryResponse(tc.goCtx, &tc.validMsg)
	s.Require().ErrorContains(err, "unable to determine balance from query response")
}
