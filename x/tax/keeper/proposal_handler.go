package keeper

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/tendermint/tax/x/tax/types"
)

// HandlePublicPlanProposal is a handler for executing a fixed amount plan creation proposal.
func HandlePublicPlanProposal(ctx sdk.Context, k Keeper, taxesAny []*codectypes.Any) error {
	taxes, err := types.UnpackTaxes(taxesAny)
	if err != nil {
		return err
	}

	for _, plan := range taxes {
		switch p := plan.(type) {
		case *types.FixedAmountPlan:
			msg := types.NewMsgCreateFixedAmountPlan(
				p.GetTaxPoolAddress(),
				p.GetStakingCoinWeights(),
				p.GetStartTime(),
				p.GetEndTime(),
				p.EpochAmount,
			)

			fixedPlan := k.CreateFixedAmountPlan(ctx, msg, types.PlanTypePublic)

			logger := k.Logger(ctx)
			logger.Info("created public fixed amount plan", "fixed_amount_plan", fixedPlan)

		case *types.RatioPlan:
			msg := types.NewMsgCreateRatioPlan(
				p.GetTaxPoolAddress(),
				p.GetStakingCoinWeights(),
				p.GetStartTime(),
				p.GetEndTime(),
				p.EpochRatio,
			)

			ratioPlan := k.CreateRatioPlan(ctx, msg, types.PlanTypePublic)

			logger := k.Logger(ctx)
			logger.Info("created public fixed amount plan", "ratio_plan", ratioPlan)

		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized tax proposal plan type: %T", p)
		}
	}

	return nil
}
