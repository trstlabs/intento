package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	// The connection end identifier on the controller chain
	FlagConnectionID = "connection-id"
	// The connection end identifier on the host chain
	//FlagCounterpartyConnectionID = "counterparty-connection-id"
	flagDuration  = "duration"
	flagInterval  = "interval"
	flagStartAt   = "start_at"
	flagDependsOn = "depends_on"
	flagRetries   = "retries"
)

// common flagsets to add to various functions
var (
	fsConnectionPair = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	fsConnectionPair.String(FlagConnectionID, "", "Connection ID")
	//fsConnectionPair.String(FlagCounterpartyConnectionID, "", "Counterparty Connection ID")
}
