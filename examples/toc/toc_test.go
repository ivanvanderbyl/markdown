package toc_test

import (
	"bytes"
	"testing"

	"github.com/ivanvanderbyl/markdown"
)

func TestTOCRendering(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	md := markdown.NewMarkdown(buf)
	md.H1("Guide to markdown").
		PlainText("Markdown built through goldmark AST.").
		Table(markdown.TableSet{
			Header: []string{"Feature", "Description"},
			Rows: [][]string{
				{"TOC", "Generate nested table of contents"},
				{"Tables", "Alignment-aware rendering without tablewriter"},
			},
		}).
		Build()

	expected := `# Guide to markdown
Markdown built through goldmark AST.
| Feature | Description                                   |
| ------- | --------------------------------------------- |
| TOC     | Generate nested table of contents             |
| Tables  | Alignment-aware rendering without tablewriter |
`

	if buf.String() != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, buf.String())
	}
}
