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

// New returns a base new error
func New(opts ...bear.ErrOption) *bear.Error {
	return Base.New(opts...)
}

// Wrap returns a base wrapped error
func Wrap(err error, opts ...bear.ErrOption) *bear.Error {
	return Base.Wrap(err, opts...)
}
