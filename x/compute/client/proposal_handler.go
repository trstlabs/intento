package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/trstlabs/trst/x/compute/client/cli"
	"github.com/trstlabs/trst/x/compute/client/rest"
)

// ProposalHandlers define the wasm cli proposal types and rest handler.
var ProposalHandlers = []govclient.ProposalHandler{
	govclient.NewProposalHandler(cli.ProposalStoreCodeCmd, rest.StoreCodeProposalHandler),
}
