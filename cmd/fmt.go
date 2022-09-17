package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(fmtCmd)
}

var fmtCmd = &cobra.Command{
	Use:   "fmt [source file]",
	Short: "format and re-write the source file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: fmt the code

		return nil
	},
}
