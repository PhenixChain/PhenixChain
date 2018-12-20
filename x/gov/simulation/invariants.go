package simulation

import (
	sdk "github.com/PhenixChain/PhenixChain/types"
	"github.com/PhenixChain/PhenixChain/x/mock/simulation"
)

// AllInvariants tests all governance invariants
func AllInvariants() simulation.Invariant {
	return func(ctx sdk.Context) error {
		// TODO Add some invariants!
		// Checking proposal queues, no passed-but-unexecuted proposals, etc.
		return nil
	}
}
