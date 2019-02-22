package simulation

import (
	sdk "github.com/PhenixChain/PhenixChain/types"
)

// TODO Any invariants to check here?
// AllInvariants tests all slashing invariants
func AllInvariants() sdk.Invariant {
	return func(_ sdk.Context) error {
		return nil
	}
}
