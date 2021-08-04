package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/tax/x/tax/types"
)

func (suite *KeeperTestSuite) TestTaxCollection() {
	taxes1 := []types.Tax{
		{
			Name:                  "test",
			Rate:                  sdk.MustNewDecFromStr("0.5"),
			CollectionAddress:     suite.collectionAddrs[0].String(),
			CollectionAccountName: "",
			TaxSourceAddress:      suite.taxSourceAddrs[0].String(),
			TaxSourceAccountName:  "",
			StartTime:             mustParseRFC3339("0000-01-01T00:00:00Z"),
			EndTime:               mustParseRFC3339("9999-12-31T00:00:00Z"),
		},
		{
			Name:                  "test2",
			Rate:                  sdk.MustNewDecFromStr("0.5"),
			CollectionAddress:     suite.collectionAddrs[1].String(),
			CollectionAccountName: "",
			TaxSourceAddress:      suite.taxSourceAddrs[0].String(),
			TaxSourceAccountName:  "",
			StartTime:             mustParseRFC3339("0000-01-01T00:00:00Z"),
			EndTime:               mustParseRFC3339("9999-12-31T00:00:00Z"),
		},
		{
			Name:                  "test3",
			Rate:                  sdk.MustNewDecFromStr("1.0"),
			CollectionAddress:     suite.collectionAddrs[2].String(),
			CollectionAccountName: "",
			TaxSourceAddress:      suite.taxSourceAddrs[1].String(),
			TaxSourceAccountName:  "",
			StartTime:             mustParseRFC3339("0000-01-01T00:00:00Z"),
			EndTime:               mustParseRFC3339("9999-12-31T00:00:00Z"),
		},
		{
			Name:                  "test4",
			Rate:                  sdk.MustNewDecFromStr("1"),
			CollectionAddress:     suite.collectionAddrs[3].String(),
			CollectionAccountName: "",
			TaxSourceAddress:      suite.taxSourceAddrs[2].String(),
			TaxSourceAccountName:  "",
			StartTime:             mustParseRFC3339("0000-01-01T00:00:00Z"),
			EndTime:               mustParseRFC3339("0000-01-01T00:00:00Z"),
		},
	}

	taxes3 := []types.Tax{
		{
			Name:                  "test",
			Rate:                  sdk.MustNewDecFromStr("0.5"),
			CollectionAddress:     suite.collectionAddrs[0].String(),
			CollectionAccountName: "",
			TaxSourceAddress:      suite.taxSourceAddrs[3].String(),
			TaxSourceAccountName:  "",
			StartTime:             mustParseRFC3339("0000-01-01T00:00:00Z"),
			EndTime:               mustParseRFC3339("9999-12-31T00:00:00Z"),
		},
		{
			Name:                  "test2",
			Rate:                  sdk.MustNewDecFromStr("0.5"),
			CollectionAddress:     suite.collectionAddrs[1].String(),
			CollectionAccountName: "",
			TaxSourceAddress:      suite.taxSourceAddrs[3].String(),
			TaxSourceAccountName:  "",
			StartTime:             mustParseRFC3339("0000-01-01T00:00:00Z"),
			EndTime:               mustParseRFC3339("9999-12-31T00:00:00Z"),
		},
	}

	for _, tc := range []struct {
		name           string
		taxes          []types.Tax
		accAsserts     []sdk.AccAddress
		balanceAsserts []sdk.Coins
		expectErr      bool
	}{
		{
			"basic taxes case",
			taxes1,
			[]sdk.AccAddress{
				suite.collectionAddrs[0],
				suite.collectionAddrs[1],
				suite.collectionAddrs[2],
				suite.collectionAddrs[3],
				suite.taxSourceAddrs[0],
				suite.taxSourceAddrs[1],
				suite.taxSourceAddrs[2],
			},
			[]sdk.Coins{
				mustParseCoinsNormalized("500000000denom1,500000000denom2,500000000denom3,500000000stake"),
				mustParseCoinsNormalized("500000000denom1,500000000denom2,500000000denom3,500000000stake"),
				mustParseCoinsNormalized("1000000000denom1,1000000000denom2,1000000000denom3,1000000000stake"),
				sdk.Coins{},
				sdk.Coins{},
				sdk.Coins{},
				mustParseCoinsNormalized("1000000000denom1,1000000000denom2,1000000000denom3,1000000000stake"),
			},
			false,
		},
		{
			"only expired tax case",
			[]types.Tax{taxes1[3]},
			[]sdk.AccAddress{
				suite.collectionAddrs[3],
				suite.taxSourceAddrs[2],
			},
			[]sdk.Coins{
				sdk.Coins{},
				mustParseCoinsNormalized("1000000000denom1,1000000000denom2,1000000000denom3,1000000000stake"),
			},
			false,
		},
		{
			"tax source has small balances case",
			taxes3,
			[]sdk.AccAddress{
				suite.collectionAddrs[0],
				suite.collectionAddrs[1],
				suite.taxSourceAddrs[3],
			},
			[]sdk.Coins{
				mustParseCoinsNormalized("1denom2,1denom3,500000000stake"),
				mustParseCoinsNormalized("1denom2,1denom3,500000000stake"),
				mustParseCoinsNormalized("1denom1,1denom3"),
			},
			false,
		},
	} {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			suite.keeper.SetParams(suite.ctx, types.Params{Taxes: tc.taxes})
			//err := simapp.FundAccount(suite.app.BankKeeper, suite.ctx, suite.taxSourceAddrs[3], mustParseCoinsNormalized("1denom1,1denom2,1denom3,1stake"))
			//suite.Require().NoError(err)

			err := suite.keeper.TaxCollection(suite.ctx)
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
