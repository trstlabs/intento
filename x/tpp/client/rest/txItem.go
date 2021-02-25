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

type createItemRequest struct {
	BaseReq                     rest.BaseReq `json:"base_req"`
	Creator                     string       `json:"creator"`
	Title                       string       `json:"title"`
	Description                 string       `json:"description"`
	Shippingcost                string       `json:"shippingcost"`
	Localpickup                 string       `json:"localpickup"`
	Estimationcounthash         string       `json:"estimationcounthash"`
	Bestestimator               string       `json:"bestestimator"`
	Lowestestimator             string       `json:"lowestestimator"`
	Highestestimator            string       `json:"highestestimator"`
	Estimationprice             string       `json:"estimationprice"`
	Estimatorlist               string       `json:"estimatorlist"`
	Estimatorestimationhashlist string       `json:"estimatorestimationhashlist"`
	Transferable                string       `json:"transferable"`
	Buyer                       string       `json:"buyer"`
	Tracking                    string       `json:"tracking"`
	Status                      string       `json:"status"`
	Comments                    string       `json:"comments"`
	Tags                        string       `json:"tags"`
	Flags                       string       `json:"flags"`
	Condition                   string       `json:"condition"`
	Shippingregion              string       `json:"shippingregion"`
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

		parsedShippingcost := sdk.ParseCoinNormalized(req.Shippingcost)

		parsedLocalpickup := req.Localpickup

		parsedEstimationcounthash := req.Estimationcounthash

		parsedBestestimator := req.Bestestimator

		parsedLowestestimator := req.Lowestestimator

		parsedHighestestimator := req.Highestestimator

		parsedEstimationprice := req.Estimationprice

		parsedEstimatorlist := req.Estimatorlist

		parsedEstimatorestimationhashlist := req.Estimatorestimationhashlist

		parsedTransferable := req.Transferable

		parsedBuyer := req.Buyer

		parsedTracking := req.Tracking

		parsedStatus := req.Status

		parsedComments := req.Comments

		parsedTags := req.Tags

		parsedFlags := req.Flags

		parsedCondition := req.Condition

		parsedShippingregion := req.Shippingregion

		msg := types.NewMsgCreateItem(
			req.Creator,
			parsedTitle,
			parsedDescription,
			parsedShippingcost,
			parsedLocalpickup,
			parsedEstimationcounthash,
			parsedBestestimator,
			parsedLowestestimator,
			parsedHighestestimator,
			parsedEstimationprice,
			parsedEstimatorlist,
			parsedEstimatorestimationhashlist,
			parsedTransferable,
			parsedBuyer,
			parsedTracking,
			parsedStatus,
			parsedComments,
			parsedTags,
			parsedFlags,
			parsedCondition,
			parsedShippingregion,
		)

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

type updateItemRequest struct {
	BaseReq                     rest.BaseReq `json:"base_req"`
	Creator                     string       `json:"creator"`
	Title                       string       `json:"title"`
	Description                 string       `json:"description"`
	Shippingcost                string       `json:"shippingcost"`
	Localpickup                 string       `json:"localpickup"`
	Estimationcounthash         string       `json:"estimationcounthash"`
	Bestestimator               string       `json:"bestestimator"`
	Lowestestimator             string       `json:"lowestestimator"`
	Highestestimator            string       `json:"highestestimator"`
	Estimationprice             string       `json:"estimationprice"`
	Estimatorlist               string       `json:"estimatorlist"`
	Estimatorestimationhashlist string       `json:"estimatorestimationhashlist"`
	Transferable                string       `json:"transferable"`
	Buyer                       string       `json:"buyer"`
	Tracking                    string       `json:"tracking"`
	Status                      string       `json:"status"`
	Comments                    string       `json:"comments"`
	Tags                        string       `json:"tags"`
	Flags                       string       `json:"flags"`
	Condition                   string       `json:"condition"`
	Shippingregion              string       `json:"shippingregion"`
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

		_, err := sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		parsedTitle := req.Title

		parsedDescription := req.Description

		parsedShippingcost := sdk.ParseCoinNormalized(eq.Shippingcost)

		parsedLocalpickup := req.Localpickup

		parsedEstimationcounthash := req.Estimationcounthash

		parsedBestestimator := req.Bestestimator

		parsedLowestestimator := req.Lowestestimator

		parsedHighestestimator := req.Highestestimator

		parsedEstimationprice := req.Estimationprice

		parsedEstimatorlist := req.Estimatorlist

		parsedEstimatorestimationhashlist := req.Estimatorestimationhashlist

		parsedTransferable := req.Transferable

		parsedBuyer := req.Buyer

		parsedTracking := req.Tracking

		parsedStatus := req.Status

		parsedComments := req.Comments

		parsedTags := req.Tags

		parsedFlags := req.Flags

		parsedCondition := req.Condition

		parsedShippingregion := req.Shippingregion

		msg := types.NewMsgUpdateItem(
			req.Creator,
			id,
			parsedTitle,
			parsedDescription,
			parsedShippingcost,
			parsedLocalpickup,
			parsedEstimationcounthash,
			parsedBestestimator,
			parsedLowestestimator,
			parsedHighestestimator,
			parsedEstimationprice,
			parsedEstimatorlist,
			parsedEstimatorestimationhashlist,
			parsedTransferable,
			parsedBuyer,
			parsedTracking,
			parsedStatus,
			parsedComments,
			parsedTags,
			parsedFlags,
			parsedCondition,
			parsedShippingregion,
		)

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

type deleteItemRequest struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Creator string       `json:"creator"`
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

		_, err := sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgDeleteItem(
			req.Creator,
			id,
		)

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}
*/