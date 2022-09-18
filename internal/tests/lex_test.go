package tests

import (
	"reflect"
	"testing"

	"github.com/bjatkin/blow-k/internal/lex"
	"github.com/bjatkin/blow-k/internal/tok"
)

func TestLexClient(t *testing.T) {
	tests := getTestFiles(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := tok.NewClient().Tokenize(tt.fileName)
			if err != nil {
				t.Fatalf("LexClient() unexpected error %v", err)
			}

			got := lex.NewClient().Lex(tokens)
			if !reflect.DeepEqual(got, tt.wantLexs) {
				t.Fatalf("LexClient() = \n%v, but wanted \n%v", got, tt.wantLexs)
			}
		})
	}
}
