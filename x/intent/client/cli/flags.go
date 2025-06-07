package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	// The connection end identifier on the controller chain
	flagConnectionID = "connection-id"
	// The controller chain channel version
	flagHostConnectionID      = "host-connection-id"
	flagLabel                 = "label"
	flagDuration              = "duration"
	flagInterval              = "interval"
	flagStartAt               = "start-at"
	flagFeeFunds              = "fee-funds"
	flagEndTime               = "end-at"
	flagHostedAccount         = "hosted-account"
	flagHostedAccountFeeLimit = "hosted-account-fee-limit"

	//Execution conditions
	flagUpdatingDisabled       = "updating-disabled"
	flagSaveResponses          = "save-responses"
	flagFallbackToOwnerBalance = "fallback-to-owner-balance"
	flagStopOnSuccess          = "stop-on-success"
	flagStopOnFailure          = "stop-on-failure"
	flagStopOnTimeout          = "stop-on-timeout"
	flagStopOnSuccessOf        = "stop-on-success-of"
	flagStopOnFailureOf        = "stop-on-failure-of"
	flagSkipOnSuccessOf        = "skip-on-success-of"
	flagSkipOnFailureOf        = "skip-on-failure-of"

	flagFeeCoinsSupported = "fee-coins-supported"
	flagNewAdmin          = "new-admin"
	flagConditions        = "conditions"
	flagICQConfig         = "icq-config"
)

// common flagsets to add to various functions
var (
	fsFlow = flag.NewFlagSet("", flag.ContinueOnError)
	fsIBC  = flag.NewFlagSet("", flag.ContinueOnError)
	// fsVersion = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	fsFlow.String(flagLabel, "", "A custom label for the flow e.g. AutoTransfer, UpdateContractParams, optional")
	fsFlow.String(flagInterval, "", "A custom interval for the flow e.g. 2h, 6000s, 72h3m0.5s, optional")
	fsFlow.String(flagStartAt, "0", "A custom start time in UNIX time, optional")
	fsFlow.String(flagFeeFunds, "", "Coin to attach to the flow, optional")
	fsFlow.String(flagConnectionID, "", "Connection ID from this chain to the host chain, optional")
	fsFlow.String(flagHostConnectionID, "", "Connection ID from host chain to this chain, optional")
	fsFlow.String(flagHostedAccount, "", "A hosted account to execute actions on a host, optional")
	fsFlow.String(flagConditions, "", "intent conditions in JSON format, optional")
	fsFlow.String(flagHostedAccountFeeLimit, "", "Coin to sent to limit the hosted fees, optional")
	fsFlow.String(flagICQConfig, "", "A config to query keyvalue store on a host, optional")

	fsFlow.Bool(flagUpdatingDisabled, false, "disable future updates to the configuration'")
	fsFlow.Bool(flagSaveResponses, true, "save message and query responses to tx history (Cosmos SDK v0.46+ chains only), true on default'")
	fsFlow.Bool(flagStopOnSuccess, false, "stop execution after success'")
	fsFlow.Bool(flagStopOnFailure, false, "stop execution after failure'")
	fsFlow.Bool(flagStopOnTimeout, false, " If true, allows the flow to continue execution after an ibc channel times out'")
	fsFlow.Bool(flagFallbackToOwnerBalance, false, "fallback to owner balance'")

	fsIBC.String(flagConnectionID, "", "Connection ID from this chain to the host chain, optional")
	fsIBC.String(flagHostConnectionID, "", "Connection ID from host chain to this chain, optional")

	// fsFlow.StringArray(flagSkipOnSuccessOf, []string{}, "array of ids that should fail, e.g. '5, 623'")
	// fsFlow.StringArray(flagSkipOnFailureOf, []string{}, "array of ids that should execute successfully, e.g. '5, 623'")
	// fsFlow.StringArray(flagStopOnSuccessOf, []string{}, "array of ids that should execute successfully, will otherwise stop execution'")
	// fsFlow.StringArray(flagStopOnFailureOf, []string{}, "array of ids that should fail, will otherwise stop execution'")
}
