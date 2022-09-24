package lang

import (
	"reflect"
	"testing"

	"github.com/bjatkin/blow-k/internal/lex"
)

func Test_getEpresssions(t *testing.T) {
	type args struct {
		tokens []lex.Token
	}
	tests := []struct {
		name string
		args args
		want [][]lex.Token
	}{
		{
			"simple blocks",
			args{
				tokens: []lex.Token{
					{T: lex.ImportKeyword, Value: "import"},
					{T: lex.Identifyer, Value: "echo"},
					{T: lex.SemiColon, Value: ";"},
					{T: lex.ImportKeyword, Value: "import"},
					{T: lex.Identifyer, Value: "grep"},
					{T: lex.AsKeyword, Value: "as"},
					{T: lex.Identifyer, Value: "g"},
					{T: lex.SemiColon, Value: ";"},
				},
			},
			[][]lex.Token{
				{
					{T: lex.ImportKeyword, Value: "import"},
					{T: lex.Identifyer, Value: "echo"},
					{T: lex.SemiColon, Value: ";"},
				},
				{
					{T: lex.ImportKeyword, Value: "import"},
					{T: lex.Identifyer, Value: "grep"},
					{T: lex.AsKeyword, Value: "as"},
					{T: lex.Identifyer, Value: "g"},
					{T: lex.SemiColon, Value: ";"},
				},
			},
		},
		{
			"with functions",
			args{
				tokens: []lex.Token{
					{T: lex.Identifyer, Value: "fiz"},
					{T: lex.Colon, Value: ":"},
					{T: lex.Colon, Value: ":"},
					{T: lex.String, Value: "buz"},
					{T: lex.SemiColon, Value: ";"},
					{T: lex.Identifyer, Value: "main"},
					{T: lex.Colon, Value: ":"},
					{T: lex.OpenParen, Value: "("},
					{T: lex.Identifyer, Value: "args"},
					{T: lex.StringArrayType, Value: "[]string"},
					{T: lex.CloseParen, Value: ")"},
					{T: lex.Colon, Value: ":"},
					{T: lex.OpenBrace, Value: "{"},
					{T: lex.CloseBrace, Value: "}"},
					{T: lex.SemiColon, Value: ";"},
				},
			},
			[][]lex.Token{
				{
					{T: lex.Identifyer, Value: "fiz"},
					{T: lex.Colon, Value: ":"},
					{T: lex.Colon, Value: ":"},
					{T: lex.String, Value: "buz"},
					{T: lex.SemiColon, Value: ";"},
				},
				{
					{T: lex.Identifyer, Value: "main"},
					{T: lex.Colon, Value: ":"},
					{T: lex.OpenParen, Value: "("},
					{T: lex.Identifyer, Value: "args"},
					{T: lex.StringArrayType, Value: "[]string"},
					{T: lex.CloseParen, Value: ")"},
					{T: lex.Colon, Value: ":"},
					{T: lex.OpenBrace, Value: "{"},
					{T: lex.CloseBrace, Value: "}"},
					{T: lex.SemiColon, Value: ";"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Client{}
			if got := c.getExpressions(tt.args.tokens); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getExpressions() = %v, want %v", got, tt.want)
			}
		})
	}
}
