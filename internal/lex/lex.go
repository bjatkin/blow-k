package lex

import (
	"regexp"
	"strings"

	"github.com/bjatkin/bear"
	"github.com/bjatkin/blow-k/internal/errors"
	"github.com/bjatkin/blow-k/internal/tok"
)

// Token is a valid token string
type Token struct {
	T          TokType
	Value      string
	FileName   string
	LineNumber int
	ColNumber  int
}

// Client is a lex client that takes a slice of tok.Tokens and returns a slice of lex.Tokens
type Client struct {
	matchers     []matcher
	transformers []transformer
}

// NewClient creates a new default lex.Client
func NewClient() *Client {
	return &Client{
		matchers: []matcher{
			newSMatcher("import", ImportKeyword),
			newSMatcher("as", AsKeyword),
			newSMatcher("#", StartComment),
			newSMatcher(":", Colon),
			newSMatcher(",", Comma),
			newSMatcher("$", Exec),
			newSMatcher(";", SemiColon),
			newSMatcher("\n", NewLine),
			newSMatcher(" ", WhiteSpace),
			newSMatcher("(", OpenParen),
			newSMatcher(")", CloseParen),
			newSMatcher("[", OpenSquare),
			newSMatcher("]", CloseSquare),
			newSMatcher("{", OpenBrace),
			newSMatcher("}", CloseBrace),
			newSMatcher("\"", StartEndString),
			newSMatcher("string", StringType),
			newRMatcher(`[a-zA-Z][a-zA-Z0-9_]+`, Identifyer),
		},
		transformers: []transformer{
			coalesceStrings,
			coalesceArrayTypes,
			coalesceComments,
			insertSemicolons,
			filter,
		},
	}
}

// Lex converts a slice of tok.Token into a slice of lex.Tokens
func (c *Client) Lex(tokens []tok.Token) []Token {
	var ret []Token

	for _, tok := range tokens {
		t := Token{
			Value:      tok.Value,
			FileName:   tok.FileName,
			LineNumber: tok.LineNumber,
			ColNumber:  tok.ColNumber,
		}
		for _, matcher := range c.matchers {
			if matcher.match(tok) {
				t.T = matcher.t
				break
			}
		}
		ret = append(ret, t)
	}

	// do token transformations
	for _, t := range c.transformers {
		ret = t(ret)
	}

	return ret
}

// matcher matches a token with a token type
type matcher struct {
	str string
	reg *regexp.Regexp
	t   TokType
}

// newSMatcher creates a new string matcher
func newSMatcher(str string, t TokType) matcher {
	return matcher{str: str, t: t}
}

// newRMatcher creates a new regex matcher
// if the provided regex does not complie this function will panic
func newRMatcher(reg string, t TokType) matcher {
	if !strings.HasPrefix("^", reg) {
		reg = "^" + reg
	}
	if !strings.HasSuffix("$", reg) {
		reg += "$"
	}

	return matcher{reg: regexp.MustCompile(reg), t: t}
}

// match returns true if the matcher matches the token
func (m matcher) match(token tok.Token) bool {
	sToken := token.Value
	if m.str != "" {
		return sToken == m.str
	}

	if m.reg != nil {
		return m.reg.MatchString(sToken)
	}

	// panic if str and reg are both unset
	errors.New(
		bear.WithExitCode(errors.LexerFailed),
		bear.WithErrType(errors.InvalidLexMatcher),
	).Panic(true)

	return false
}
