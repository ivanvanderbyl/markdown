# markdown

`markdown` is a lightweight builder for generating Markdown programmatically using the [goldmark](https://github.com/yuin/goldmark) AST. It mirrors the API provided by [`github.com/nao1215/markdown`](https://github.com/nao1215/markdown) while dropping external dependencies (such as `tablewriter`) and relying entirely on goldmark nodes for construction and rendering.

## Features

- Chainable builder API for headings, lists, blockquotes, tables, callouts, badges, links, and more
- Table of Contents generation with configurable depth
- Table rendering with per-column alignment and auto-formatting helpers
- Custom table handling without `tablewriter`
- Goldmark-backed internal representation ensures consistent Markdown output across platforms
- Simple syntax sugar helpers for inline formatting

## Installation

```bash
go get github.com/alcova-ai/markdown
```

## Quick Start

```go
package main

import (
    "fmt"
    "os"

    "github.com/alcova-ai/markdown"
)

func main() {
    md := markdown.NewMarkdown(os.Stdout)
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
}
```

Output:

```markdown
# Guide to markdown
Markdown built through goldmark AST.
| Feature | Description                                   |
| ------- | --------------------------------------------- |
| TOC     | Generate nested table of contents             |
| Tables  | Alignment-aware rendering without tablewriter |
```

## Building Documents

Every method on `*Markdown` returns the same builder, enabling fluent composition. When you’re done, call `Build()` to write the rendered Markdown to the provided `io.Writer`.

```go
md := markdown.NewMarkdown(os.Stdout)
md.H1("Release Notes").
    H2("v1.0.0").
    BulletList("Initial release", "Markdown builder", "Table support").
    LF().
    Important("Remember to pin dependencies").
    Build()
```

## Adding a Table of Contents

`TableOfContents` consumes the recorded heading metadata and writes a Markdown TOC up to a specified depth.

```go
md := markdown.NewMarkdown(os.Stdout)
md.H1("Project").
    H2("Overview").
    H2("Usage").
    TableOfContents(markdown.TableOfContentsDepthH2).
    Build()
```

The generated TOC uses bullet indentation to reflect heading levels.

## Working with Tables

Tables are defined through `TableSet`. The renderer automatically pads columns to fit the widest cell and emits separators honoring column alignment.

```go
md.Table(markdown.TableSet{
    Header: []string{"Left", "Center", "Right"},
    Rows: [][]string{
        {"L", "C", "R"},
    },
    Alignment: []markdown.TableAlignment{
        markdown.AlignLeft,
        markdown.AlignCenter,
        markdown.AlignRight,
    },
})
```

Output:

```markdown
| Left | Center | Right |
| :--- | :----: | ----: |
| L    |   C    |     R |
```

### Custom Table Helpers

`CustomTable` applies optional formatting on top of standard rendering. Currently, it supports:

- `AutoFormatHeaders`: Title-cases header cells by splitting on whitespace

```go
md.CustomTable(markdown.TableSet{
    Header: []string{"first name", "status"},
    Rows: [][]string{{"Alice", "active"}},
}, markdown.TableOptions{AutoFormatHeaders: true})
```

## Inline Formatting Helpers

Use the standalone helpers for inline Markdown strings:

```go
markdown.Bold("text")       // **text**
markdown.Italic("text")     // *text*
markdown.Link("Docs", "https://example.com")
markdown.Image("Logo", "https://example.com/logo.png")
markdown.Highlight("Note") // ==Note==
```

## Callouts and Badges

The builder supports GitHub-style callouts and shield badges:

```go
md.Note("Heads up!")
md.Tip("Try the new API.")
md.BlueBadge("stable")
```

Each callout renders a blockquote with the appropriate label (e.g., `[!NOTE]`).

## Rendering Programmatically Generated Data

The powered example below demonstrates building a weekly price table from structs:

```go
func ExampleMarkdown_Table_bars() {
    bars := []Bar{ /* ... seven days ... */ }

    rows := make([][]string, len(bars))
    for i, bar := range bars {
        rows[i] = []string{
            bar.Timestamp.Format("2006-01-02"),
            fmt.Sprintf("%.2f", bar.Open),
            fmt.Sprintf("%.2f", bar.High),
            fmt.Sprintf("%.2f", bar.Low),
            fmt.Sprintf("%.2f", bar.Close),
            fmt.Sprintf("%d", bar.Volume),
            fmt.Sprintf("%d", bar.TradeCount),
            fmt.Sprintf("%.2f", bar.VWAP),
        }
    }

    md := markdown.NewMarkdown(os.Stdout)
    md.H2("Daily Bars")
    md.Table(markdown.TableSet{
        Header: []string{"Day", "Open", "High", "Low", "Close", "Volume", "Trades", "VWAP"},
        Rows:   rows,
    })
    md.Build()
}
```

## Error Handling

Most builder methods return the builder and only record errors internally. Retrieve the combined error from `Error()` or defer the check to `Build()`:

```go
md.Table(markdown.TableSet{Header: []string{"A"}, Rows: [][]string{{"x", "y"}}})
if err := md.Build(); err != nil {
    log.Fatalf("build failed: %v", err)
}
```

## Testing

Run project tests with:

```bash
go test ./...
```

## License

MIT Licensed. See [LICENSE](LICENSE) for details.
