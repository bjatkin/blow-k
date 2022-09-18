package lex

import (
	"fmt"
	"strings"

	"github.com/bjatkin/bear"
	"github.com/bjatkin/blow-k/internal/errors"
)

// TokType is a token type
type TokType int

// String converts a token type into a string representation
func (t TokType) String() string {
	return tokenStrings[t]
}

// MarshalJSON converts the token into a json string
func (t *TokType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t)), nil
}

// UnmarshalJSON converts a json string into a token type
func (t *TokType) UnmarshalJSON(data []byte) error {
	json := strings.Trim(string(data), `"`)
	for i := range tokenStrings {
		if tokenStrings[i] == json {
			*t = TokType(i)
			return nil
		}
	}
	return bear.New(
		bear.WithErrType(errors.InvalidJSON),
		bear.WithTag("data", string(data)),
	)
}

// All valid token types
const (
	Unknown TokType = iota

	ImportKeyword
	AsKeyword

	StartComment
	Comment
	String

	Colon
	Comma
	Exec
	SemiColon
	NewLine
	WhiteSpace

	OpenParen
	CloseParen
	OpenSquare
	CloseSquare
	OpenBrace
	CloseBrace

	StartEndString

	StringType
	StringArrayType

	Identifyer
)

// All valid token strings
var tokenStrings = [...]string{
	"Unknown",
	"ImportKeyword",
	"AsKeyword",
	"StartComment",
	"Comment",
	"String",
	"Colon",
	"Comma",
	"Exec",
	"SemiColon",
	"NewLine",
	"WhiteSpace",
	"OpenParen",
	"CloseParen",
	"OpenSquare",
	"CloseSquare",
	"OpenBrace",
	"CloseBrace",
	"StartEndString",
	"StringType",
	"StringArrayType",
	"Identifyer",
}
