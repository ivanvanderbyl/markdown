package markdown

import "fmt"

// RedBadge set text with red badge format.
func (m *Markdown) RedBadge(text string) *Markdown {
	return m.PlainText(fmt.Sprintf("![Badge](https://img.shields.io/badge/%s-red)", text))
}

// RedBadgef set text with red badge format.
func (m *Markdown) RedBadgef(format string, args ...interface{}) *Markdown {
	return m.RedBadge(fmt.Sprintf(format, args...))
}

// YellowBadge set text with yellow badge format.
func (m *Markdown) YellowBadge(text string) *Markdown {
	return m.PlainText(fmt.Sprintf("![Badge](https://img.shields.io/badge/%s-yellow)", text))
}

// YellowBadgef set text with yellow badge format.
func (m *Markdown) YellowBadgef(format string, args ...interface{}) *Markdown {
	return m.YellowBadge(fmt.Sprintf(format, args...))
}

// GreenBadge set text with green badge format.
func (m *Markdown) GreenBadge(text string) *Markdown {
	return m.PlainText(fmt.Sprintf("![Badge](https://img.shields.io/badge/%s-green)", text))
}

// GreenBadgef set text with green badge format.
func (m *Markdown) GreenBadgef(format string, args ...interface{}) *Markdown {
	return m.GreenBadge(fmt.Sprintf(format, args...))
}

// BlueBadge set text with blue badge format.
func (m *Markdown) BlueBadge(text string) *Markdown {
	return m.PlainText(fmt.Sprintf("![Badge](https://img.shields.io/badge/%s-blue)", text))
}

// BlueBadgef set text with blue badge format.
func (m *Markdown) BlueBadgef(format string, args ...interface{}) *Markdown {
	return m.BlueBadge(fmt.Sprintf(format, args...))
}
