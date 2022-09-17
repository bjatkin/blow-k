package ast

import (
	"encoding/json"
	"fmt"

	"github.com/bjatkin/bear"
	"github.com/bjatkin/blowk/internal/lex"
)

type Client struct {
	matchers []matcher
}

func NewTopLevelClient() *Client {
	return &Client{
		matchers: []matcher{
			importMatcher,
			commentNodeMatcher,
			fnMatcher,
		},
	}
}

func newFnClient() *Client {
	return &Client{
		matchers: []matcher{
			cmdExecMatcher,
			commentNodeMatcher,
		},
	}
}

func (c *Client) Parse(tokens []lex.Token) Node {
	root := &RootNode{}
	var i int
	for i < len(tokens) {
		for _, m := range c.matchers {
			if offset, ok := m.pat.Match(tokens[i:]); ok {
				node := m.new(tokens[i : i+offset])
				root.children = append(root.children, node)

				// offset - 1 because every loop adds 1 to i
				i += (offset - 1)
				break
			}
		}
		i++
	}

	data, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		fmt.Println("failed to marshal AST")
	}

	imports := getImports(root)
	err = checkCmds(root, imports)
	if err != nil {
		fmt.Println("failed to process ast")
	}

	fmt.Println("AST: ", string(data))
	return root
}

func getImports(ast Node) []*ImportNode {
	var nodes []*ImportNode
	if n, ok := ast.(*ImportNode); ok {
		nodes = append(nodes, n)
	}

	for _, child := range ast.Children() {
		nodes = append(nodes, getImports(child)...)
	}

	return nodes
}

func checkCmds(ast Node, imports []*ImportNode) error {
	if n, ok := ast.(*CmdExpressionNode); ok {
		// check to make sure this is a valid cmd
		for _, i := range imports {
			if n.cmd == i.Ident {
				goto success
			}
			if n.cmd == i.Alias {
				n.cmd = i.Ident
				goto success
			}
		}

		return bear.New(bear.WithMsg("cmd does not exist"), bear.WithTag("cmd", n.cmd))
	success:
	}

	for _, child := range ast.Children() {
		err := checkCmds(child, imports)
		if err != nil {
			return err
		}
	}

	return nil
}
