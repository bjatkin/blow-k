package errors

import "github.com/bjatkin/bear"

// Error Types
var (
	FileNotFound      = bear.NewType("File Not Found")
	InvalidJSON       = bear.NewType("Invalid JSON")
	InvalidLexMatcher = bear.NewType("Invalid Lex Matcher")
)

// Exit Codes
const (
	BuildFailed = iota + 1
	TokenizerFailed
	LexerFailed
)

// base error template
var Base = bear.NewTemplate(bear.FmtPrettyPrint(true), bear.FmtNoID(true))
