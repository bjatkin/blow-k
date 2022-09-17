package cmd

import (
	"io/ioutil"

	"github.com/bjatkin/bear"
	"github.com/spf13/cobra"

	"github.com/bjatkin/blow-k/internal/errors"
)

func init() {
	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:   "build [source file]",
	Short: "transpile blowK source code into bash",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srcFile := args[0]
		_, err := ioutil.ReadFile(srcFile)
		if err != nil {
			return bear.Wrap(err,
				bear.WithErrType(errors.FileNotFound),
				bear.WithExitCode(errors.BuildFailed),
				bear.WithTag("src name", srcFile),
			)
		}

		// TODO: build the code

		return nil
	},
}
