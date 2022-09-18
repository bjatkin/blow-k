package tests

import (
	"reflect"
	"testing"

	"github.com/bjatkin/blow-k/internal/tok"
)

func TestTokClient(t *testing.T) {
	tests := getTestFiles(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tok.NewClient().Tokenize(tt.fileName)
			if err != nil {
				t.Fatalf("TokClient() unexpected error %v", err)
			}

			if !reflect.DeepEqual(got, tt.wantTokens) {
				t.Fatalf("TokClient() got = %v, wanted %v", got, tt.wantTokens)
			}
		})
	}
}
