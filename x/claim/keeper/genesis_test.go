package keeper_test

import (
	"github.com/trstlabs/intento/x/claim"
	"github.com/trstlabs/intento/x/claim/types"
)

func (s *KeeperTestSuite) TestExportGenesis() {
	app, ctx := s.app, s.ctx

	claim.InitGenesis(ctx, app.ClaimKeeper, *types.DefaultGenesis())
	// app.ClaimKeeper.SetParams(ctx, types.DefaultParams())
	exported := claim.ExportGenesis(ctx, app.ClaimKeeper)
	gen := types.DefaultGenesis()
	gen.Params.AirdropStartTime = ctx.BlockTime()
	s.Require().Equal(gen.Params.DurationOfDecay, exported.Params.DurationOfDecay)
}
