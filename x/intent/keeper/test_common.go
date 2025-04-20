package keeper

import (
	"encoding/binary"
	"path/filepath"
	"testing"
	"time"

	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
	storemetrics "cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/feegrant"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/std"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	ica "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	claimkeeper "github.com/trstlabs/intento/x/claim/keeper"
	claimtypes "github.com/trstlabs/intento/x/claim/types"
	interchainquerykeeper "github.com/trstlabs/intento/x/interchainquery/keeper"
	interchainquerytypes "github.com/trstlabs/intento/x/interchainquery/types"

	//icacontroller "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
	ibctransferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v8/modules/core"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"

	///intent "github.com/trstlabs/intento/x/intent"

	"github.com/stretchr/testify/require"

	"cosmossdk.io/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"

	"github.com/cosmos/cosmos-sdk/baseapp"

	"cosmossdk.io/store"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	crypto "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/cosmos/cosmos-sdk/x/auth"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/cosmos/ibc-go/modules/capability"

	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"

	"github.com/cosmos/cosmos-sdk/x/distribution"

	"cosmossdk.io/x/evidence"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
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

	"cosmossdk.io/x/upgrade"

	"github.com/trstlabs/intento/x/intent/types"
	intenttypes "github.com/trstlabs/intento/x/intent/types"
	// "github.com/trstlabs/intento/x/registration"
)

func setupTest(t *testing.T, additionalCoinsInWallets sdk.Coins) (sdk.Context, Keeper, sdk.AccAddress, crypto.PrivKey, sdk.AccAddress, crypto.PrivKey) {

	ctx, keepers, _ := CreateTestInput(t, false)
	accKeeper, keeper := keepers.AccountKeeper, keepers.IntentKeeper

	walletA, privKeyA := CreateFakeFundedAccount(ctx, accKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("denom", 200000)).Add(additionalCoinsInWallets...))
	walletB, privKeyB := CreateFakeFundedAccount(ctx, accKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("denom", 5000)).Add(additionalCoinsInWallets...))

	keeper.SetParams(ctx, intenttypes.Params{
		FlowFundsCommission: 2,
		BurnFeePerMsg:       1_000_000,
		FlowFlexFeeMul:      100,
		GasFeeCoins:         sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(1))),
		MaxFlowDuration:     time.Hour * 24 * 366 * 10,
		MinFlowDuration:     time.Second * 60,
		MinFlowInterval:     time.Second * 20,
		RelayerRewards:      []int64{10_000, 10_000, 10_000, 10_000},
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

	val, err := stakingtypes.NewValidator(sdk.ValAddress(valConsAddr).String(), pk, stakingtypes.Description{})

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
		[]govclient.ProposalHandler{paramsclient.ProposalHandler},
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
	MinCommissionRate: math.LegacyNewDec(0),
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
	tempDir := t.TempDir()

	keys := storetypes.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		minttypes.StoreKey, distrtypes.StoreKey, slashingtypes.StoreKey,
		govtypes.StoreKey, paramstypes.StoreKey, ibcexported.StoreKey, upgradetypes.StoreKey,
		evidencetypes.StoreKey, ibctransfertypes.StoreKey,
		capabilitytypes.StoreKey, feegrant.StoreKey, authzkeeper.StoreKey,
		intenttypes.StoreKey, icacontrollertypes.StoreKey,
	)

	db := dbm.NewMemDB()
	logger := log.NewTestLogger(t)
	ms := store.NewCommitMultiStore(db, logger, storemetrics.NewNoOpMetrics())
	for _, v := range keys {
		ms.MountStoreWithDB(v, storetypes.StoreTypeIAVL, db)
	}

	tkeys := storetypes.NewTransientStoreKeys(paramstypes.TStoreKey)
	for _, v := range tkeys {
		ms.MountStoreWithDB(v, storetypes.StoreTypeTransient, db)
	}

	memKeys := storetypes.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)
	for _, v := range memKeys {
		ms.MountStoreWithDB(v, storetypes.StoreTypeMemory, db)
	}

	require.NoError(t, ms.LoadLatestVersion())
	_, valConsPk0, _ := keyPubAddr()
	valCons := sdk.ConsAddress(valConsPk0.Address())
	val, _ := CreateValidator(valConsPk0, math.NewInt(100_000))

	ctx := sdk.NewContext(ms, tmproto.Header{
		Height:          1234567,
		Time:            time.Date(2020, time.April, 22, 12, 0, 0, 0, time.UTC),
		ChainID:         TestConfig.ChainID,
		ProposerAddress: valCons,
	}, isCheckTx, log.NewNopLogger())

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(sdk.Bech32MainPrefix, sdk.PrefixPublic)
	// config.SetBech32PrefixForValidator(sdk.Bech32MainPrefix, sdk.PrefixPublic)
	// config.SetBech32PrefixForConsensusNode(sdk.Bech32MainPrefix, sdk.PrefixPublic)

	encodingConfig := MakeEncodingConfig()
	paramsKeeper := paramskeeper.NewKeeper(
		encodingConfig.Codec,
		encodingConfig.Amino,
		keys[paramstypes.StoreKey],
		tkeys[paramstypes.TStoreKey],
	)

	//TrstApp := app.NewTrstApp(log.NewNopLogger(), dbm.NewMemDB(), nil, true, app.EmptyAppOptions{})

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
		intenttypes.ModuleName:         {authtypes.Minter, authtypes.Burner},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
	}

	accountKeeper := authkeeper.NewAccountKeeper(
		encodingConfig.Codec,
		runtime.NewKVStoreService(keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		maccPerms,
		authcodec.NewBech32Codec(sdk.Bech32MainPrefix),
		sdk.Bech32MainPrefix,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	require.NoError(t, accountKeeper.Params.Set(ctx, authtypes.DefaultParams()))
	blockedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		allowReceivingFunds := acc != distrtypes.ModuleName
		blockedAddrs[authtypes.NewModuleAddress(acc).String()] = allowReceivingFunds
	}
	bankKeeper := bankkeeper.NewBaseKeeper(
		encodingConfig.Codec,
		runtime.NewKVStoreService(keys[banktypes.StoreKey]),
		accountKeeper,
		blockedAddrs,
		authtypes.NewModuleAddress(banktypes.ModuleName).String(),
		logger,
	)

	bankKeeper.SetParams(ctx, banktypes.DefaultParams())

	stakingKeeper := stakingkeeper.NewKeeper(
		encodingConfig.Codec,
		runtime.NewKVStoreService(keys[stakingtypes.StoreKey]),
		accountKeeper,
		bankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	)
	err := stakingKeeper.SetParams(ctx, TestingStakeParams)
	require.NoError(t, err)
	val = stakingkeeper.TestingUpdateValidator(stakingKeeper, ctx, val, true)
	stakingKeeper.SetValidator(ctx, val)
	stakingKeeper.SetValidatorByConsAddr(ctx, val)

	stakingKeeper.Hooks().AfterValidatorCreated(ctx, sdk.ValAddress(val.GetOperator()))
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
		runtime.NewKVStoreService(keys[distrtypes.StoreKey]), accountKeeper, bankKeeper, stakingKeeper, authtypes.FeeCollectorName, authtypes.NewModuleAddress(govtypes.ModuleName).String())
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
	distKeeper.Params.Set(ctx, distrtypes.DefaultParams())
	distKeeper.FeePool.Set(ctx, distrtypes.InitialFeePool())
	stakingKeeper.SetHooks(stakingtypes.NewMultiStakingHooks(distKeeper.Hooks()))

	// set some funds ot pay out validators, based on code from:
	// https://github.com/cosmos/cosmos-sdk/blob/fea231556aee4d549d7551a6190389c4328194eb/x/distribution/keeper/keeper_test.go#L50-L57
	// distrAcc := distKeeper.GetDistributionAccount(ctx)

	testAccsSupply := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(2000000)))
	testIntentSupply := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100000000)))
	err = bankKeeper.MintCoins(ctx, faucetAccountName, testAccsSupply.Add(testIntentSupply[0]))
	require.NoError(t, err)
	err = bankKeeper.MintCoins(ctx, (distrtypes.ModuleName), testAccsSupply.Add(testIntentSupply[0]))
	require.NoError(t, err)

	// err = bankKeeper.SendCoinsFromModuleToAccount(ctx, faucetAccountName, distrAcc.GetAddress(), totalSupply)
	// require.NoError(t, err)
	// distrAcc := authtypes.NewEmptyModuleAccount(distrtypes.ModuleName, authtypes.Minter)

	// notBondedPool := authtypes.NewEmptyModuleAccount(stakingtypes.NotBondedPoolName, authtypes.Burner, authtypes.Staking)
	// bondPool := authtypes.NewEmptyModuleAccount(stakingtypes.BondedPoolName, authtypes.Burner, authtypes.Staking)
	// feeCollectorAcc := authtypes.NewEmptyModuleAccount(authtypes.FeeCollectorName)
	// intentAcc := authtypes.NewEmptyModuleAccount(intenttypes.ModuleName)

	// fmt.Printf("MOD %s \n", accountKeeper.GetModuleAccount(ctx, distrtypes.ModuleName))
	// fmt.Printf("MOD %s \n", accountKeeper.GetModuleAccount(ctx, stakingtypes.NotBondedPoolName))
	// fmt.Printf("MOD %s \n", accountKeeper.GetModuleAccount(ctx, stakingtypes.BondedPoolName))
	// fmt.Printf("MOD %s \n", accountKeeper.GetModuleAccount(ctx, authtypes.FeeCollectorName))
	// fmt.Printf("MOD %s \n", accountKeeper.GetModuleAccount(ctx, intenttypes.ModuleName))
	// accountKeeper.SetModuleAccount(ctx, distrAcc)
	// accountKeeper.SetModuleAccount(ctx, bondPool)
	// accountKeeper.SetModuleAccount(ctx, notBondedPool)
	// accountKeeper.SetModuleAccount(ctx, feeCollectorAcc)
	// accountKeeper.SetModuleAccount(ctx, intentAcc)

	err = bankKeeper.SendCoinsFromModuleToModule(ctx, faucetAccountName, stakingtypes.NotBondedPoolName, testAccsSupply)
	require.NoError(t, err)

	err = bankKeeper.SendCoinsFromModuleToModule(ctx, faucetAccountName, intenttypes.ModuleName, testIntentSupply)
	require.NoError(t, err)

	mintKeeper := mintkeeper.NewKeeper(encodingConfig.Codec, runtime.NewKVStoreService(keys[minttypes.StoreKey]), accountKeeper, bankKeeper, authtypes.FeeCollectorName, authtypes.NewModuleAddress(govtypes.ModuleName).String())
	mintKeeper.SetMinter(ctx, minttypes.DefaultInitialMinter())

	upgradeKeeper := upgradekeeper.NewKeeper(
		map[int64]bool{},
		runtime.NewKVStoreService(keys[upgradetypes.StoreKey]),
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
	scopedintentKeeper := capabilityKeeper.ScopeToModule(intenttypes.ModuleName)
	scopedTransferKeeper := capabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)

	ibchostSubSp, _ := paramsKeeper.GetSubspace(ibcexported.ModuleName)
	ibcKeeper := ibckeeper.NewKeeper(
		encodingConfig.Codec,
		keys[ibcexported.StoreKey],
		ibchostSubSp,
		stakingKeeper,
		upgradeKeeper,
		scopedIBCKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	ibcControllerSubSp, _ := paramsKeeper.GetSubspace(icacontrollertypes.SubModuleName)
	icacontrollerKeeper := icacontrollerkeeper.NewKeeper(encodingConfig.Codec, keys[icacontrollertypes.StoreKey], ibcControllerSubSp, ibcKeeper.ChannelKeeper, ibcKeeper.ChannelKeeper, ibcKeeper.PortKeeper, scopedIBCControllerKeeper, baseapp.NewMsgServiceRouter(), authtypes.NewModuleAddress(govtypes.ModuleName).String())

	// add keepers

	ibctransferSubSp, _ := paramsKeeper.GetSubspace(ibctransfertypes.ModuleName)
	ibctransferKeeper := ibctransferkeeper.NewKeeper(
		encodingConfig.Codec, keys[ibctransfertypes.StoreKey], ibctransferSubSp,
		ibcKeeper.ChannelKeeper, ibcKeeper.ChannelKeeper, ibcKeeper.PortKeeper,
		accountKeeper, bankKeeper, scopedTransferKeeper, authtypes.NewModuleAddress(govtypes.ModuleName).String())

	ibctransferKeeper.SetParams(ctx, ibctransfertypes.Params{
		SendEnabled: true,
	})

	govConfig := govtypes.DefaultConfig()

	claimKeeper := claimkeeper.NewKeeper(
		encodingConfig.Codec,
		runtime.NewKVStoreService(keys[claimtypes.StoreKey]),
		accountKeeper,
		bankKeeper,
		stakingKeeper,
		distKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	queryRouter := baseapp.NewGRPCQueryRouter()
	queryRouter.SetInterfaceRegistry(encodingConfig.InterfaceRegistry)
	msgServiceRouter := baseapp.NewMsgServiceRouter()
	msgServiceRouter.SetInterfaceRegistry(encodingConfig.InterfaceRegistry)

	govKeeper := govkeeper.NewKeeper(
		encodingConfig.Codec, runtime.NewKVStoreService(keys[govtypes.StoreKey]), accountKeeper, bankKeeper,
		stakingKeeper, distKeeper, msgServiceRouter, govConfig, authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	interchainQueryKeeper := interchainquerykeeper.NewKeeper(encodingConfig.Codec, runtime.NewKVStoreService(keys[interchainquerytypes.StoreKey]), ibcKeeper)

	intentKeeper := NewKeeper(
		encodingConfig.Codec,
		runtime.NewKVStoreService(keys[types.StoreKey]),
		icacontrollerKeeper,
		scopedintentKeeper,
		bankKeeper,
		distKeeper,
		*stakingKeeper,
		ibctransferKeeper,
		accountKeeper,
		interchainQueryKeeper,
		NewMultiIntentHooks(claimKeeper.Hooks()),
		msgServiceRouter,
		encodingConfig.InterfaceRegistry,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	intentKeeper.SetParams(ctx, intenttypes.DefaultParams())

	am := module.NewManager( // minimal module set that we use for message/ query tests
		bank.NewAppModule(encodingConfig.Codec, bankKeeper, accountKeeper, GetSubspace(banktypes.ModuleName, paramsKeeper)),
		staking.NewAppModule(encodingConfig.Codec, stakingKeeper, accountKeeper, bankKeeper, GetSubspace(stakingtypes.ModuleName, paramsKeeper)),
		distribution.NewAppModule(encodingConfig.Codec, distKeeper, accountKeeper, bankKeeper, stakingKeeper, GetSubspace(distrtypes.ModuleName, paramsKeeper)),
		gov.NewAppModule(encodingConfig.Codec, govKeeper, accountKeeper, bankKeeper, GetSubspace(govtypes.ModuleName, paramsKeeper)),
	)
	am.RegisterServices(module.NewConfigurator(encodingConfig.Codec, msgServiceRouter, queryRouter))
	// intenttypes.RegisterMsgServer(msgServiceRouter, NewMsgServerImpl(keeper))
	// intenttypes.RegisterQueryServer(queryRouter, intenttypes.QueryServer(keeper))

	keepers := TestKeepers{
		AccountKeeper:             accountKeeper,
		StakingKeeper:             *stakingKeeper,
		DistKeeper:                distKeeper,
		IntentKeeper:              intentKeeper,
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
	priv, _, addr := keyPubAddr()
	baseAcct := am.NewAccountWithAddress(ctx, addr)

	fundAccounts(ctx, am, bk, baseAcct, coins)
	return addr, priv
}

const faucetAccountName = "faucet"

func fundAccounts(ctx sdk.Context, am authkeeper.AccountKeeper, bk bankkeeper.Keeper, addr sdk.AccountI, coins sdk.Coins) {
	if err := bk.MintCoins(ctx, faucetAccountName, coins); err != nil {
		panic(err)
	}

	_ = bk.SendCoinsFromModuleToAccount(ctx, faucetAccountName, addr.GetAddress(), coins)

	am.NewAccount(ctx, addr)
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

// MakeEncodingConfig creates an EncodingConfig for testing
func MakeEncodingConfig() moduletestutil.TestEncodingConfig {
	encodingConfig := moduletestutil.MakeTestEncodingConfig(
		auth.AppModule{},
		bank.AppModule{},
		staking.AppModule{},
		mint.AppModule{},
		slashing.AppModule{},
		gov.AppModule{},
		crisis.AppModule{},
		ibc.AppModule{},
	)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	//ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	// add  types
	types.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	types.RegisterLegacyAminoCodec(encodingConfig.Amino)

	return encodingConfig
}

func GetSubspace(moduleName string, paramsKeeper paramskeeper.Keeper) paramstypes.Subspace {
	subspace, _ := paramsKeeper.GetSubspace(moduleName)
	return subspace
}
