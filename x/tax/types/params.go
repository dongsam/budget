package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"gopkg.in/yaml.v2"

	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys
var (
	KeyTaxes = []byte("Taxes")
)

var _ paramstypes.ParamSet = (*Params)(nil)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns the default tax module parameters.
func DefaultParams() Params {
	return Params{
		Taxes: []Tax{},
	}
}

// ParamSetPairs implements paramstypes.ParamSet.
func (p *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyTaxes, &p.Taxes, validateTaxes),
	}
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Validate validates parameters.
func (p Params) Validate() error {
	for _, v := range []struct {
		value     interface{}
		validator func(interface{}) error
	}{
		{p.Taxes, validateTaxes},
	} {
		if err := v.validator(v.value); err != nil {
			return err
		}
	}
	return nil
}

func validateTaxes(i interface{}) error {
	fmt.Println("@@@@@@@ validateTaxes", i)
	taxes, ok := i.([]Tax)
	fmt.Println(taxes)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for _, tax := range taxes {
		taxSourceAcc, err := sdk.AccAddressFromBech32(tax.TaxSourceAddress)
		if err != nil {
			return fmt.Errorf("invalid TaxSourceAddress type: %T", tax)
		}
		collectionAcc, err := sdk.AccAddressFromBech32(tax.CollectionAddress)
		if err != nil {
			return fmt.Errorf("invalid CollectionAddress type: %T", tax)
		}
		fmt.Println(taxSourceAcc, collectionAcc)
	}
	// TODO: check added, fixed, deleted

	// TODO: unimplemented
	return nil
}
