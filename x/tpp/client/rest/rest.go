package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	// this line is used by starport scaffolding # 1
)

const (
	MethodGet = "GET"
)

// RegisterRoutes registers tpp-related REST handlers to a router
func RegisterRoutes(clientCtx client.Context, r *mux.Router) {
	// this line is used by starport scaffolding # 2
	registerQueryRoutes(clientCtx, r)
	registerTxHandlers(clientCtx, r)


}

func registerQueryRoutes(clientCtx client.Context, r *mux.Router) {
	// this line is used by starport scaffolding # 3
	r.HandleFunc("/tpp/estimator/{id}", getEstimatorHandler(clientCtx)).Methods("GET")
	r.HandleFunc("/tpp/estimator", listEstimatorHandler(clientCtx)).Methods("GET")

	r.HandleFunc("/tpp/buyer/{id}", getBuyerHandler(clientCtx)).Methods("GET")
	r.HandleFunc("/tpp/buyer", listBuyerHandler(clientCtx)).Methods("GET")

	r.HandleFunc("/tpp/item/{id}", getItemHandler(clientCtx)).Methods("GET")
	r.HandleFunc("/tpp/item", listItemHandler(clientCtx)).Methods("GET")
	r.HandleFunc("/tpp/inactiveitems", listInactiveItemsHandler(clientCtx)).Methods("GET")

}

func registerTxHandlers(clientCtx client.Context, r *mux.Router) {
	// this line is used by starport scaffolding # 4
	r.HandleFunc("/tpp/estimator", createEstimatorHandler(clientCtx)).Methods("POST")
	r.HandleFunc("/tpp/estimator/{id}", updateEstimatorHandler(clientCtx)).Methods("POST")
	r.HandleFunc("/tpp/estimator/{id}", deleteEstimatorHandler(clientCtx)).Methods("POST")


	r.HandleFunc("/tpp/buyer", createBuyerHandler(clientCtx)).Methods("POST")
	r.HandleFunc("/tpp/buyer/{id}", updateBuyerHandler(clientCtx)).Methods("POST")
	r.HandleFunc("/tpp/buyer/{id}", deleteBuyerHandler(clientCtx)).Methods("POST")
	
	r.HandleFunc("/tpp/item", createItemHandler(clientCtx)).Methods("POST")
	r.HandleFunc("/tpp/item/{id}", updateItemHandler(clientCtx)).Methods("POST")
	r.HandleFunc("/tpp/item/{id}", deleteItemHandler(clientCtx)).Methods("POST")
}
