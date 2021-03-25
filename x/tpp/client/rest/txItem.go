package rest

import (
	//"crypto/sha256"
	//"encoding/hex"
	///"fmt"
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

type createItemRequest struct {
	BaseReq     rest.BaseReq `json:"base_req"`
	Creator     string       `json:"creator"`
	Title       string       `json:"title"`
	Description string       `json:"description"`

	Shippingcost    int64    `json:"shippingcost"`
	Localpickup     bool     `json:"localpickup"`
	Estimationcount int64    `json:"estimationcount"`
	Tags            []string `json:"tags"`
	Condition       int64    `json:"condition"`
	Shippingregion  []string `json:"shippingregion"`
	Depositamount   int64    `json:"Depositamount"`
}

func createItemHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createItemRequest
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

		parsedTitle := req.Title

		parsedDescription := req.Description

		parsedShippingcost := req.Shippingcost

		parsedLocalpickup := req.Localpickup

		parsedEstimationcount := req.Estimationcount
		//var estimationcount = fmt.Sprint(req.Estimationcount)
		//var estimationcountHash = sha256.Sum256([]byte(estimationcount))
		//var estimationcountHashString = hex.EncodeToString(estimationcountHash[:])

		parsedTags := req.Tags

		parsedCondition := req.Condition

		parsedShippingregion := req.Shippingregion

		parsedDepositAmount := req.Depositamount

		msg := types.NewMsgCreateItem(
			req.Creator,
			parsedTitle,
			parsedDescription,
			parsedShippingcost,
			parsedLocalpickup,
			parsedEstimationcount,

			parsedTags,

			parsedCondition,
			parsedShippingregion,
			parsedDepositAmount,
		)

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

type updateItemRequest struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Seller string       `json:"creator"`

	Shippingcost int64 `json:"shippingcost"`
	Localpickup  bool  `json:"localpickup"`

	Shippingregion []string `json:"shippingregion"`
}

func updateItemHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		var req updateItemRequest
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		_, err := sdk.AccAddressFromBech32(req.Seller)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		parsedShippingcost := req.Shippingcost

		parsedLocalpickup := req.Localpickup

		parsedShippingregion := req.Shippingregion

		msg := types.NewMsgUpdateItem(
			req.Seller,
			id,

			parsedShippingcost,
			parsedLocalpickup,

			parsedShippingregion,
		)

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

type deleteItemRequest struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Seller string       `json:"creator"`
}

func deleteItemHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		var req deleteItemRequest
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		_, err := sdk.AccAddressFromBech32(req.Seller)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgDeleteItem(
			req.Seller,
			id,
		)

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}
