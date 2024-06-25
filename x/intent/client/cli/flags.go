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
	flagUpdatingDisabled          = "updating-disabled"
	flagSaveMsgResponses          = "save-msg-responses"
	flagFallbackToOwnerBalance    = "fallback-to-owner-balance"
	flagStopOnSuccess             = "stop-on-success"
	flagStopOnFailure             = "stop-on-failure"
	flagStopOnSuccessOf           = "stop-on-success-of"
	flagStopOnFailureOf           = "stop-on-failure-of"
	flagSkipOnSuccessOf           = "skip-on-success-of"
	flagSkipOnFailureOf           = "skip-on-failure-of"
	flagReregisterICAAfterTimeout = "reregister-ica-after-timeout"

	flagFeeCoinsSupported = "fee-coins-suported"
	flagNewAdmin          = "new-admin"
)

// common flagsets to add to various functions
var (
	fsAction = flag.NewFlagSet("", flag.ContinueOnError)
	fsIBC    = flag.NewFlagSet("", flag.ContinueOnError)
	// fsVersion = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	fsAction.String(flagLabel, "", "A custom label for the action e.g. AutoTransfer, UpdateContractParams, optional")
	fsAction.String(flagInterval, "", "A custom interval for the action e.g. 2h, 6000s, 72h3m0.5s, optional")
	fsAction.String(flagStartAt, "0", "A custom start time in UNIX time, optional")
	fsAction.String(flagFeeFunds, "", "Coin to attach to the action, optional")
	fsAction.String(flagConnectionID, "", "Connection ID from this chain to the host chain, optional")
	fsAction.String(flagHostConnectionID, "", "Connection ID from host chain to this chain, optional")
	fsAction.String(flagHostedAccount, "", "A hosted account to execute actions on a host, optional")
	fsAction.String(flagHostedAccountFeeLimit, "", "Coin to sent to limit the hosted fees, optional")
	fsAction.Bool(flagUpdatingDisabled, false, "disable future updates to the configuration'")
	fsAction.Bool(flagSaveMsgResponses, true, "save message responses to tx history (Cosmos SDK v0.46+ chains only)'")
	fsAction.Bool(flagStopOnSuccess, false, "stop execution after success'")
	fsAction.Bool(flagStopOnFailure, false, "stop execution after failure'")
	fsAction.Bool(flagFallbackToOwnerBalance, false, "fallback to owner balance'")
	fsAction.Bool(flagReregisterICAAfterTimeout, true, " If true, allows the action to continue execution after an ibc channel times out (recommended)'")

	fsIBC.String(flagConnectionID, "", "Connection ID from this chain to the host chain, optional")
	fsIBC.String(flagHostConnectionID, "", "Connection ID from host chain to this chain, optional")

	// fsAction.StringArray(flagSkipOnSuccessOf, []string{}, "array of ids that should fail, e.g. '5, 623'")
	// fsAction.StringArray(flagSkipOnFailureOf, []string{}, "array of ids that should execute successfully, e.g. '5, 623'")
	// fsAction.StringArray(flagStopOnSuccessOf, []string{}, "array of ids that should execute successfully, will otherwise stop execution'")
	// fsAction.StringArray(flagStopOnFailureOf, []string{}, "array of ids that should fail, will otherwise stop execution'")
}
