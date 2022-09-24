package lex

import (
	"github.com/bjatkin/bear"
	"github.com/bjatkin/blow-k/internal/errors"
)

// transformer is a function that transforms the lex.Token slice into a different slice
type transformer func([]Token) []Token

// coalesceStrings combines tokens into string literals
func coalesceStrings(tokens []Token) []Token {
	var ret []Token
	var collect []Token
	var open bool

	for _, tok := range tokens {
		if tok.T == StartEndString && open {
			ret = append(ret, combineTokens(String, collect))
			open = false
			collect = []Token{}
			continue
		}

		if tok.T == StartEndString && !open {
			open = true
			continue
		}

		if open {
			collect = append(collect, tok)
			continue
		}

		ret = append(ret, tok)
	}

	return ret
}

// coalesceComments combines tokens in to comments
func coalesceComments(tokens []Token) []Token {
	var ret []Token
	var collect []Token
	var open bool

	for _, tok := range tokens {
		if tok.T == NewLine && open {
			ret = append(ret, combineTokens(Comment, collect))
			collect = []Token{}
			open = false
		}

		if tok.T == StartComment {
			open = true
			continue
		}

		if open {
			collect = append(collect, tok)
			continue
		}

		ret = append(ret, tok)
	}

	return ret
}

// coalesceArrayTypes combines tokens into a single array type
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
			ret = append(ret, combineTokens(StringArrayType, tokens[i:i+3]))
			i += 2
			continue
		}

		ret = append(ret, tok)
	}

	return ret
}

// combineTokens combines tok.Tokens into a single Token
func combineTokens(t TokType, tokens []Token) Token {
	if len(tokens) == 0 {
		errors.New(
			bear.WithExitCode(errors.LexerFailed),
			bear.WithTag("tokens", 0),
		).Panic(true)
	}

	token := Token{
		T:          t,
		FileName:   tokens[0].FileName,
		ColNumber:  tokens[0].ColNumber,
		LineNumber: tokens[0].LineNumber,
	}

	for _, t := range tokens {
		token.Value += t.Value
	}

	return token
}

// insertSemicolons add semi-colons to all eligible new lines
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

		tokens[i] = Token{T: SemiColon, Value: ";", FileName: tok.FileName, LineNumber: tok.LineNumber, ColNumber: tok.ColNumber}
	}

	return tokens
}

// filter removes whitespace and new line tokens
func filter(tokens []Token) []Token {
	var ret []Token

	for _, tok := range tokens {
		if tok.T != WhiteSpace &&
			tok.T != NewLine {
			ret = append(ret, tok)
		}

	}

	return ret
}
