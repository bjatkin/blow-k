package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bjatkin/bear"
	"github.com/bjatkin/blowk/internal/ast"
	"github.com/bjatkin/blowk/internal/errors"
	"github.com/bjatkin/blowk/internal/lex"
	"github.com/bjatkin/blowk/internal/tok"
	"github.com/spf13/cobra"
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
		srcCode, err := ioutil.ReadFile(srcFile)
		if err != nil {
			return bear.Wrap(err,
				bear.WithErrType(errors.FileNotFound),
				bear.WithExitCode(errors.BuildFailed),
				bear.WithTag("file name", srcFile),
			)
		}

		wd, err := os.Getwd()
		if err != nil {
			return bear.Wrap(err,
				bear.WithExitCode(errors.BuildFailed),
			)
		}

		fmt.Printf("Compiling %s\n", srcFile)

		tokens := tok.Tokenize(string(srcCode))

		fmt.Println("TOKENS: ", tokens)

		lexClient := lex.NewClient()
		lexedTokens := lexClient.Lex(tokens)

		fmt.Println("LEXED TOKENS: ", lexedTokens)

		astClient := ast.NewTopLevelClient()
		ast := astClient.Parse(lexedTokens)
		code, err := ast.BashGen()
		if err != nil {
			return bear.Wrap(err)
		}

		fileName := filepath.Base(wd)
		fmt.Println("GENERATED CODE: ", code)
		os.WriteFile(filepath.Join(wd, fileName+".sh"), []byte(code), 0777)

		return nil
	},
}
