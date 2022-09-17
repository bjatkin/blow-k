package tok

import (
	"os"

	"github.com/bjatkin/bear"
	"github.com/bjatkin/blow-k/internal/errors"
)

// Token is a tokenizer token
type Token struct {
	Value      string
	FileName   string
	LineNumber int
	ColNumber  int
}

// Client is a tokenizer client that tokenizes a file
type Client struct {
	sperators []rune
}

// NewClient creates a new default Tok.Client
func NewClient() *Client {
	return &Client{
		sperators: []rune{
			' ', '\t', '\n', // white space tokens
			'(', ')', '{', '}', '[', ']', // parens etc. tokens
			':', '.', ',', '$', '"', // punctuation tokens
			'+', '-', '*', '/', '&', '|', '=', // math tokens
		},
	}
}

// Parse reads in a file and converts it into a token slice
func (c *Client) Parse(fileName string) ([]Token, error) {
	src, err := os.ReadFile(fileName)
	if err != nil {
		return nil, bear.Wrap(err,
			bear.WithErrType(errors.FileNotFound),
		)
	}

	var collect []rune
	var tokens []Token
	var lineNumber, colNumber int
	for _, r := range string(src) {
		if !c.split(r) {
			collect = append(collect, r)
		} else {
			tokens = append(tokens, Token{
				Value:      string(collect),
				FileName:   fileName,
				LineNumber: lineNumber,
				ColNumber:  colNumber - len(collect),
			})
			tokens = append(tokens, Token{
				Value:      string(r),
				FileName:   fileName,
				LineNumber: lineNumber,
				ColNumber:  colNumber,
			})

			collect = []rune{}
		}

		colNumber++
		if r == '\n' {
			lineNumber++
			colNumber = 0
		}
	}

	// make sure to get the final token
	tokens = append(tokens, Token{
		Value:      string(collect),
		FileName:   fileName,
		LineNumber: lineNumber,
		ColNumber:  colNumber - len(collect),
	})

	// filter out any empty tokens that were added
	var filtered []Token
	for _, tok := range tokens {
		if len(tok.Value) > 0 {
			filtered = append(filtered, tok)
		}
	}

	return filtered, nil
}

// split returns true if the check rune should split 2 tokens apart
func (c *Client) split(check rune) bool {
	for _, r := range c.sperators {
		if r == check {
			return true
		}
	}

	return false
}
