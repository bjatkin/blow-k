package lex

import (
	"embed"
	"reflect"
	"testing"

	"github.com/bjatkin/blowk/internal/tok"
)

// codeDir contains a dir of src code for testing
//
//go:embed test_data
var codeDir embed.FS

func TestClient_Lex(t *testing.T) {
	type args struct {
		tokens []tok.Token
	}
	tests := []struct {
		name string
		args args
		want []Token
	}{
		{
			"example 1",
			args{
				tokens: func() []tok.Token {
					data, err := codeDir.ReadFile("test_data/example1.bk")
					if err != nil {
						t.Fatalf("Tokenize() failed to setup example: %s", err)
					}
					return tok.Tokenize(string(data))
				}(),
			},
			[]Token{
				{T: ImportKeyword, V: "import"},
				{T: Identifyer, V: "echo"},
				{T: AsKeyword, V: "as"},
				{T: Identifyer, V: "print"},
				{T: SemiColon, V: ";"},
				{T: Comment, V: " main is the entry point for any bk script"},
				{T: SemiColon, V: ";"},
				{T: Identifyer, V: "main"},
				{T: Colon, V: ":"},
				{T: OpenParen, V: "("},
				{T: Identifyer, V: "args"},
				{T: Colon, V: ":"},
				{T: StringArrayType, V: "[]string"},
				{T: CloseParen, V: ")"},
				{T: Colon, V: ":"},
				{T: OpenBrace, V: "{"},
				{T: Exec, V: "$"},
				{T: Identifyer, V: "print"},
				{T: OpenSquare, V: "["},
				{T: String, V: "Hello"},
				{T: Comma, V: ","},
				{T: String, V: "World"},
				{T: CloseSquare, V: "]"},
				{T: SemiColon, V: ";"},
				{T: CloseBrace, V: "}"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient()
			if got := c.Lex(tt.args.tokens); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.Lex() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func TestClient_coalesceStrings(t *testing.T) {
	type args struct {
		tokens []Token
	}
	tests := []struct {
		name string
		args args
		want []Token
	}{
		{
			"single string",
			args{
				tokens: []Token{
					{T: StartEndString, V: "\""},
					{T: WhiteSpace, V: " "},
					{T: Identifyer, V: "Test_"},
					{T: StartEndString, V: "\""},
					{T: NewLine, V: "\n"},
				},
			},
			[]Token{
				{T: String, V: " Test_"},
				{T: NewLine, V: "\n"},
			},
		},
		{
			"multipe strings",
			args{
				tokens: []Token{
					{T: StartEndString, V: "\""},
					{T: WhiteSpace, V: " "},
					{T: Identifyer, V: "Test1"},
					{T: StartEndString, V: "\""},
					{T: Comma, V: ","},
					{T: StartEndString, V: "\""},
					{T: Identifyer, V: "Test"},
					{T: WhiteSpace, V: " "},
					{T: Unknown, V: "2"},
					{T: StartEndString, V: "\""},
					{T: NewLine, V: "\n"},
				},
			},
			[]Token{
				{T: String, V: " Test1"},
				{T: Comma, V: ","},
				{T: String, V: "Test 2"},
				{T: NewLine, V: "\n"},
			},
		},
		{
			"no strings",
			args{
				tokens: []Token{
					{T: Identifyer, V: "main"},
					{T: Colon, V: ":"},
					{T: OpenParen, V: "("},
					{T: CloseParen, V: ")"},
					{T: Colon, V: ":"},
					{T: OpenBrace, V: "{"},
					{T: CloseBrace, V: "}"},
				},
			},
			[]Token{
				{T: Identifyer, V: "main"},
				{T: Colon, V: ":"},
				{T: OpenParen, V: "("},
				{T: CloseParen, V: ")"},
				{T: Colon, V: ":"},
				{T: OpenBrace, V: "{"},
				{T: CloseBrace, V: "}"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := coalesceStrings(tt.args.tokens); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.coalesceStrings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_coalesceComments(t *testing.T) {
	type args struct {
		tokens []Token
	}
	tests := []struct {
		name string
		args args
		want []Token
	}{
		{
			"single comment",
			args{
				tokens: []Token{
					{T: StartComment, V: "#"},
					{T: WhiteSpace, V: " "},
					{T: Identifyer, V: "Test_"},
					{T: NewLine, V: "\n"},
				},
			},
			[]Token{
				{T: Comment, V: " Test_"},
				{T: NewLine, V: "\n"},
			},
		},
		{
			"multipe comments",
			args{
				tokens: []Token{
					{T: StartComment, V: "#"},
					{T: WhiteSpace, V: " "},
					{T: Identifyer, V: "Test1"},
					{T: NewLine, V: "\n"},
					{T: StartComment, V: "#"},
					{T: Identifyer, V: "Test"},
					{T: WhiteSpace, V: " "},
					{T: Unknown, V: "2"},
					{T: NewLine, V: "\n"},
				},
			},
			[]Token{
				{T: Comment, V: " Test1"},
				{T: NewLine, V: "\n"},
				{T: Comment, V: "Test 2"},
				{T: NewLine, V: "\n"},
			},
		},
		{
			"no comments",
			args{
				tokens: []Token{
					{T: Identifyer, V: "main"},
					{T: Colon, V: ":"},
					{T: OpenParen, V: "("},
					{T: CloseParen, V: ")"},
					{T: Colon, V: ":"},
					{T: OpenBrace, V: "{"},
					{T: CloseBrace, V: "}"},
				},
			},
			[]Token{
				{T: Identifyer, V: "main"},
				{T: Colon, V: ":"},
				{T: OpenParen, V: "("},
				{T: CloseParen, V: ")"},
				{T: Colon, V: ":"},
				{T: OpenBrace, V: "{"},
				{T: CloseBrace, V: "}"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := coalesceComments(tt.args.tokens); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.coalesceComments() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_insertSemicolons(t *testing.T) {
	type args struct {
		tokens []Token
	}
	tests := []struct {
		name string
		args args
		want []Token
	}{
		{
			"replace new lines",
			args{
				tokens: []Token{
					{T: ImportKeyword, V: "import"},
					{T: Identifyer, V: "echo"},
					{T: NewLine, V: "\n"},
					{T: Comment, V: "Hello World"},
					{T: NewLine, V: "\n"},
				},
			},
			[]Token{
				{T: ImportKeyword, V: "import"},
				{T: Identifyer, V: "echo"},
				{T: SemiColon, V: ";"},
				{T: Comment, V: "Hello World"},
				{T: SemiColon, V: ";"},
			},
		},
		{
			"replace some new lines",
			args{
				tokens: []Token{
					{T: Identifyer, V: "foo"},
					{T: OpenParen, V: "("},
					{T: NewLine, V: "\n"},
					{T: Identifyer, V: "id"},
					{T: Comma, V: ","},
					{T: NewLine, V: "\n"},
					{T: CloseParen, V: ")"},
					{T: NewLine, V: "\n"},
					{T: Comment, V: "Hello World"},
					{T: NewLine, V: "\n"},
				},
			},
			[]Token{
				{T: Identifyer, V: "foo"},
				{T: OpenParen, V: "("},
				{T: NewLine, V: "\n"},
				{T: Identifyer, V: "id"},
				{T: Comma, V: ","},
				{T: NewLine, V: "\n"},
				{T: CloseParen, V: ")"},
				{T: SemiColon, V: ";"},
				{T: Comment, V: "Hello World"},
				{T: SemiColon, V: ";"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := insertSemicolons(tt.args.tokens); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.insertSemicolons() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func TestClient_filter(t *testing.T) {
	type args struct {
		tokens []Token
	}
	tests := []struct {
		name string
		args args
		want []Token
	}{
		{
			"nothing to remove",
			args{
				tokens: []Token{
					{T: Identifyer, V: "main"},
					{T: Colon, V: ":"},
					{T: OpenParen, V: "("},
					{T: CloseParen, V: ")"},
					{T: Colon, V: ":"},
					{T: OpenBrace, V: "{"},
					{T: CloseBrace, V: "}"},
				},
			},
			[]Token{
				{T: Identifyer, V: "main"},
				{T: Colon, V: ":"},
				{T: OpenParen, V: "("},
				{T: CloseParen, V: ")"},
				{T: Colon, V: ":"},
				{T: OpenBrace, V: "{"},
				{T: CloseBrace, V: "}"},
			},
		},
		{
			"remove new lines and white space",
			args{
				tokens: []Token{
					{T: WhiteSpace, V: " "},
					{T: WhiteSpace, V: " "},
					{T: Identifyer, V: "main"},
					{T: WhiteSpace, V: " "},
					{T: Colon, V: ":"},
					{T: WhiteSpace, V: " "},
					{T: OpenParen, V: "("},
					{T: NewLine, V: "\n"},
					{T: WhiteSpace, V: " "},
					{T: WhiteSpace, V: " "},
					{T: CloseParen, V: ")"},
					{T: Colon, V: ":"},
					{T: WhiteSpace, V: " "},
					{T: NewLine, V: "\n"},
				},
			},
			[]Token{
				{T: Identifyer, V: "main"},
				{T: Colon, V: ":"},
				{T: OpenParen, V: "("},
				{T: CloseParen, V: ")"},
				{T: Colon, V: ":"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filter(tt.args.tokens); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_coalesceArrayTypes(t *testing.T) {
	type args struct {
		tokens []Token
	}
	tests := []struct {
		name string
		args args
		want []Token
	}{
		{
			"string array",
			args{
				tokens: []Token{
					{T: OpenParen, V: "("},
					{T: OpenSquare, V: "["},
					{T: CloseSquare, V: "]"},
					{T: StringType, V: "string"},
					{T: Identifyer, V: "str"},
				},
			},
			[]Token{
				{T: OpenParen, V: "("},
				{T: StringArrayType, V: "[]string"},
				{T: Identifyer, V: "str"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := coalesceArrayTypes(tt.args.tokens); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("coalesceArrayTypes() = %v, want %v", got, tt.want)
			}
		})
	}
}
