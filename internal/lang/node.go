package lang

import (
	"github.com/bjatkin/bear"
	"github.com/bjatkin/blow-k/internal/errors"
	"github.com/bjatkin/blow-k/internal/lex"
)

type Node interface {
	Children() []Node
}

type Root struct {
	Imports     []Node
	Main        *Var
	Expressions []Node
}

func (n *Root) Children() []Node {
	var children []Node
	children = append(children, n.Imports...)
	children = append(children, n.Main)
	children = append(children, n.Expressions...)

	return children
}

type Import struct {
	Name string
	As   string
	From string
}

func (n *Import) Children() []Node {
	return nil
}

func MatchImport(tokens []lex.Token) bool {
	return len(tokens) > 0 && tokens[0].T == lex.ImportKeyword
}

func NewImport(tokens []lex.Token) (Node, error) {
	switch len(tokens) {
	case 3:
		if tokens[0].T == lex.ImportKeyword &&
			tokens[1].T == lex.Identifyer &&
			tokens[2].T == lex.SemiColon {
			return &Import{
				Name: tokens[1].Value,
			}, nil
		}
	case 5:
		if tokens[0].T == lex.ImportKeyword &&
			tokens[1].T == lex.Identifyer &&
			tokens[2].T == lex.AsKeyword &&
			tokens[3].T == lex.Identifyer &&
			tokens[4].T == lex.SemiColon {
			return &Import{
				Name: tokens[1].Value,
				As:   tokens[3].Value,
			}, nil
		}

		if tokens[0].T == lex.ImportKeyword &&
			tokens[1].T == lex.Identifyer &&
			tokens[2].T == lex.FromKeyword &&
			tokens[3].T == lex.String &&
			tokens[4].T == lex.SemiColon {
			return &Import{
				Name: tokens[1].Value,
				From: tokens[3].Value,
			}, nil
		}
	case 7:
		if tokens[0].T == lex.ImportKeyword &&
			tokens[1].T == lex.Identifyer &&
			tokens[2].T == lex.AsKeyword &&
			tokens[3].T == lex.Identifyer &&
			tokens[4].T == lex.FromKeyword &&
			tokens[5].T == lex.String &&
			tokens[6].T == lex.SemiColon {
			return &Import{
				Name: tokens[1].Value,
				As:   tokens[3].Value,
				From: tokens[5].Value,
			}, nil
		}
	}

	return nil, errors.New(
		bear.WithErrType(errors.SyntaxError),
		bear.WithLabels("invalid import"),
	)
}

type Var struct {
	Name    string
	Type    Node
	Default Node
}

func MatchVar(tokens []lex.Token) bool {
	return false
}

func NewVar(tokens []lex.Token) (Node, error) {
	return nil, nil
}

func (n *Var) Children() []Node {
	return []Node{n.Type, n.Default}
}

type Exec struct {
	Expr Node
}

type Comment struct {
	Value string
}

type String struct {
	Value Node
}

type StringArray struct {
	Value Node
}
