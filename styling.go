package tables

// rowcol is a compact composite key for the per-cell color map.
// Using a struct as a map key is zero-allocation.
type rowcol struct {
	row, col int
}

// Apply wraps text with the receiver's ANSI codes and returns the result.
// If DisableColors is set or the Color is nil, the original text is returned.
func (c *Color) Apply(text string) string {
	if c == nil || DisableColors {
		return text
	}

	codes := make([]string, 0, 1+1+len(c.styles))
	if c.fg != "" {
		codes = append(codes, c.fg)
	}

	if c.bg != "" {
		codes = append(codes, c.bg)
	}

	codes = append(codes, c.styles...)
	if len(codes) == 0 {
		return text
	}

	return Colorize(text, codes...)
}

// --- Header styling ----------------------------------------------------------

// SetHeaderColor sets the ANSI color/style applied to every cell in the header
// row. It does not affect data rows. Pass nil to clear any existing header style.
//
// Example:
//
//	t.SetHeaderColor(
//	    tables.NewColor().WithFg(tables.FgCyan).WithStyle(tables.Bold),
//	)
func (t *Table) SetHeaderColor(c *Color) *Table {
	t.headerColor = c
	return t
}

// --- Row / column / cell coloring --------------------------------------------

// SetRowColor applies a color to every cell in the given data row (0-indexed,
// not counting separator rows). If the row index is out of range the call is a
// no-op, consistent with how SetAlign and SetMaxWidth behave.
func (t *Table) SetRowColor(row int, c *Color) *Table {
	if row < 0 {
		return t
	}

	if t.rowColors == nil {
		t.rowColors = make(map[int]*Color)
	}
	t.rowColors[row] = c
	return t
}

// SetColumnColor applies a color to every data cell in the given column
// (0-indexed). The header cell is NOT affected — use SetHeaderColor for that.
func (t *Table) SetColumnColor(col int, c *Color) *Table {
	if col < 0 || col >= len(t.headers) {
		return t
	}

	if t.colColors == nil {
		t.colColors = make(map[int]*Color)
	}
	t.colColors[col] = c
	return t
}

// SetCellColor applies a color to a single data cell at (row, col), both
// 0-indexed. Cell color takes priority over row and column colors.
func (t *Table) SetCellColor(row, col int, c *Color) *Table {
	if row < 0 || col < 0 || col >= len(t.headers) {
		return t
	}
	if t.cellColors == nil {
		t.cellColors = make(map[rowcol]*Color)
	}
	t.cellColors[rowcol{row, col}] = c
	return t
}

// cellColor resolves the effective color for a data cell, applying the
// priority: cell > row > column > nil.
func (t *Table) cellColor(row, col int) *Color {
	if t.cellColors != nil {
		if c, ok := t.cellColors[rowcol{row, col}]; ok {
			return c
		}
	}
	if t.rowColors != nil {
		if c, ok := t.rowColors[row]; ok {
			return c
		}
	}
	if t.colColors != nil {
		if c, ok := t.colColors[col]; ok {
			return c
		}
	}
	return nil
}

