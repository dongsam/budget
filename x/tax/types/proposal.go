package types

import (
	"fmt"

	"github.com/gogo/protobuf/proto"

	"github.com/cosmos/cosmos-sdk/codec/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypePublicPlan string = "PublicPlan"
)

// Implements Proposal Interface
var _ gov.Content = &PublicPlanProposal{}

func init() {
	gov.RegisterProposalType(ProposalTypePublicPlan)
	gov.RegisterProposalTypeCodec(&PublicPlanProposal{}, "cosmos-sdk/PublicPlanProposal")
}

func NewPublicPlanProposal(title, description string, taxes []PlanI) (gov.Content, error) {
	taxesAny, err := PackTaxes(taxes)
	if err != nil {
		panic(err)
	}

	return &PublicPlanProposal{
		Title:       title,
		Description: description,
		Taxes:       taxesAny,
	}, nil
}

func (p *PublicPlanProposal) GetTitle() string { return p.Title }

func (p *PublicPlanProposal) GetDescription() string { return p.Description }

func (p *PublicPlanProposal) ProposalRoute() string { return RouterKey }

func (p *PublicPlanProposal) ProposalType() string { return ProposalTypePublicPlan }

func (p *PublicPlanProposal) ValidateBasic() error {
	for _, plan := range p.Taxes {
		_, ok := plan.GetCachedValue().(PlanI)
		if !ok {
			return fmt.Errorf("expected planI")
		}
		// TODO: PlanI needs ValidateBasic()?
		// if err := p.ValidateBasic(); err != nil {
		// 	return err
		// }
	}
	return gov.ValidateAbstract(p)
}

func (p PublicPlanProposal) String() string {
	return fmt.Sprintf(`Create FixedAmountPlan Proposal:
  Title:       %s
  Description: %s
  Taxes: 	   %s
`, p.Title, p.Description, p.Taxes)
}

// PackTaxes converts PlanIs to Any slice.
func PackTaxes(taxes []PlanI) ([]*types.Any, error) {
	taxesAny := make([]*types.Any, len(taxes))
	for i, plan := range taxes {
		msg, ok := plan.(proto.Message)
		if !ok {
			return nil, fmt.Errorf("cannot proto marshal %T", plan)
		}
		any, err := types.NewAnyWithValue(msg)
		if err != nil {
			return nil, err
		}
		taxesAny[i] = any
	}

	return taxesAny, nil
}

// UnpackTaxes converts Any slice to PlanIs.
func UnpackTaxes(taxesAny []*types.Any) ([]PlanI, error) {
	taxes := make([]PlanI, len(taxesAny))
	for i, any := range taxesAny {
		p, ok := any.GetCachedValue().(PlanI)
		if !ok {
			return nil, fmt.Errorf("expected planI")
		}
		taxes[i] = p
	}

	return taxes, nil
}

// UnpackPlan converts Any slice to PlanI.
func UnpackPlan(any *types.Any) (PlanI, error) {
	p, ok := any.GetCachedValue().(PlanI)
	if !ok {
		return nil, fmt.Errorf("expected planI")
	}

	return p, nil
}
