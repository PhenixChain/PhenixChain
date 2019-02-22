package rest

import (
	"io/ioutil"
	"net/http"

	"github.com/PhenixChain/PhenixChain/client/context"
	"github.com/PhenixChain/PhenixChain/codec"
	"github.com/PhenixChain/PhenixChain/types/rest"
	"github.com/PhenixChain/PhenixChain/x/auth"
)

type broadcastBody struct {
	Tx auth.StdTx `json:"tx"`
}

// BroadcastTxRequestHandlerFn returns the broadcast tx REST handler
func BroadcastTxRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m broadcastBody
		if ok := unmarshalBodyOrReturnBadRequest(cliCtx, w, r, &m); !ok {
			return
		}

		// nikolas
		txBytes, err := cliCtx.Codec.MarshalJSON(m.Tx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		res, err := cliCtx.BroadcastTx(txBytes)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func unmarshalBodyOrReturnBadRequest(cliCtx context.CLIContext, w http.ResponseWriter, r *http.Request, m interface{}) bool {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return false
	}
	err = cliCtx.Codec.UnmarshalJSON(body, m)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return false
	}
	return true
}
