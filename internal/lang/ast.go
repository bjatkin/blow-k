package lang

import (
	"github.com/bjatkin/bear"
	"github.com/bjatkin/blow-k/internal/errors"
	"github.com/bjatkin/blow-k/internal/lex"
)

type Client struct {
	matchers []matcher
}

func NewClient() *Client {
	return &Client{
		matchers: []matcher{
			{
				match: MatchImport,
				new:   NewImport,
			},
		},
	}
}

func (c *Client) Build(tokens []lex.Token) (Node, error) {
	root := &Root{}

	exprs := c.getExpressions(tokens)
	for _, expr := range exprs {
		for _, matcher := range c.matchers {
			if matcher.match(expr) {
				node, err := matcher.new(expr)
				if err != nil {
					return nil, errors.Wrap(err,
						bear.WithErrType(errors.SyntaxError),
						bear.WithExitCode(errors.ASTFailed),
					)
				}

				switch v := node.(type) {
				case *Import:
					root.Imports = append(root.Imports, node)
				case *Var:
					if v.Name == "main" {
						root.Main = v
					} else {
						root.Expressions = append(root.Expressions, v)
					}
				default:
					return nil, errors.New(
						bear.WithErrType(errors.SyntaxError),
						bear.WithExitCode(errors.ASTFailed),
					)
				}
			}
		}
	}

	return root, nil
}

// getExpressions converts a slice of lex.Tokens into a [][]lex.Token
// where each sub-slic is an expression
func (c *Client) getExpressions(tokens []lex.Token) [][]lex.Token {
	var blocks [][]lex.Token
	var collect []lex.Token
	var openParenCount, openBraceCount, openSquareCount int

	for _, token := range tokens {
		collect = append(collect, token)

		if token.T == lex.OpenParen {
			openParenCount++
		}
		if token.T == lex.OpenBrace {
			openBraceCount++
		}
		if token.T == lex.OpenSquare {
			openSquareCount++
		}
		if token.T == lex.CloseParen {
			openParenCount--
		}
		if token.T == lex.CloseBrace {
			openBraceCount--
		}
		if token.T == lex.CloseSquare {
			openSquareCount--
		}

		if token.T == lex.SemiColon &&
			openParenCount == 0 &&
			openBraceCount == 0 &&
			openSquareCount == 0 {
			blocks = append(blocks, collect)
			collect = []lex.Token{}
		}
	}

	return blocks
}

// match, matches an express block and uses it to create a node
type matcher struct {
	match func([]lex.Token) bool
	new   func([]lex.Token) (Node, error)
}
