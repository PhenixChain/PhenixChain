package ibc

import (
	"github.com/PhenixChain/PhenixChain/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(IBCTransferMsg{}, "cosmos-sdk/IBCTransferMsg", nil)
	cdc.RegisterConcrete(IBCReceiveMsg{}, "cosmos-sdk/IBCReceiveMsg", nil)
}
