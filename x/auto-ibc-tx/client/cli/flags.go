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
	flagStartAt                  = "start_at"
	flagFeeFunds                 = "fee_funds"
	flagEndTime                  = "end_at"

	//Execution conditions
	flagUpdatingDisabled = "updating_disabled"
	flagSaveMsgResponses = "save_msg_responses"
	flagStopOnSuccess    = "stop_on_success"
	flagStopOnFailure    = "stop_on_failure"
	flagStopOnSuccessOf  = "stop_on_success_of"
	flagStopOnFailureOf  = "stop_on_failure_of"
	flagSkipOnSuccessOf  = "skip_on_success_of"
	flagSkipOnFailureOf  = "skip_on_failure_of"
)

// common flagsets to add to various functions
var (
	fsAutoTx  = flag.NewFlagSet("", flag.ContinueOnError)
	fsVersion = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	fsAutoTx.String(flagLabel, "", "A custom label for the AutoTx e.g. AutoTransfer, UpdateContractParams, optional")
	fsAutoTx.String(flagDuration, "", "A custom duration for the AutoTx e.g. 2h, 6000s, 72h3m0.5s, optional")
	fsAutoTx.String(flagInterval, "", "A custom interval for the AutoTx e.g. 2h, 6000s, 72h3m0.5s, optional")
	fsAutoTx.String(flagStartAt, "0", "A custom start time in UNIX time, optional")
	fsAutoTx.String(flagFeeFunds, "", "Coins to sent to limit the fees incurred, optional")
	fsAutoTx.String(flagConnectionID, "", "Connection ID, an IBC ID from this chain to the host chain, optional")
	fsAutoTx.Bool(flagUpdatingDisabled, false, "disable future updates to the configuration'")
	fsAutoTx.Bool(flagSaveMsgResponses, true, "save message responses to tx history (Cosmos SDK v0.46+ chains only)'")
	fsAutoTx.Bool(flagStopOnSuccess, false, "stop execution after success'")
	fsAutoTx.Bool(flagStopOnFailure, false, "stop execution after failure'")
	// fsAutoTx.StringArray(flagSkipOnSuccessOf, []string{}, "array of ids that should fail, e.g. '5, 623'")
	// fsAutoTx.StringArray(flagSkipOnFailureOf, []string{}, "array of ids that should execute successfully, e.g. '5, 623'")
	// fsAutoTx.StringArray(flagStopOnSuccessOf, []string{}, "array of ids that should execute successfully, will otherwise stop execution'")
	// fsAutoTx.StringArray(flagStopOnFailureOf, []string{}, "array of ids that should fail, will otherwise stop execution'")
}
