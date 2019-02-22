package keeper

import (
	sdk "github.com/PhenixChain/PhenixChain/types"
	"github.com/PhenixChain/PhenixChain/x/staking/types"
)

// Return all validators that a delegator is bonded to. If maxRetrieve is supplied, the respective amount will be returned.
func (k Keeper) GetDelegatorValidators(ctx sdk.Context, delegatorAddr sdk.AccAddress,
	maxRetrieve uint16) (validators []types.Validator) {
	validators = make([]types.Validator, maxRetrieve)

	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := GetDelegationsKey(delegatorAddr)
	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey) // smallest to largest
	defer iterator.Close()

	i := 0
	for ; iterator.Valid() && i < int(maxRetrieve); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())

		validator, found := k.GetValidator(ctx, delegation.ValidatorAddr)
		if !found {
			panic(types.ErrNoValidatorFound(types.DefaultCodespace))
		}
		validators[i] = validator
		i++
	}
	return validators[:i] // trim
}

// return a validator that a delegator is bonded to
func (k Keeper) GetDelegatorValidator(ctx sdk.Context, delegatorAddr sdk.AccAddress,
	validatorAddr sdk.ValAddress) (validator types.Validator, err sdk.Error) {

	delegation, found := k.GetDelegation(ctx, delegatorAddr, validatorAddr)
	if !found {
		return validator, types.ErrNoDelegation(types.DefaultCodespace)
	}

	validator, found = k.GetValidator(ctx, delegation.ValidatorAddr)
	if !found {
		panic(types.ErrNoValidatorFound(types.DefaultCodespace))
	}
	return
}

//_____________________________________________________________________________________

// return all delegations for a delegator
func (k Keeper) GetAllDelegatorDelegations(ctx sdk.Context, delegator sdk.AccAddress) (
	delegations []types.Delegation) {

	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := GetDelegationsKey(delegator)
	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey) //smallest to largest
	defer iterator.Close()

	i := 0
	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		delegations = append(delegations, delegation)
		i++
	}
	return delegations
}

// return all unbonding-delegations for a delegator
func (k Keeper) GetAllUnbondingDelegations(ctx sdk.Context, delegator sdk.AccAddress) (
	unbondingDelegations []types.UnbondingDelegation) {

	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := GetUBDsKey(delegator)
	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey) // smallest to largest
	defer iterator.Close()

	i := 0
	for ; iterator.Valid(); iterator.Next() {
		unbondingDelegation := types.MustUnmarshalUBD(k.cdc, iterator.Value())
		unbondingDelegations = append(unbondingDelegations, unbondingDelegation)
		i++
	}
	return unbondingDelegations
}

// return all redelegations for a delegator
func (k Keeper) GetAllRedelegations(ctx sdk.Context, delegator sdk.AccAddress,
	srcValAddress, dstValAddress sdk.ValAddress) (
	redelegations []types.Redelegation) {

	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := GetREDsKey(delegator)
	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey) // smallest to largest
	defer iterator.Close()

	srcValFilter := !(srcValAddress.Empty())
	dstValFilter := !(dstValAddress.Empty())

	for ; iterator.Valid(); iterator.Next() {
		redelegation := types.MustUnmarshalRED(k.cdc, iterator.Value())
		if srcValFilter && !(srcValAddress.Equals(redelegation.ValidatorSrcAddr)) {
			continue
		}
		if dstValFilter && !(dstValAddress.Equals(redelegation.ValidatorDstAddr)) {
			continue
		}
		redelegations = append(redelegations, redelegation)
	}
	return redelegations
}