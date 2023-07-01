package main

import (
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/trstlabs/trst/app"
	cmd "github.com/trstlabs/trst/cmd/trstd/cmd"
	cmdcfg "github.com/trstlabs/trst/cmd/trstd/cmd/config"
)

func main() {

	cmdcfg.SetupConfig()
	cmdcfg.RegisterDenoms()

	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "TRSTD", app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
