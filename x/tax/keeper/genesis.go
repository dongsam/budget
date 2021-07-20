package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/tax/x/tax/types"
)

// InitGenesis initializes the tax module's state from a given genesis state.
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	if err := k.ValidateGenesis(ctx, genState); err != nil {
		panic(err)
	}

	k.SetParams(ctx, genState.Params)
	moduleAcc := k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
	k.accountKeeper.SetModuleAccount(ctx, moduleAcc)
	// TODO: unimplemented
	//for _, record := range genState.PlanRecords {
	//	k.SetPlanRecord(ctx, record)
	//}
	//for _, staking := range genState.Stakings {
	//	k.SetStaking(ctx, staking)
	//}
	//for _, reward := range genState.Rewards {
	//	k.SetReward(ctx, reword)
	//}
}

// ExportGenesis returns the tax module's genesis state.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	params := k.GetParams(ctx)

	// TODO: unimplemented
	var planRecords []types.PlanRecord

	//taxes := k.GetAllTaxes(ctx)
	//stakings := k.GetAllStakings(ctx)
	//rewards := k.GetAllRewards(ctx)

	//for _, plan := range taxes {
	//	record, found := k.GetPlanRecord(ctx, plan)
	//	if found {
	//		planRecords = append(planRecords, record)
	//	}
	//}
	//
	//if len(planRecords) == 0 {
	//	planRecords = []types.PlanRecord{}
	//}

	return types.NewGenesisState(params, planRecords, nil, nil)
}

// ValidateGenesis validates the tax module's genesis state.
func (k Keeper) ValidateGenesis(ctx sdk.Context, genState types.GenesisState) error {
	if err := genState.Params.Validate(); err != nil {
		return err
	}

	cc, _ := ctx.CacheContext()
	k.SetParams(cc, genState.Params)

	// TODO: unimplemented
	//for _, record := range genState.PlanRecords {
	//	record = k.SetPlanRecord(cc, record)
	//	if err := k.ValidatePlanRecord(cc, record); err != nil {
	//		return err
	//	}
	//}

	return nil
}
