package markdown

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/yuin/goldmark/ast"
	tableast "github.com/yuin/goldmark/extension/ast"
)

func (m *Markdown) renderMarkdown() string {
	lines := collectDocumentLines(m.doc)
	return strings.Join(lines, lineFeed())
}

func collectDocumentLines(doc *ast.Document) []string {
	var lines []string
	for node := doc.FirstChild(); node != nil; node = node.NextSibling() {
		lines = append(lines, renderNodeLines(node, 0)...)
	}
	return lines
}

func renderNodeLines(node ast.Node, indentLevel int) []string {
	switch n := node.(type) {
	case *ast.Heading:
		return []string{renderHeadingLine(n)}
	case *ast.Paragraph:
		return []string{collectInlineText(n)}
	case *ast.Blockquote:
		return renderBlockquoteLines(n)
	case *ast.List:
		return renderListLines(n, indentLevel)
	case *ast.ThematicBreak:
		return []string{"---"}
	case *literalBlock:
		return []string{n.value}
	case *codeBlockNode:
		return renderCodeBlockLines(n)
	case *tableast.Table:
		return renderTableLines(n)
	default:
		return nil
	}
}

func renderHeadingLine(h *ast.Heading) string {
	prefix := strings.Repeat("#", h.Level)
	content := collectInlineText(h)
	if content == "" {
		return prefix
	}
	return fmt.Sprintf("%s %s", prefix, content)
}

func collectInlineText(node ast.Node) string {
	var buf strings.Builder
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		switch c := child.(type) {
		case *ast.String:
			buf.Write(c.Value)
		}
	}
	return buf.String()
}

func renderBlockquoteLines(bq *ast.Blockquote) []string {
	var lines []string
	for child := bq.FirstChild(); child != nil; child = child.NextSibling() {
		childLines := renderNodeLines(child, 0)
		if len(childLines) == 0 {
			lines = append(lines, ">")
			continue
		}
		for _, line := range childLines {
			if line == "" {
				lines = append(lines, ">")
				continue
			}
			lines = append(lines, "> "+line)
		}
	}
	return lines
}

func renderListLines(list *ast.List, indentLevel int) []string {
	var lines []string
	ordered := list.IsOrdered()
	counter := list.Start
	if !ordered {
		counter = 1
	} else if counter == 0 {
		counter = 1
	}
	for item := list.FirstChild(); item != nil; item = item.NextSibling() {
		li, ok := item.(*ast.ListItem)
		if !ok {
			continue
		}
		primary := ""
		var nested []ast.Node
		for child := li.FirstChild(); child != nil; child = child.NextSibling() {
			switch c := child.(type) {
			case *ast.Paragraph:
				if primary == "" {
					primary = collectInlineText(c)
				} else {
					primary += lineFeed() + collectInlineText(c)
				}
			case *ast.List:
				nested = append(nested, c)
			case *literalBlock:
				if primary == "" {
					primary = c.value
				} else {
					primary += lineFeed() + c.value
				}
			}
		}
		if primary == "" {
			primary = ""
		}
		indent := strings.Repeat("  ", indentLevel)
		if ordered {
			lines = append(lines, fmt.Sprintf("%s%d. %s", indent, counter, primary))
			counter++
		} else {
			lines = append(lines, fmt.Sprintf("%s- %s", indent, primary))
		}
		for _, nestedList := range nested {
			lines = append(lines, renderListLines(nestedList.(*ast.List), indentLevel+1)...)
		}
	}
	return lines
}

func renderCodeBlockLines(cb *codeBlockNode) []string {
	lf := lineFeed()
	var buf strings.Builder
	buf.WriteString("```")
	buf.WriteString(string(cb.language))
	buf.WriteString(lf)
	buf.WriteString(cb.value)
	buf.WriteString(lf)
	buf.WriteString("```")
	return []string{buf.String()}
}

func renderTableLines(table *tableast.Table) []string {
	var headerCells []string
	var header *tableast.TableHeader
	if h, ok := table.FirstChild().(*tableast.TableHeader); ok {
		header = h
	}
	if header != nil {
		for cell := header.FirstChild(); cell != nil; cell = cell.NextSibling() {
			c, ok := cell.(*tableast.TableCell)
			if !ok {
				continue
			}
			headerCells = append(headerCells, collectCellText(c))
		}
	}

	alignments := table.Alignments
	if len(alignments) == 0 {
		alignments = make([]tableast.Alignment, len(headerCells))
	}

	lf := lineFeed()
	var buf strings.Builder

	if len(headerCells) > 0 {
		buf.WriteString("|")
		for _, cell := range headerCells {
			buf.WriteString(" ")
			buf.WriteString(cell)
			buf.WriteString(" |")
		}
		buf.WriteString(lf)

		buf.WriteString("|")
		for i := 0; i < len(headerCells); i++ {
			align := tableast.AlignNone
			if i < len(alignments) {
				align = alignments[i]
			}
			buf.WriteString(separatorForAlignment(align))
		}
		buf.WriteString(lf)
	}

	for rowNode := table.FirstChild(); rowNode != nil; rowNode = rowNode.NextSibling() {
		row, ok := rowNode.(*tableast.TableRow)
		if !ok {
			continue
		}
		buf.WriteString("|")
		for _, cellText := range collectRowTexts(row) {
			buf.WriteString(" ")
			buf.WriteString(cellText)
			buf.WriteString(" |")
		}
		buf.WriteString(lf)
	}

	return []string{buf.String()}
}

func collectCellText(cell *tableast.TableCell) string {
	var buf strings.Builder
	for child := cell.FirstChild(); child != nil; child = child.NextSibling() {
		switch c := child.(type) {
		case *ast.Paragraph:
			buf.WriteString(collectInlineText(c))
		case *literalBlock:
			buf.WriteString(c.value)
		case *ast.String:
			buf.Write(c.Value)
		}
	}
	return buf.String()
}

func collectRowTexts(row *tableast.TableRow) []string {
	var cells []string
	for cell := row.FirstChild(); cell != nil; cell = cell.NextSibling() {
		if c, ok := cell.(*tableast.TableCell); ok {
			cells = append(cells, collectCellText(c))
		}
	}
	return cells
}

func separatorForAlignment(a tableast.Alignment) string {
	switch a {
	case tableast.AlignLeft:
		return ":--------|"
	case tableast.AlignCenter:
		return ":-------:|"
	case tableast.AlignRight:
		return "--------:|"
	default:
		return "---------|"
	}
}

func lineFeed() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}
