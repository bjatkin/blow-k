package pat

import (
	"testing"

	"github.com/bjatkin/blowk/internal/lex"
)

func TestExactMatch_match(t *testing.T) {
	type fields struct {
		tok []lex.TokType
	}
	type args struct {
		check []lex.Token
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  bool
	}{
		{
			"success",
			fields{
				tok: []lex.TokType{lex.ImportKeyword, lex.WhiteSpace, lex.Identifyer},
			},
			args{
				check: []lex.Token{
					{T: lex.ImportKeyword, V: "import"},
					{T: lex.WhiteSpace, V: " "},
					{T: lex.Identifyer, V: "echo"},
					{T: lex.NewLine, V: "\n"},
					{T: lex.StartComment, V: "#"},
					{T: lex.WhiteSpace, V: " "},
					{T: lex.Identifyer, V: "comment"},
				},
			},
			3,
			true,
		},
		{
			"failure",
			fields{
				tok: []lex.TokType{lex.ImportKeyword, lex.WhiteSpace, lex.Identifyer},
			},
			args{
				check: []lex.Token{
					{T: lex.Identifyer, V: "main"},
					{T: lex.WhiteSpace, V: " "},
					{T: lex.Colon, V: ":"},
					{T: lex.OpenParen, V: "("},
					{T: lex.CloseParen, V: ")"},
					{T: lex.Colon, V: ":"},
					{T: lex.OpenBrace, V: "{"},
					{T: lex.CloseBrace, V: "}"},
				},
			},
			0,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewExact(tt.fields.tok...)
			got, got1 := m.Match(tt.args.check)
			if got != tt.want {
				t.Errorf("ExactMatch.match() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ExactMatch.match() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestOneOfMatch_match(t *testing.T) {
	type fields struct {
		tok []lex.TokType
	}
	type args struct {
		check []lex.Token
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  bool
	}{
		{
			"success",
			fields{
				tok: []lex.TokType{lex.WhiteSpace, lex.NewLine, lex.ImportKeyword},
			},
			args{
				check: []lex.Token{
					{T: lex.NewLine, V: "\n"},
				},
			},
			1,
			true,
		},
		{
			"failure",
			fields{
				tok: []lex.TokType{lex.WhiteSpace, lex.NewLine, lex.ImportKeyword},
			},
			args{
				check: []lex.Token{
					{T: lex.Identifyer, V: "main"},
				},
			},
			0,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewOneOf(tt.fields.tok...)
			got, got1 := m.Match(tt.args.check)
			if got != tt.want {
				t.Errorf("OneOfMatch.match() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("OneOfMatch.match() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestOneOrMore_match(t *testing.T) {
	type fields struct {
		mat Pattern
	}
	type args struct {
		check []lex.Token
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  bool
	}{
		{
			"single match",
			fields{
				NewOneOf(lex.StartComment, lex.StringType),
			},
			args{
				check: []lex.Token{
					{T: lex.StartComment},
					{T: lex.WhiteSpace},
					{T: lex.Identifyer},
				},
			},
			1,
			true,
		},
		{
			"multiple match",
			fields{
				NewOneOf(lex.StartComment, lex.StringType),
			},
			args{
				check: []lex.Token{
					{T: lex.StartComment},
					{T: lex.StringType},
					{T: lex.WhiteSpace},
				},
			},
			2,
			true,
		},
		{
			"no match",
			fields{
				NewExact(lex.StartComment, lex.StringType),
			},
			args{
				check: []lex.Token{
					{T: lex.StartComment},
					{T: lex.WhiteSpace},
					{T: lex.Identifyer},
				},
			},
			0,
			false,
		},
		{
			"no args",
			fields{
				NewExact(lex.StartComment, lex.StringType),
			},
			args{},
			0,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewOneOrMore(tt.fields.mat)
			got, got1 := m.Match(tt.args.check)
			if got != tt.want {
				t.Errorf("OneOrMore.match() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("OneOrMore.match() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestCompositMatcher_Match(t *testing.T) {
	type fields struct {
		mats []Pattern
	}
	type args struct {
		check []lex.Token
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  bool
	}{
		{
			"single matcher",
			fields{
				mats: []Pattern{NewExact(lex.ImportKeyword, lex.WhiteSpace, lex.Identifyer)},
			},
			args{
				check: []lex.Token{
					{T: lex.ImportKeyword},
					{T: lex.WhiteSpace},
					{T: lex.Identifyer},
				},
			},
			3,
			true,
		},
		{
			"multiple patterns",
			fields{
				mats: []Pattern{NewExact(lex.ImportKeyword, lex.WhiteSpace, lex.Identifyer, lex.NewLine)},
			},
			args{
				check: []lex.Token{
					{T: lex.ImportKeyword},
					{T: lex.WhiteSpace},
					{T: lex.Identifyer},
					{T: lex.NewLine},
					{T: lex.ImportKeyword},
					{T: lex.WhiteSpace},
					{T: lex.Identifyer},
					{T: lex.NewLine},
				},
			},
			4,
			true,
		},
		{
			"failed matcher",
			fields{
				mats: []Pattern{NewOneOf(lex.Identifyer, lex.StringType)},
			},
			args{
				check: []lex.Token{
					{T: lex.AsKeyword},
					{T: lex.WhiteSpace},
					{T: lex.Identifyer},
				},
			},
			0,
			false,
		},
		{
			"empty args",
			fields{
				mats: []Pattern{NewOneOf(lex.Identifyer, lex.StringType)},
			},
			args{},
			0,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewComposit(tt.fields.mats...)
			got, got1 := m.Match(tt.args.check)
			if got != tt.want {
				t.Errorf("CompositMatcher.Match() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CompositMatcher.Match() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestZeroOrMore_Match(t *testing.T) {
	type fields struct {
		mat Pattern
	}
	type args struct {
		check []lex.Token
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  bool
	}{
		{
			"single match",
			fields{
				NewOneOf(lex.StartComment, lex.StringType),
			},
			args{
				check: []lex.Token{
					{T: lex.StartComment},
					{T: lex.WhiteSpace},
					{T: lex.Identifyer},
				},
			},
			1,
			true,
		},
		{
			"multiple match",
			fields{
				NewOneOf(lex.StartComment, lex.StringType),
			},
			args{
				check: []lex.Token{
					{T: lex.StartComment},
					{T: lex.StringType},
					{T: lex.WhiteSpace},
				},
			},
			2,
			true,
		},
		{
			"no match",
			fields{
				NewExact(lex.StartComment, lex.StringType),
			},
			args{
				check: []lex.Token{
					{T: lex.StartComment},
					{T: lex.WhiteSpace},
					{T: lex.Identifyer},
				},
			},
			0,
			true,
		},
		{
			"no args",
			fields{
				NewExact(lex.StartComment, lex.StringType),
			},
			args{},
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewZeroOrMore(tt.fields.mat)
			got, got1 := m.Match(tt.args.check)
			if got != tt.want {
				t.Errorf("ZeroOrMore.Match() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ZeroOrMore.Match() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNot_Match(t *testing.T) {
	type fields struct {
		t lex.TokType
	}
	type args struct {
		check []lex.Token
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  bool
	}{
		{
			"match",
			fields{
				t: lex.NewLine,
			},
			args{
				check: []lex.Token{
					{T: lex.AsKeyword},
				},
			},
			1,
			true,
		},
		{
			"not matched",
			fields{
				t: lex.WhiteSpace,
			},
			args{
				check: []lex.Token{
					{T: lex.WhiteSpace},
				},
			},
			0,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewNot(tt.fields.t)
			got, got1 := m.Match(tt.args.check)
			if got != tt.want {
				t.Errorf("Not.Match() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Not.Match() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestBlock_Match(t *testing.T) {
	type fields struct {
		open  lex.TokType
		close lex.TokType
	}
	type args struct {
		check []lex.Token
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  bool
	}{
		{
			"no block",
			fields{
				open:  lex.StartEndString,
				close: lex.StartEndString,
			},
			args{
				check: []lex.Token{
					{T: lex.AsKeyword},
					{T: lex.WhiteSpace},
					{T: lex.Identifyer},
				},
			},
			0,
			false,
		},
		{
			"empty check",
			fields{
				open:  lex.StartEndString,
				close: lex.StartEndString,
			},
			args{},
			0,
			false,
		},
		{
			"string block",
			fields{
				open:  lex.StartEndString,
				close: lex.StartEndString,
			},
			args{
				check: []lex.Token{
					{T: lex.StartEndString},
					{T: lex.AsKeyword},
					{T: lex.StartEndString},
				},
			},
			3,
			true,
		},
		{
			"nested function block",
			fields{
				open:  lex.OpenBrace,
				close: lex.CloseBrace,
			},
			args{
				check: []lex.Token{
					{T: lex.OpenBrace},
					{T: lex.OpenBrace},
					{T: lex.OpenBrace},
					{T: lex.ImportKeyword},
					{T: lex.CloseBrace},
					{T: lex.CloseBrace},
					{T: lex.CloseBrace},
				},
			},
			7,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewBlock(tt.fields.open, tt.fields.close)
			got, got1 := m.Match(tt.args.check)
			if got != tt.want {
				t.Errorf("Block.Match() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Block.Match() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
