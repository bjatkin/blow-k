package ast

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bjatkin/bear"
	"github.com/bjatkin/blowk/internal/errors"
	"github.com/bjatkin/blowk/internal/lex"
	"github.com/bjatkin/blowk/internal/pat"
)

type Node interface {
	Children() []Node
	MarshalJSON() ([]byte, error)
	BashGen() (string, error)
}

type NewNode func([]lex.Token) Node

type RootNode struct {
	children []Node
}

func (n *RootNode) Children() []Node {
	return n.children
}

func (n *RootNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type     string
		Children []Node
	}{
		"root",
		n.children,
	})
}

func (n *RootNode) BashGen() (string, error) {
	var other []string
	var imports []string

	for _, c := range n.children {
		gen, err := c.BashGen()
		if err != nil {
			return "", bear.Wrap(
				err,
				bear.WithErrType(errors.CodeGenFailed),
				bear.WithExitCode(errors.BashGenFailed),
			)
		}

		if _, ok := c.(*ImportNode); ok {
			imports = append(imports, gen)
		} else {
			other = append(other, gen)
		}

	}
	return fmt.Sprintf(`# check for which and echo by default
if [[ -z "$( which which )" ]]; then
	exit 213
fi

if [[ -z "$( which echo )" ]]; then
	exit 214
fi

%s

%s`, strings.Join(imports, "\n"), strings.Join(other, "\n")), nil
}

type ImportNode struct {
	Ident string
	Alias string
}

func NewImportNode(toks []lex.Token) Node {
	node := &ImportNode{
		Ident: toks[1].V,
	}

	if len(toks) > 3 {
		node.Alias = toks[3].V
	}

	return node
}

func (i *ImportNode) Children() []Node {
	return nil
}

func (n *ImportNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type   string
		Import string
	}{
		"import",
		n.Ident,
	})
}

func (n *ImportNode) BashGen() (string, error) {
	if n.Ident == "echo" || n.Ident == "which" {
		// echo and which are both imported by default
		return "", nil
	}

	return fmt.Sprintf(`if [[ -z "$( which %s )" ]]; then
	echo "imported command %s could not be found"
	exit 215
fi`, n.Ident, n.Ident), nil
}

type FnNode struct {
	name      string
	args      []Node
	ret       lex.TokType
	statments []Node
}

func NewFnNode(toks []lex.Token) Node {
	node := &FnNode{
		name: toks[0].V,
	}

	i := 3
	argsMatch := pat.NewComposit(
		pat.NewExact(lex.Identifyer),
		pat.NewOneOf(lex.StringType, lex.StringArrayType),
	)
	for toks[i].T != lex.CloseParen {
		if _, ok := argsMatch.Match(toks[i:]); ok {
			node.args = append(node.args, NewVarNode(toks[i:]))
		}
		i++
	}

	node.ret = toks[i+1].T

	j, _ := pat.NewBlock(lex.OpenBrace, lex.CloseBrace).Match(toks[i+2:])

	root := newFnClient().Parse(toks[i+2 : i+2+j])
	node.statments = root.Children()

	return node
}

func (i *FnNode) Children() []Node {
	return i.statments
}

func (n *FnNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type      string
		FnName    string
		Args      []Node
		Return    int
		Statments []Node
	}{
		"function",
		n.name,
		n.args,
		int(n.ret),
		n.statments,
	})
}

func (n *FnNode) BashGen() (string, error) {
	var statments []string
	for _, s := range n.statments {
		gen, err := s.BashGen()
		if err != nil {
			return "", bear.Wrap(err, bear.WithErrType(errors.CodeGenFailed))
		}
		statments = append(statments, gen)
	}
	if n.name == "main" {
		return fmt.Sprintf("# main code\n%s", strings.Join(statments, "\n")), nil
	}

	// TODO: need to add other functions
	return "", nil
}

type VarNode struct {
	Ident string
	Type  lex.TokType
	Value string
}

func NewVarNode(toks []lex.Token) Node {
	// get identifyer and type
	node := &VarNode{
		Ident: toks[0].V,
		Type:  toks[2].T,
	}

	return node
}

func (n *VarNode) Children() []Node {
	return nil
}

func (n *VarNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type         string
		VarName      string
		VarType      int
		InitialValue string
	}{
		"Var",
		n.Ident,
		int(n.Type),
		n.Value,
	})
}

func (n *VarNode) BashGen() (string, error) {
	return fmt.Sprintf(`%s="%s"`, n.Ident, n.Value), nil
}

type CmdExeNode struct {
	cmd Node
}

func NewCmdExecNode(toks []lex.Token) Node {
	return &CmdExeNode{
		cmd: NewCmdExpressionNode(toks[1:]),
	}
}

func (n *CmdExeNode) Children() []Node {
	return []Node{n.cmd}
}

func (n *CmdExeNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string
		Exec Node
	}{
		"cmd exec",
		n.cmd,
	})
}

func (n *CmdExeNode) BashGen() (string, error) {
	return n.cmd.BashGen()
}

type CmdExpressionNode struct {
	cmd  string
	args []string
}

func NewCmdExpressionNode(toks []lex.Token) Node {
	node := &CmdExpressionNode{cmd: toks[0].V}

	i := 2
	for toks[i].T != lex.CloseSquare {
		if toks[i].T != lex.WhiteSpace && toks[i].T != lex.Comma {
			node.args = append(node.args, toks[i].V)
		}
		i++
	}

	return node
}

func (n *CmdExpressionNode) Children() []Node {
	return nil
}

func (n *CmdExpressionNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string
		Cmd  string
		Args []string
	}{
		"cmd expression",
		n.cmd,
		n.args,
	})
}

func (n *CmdExpressionNode) BashGen() (string, error) {
	cmd := fmt.Sprintf("%s ", n.cmd)

	for _, arg := range n.args {
		cmd += fmt.Sprintf(`"%s" `, arg)
	}

	return cmd, nil
}

type CommentNode struct {
	comment string
}

func NewCommentNode(toks []lex.Token) Node {
	return &CommentNode{
		comment: toks[0].V,
	}
}

func (n *CommentNode) Children() []Node {
	return nil
}

func (n *CommentNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type    string
		Comment string
	}{
		"comment",
		n.comment,
	})
}

func (n *CommentNode) BashGen() (string, error) {
	return "", nil
}
