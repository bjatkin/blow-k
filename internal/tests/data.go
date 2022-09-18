package tests

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/bjatkin/bear"
	"github.com/bjatkin/blow-k/internal/errors"
	"github.com/bjatkin/blow-k/internal/lex"
	"github.com/bjatkin/blow-k/internal/tok"
)

// testFiles groups src files with the respective test data
type testCase struct {
	name       string
	fileName   string
	wantTokens []tok.Token
	wantLexs   []lex.Token
}

// getTestFiles creates test files from the data directory
func getTestFiles(t *testing.T) []testCase {
	var files []testCase
	err := filepath.Walk("data", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		if path == "data" {
			return nil
		}

		bkPath := filepath.Join(path, info.Name()+".bk")
		if _, err := os.Stat(bkPath); err != nil {
			return errors.Base.Wrap(err,
				bear.WithErrType(errors.FileNotFound),
				bear.WithTag("file path", bkPath),
				bear.FmtPrettyPrint(true),
			)
		}

		tokens, err := getTokenFile(path, info.Name())
		if err != nil {
			return err
		}

		lexs, err := getLexFile(path, info.Name())
		if err != nil {
			return err
		}

		files = append(files, testCase{
			name:       info.Name(),
			fileName:   bkPath,
			wantTokens: tokens,
			wantLexs:   lexs,
		})

		return nil
	})
	if err != nil {
		t.Fatalf("getTestFiles() failed to get test files %s", err)
	}

	return files
}

// getTokenFile reads a token file from the given directory
func getTokenFile(path string, name string) ([]tok.Token, error) {
	tokPath := filepath.Join(path, name+"_tok")
	tokFile, err := os.ReadFile(tokPath)
	if err != nil {
		return nil, errors.Base.Wrap(err,
			bear.WithErrType(errors.FileNotFound),
			bear.WithTag("file path", tokPath),
			bear.FmtPrettyPrint(true),
		)
	}
	tokens := &[]tok.Token{}
	err = json.Unmarshal(tokFile, tokens)
	if err != nil {
		return nil, errors.Base.Wrap(err,
			bear.WithErrType(errors.InvalidJSON),
			bear.WithTag("JSON", string(tokFile)),
		)
	}

	return *tokens, nil
}

// getLexFile reads a token file from the given directory
func getLexFile(path string, name string) ([]lex.Token, error) {
	lexPath := filepath.Join(path, name+"_lex")
	lexFile, err := os.ReadFile(lexPath)
	if err != nil {
		return nil, errors.Base.Wrap(err,
			bear.WithErrType(errors.FileNotFound),
			bear.WithTag("file path", lexPath),
			bear.FmtPrettyPrint(true),
		)
	}
	tokens := &[]lex.Token{}
	err = json.Unmarshal(lexFile, tokens)
	if err != nil {
		return nil, errors.Base.Wrap(err,
			bear.WithErrType(errors.InvalidJSON),
			bear.WithTag("JSON", string(lexFile)),
		)
	}

	fmt.Println("got lex", tokens)

	return *tokens, nil
}
