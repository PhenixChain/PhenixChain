package tx

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"

	"github.com/PhenixChain/PhenixChain/client"
	"github.com/PhenixChain/PhenixChain/client/context"
	"github.com/PhenixChain/PhenixChain/client/utils"
	"github.com/PhenixChain/PhenixChain/codec"
	"github.com/PhenixChain/PhenixChain/types/rest"
	"github.com/PhenixChain/PhenixChain/x/auth"
)

type (
	// EncodeReq defines a tx encoding request.
	EncodeReq struct {
		Tx auth.StdTx `json:"tx"`
	}

	// EncodeResp defines a tx encoding response.
	EncodeResp struct {
		Tx string `json:"tx"`
	}
)

// EncodeTxRequestHandlerFn returns the encode tx REST handler. In particular,
// it takes a json-formatted transaction, encodes it to the Amino wire protocol,
// and responds with base64-encoded bytes.
func EncodeTxRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req EncodeReq

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = cdc.UnmarshalJSON(body, &req)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// nikolas
		// Re-encode it to the wire protocol
		txBytes, err := cliCtx.Codec.MarshalJSON(req.Tx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// base64 encode the encoded tx bytes
		txBytesBase64 := base64.StdEncoding.EncodeToString(txBytes)

		response := EncodeResp{Tx: txBytesBase64}
		rest.PostProcessResponse(w, cdc, response, cliCtx.Indent)
	}
}

// txEncodeRespStr implements a simple Stringer wrapper for a encoded tx.
type txEncodeRespStr string

func (txr txEncodeRespStr) String() string {
	return string(txr)
}

// GetEncodeCommand returns the encode command to take a JSONified transaction and turn it into
// Amino-serialized bytes
func GetEncodeCommand(codec *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "encode [file]",
		Short: "Encode transactions generated offline",
		Long: `Encode transactions created with the --generate-only flag and signed with the sign command.
Read a transaction from <file>, serialize it to the Amino wire protocol, and output it as base64.
If you supply a dash (-) argument in place of an input filename, the command reads from standard input.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cliCtx := context.NewCLIContext().WithCodec(codec)

			stdTx, err := utils.ReadStdTxFromFile(cliCtx.Codec, args[0])
			if err != nil {
				return
			}

			// nikolas
			txBytes, err := cliCtx.Codec.MarshalJSON(stdTx)
			if err != nil {
				return err
			}

			// base64 encode the encoded tx bytes
			txBytesBase64 := base64.StdEncoding.EncodeToString(txBytes)

			response := txEncodeRespStr(txBytesBase64)
			cliCtx.PrintOutput(response) // nolint:errcheck

			return nil
		},
	}

	return client.PostCommands(cmd)[0]
}
