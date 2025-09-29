package markdown

import (
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"
	tableast "github.com/yuin/goldmark/extension/ast"
)

// TableOfContents generates a table of contents from the recorded headers.
func (m *Markdown) TableOfContents(depth TableOfContentsDepth) *Markdown {
	if len(m.headers) == 0 {
		return m
	}

	for _, header := range m.headers {
		if header.level > depth {
			continue
		}
		indent := strings.Repeat("  ", int(header.level)-1)
		anchor := buildAnchor(header.text)
		line := fmt.Sprintf("%s- [%s](#%s)", indent, header.text, anchor)
		m.appendBlock(newLiteralBlock(line))
	}
	m.appendBlock(newLiteralBlock(""))
	return m
}

func buildAnchor(text string) string {
	anchor := strings.ToLower(text)
	anchor = strings.ReplaceAll(anchor, " ", "-")
	anchor = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, anchor)
	return anchor
}

// Details renders an HTML <details> block.
func (m *Markdown) Details(summary, text string) *Markdown {
	block := fmt.Sprintf("<details><summary>%s</summary>%s%s%s</details>", summary, lineFeed(), text, lineFeed())
	m.appendBlock(newLiteralBlock(block))
	return m
}

// Detailsf renders formatted <details> block content.
func (m *Markdown) Detailsf(summary, format string, args ...interface{}) *Markdown {
	return m.Details(summary, fmt.Sprintf(format, args...))
}

// BulletList appends an unordered list.
func (m *Markdown) BulletList(items ...string) *Markdown {
	if len(items) == 0 {
		return m
	}
	list := ast.NewList('-')
	list.IsTight = true
	for _, item := range items {
		listItem := ast.NewListItem(0)
		paragraph := ast.NewParagraph()
		paragraph.AppendChild(paragraph, ast.NewString([]byte(item)))
		listItem.AppendChild(listItem, paragraph)
		list.AppendChild(list, listItem)
	}
	m.appendBlock(list)
	return m
}

// OrderedList appends an ordered list.
func (m *Markdown) OrderedList(items ...string) *Markdown {
	if len(items) == 0 {
		return m
	}
	list := ast.NewList('.')
	list.IsTight = true
	list.Start = 1
	for _, item := range items {
		listItem := ast.NewListItem(0)
		paragraph := ast.NewParagraph()
		paragraph.AppendChild(paragraph, ast.NewString([]byte(item)))
		listItem.AppendChild(listItem, paragraph)
		list.AppendChild(list, listItem)
	}
	m.appendBlock(list)
	return m
}

// CheckBox appends a checkbox list.
func (m *Markdown) CheckBox(set []CheckBoxSet) *Markdown {
	if len(set) == 0 {
		return m
	}
	list := ast.NewList('-')
	list.IsTight = true
	for _, entry := range set {
		prefix := "[ ] "
		if entry.Checked {
			prefix = "[x] "
		}
		paragraph := ast.NewParagraph()
		paragraph.AppendChild(paragraph, ast.NewString([]byte(prefix+entry.Text)))
		item := ast.NewListItem(0)
		item.AppendChild(item, paragraph)
		list.AppendChild(list, item)
	}
	m.appendBlock(list)
	return m
}

// Blockquote appends a blockquote block.
func (m *Markdown) Blockquote(text string) *Markdown {
	block := ast.NewBlockquote()
	normalized := strings.ReplaceAll(text, "\r\n", "\n")
	lines := strings.Split(normalized, "\n")
	for _, line := range lines {
		paragraph := ast.NewParagraph()
		paragraph.AppendChild(paragraph, ast.NewString([]byte(line)))
		block.AppendChild(block, paragraph)
	}
	m.appendBlock(block)
	return m
}

// CodeBlocks appends a fenced code block.
func (m *Markdown) CodeBlocks(lang SyntaxHighlight, text string) *Markdown {
	m.appendBlock(newCodeBlockNode(lang, text))
	return m
}

// HorizontalRule appends a thematic break.
func (m *Markdown) HorizontalRule() *Markdown {
	m.appendBlock(ast.NewThematicBreak())
	return m
}

// Table renders a markdown table using goldmark table AST nodes.
func (m *Markdown) Table(set TableSet) *Markdown {
	if err := set.ValidateColumns(); err != nil {
		if m.err != nil {
			m.err = fmt.Errorf("failed to validate columns: %w: %s", err, m.err)
		} else {
			m.err = fmt.Errorf("failed to validate columns: %w", err)
		}
		return m
	}

	if len(set.Header) == 0 {
		return m
	}

	table := tableast.NewTable()
	table.Alignments = convertAlignments(set)

	headerRow := tableast.NewTableRow(table.Alignments)
	for idx, cellText := range set.Header {
		cell := tableast.NewTableCell()
		if idx < len(table.Alignments) {
			cell.Alignment = table.Alignments[idx]
		}
		paragraph := ast.NewParagraph()
		paragraph.AppendChild(paragraph, ast.NewString([]byte(cellText)))
		cell.AppendChild(cell, paragraph)
		headerRow.AppendChild(headerRow, cell)
	}
	table.AppendChild(table, tableast.NewTableHeader(headerRow))

	for _, row := range set.Rows {
		rowNode := tableast.NewTableRow(table.Alignments)
		for idx, cellText := range row {
			cell := tableast.NewTableCell()
			if idx < len(table.Alignments) {
				cell.Alignment = table.Alignments[idx]
			}
			paragraph := ast.NewParagraph()
			paragraph.AppendChild(paragraph, ast.NewString([]byte(cellText)))
			cell.AppendChild(cell, paragraph)
			rowNode.AppendChild(rowNode, cell)
		}
		table.AppendChild(table, rowNode)
	}

	m.appendBlock(table)
	return m
}

func convertAlignments(set TableSet) []tableast.Alignment {
	aligned := make([]tableast.Alignment, len(set.Header))
	for i := 0; i < len(set.Header); i++ {
		align := AlignDefault
		if i < len(set.Alignment) {
			align = set.Alignment[i]
		}
		switch align {
		case AlignLeft:
			aligned[i] = tableast.AlignLeft
		case AlignCenter:
			aligned[i] = tableast.AlignCenter
		case AlignRight:
			aligned[i] = tableast.AlignRight
		default:
			aligned[i] = tableast.AlignNone
		}
	}
	return aligned
}

// CustomTable renders a table with optional formatting behaviors.
func (m *Markdown) CustomTable(set TableSet, options TableOptions) *Markdown {
	if options.AutoFormatHeaders {
		set.Header = formatHeaders(set.Header)
	}
	return m.Table(set)
}

func formatHeaders(headers []string) []string {
	formatted := make([]string, len(headers))
	for i, header := range headers {
		words := strings.Fields(strings.ToLower(header))
		for j, word := range words {
			if len(word) == 0 {
				continue
			}
			runes := []rune(word)
			runes[0] = toUpper(runes[0])
			words[j] = string(runes)
		}
		formatted[i] = strings.Join(words, " ")
	}
	return formatted
}

func toUpper(r rune) rune {
	if r >= 'a' && r <= 'z' {
		return r - ('a' - 'A')
	}
	return r
}

// LF appends a markdown line feed (two spaces).
func (m *Markdown) LF() *Markdown {
	m.appendBlock(newLiteralBlock("  "))
	return m
}
