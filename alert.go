package markdown

import (
	"fmt"

	"github.com/yuin/goldmark/ast"
)

func (m *Markdown) callout(label, text string) *Markdown {
	blockquote := ast.NewBlockquote()
	head := ast.NewParagraph()
	head.AppendChild(head, ast.NewString([]byte(label+"  ")))
	blockquote.AppendChild(blockquote, head)

	body := ast.NewParagraph()
	body.AppendChild(body, ast.NewString([]byte(text)))
	blockquote.AppendChild(blockquote, body)

	m.appendBlock(blockquote)
	return m
}

// Note set text with note format.
func (m *Markdown) Note(text string) *Markdown { return m.callout("[!NOTE]", text) }

// Notef set text with note format.
func (m *Markdown) Notef(format string, args ...interface{}) *Markdown {
	return m.Note(fmt.Sprintf(format, args...))
}

// Tip set text with tip format.
func (m *Markdown) Tip(text string) *Markdown { return m.callout("[!TIP]", text) }

// Tipf set text with tip format.
func (m *Markdown) Tipf(format string, args ...interface{}) *Markdown {
	return m.Tip(fmt.Sprintf(format, args...))
}

// Important set text with important format.
func (m *Markdown) Important(text string) *Markdown { return m.callout("[!IMPORTANT]", text) }

// Importantf set text with important format.
func (m *Markdown) Importantf(format string, args ...interface{}) *Markdown {
	return m.Important(fmt.Sprintf(format, args...))
}

// Warning set text with warning format.
func (m *Markdown) Warning(text string) *Markdown { return m.callout("[!WARNING]", text) }

// Warningf set text with warning format.
func (m *Markdown) Warningf(format string, args ...interface{}) *Markdown {
	return m.Warning(fmt.Sprintf(format, args...))
}

// Caution set text with caution format.
func (m *Markdown) Caution(text string) *Markdown { return m.callout("[!CAUTION]", text) }

// Cautionf set text with caution format.
func (m *Markdown) Cautionf(format string, args ...interface{}) *Markdown {
	return m.Caution(fmt.Sprintf(format, args...))
}
