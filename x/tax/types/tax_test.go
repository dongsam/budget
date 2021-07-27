package types_test


import (
	"fmt"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/types/address"
	jsoniter "github.com/json-iterator/go"
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/tax/x/tax/types"
)

func TestGetPoolInformation(t *testing.T) {
	commonTerminationAcc := sdk.AccAddress("terminationAddr")
	//commonStartTime := time.Now().UTC()
	//commonEndTime := commonStartTime.AddDate(1, 0, 0)
	//commonCoinWeights := sdk.NewDecCoins(
	//	sdk.DecCoin{Denom: "testFarmStakingCoinDenom", Amount: sdk.MustNewDecFromStr("1.0")},
	//)

	a := sdk.AccAddress(address.Module(types.ModuleName, []byte("StakingReserveAcc")))
	fmt.Println(a.String())
	fmt.Println(a)
	fmt.Println(address.Module(types.ModuleName, []byte("StakingReserveAcc")))
	fmt.Println(string(address.Module(types.ModuleName, []byte("StakingReserveAcc"))))
	fmt.Println(sdk.AccAddress(crypto.AddressHash((address.Module(types.ModuleName, []byte("StakingReserveAcc"))))))
	fmt.Println([]byte(string(address.Module(types.ModuleName, []byte("StakingReserveAcc")))))
	fmt.Println(sdk.AccAddress([]byte(string(address.Module(types.ModuleName, []byte("StakingReserveAcc"))))))

	tax := types.Tax{
		Name:                  "testTax",
		Rate:                  sdk.MustNewDecFromStr("0.1"),
		CollectionAddress:     commonTerminationAcc.String(),
		CollectionAccountName: "",
		TaxSourceAddress:      commonTerminationAcc.String(),
		TaxSourceAccountName:  "",
		StartTime:             time.Now(),
		EndTime:               time.Now(),
	}
	fmt.Println(types.Tax{})
	fmt.Println(tax.String())
	bz, _ := tax.Marshal()
	j := sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&bz))
	fmt.Println(string(j))

	jsoni := jsoniter.ConfigCompatibleWithStandardLibrary
	b, _ := jsoni.Marshal(tax)
	fmt.Println(string(b))
}
