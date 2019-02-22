package keeper

import (
	sdk "github.com/PhenixChain/PhenixChain/types"
)

// get outstanding rewards
func (k Keeper) GetOutstandingRewardsCoins(ctx sdk.Context) sdk.DecCoins {
	return k.GetOutstandingRewards(ctx)
}

// get the community coins
func (k Keeper) GetFeePoolCommunityCoins(ctx sdk.Context) sdk.DecCoins {
	return k.GetFeePool(ctx).CommunityPool
}
