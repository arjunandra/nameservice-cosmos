package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/arjunandra/nameservice-cosmos/x/nameservice/internal/types"
)

// Declaring Request Structures For Transaction Commands (From /client/cli/tx.go)

type buyNameReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Name    string       `json:"name"`
	Amount  string       `json:"amount"`
	Buyer   string       `json:"buyer"`
}

type setNameReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Name    string       `json:"name"`
	Value   string       `json:"value"`
	Owner   string       `json:"owner"`
}

type deleteNameReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Name    string       `json:"name"`
	Owner   string       `json:"owner"`
}

// Defining Handlers For Transaction Commands (From /client/cli/tx.go)

func buyNameHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req buyNameReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Failed To Parse Request")
			return
		}

		baseReq := req.BaseReq.Sanitize()

		// State-less Checks
		if !baseReq.ValidateBasic(w) {
			return
		}

		// Retrieve Account
		addr, err := sdk.AccAddressFromBech32(req.Buyer)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Send Coins
		coins, err := sdk.ParseCoins(req.Amount)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Create Message
		msg := types.NewMsgBuyName(req.Name, coins, addr)

		// State-less Checks
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Generate Response
		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func setNameHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req setNameReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// Retrieve Address
		addr, err := sdk.AccAddressFromBech32(req.Owner)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Create Message
		msg := types.NewMsgSetName(req.Name, req.Value, addr)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Generate Response
		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func deleteNameHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req deleteNameReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}	

		// Retrieve Address
		addr, err := sdk.AccAddressFromBech32(req.Owner)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Create Message
		msg := types.NewMsgDeleteName(req.Name, addr)
		err = msg.ValidateBasic()
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Generate Response
		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}