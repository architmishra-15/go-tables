# Go Tables Library

A high-performance Go library for creating beautiful, customizable terminal tables with Unicode support, ANSI color compatibility, and optimized rendering.

## Features

- 🎨 **Multiple Border Styles**: Single, double, rounded, ASCII, and borderless styles
- 🌍 **Full Unicode Support**: Proper width calculation for CJK characters, emojis, and combining marks
- 🎯 **Column Alignment**: Left, center, and right alignment for each column
- 🚀 **High Performance**: Byte-first API with buffer pooling for minimal allocations
- 🌈 **ANSI Color Support**: Full compatibility with colored text and escape sequences
- 📝 **Flexible Input**: Support for strings, bytes, numbers, and any type
- 📊 **Width Management**: Automatic sizing with optional maximum width constraints
- 💾 **Multiple Output Formats**: Print to stdout, get as string, or write to any io.Writer

## Installation

```bash
go get github.com/architmishra-15/go-table
```

## Quick Start

```go
package main

import "github.com/architmishra-15/go-table"

func main() {
    // Simple table creation
    table := NewFromStrings("Name", "Age", "City").
        AddRow("Alice", 25, "New York").
        AddRow("Bob", 30, "Los Angeles").
        SetStyle(StyleRounded).
        Print()
}
```

## API Reference

### Table Creation

#### `New(headers ...[]byte) *Table`
Creates a new table with byte slice headers (most performant).

```go
table := New([]byte("Name"), []byte("Age"), []byte("City"))
```

#### `NewFromStrings(headers ...string) *Table`
Creates a new table with string headers (convenience function).

```go
table := NewFromStrings("Name", "Age", "City")
```

### Adding Data

#### `AddRow(values ...interface{}) *Table`
Adds a row with mixed data types. Supports strings, numbers, booleans, and byte slices.

```go
table.AddRow("Alice", 25, true, []byte("data"))
```

#### `AddRowBytes(values ...[]byte) *Table`
Adds a row from byte slices directly (highest performance).

```go
table.AddRowBytes([]byte("Alice"), []byte("25"), []byte("Active"))
```

### Styling

#### `SetStyle(style Style) *Table`
Sets the border style for the table.

Available styles:
- `StyleSingle` - Single line borders (┌─┐│└┘)
- `StyleDouble` - Double line borders (╔═╗║╚╝)
- `StyleRounded` - Rounded corners (╭─╮│╰╯)
- `StyleASCII` - ASCII-only borders (+|-+)
- `StyleNone` - No borders, spacing only

```go
table.SetStyle(StyleDouble)
```

#### `SetAlign(col int, align Align) *Table`
Sets alignment for a specific column.

Available alignments:
- `AlignLeft` - Left-aligned text
- `AlignCenter` - Center-aligned text
- `AlignRight` - Right-aligned text

```go
table.SetAlign(0, AlignLeft).
      SetAlign(1, AlignCenter).
      SetAlign(2, AlignRight)
```

#### `SetMaxWidth(col int, width int) *Table`
Sets maximum width for a specific column (0 = unlimited).

```go
table.SetMaxWidth(0, 20) // Limit first column to 20 characters
```

### Advanced Configuration

#### `SetWidthFunc(fn WidthFunc) *Table`
Sets a custom width calculation function for special character handling.

```go
customWidthFunc := func(r rune) int {
    // Custom width logic
    return 1
}
table.SetWidthFunc(customWidthFunc)
```

### Output Methods

#### `String() string`
Returns the formatted table as a string.

```go
tableString := table.String()
fmt.Print(tableString)
```

#### `Print()`
Prints the table directly to stdout.

```go
table.Print()
```

#### `WriteTo(w io.Writer) (int64, error)`
Writes the table to any io.Writer interface.

```go
file, _ := os.Create("output.txt")
bytesWritten, err := table.WriteTo(file)
```

## Color Support

The library includes built-in ANSI color support:

### Color Functions
```go
Info("text")        // Blue text
Success("text")     // Green bold text  
Warning("text")     // Yellow text
Error("text")       // Red bold text
Sprint("text", Bold, FgRed) // Custom styling
```

### Color Constants
```go
// Foreground colors
FgBlack, FgRed, FgGreen, FgYellow, FgBlue, FgMagenta, FgCyan, FgWhite

// Background colors  
BgBlack, BgRed, BgGreen, BgYellow, BgBlue, BgMagenta, BgCyan, BgWhite

// Text styles
Bold, Dim, Underline, Reverse, Strike

// Advanced colors
Color256(42)           // 256-color palette
TrueColor(255, 128, 0) // 24-bit RGB color
```

### Colored Table Example
```go
table := NewFromStrings("Status", "Service", "Uptime").
    AddRow(Success("ONLINE"), "API", "99.9%").
    AddRow(Warning("DEGRADED"), "Cache", "87.2%").
    AddRow(Error("OFFLINE"), "DB", "0.0%").
    Print()
```

## Unicode Support

The library properly handles:
- **CJK Characters**: Chinese, Japanese, Korean text with correct 2-character width
- **Emojis**: Full emoji support with proper width calculation
- **Combining Characters**: Diacritical marks and zero-width characters
- **RTL Text**: Arabic and Hebrew text support

### Unicode Example
```go
table := NewFromStrings("Language", "Greeting", "Flag").
    AddRow("English", "Hello World", "🇺🇸").
    AddRow("Japanese", "こんにちは世界", "🗾").
    AddRow("Chinese", "你好世界", "🇨🇳").
    AddRow("Arabic", "مرحبا بالعالم", "🇸🇦").
    Print()
```

## Performance Optimizations

### Byte-First Approach
For maximum performance, use byte slices directly:

```go
// Fastest - no string conversions
table := New([]byte("Name"), []byte("Age")).
    AddRowBytes([]byte("Alice"), []byte("25"))

// Good - minimal conversions  
table := NewFromStrings("Name", "Age").
    AddRow("Alice", 25)
```

### Buffer Pooling
The library uses sync.Pool for buffer reuse, minimizing garbage collection pressure during table rendering.

### ANSI Sequence Handling
Width calculations intelligently strip ANSI escape sequences for accurate column sizing while preserving color formatting.

## Width Calculation Functions

### String Width Functions
```go
StringWidth(s string) int                    // Basic width calculation
StringWidthANSI(s string) int               // ANSI-aware width calculation  
StringWidthCustom(s string, WidthFunc) int  // Custom width function
```

### Byte Width Functions  
```go
StringWidthBytes(b []byte) int                    // Byte slice width
StringWidthBytesANSI(b []byte) int               // ANSI-aware byte width
StringWidthBytesCustom(b []byte, WidthFunc) int  // Custom byte width
```

### Character Classification
```go
RuneWidth(r rune) int      // Returns 0, 1, or 2 for character width
IsWideRune(r rune) bool    // True for double-width characters
IsZeroWidthRune(r rune) bool // True for combining/zero-width characters
```

### Text Manipulation
```go
TruncateToWidth(s string, maxWidth int) string     // Truncate with ellipsis
TruncateToWidthBytes(b []byte, maxWidth int) []byte // Byte version
PadToWidth(s string, width int, align Align) string // Pad to specific width
```

## ANSI Sequence Utilities

```go
StripANSI(s string) string          // Remove all ANSI escape sequences
StripANSIBytes(b []byte) []byte     // Byte version
HasANSI(s string) bool              // Check if string contains ANSI sequences
```

## Error Handling

The library is designed to be robust:
- Invalid UTF-8 sequences are handled gracefully
- Missing columns are filled with empty values
- Out-of-range column operations are ignored silently
- All functions prioritize continuing execution over panicking

## Examples

### Basic Table
```go
NewFromStrings("Name", "Score").
    AddRow("Alice", 95).
    AddRow("Bob", 87).
    SetStyle(StyleSingle).
    Print()
```

### Advanced Styling
```go
NewFromStrings("Server", "Status", "Load").
    SetStyle(StyleRounded).
    SetAlign(1, AlignCenter).
    SetAlign(2, AlignRight).
    SetMaxWidth(0, 15).
    AddRow("web-01", Success("OK"), "23%").
    AddRow("db-primary", Warning("HIGH"), "89%").
    Print()
```

### File Output
```go
table := NewFromStrings("Item", "Price").
    AddRow("Laptop", "$999").
    AddRow("Mouse", "$25")

file, _ := os.Create("report.txt")
defer file.Close()
table.WriteTo(file)
```

## Best Practices

1. **Use byte slices** for high-performance scenarios
2. **Set column alignments** for better visual presentation  
3. **Apply max widths** to prevent overly wide tables
4. **Use appropriate styles** for your terminal environment
5. **Leverage color functions** for status indicators and emphasis
6. **Consider ANSI compatibility** when outputting to files or pipes

## License

This project is licensed under the MIT License - see the LICENSE file for details.
