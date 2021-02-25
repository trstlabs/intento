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
	//registerTxHandlers(clientCtx, r)


}

func registerQueryRoutes(clientCtx client.Context, r *mux.Router) {
	// this line is used by starport scaffolding # 3
	//r.HandleFunc("/tpp/estimators/{id}", getEstimatorHandler(clientCtx)).Methods("GET")
	r.HandleFunc("/tpp/estimators", listEstimatorHandler(clientCtx)).Methods("GET")

	//r.HandleFunc("/tpp/buyers/{id}", getBuyerHandler(clientCtx)).Methods("GET")
	r.HandleFunc("/tpp/buyers", listBuyerHandler(clientCtx)).Methods("GET")

	//r.HandleFunc("/tpp/items/{id}", getItemHandler(clientCtx)).Methods("GET")
	r.HandleFunc("/tpp/items", listItemHandler(clientCtx)).Methods("GET")

}

/*func registerTxHandlers(clientCtx client.Context, r *mux.Router) {
	// this line is used by starport scaffolding # 4
	r.HandleFunc("/tpp/estimators", createEstimatorHandler(clientCtx)).Methods("POST")
	r.HandleFunc("/tpp/estimators/{id}", updateEstimatorHandler(clientCtx)).Methods("POST")
	r.HandleFunc("/tpp/estimators/{id}", deleteEstimatorHandler(clientCtx)).Methods("POST")

	r.HandleFunc("/tpp/buyers", createBuyerHandler(clientCtx)).Methods("POST")
	r.HandleFunc("/tpp/buyers/{id}", updateBuyerHandler(clientCtx)).Methods("POST")
	r.HandleFunc("/tpp/buyers/{id}", deleteBuyerHandler(clientCtx)).Methods("POST")

	r.HandleFunc("/tpp/items", createItemHandler(clientCtx)).Methods("POST")
	r.HandleFunc("/tpp/items/{id}", updateItemHandler(clientCtx)).Methods("POST")
	r.HandleFunc("/tpp/items/{id}", deleteItemHandler(clientCtx)).Methods("POST")

}*/
