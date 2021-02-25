
package rest
/*
import (
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/danieljdd/tpp/x/tpp/types"
	"github.com/gorilla/mux"
)

// Used to not have an error if strconv is unused
var _ = strconv.Itoa(42)

type createEstimatorRequest struct {
	BaseReq                 rest.BaseReq `json:"base_req"`
	Creator                 string       `json:"creator"`
	Estimation              string       `json:"estimation"`
	Estimatorestimationhash string       `json:"estimatorestimationhash"`
	Itemid                  string       `json:"itemid"`
	Deposit                 string       `json:"deposit"`
	Interested              string       `json:"interested"`
	Comment                 string       `json:"comment"`
	Flag                    string       `json:"flag"`
}

func createEstimatorHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createEstimatorRequest
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		_, err := sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		parsedEstimation := req.Estimation

		parsedEstimatorestimationhash := req.Estimatorestimationhash

		parsedItemid := req.Itemid

		parsedDeposit := sdk.ParseCoinNormalized(req.Deposit)

		parsedInterested := req.Interested

		parsedComment := req.Comment

		parsedFlag := req.Flag

		msg := types.NewMsgCreateEstimator(
			req.Creator,
			parsedEstimation,
			parsedEstimatorestimationhash,
			parsedItemid,
			parsedDeposit,
			parsedInterested,
			parsedComment,
			parsedFlag,
		)

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

type updateEstimatorRequest struct {
	BaseReq                 rest.BaseReq `json:"base_req"`
	Creator                 string       `json:"creator"`
	Estimation              string       `json:"estimation"`
	Estimatorestimationhash string       `json:"estimatorestimationhash"`
	Itemid                  string       `json:"itemid"`
	Deposit                 string       `json:"deposit"`
	Interested              string       `json:"interested"`
	Comment                 string       `json:"comment"`
	Flag                    string       `json:"flag"`
}

func updateEstimatorHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		var req updateEstimatorRequest
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		_, err := sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		parsedEstimation := req.Estimation

		parsedEstimatorestimationhash := req.Estimatorestimationhash

		parsedItemid := req.Itemid

		parsedDeposit := sdk.ParseCoinNormalized(req.Deposit)

		parsedInterested := req.Interested

		parsedComment := req.Comment

		parsedFlag := req.Flag

		msg := types.NewMsgUpdateEstimator(
			req.Creator,
			id,
			parsedEstimation,
			parsedEstimatorestimationhash,
			parsedItemid,
			parsedDeposit,
			parsedInterested,
			parsedComment,
			parsedFlag,
		)

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

type deleteEstimatorRequest struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Creator string       `json:"creator"`
}

func deleteEstimatorHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		var req deleteEstimatorRequest
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		_, err := sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgDeleteEstimator(
			req.Creator,
			id,
		)

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}
*/