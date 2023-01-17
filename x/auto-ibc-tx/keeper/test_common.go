package keeper

import (
	"encoding/binary"

	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ica "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts"

	//icacontroller "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/controller"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/controller/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v3/modules/core"
	ibchost "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v3/modules/core/keeper"

	///auto-ibc-tx "github.com/trstlabs/trst/x/auto-ibc-tx"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	crypto "github.com/cosmos/cosmos-sdk/crypto/types"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/module"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
	sdksigning "github.com/cosmos/cosmos-sdk/types/tx/signing"

	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/cosmos/cosmos-sdk/x/capability"

	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"

	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	distrclient "github.com/cosmos/cosmos-sdk/x/distribution/client"
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
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"

	mintkeeper "github.com/trstlabs/trst/x/mint/keeper"
	minttypes "github.com/trstlabs/trst/x/mint/types"

	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"

	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"

	autoibctxtypes "github.com/trstlabs/trst/x/auto-ibc-tx/types"
	"github.com/trstlabs/trst/x/registration"
)

func setupTest(t *testing.T, additionalCoinsInWallets sdk.Coins) (sdk.Context, Keeper, sdk.AccAddress, crypto.PrivKey, sdk.AccAddress, crypto.PrivKey) {

	ctx, keepers := CreateTestInput(t, false)
	accKeeper, keeper := keepers.AccountKeeper, keepers.AutoIbcTxKeeper

	walletA, privKeyA := CreateFakeFundedAccount(ctx, accKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("denom", 200000)).Add(additionalCoinsInWallets...))
	walletB, privKeyB := CreateFakeFundedAccount(ctx, accKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("denom", 5000)).Add(additionalCoinsInWallets...))

	keeper.SetParams(ctx, autoibctxtypes.Params{
		AutoTxFundsCommission:      2,
		AutoTxConstantFee:          1_000_000,                 // 1trst
		AutoTxFlexFeeMul:           100,                       // 100/100 = 1 = gasUsed
		RecurringAutoTxConstantFee: 1_000_000,                 // 1trst
		MaxAutoTxDuration:          time.Hour * 24 * 366 * 10, // a little over 10 years
		MinAutoTxDuration:          time.Second * 40,
		MinAutoTxInterval:          time.Second * 20,
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

func CreateValidator(pk crypto.PubKey, stake sdk.Int) (stakingtypes.Validator, error) {
	valConsAddr := sdk.GetConsAddress(pk)
	val, err := stakingtypes.NewValidator(sdk.ValAddress(valConsAddr), pk, stakingtypes.Description{})
	val.Tokens = stake
	val.DelegatorShares = sdk.NewDecFromInt(val.Tokens)
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
		paramsclient.ProposalHandler, distrclient.ProposalHandler, upgradeclient.ProposalHandler,
	),
	params.AppModuleBasic{},
	crisis.AppModuleBasic{},
	slashing.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	evidence.AppModuleBasic{},
	// transfer.AppModuleBasic{},
	registration.AppModuleBasic{},
	//auto-ibc-tx.AppModuleBasic{},
	ica.AppModuleBasic{},
	ibc.AppModuleBasic{},
)

func MakeTestCodec() codec.Codec {
	return MakeEncodingConfig().Marshaler
}

func MakeEncodingConfig() simappparams.EncodingConfig {
	amino := codec.NewLegacyAmino()
	interfaceRegistry := types.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	txCfg := authtx.NewTxConfig(marshaler, authtx.DefaultSignModes)

	std.RegisterInterfaces(interfaceRegistry)
	std.RegisterLegacyAminoCodec(amino)

	ModuleBasics.RegisterInterfaces(interfaceRegistry)
	ModuleBasics.RegisterLegacyAminoCodec(amino)
	autoibctxtypes.RegisterInterfaces(interfaceRegistry)
	autoibctxtypes.RegisterCodec(amino)
	return simappparams.EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TxConfig:          txCfg,
		Amino:             amino,
	}
}

var TestingStakeParams = stakingtypes.Params{
	UnbondingTime:     100,
	MaxValidators:     10,
	MaxEntries:        10,
	HistoricalEntries: 10,
	BondDenom:         sdk.DefaultBondDenom,
}

type TestKeepers struct {
	AccountKeeper             authkeeper.AccountKeeper
	StakingKeeper             stakingkeeper.Keeper
	AutoIbcTxKeeper           Keeper
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
func CreateTestInput(t *testing.T, isCheckTx bool) (sdk.Context, TestKeepers) {
	tempDir, err := os.MkdirTemp("", "wasm")
	require.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(tempDir) })

	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		minttypes.StoreKey, distrtypes.StoreKey, slashingtypes.StoreKey,
		govtypes.StoreKey, paramstypes.StoreKey, ibchost.StoreKey, upgradetypes.StoreKey,
		evidencetypes.StoreKey, ibctransfertypes.StoreKey,
		capabilitytypes.StoreKey, feegrant.StoreKey, authzkeeper.StoreKey,
		autoibctxtypes.StoreKey, icacontrollertypes.StoreKey,
	)

	db := dbm.NewMemDB()

	ms := store.NewCommitMultiStore(db)
	for _, v := range keys {
		ms.MountStoreWithDB(v, sdk.StoreTypeIAVL, db)
	}

	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)
	for _, v := range tkeys {
		ms.MountStoreWithDB(v, sdk.StoreTypeTransient, db)
	}

	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)
	for _, v := range memKeys {
		ms.MountStoreWithDB(v, sdk.StoreTypeMemory, db)
	}

	require.NoError(t, ms.LoadLatestVersion())
	_, valConsPk0, _ := keyPubAddr()
	valCons := sdk.ConsAddress(valConsPk0.Address())
	val, _ := CreateValidator(valConsPk0, sdk.NewInt(100))
	//val, err := stakingkeeper.Keeper.SetValidatorByConsAddr(ctx,valConsPk0//, math.NewInt(100))

	ctx := sdk.NewContext(ms, tmproto.Header{
		Height:          1234567,
		Time:            time.Date(2020, time.April, 22, 12, 0, 0, 0, time.UTC),
		ChainID:         TestConfig.ChainID,
		ProposerAddress: valCons,
	}, isCheckTx, log.NewNopLogger())
	encodingConfig := MakeEncodingConfig()
	paramsKeeper := paramskeeper.NewKeeper(
		encodingConfig.Marshaler,
		encodingConfig.Amino,
		keys[paramstypes.StoreKey],
		tkeys[paramstypes.TStoreKey],
	)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibchost.ModuleName)
	paramsKeeper.Subspace(autoibctxtypes.ModuleName)
	paramsKeeper.Subspace(icacontrollertypes.SubModuleName)

	// this is also used to initialize module accounts (so nil is meaningful here)
	maccPerms := map[string][]string{
		faucetAccountName:              {authtypes.Burner, authtypes.Minter},
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          nil,
		minttypes.ModuleName:           {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		autoibctxtypes.ModuleName:      {authtypes.Minter},
	}
	authSubsp, _ := paramsKeeper.GetSubspace(authtypes.ModuleName)
	authKeeper := authkeeper.NewAccountKeeper(
		encodingConfig.Marshaler,
		keys[authtypes.StoreKey], // target store
		authSubsp,
		authtypes.ProtoBaseAccount, // prototype
		maccPerms,
	)
	blockedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		allowReceivingFunds := acc != distrtypes.ModuleName
		blockedAddrs[authtypes.NewModuleAddress(acc).String()] = allowReceivingFunds
	}

	bankSubsp, _ := paramsKeeper.GetSubspace(banktypes.ModuleName)
	bankKeeper := bankkeeper.NewBaseKeeper(
		encodingConfig.Marshaler,
		keys[banktypes.StoreKey],
		authKeeper,
		bankSubsp,
		blockedAddrs,
	)

	// bankParams = bankParams.SetSendEnabledParam(sdk.DefaultBondDenom, true)
	bankKeeper.SetParams(ctx, banktypes.DefaultParams())

	stakingSubsp, _ := paramsKeeper.GetSubspace(stakingtypes.ModuleName)
	stakingKeeper := stakingkeeper.NewKeeper(
		encodingConfig.Marshaler,
		keys[stakingtypes.StoreKey],
		authKeeper,
		bankKeeper,
		stakingSubsp,
	)
	stakingKeeper.SetValidator(ctx, val)
	stakingKeeper.SetValidatorByConsAddr(ctx, val)
	stakingKeeper.SetParams(ctx, TestingStakeParams)

	// mintSubsp, _ := paramsKeeper.GetSubspace(minttypes.ModuleName)

	// mintKeeper := mintkeeper.NewKeeper(encodingConfig.Marshaler,
	//	keyBank,
	//	mintSubsp,
	//	stakingKeeper,
	//	authKeeper,
	//	bankKeeper,
	//	authtypes.FeeCollectorName,
	//	)
	//
	// bankkeeper.SetSupply(ctx, banktypes.NewSupply(sdk.NewCoins((sdk.NewInt64Coin("stake", 1)))))

	distSubsp, _ := paramsKeeper.GetSubspace(distrtypes.ModuleName)
	distKeeper := distrkeeper.NewKeeper(
		encodingConfig.Marshaler,
		keys[distrtypes.StoreKey],
		distSubsp,
		authKeeper,
		bankKeeper,
		stakingKeeper,
		authtypes.FeeCollectorName,
		nil,
	)

	// set genesis items required for distribution
	distKeeper.SetParams(ctx, distrtypes.DefaultParams())
	distKeeper.SetFeePool(ctx, distrtypes.InitialFeePool())
	stakingKeeper.SetHooks(stakingtypes.NewMultiStakingHooks(distKeeper.Hooks()))

	// set some funds ot pay out validators, based on code from:
	// https://github.com/cosmos/cosmos-sdk/blob/fea231556aee4d549d7551a6190389c4328194eb/x/distribution/keeper/keeper_test.go#L50-L57
	// distrAcc := distKeeper.GetDistributionAccount(ctx)
	distrAcc := authtypes.NewEmptyModuleAccount(distrtypes.ModuleName)

	testAccsSupply := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(2000000)))
	testAutoIbcTxSupply := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100000000)))
	err = bankKeeper.MintCoins(ctx, faucetAccountName, testAccsSupply.Add(testAutoIbcTxSupply[0]))
	require.NoError(t, err)

	// err = bankKeeper.SendCoinsFromModuleToAccount(ctx, faucetAccountName, distrAcc.GetAddress(), totalSupply)
	// require.NoError(t, err)

	notBondedPool := authtypes.NewEmptyModuleAccount(stakingtypes.NotBondedPoolName, authtypes.Burner, authtypes.Staking)
	bondPool := authtypes.NewEmptyModuleAccount(stakingtypes.BondedPoolName, authtypes.Burner, authtypes.Staking)
	feeCollectorAcc := authtypes.NewEmptyModuleAccount(authtypes.FeeCollectorName)
	autoIbcTxAcc := authtypes.NewEmptyModuleAccount(autoibctxtypes.ModuleName)

	authKeeper.SetModuleAccount(ctx, autoIbcTxAcc)
	authKeeper.SetModuleAccount(ctx, distrAcc)
	authKeeper.SetModuleAccount(ctx, bondPool)
	authKeeper.SetModuleAccount(ctx, notBondedPool)
	authKeeper.SetModuleAccount(ctx, feeCollectorAcc)

	err = bankKeeper.SendCoinsFromModuleToModule(ctx, faucetAccountName, stakingtypes.NotBondedPoolName, testAccsSupply)
	require.NoError(t, err)

	err = bankKeeper.SendCoinsFromModuleToModule(ctx, faucetAccountName, autoibctxtypes.ModuleName, testAutoIbcTxSupply)
	require.NoError(t, err)

	router := baseapp.NewRouter()
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

	govKeeper := govkeeper.NewKeeper(
		encodingConfig.Marshaler, keys[govtypes.StoreKey], paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govtypes.ParamKeyTable()), authKeeper, bankKeeper, stakingKeeper, govRouter,
	)

	// bank := bankKeeper.
	// bk := bank.Keeper(bankKeeper)

	mintSubsp, _ := paramsKeeper.GetSubspace(minttypes.ModuleName)
	mintKeeper := mintkeeper.NewKeeper(encodingConfig.Marshaler, keys[minttypes.StoreKey], mintSubsp, authKeeper, bankKeeper, authtypes.FeeCollectorName)
	mintKeeper.SetMinter(ctx, minttypes.DefaultInitialMinter())

	// keeper := NewKeeper(cdc, keyContract, accountKeeper, &bk, &govKeeper, &distKeeper, &mintKeeper, &stakingKeeper, router, tempDir, wasmConfig, supportedFeatures, encoders, queriers)
	//// add wasm handler so we can loop-back (contracts calling contracts)
	// router.AddRoute(wasmTypes.RouterKey, TestHandler(keeper))

	govKeeper.SetProposalID(ctx, govtypes.DefaultStartingProposalID)
	govKeeper.SetDepositParams(ctx, govtypes.DefaultDepositParams())
	govKeeper.SetVotingParams(ctx, govtypes.DefaultVotingParams())
	govKeeper.SetTallyParams(ctx, govtypes.DefaultTallyParams())
	gh := gov.NewHandler(govKeeper)
	router.AddRoute(sdk.NewRoute(govtypes.RouterKey, gh))

	upgradeKeeper := upgradekeeper.NewKeeper(
		map[int64]bool{},
		keys[upgradetypes.StoreKey],
		encodingConfig.Marshaler,
		tempDir,
		nil,
	)

	capabilityKeeper := capabilitykeeper.NewKeeper(
		encodingConfig.Marshaler,
		keys[capabilitytypes.StoreKey],
		memKeys[capabilitytypes.MemStoreKey],
	)

	scopedIBCKeeper := capabilityKeeper.ScopeToModule(ibchost.ModuleName)
	scopedIBCControllerKeeper := capabilityKeeper.ScopeToModule(icacontrollertypes.SubModuleName)
	scopedICAAuthKeeper := capabilityKeeper.ScopeToModule(autoibctxtypes.ModuleName)

	ibchostSubSp, _ := paramsKeeper.GetSubspace(ibchost.ModuleName)
	ibcKeeper := ibckeeper.NewKeeper(
		encodingConfig.Marshaler,
		keys[ibchost.StoreKey],
		ibchostSubSp,
		stakingKeeper,
		upgradeKeeper,
		scopedIBCKeeper,
	)
	ibcControllerSubSp, _ := paramsKeeper.GetSubspace(icacontrollertypes.SubModuleName)
	icacontrollerKeeper := icacontrollerkeeper.NewKeeper(encodingConfig.Marshaler, keys[icacontrollertypes.StoreKey], ibcControllerSubSp, ibcKeeper.ChannelKeeper, ibcKeeper.ChannelKeeper, &ibcKeeper.PortKeeper, scopedIBCControllerKeeper, baseapp.NewMsgServiceRouter())

	// add keepers
	accSubsp, _ := paramsKeeper.GetSubspace(authtypes.ModuleName)
	accountKeeper := authkeeper.NewAccountKeeper(
		encodingConfig.Marshaler, keys[authtypes.StoreKey], accSubsp, authtypes.ProtoBaseAccount, maccPerms)

	queryRouter := baseapp.NewGRPCQueryRouter()
	queryRouter.SetInterfaceRegistry(encodingConfig.InterfaceRegistry)
	msgRouter := baseapp.NewMsgServiceRouter()
	msgRouter.SetInterfaceRegistry(encodingConfig.InterfaceRegistry)
	//icaAuthKeeper := icaauthkeeper.NewKeeper(appCodec, ak.keys[icaauthtypes.StoreKey], *ak.ICAControllerKeeper, ak.ScopedICAAuthKeeper, ak.BankKeeper, *ak.DistrKeeper, *ak.StakingKeeper, *ak.AccountKeeper, ak.GetSubspace(icaauthtypes.ModuleName))
	autoIbcTxSubsp, _ := paramsKeeper.GetSubspace(autoibctxtypes.ModuleName)

	keeper := NewKeeper(
		encodingConfig.Marshaler,
		keys[autoibctxtypes.StoreKey],
		icacontrollerKeeper,
		scopedICAAuthKeeper,
		bankKeeper,
		distKeeper,
		stakingKeeper,
		accountKeeper,
		autoIbcTxSubsp,
	)
	//keeper.SetParams(ctx, autoibctxtypes.DefaultParams())

	am := module.NewManager( // minimal module set that we use for message/ query tests
		bank.NewAppModule(encodingConfig.Marshaler, bankKeeper, authKeeper),
		staking.NewAppModule(encodingConfig.Marshaler, stakingKeeper, authKeeper, bankKeeper),
		distribution.NewAppModule(encodingConfig.Marshaler, distKeeper, authKeeper, bankKeeper, stakingKeeper),
		gov.NewAppModule(encodingConfig.Marshaler, govKeeper, authKeeper, bankKeeper),
	)
	am.RegisterServices(module.NewConfigurator(encodingConfig.Marshaler, msgRouter, queryRouter))
	autoibctxtypes.RegisterMsgServer(msgRouter, NewMsgServerImpl(keeper))
	autoibctxtypes.RegisterQueryServer(queryRouter, autoibctxtypes.QueryServer(keeper))

	keepers := TestKeepers{
		AccountKeeper:             authKeeper,
		StakingKeeper:             stakingKeeper,
		DistKeeper:                distKeeper,
		AutoIbcTxKeeper:           keeper,
		GovKeeper:                 govKeeper,
		BankKeeper:                bankKeeper,
		MintKeeper:                mintKeeper,
		ParamsKeeper:              paramsKeeper,
		IbcKeeper:                 *ibcKeeper,
		scopedIBCControllerKeeper: scopedIBCControllerKeeper,
		ICAControllerKeeper:       &icacontrollerKeeper,
	}

	return ctx, keepers
}

/*
// TestHandler returns a wasm handler for tests (to avoid circular imports)
func TestHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *autoibctxtypes,:
			return handleInstantiate(ctx, k, msg)

		case *wasmTypes.MsgExecuteContract:
			return handleExecute(ctx, k, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized wasm message type: %T", msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func PrepareInitSignedTx(t *testing.T, keeper Keeper, ctx sdk.Context, creator sdk.AccAddress, privKey crypto.PrivKey, encMsg []byte, codeID uint64, funds sdk.Coins) sdk.Context {
	creatorAcc, err := ante.GetSignerAcc(ctx, keeper.accountKeeper, creator)
	require.NoError(t, err)

	initMsg := wasmTypes.MsgInstantiateContract{
		Sender:     creator.String(),
		CodeID:     codeID,
		ContractId: "demo contract 1",
		Msg:        encMsg,
		Funds:      funds,
	}
	tx := NewTestTx(&initMsg, creatorAcc, privKey)

	txBytes, err := tx.Marshal()
	require.NoError(t, err)

	return ctx.WithTxBytes(txBytes)
}

func PrepareExecSignedTx(t *testing.T, keeper Keeper, ctx sdk.Context, sender sdk.AccAddress, privKey crypto.PrivKey, encMsg []byte, contract sdk.AccAddress, funds sdk.Coins) sdk.Context {
	creatorAcc, err := ante.GetSignerAcc(ctx, keeper.accountKeeper, sender)
	require.NoError(t, err)

	executeMsg := wasmTypes.MsgExecuteContract{
		Sender:   sender.String(),
		Contract: contract.String(),
		Msg:      encMsg,
		Funds:    funds,
	}
	tx := NewTestTx(&executeMsg, creatorAcc, privKey)

	txBytes, err := tx.Marshal()
	require.NoError(t, err)

	return ctx.WithTxBytes(txBytes)
} */

func NewTestTx(msg sdk.Msg, creatorAcc authtypes.AccountI, privKey crypto.PrivKey) *sdktx.Tx {
	return NewTestTxMultiple([]sdk.Msg{msg}, []authtypes.AccountI{creatorAcc}, []crypto.PrivKey{privKey})
}

func NewTestTxMultiple(msgs []sdk.Msg, creatorAccs []authtypes.AccountI, privKeys []crypto.PrivKey) *sdktx.Tx {
	if len(msgs) != len(creatorAccs) || len(msgs) != len(privKeys) {
		panic("length of `msgs` `creatorAccs` and `privKeys` must be the same")
	}

	// There's no need to pass values to `NewTxConfig` because they get ignored by `NewTxBuilder` anyways,
	// and we just need the builder, which can not be created any other way, apparently.
	txConfig := authtx.NewTxConfig(nil, authtx.DefaultSignModes)
	signModeHandler := txConfig.SignModeHandler()
	builder := txConfig.NewTxBuilder()
	builder.SetFeeAmount(nil)
	builder.SetGasLimit(0)
	builder.SetTimeoutHeight(0)

	err := builder.SetMsgs(msgs...)
	if err != nil {
		panic(err)
	}

	// This code is based on `cosmos-sdk/client/tx/tx.go::Sign()`
	var sigs []sdksigning.SignatureV2
	for _, creatorAcc := range creatorAccs {
		sig := sdksigning.SignatureV2{
			PubKey: creatorAcc.GetPubKey(),
			Data: &sdksigning.SingleSignatureData{
				SignMode:  sdksigning.SignMode_SIGN_MODE_DIRECT,
				Signature: nil,
			},
			Sequence: creatorAcc.GetSequence(),
		}
		sigs = append(sigs, sig)
	}
	err = builder.SetSignatures(sigs...)
	if err != nil {
		panic(err)
	}

	sigs = []sdksigning.SignatureV2{}
	for i, creatorAcc := range creatorAccs {
		privKey := privKeys[i]
		signerData := authsigning.SignerData{
			ChainID:       TestConfig.ChainID,
			AccountNumber: creatorAcc.GetAccountNumber(),
			Sequence:      creatorAcc.GetSequence(),
		}
		bytesToSign, err := signModeHandler.GetSignBytes(sdksigning.SignMode_SIGN_MODE_DIRECT, signerData, builder.GetTx())

		signBytes, err := privKey.Sign(bytesToSign)
		if err != nil {
			panic(err)
		}
		sig := sdksigning.SignatureV2{
			PubKey: creatorAcc.GetPubKey(),
			Data: &sdksigning.SingleSignatureData{
				SignMode:  sdksigning.SignMode_SIGN_MODE_DIRECT,
				Signature: signBytes,
			},
			Sequence: creatorAcc.GetSequence(),
		}
		sigs = append(sigs, sig)
	}

	err = builder.SetSignatures(sigs...)
	if err != nil {
		panic(err)
	}

	tx, ok := builder.(protoTxProvider)
	if !ok {
		panic("failed to unwrap tx builder to protobuf tx")
	}
	return tx.GetProtoTx()
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
	GetProtoTx() *tx.Tx
}

func txBuilderToProtoTx(txBuilder client.TxBuilder) (*tx.Tx, error) { // nolint
	protoProvider, ok := txBuilder.(protoTxProvider)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "expected proto tx builder, got %T", txBuilder)
	}

	return protoProvider.GetProtoTx(), nil
}
