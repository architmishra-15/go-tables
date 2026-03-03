# Go Tables Library

A dependency-free, high-performance Go library for rendering terminal tables. It started as a small utility to avoid copy-pasting table-rendering logic across CLI projects, and grew into something general enough to be worth publishing. The priorities, in order, are: correctness, speed, and a small API surface. Nothing here requires a third-party dependency.

---

## Installation

```bash
go get github.com/architmishra-15/go-tables
```

Requires Go 1.21 or later. No external dependencies — only the standard library.

---

## Quick Start

```go
package main

import tables "github.com/architmishra-15/go-tables"

func main() {
    tables.NewFromStrings("Name", "Age", "City").
        AddRow("Alice", 25, "New York").
        AddRow("Bob", 30, "Los Angeles").
        SetStyle(tables.StyleRounded).
        Print()
}
```

Output:
```
╭───────┬─────┬─────────────╮
│ Name  │ Age │ City        │
├───────┼─────┼─────────────┤
│ Alice │ 25  │ New York    │
│ Bob   │ 30  │ Los Angeles │
╰───────┴─────┴─────────────╯
```

---

## Table Creation

There are two constructors. Use `New` when you already have `[]byte` headers and want to avoid allocations. Use `NewFromStrings` everywhere else — it's just more convenient.

```go
// Byte-first, fastest path
t := tables.New([]byte("Name"), []byte("Age"))

// String convenience wrapper
t := tables.NewFromStrings("Name", "Age")
```

Both return `*Table` and support method chaining.

---

## Adding Data

### `AddRow(values ...interface{}) *Table`

Accepts strings, ints, int64, float64, booleans, byte slices, or anything that implements `fmt.Stringer`. Type conversion happens once during insertion, not during rendering.

```go
t.AddRow("Alice", 95, true, 3.14, []byte("raw"))
```

### `AddRowBytes(values ...[]byte) *Table`

The fastest way to add rows — no type switching, no conversions, just a direct copy into the internal buffer.

```go
t.AddRowBytes([]byte("Alice"), []byte("95"))
```

### `AddSeparator() *Table`

Inserts a horizontal border line at the current position in the table. Useful for grouping rows visually. You can add as many as you want and they compose correctly with all border styles.

```go
t.AddRow("Alice", 95).
 AddRow("Bob",   87).
 AddSeparator().
 AddRow("Total", 182)
```

```
╭───────┬───────╮
│ Name  │ Score │
├───────┼───────┤
│ Alice │ 95    │
│ Bob   │ 87    │
├───────┼───────┤
│ Total │ 182   │
╰───────┴───────╯
```

Note: separator rows are stripped when you call `SortByColumn` because their positions become meaningless after reordering. Add them again after sorting if you need them.

---

## Border Styles

Five styles are available out of the box:

| Constant | Description |
| --- | --- |
| `StyleSingle` | Single-line Unicode box drawing (`┌─┬┐│├┼┤└┴┘`) |
| `StyleDouble` | Double-line Unicode box drawing (`╔═╦╗║╠╬╣╚╩╝`) |
| `StyleRounded` | Rounded corners (`╭─╮│╰╯`) |
| `StyleASCII` | Plain ASCII (`+-|`) — safe on any terminal or log file |
| `StyleNone` | No borders, just spacing — good for copying into documents |

```go
t.SetStyle(tables.StyleDouble)
```

You can also define your own style by filling in a `Style` struct directly:

```go
custom := tables.Style{
    TopLeft: '╔', TopRight: '╗',
    // ... fill in all 11 fields
}
t.SetStyle(custom)
```

Use `PrintStyles()` to render a live preview of all built-in styles to stdout.

---

## Column Options

### Alignment

```go
t.SetAlign(0, tables.AlignLeft)    // default
t.SetAlign(1, tables.AlignCenter)
t.SetAlign(2, tables.AlignRight)
```

Alignment is respected in terminal output, Markdown export, and HTML export.

### Max Width

```go
t.SetMaxWidth(0, 20) // truncates column 0 to 20 display characters
```

Long values are truncated with an ellipsis (`...`). Set to `0` for unlimited (the default).

### Custom Width Function

By default, width is calculated using a compact embedded Unicode range table covering CJK, Hangul, Hiragana, Katakana, and common emoji. If you need a different heuristic, you can swap it out:

```go
t.SetWidthFunc(func(r rune) int {
    // your logic here
    return 1
})
```

---

## Coloring

The library has two layers of color: convenience functions for wrapping individual strings, and structural color that gets applied automatically during rendering.

### Convenience Functions

These wrap a string with ANSI codes and return it. They work in any context, not just tables.

```go
tables.Info("text")     // blue
tables.Success("text")  // green + bold
tables.Warning("text")  // yellow
tables.Error("text")    // red + bold
tables.Sprint("text", tables.Bold, tables.FgCyan) // custom
```

You can pass the result directly into `AddRow`:

```go
t.AddRow(tables.Success("ONLINE"), "API", "99.9%")
```

### Color Constants

```go
// Foreground
FgBlack, FgRed, FgGreen, FgYellow, FgBlue, FgMagenta, FgCyan, FgWhite

// Background
BgBlack, BgRed, BgGreen, BgYellow, BgBlue, BgMagenta, BgCyan, BgWhite

// Text styles
Bold, Dim, Underline, Blink, Reverse, Hidden, Strike

// Extended palettes
tables.Color256(214)           // 256-color foreground
tables.BgColor256(214)         // 256-color background
tables.TrueColor(255, 128, 0)  // 24-bit RGB foreground
tables.BgTrueColor(255, 128, 0) // 24-bit RGB background
```

To build a reusable color style:

```go
highlight := tables.NewColor().
    WithFg(tables.FgCyan).
    WithBg(tables.BgBlack).
    WithStyle(tables.Bold)

// then apply it to any string:
formatted := highlight.Apply("some text")
```

Set `tables.DisableColors = true` to strip all ANSI output globally — useful when piping to a file or a log aggregator.

### Structural Coloring

Instead of colorizing individual cells manually, you can attach color to entire rows, columns, or individual cells. These are applied automatically during rendering.

#### Header Color

```go
t.SetHeaderColor(tables.NewColor().WithFg(tables.FgCyan).WithStyle(tables.Bold))
```

#### Footer Color

```go
t.SetFooterColor(tables.NewColor().WithStyle(tables.Bold))
```

#### Row, Column, and Cell Color

```go
// Color an entire data row (0-indexed)
t.SetRowColor(2, tables.NewColor().WithFg(tables.FgRed))

// Color an entire column
t.SetColumnColor(1, tables.NewColor().WithFg(tables.FgYellow))

// Color a single cell — takes priority over row and column colors
t.SetCellColor(0, 1, tables.NewColor().WithFg(tables.FgGreen).WithStyle(tables.Bold))
```

Priority when multiple colors apply to the same cell: **cell > row > column**. The most specific one wins.

```go
tables.NewFromStrings("Name", "Score", "Grade").
    SetStyle(tables.StyleDouble).
    SetHeaderColor(tables.NewColor().WithStyle(tables.Bold, tables.Underline)).
    SetColumnColor(1, tables.NewColor().WithFg(tables.FgYellow)).
    SetRowColor(2, tables.NewColor().WithFg(tables.FgRed).WithStyle(tables.Dim)).
    SetCellColor(0, 1, tables.NewColor().WithFg(tables.FgGreen).WithStyle(tables.Bold)).
    AddRow("Alice", 99, "A+").   // cell (0,1) → green bold, beats column yellow
    AddRow("Bob",   87, "B").    // cell (1,1) → column yellow
    AddRow("Charlie", 45, "F").  // entire row 2 → red dim, beats column yellow
    Print()
```

---

## Footer

A footer row is rendered after all data rows, separated from them by a border line. It's intended for totals, averages, or any kind of summary.

```go
t.SetFooter("Total", 9, "$51.99")
t.SetFooterColor(tables.NewColor().WithStyle(tables.Bold))

// Remove the footer
t.ClearFooter()
```

Footer cells participate in column width measurement, so a wide footer value will expand the column correctly. The footer is never affected by `SortByColumn` — it always stays pinned at the bottom.

---

## Sorting

`SortByColumn` sorts data rows by the values in a given column. The sort is stable, so rows with equal values preserve their original relative order.

```go
t.SortByColumn(1, false) // column 1, descending
t.SortByColumn(0, true)  // column 0, ascending
```

The library detects automatically whether the column contains numeric data. If every non-empty value in the column parses as a `float64`, numeric comparison is used so that `"10"` sorts after `"9"`. Otherwise, it falls back to lexicographic comparison.

Separator rows are removed during sorting (their positions would be meaningless after reordering). If you need them, add them again after the sort call.

---

## Exporting

All three export formats strip ANSI escape sequences from cell content — color codes are a terminal concept and would corrupt a CSV file or break HTML rendering.

### CSV

```go
csv := t.ToCSV()
```

Follows RFC 4180: fields containing commas, double-quotes, or newlines are wrapped in double-quotes, and any inner double-quotes are escaped by doubling them. Separator rows are skipped. The footer row, if set, is appended as the last line.

### Markdown

```go
md := t.ToMarkdown()
```

Outputs a GitHub Flavored Markdown pipe table. Column alignment is expressed with the standard colon syntax in the separator row (`:---:` for center, `---:` for right). Separator rows added via `AddSeparator` are dropped since GFM has no equivalent concept. The footer, if set, is appended as a plain data row — GFM has no `<tfoot>`.

The output is padded to be readable as plain text, not just spec-valid. Each column is wide enough to accommodate its widest value without truncation (unless `SetMaxWidth` constrains it).

### HTML

```go
html := t.ToHTML()
```

Returns a `<table>` fragment — no `<html>`, `<head>`, or `<body>` wrapper, just the block you'd paste into an existing page. Structure:

- Headers go in `<thead>` with `<th>` elements
- The footer, if set, goes in `<tfoot>` (placed before `<tbody>` in source order, which is correct per the HTML spec and lets browsers render a sticky footer on long tables)
- Data rows go in `<tbody>` with `<td>` elements
- Separator rows become `<tr class="separator">` — style them however you want with CSS:

```css
tr.separator td { border-top: 2px solid #e0e0e0; padding: 0; height: 0; }
```

Alignment is expressed as an inline `style="text-align:..."` attribute on each cell. Footer cells are wrapped in `<strong>` by default.

---

## Output Methods

```go
t.Print()               // writes to stdout
s := t.String()         // returns as string
t.WriteTo(w io.Writer)  // writes to any io.Writer (file, buffer, HTTP response, etc.)
```

`WriteTo` implements `io.WriterTo`, so it works directly with `bufio.Writer`, `http.ResponseWriter`, `os.File`, and anything else that satisfies the interface.

---

## Unicode Support

Width calculation uses a compact, embedded range table — no external dependencies. The following are handled correctly:

- **CJK Unified Ideographs** — Chinese, Japanese, Korean characters render as width 2
- **Hangul, Hiragana, Katakana** — all width 2
- **Emoji** — common emoji ranges (U+1F300–U+1FAFF) treated as width 2
- **Combining characters** — diacritical marks and zero-width joiners are width 0
- **Full-width ASCII variants** — width 2

The heuristic covers the vast majority of real-world cases. Edge cases like complex emoji sequences (multi-codepoint ZWJ sequences) may not be measured perfectly — this is a documented limitation of going dependency-free. If you need perfect accuracy for an unusual script, swap in your own function with `SetWidthFunc`.

```go
tables.NewFromStrings("Language", "Greeting", "Flag").
    SetStyle(tables.StyleRounded).
    AddRow("Japanese", "こんにちは世界", "🗾").
    AddRow("Chinese",  "你好世界",       "🇨🇳").
    AddRow("Korean",   "안녕하세요",     "🇰🇷").
    Print()
```

---

## Performance

The library is built around a byte-first internal representation. Strings passed to `AddRow` are converted to `[]byte` once during insertion and never converted back until output. `AddRowBytes` skips even that conversion.

Rendering uses `sync.Pool` to reuse `bytes.Buffer` instances across calls, which keeps GC pressure low in tight loops or repeated renders. The buffer is sized to avoid reallocations for typical tables.

For high-throughput scenarios:

```go
// fastest possible path
t := tables.New([]byte("ID"), []byte("Value"))
for _, item := range items {
    t.AddRowBytes(item.IDBytes, item.ValueBytes)
}
t.WriteTo(w)
```

---

## Width Utility Functions

These are exported because they're useful outside of table rendering too — measuring terminal output width, padding strings for aligned output, etc.

```go
// Display width of a string (no ANSI handling)
tables.StringWidth(s string) int

// Display width, ignoring ANSI escape sequences
tables.StringWidthANSI(s string) int

// Same, but operating on []byte directly
tables.StringWidthBytes(b []byte) int
tables.StringWidthBytesANSI(b []byte) int

// Width of a single rune: 0, 1, or 2
tables.RuneWidth(r rune) int
tables.IsWideRune(r rune) bool
tables.IsZeroWidthRune(r rune) bool

// Truncate to a display width, adding ellipsis if truncated
tables.TruncateToWidth(s string, maxWidth int) string
tables.TruncateToWidthBytes(b []byte, maxWidth int) []byte

// Pad to a display width with alignment
tables.PadToWidth(s string, width int, align tables.Align) string
tables.PadToWidthBytes(b []byte, width int, align tables.Align) []byte

// ANSI utilities
tables.StripANSI(s string) string
tables.StripANSIBytes(b []byte) []byte
tables.HasANSI(s string) bool
tables.HasANSIBytes(b []byte) bool
```

---

## Architectural Notes

A few decisions that aren't obvious from the API:

**Byte-first internals.** All cell data is stored as `[][]byte` internally. This avoids repeated string ↔ `[]byte` round-trips during rendering and makes ANSI stripping cheap (byte scanning rather than rune decoding).

**`rowKind` sentinel.** Rather than storing separator rows as a special cell value or a separate list of indices, each row has a parallel `rowKind` tag (`rowData` or `rowSeparator`). This keeps the rows and their metadata in sync automatically when rows are appended or reordered, without needing to update a separate index.

**Sparse maps for structural colors.** `rowColors`, `colColors`, and `cellColors` are `nil` until first use, and keyed by index rather than stored as full-length slices. Most tables only color a handful of rows or columns — allocating a full-length slice for every table regardless of whether coloring is used would waste memory for the common case.

**`rowcol` struct as map key.** The per-cell color map uses `map[rowcol]*Color` where `rowcol` is a plain `struct{ row, col int }`. Go can hash fixed-size structs without boxing them, so lookups are allocation-free.

**Footer is not a row.** The footer is stored separately from `t.rows` and rendered after the main loop. This means `SortByColumn` cannot accidentally reorder it, and it doesn't interfere with row index calculations used by `SetRowColor` and `SetCellColor`.

**Export strips ANSI.** All three export formats (`ToCSV`, `ToMarkdown`, `ToHTML`) call `StripANSI` on every cell. You can freely pass colored strings into `AddRow` and export the same table to both terminal and file formats without having to maintain two versions of the data.

**`render` is the single source of truth.** Both `String()` and `WriteTo()` delegate to a private `render(buf *bytes.Buffer)` method. This means the rendering logic only exists in one place, and the two output methods are just thin wrappers that handle buffer pool lifecycle.

---

## License

AGPL-3.0 — see the [LICENSE](./LICENSE) file for details.
