package rest

import (
	"github.com/PhenixChain/PhenixChain/client/context"
	"github.com/PhenixChain/PhenixChain/codec"
	"github.com/gorilla/mux"
)

// RegisterRoutes registers minting module REST handlers on the provided router.
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	registerQueryRoutes(cliCtx, r, cdc)
}
