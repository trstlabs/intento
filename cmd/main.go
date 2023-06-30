package main

import (
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/trstlabs/trst/app"
	cmd "github.com/trstlabs/trst/cmd/trstd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "TRSTD", app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
