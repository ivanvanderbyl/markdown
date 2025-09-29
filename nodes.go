package markdown

import "github.com/yuin/goldmark/ast"

var (
	kindLiteralBlock = ast.NewNodeKind("MarkdownLiteralBlock")
	kindCodeBlock    = ast.NewNodeKind("MarkdownCodeBlock")
)

type literalBlock struct {
	ast.BaseBlock
	value string
}

func newLiteralBlock(value string) *literalBlock {
	return &literalBlock{value: value}
}

func (n *literalBlock) Kind() ast.NodeKind {
	return kindLiteralBlock
}

func (n *literalBlock) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, map[string]string{"Value": n.value}, nil)
}

type codeBlockNode struct {
	ast.BaseBlock
	language SyntaxHighlight
	value    string
}

func newCodeBlockNode(language SyntaxHighlight, value string) *codeBlockNode {
	return &codeBlockNode{language: language, value: value}
}

func (n *codeBlockNode) Kind() ast.NodeKind {
	return kindCodeBlock
}

func (n *codeBlockNode) Dump(source []byte, level int) {
	meta := map[string]string{}
	if n.language != "" {
		meta["Language"] = string(n.language)
	}
	ast.DumpHelper(n, source, level, meta, nil)
}
