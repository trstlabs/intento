package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	// this line is used by starport scaffolding # 1
)

const (
	MethodGet = "GET"
)

// RegisterRoutes registers trst-related REST handlers to a router
func RegisterRoutes(clientCtx client.Context, r *mux.Router) {
	// this line is used by starport scaffolding # 2
	registerQueryRoutes(clientCtx, r)
	registerTxHandlers(clientCtx, r)

}

func registerQueryRoutes(clientCtx client.Context, r *mux.Router) {
	// this line is used by starport scaffolding # 3
	r.HandleFunc("/trst/estimator/{id}", getEstimatorHandler(clientCtx)).Methods("GET")
	r.HandleFunc("/trst/estimator", listEstimatorHandler(clientCtx)).Methods("GET")

	//	r.HandleFunc("/trst/buyer/{id}", getBuyerHandler(clientCtx)).Methods("GET")
	//	r.HandleFunc("/trst/buyer", listBuyerHandler(clientCtx)).Methods("GET")
	r.HandleFunc("/trst/buyeritems/{buyer}", buyerItemsHandler(clientCtx)).Methods("GET")

	r.HandleFunc("/trst/item/{id}", getItemHandler(clientCtx)).Methods("GET")
	r.HandleFunc("/trst/item", listItemHandler(clientCtx)).Methods("GET")
	r.HandleFunc("/trst/listeditems", listListedItemsHandler(clientCtx)).Methods("GET")
	r.HandleFunc("/trst/selleritems/{seller}", sellerItemsHandler(clientCtx)).Methods("GET")

}

func registerTxHandlers(clientCtx client.Context, r *mux.Router) {
	// this line is used by starport scaffolding # 4
	r.HandleFunc("/trst/estimator", createEstimatorHandler(clientCtx)).Methods("POST")
	r.HandleFunc("/trst/estimator/{id}", updateEstimatorHandler(clientCtx)).Methods("POST")
	r.HandleFunc("/trst/estimator/{id}", deleteEstimatorHandler(clientCtx)).Methods("POST")

	r.HandleFunc("/trst/buyer", createBuyerHandler(clientCtx)).Methods("POST")

	r.HandleFunc("/trst/buyer/{id}", deleteBuyerHandler(clientCtx)).Methods("POST")

	r.HandleFunc("/trst/item", createItemHandler(clientCtx)).Methods("POST")

	r.HandleFunc("/trst/item/{id}", deleteItemHandler(clientCtx)).Methods("POST")
}
