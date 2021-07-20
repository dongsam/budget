package keeper

import (
	gogotypes "github.com/gogo/protobuf/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/tax/x/tax/types"
)

// NewPlan sets the next plan number to a given plan interface
func (k Keeper) NewPlan(ctx sdk.Context, plan types.PlanI) types.PlanI {
	if err := plan.SetId(k.GetNextPlanIDWithUpdate(ctx)); err != nil {
		panic(err)
	}

	return plan
}

// GetPlan implements PlanI.
func (k Keeper) GetPlan(ctx sdk.Context, id uint64) (plan types.PlanI, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetPlanKey(id))
	if bz == nil {
		return plan, false
	}

	return k.decodePlan(bz), true
}

// GetAllTaxes returns all taxes in the Keeper.
func (k Keeper) GetAllTaxes(ctx sdk.Context) (taxes []types.PlanI) {
	k.IterateAllTaxes(ctx, func(plan types.PlanI) (stop bool) {
		taxes = append(taxes, plan)
		return false
	})

	return taxes
}

// SetPlan implements PlanI.
func (k Keeper) SetPlan(ctx sdk.Context, plan types.PlanI) {
	id := plan.GetId()
	store := ctx.KVStore(k.storeKey)

	bz, err := k.MarshalPlan(plan)
	if err != nil {
		panic(err)
	}

	store.Set(types.GetPlanKey(id), bz)
}

// RemovePlan removes an plan for the plan mapper store.
// NOTE: this will cause supply invariant violation if called
func (k Keeper) RemovePlan(ctx sdk.Context, plan types.PlanI) {
	id := plan.GetId()
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetPlanKey(id))
}

// IterateAllTaxes iterates over all the stored taxes and performs a callback function.
// Stops iteration when callback returns true.
func (k Keeper) IterateAllTaxes(ctx sdk.Context, cb func(plan types.PlanI) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PlanKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		plan := k.decodePlan(iterator.Value())

		if cb(plan) {
			break
		}
	}
}

// GetTaxesByFarmerAddrIndex reads from kvstore and return a specific Plan indexed by given farmer address
func (k Keeper) GetTaxesByFarmerAddrIndex(ctx sdk.Context, farmerAcc sdk.AccAddress) (taxes []types.PlanI) {
	k.IterateTaxesByFarmerAddr(ctx, farmerAcc, func(plan types.PlanI) bool {
		taxes = append(taxes, plan)
		return false
	})

	return taxes
}

// IterateTaxesByFarmerAddr iterates over all the stored taxes and performs a callback function.
// Stops iteration when callback returns true.
func (k Keeper) IterateTaxesByFarmerAddr(ctx sdk.Context, farmerAcc sdk.AccAddress, cb func(plan types.PlanI) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetTaxesByFarmerAddrIndexKey(farmerAcc))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		planID := gogotypes.UInt64Value{}

		err := k.cdc.Unmarshal(iterator.Value(), &planID)
		if err != nil {
			panic(err)
		}
		plan, _ := k.GetPlan(ctx, planID.GetValue())
		if cb(plan) {
			break
		}
	}
}

// SetPlanIdByFarmerAddrIndex sets Index by FarmerAddr
// TODO: need to gas cost check for existing check or update everytime
func (k Keeper) SetPlanIdByFarmerAddrIndex(ctx sdk.Context, farmerAcc sdk.AccAddress, planID uint64) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&gogotypes.UInt64Value{Value: planID})
	store.Set(types.GetPlanByFarmerAddrIndexKey(farmerAcc, planID), b)
}

// CreateFixedAmountPlan sets fixed amount plan.
func (k Keeper) CreateFixedAmountPlan(ctx sdk.Context, msg *types.MsgCreateFixedAmountPlan, typ types.PlanType) *types.FixedAmountPlan {
	nextId := k.GetNextPlanIDWithUpdate(ctx)
	taxPoolAddr := msg.TaxPoolAddress
	terminationAddr := taxPoolAddr

	basePlan := types.NewBasePlan(
		nextId,
		typ,
		taxPoolAddr,
		terminationAddr,
		msg.StakingCoinWeights,
		msg.StartTime,
		msg.EndTime,
	)

	fixedPlan := types.NewFixedAmountPlan(basePlan, msg.EpochAmount)

	k.SetPlan(ctx, fixedPlan)

	return fixedPlan
}

// CreateRatioPlan sets ratio plan.
func (k Keeper) CreateRatioPlan(ctx sdk.Context, msg *types.MsgCreateRatioPlan, typ types.PlanType) *types.RatioPlan {
	nextId := k.GetNextPlanIDWithUpdate(ctx)
	taxPoolAddr := msg.TaxPoolAddress
	terminationAddr := taxPoolAddr

	basePlan := types.NewBasePlan(
		nextId,
		typ,
		taxPoolAddr,
		terminationAddr,
		msg.StakingCoinWeights,
		msg.StartTime,
		msg.EndTime,
	)

	ratioPlan := types.NewRatioPlan(basePlan, msg.EpochRatio)

	k.SetPlan(ctx, ratioPlan)

	return ratioPlan
}
