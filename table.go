// table.go
package tables

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"sync"
)

// Constants
type Align   int

// rowKind distinguishes between normal data rows and separator rows.
type rowKind int	// empty structs (struct{}) could also be used here, but since it'll be later used as a slice, it'd still end up taking atleast a byte.

const (
	AlignLeft Align = iota
	AlignCenter
	AlignRight
)

const (
	rowData rowKind = iota
	rowSeparator
)

// Style represents border characters for table rendering
type Style struct {
	TopLeft     rune
	TopRight    rune
	BottomLeft  rune
	BottomRight rune
	Horizontal  rune
	Vertical    rune
	Cross       rune
	TopTee      rune
	BottomTee   rune
	LeftTee     rune
	RightTee    rune
}

// Table represents a table with headers and rows stored as bytes
type Table struct {
	headers   [][]byte   // Column headers as bytes
	rows      [][][]byte // Each row contains multiple cells, each cell is []byte
	rowKinds  []rowKind	 // Parallel to rows — rowData or rowSeparator
	style     Style
	aligns    []Align   // Alignment per column
	maxWidths []int     // Max width per column (0 = unlimited)
	widthFunc WidthFunc // Pluggable width calculation function

	// Styling
	headerColor *Color
	rowColors   map[int]*Color
	colColors   map[int]*Color
	cellColors  map[rowcol]*Color

	// Buffer pool for performance
	bufPool *sync.Pool
}

// Buffer pool for reusing byte buffers
var defaultBufPool = &sync.Pool{
	New: func() any {
		return &bytes.Buffer{}
	},
}

// PrintStyles displays all available table styles with visual examples
func PrintStyles() {
	fmt.Println("Available Table Styles")
	fmt.Println("======================")
	fmt.Println()

	styles := []struct {
		name        string
		style       Style
		description string
	}{
		{
			name:        "StyleSingle",
			style:       StyleSingle,
			description: "Single line Unicode box drawing characters",
		},
		{
			name:        "StyleDouble", 
			style:       StyleDouble,
			description: "Double line Unicode box drawing characters",
		},
		{
			name:        "StyleRounded",
			style:       StyleRounded,
			description: "Rounded corner Unicode characters",
		},
		{
			name:        "StyleASCII",
			style:       StyleASCII,
			description: "ASCII-only characters for maximum compatibility",
		},
		{
			name:        "StyleNone",
			style:       StyleNone,
			description: "No borders, spacing only for clean text output",
		},
	}

	for i, styleInfo := range styles {
		fmt.Printf("%d. %s\n", i+1, styleInfo.name)
		fmt.Printf("   %s\n", styleInfo.description)
		fmt.Println()

		sampleTable := NewFromStrings("Header1", "Header2", "Header3").
			SetStyle(styleInfo.style).
			AddRow("Row 1", "Data A", "Value X").
			AddRow("Row 2", "Data B", "Value Y")

		fmt.Print(sampleTable.String())
		fmt.Println()

		// Add separator between styles (except for the last one)
		if i < len(styles)-1 {
			fmt.Println("---")
			fmt.Println()
		}
	}

	fmt.Println("Usage Examples:")
	fmt.Println("---------------")
	fmt.Println()
	fmt.Println("// Set style during table creation")
	fmt.Println("table := NewFromStrings(\"Name\", \"Age\").SetStyle(StyleSingle)")
	fmt.Println()
	fmt.Println("// Change style after creation")
	fmt.Println("table.SetStyle(StyleDouble)")
	fmt.Println()
	fmt.Println("// Available styles:")
	for _, styleInfo := range styles {
		fmt.Printf("// - %s: %s\n", styleInfo.name, styleInfo.description)
	}
}

// New creates a new table with the given headers (accepts bytes directly)
func New(headers ...[]byte) *Table {
	t := &Table{
		headers:   make([][]byte, len(headers)),
		rows:      make([][][]byte, 0),
		rowKinds:  make([]rowKind, 0),
		style:     StyleSingle, // Default to single line style
		aligns:    make([]Align, len(headers)),
		maxWidths: make([]int, len(headers)),
		widthFunc: DefaultWidthFunc, // Default width calculation
		bufPool:   defaultBufPool,
	}

	// Copy headers to avoid shared slice issues
	for i, header := range headers {
		t.headers[i] = make([]byte, len(header))
		copy(t.headers[i], header)
	}

	return t
}

// NewFromStrings creates a new table from string headers (convenience function)
func NewFromStrings(headers ...string) *Table {
	byteHeaders := make([][]byte, len(headers))
	for i, header := range headers {
		byteHeaders[i] = []byte(header)
	}
	return New(byteHeaders...)
}

// AddRow adds a row to the table, preferring byte inputs for performance
func (t *Table) AddRow(values ...any) *Table {
	if len(values) == 0 {
		return t
	}

	row := make([][]byte, len(t.headers))

	for i, val := range values {
		if i >= len(t.headers) {
			break // Don't exceed header count
		}

		// Convert interface{} to []byte efficiently - prioritize []byte inputs
		switch v := val.(type) {
		case []byte:
			// Direct byte slice - make a copy to avoid shared slice issues
			row[i] = make([]byte, len(v))
			copy(row[i], v)
		case string:
			row[i] = []byte(v) // Only convert when necessary
		case int:
			row[i] = strconv.AppendInt(nil, int64(v), 10)
		case int64:
			row[i] = strconv.AppendInt(nil, v, 10)
		case float64:
			row[i] = strconv.AppendFloat(nil, v, 'f', -1, 64)
		case bool:
			row[i] = strconv.AppendBool(nil, v)
		default:
			// Fallback to string conversion (avoid this path for performance)
			row[i] = fmt.Appendf(nil, "%v", v)
		}
	}

	// Fill remaining columns with empty bytes if row is shorter
	for i := len(values); i < len(t.headers); i++ {
		row[i] = []byte{}
	}

	t.rows = append(t.rows, row)
	t.rowKinds = append(t.rowKinds, rowData)
	return t
}

// AddRowBytes adds a row from byte slices directly (fastest method)
func (t *Table) AddRowBytes(values ...[]byte) *Table {
	if len(values) == 0 {
		return t
	}

	row := make([][]byte, len(t.headers))

	for i, val := range values {
		if i >= len(t.headers) {
			break
		}
		// Make a copy to avoid shared slice issues
		row[i] = make([]byte, len(val))
		copy(row[i], val)
	}

	// Fill remaining columns with empty bytes
	for i := len(values); i < len(t.headers); i++ {
		row[i] = []byte{}
	}

	t.rows = append(t.rows, row)
	t.rowKinds = append(t.rowKinds, rowData)
	return t
}

// AddSeparator inserts a horizontal separator line at the current position in
// the table. The separator is rendered using the table's active border style,
// identical to the header/body divider. Multiple separators can be added to
// create logical groups of rows.
//
// Example:
//
//	t := NewFromStrings("Name", "Score").
//	    AddRow("Alice", 95).
//	    AddRow("Bob",   87).
//	    AddSeparator().
//	    AddRow("Total", 182)
func (t *Table) AddSeparator() *Table {
	// The actual cell content doesn't matter for separators; use a nil row so
	// the renderer can identify it quickly without allocating a full cell slice.
	t.rows = append(t.rows, nil)
	t.rowKinds = append(t.rowKinds, rowSeparator)
	return t
}

// SetStyle sets the border style for the table
func (t *Table) SetStyle(style Style) *Table {
	t.style = style
	return t
}

// SetAlign sets alignment for a specific column
func (t *Table) SetAlign(col int, align Align) *Table {
	if col >= 0 && col < len(t.aligns) {
		t.aligns[col] = align
	}
	return t
}

// SetMaxWidth sets maximum width for a specific column
func (t *Table) SetMaxWidth(col int, width int) *Table {
	if col >= 0 && col < len(t.maxWidths) {
		t.maxWidths[col] = width
	}
	return t
}

// SetWidthFunc sets a custom width calculation function
func (t *Table) SetWidthFunc(fn WidthFunc) *Table {
	t.widthFunc = fn
	return t
}

// measureColumns calculates the width needed for each column
func (t *Table) measureColumns() []int {
	if len(t.headers) == 0 {
		return nil
	}

	widths := make([]int, len(t.headers))

	// Measure header widths using ANSI-aware width calculation
	for i, header := range t.headers {
		widths[i] = MeasureWidthIgnoreANSIBytesCustom(header, t.widthFunc)
	}

	// Measure row widths
	for i, row := range t.rows {

		if t.rowKinds[i] == rowSeparator {
			continue // separators don't affect column widths
		}

		for i, cell := range row {
			if i < len(widths) {
				cellWidth := MeasureWidthIgnoreANSIBytesCustom(cell, t.widthFunc)
				// if cellWidth > widths[i] {
				// 	widths[i] = cellWidth
				// }
				widths[i] = max(widths[i], cellWidth)
			}
		}
	}

	// Apply max width constraints
	for i, maxWidth := range t.maxWidths {
		if maxWidth > 0 && widths[i] > maxWidth {
			widths[i] = maxWidth
		}
	}

	return widths
}

// alignCell aligns a cell's content within the given width
func (t *Table) alignCell(cell []byte, width int, align Align) []byte {
	cellWidth := MeasureWidthIgnoreANSIBytesCustom(cell, t.widthFunc)

	if cellWidth >= width {
		// Truncate if too long - need to preserve ANSI sequences
		return t.truncateWithANSI(cell, width)
	}

	// Pad the cell while preserving ANSI sequences
	return t.padWithANSI(cell, width, cellWidth, align)
}

// truncateWithANSI truncates text while preserving ANSI sequences
func (t *Table) truncateWithANSI(cell []byte, maxWidth int) []byte {
	if !HasANSIBytes(cell) {
		return TruncateToWidthBytes(cell, maxWidth)
	}

	// For ANSI text, we need to be more careful
	// This is a simplified version - could be optimized further
	stripped := StripANSIBytes(cell)
	if StringWidthBytesCustom(stripped, t.widthFunc) <= maxWidth {
		return cell // Fits even with ANSI codes
	}

	// Truncate the stripped version and add ellipsis
	truncated := TruncateToWidthBytes(stripped, maxWidth)
	return truncated
}

// padWithANSI pads text while preserving ANSI sequences
func (t *Table) padWithANSI(cell []byte, targetWidth, currentWidth int, align Align) []byte {
	padding := targetWidth - currentWidth
	if padding <= 0 {
		return cell
	}

	switch align {
	case AlignCenter:
		leftPad := padding / 2
		rightPad := padding - leftPad
		result := make([]byte, len(cell)+padding)

		// Left padding
		for i := range leftPad {
			result[i] = ' '
		}

		// Original content
		copy(result[leftPad:], cell)

		// Right padding
		for i := range rightPad {
			result[leftPad+len(cell)+i] = ' '
		}

		return result

	case AlignRight:
		result := make([]byte, len(cell)+padding)

		// Left padding
		for i := range padding {
			result[i] = ' '
		}

		// Original content
		copy(result[padding:], cell)
		return result

	default: // AlignLeft
		result := make([]byte, len(cell)+padding)

		// Original content
		copy(result, cell)

		// Right padding
		for i := range padding {
			result[len(cell)+i] = ' '
		}

		return result
	}
}

// renderBorder renders a border line using the table's style
func (t *Table) renderBorder(buf *bytes.Buffer, widths []int, borderType string) {
	if len(widths) == 0 {
		return
	}

	// Use the style to render the border
	borderBytes := t.style.renderBorderLine(widths, borderType)
	buf.Write(borderBytes)
}

// renderRow renders a single data row using the table's style
func (t *Table) renderRow(buf *bytes.Buffer, row [][]byte, widths []int, rowIdx int) {
	if len(widths) == 0 {
		return
	}

	// Use vertical character from style
	verticalChar := t.style.Vertical

	buf.WriteRune(verticalChar) // Left border

	for i, width := range widths {
		buf.WriteByte(' ') // Left padding

		var cell []byte
		if i < len(row) {
			cell = row[i]
		}

		align := AlignLeft
		if i < len(t.aligns) {
			align = t.aligns[i]
		}

		aligned := string(t.alignCell(cell, width, align))

		// apply color — header vs data row
		if rowIdx == -1 {
			aligned = t.headerColor.Apply(aligned)
		} else {
			aligned = t.cellColor(rowIdx, i).Apply(aligned)
		}

		buf.WriteString(aligned)
		buf.WriteByte(' ')          // Right padding
		buf.WriteRune(verticalChar) // Column separator / Right border
	}

	buf.WriteByte('\n')
}

// render writes the complete table into buf.
func (t *Table) render(buf *bytes.Buffer) {
	if len(t.headers) == 0 {
		return
	}

	widths := t.measureColumns()

	t.renderBorder(buf, widths, "top")
	t.renderRow(buf, t.headers, widths, -1)      // -1 = header
	t.renderBorder(buf, widths, "middle")

	dataIdx := 0
	for i, row := range t.rows {
		if t.rowKinds[i] == rowSeparator {
			t.renderBorder(buf, widths, "middle")
			continue
		}
		t.renderRow(buf, row, widths, dataIdx)
		dataIdx++
	}

	t.renderBorder(buf, widths, "bottom")
}

// String returns the formatted table as a string
func (t *Table) String() string {
	if len(t.headers) == 0 {
		return ""
	}

	// Get buffer from pool
	buf := t.bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer t.bufPool.Put(buf)

	t.render(buf)

	// Create a copy of the buffer content to return
	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())
	return string(result)
}

// Print prints the table directly to stdout
func (t *Table) Print() {
	fmt.Print(t.String())
}

// WriteTo writes the table to any io.Writer
func (t *Table) WriteTo(w io.Writer) (int64, error) {
	if len(t.headers) == 0 {
		return 0, nil
	}

	// Get buffer from pool
	buf := t.bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer t.bufPool.Put(buf)

	t.render(buf)
	// Write directly from buffer to avoid string conversion
	return buf.WriteTo(w)
}
