package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "blowk <command> [arguments]",
	Short: "blowk is the tool for managing blowK source code",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("blowk encountered an error, exiting")
		panic(err)
	}
}
