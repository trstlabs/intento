package rest

import (
	//"crypto/sha256"
	//"encoding/hex"
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
	BaseReq     rest.BaseReq `json:"base_req"`
	Estimator   string       `json:"creator"`
	Estimatemsg []byte       `json:"estimatemsg"`
	Itemid      uint64       `json:"itemid"`
	Deposit     int64        `json:"deposit"`
	Interested  bool         `json:"interested"`
	Comment     string       `json:"comment"`
	//Flag                    string       `json:"flag"`
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

		_, err := sdk.AccAddressFromBech32(req.Estimator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		parsedMsg := req.Estimatemsg

		//parsedFlag := req.Flag
		parsedItemid := req.Itemid
		/*parsedItemid, e := strconv.ParseUint(req.Itemid, 10, 64)
		if e != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, e.Error())
			return
		}*/
		parsedInterested := req.Interested
		//parsedComment := req.Comment

		//parsedDeposit := req.Deposit

		//var estimatorestimation = strconv.FormatInt(parsedEstimation, 10)
		//var estimatorestimationhash = sha256.Sum256([]byte(estimatorestimation + req.Estimator))
		//var estimatorestimationhashstring = hex.EncodeToString(estimatorestimationhash[:])

		//depositamount := "5tpp"
		deposit := req.Deposit

		msg := types.NewMsgCreateEstimation(
			req.Estimator,
			parsedMsg,

			parsedItemid,
			deposit,
			parsedInterested,
		//	parsedComment,
		)

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

type updateEstimatorRequest struct {
	BaseReq    rest.BaseReq `json:"base_req"`
	Estimator  string       `json:"creator"`
	Itemid     uint64       `json:"itemid"`
	Interested bool         `json:"interested"`
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

		_, err := sdk.AccAddressFromBech32(req.Estimator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		parsedItemid, e := strconv.ParseUint(id, 10, 64)
		if e != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, e.Error())
			return
		}

		parsedInterested := req.Interested

		msg := types.NewMsgUpdateLike(
			req.Estimator,

			parsedItemid,

			parsedInterested,
		)

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

type deleteEstimatorRequest struct {
	BaseReq   rest.BaseReq `json:"base_req"`
	Estimator string       `json:"creator"`
}

func deleteEstimatorHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, e := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
		if e != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, e.Error())
			return
		}

		var req deleteEstimatorRequest
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		_, err := sdk.AccAddressFromBech32(req.Estimator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgDeleteEstimation(
			req.Estimator,
			id,
		)

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}
