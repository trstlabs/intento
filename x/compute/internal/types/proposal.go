package types

import (
	"encoding/base64"
	fmt "fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

type ProposalType string

const (
	ProposalTypeStoreCode           ProposalType = "StoreCode"
	ProposalTypeInstantiateContract ProposalType = "InstantiateContract"
	ProposalTypeExecuteContract     ProposalType = "ExecuteContract"
)

// DisableAllProposals contains no wasm gov types.
var DisableAllProposals []ProposalType

// EnableAllProposals contains all wasm gov types as keys.
var EnableAllProposals = []ProposalType{
	ProposalTypeStoreCode,
	ProposalTypeInstantiateContract,
	ProposalTypeExecuteContract,
}

// ConvertToProposals maps each key to a ProposalType and returns a typed list.
// If any string is not a valid type (in this file), then return an error
func ConvertToProposals(keys []string) ([]ProposalType, error) {
	valid := make(map[string]bool, len(EnableAllProposals))
	for _, key := range EnableAllProposals {
		valid[string(key)] = true
	}

	proposals := make([]ProposalType, len(keys))
	for i, key := range keys {
		if _, ok := valid[key]; !ok {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "'%s' is not a valid ProposalType", key)
		}
		proposals[i] = ProposalType(key)
	}
	return proposals, nil
}

func init() { // register new content types with the sdk
	govtypes.RegisterProposalType(string(ProposalTypeStoreCode))
	govtypes.RegisterProposalType(string(ProposalTypeInstantiateContract))
	govtypes.RegisterProposalType(string(ProposalTypeExecuteContract))
	govtypes.RegisterProposalTypeCodec(&StoreCodeProposal{}, "wasm/StoreCodeProposal")
	govtypes.RegisterProposalTypeCodec(&InstantiateContractProposal{}, "wasm/InstantiateContractProposal")
	govtypes.RegisterProposalTypeCodec(&ExecuteContractProposal{}, "wasm/ExecuteContractProposal")
	govtypes.RegisterProposalTypeCodec(&UpdateCodeProposal{}, "wasm/UpdateCodeProposal")
	govtypes.RegisterProposalTypeCodec(&UpdateContractProposal{}, "wasm/UpdateContractProposal")

}

// String implements the Stringer interface.
func (p StoreCodeProposal) String() string {
	return fmt.Sprintf(`Store Code Proposal:
  Title:       %s
  Description: %s
  Contract Title:       %s
  Contract Description: %s
  Creator:      %s
  WasmCode:    %X
  Contract Duration: %s
`, p.Title, p.Description, p.ContractTitle, p.ContractDescription, p.Creator, p.WASMByteCode, p.DefaultDuration)
}

// GetTitle returns the title of the proposal
func (p *StoreCodeProposal) GetTitle() string { return p.Title }

// GetDescription returns the human readable description of the proposal
func (p StoreCodeProposal) GetDescription() string { return p.Description }

// ProposalRoute returns the routing key of a parameter change proposal.
func (p StoreCodeProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type
func (p StoreCodeProposal) ProposalType() string { return string(ProposalTypeStoreCode) }

// ValidateBasic validates the proposal
func (p StoreCodeProposal) ValidateBasic() error {
	if err := validateProposalCommons(p.Title, p.Description); err != nil {
		return err
	}
	if _, err := sdk.AccAddressFromBech32(p.Creator); err != nil {
		return sdkerrors.Wrap(err, "run as")
	}

	if err := validateWasmCode(p.WASMByteCode); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "code bytes %s", err.Error())
	}

	return nil
}

// MarshalYAML pretty prints the wasm byte code
func (p StoreCodeProposal) MarshalYAML() (interface{}, error) {
	return struct {
		Title               string `yaml:"title"`
		Description         string `yaml:"description"`
		Creator             string `yaml:"creator"`
		WASMByteCode        string `yaml:"wasm_byte_code"`
		ContractTitle       string `yaml:"contract_title"`
		ContractDescription string `yaml:"contract_description"`
	}{
		Title:               p.Title,
		Description:         p.Description,
		Creator:             p.Creator,
		WASMByteCode:        base64.StdEncoding.EncodeToString(p.WASMByteCode),
		ContractTitle:       p.ContractTitle,
		ContractDescription: p.ContractDescription,
	}, nil
}

func validateProposalCommons(title, description string) error {
	if strings.TrimSpace(title) != title {
		return sdkerrors.Wrap(govtypes.ErrInvalidProposalContent, "proposal title must not start/end with white spaces")
	}
	if len(title) == 0 {
		return sdkerrors.Wrap(govtypes.ErrInvalidProposalContent, "proposal title cannot be blank")
	}
	if len(title) > govtypes.MaxTitleLength {
		return sdkerrors.Wrapf(govtypes.ErrInvalidProposalContent, "proposal title is longer than max length of %d", govtypes.MaxTitleLength)
	}
	if strings.TrimSpace(description) != description {
		return sdkerrors.Wrap(govtypes.ErrInvalidProposalContent, "proposal description must not start/end with white spaces")
	}
	if len(description) == 0 {
		return sdkerrors.Wrap(govtypes.ErrInvalidProposalContent, "proposal description cannot be blank")
	}
	if len(description) > govtypes.MaxDescriptionLength {
		return sdkerrors.Wrapf(govtypes.ErrInvalidProposalContent, "proposal description is longer than max length of %d", govtypes.MaxDescriptionLength)
	}
	if len(description) > govtypes.MaxDescriptionLength {
		return sdkerrors.Wrapf(govtypes.ErrInvalidProposalContent, "proposal description is longer than max length of %d", govtypes.MaxDescriptionLength)
	}
	return nil
}

// ProposalRoute returns the routing key of a parameter change proposal.
func (p InstantiateContractProposal) ProposalRoute() string { return RouterKey }

// GetTitle returns the title of the proposal
func (p *InstantiateContractProposal) GetTitle() string { return p.Title }

// GetDescription returns the human readable description of the proposal
func (p InstantiateContractProposal) GetDescription() string { return p.Description }

// ProposalType returns the type
func (p InstantiateContractProposal) ProposalType() string {
	return string(ProposalTypeInstantiateContract)
}

// ValidateBasic validates the proposal
func (p InstantiateContractProposal) ValidateBasic() error {
	if err := validateProposalCommons(p.Title, p.Description); err != nil {
		return err
	}
	/*if _, err := sdk.AccAddressFromBech32(p.RunAs); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "run as")
	}*/

	if p.CodeID == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "code id is required")
	}

	if err := validateContractId(p.ContractId); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p InstantiateContractProposal) String() string {
	return fmt.Sprintf(`Instantiate Code Proposal:
  Title:       %s
  Description: %s
  Code id:     %d
  Contract id:     %s
  Msg:         %q
  Funds:         %q
`, p.Title, p.Description, p.CodeID, p.ContractId, p.Msg, p.Funds)
}

// MarshalYAML pretty prints the init message
func (p InstantiateContractProposal) MarshalYAML() (interface{}, error) {
	return struct {
		Title       string `yaml:"title"`
		Description string `yaml:"description"`
		//RunAs       string    `yaml:"run_as"`
		CodeID     uint64    `yaml:"code_id"`
		ContractId string    `yaml:"contract_id"`
		Msg        string    `yaml:"msg"`
		Funds      sdk.Coins `yaml:"funds"`
	}{
		Title:       p.Title,
		Description: p.Description,
		//RunAs:       p.RunAs,
		CodeID:     p.CodeID,
		ContractId: p.ContractId,
		Msg:        string(p.Msg),
		Funds:      p.Funds,
	}, nil
}

// ProposalRoute returns the routing key of a parameter change proposal.
func (p ExecuteContractProposal) ProposalRoute() string { return RouterKey }

// GetTitle returns the title of the proposal
func (p *ExecuteContractProposal) GetTitle() string { return p.Title }

// GetDescription returns the human readable description of the proposal
func (p ExecuteContractProposal) GetDescription() string { return p.Description }

// ProposalType returns the type
func (p ExecuteContractProposal) ProposalType() string { return string(ProposalTypeExecuteContract) }

// ValidateBasic validates the proposal
func (p ExecuteContractProposal) ValidateBasic() error {
	if err := validateProposalCommons(p.Title, p.Description); err != nil {
		return err
	}
	if _, err := sdk.AccAddressFromBech32(p.Contract); err != nil {
		return sdkerrors.Wrap(err, "contract")
	}
	/*if _, err := sdk.AccAddressFromBech32(p.RunAs); err != nil {
		return sdkerrors.Wrap(err, "run as")
	}
	if !p.Funds.IsValid() {
		return sdkerrors.ErrInvalidCoins
	}*/

	return nil
}

// String implements the Stringer interface.
func (p ExecuteContractProposal) String() string {
	return fmt.Sprintf(`Execute Contract Proposal:
  Title:       %s
  Description: %s
  Contract:    %s

  Msg:         %q

`, p.Title, p.Description, p.Contract, p.Msg)
}

// MarshalYAML pretty prints the migrate message
func (p ExecuteContractProposal) MarshalYAML() (interface{}, error) {
	return struct {
		Title       string    `yaml:"title"`
		Description string    `yaml:"description"`
		Contract    string    `yaml:"contract"`
		Msg         string    `yaml:"msg"`
		Funds       sdk.Coins `yaml:"funds"`
		/*RunAs       string    `yaml:"run_as"`
		Funds   sdk.Coins `yaml:"funds"`*/
	}{
		Title:       p.Title,
		Description: p.Description,
		Contract:    p.Contract,
		Msg:         string(p.Msg),
		Funds:       p.Funds,
		/*RunAs:       p.RunAs,*/
	}, nil
}

// ProposalRoute returns the routing key of a parameter change proposal.
func (p UpdateCodeProposal) ProposalRoute() string { return RouterKey }

// GetTitle returns the title of the proposal
func (p *UpdateCodeProposal) GetTitle() string { return p.Title }

// GetDescription returns the human readable description of the proposal
func (p UpdateCodeProposal) GetDescription() string { return p.Description }

// ProposalType returns the type
func (p UpdateCodeProposal) ProposalType() string { return string(ProposalTypeExecuteContract) }

// ValidateBasic validates the proposal
func (p UpdateCodeProposal) ValidateBasic() error {
	if err := validateProposalCommons(p.Title, p.Description); err != nil {
		return err
	}

	return nil
}

/*
// String implements the Stringer interface.
func (p UpdateCodeProposal) String() string {
	return fmt.Sprintf(`Update Code Proposal:
  Title:       %s
  Description: %s
  Contract:    %s
  DefaultDuration: %s
  DefaultInterval:         %s

`, p.Title, p.Description, p.CodeId, p.DefaultDuration, p.DefaultInterval)
}
*/
/*
// MarshalYAML pretty prints the migrate message
func (p UpdateCodeProposal) MarshalYAML() (interface{}, error) {
	return struct {
		Title       string    `yaml:"title"`
		Description string    `yaml:"description"`
		CodeId    string    `yaml:"code_id"`
		Msg         string    `yaml:"msg"`
		Funds       sdk.Coins `yaml:"funds"`
		/*RunAs       string    `yaml:"run_as"`
		Funds   sdk.Coins `yaml:"funds"`
	}{
		Title:       p.Title,
		Description: p.Description,
		Contract:    p.Contract,
		Msg:         string(p.Msg),
		Funds:       p.Funds,
		/*RunAs:       p.RunAs,
	}, nil
}
*/
// ProposalRoute returns the routing key of a parameter change proposal.
func (p UpdateContractProposal) ProposalRoute() string { return RouterKey }

// GetTitle returns the title of the proposal
func (p *UpdateContractProposal) GetTitle() string { return p.Title }

// GetDescription returns the human readable description of the proposal
func (p UpdateContractProposal) GetDescription() string { return p.Description }

// ProposalType returns the type
func (p UpdateContractProposal) ProposalType() string { return string(ProposalTypeExecuteContract) }

// ValidateBasic validates the proposal
func (p UpdateContractProposal) ValidateBasic() error {
	if err := validateProposalCommons(p.Title, p.Description); err != nil {
		return err
	}
	if _, err := sdk.AccAddressFromBech32(p.Contract); err != nil {
		return sdkerrors.Wrap(err, "contract")
	}

	return nil
}

/*
// String implements the Stringer interface.
func (p UpdateContractProposal) String() string {
	return fmt.Sprintf(`Update Contract Proposal:
  Title:       %s
  Description: %s
  Contract:    %s
StartTime: %s
EndTime: %s
Interval: %s
  Msg:         %q

`, p.Title, p.Description, p.Contract, p.StartTime, p.EndTime, p.Interval)
}*/

/*
// MarshalYAML pretty prints the migrate message
func (p UpdateContractProposal) MarshalYAML() (interface{}, error) {
	return struct {
		Title       string    `yaml:"title"`
		Description string    `yaml:"description"`
		Contract    string    `yaml:"contract"`
		Msg         string    `yaml:"msg"`
		Funds       sdk.Coins `yaml:"funds"`
		/*RunAs       string    `yaml:"run_as"`
		Funds   sdk.Coins `yaml:"funds"`
	}{
		Title:       p.Title,
		Description: p.Description,
		Contract:    p.Contract,
		Msg:         string(p.Msg),
		Funds:       p.Funds,
		/*RunAs:       p.RunAs,
	}, nil
}
*/
