package keeper_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/tax/app"
	"github.com/tendermint/tax/x/tax/types"
)

func TestGetSetNewPlan(t *testing.T) {
	simapp, ctx := createTestApp(true)

	taxPoolAddr := sdk.AccAddress([]byte("taxPoolAddr"))
	terminationAddr := sdk.AccAddress([]byte("terminationAddr"))
	stakingCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1000000)))
	coinWeights := sdk.NewDecCoins(
		sdk.DecCoin{Denom: "testFarmStakingCoinDenom", Amount: sdk.MustNewDecFromStr("1.0")},
	)

	addrs := app.AddTestAddrs(simapp, ctx, 2, sdk.NewInt(2000000))
	farmerAddr := addrs[0]

	startTime := time.Now().UTC()
	endTime := startTime.AddDate(1, 0, 0)
	basePlan := types.NewBasePlan(1, 1, taxPoolAddr.String(), terminationAddr.String(), coinWeights, startTime, endTime)
	fixedPlan := types.NewFixedAmountPlan(basePlan, sdk.NewCoins(sdk.NewCoin("testFarmCoinDenom", sdk.NewInt(1000000))))
	simapp.TaxKeeper.SetPlan(ctx, fixedPlan)

	planGet, found := simapp.TaxKeeper.GetPlan(ctx, 1)
	require.True(t, found)
	require.Equal(t, fixedPlan, planGet)

	taxes := simapp.TaxKeeper.GetAllTaxes(ctx)
	require.Len(t, taxes, 1)
	require.Equal(t, fixedPlan, taxes[0])

	// TODO: tmp test codes for testing functionality, need to separated
	err := simapp.TaxKeeper.Stake(ctx, farmerAddr, stakingCoins)
	require.NoError(t, err)

	stakings := simapp.TaxKeeper.GetAllStakings(ctx)
	fmt.Println(stakings)
	stakingByFarmer, found := simapp.TaxKeeper.GetStakingByFarmer(ctx, farmerAddr)
	stakingsByDenom := simapp.TaxKeeper.GetStakingsByStakingCoinDenom(ctx, sdk.DefaultBondDenom)

	require.True(t, found)
	require.Equal(t, stakings[0], stakingByFarmer)
	require.Equal(t, stakings, stakingsByDenom)

	simapp.TaxKeeper.SetReward(ctx, sdk.DefaultBondDenom, farmerAddr, stakingCoins)

	rewards := simapp.TaxKeeper.GetAllRewards(ctx)
	rewardsByFarmer := simapp.TaxKeeper.GetRewardsByFarmer(ctx, farmerAddr)
	rewardsByDenom := simapp.TaxKeeper.GetRewardsByStakingCoinDenom(ctx, sdk.DefaultBondDenom)

	require.Equal(t, rewards, rewardsByFarmer)
	require.Equal(t, rewards, rewardsByDenom)
}
