package rest

import (
	"encoding/json"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/trstlabs/trst/x/compute/internal/types"
)

type StoreCodeProposalJSONReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

	Title       string    `json:"title" yaml:"title"`
	Description string    `json:"description" yaml:"description"`
	Proposer    string    `json:"proposer" yaml:"proposer"`
	Deposit     sdk.Coins `json:"deposit" yaml:"deposit"`

	RunAs string `json:"run_as" yaml:"run_as"`
	// WASMByteCode can be raw or gzip compressed
	WASMByteCode []byte `json:"wasm_byte_code" yaml:"wasm_byte_code"`
}

func (s StoreCodeProposalJSONReq) Content() govtypes.Content {
	return &types.StoreCodeProposal{
		Title:        s.Title,
		Description:  s.Description,
		RunAs:        s.RunAs,
		WASMByteCode: s.WASMByteCode,
	}
}
func (s StoreCodeProposalJSONReq) GetProposer() string {
	return s.Proposer
}
func (s StoreCodeProposalJSONReq) GetDeposit() sdk.Coins {
	return s.Deposit
}
func (s StoreCodeProposalJSONReq) GetBaseReq() rest.BaseReq {
	return s.BaseReq
}

func StoreCodeProposalHandler(cliCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "wasm_store_code",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			var req StoreCodeProposalJSONReq
			if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
				return
			}
			toStdTxResponse(cliCtx, w, req)
		},
	}
}

type InstantiateProposalJSONReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`

	Proposer string    `json:"proposer" yaml:"proposer"`
	Deposit  sdk.Coins `json:"deposit" yaml:"deposit"`

	//RunAs      string          `json:"run_as" yaml:"run_as"`
	Code       uint64          `json:"code_id" yaml:"code_id"`
	ContractId string          `json:"contract_id" yaml:"contract_id"`
	InitMsg    json.RawMessage `json:"init_msg" yaml:"init_msg"`
	AutoMsg    json.RawMessage `json:"auto_msg" yaml:"auto_msg"`
	Funds      sdk.Coins       `json:"funds" yaml:"funds"`
}

func (s InstantiateProposalJSONReq) Content() govtypes.Content {
	return &types.InstantiateContractProposal{
		Title:       s.Title,
		Description: s.Description,
		//RunAs:       s.RunAs,
		//Proposer:   s.Proposer,
		CodeID:     s.Code,
		ContractId: s.ContractId,
		InitMsg:    []byte(s.InitMsg),
		Funds:      s.Funds,
	}
}
func (s InstantiateProposalJSONReq) GetProposer() string {
	return s.Proposer
}
func (s InstantiateProposalJSONReq) GetDeposit() sdk.Coins {
	return s.Deposit
}
func (s InstantiateProposalJSONReq) GetBaseReq() rest.BaseReq {
	return s.BaseReq
}

func InstantiateProposalHandler(cliCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "wasm_instantiate",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			var req InstantiateProposalJSONReq
			if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
				return
			}
			toStdTxResponse(cliCtx, w, req)
		},
	}
}

type ExecuteProposalJSONReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`

	Proposer string    `json:"proposer" yaml:"proposer"`
	Deposit  sdk.Coins `json:"deposit" yaml:"deposit"`

	Contract string          `json:"contract" yaml:"contract"`
	Msg      json.RawMessage `json:"msg" yaml:"msg"`
	// RunAs is the role that is passed to the contract's environment
	//RunAs     string    `json:"run_as" yaml:"run_as"`
	Funds sdk.Coins `json:"funds" yaml:"funds"`
}

func (s ExecuteProposalJSONReq) Content() govtypes.Content {
	return &types.ExecuteContractProposal{
		Title:       s.Title,
		Description: s.Description,
		Contract:    s.Contract,
		Msg:         []byte(s.Msg),
		//RunAs:       s.RunAs,
		//Proposer:  s.Proposer,
		Funds: s.Funds,
	}
}
func (s ExecuteProposalJSONReq) GetProposer() string {
	return s.Proposer
}
func (s ExecuteProposalJSONReq) GetDeposit() sdk.Coins {
	return s.Deposit
}
func (s ExecuteProposalJSONReq) GetBaseReq() rest.BaseReq {
	return s.BaseReq
}
func ExecuteProposalHandler(cliCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "wasm_execute",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			var req ExecuteProposalJSONReq
			if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
				return
			}
			toStdTxResponse(cliCtx, w, req)
		},
	}
}

type wasmProposalData interface {
	Content() govtypes.Content
	GetProposer() string
	GetDeposit() sdk.Coins
	GetBaseReq() rest.BaseReq
}

func toStdTxResponse(cliCtx client.Context, w http.ResponseWriter, data wasmProposalData) {
	proposerAddr, err := sdk.AccAddressFromBech32(data.GetProposer())
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	msg, err := govtypes.NewMsgSubmitProposal(data.Content(), data.GetDeposit(), proposerAddr)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := msg.ValidateBasic(); err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	baseReq := data.GetBaseReq().Sanitize()
	if !baseReq.ValidateBasic(w) {
		return
	}
	tx.WriteGeneratedTxResponse(cliCtx, w, baseReq, msg)
}
