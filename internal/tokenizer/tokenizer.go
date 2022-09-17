package tok

import (
	"regexp"
	"strings"

	"github.com/bjatkin/bear"
	"github.com/bjatkin/blowk/internal/errors"
)

// Error Types
var (
	MatchErr        = bear.NewType("Matcher Error")
	UnknownTokenErr = bear.NewType("Unknown Token Error")
)

type TokType int

// All token types
const (
	NewLine TokType = iota
	WhiteSpace
	Comment
	FN
	Let
	Import
	Identifyer
	OpenParen
	CloseParen
	OpenBracket
	CloseBracket
	String
)

// All the valid matchers
var (
	matchers = []*matcher{
		newRegMatcher("white space", WhiteSpace, `[\t ]+`),
		newRegMatcher("comment", Comment, `#[^\n]*`),
		newRegMatcher("string", String, `"[^"]*"`),

		newStrMatcher("new line", NewLine, "\n"),
		newStrMatcher("function", FN, "fn"),
		newStrMatcher("let", Let, "let"),
		newStrMatcher("import", Import, "import"),
		newStrMatcher("open paren", OpenParen, "("),
		newStrMatcher("open paren", CloseParen, ")"),
		newStrMatcher("open bracket", OpenBracket, "{"),
		newStrMatcher("close bracket", OpenBracket, "}"),

		// keep this matcher last
		newRegMatcher("identifyer", Identifyer, `[a-zA-Z][a-zA-Z0-9]*`),
	}
)

// Token is a blowK token
type Token struct {
	T     TokType
	Value string
}

// rlen returns the len of the tokens value in runes
func (t Token) rlen() int {
	return len([]rune(t.Value))
}

// matcher is a matcher for a specific token
type matcher struct {
	name     string
	strMatch string
	regMatch *regexp.Regexp
	t        TokType
}

// newStrMatcher creates a new tokenMatcher with a string match
func newStrMatcher(name string, t TokType, match string) *matcher {
	return &matcher{
		name:     name,
		t:        t,
		strMatch: match,
	}
}

// newStrMatcher creates a new tokenMatcher with a regexp match
// if the regexp does not complie newRegMatcher will panic
func newRegMatcher(name string, t TokType, match string) *matcher {
	if !strings.HasPrefix(match, "^") {
		match = "^" + match
	}

	return &matcher{
		name:     name,
		t:        t,
		regMatch: regexp.MustCompile(match),
	}
}

// match returns true if the matcher matches the string
// the matching string is also returned
func (t *matcher) match(check string) (string, bool) {
	switch {
	case t.strMatch != "":
		l := len(t.strMatch)
		if l > len(check) {
			break
		}

		if t.strMatch == check[:l] {
			return check[:l], true
		}
	case t.regMatch != nil:
		if match := t.regMatch.FindString(check); match != "" {
			return match, true
		}
	default:
		// if it's not a str or reg matcher panic
		errors.Base.New(
			bear.WithErrType(MatchErr),
			bear.WithExitCode(errors.TokenizerFailed),
			bear.WithTag("name", t.name),
		).Panic(true)
	}

	return "", false
}

// token creates a new token from the matcher
func (t *matcher) token(value string) Token {
	return Token{
		T:     t.t,
		Value: value,
	}
}

// TOK is a tokenizer for blowK source code
type TOK struct {
	Tokens []Token

	matchers []*matcher
	runes    []rune

	offset     int
	line       int
	lineOffset int
}

// Parse parses the given code and adds it to the string of TOK.tokens
func New(code string) *TOK {
	t := &TOK{
		matchers: matchers,
	}

	// load all the runes into the tokenizer
	for _, r := range code {
		t.runes = append(t.runes, r)
	}

	return t
}

// Parse parses through all the tokens in the source code
func (t *TOK) Parse() error {
	for t.next() {
		token, err := t.match()
		if err != nil {
			return bear.Wrap(err,
				bear.WithTag("line", t.line),
				bear.WithTag("line index", t.lineOffset),
			)
		}

		// update the offsets
		if token.T == NewLine {
			t.line++
			t.lineOffset = 0
		}

		t.offset += token.rlen()
		t.lineOffset += token.rlen()

		t.Tokens = append(t.Tokens, token)
	}

	return nil
}

// next returns true if there are still tokens to be parsed
func (t *TOK) next() bool {
	return t.offset < len(t.runes)
}

// match gets the next valid token in the string
// if the token can not be matched this function will panic
func (t *TOK) match() (Token, error) {
	rest := t.runes[t.offset:]
	for _, matcher := range t.matchers {
		if match, ok := matcher.match(string(rest)); ok {
			return matcher.token(match), nil
		}
	}

	return Token{}, bear.New(
		bear.WithExitCode(errors.TokenizerFailed),
		bear.WithErrType(UnknownTokenErr),
		bear.WithTag("token", string(rest)),
	)
}
