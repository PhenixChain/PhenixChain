package rest

import (
	"github.com/gorilla/mux"

	"github.com/PhenixChain/PhenixChain/client/context"
	"github.com/PhenixChain/PhenixChain/codec"
	"github.com/PhenixChain/PhenixChain/crypto/keys"
)

// RegisterRoutes registers staking-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, kb keys.Keybase) {
	registerQueryRoutes(cliCtx, r, cdc)
	registerTxRoutes(cliCtx, r, cdc, kb)
}
