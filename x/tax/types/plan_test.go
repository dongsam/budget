package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/tax/x/tax/types"
)

func TestGetPoolInformation(t *testing.T) {
	commonTerminationAcc := sdk.AccAddress([]byte("terminationAddr"))
	commonStartTime := time.Now().UTC()
	commonEndTime := commonStartTime.AddDate(1, 0, 0)
	commonCoinWeights := sdk.NewDecCoins(
		sdk.DecCoin{Denom: "testFarmStakingCoinDenom", Amount: sdk.MustNewDecFromStr("1.0")},
	)

	testCases := []struct {
		planId          uint64
		planType        types.PlanType
		taxPoolAddr string
		rewardPoolAddr  string
		terminationAddr string
		reserveAddr     string
		coinWeights     sdk.DecCoins
	}{
		{
			planId:          uint64(1),
			planType:        types.PlanTypePublic,
			taxPoolAddr: sdk.AccAddress([]byte("taxPoolAddr1")).String(),
			rewardPoolAddr:  "cosmos1yqurgw7xa94psk95ctje76ferlddg8vykflaln6xsgarj5w6jkrsuvh9dj",
			reserveAddr:     "cosmos18f2zl0q0gpexruasqzav2vfwdthl4779gtmdxgqdpdl03sq9eygq42ff0u",
		},
	}

	for _, tc := range testCases {
		planName := types.PlanName(tc.planId, tc.planType, tc.taxPoolAddr)
		rewardPoolAcc := types.GenerateRewardPoolAcc(planName)
		basePlan := types.NewBasePlan(tc.planId, tc.planType, tc.taxPoolAddr, commonTerminationAcc.String(), commonCoinWeights, commonStartTime, commonEndTime)
		require.Equal(t, basePlan.RewardPoolAddress, rewardPoolAcc.String())
	}
}
