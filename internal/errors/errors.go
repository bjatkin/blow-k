package errors

import "github.com/bjatkin/bear"

// Error Types
var (
	FileNotFound      = bear.NewType("File Not Found")
	InvalidLexMatcher = bear.NewType("Invalid Lex Matcher")
	CodeGenFailed     = bear.NewType("Code Gen Failed")
)

// Exit Codes
const (
	BuildFailed = iota + 1
	LexerFailed
	TokenizerFailed
	BashGenFailed
)

// base error template
var Base = bear.NewTemplate(bear.FmtPrettyPrint(true), bear.FmtNoID(true))
