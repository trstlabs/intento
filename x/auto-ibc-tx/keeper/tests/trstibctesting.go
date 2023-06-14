package trstibctesting

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctesting "github.com/cosmos/ibc-go/v4/testing"
	simapp "github.com/cosmos/ibc-go/v4/testing/simapp"
	"github.com/cosmos/ibc-go/v4/testing/simapp/helpers"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"github.com/trstlabs/trst/app"
	// "github.com/trstlabs/trst/x/compute"
)

type TestChain struct {
	*ibctesting.TestChain
}

func SetupTestingApp() (ibctesting.TestingApp, map[string]json.RawMessage) {
	TrstApp := app.NewTrstApp(log.NewNopLogger(), dbm.NewMemDB(), nil, true, map[int64]bool{}, app.DefaultNodeHome, 5, true, simapp.EmptyAppOptions{} /* , compute.DefaultWasmConfig(), app.GetEnabledProposals() */)
	//TrstApp.AppKeepers.AutoIBCTXKeeper.SetParams(sdk.NewContext(), autoibctxtypes.DefaultParams())
	//encCdc := app.MakeEncodingConfig()
	return TrstApp, app.NewDefaultGenesisState(TrstApp.AppCodec())
}

// SendMsgsNoCheck is an alternative to ibctesting.TestChain.SendMsgs so that it doesn't check for errors. That should be handled by the caller
func (chain *TestChain) SendMsgsNoCheck(msgs ...sdk.Msg) (*sdk.Result, error) {
	// ensure the chain has the latest time
	chain.Coordinator.UpdateTimeForChain(chain.TestChain)

	_, r, err := SignAndDeliver(
		chain.TxConfig,
		chain.App.GetBaseApp(),
		chain.GetContext().BlockHeader(),
		msgs,
		chain.ChainID,
		[]uint64{chain.SenderAccount.GetAccountNumber()},
		[]uint64{chain.SenderAccount.GetSequence()},
		chain.SenderPrivKey,
	)
	if err != nil {
		return nil, err
	}

	// SignAndDeliver calls app.Commit()
	chain.NextBlock()

	// increment sequence for successful transaction execution
	err = chain.SenderAccount.SetSequence(chain.SenderAccount.GetSequence() + 1)
	if err != nil {
		return nil, err
	}

	chain.Coordinator.IncrementTime()

	return r, nil
}

// SignAndDeliver signs and delivers a transaction without asserting the results. This overrides the function
// from ibctesting
func SignAndDeliver(
	txCfg client.TxConfig, app *baseapp.BaseApp, header tmproto.Header, msgs []sdk.Msg,
	chainID string, accNums, accSeqs []uint64, priv ...cryptotypes.PrivKey,
) (sdk.GasInfo, *sdk.Result, error) {
	tx, _ := helpers.GenTx(
		txCfg,
		msgs,
		sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 2500)},
		helpers.DefaultGenTxGas,
		chainID,
		accNums,
		accSeqs,
		priv...,
	)

	// Simulate a sending a transaction and committing a block
	gInfo, res, err := app.Deliver(txCfg.TxEncoder(), tx)

	return gInfo, res, err
}

/*
// Move epochs to the future to avoid issues with minting
func (chain *TestChain) MoveEpochsToTheFuture() error {
	epochsKeeper := chain.GetTrstApp().EpochsKeeper
	ctx := chain.GetContext()
	for _, epoch := range epochsKeeper.AllEpochInfos(ctx) {
		epoch.StartTime = ctx.BlockTime().Add(time.Hour * 24 * 30)
		epochsKeeper.DeleteEpochInfo(chain.GetContext(), epoch.Identifier)
		err := epochsKeeper.AddEpochInfo(ctx, epoch)
		if err != nil {
			return err
		}
	}
	return nil
} */

// GetTrstApp returns the current chain's app as an TrstApp
func (chain *TestChain) GetTrstApp() *app.TrstApp {
	v, _ := chain.App.(*app.TrstApp)

	return v
}
