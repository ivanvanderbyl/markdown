package markdown

import (
	"io"
	"testing"
)

func TestMarkdownHeadingsAndTOC(t *testing.T) {
	t.Parallel()

	md := NewMarkdown(io.Discard)
	md.H1("Introduction").
		H2("Overview").
		H3("Details").
		TableOfContents(TableOfContentsDepthH2)

	lf := lineFeed()
	want := "# Introduction" + lf +
		"## Overview" + lf +
		"### Details" + lf +
		"- [Introduction](#introduction)" + lf +
		"  - [Overview](#overview)" + lf +
		""

	if got := md.String(); got != want {
		t.Fatalf("unexpected markdown output\nwant: %q\ngot:  %q", want, got)
	}
}

func TestMarkdownLists(t *testing.T) {
	t.Parallel()

	lf := lineFeed()

	t.Run("bullet list", func(t *testing.T) {
		md := NewMarkdown(io.Discard)
		md.BulletList("Hello", "World")
		want := "- Hello" + lf + "- World"
		if got := md.String(); got != want {
			t.Fatalf("unexpected bullet list output\nwant: %q\ngot:  %q", want, got)
		}
	})

	t.Run("ordered list", func(t *testing.T) {
		md := NewMarkdown(io.Discard)
		md.OrderedList("First", "Second")
		want := "1. First" + lf + "2. Second"
		if got := md.String(); got != want {
			t.Fatalf("unexpected ordered list output\nwant: %q\ngot:  %q", want, got)
		}
	})

	t.Run("checkbox list", func(t *testing.T) {
		md := NewMarkdown(io.Discard)
		md.CheckBox([]CheckBoxSet{{Text: "Task", Checked: true}, {Text: "Review", Checked: false}})
		want := "- [x] Task" + lf + "- [ ] Review"
		if got := md.String(); got != want {
			t.Fatalf("unexpected checkbox list output\nwant: %q\ngot:  %q", want, got)
		}
	})
}

func TestMarkdownTable(t *testing.T) {
	t.Parallel()

	lf := lineFeed()
	md := NewMarkdown(io.Discard)
	md.Table(TableSet{
		Header: []string{"Name", "Age"},
		Rows:   [][]string{{"Alice", "24"}},
	})

	want := "| Name | Age |" + lf +
		"|---------|---------|" + lf +
		"| Alice | 24 |" + lf

	if got := md.String(); got != want {
		t.Fatalf("unexpected table output\nwant: %q\ngot:  %q", want, got)
	}
}

func TestMarkdownTableRenderingVariants(t *testing.T) {
	t.Parallel()

	lf := lineFeed()

	t.Run("alignment markers", func(t *testing.T) {
		md := NewMarkdown(io.Discard)
		md.Table(TableSet{
			Header:    []string{"Left", "Center", "Right"},
			Rows:      [][]string{{"L", "C", "R"}},
			Alignment: []TableAlignment{AlignLeft, AlignCenter, AlignRight},
		})

		want := "| Left | Center | Right |" + lf +
			"|:--------|:-------:|--------:|" + lf +
			"| L | C | R |" + lf

		if got := md.String(); got != want {
			t.Fatalf("unexpected alignment table output\nwant: %q\ngot:  %q", want, got)
		}
	})

	t.Run("custom table header formatting", func(t *testing.T) {
		md := NewMarkdown(io.Discard)
		md.CustomTable(TableSet{
			Header: []string{"first name", "status"},
			Rows:   [][]string{{"Alice", "active"}},
		}, TableOptions{AutoFormatHeaders: true})

		want := "| First Name | Status |" + lf +
			"|---------|---------|" + lf +
			"| Alice | active |" + lf

		if got := md.String(); got != want {
			t.Fatalf("unexpected custom table output\nwant: %q\ngot:  %q", want, got)
		}
	})
}
