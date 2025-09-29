package markdown

import "fmt"

// Link return text with link format.
func Link(text, url string) string {
	return fmt.Sprintf("[%s](%s)", text, url)
}

// Image return text with image format.
func Image(text, url string) string {
	return fmt.Sprintf("![%s](%s)", text, url)
}

// Strikethrough return text with strikethrough format.
func Strikethrough(text string) string {
	return fmt.Sprintf("~~%s~~", text)
}

// Bold return text with bold format.
func Bold(text string) string {
	return fmt.Sprintf("**%s**", text)
}

// Italic return text with italic format.
func Italic(text string) string {
	return fmt.Sprintf("*%s*", text)
}

// BoldItalic return text with bold and italic format.
func BoldItalic(text string) string {
	return fmt.Sprintf("***%s***", text)
}

// Code return text with code format.
func Code(text string) string {
	return fmt.Sprintf("`%s`", text)
}

// Highlight return text with highlight format.
func Highlight(text string) string {
	return fmt.Sprintf("==%s==", text)
}
