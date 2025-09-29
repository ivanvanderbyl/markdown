package markdown

import (
	"fmt"
	"io"

	"github.com/yuin/goldmark/ast"
)

// SyntaxHighlight is syntax highlight language.
type SyntaxHighlight string

const (
	SyntaxHighlightNone         SyntaxHighlight = ""
	SyntaxHighlightText         SyntaxHighlight = "text"
	SyntaxHighlightAPIBlueprint SyntaxHighlight = "markdown"
	SyntaxHighlightShell        SyntaxHighlight = "shell"
	SyntaxHighlightGo           SyntaxHighlight = "go"
	SyntaxHighlightJSON         SyntaxHighlight = "json"
	SyntaxHighlightYAML         SyntaxHighlight = "yaml"
	SyntaxHighlightXML          SyntaxHighlight = "xml"
	SyntaxHighlightHTML         SyntaxHighlight = "html"
	SyntaxHighlightCSS          SyntaxHighlight = "css"
	SyntaxHighlightJavaScript   SyntaxHighlight = "javascript"
	SyntaxHighlightTypeScript   SyntaxHighlight = "typescript"
	SyntaxHighlightSQL          SyntaxHighlight = "sql"
	SyntaxHighlightC            SyntaxHighlight = "c"
	SyntaxHighlightCSharp       SyntaxHighlight = "csharp"
	SyntaxHighlightCPlusPlus    SyntaxHighlight = "cpp"
	SyntaxHighlightJava         SyntaxHighlight = "java"
	SyntaxHighlightKotlin       SyntaxHighlight = "kotlin"
	SyntaxHighlightPHP          SyntaxHighlight = "php"
	SyntaxHighlightPython       SyntaxHighlight = "python"
	SyntaxHighlightRuby         SyntaxHighlight = "ruby"
	SyntaxHighlightSwift        SyntaxHighlight = "swift"
	SyntaxHighlightScala        SyntaxHighlight = "scala"
	SyntaxHighlightRust         SyntaxHighlight = "rust"
	SyntaxHighlightObjectiveC   SyntaxHighlight = "objectivec"
	SyntaxHighlightPerl         SyntaxHighlight = "perl"
	SyntaxHighlightLua          SyntaxHighlight = "lua"
	SyntaxHighlightDart         SyntaxHighlight = "dart"
	SyntaxHighlightClojure      SyntaxHighlight = "clojure"
	SyntaxHighlightGroovy       SyntaxHighlight = "groovy"
	SyntaxHighlightR            SyntaxHighlight = "r"
	SyntaxHighlightHaskell      SyntaxHighlight = "haskell"
	SyntaxHighlightErlang       SyntaxHighlight = "erlang"
	SyntaxHighlightElixir       SyntaxHighlight = "elixir"
	SyntaxHighlightOCaml        SyntaxHighlight = "ocaml"
	SyntaxHighlightJulia        SyntaxHighlight = "julia"
	SyntaxHighlightScheme       SyntaxHighlight = "scheme"
	SyntaxHighlightFSharp       SyntaxHighlight = "fsharp"
	SyntaxHighlightCoffeeScript SyntaxHighlight = "coffeescript"
	SyntaxHighlightVBNet        SyntaxHighlight = "vbnet"
	SyntaxHighlightTeX          SyntaxHighlight = "tex"
	SyntaxHighlightDiff         SyntaxHighlight = "diff"
	SyntaxHighlightApache       SyntaxHighlight = "apache"
	SyntaxHighlightDockerfile   SyntaxHighlight = "dockerfile"
	SyntaxHighlightMermaid      SyntaxHighlight = "mermaid"
)

// TableOfContentsDepth represents the depth level for table of contents.
type TableOfContentsDepth int

const (
	TableOfContentsDepthH1 TableOfContentsDepth = 1
	TableOfContentsDepthH2 TableOfContentsDepth = 2
	TableOfContentsDepthH3 TableOfContentsDepth = 3
	TableOfContentsDepthH4 TableOfContentsDepth = 4
	TableOfContentsDepthH5 TableOfContentsDepth = 5
	TableOfContentsDepthH6 TableOfContentsDepth = 6
)

type headerInfo struct {
	level TableOfContentsDepth
	text  string
}

// TableAlignment represents column alignment in markdown tables.
type TableAlignment int

const (
	// AlignDefault represents no specific alignment (left by default).
	AlignDefault TableAlignment = iota
	// AlignLeft represents left alignment (:------).
	AlignLeft
	// AlignCenter represents center alignment (:-----:).
	AlignCenter
	// AlignRight represents right alignment (------:).
	AlignRight
)

// TableSet describes the content and layout for a markdown table.
type TableSet struct {
	Header    []string
	Rows      [][]string
	Alignment []TableAlignment
}

// ValidateColumns checks if the number of columns in the header and records match.
func (t *TableSet) ValidateColumns() error {
	headerColumns := len(t.Header)
	for _, record := range t.Rows {
		if len(record) != headerColumns {
			return ErrMismatchColumn
		}
	}
	return nil
}

// TableOptions controls formatting when rendering custom tables.
type TableOptions struct {
	AutoWrapText      bool
	AutoFormatHeaders bool
}

// CheckBoxSet configures a single checkbox entry.
type CheckBoxSet struct {
	Checked bool
	Text    string
}

// Markdown is markdown text.
type Markdown struct {
	doc     *ast.Document
	dest    io.Writer
	err     error
	headers []headerInfo
}

func (m *Markdown) appendBlock(node ast.Node) {
	m.doc.AppendChild(m.doc, node)
}

// NewMarkdown returns new Markdown.
func NewMarkdown(w io.Writer) *Markdown {
	return &Markdown{
		doc:     ast.NewDocument(),
		dest:    w,
		headers: []headerInfo{},
	}
}

// String returns markdown text.
func (m *Markdown) String() string {
	return m.renderMarkdown()
}

// Error returns error.
func (m *Markdown) Error() error {
	return m.err
}

// PlainText set plain text
func (m *Markdown) PlainText(text string) *Markdown {
	para := ast.NewParagraph()
	para.AppendChild(para, ast.NewString([]byte(text)))
	m.appendBlock(para)
	return m
}

// PlainTextf set plain text with format
func (m *Markdown) PlainTextf(format string, args ...interface{}) *Markdown {
	return m.PlainText(fmt.Sprintf(format, args...))
}

// Build writes markdown text to output destination.
func (m *Markdown) Build() error {
	if _, err := fmt.Fprint(m.dest, m.String()); err != nil {
		if m.err != nil {
			return fmt.Errorf("failed to write markdown text: %w: %s", err, m.err.Error())
		}
		return fmt.Errorf("failed to write markdown text: %w", err)
	}
	return m.err
}

func (m *Markdown) addHeading(level int, text string) *Markdown {
	heading := ast.NewHeading(level)
	heading.AppendChild(heading, ast.NewString([]byte(text)))
	m.headers = append(m.headers, headerInfo{level: TableOfContentsDepth(level), text: text})
	m.appendBlock(heading)
	return m
}

// H1 is markdown header.
func (m *Markdown) H1(text string) *Markdown { return m.addHeading(1, text) }

// H1f is markdown header with format.
func (m *Markdown) H1f(format string, args ...interface{}) *Markdown {
	return m.H1(fmt.Sprintf(format, args...))
}

// H2 is markdown header.
func (m *Markdown) H2(text string) *Markdown { return m.addHeading(2, text) }

// H2f is markdown header with format.
func (m *Markdown) H2f(format string, args ...interface{}) *Markdown {
	return m.H2(fmt.Sprintf(format, args...))
}

// H3 is markdown header.
func (m *Markdown) H3(text string) *Markdown { return m.addHeading(3, text) }

// H3f is markdown header with format.
func (m *Markdown) H3f(format string, args ...interface{}) *Markdown {
	return m.H3(fmt.Sprintf(format, args...))
}

// H4 is markdown header.
func (m *Markdown) H4(text string) *Markdown { return m.addHeading(4, text) }

// H4f is markdown header with format.
func (m *Markdown) H4f(format string, args ...interface{}) *Markdown {
	return m.H4(fmt.Sprintf(format, args...))
}

// H5 is markdown header.
func (m *Markdown) H5(text string) *Markdown { return m.addHeading(5, text) }

// H5f is markdown header with format.
func (m *Markdown) H5f(format string, args ...interface{}) *Markdown {
	return m.H5(fmt.Sprintf(format, args...))
}

// H6 is markdown header.
func (m *Markdown) H6(text string) *Markdown { return m.addHeading(6, text) }

// H6f is markdown header with format.
func (m *Markdown) H6f(format string, args ...interface{}) *Markdown {
	return m.H6(fmt.Sprintf(format, args...))
}
