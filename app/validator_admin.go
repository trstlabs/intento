package app

import (
	"fmt"
	"strings"

	math "cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

const (
	// ProposalTypeValidatorAdd defines the type for a ValidatorAddProposal
	ProposalTypeValidatorAdd = "ValidatorAdd"

	// ProposalTypeValidatorRemove defines the type for a ValidatorRemoveProposal
	ProposalTypeValidatorRemove = "ValidatorRemove"
)

func init() {
	govv1beta1.RegisterProposalType("ValidatorAdd")
	govv1beta1.RegisterProposalType("ValidatorRemove")
}

// ValidatorAddProposal defines a proposal to add a new validator
type ValidatorAddProposal struct {
	Title       string          `json:"title" yaml:"title"`
	Description string          `json:"description" yaml:"description"`
	Valoper     string          `json:"valoper" yaml:"valoper"`
	PubKey      *codectypes.Any `json:"pub_key" yaml:"pub_key"`
	Moniker     string          `json:"moniker" yaml:"moniker"`
}

// NewValidatorAddProposal creates a new ValidatorAddProposal
func NewValidatorAddProposal(title, description, valoper string, pubKey codectypes.Any, moniker string) *ValidatorAddProposal {
	return &ValidatorAddProposal{
		Title:       title,
		Description: description,
		Valoper:     valoper,
		PubKey:      &pubKey,
		Moniker:     moniker,
	}
}

// GetTitle returns the title of the proposal
func (p *ValidatorAddProposal) GetTitle() string { return p.Title }

// GetDescription returns the description of the proposal
func (p *ValidatorAddProposal) GetDescription() string { return p.Description }

// RouterKey returns the routing key of the module
func (p *ValidatorAddProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of the proposal
func (p *ValidatorAddProposal) ProposalType() string { return ProposalTypeValidatorAdd }

// ValidateBasic validates the proposal
func (p *ValidatorAddProposal) ValidateBasic() error {
	if p.Title == "" || p.Valoper == "" || p.PubKey == nil {
		return fmt.Errorf("missing required fields")
	}
	return nil
}

// String implements the Stringer interface
func (p ValidatorAddProposal) String() string {
	b := new(strings.Builder)
	b.WriteString(fmt.Sprintf(`Validator Add Proposal:
  Title:       %s
  Description: %s
  Valoper:     %s
  Moniker:     %s
`, p.Title, p.Description, p.Valoper, p.Moniker))
	return b.String()
}

// ValidatorRemoveProposal defines a proposal to remove an existing validator
type ValidatorRemoveProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Valoper     string `json:"valoper" yaml:"valoper"`
}

// NewValidatorRemoveProposal creates a new ValidatorRemoveProposal
func NewValidatorRemoveProposal(title, description, valoper string) *ValidatorRemoveProposal {
	return &ValidatorRemoveProposal{
		Title:       title,
		Description: description,
		Valoper:     valoper,
	}
}

// GetTitle returns the title of the proposal
func (p *ValidatorRemoveProposal) GetTitle() string { return p.Title }

// GetDescription returns the description of the proposal
func (p *ValidatorRemoveProposal) GetDescription() string { return p.Description }

// RouterKey returns the routing key of the module
func (p *ValidatorRemoveProposal) ProposalRoute() string { return RouterKey }

const RouterKey = "validatoradmin"

// ProposalType returns the type of the proposal
func (p *ValidatorRemoveProposal) ProposalType() string { return ProposalTypeValidatorRemove }

// ValidateBasic validates the proposal
func (p *ValidatorRemoveProposal) ValidateBasic() error {
	if p.Title == "" || p.Valoper == "" {
		return fmt.Errorf("missing required fields")
	}
	return nil
}

// String implements the Stringer interface
func (p ValidatorRemoveProposal) String() string {
	b := new(strings.Builder)
	b.WriteString(fmt.Sprintf(`Validator Remove Proposal:
  Title:       %s
  Description: %s
  Valoper:     %s
`, p.Title, p.Description, p.Valoper))
	return b.String()
}

// NewValidatorAdminProposalHandler creates a governance handler to manage validator admin proposals
func NewValidatorAdminProposalHandler(
	stakingKeeper *stakingkeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
) govv1beta1.Handler {
	return func(ctx sdk.Context, content govv1beta1.Content) error {
		switch p := content.(type) {
		case *ValidatorAddProposal:
			valAddr, err := sdk.ValAddressFromBech32(p.Valoper)
			if err != nil {
				return err
			}

			// Create validator message
			msg := stakingtypes.MsgCreateValidator{
				Description: stakingtypes.Description{
					Moniker: p.Moniker,
				},
				Commission: stakingtypes.CommissionRates{
					Rate:          math.LegacyMustNewDecFromStr("0.1"),
					MaxRate:       math.LegacyMustNewDecFromStr("0.1"),
					MaxChangeRate: math.LegacyMustNewDecFromStr("0.1"),
				},
				MinSelfDelegation: math.NewInt(1),
				DelegatorAddress:  sdk.AccAddress(valAddr).String(),
				ValidatorAddress:  p.Valoper,
				Pubkey:            p.PubKey,
				Value: sdk.Coin{
					Denom:  sdk.DefaultBondDenom,
					Amount: math.NewInt(1000000),
				},
			}
			if !bankKeeper.HasBalance(ctx, sdk.AccAddress(valAddr), msg.Value) {
				return fmt.Errorf("validator account lacks funds for self-bond: need %s", msg.Value.String())
			}

			msgServer := stakingkeeper.NewMsgServerImpl(stakingKeeper)
			_, err = msgServer.CreateValidator(ctx, &msg)
			return err

		case *ValidatorRemoveProposal:
			valAddr, err := sdk.ValAddressFromBech32(p.Valoper)
			if err != nil {
				return err
			}
			delAddr := sdk.AccAddress(valAddr)

			// Get self-delegation
			delegation, err := stakingKeeper.GetDelegation(ctx, delAddr, valAddr)
			if err != nil {
				return fmt.Errorf("validator self-delegation not found for %s: %w", p.Valoper, err)
			}

			// Get validator to calculate tokens
			validator, err := stakingKeeper.GetValidator(ctx, valAddr)
			if err != nil {
				return fmt.Errorf("validator not found %s: %w", p.Valoper, err)
			}

			// Calculate tokens from shares
			// We want to unbond ALL shares
			tokens := validator.TokensFromShares(delegation.Shares).TruncateInt()
			coin := sdk.NewCoin(sdk.DefaultBondDenom, tokens)

			msg := stakingtypes.MsgUndelegate{
				DelegatorAddress: delAddr.String(),
				ValidatorAddress: p.Valoper,
				Amount:           coin,
			}

			msgServer := stakingkeeper.NewMsgServerImpl(stakingKeeper)
			_, err = msgServer.Undelegate(ctx, &msg)
			return err
		default:
			return fmt.Errorf("unrecognized validator admin proposal type: %T", content)
		}
	}
}
