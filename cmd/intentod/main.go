package main

import (
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/trstlabs/intento/app"
	cmd "github.com/trstlabs/intento/cmd/intentod/cmd"
	cmdcfg "github.com/trstlabs/intento/cmd/intentod/cmd/config"
)

func main() {

	cmdcfg.SetupConfig()
	cmdcfg.RegisterDenoms()

	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "INTENTOD", app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
