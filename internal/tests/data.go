package tests

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/tabwriter"

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
			return errors.Wrap(err,
				bear.WithErrType(errors.FileNotFound),
				bear.WithTag("file path", bkPath),
			)
		}

		tokens, err := getTokenFile(path, info.Name())
		if err != nil {
			return errors.Wrap(err, bear.WithLabels("token file"))
		}

		lexs, err := getLexFile(path, info.Name())
		if err != nil {
			return errors.Wrap(err, bear.WithLabels("lex file"))
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
		return nil, errors.Wrap(err,
			bear.WithErrType(errors.FileNotFound),
			bear.WithTag("file path", tokPath),
			bear.FmtPrettyPrint(true),
		)
	}
	tokens := &[]tok.Token{}
	err = json.Unmarshal(tokFile, tokens)
	if err != nil {
		return nil, errors.Wrap(err,
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
		return nil, errors.Wrap(err,
			bear.WithErrType(errors.FileNotFound),
			bear.WithTag("file path", lexPath),
			bear.FmtPrettyPrint(true),
		)
	}
	tokens := &[]lex.Token{}
	err = json.Unmarshal(lexFile, tokens)
	if err != nil {
		return nil, errors.Wrap(err,
			bear.WithErrType(errors.InvalidJSON),
			bear.WithTag("JSON", string(lexFile)),
		)
	}

	return *tokens, nil
}

// buildCompTable creates a table with a as the left column and b as the right
func buildCompTable[T any](a, b []T) string {
	buf := &strings.Builder{}
	w := tabwriter.NewWriter(buf, 1, 2, 1, ' ', 0)

	max := len(a)
	if len(b) > max {
		max = len(b)
	}

	for i := 0; i < max; i++ {
		var left, right string
		if len(a) > i {
			left = fmt.Sprintf("%v", a[i])
			left = strings.ReplaceAll(left, "\n", "\\n")
		}
		if len(b) > i {
			right = fmt.Sprintf("%v", b[i])
			right = strings.ReplaceAll(right, "\n", "\\n")
		}

		match := "-"
		if left != right {
			match = "X"
		}

		fmt.Fprintf(w, "%s\t%s\t%s\n", match, left, right)
	}
	w.Flush()

	return buf.String()
}
