package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	// The connection end identifier on the controller chain
	flagConnectionID = "connection-id"
	// The controller chain channel version
	flagCounterpartyConnectionID = "counterparty-connection-id"
	flagLabel                    = "label"
	flagDuration                 = "duration"
	flagInterval                 = "interval"
	flagStartAt                  = "start-at"
	flagFeeFunds                 = "fee-funds"
	flagEndTime                  = "end-at"

	//Execution conditions
	flagUpdatingDisabled       = "updating-disabled"
	flagSaveMsgResponses       = "save-msg-responses"
	flagFallbackToOwnerBalance = "fallback-to-owner-balance"
	flagStopOnSuccess          = "stop-on-success"
	flagStopOnFailure          = "stop-on-failure"
	flagStopOnSuccessOf        = "stop-on-success-of"
	flagStopOnFailureOf        = "stop-on-failure-of"
	flagSkipOnSuccessOf        = "skip-on-success-of"
	flagSkipOnFailureOf        = "skip-on-failure-of"
)

// common flagsets to add to various functions
var (
	fsAutoTx  = flag.NewFlagSet("", flag.ContinueOnError)
	fsVersion = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	fsAutoTx.String(flagLabel, "", "A custom label for the AutoTx e.g. AutoTransfer, UpdateContractParams, optional")
	fsAutoTx.String(flagInterval, "", "A custom interval for the AutoTx e.g. 2h, 6000s, 72h3m0.5s, optional")
	fsAutoTx.String(flagStartAt, "0", "A custom start time in UNIX time, optional")
	fsAutoTx.String(flagFeeFunds, "", "Coins to sent to limit the fees incurred, optional")
	fsAutoTx.String(flagConnectionID, "", "Connection ID, an IBC ID from this chain to the host chain, optional")
	fsAutoTx.Bool(flagUpdatingDisabled, false, "disable future updates to the configuration'")
	fsAutoTx.Bool(flagSaveMsgResponses, true, "save message responses to tx history (Cosmos SDK v0.46+ chains only)'")
	fsAutoTx.Bool(flagStopOnSuccess, false, "stop execution after success'")
	fsAutoTx.Bool(flagStopOnFailure, false, "stop execution after failure'")
	fsAutoTx.Bool(flagFallbackToOwnerBalance, false, "fallback to owner balance'")
	// fsAutoTx.StringArray(flagSkipOnSuccessOf, []string{}, "array of ids that should fail, e.g. '5, 623'")
	// fsAutoTx.StringArray(flagSkipOnFailureOf, []string{}, "array of ids that should execute successfully, e.g. '5, 623'")
	// fsAutoTx.StringArray(flagStopOnSuccessOf, []string{}, "array of ids that should execute successfully, will otherwise stop execution'")
	// fsAutoTx.StringArray(flagStopOnFailureOf, []string{}, "array of ids that should fail, will otherwise stop execution'")
}
