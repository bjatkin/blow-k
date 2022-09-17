package ast

import (
	"github.com/bjatkin/blowk/internal/lex"
	"github.com/bjatkin/blowk/internal/pat"
)

type matcher struct {
	pat pat.Pattern
	new NewNode
}

var (
	importMatcher = matcher{
		pat: pat.NewComposit(
			pat.NewExact(lex.ImportKeyword, lex.Identifyer),
			pat.NewZeroOrMore(pat.NewExact(lex.AsKeyword, lex.Identifyer)),
			pat.NewExact(lex.SemiColon),
		),
		new: NewImportNode,
	}
	fnMatcher = matcher{
		pat: pat.NewComposit(
			pat.NewExact(lex.Identifyer, lex.Colon, lex.OpenParen),
			pat.NewZeroOrMore(pat.NewNot(lex.CloseParen)),
			pat.NewExact(lex.CloseParen, lex.Colon),
			pat.NewBlock(lex.OpenBrace, lex.CloseBrace),
		),
		new: NewFnNode,
	}
	cmdExecMatcher = matcher{
		pat: pat.NewComposit(
			pat.NewExact(lex.Exec, lex.Identifyer, lex.OpenSquare),
			pat.NewZeroOrMore(pat.NewNot(lex.CloseSquare)),
			pat.NewExact(lex.CloseSquare),
		),
		new: NewCmdExecNode,
	}
	commentNodeMatcher = matcher{
		pat: pat.NewExact(lex.Comment),
		new: NewCommentNode,
	}
)
