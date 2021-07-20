package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/tax/x/tax/types"
)

func TestMsgCreateFixedAmountPlan(t *testing.T) {
	taxPoolAddr := sdk.AccAddress(crypto.AddressHash([]byte("taxPoolAddr")))
	stakingCoinWeights := sdk.NewDecCoins(
		sdk.DecCoin{Denom: "testFarmStakingCoinDenom", Amount: sdk.MustNewDecFromStr("1.0")},
	)
	// needs to be deterministic for test
	startTime, _ := time.Parse(time.RFC3339, "2021-11-01T22:08:41+00:00")
	endTime := startTime.AddDate(1, 0, 0)

	testCases := []struct {
		expectedErr string
		msg         *types.MsgCreateFixedAmountPlan
	}{
		{
			"", // empty means no error expected
			types.NewMsgCreateFixedAmountPlan(
				taxPoolAddr, stakingCoinWeights, startTime,
				endTime, sdk.Coins{sdk.NewCoin("uatom", sdk.NewInt(1))},
			),
		},
		{
			"invalid tax pool address \"\": empty address string is not allowed: invalid address",
			types.NewMsgCreateFixedAmountPlan(
				sdk.AccAddress{}, stakingCoinWeights, startTime,
				endTime, sdk.Coins{sdk.NewCoin("uatom", sdk.NewInt(1))},
			),
		},
		{
			"end time 2020-11-01 22:08:41 +0000 +0000 must be greater than start time 2021-11-01 22:08:41 +0000 +0000: invalid plan end time",
			types.NewMsgCreateFixedAmountPlan(
				taxPoolAddr, stakingCoinWeights, startTime,
				startTime.AddDate(-1, 0, 0), sdk.Coins{sdk.NewCoin("uatom", sdk.NewInt(1))},
			),
		},
		{
			"staking coin weights must not be empty",
			types.NewMsgCreateFixedAmountPlan(
				taxPoolAddr, sdk.NewDecCoins(), startTime,
				endTime, sdk.Coins{sdk.NewCoin("uatom", sdk.NewInt(1))},
			),
		},
		{
			"epoch amount must not be empty",
			types.NewMsgCreateFixedAmountPlan(
				taxPoolAddr, stakingCoinWeights, startTime,
				endTime, sdk.Coins{},
			),
		},
	}

	for _, tc := range testCases {
		require.IsType(t, &types.MsgCreateFixedAmountPlan{}, tc.msg)
		require.Equal(t, types.TypeMsgCreateFixedAmountPlan, tc.msg.Type())
		require.Equal(t, types.RouterKey, tc.msg.Route())
		require.Equal(t, sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(tc.msg)), tc.msg.GetSignBytes())

		err := tc.msg.ValidateBasic()
		if tc.expectedErr == "" {
			require.Nil(t, err)
			signers := tc.msg.GetSigners()
			require.Len(t, signers, 1)
			require.Equal(t, tc.msg.GetCreator(), signers[0])
		} else {
			require.EqualError(t, err, tc.expectedErr)
		}
	}
}

func TestMsgCreateRatioPlan(t *testing.T) {
	taxPoolAddr := sdk.AccAddress(crypto.AddressHash([]byte("taxPoolAddr")))
	stakingCoinWeights := sdk.NewDecCoins(
		sdk.DecCoin{Denom: "testFarmStakingCoinDenom", Amount: sdk.MustNewDecFromStr("1.0")},
	)
	// needs to be deterministic for test
	startTime, _ := time.Parse(time.RFC3339, "2021-11-01T22:08:41+00:00")
	endTime := startTime.AddDate(1, 0, 0)

	testCases := []struct {
		expectedErr string
		msg         *types.MsgCreateRatioPlan
	}{
		{
			"", // empty means no error expected
			types.NewMsgCreateRatioPlan(
				taxPoolAddr, stakingCoinWeights, startTime,
				endTime, sdk.NewDec(1),
			),
		},
		{
			"invalid tax pool address \"\": empty address string is not allowed: invalid address",
			types.NewMsgCreateRatioPlan(
				sdk.AccAddress{}, stakingCoinWeights, startTime,
				endTime, sdk.NewDec(1),
			),
		},
		{
			"end time 2020-11-01 22:08:41 +0000 +0000 must be greater than start time 2021-11-01 22:08:41 +0000 +0000: invalid plan end time",
			types.NewMsgCreateRatioPlan(
				taxPoolAddr, stakingCoinWeights, startTime,
				startTime.AddDate(-1, 0, 0), sdk.NewDec(1),
			),
		},
		{
			"staking coin weights must not be empty",
			types.NewMsgCreateRatioPlan(
				taxPoolAddr, sdk.NewDecCoins(), startTime,
				endTime, sdk.NewDec(1),
			),
		},
		{
			"invalid plan epoch ratio",
			types.NewMsgCreateRatioPlan(
				taxPoolAddr, stakingCoinWeights, startTime,
				endTime, sdk.NewDec(-1),
			),
		},
	}

	for _, tc := range testCases {
		require.IsType(t, &types.MsgCreateRatioPlan{}, tc.msg)
		require.Equal(t, types.TypeMsgCreateRatioPlan, tc.msg.Type())
		require.Equal(t, types.RouterKey, tc.msg.Route())
		require.Equal(t, sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(tc.msg)), tc.msg.GetSignBytes())

		err := tc.msg.ValidateBasic()
		if tc.expectedErr == "" {
			require.Nil(t, err)
			signers := tc.msg.GetSigners()
			require.Len(t, signers, 1)
			require.Equal(t, tc.msg.GetCreator(), signers[0])
		} else {
			require.EqualError(t, err, tc.expectedErr)
		}
	}
}

func TestMsgStake(t *testing.T) {
	taxPoolAddr := sdk.AccAddress(crypto.AddressHash([]byte("taxPoolAddr")))
	stakingCoins := sdk.NewCoins(
		sdk.NewCoin("testFarmStakingCoinDenom", sdk.NewInt(1)),
	)

	testCases := []struct {
		expectedErr string
		msg         *types.MsgStake
	}{
		{
			"", // empty means no error expected
			types.NewMsgStake(taxPoolAddr, stakingCoins),
		},
		// TODO" not implemented yet
	}

	for _, tc := range testCases {
		require.IsType(t, &types.MsgStake{}, tc.msg)
		require.Equal(t, types.TypeMsgStake, tc.msg.Type())
		require.Equal(t, types.RouterKey, tc.msg.Route())
		require.Equal(t, sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(tc.msg)), tc.msg.GetSignBytes())

		err := tc.msg.ValidateBasic()
		if tc.expectedErr == "" {
			require.Nil(t, err)
			signers := tc.msg.GetSigners()
			require.Len(t, signers, 1)
			require.Equal(t, tc.msg.GetFarmer(), signers[0])
		} else {
			require.EqualError(t, err, tc.expectedErr)
		}
	}
}

func TestMsgUnstake(t *testing.T) {
	taxPoolAddr := sdk.AccAddress(crypto.AddressHash([]byte("taxPoolAddr")))
	stakingCoins := sdk.NewCoins(
		sdk.NewCoin("testFarmStakingCoinDenom", sdk.NewInt(1)),
	)

	testCases := []struct {
		expectedErr string
		msg         *types.MsgUnstake
	}{
		{
			"", // empty means no error expected
			types.NewMsgUnstake(taxPoolAddr, stakingCoins),
		},
		// TODO" not implemented yet
	}

	for _, tc := range testCases {
		require.IsType(t, &types.MsgUnstake{}, tc.msg)
		require.Equal(t, types.TypeMsgUnstake, tc.msg.Type())
		require.Equal(t, types.RouterKey, tc.msg.Route())
		require.Equal(t, sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(tc.msg)), tc.msg.GetSignBytes())

		err := tc.msg.ValidateBasic()
		if tc.expectedErr == "" {
			require.Nil(t, err)
			signers := tc.msg.GetSigners()
			require.Len(t, signers, 1)
			require.Equal(t, tc.msg.GetFarmer(), signers[0])
		} else {
			require.EqualError(t, err, tc.expectedErr)
		}
	}
}

func TestMsgHarvest(t *testing.T) {
	stakingCoinDenoms := []string{""}
	taxPoolAddr := sdk.AccAddress(crypto.AddressHash([]byte("taxPoolAddr")))

	testCases := []struct {
		expectedErr string
		msg         *types.MsgHarvest
	}{
		{
			"", // empty means no error expected
			types.NewMsgHarvest(taxPoolAddr, stakingCoinDenoms),
		},
		// TODO" not implemented yet
	}

	for _, tc := range testCases {
		require.IsType(t, &types.MsgHarvest{}, tc.msg)
		require.Equal(t, types.TypeMsgHarvest, tc.msg.Type())
		require.Equal(t, types.RouterKey, tc.msg.Route())
		require.Equal(t, sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(tc.msg)), tc.msg.GetSignBytes())

		err := tc.msg.ValidateBasic()
		if tc.expectedErr == "" {
			require.Nil(t, err)
			signers := tc.msg.GetSigners()
			require.Len(t, signers, 1)
			require.Equal(t, tc.msg.GetFarmer(), signers[0])
		} else {
			require.EqualError(t, err, tc.expectedErr)
		}
	}
}
