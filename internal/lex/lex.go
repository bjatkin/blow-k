package lex

import (
	"regexp"
	"strings"

	"github.com/bjatkin/bear"
	"github.com/bjatkin/blowk/internal/errors"
	"github.com/bjatkin/blowk/internal/tok"
)

// TokType is a token type
type TokType int

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

// Token is a valid token string
type Token struct {
	T TokType
	V string
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

func (m matcher) match(token tok.Token) bool {
	sToken := string(token)
	if m.str != "" {
		return sToken == m.str
	}

	if m.reg != nil {
		return m.reg.MatchString(sToken)
	}

	// panic if str and reg are both unset
	errors.Base.New(
		bear.WithExitCode(errors.LexerFailed),
		bear.WithErrType(errors.InvalidLexMatcher),
	).Panic(true)

	return false
}

// Client is a lex client that lexes an array of tokens
type Client struct {
	matchers     []matcher
	transformers []transform
}

// NewClient creates a new lex client
func NewClient() *Client {
	return &Client{
		// default matchers
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

		// default transformers
		transformers: []transform{
			coalesceStrings,
			coalesceArrayTypes,
			coalesceComments,
			insertSemicolons,
			filter,
		},
	}
}

// Lex converts a []tok.Token into a []lex.Tokens
func (c *Client) Lex(tokens []tok.Token) []Token {
	var ret []Token

	for _, tok := range tokens {
		t := Token{V: string(tok)}
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

type transform func([]Token) []Token

func coalesceStrings(tokens []Token) []Token {
	var ret []Token
	var current string
	var open bool

	for _, tok := range tokens {
		if tok.T == StartEndString && open {
			ret = append(ret, Token{T: String, V: current})
			open = false
			current = ""
			continue
		}

		if tok.T == StartEndString && !open {
			open = true
			continue
		}

		if open {
			current += tok.V
			continue
		}

		ret = append(ret, tok)
	}

	return ret
}

func coalesceComments(tokens []Token) []Token {
	var ret []Token
	var current string
	var open bool

	for _, tok := range tokens {
		if tok.T == NewLine && open {
			ret = append(ret, Token{T: Comment, V: current})
			current = ""
			open = false
		}

		if tok.T == StartComment {
			open = true
			continue
		}

		if open {
			current += tok.V
			continue
		}

		ret = append(ret, tok)
	}

	return ret
}

func coalesceArrayTypes(tokens []Token) []Token {
	var ret []Token
	for i := 0; i < len(tokens); i++ {
		tok := tokens[i]
		if i+2 > len(tokens) {
			ret = append(ret, tok)
			continue
		}

		// check for leading []
		if tokens[i].T != OpenSquare ||
			tokens[i+1].T != CloseSquare {
			ret = append(ret, tok)
			continue
		}

		// check for type
		if tokens[i+2].T == StringType {
			ret = append(ret, Token{T: StringArrayType, V: "[]string"})
			i += 2
			continue
		}

		ret = append(ret, tok)
	}

	return ret
}

func insertSemicolons(tokens []Token) []Token {
	for i, tok := range tokens {
		if tok.T != NewLine || i == 0 {
			continue
		}

		prev := tokens[i-1].T
		if prev == OpenBrace ||
			prev == StartEndString ||
			prev == OpenParen ||
			prev == OpenSquare ||
			prev == Comma ||
			prev == SemiColon {
			continue
		}

		tokens[i] = Token{T: SemiColon, V: ";"}
	}

	return tokens
}

func filter(tokens []Token) []Token {
	var ret []Token

	for _, tok := range tokens {
		if tok.T == WhiteSpace ||
			tok.T == NewLine {
			continue
		}

		ret = append(ret, tok)
	}

	return ret
}
