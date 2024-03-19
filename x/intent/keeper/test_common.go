package keeper

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"
	"time"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/std"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ica "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	appparams "github.com/trstlabs/intento/app/params"
	claimkeeper "github.com/trstlabs/intento/x/claim/keeper"
	claimtypes "github.com/trstlabs/intento/x/claim/types"

	//icacontroller "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	ibctransferkeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"

	///intent "github.com/trstlabs/intento/x/intent"

	"github.com/stretchr/testify/require"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	crypto "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/module"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"

	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/cosmos/cosmos-sdk/x/capability"

	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"

	"github.com/cosmos/cosmos-sdk/x/distribution"

	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/cosmos/cosmos-sdk/x/mint"

	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	mintkeeper "github.com/trstlabs/intento/x/mint/keeper"
	minttypes "github.com/trstlabs/intento/x/mint/types"

	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"

	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"

	intenttypes "github.com/trstlabs/intento/x/intent/types"
	// "github.com/trstlabs/intento/x/registration"
)

func setupTest(t *testing.T, additionalCoinsInWallets sdk.Coins) (sdk.Context, Keeper, sdk.AccAddress, crypto.PrivKey, sdk.AccAddress, crypto.PrivKey) {

	ctx, keepers, _ := CreateTestInput(t, false)
	accKeeper, keeper := keepers.AccountKeeper, keepers.IntentKeeper

	walletA, privKeyA := CreateFakeFundedAccount(ctx, accKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("denom", 200000)).Add(additionalCoinsInWallets...))
	walletB, privKeyB := CreateFakeFundedAccount(ctx, accKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("denom", 5000)).Add(additionalCoinsInWallets...))

	keeper.SetParams(ctx, intenttypes.Params{
		ActionFundsCommission:      2,
		ActionConstantFee:          1_000_000,
		ActionFlexFeeMul:           100,
		RecurringActionConstantFee: 1_000_000,
		MaxActionDuration:          time.Hour * 24 * 366 * 10,
		MinActionDuration:          time.Second * 60,
		MinActionInterval:          time.Second * 20,
		RelayerRewards:             []int64{10_000, 10_000, 10_000, 10_000},
	})
	return ctx, keeper, walletA, privKeyA, walletB, privKeyB
}

const (
	flagLRUCacheSize  = "lru_size"
	flagQueryGasLimit = "query_gas_limit"
	ibcContract       = "ibc.wasm"
)

const contractPath = "testdata"

var TestContractPaths = map[string]string{

	ibcContract: filepath.Join(".", contractPath, ibcContract),
}

func CreateValidator(pk crypto.PubKey, stake sdkmath.Int) (stakingtypes.Validator, error) {
	valConsAddr := sdk.GetConsAddress(pk)
	val, err := stakingtypes.NewValidator(sdk.ValAddress(valConsAddr), pk, stakingtypes.Description{})

	return val, err
}

type MockIBCTransferKeeper struct {
	GetPortFn func(ctx sdk.Context) string
}

func (m MockIBCTransferKeeper) GetPort(ctx sdk.Context) string {
	if m.GetPortFn == nil {
		panic("not expected to be called")
	}
	return m.GetPortFn(ctx)
}

var ModuleBasics = module.NewBasicManager(
	auth.AppModuleBasic{},
	bank.AppModuleBasic{},
	capability.AppModuleBasic{},
	staking.AppModuleBasic{},
	mint.AppModuleBasic{},
	distribution.AppModuleBasic{},
	gov.NewAppModuleBasic(
		[]govclient.ProposalHandler{paramsclient.ProposalHandler, upgradeclient.LegacyProposalHandler, upgradeclient.LegacyCancelProposalHandler},
	),
	params.AppModuleBasic{},
	crisis.AppModuleBasic{},
	slashing.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	evidence.AppModuleBasic{},
	// transfer.AppModuleBasic{},
	// registration.AppModuleBasic{},
	//intent.AppModuleBasic{},
	ica.AppModuleBasic{},
	ibc.AppModuleBasic{},
)

/*
	func MakeTestCodec() codec.Codec {
		return MakeEncodingConfig().Codec
	}

func MakeEncodingConfig() simappparams.EncodingConfig {

		interfaceRegistry := types.NewInterfaceRegistry()
		marshaler := codec.NewProtoCodec(interfaceRegistry)
		txCfg := authtx.NewTxConfig(marshaler, authtx.DefaultSignModes)

		std.RegisterInterfaces(interfaceRegistry)

		ModuleBasics.RegisterInterfaces(interfaceRegistry)

		intenttypes.RegisterInterfaces(interfaceRegistry)

		return simappparams.EncodingConfig{
			InterfaceRegistry: interfaceRegistry,
			Codec:             MakeTestCodec(),
			TxConfig:          txCfg,
		}
	}
*/
var TestingStakeParams = stakingtypes.Params{
	UnbondingTime:     100,
	MaxValidators:     10,
	MaxEntries:        10,
	HistoricalEntries: 10,
	BondDenom:         sdk.DefaultBondDenom,
	MinCommissionRate: sdk.NewDec(0),
}

type TestKeepers struct {
	AccountKeeper             authkeeper.AccountKeeper
	StakingKeeper             stakingkeeper.Keeper
	IntentKeeper              Keeper
	DistKeeper                distrkeeper.Keeper
	GovKeeper                 govkeeper.Keeper
	BankKeeper                bankkeeper.Keeper
	MintKeeper                mintkeeper.Keeper
	ParamsKeeper              paramskeeper.Keeper
	IbcKeeper                 ibckeeper.Keeper
	ICAControllerKeeper       *icacontrollerkeeper.Keeper
	scopedIBCControllerKeeper capabilitykeeper.ScopedKeeper
}

var TestConfig = TestConfigType{
	ChainID: "test-chain",
}

type TestConfigType struct {
	ChainID string
}

// encoders can be nil to accept the defaults, or set it to override some of the message handlers (like default)
func CreateTestInput(t *testing.T, isCheckTx bool) (sdk.Context, TestKeepers, codec.Codec) {
	tempDir, err := os.MkdirTemp("", "test")
	require.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(tempDir) })

	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		minttypes.StoreKey, distrtypes.StoreKey, slashingtypes.StoreKey,
		govtypes.StoreKey, paramstypes.StoreKey, ibcexported.StoreKey, upgradetypes.StoreKey,
		evidencetypes.StoreKey, ibctransfertypes.StoreKey,
		capabilitytypes.StoreKey, feegrant.StoreKey, authzkeeper.StoreKey,
		intenttypes.StoreKey, icacontrollertypes.StoreKey,
	)

	db := dbm.NewMemDB()

	ms := store.NewCommitMultiStore(db)
	for _, v := range keys {
		ms.MountStoreWithDB(v, storetypes.StoreTypeIAVL, db)
	}

	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)
	for _, v := range tkeys {
		ms.MountStoreWithDB(v, storetypes.StoreTypeTransient, db)
	}

	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)
	for _, v := range memKeys {
		ms.MountStoreWithDB(v, storetypes.StoreTypeMemory, db)
	}

	require.NoError(t, ms.LoadLatestVersion())
	_, valConsPk0, _ := keyPubAddr()
	valCons := sdk.ConsAddress(valConsPk0.Address())
	val, _ := CreateValidator(valConsPk0, sdk.NewInt(100_000))

	ctx := sdk.NewContext(ms, tmproto.Header{
		Height:          1234567,
		Time:            time.Date(2020, time.April, 22, 12, 0, 0, 0, time.UTC),
		ChainID:         TestConfig.ChainID,
		ProposerAddress: valCons,
	}, isCheckTx, log.NewNopLogger())

	encodingConfig := MakeEncodingConfig()
	paramsKeeper := paramskeeper.NewKeeper(
		encodingConfig.Codec,
		encodingConfig.Amino,
		keys[paramstypes.StoreKey],
		tkeys[paramstypes.TStoreKey],
	)

	//IntoApp := app.NewIntoApp(log.NewNopLogger(), dbm.NewMemDB(), nil, true, app.EmptyAppOptions{})

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibcexported.ModuleName)
	paramsKeeper.Subspace(intenttypes.ModuleName)
	paramsKeeper.Subspace(icacontrollertypes.SubModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)

	// this is also used to initialize module accounts (so nil is meaningful here)
	maccPerms := map[string][]string{
		faucetAccountName:              {authtypes.Burner, authtypes.Minter},
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          {authtypes.Minter},
		minttypes.ModuleName:           {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		intenttypes.ModuleName:         {authtypes.Minter},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
	}

	authKeeper := authkeeper.NewAccountKeeper(
		encodingConfig.Codec,
		keys[authtypes.StoreKey], // target store
		authtypes.ProtoBaseAccount,
		maccPerms,
		sdk.Bech32MainPrefix, authtypes.NewModuleAddress(govtypes.ModuleName).String())
	blockedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		allowReceivingFunds := acc != distrtypes.ModuleName
		blockedAddrs[authtypes.NewModuleAddress(acc).String()] = allowReceivingFunds
	}

	bankKeeper := bankkeeper.NewBaseKeeper(
		encodingConfig.Codec,
		keys[banktypes.StoreKey],
		authKeeper,
		blockedAddrs,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	bankParams := banktypes.DefaultParams()
	bankParams.DefaultSendEnabled = true
	bankKeeper.SetParams(ctx, banktypes.DefaultParams())

	stakingKeeper := stakingkeeper.NewKeeper(
		encodingConfig.Codec,
		keys[stakingtypes.StoreKey],
		authKeeper,
		bankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	err = stakingKeeper.SetParams(ctx, TestingStakeParams)
	require.NoError(t, err)
	val = stakingkeeper.TestingUpdateValidator(stakingKeeper, ctx, val, true)
	stakingKeeper.SetValidator(ctx, val)
	stakingKeeper.SetValidatorByConsAddr(ctx, val)

	stakingKeeper.Hooks().AfterValidatorCreated(ctx, val.GetOperator())
	//val, _ = val.AddTokensFromDel(sdk.TokensFromConsensusPower(1, sdk.DefaultPowerReduction))
	// mintSubsp, _ := paramsKeeper.GetSubspace(minttypes.ModuleName)

	// mintKeeper := mintkeeper.NewKeeper(encodingConfig.Codec,
	//	keyBank,
	//	mintSubsp,
	//	stakingKeeper,
	//	authKeeper,
	//	bankKeeper,
	//	authtypes.FeeCollectorName,
	//	)
	//
	// bankkeeper.SetSupply(ctx, banktypes.NewSupply(sdk.NewCoins((sdk.NewInt64Coin("stake", 1)))))

	distKeeper := distrkeeper.NewKeeper(
		encodingConfig.Codec,
		keys[distrtypes.StoreKey], authKeeper, bankKeeper, stakingKeeper, authtypes.FeeCollectorName, authtypes.NewModuleAddress(govtypes.ModuleName).String())
	// // set some baseline - this seems to be needed
	// distKeeper.SetValidatorHistoricalRewards(ctx, val.GetOperator(), 2, distrtypes.ValidatorHistoricalRewards{
	// 	CumulativeRewardRatio: sdk.DecCoins{},
	// 	ReferenceCount:        2,
	// })
	// distKeeper.SetValidatorCurrentRewards(ctx, val.GetOperator(), distrtypes.ValidatorCurrentRewards{
	// 	Rewards: sdk.DecCoins{},
	// 	Period:  3,
	// })
	// set genesis items required for distribution
	distKeeper.SetParams(ctx, distrtypes.DefaultParams())
	distKeeper.SetFeePool(ctx, distrtypes.InitialFeePool())
	stakingKeeper.SetHooks(stakingtypes.NewMultiStakingHooks(distKeeper.Hooks()))

	// set some funds ot pay out validators, based on code from:
	// https://github.com/cosmos/cosmos-sdk/blob/fea231556aee4d549d7551a6190389c4328194eb/x/distribution/keeper/keeper_test.go#L50-L57
	// distrAcc := distKeeper.GetDistributionAccount(ctx)
	distrAcc := authtypes.NewEmptyModuleAccount(distrtypes.ModuleName)

	testAccsSupply := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(2000000)))
	testIntentSupply := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100000000)))
	err = bankKeeper.MintCoins(ctx, faucetAccountName, testAccsSupply.Add(testIntentSupply[0]))
	require.NoError(t, err)
	err = bankKeeper.MintCoins(ctx, (distrtypes.ModuleName), testAccsSupply.Add(testIntentSupply[0]))
	require.NoError(t, err)

	// err = bankKeeper.SendCoinsFromModuleToAccount(ctx, faucetAccountName, distrAcc.GetAddress(), totalSupply)
	// require.NoError(t, err)

	notBondedPool := authtypes.NewEmptyModuleAccount(stakingtypes.NotBondedPoolName, authtypes.Burner, authtypes.Staking)
	bondPool := authtypes.NewEmptyModuleAccount(stakingtypes.BondedPoolName, authtypes.Burner, authtypes.Staking)
	feeCollectorAcc := authtypes.NewEmptyModuleAccount(authtypes.FeeCollectorName)
	autoIbcTxAcc := authtypes.NewEmptyModuleAccount(intenttypes.ModuleName)

	authKeeper.SetModuleAccount(ctx, autoIbcTxAcc)
	authKeeper.SetModuleAccount(ctx, distrAcc)
	authKeeper.SetModuleAccount(ctx, bondPool)
	authKeeper.SetModuleAccount(ctx, notBondedPool)
	authKeeper.SetModuleAccount(ctx, feeCollectorAcc)

	err = bankKeeper.SendCoinsFromModuleToModule(ctx, faucetAccountName, stakingtypes.NotBondedPoolName, testAccsSupply)
	require.NoError(t, err)

	err = bankKeeper.SendCoinsFromModuleToModule(ctx, faucetAccountName, intenttypes.ModuleName, testIntentSupply)
	require.NoError(t, err)

	/* router := baseapp.NewRouter()
	bh := bank.NewHandler(bankKeeper)
	router.AddRoute(sdk.NewRoute(banktypes.RouterKey, bh))
	sh := staking.NewHandler(stakingKeeper)
	router.AddRoute(sdk.NewRoute(stakingtypes.RouterKey, sh))
	dh := distribution.NewHandler(distKeeper)
	router.AddRoute(sdk.NewRoute(distrtypes.RouterKey, dh))

	govRouter := govtypes.NewRouter().
		AddRoute(govtypes.RouterKey, govtypes.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(paramsKeeper)).
		AddRoute(distrtypes.RouterKey, distribution.NewCommunityPoolSpendProposalHandler(distKeeper))
	// AddRoute(wasmTypes.RouterKey, NewWasmProposalHandler(keeper, wasmTypes.EnableAllProposals))
	*/

	// bank := bankKeeper.
	// bk := bank.Keeper(bankKeeper)

	mintSubsp, _ := paramsKeeper.GetSubspace(minttypes.ModuleName)
	mintKeeper := mintkeeper.NewKeeper(encodingConfig.Codec, keys[minttypes.StoreKey], mintSubsp, authKeeper, bankKeeper, authtypes.FeeCollectorName)
	mintKeeper.SetMinter(ctx, minttypes.DefaultInitialMinter())

	// keeper := NewKeeper(cdc, keyContract, accountKeeper, &bk, &govKeeper, &distKeeper, &mintKeeper, &stakingKeeper, router, tempDir, wasmConfig, supportedFeatures, encoders, queriers)
	//// add wasm handler so we can loop-back (contracts calling contracts)
	// router.AddRoute(wasmTypes.RouterKey, TestHandler(keeper))

	// govKeeper.SetProposalID(ctx, govtypes.DefaultStartingProposalID)
	// govKeeper.SetDepositParams(ctx, govtypes.DefaultDepositParams())
	// govKeeper.SetVotingParams(ctx, govtypes.DefaultVotingParams())
	// govKeeper.SetTallyParams(ctx, govtypes.DefaultTallyParams())
	// gh := gov.NewHandler(govKeeper)
	// router.AddRoute(sdk.NewRoute(govtypes.RouterKey, gh))

	upgradeKeeper := upgradekeeper.NewKeeper(
		map[int64]bool{},
		keys[upgradetypes.StoreKey],
		encodingConfig.Codec,
		tempDir,
		&baseapp.BaseApp{},
		authtypes.NewModuleAddress(govtypes.ModuleName).String())

	capabilityKeeper := capabilitykeeper.NewKeeper(
		encodingConfig.Codec,
		keys[capabilitytypes.StoreKey],
		memKeys[capabilitytypes.MemStoreKey],
	)

	scopedIBCKeeper := capabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	scopedIBCControllerKeeper := capabilityKeeper.ScopeToModule(icacontrollertypes.SubModuleName)
	scopedAutoIBCTXKeeper := capabilityKeeper.ScopeToModule(intenttypes.ModuleName)
	scopedTransferKeeper := capabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)

	ibchostSubSp, _ := paramsKeeper.GetSubspace(ibcexported.ModuleName)
	ibcKeeper := ibckeeper.NewKeeper(
		encodingConfig.Codec,
		keys[ibcexported.StoreKey],
		ibchostSubSp,
		stakingKeeper,
		upgradeKeeper,
		scopedIBCKeeper,
	)
	ibcControllerSubSp, _ := paramsKeeper.GetSubspace(icacontrollertypes.SubModuleName)
	icacontrollerKeeper := icacontrollerkeeper.NewKeeper(encodingConfig.Codec, keys[icacontrollertypes.StoreKey], ibcControllerSubSp, ibcKeeper.ChannelKeeper, ibcKeeper.ChannelKeeper, &ibcKeeper.PortKeeper, scopedIBCControllerKeeper, baseapp.NewMsgServiceRouter())

	// add keepers

	accountKeeper := authkeeper.NewAccountKeeper(encodingConfig.Codec, keys[authtypes.StoreKey], authtypes.ProtoBaseAccount, maccPerms, sdk.Bech32MainPrefix, authtypes.NewModuleAddress(govtypes.ModuleName).String())

	ibctransferSubSp, _ := paramsKeeper.GetSubspace(ibctransfertypes.ModuleName)
	ibctransferKeeper := ibctransferkeeper.NewKeeper(
		encodingConfig.Codec, keys[ibctransfertypes.StoreKey], ibctransferSubSp,
		ibcKeeper.ChannelKeeper, ibcKeeper.ChannelKeeper, &ibcKeeper.PortKeeper,
		accountKeeper, bankKeeper, scopedTransferKeeper)

	govConfig := govtypes.DefaultConfig()

	claimKeeper := claimkeeper.NewKeeper(
		encodingConfig.Codec,
		keys[claimtypes.StoreKey],
		accountKeeper,
		bankKeeper,
		stakingKeeper,
		distKeeper,
	)

	queryRouter := baseapp.NewGRPCQueryRouter()
	queryRouter.SetInterfaceRegistry(encodingConfig.InterfaceRegistry)
	msgServiceRouter := baseapp.NewMsgServiceRouter()
	msgServiceRouter.SetInterfaceRegistry(encodingConfig.InterfaceRegistry)
	//AutoIBCTXKeeper := AutoIBCTXKeeper.NewKeeper(appCodec, ak.keys[icaauthtypes.StoreKey], *ak.ICAControllerKeeper, ak.ScopedAutoIBCTXKeeper, ak.BankKeeper, *ak.DistrKeeper, *ak.StakingKeeper, *ak.AccountKeeper, ak.GetSubspace(icaauthtypes.ModuleName))
	autoIbcTxSubsp, _ := paramsKeeper.GetSubspace(intenttypes.ModuleName)
	govKeeper := govkeeper.NewKeeper(
		encodingConfig.Codec, keys[govtypes.StoreKey], accountKeeper, bankKeeper,
		stakingKeeper, msgServiceRouter, govConfig, authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	autoIbcTxKeeper := NewKeeper(
		encodingConfig.Codec,
		keys[intenttypes.StoreKey],
		icacontrollerKeeper,
		scopedAutoIBCTXKeeper,
		bankKeeper,
		distKeeper,
		*stakingKeeper,
		ibctransferKeeper,
		accountKeeper,
		autoIbcTxSubsp,
		NewMultiIntentHooks(claimKeeper.Hooks()),
		msgServiceRouter,
	)
	autoIbcTxKeeper.SetParams(ctx, intenttypes.DefaultParams())

	am := module.NewManager( // minimal module set that we use for message/ query tests
		bank.NewAppModule(encodingConfig.Codec, bankKeeper, authKeeper, GetSubspace(banktypes.ModuleName, paramsKeeper)),
		staking.NewAppModule(encodingConfig.Codec, stakingKeeper, authKeeper, bankKeeper, GetSubspace(stakingtypes.ModuleName, paramsKeeper)),
		distribution.NewAppModule(encodingConfig.Codec, distKeeper, authKeeper, bankKeeper, stakingKeeper, GetSubspace(distrtypes.ModuleName, paramsKeeper)),
		gov.NewAppModule(encodingConfig.Codec, govKeeper, authKeeper, bankKeeper, GetSubspace(govtypes.ModuleName, paramsKeeper)),
	)
	am.RegisterServices(module.NewConfigurator(encodingConfig.Codec, msgServiceRouter, queryRouter))
	// intenttypes.RegisterMsgServer(msgServiceRouter, NewMsgServerImpl(keeper))
	// intenttypes.RegisterQueryServer(queryRouter, intenttypes.QueryServer(keeper))

	keepers := TestKeepers{
		AccountKeeper:             authKeeper,
		StakingKeeper:             *stakingKeeper,
		DistKeeper:                distKeeper,
		IntentKeeper:              autoIbcTxKeeper,
		GovKeeper:                 *govKeeper,
		BankKeeper:                bankKeeper,
		MintKeeper:                mintKeeper,
		ParamsKeeper:              paramsKeeper,
		IbcKeeper:                 *ibcKeeper,
		scopedIBCControllerKeeper: scopedIBCControllerKeeper,
		ICAControllerKeeper:       &icacontrollerKeeper,
	}

	return ctx, keepers, encodingConfig.Codec
}

func CreateFakeFundedAccount(ctx sdk.Context, am authkeeper.AccountKeeper, bk bankkeeper.Keeper, coins sdk.Coins) (sdk.AccAddress, crypto.PrivKey) {
	priv, pub, addr := keyPubAddr()
	baseAcct := authtypes.NewBaseAccountWithAddress(addr)
	_ = baseAcct.SetPubKey(pub)
	am.SetAccount(ctx, baseAcct)

	fundAccounts(ctx, am, bk, addr, coins)
	return addr, priv
}

const faucetAccountName = "faucet"

func fundAccounts(ctx sdk.Context, am authkeeper.AccountKeeper, bk bankkeeper.Keeper, addr sdk.AccAddress, coins sdk.Coins) {
	baseAcct := am.GetAccount(ctx, addr)
	if err := bk.MintCoins(ctx, faucetAccountName, coins); err != nil {
		panic(err)
	}

	_ = bk.SendCoinsFromModuleToAccount(ctx, faucetAccountName, addr, coins)

	am.SetAccount(ctx, baseAcct)
}

var keyCounter uint64 = 0

// we need to make this deterministic (same every test run), as encoded address size and thus gas cost,
// depends on the actual bytes (due to ugly CanonicalAddress encoding)
func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	keyCounter++
	seed := make([]byte, 8)
	binary.BigEndian.PutUint64(seed, keyCounter)

	key := secp256k1.GenPrivKeyFromSecret(seed)
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}

type protoTxProvider interface {
	GetProtoTx() *sdktx.Tx
}

func txBuilderToProtoTx(txBuilder client.TxBuilder) (*sdktx.Tx, error) { // nolint
	protoProvider, ok := txBuilder.(protoTxProvider)
	if !ok {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidType, "expected proto tx builder, got %T", txBuilder)
	}

	return protoProvider.GetProtoTx(), nil
}

// MakeEncodingConfig creates an EncodingConfig for testing
func MakeEncodingConfig() appparams.EncodingConfig {
	encodingConfig := appparams.MakeEncodingConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}

func GetSubspace(moduleName string, paramsKeeper paramskeeper.Keeper) paramstypes.Subspace {
	subspace, _ := paramsKeeper.GetSubspace(moduleName)
	return subspace
}
