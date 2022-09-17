package tests

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/bjatkin/bear"
	"github.com/bjatkin/blow-k/internal/errors"
	"github.com/bjatkin/blow-k/internal/tok"
)

// testFiles groups src files with the respective test data
type testCase struct {
	name       string
	fileName   string
	wantTokens []tok.Token
}

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
			return bear.Wrap(err,
				bear.WithErrType(errors.FileNotFound),
				bear.WithTag("file path", bkPath),
				bear.FmtPrettyPrint(true),
			)
		}

		tokPath := filepath.Join(path, info.Name()+"_tok")
		tokFile, err := os.ReadFile(tokPath)
		if err != nil {
			return bear.Wrap(err,
				bear.WithErrType(errors.FileNotFound),
				bear.WithTag("file path", tokPath),
				bear.FmtPrettyPrint(true),
			)
		}
		tokens := &[]tok.Token{}
		err = json.Unmarshal(tokFile, tokens)
		if err != nil {
			return bear.Wrap(err,
				bear.WithErrType(errors.InvalidJSON),
				bear.WithTag("JSON", string(tokFile)),
				bear.FmtPrettyPrint(true),
			)
		}

		files = append(files, testCase{
			name:       info.Name(),
			fileName:   bkPath,
			wantTokens: *tokens,
		})

		return nil
	})
	if err != nil {
		t.Fatalf("getTestFiles() failed to get test files %s", err)
	}

	return files
}
