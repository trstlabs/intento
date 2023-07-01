package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	// The connection end identifier on the controller chain
	flagConnectionID = "connection-id"
	// The controller chain channel version
	flagVersion   = "version"
	flagLabel     = "label"
	flagDuration  = "duration"
	flagInterval  = "interval"
	flagStartAt   = "start_at"
	flagDependsOn = "depends_on"
	flagRetries   = "retries"
	flagFeeFunds  = "fee_funds"
	flagEndTime   = "end_at"
)

// common flagsets to add to various functions
var (
	fsConnectionPair = flag.NewFlagSet("", flag.ContinueOnError)
	fsVersion        = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	fsConnectionPair.String(flagConnectionID, "", "Connection ID")
	fsVersion.String(flagVersion, "", "Version")
	//fsConnectionPair.String(FlagCounterpartyConnectionID, "", "Counterparty Connection ID")
}
