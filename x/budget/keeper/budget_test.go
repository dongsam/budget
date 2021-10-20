package keeper_test

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params/types/proposal"

	"github.com/tendermint/budget/x/budget/types"
)

func (suite *KeeperTestSuite) TestCollectBudgets() {
	for _, tc := range []struct {
		name           string
		budgets        []types.Budget
		epochBlocks    uint32
		accAsserts     []sdk.AccAddress
		balanceAsserts []sdk.Coins
		expectErr      bool
	}{
		{
			"basic budgets case",
			suite.budgets[:4],
			types.DefaultEpochBlocks,
			[]sdk.AccAddress{
				suite.collectionAddrs[0],
				suite.collectionAddrs[1],
				suite.collectionAddrs[2],
				suite.collectionAddrs[3],
				suite.budgetSourceAddrs[0],
				suite.budgetSourceAddrs[1],
				suite.budgetSourceAddrs[2],
			},
			[]sdk.Coins{
				mustParseCoinsNormalized("500000000denom1,500000000denom2,500000000denom3,500000000stake"),
				mustParseCoinsNormalized("500000000denom1,500000000denom2,500000000denom3,500000000stake"),
				mustParseCoinsNormalized("1000000000denom1,1000000000denom2,1000000000denom3,1000000000stake"),
				{},
				{},
				{},
				mustParseCoinsNormalized("1000000000denom1,1000000000denom2,1000000000denom3,1000000000stake"),
			},
			false,
		},
		{
			"only expired budget case",
			[]types.Budget{suite.budgets[3]},
			types.DefaultEpochBlocks,
			[]sdk.AccAddress{
				suite.collectionAddrs[3],
				suite.budgetSourceAddrs[2],
			},
			[]sdk.Coins{
				{},
				mustParseCoinsNormalized("1000000000denom1,1000000000denom2,1000000000denom3,1000000000stake"),
			},
			false,
		},
		{
			"budget source has small balances case",
			suite.budgets[4:6],
			types.DefaultEpochBlocks,
			[]sdk.AccAddress{
				suite.collectionAddrs[0],
				suite.collectionAddrs[1],
				suite.budgetSourceAddrs[3],
			},
			[]sdk.Coins{
				mustParseCoinsNormalized("1denom2,1denom3,500000000stake"),
				mustParseCoinsNormalized("1denom2,1denom3,500000000stake"),
				mustParseCoinsNormalized("1denom1,1denom3"),
			},
			false,
		},
		{
			"none budgets case",
			nil,
			types.DefaultEpochBlocks,
			[]sdk.AccAddress{
				suite.collectionAddrs[0],
				suite.collectionAddrs[1],
				suite.collectionAddrs[2],
				suite.collectionAddrs[3],
				suite.budgetSourceAddrs[0],
				suite.budgetSourceAddrs[1],
				suite.budgetSourceAddrs[2],
				suite.budgetSourceAddrs[3],
			},
			[]sdk.Coins{
				{},
				{},
				{},
				{},
				mustParseCoinsNormalized("1000000000denom1,1000000000denom2,1000000000denom3,1000000000stake"),
				mustParseCoinsNormalized("1000000000denom1,1000000000denom2,1000000000denom3,1000000000stake"),
				mustParseCoinsNormalized("1000000000denom1,1000000000denom2,1000000000denom3,1000000000stake"),
				mustParseCoinsNormalized("1denom1,2denom2,3denom3,1000000000stake"),
			},
			false,
		},
		{
			"disabled budget epoch",
			nil,
			0,
			[]sdk.AccAddress{
				suite.collectionAddrs[0],
				suite.collectionAddrs[1],
				suite.collectionAddrs[2],
				suite.collectionAddrs[3],
				suite.budgetSourceAddrs[0],
				suite.budgetSourceAddrs[1],
				suite.budgetSourceAddrs[2],
				suite.budgetSourceAddrs[3],
			},
			[]sdk.Coins{
				{},
				{},
				{},
				{},
				mustParseCoinsNormalized("1000000000denom1,1000000000denom2,1000000000denom3,1000000000stake"),
				mustParseCoinsNormalized("1000000000denom1,1000000000denom2,1000000000denom3,1000000000stake"),
				mustParseCoinsNormalized("1000000000denom1,1000000000denom2,1000000000denom3,1000000000stake"),
				mustParseCoinsNormalized("1denom1,2denom2,3denom3,1000000000stake"),
			},
			false,
		},
		{
			"disabled budget epoch with budgets",
			suite.budgets[:4],
			0,
			[]sdk.AccAddress{
				suite.collectionAddrs[0],
				suite.collectionAddrs[1],
				suite.collectionAddrs[2],
				suite.collectionAddrs[3],
				suite.budgetSourceAddrs[0],
				suite.budgetSourceAddrs[1],
				suite.budgetSourceAddrs[2],
				suite.budgetSourceAddrs[3],
			},
			[]sdk.Coins{
				{},
				{},
				{},
				{},
				mustParseCoinsNormalized("1000000000denom1,1000000000denom2,1000000000denom3,1000000000stake"),
				mustParseCoinsNormalized("1000000000denom1,1000000000denom2,1000000000denom3,1000000000stake"),
				mustParseCoinsNormalized("1000000000denom1,1000000000denom2,1000000000denom3,1000000000stake"),
				mustParseCoinsNormalized("1denom1,2denom2,3denom3,1000000000stake"),
			},
			false,
		},
	} {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			params := suite.keeper.GetParams(suite.ctx)
			params.Budgets = tc.budgets
			params.EpochBlocks = tc.epochBlocks
			suite.keeper.SetParams(suite.ctx, params)

			err := suite.keeper.CollectBudgets(suite.ctx)
			if tc.expectErr {
				suite.Error(err)
			} else {
				suite.NoError(err)

				for i, acc := range tc.accAsserts {
					suite.True(suite.app.BankKeeper.GetAllBalances(suite.ctx, acc).IsEqual(tc.balanceAsserts[i]))
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestBudgetExpiration() {
	// TODO: not implemented
}


func (suite *KeeperTestSuite) TestBudgetChangeSituation() {
	budget := types.Budget{
		Name:                "budget1",
		Rate:                sdk.NewDecWithPrec(5, 2), // 5%
		BudgetSourceAddress: suite.budgetSourceAddrs[0].String(),
		CollectionAddress:   suite.collectionAddrs[0].String(),
		StartTime:           mustParseRFC3339("0000-01-01T00:00:00Z"),
		EndTime:             mustParseRFC3339("9999-12-31T00:00:00Z"),
	}

	params := suite.keeper.GetParams(suite.ctx)
	params.Budgets = []types.Budget{budget}
	suite.keeper.SetParams(suite.ctx, params)

	// TODO: beginblock budget -> mempool -> endblock gov paramchange
	b, _ := json.Marshal(suite.budgets[6])
	//&suite.govHandler()
	//suite.app.GovKeeper.SetVotingParams()
	//proposal := govtypes.Proposal
	//suite.app.GovKeeper.SetProposal()

	testCases := []struct {
		name     string
		proposal *proposal.ParameterChangeProposal
		onHandle func()
		expErr   bool
	}{
		//{
		//	"all fields",
		//	testProposal(proposal.NewParamChange(stakingtypes.ModuleName, string(stakingtypes.KeyMaxValidators), "1")),
		//	func() {
		//		maxVals := suite.app.StakingKeeper.MaxValidators(suite.ctx)
		//		suite.Require().Equal(uint32(1), maxVals)
		//	},
		//	false,
		//},
		//{
		//	"invalid type",
		//	testProposal(proposal.NewParamChange(stakingtypes.ModuleName, string(stakingtypes.KeyMaxValidators), "-")),
		//	func() {},
		//	true,
		//},

		{
			"add budget",
			testProposal(proposal.ParamChange{
				Subspace: types.ModuleName,
				Key:      string(types.KeyBudgets),
				// TODO: add legacyAmino codec for budget object in order to call subspace.Update
				//Value:    `{"name":"gravity-dex-farming-20213Q-20313Q","rate":"0.500000000000000000","budget_source_address":"cosmos17xpfvakm2amg962yls6f84z3kell8c5lserqta","collection_address":"cosmos1228ryjucdpdv3t87rxle0ew76a56ulvnfst0hq0sscd3nafgjpqqkcxcky","start_time":"2021-09-01T00:00:00Z","end_time":"2031-09-30T00:00:00Z"}`,
				Value:    string(b),
			}),
			func() {
				params := suite.keeper.GetParams(suite.ctx)
				suite.Require().Len(params.Budgets, 1)

				//depositParams := suite.app.GovKeeper.GetDepositParams(suite.ctx)
				//suite.Require().Equal(govtypes.DepositParams{
				//	MinDeposit:       sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(64000000))),
				//	MaxDepositPeriod: govtypes.DefaultPeriod,
				//}, depositParams)
			},
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			err := suite.govHandler(suite.ctx, tc.proposal)
			if tc.expErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				tc.onHandle()
			}
		})
	}


}

func (suite *KeeperTestSuite) TestGetSetTotalCollectedCoins() {
	collectedCoins := suite.keeper.GetTotalCollectedCoins(suite.ctx, "budget1")
	suite.Require().Nil(collectedCoins)

	suite.keeper.SetTotalCollectedCoins(suite.ctx, "budget1", sdk.NewCoins(sdk.NewInt64Coin(denom1, 1000000)))
	collectedCoins = suite.keeper.GetTotalCollectedCoins(suite.ctx, "budget1")
	suite.Require().True(coinsEq(sdk.NewCoins(sdk.NewInt64Coin(denom1, 1000000)), collectedCoins))

	suite.keeper.AddTotalCollectedCoins(suite.ctx, "budget1", sdk.NewCoins(sdk.NewInt64Coin(denom2, 1000000)))
	collectedCoins = suite.keeper.GetTotalCollectedCoins(suite.ctx, "budget1")
	suite.Require().True(coinsEq(sdk.NewCoins(sdk.NewInt64Coin(denom1, 1000000), sdk.NewInt64Coin(denom2, 1000000)), collectedCoins))

	suite.keeper.AddTotalCollectedCoins(suite.ctx, "budget2", sdk.NewCoins(sdk.NewInt64Coin(denom1, 1000000)))
	collectedCoins = suite.keeper.GetTotalCollectedCoins(suite.ctx, "budget2")
	suite.Require().True(coinsEq(sdk.NewCoins(sdk.NewInt64Coin(denom1, 1000000)), collectedCoins))
}

func (suite *KeeperTestSuite) TestTotalCollectedCoins() {
	budget := types.Budget{
		Name:                "budget1",
		Rate:                sdk.NewDecWithPrec(5, 2), // 5%
		BudgetSourceAddress: suite.budgetSourceAddrs[0].String(),
		CollectionAddress:   suite.collectionAddrs[0].String(),
		StartTime:           mustParseRFC3339("0000-01-01T00:00:00Z"),
		EndTime:             mustParseRFC3339("9999-12-31T00:00:00Z"),
	}

	params := suite.keeper.GetParams(suite.ctx)
	params.Budgets = []types.Budget{budget}
	suite.keeper.SetParams(suite.ctx, params)

	balance := suite.app.BankKeeper.GetAllBalances(suite.ctx, suite.budgetSourceAddrs[0])
	expectedCoins, _ := sdk.NewDecCoinsFromCoins(balance...).MulDec(sdk.NewDecWithPrec(5, 2)).TruncateDecimal()

	collectedCoins := suite.keeper.GetTotalCollectedCoins(suite.ctx, "budget1")
	suite.Require().Equal(sdk.Coins(nil), collectedCoins)

	suite.ctx = suite.ctx.WithBlockTime(mustParseRFC3339("2021-08-31T00:00:00Z"))
	err := suite.keeper.CollectBudgets(suite.ctx)
	suite.Require().NoError(err)

	collectedCoins = suite.keeper.GetTotalCollectedCoins(suite.ctx, "budget1")
	suite.Require().True(coinsEq(expectedCoins, collectedCoins))
}
