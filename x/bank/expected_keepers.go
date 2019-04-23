package bank

import (
	sdk "github.com/PhenixChain/PhenixChain/types"
)

// expected crisis keeper
type CrisisKeeper interface {
	RegisterRoute(moduleName, route string, invar sdk.Invariant)
}
