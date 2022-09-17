package pat

import (
	"github.com/bjatkin/blowk/internal/lex"
)

type Pattern interface {
	Match([]lex.Token) (int, bool)
}

type Exact struct {
	tok []lex.Token
}

func NewExact(types ...lex.TokType) *Exact {
	match := &Exact{}
	for _, t := range types {
		match.tok = append(match.tok, lex.Token{T: t})
	}

	return match
}

func (m *Exact) Match(check []lex.Token) (int, bool) {
	if len(check) == 0 {
		return 0, false
	}

	if len(check) < len(m.tok) {
		return 0, false
	}

	for i, m := range m.tok {
		if m.T != check[i].T {
			return 0, false
		}
	}

	return len(m.tok), true
}

type Not struct {
	t lex.TokType
}

func NewNot(t lex.TokType) *Not {
	match := &Not{
		t: t,
	}

	return match
}

func (m *Not) Match(check []lex.Token) (int, bool) {
	if len(check) < 1 {
		return 0, false
	}

	if m.t == check[0].T {
		return 0, false
	}

	return 1, true
}

type OneOf struct {
	tok []lex.Token
}

//NO MERGE should this take a slice of matchers instead?
func NewOneOf(types ...lex.TokType) *OneOf {
	match := &OneOf{}
	for _, t := range types {
		match.tok = append(match.tok, lex.Token{T: t})
	}

	return match
}

func (m *OneOf) Match(check []lex.Token) (int, bool) {
	if len(check) == 0 {
		return 0, false
	}

	for _, m := range m.tok {
		if check[0].T == m.T {
			return 1, true
		}
	}

	return 0, false
}

type OneOrMore struct {
	mat Pattern
}

func NewOneOrMore(mat Pattern) *OneOrMore {
	return &OneOrMore{
		mat: mat,
	}
}

func (m *OneOrMore) Match(check []lex.Token) (int, bool) {
	if len(check) == 0 {
		return 0, false
	}

	var done bool
	var matchCount int
	for !done {
		if len(check) < matchCount {
			break
		}

		i, ok := m.mat.Match(check[matchCount:])
		if ok {
			matchCount += i
			continue
		}
		done = true
	}

	return matchCount, matchCount > 0
}

type ZeroOrMore struct {
	mat Pattern
}

func NewZeroOrMore(mat Pattern) *ZeroOrMore {
	return &ZeroOrMore{
		mat: mat,
	}
}

func (m *ZeroOrMore) Match(check []lex.Token) (int, bool) {
	if len(check) == 0 {
		return 0, true
	}

	var done bool
	var matchCount int
	for !done {
		if len(check) < matchCount {
			break
		}

		i, ok := m.mat.Match(check[matchCount:])
		if ok {
			matchCount += i
			continue
		}
		done = true
	}

	return matchCount, true
}

type Composit struct {
	mats []Pattern
}

func NewComposit(mats ...Pattern) *Composit {
	return &Composit{
		mats: mats,
	}
}

func (m *Composit) Match(check []lex.Token) (int, bool) {
	if len(check) == 0 {
		return 0, false
	}

	var matchCount int
	for _, mat := range m.mats {
		if len(check) < matchCount {
			return 0, false
		}

		i, ok := mat.Match(check[matchCount:])
		if ok {
			matchCount += i
			continue
		}

		return 0, false
	}

	return matchCount, true
}

type Block struct {
	open  lex.TokType
	close lex.TokType
}

func NewBlock(open, close lex.TokType) *Block {
	return &Block{
		open:  open,
		close: close,
	}
}

func (m *Block) Match(check []lex.Token) (int, bool) {
	if len(check) == 0 {
		return 0, false
	}

	if check[0].T != m.open {
		return 0, false
	}

	openCount := 1
	var i int
	for openCount > 0 {
		i++

		if check[i].T == m.close {
			openCount--
		}

		// if open and close token are the same, don't look for nested blocks
		if m.open == m.close {
			continue
		}

		if check[i].T == m.open {
			openCount++
		}
	}

	// +1 becaseu we're returning the len not the offset
	return i + 1, true
}
