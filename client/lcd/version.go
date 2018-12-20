package lcd

import (
	"net/http"

	"github.com/PhenixChain/PhenixChain/client/context"
	"github.com/PhenixChain/PhenixChain/client/utils"
	"github.com/PhenixChain/PhenixChain/version"
)

// cli version REST handler endpoint
func CLIVersionRequestHandler(w http.ResponseWriter, r *http.Request) {
	v := version.GetVersion()
	w.Write([]byte(v))
}

// connected node version REST handler endpoint
func NodeVersionRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		version, err := cliCtx.Query("/app/version", nil)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(version)
	}
}
