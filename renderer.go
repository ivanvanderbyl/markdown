package markdown

import (
	"fmt"
	"runtime"
	"strings"
	"unicode/utf8"

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

	bodyRows := make([][]string, 0)
	for rowNode := table.FirstChild(); rowNode != nil; rowNode = rowNode.NextSibling() {
		row, ok := rowNode.(*tableast.TableRow)
		if !ok {
			continue
		}
		bodyRows = append(bodyRows, collectRowTexts(row))
	}

	widths := computeColumnWidths(headerCells, bodyRows)
	alignments := normalizeAlignments(table.Alignments, len(widths))

	lf := lineFeed()
	var buf strings.Builder

	if len(headerCells) > 0 {
		buf.WriteString("|")
		for i, cell := range headerCells {
			buf.WriteString(" ")
			buf.WriteString(padCell(cell, widths[i], alignments[i]))
			buf.WriteString(" |")
		}
		buf.WriteString(lf)

		buf.WriteString("|")
		for i := range widths {
			buf.WriteString(" ")
			buf.WriteString(alignmentSegment(alignments[i], widths[i]))
			buf.WriteString(" |")
		}
		buf.WriteString(lf)
	}

	for _, row := range bodyRows {
		buf.WriteString("|")
		for i, cellText := range row {
			buf.WriteString(" ")
			buf.WriteString(padCell(cellText, widths[i], alignments[i]))
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

func computeColumnWidths(header []string, rows [][]string) []int {
	widths := make([]int, len(header))
	for i, cell := range header {
		widths[i] = runeWidth(cell)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i >= len(widths) {
				continue
			}
			if w := runeWidth(cell); w > widths[i] {
				widths[i] = w
			}
		}
	}
	return widths
}

func lineFeed() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}

func runeWidth(value string) int {
	return utf8.RuneCountInString(value)
}

func padCell(value string, width int, align tableast.Alignment) string {
	padding := width - runeWidth(value)
	if padding <= 0 {
		return value
	}
	switch align {
	case tableast.AlignRight:
		return strings.Repeat(" ", padding) + value
	case tableast.AlignCenter:
		left := padding / 2
		right := padding - left
		return strings.Repeat(" ", left) + value + strings.Repeat(" ", right)
	default:
		return value + strings.Repeat(" ", padding)
	}
}

func normalizeAlignments(align []tableast.Alignment, count int) []tableast.Alignment {
	if len(align) >= count {
		return align[:count]
	}
	result := make([]tableast.Alignment, count)
	copy(result, align)
	return result
}

func alignmentSegment(a tableast.Alignment, width int) string {
	segmentWidth := width
	if segmentWidth < 3 {
		segmentWidth = 3
	}
	switch a {
	case tableast.AlignLeft:
		return ":" + strings.Repeat("-", segmentWidth-1)
	case tableast.AlignCenter:
		return ":" + strings.Repeat("-", segmentWidth-2) + ":"
	case tableast.AlignRight:
		return strings.Repeat("-", segmentWidth-1) + ":"
	default:
		return strings.Repeat("-", segmentWidth)
	}
}
