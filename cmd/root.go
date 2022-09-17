package cmd

import (
	"github.com/bjatkin/bear"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "blowk <command> [arguments]",
	Short: "blowk is the tool for managing blowK source code",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		berr, _ := bear.AsBerr(err)
		berr.Add(bear.FmtNoID(true), bear.FmtPrettyPrint(true))
		berr.Panic(true)
	}
}
