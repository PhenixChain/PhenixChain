package rest

import "github.com/PhenixChain/PhenixChain/x/auth"

// BroadcastReq requests broadcasting a transaction
type BroadcastReq struct {
	Tx     auth.StdTx `json:"tx"`
	Return string     `json:"return"`
}
