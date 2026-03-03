// extra.go

package tables

import (
    "fmt"
    "sort"
    "strconv"
)

// SortByColumn sorts the table's data rows by the values in the given column
// (0-indexed), either ascending or descending. The sort is stable so rows with
// equal values preserve their original relative order.
//
// Separator rows added via AddSeparator are removed during sorting because
// their positions become meaningless after reordering. If you need separators
// after sorting, add them again once the sort is done.
//
// The sort is lexicographic (string comparison) by default. If every value in
// the column looks like a number, a numeric sort is used instead so that
// "10" sorts after "9" rather than before it.
func (t *Table) SortByColumn(col int, ascending bool) *Table {
    if col < 0 || col >= len(t.headers) || len(t.rows) == 0 {
        return t
    }

    // Strip separator rows first — their positions are meaningless post-sort.
    cleanRows := make([][][]byte, 0, len(t.rows))
    cleanKinds := make([]rowKind, 0, len(t.rows))
    for i, row := range t.rows {
        if t.rowKinds[i] == rowSeparator {
            continue
        }
        cleanRows = append(cleanRows, row)
        cleanKinds = append(cleanKinds, rowData)
    }
    t.rows = cleanRows
    t.rowKinds = cleanKinds

    // Decide whether to use numeric or lexicographic comparison.
    numeric := isNumericColumn(t.rows, col)

    sort.SliceStable(t.rows, func(i, j int) bool {
        a := cellString(t.rows[i], col)
        b := cellString(t.rows[j], col)

        var less bool
        if numeric {
            less = parseFloat(a) < parseFloat(b)
        } else {
            less = a < b
        }

        if ascending {
            return less
        }
        return !less
    })

    return t
}

// isNumericColumn returns true if every non-empty cell in the given column of
// rows parses as a float64.
func isNumericColumn(rows [][][]byte, col int) bool {
    for _, row := range rows {
        s := cellString(row, col)
        if s == "" {
            continue
        }
        if _, err := strconv.ParseFloat(s, 64); err != nil {
            return false
        }
    }
    return true
}

// cellString returns the ANSI-stripped string value of cell (row, col),
// returning "" safely if the index is out of range.
func cellString(row [][]byte, col int) string {
    if col >= len(row) {
        return ""
    }
    return StripANSI(string(row[col]))
}

// parseFloat parses s as float64, returning 0 on any error.
func parseFloat(s string) float64 {
    f, _ := strconv.ParseFloat(s, 64)
    return f
}

// --- Footer ------------------------------------------------------------------

// SetFooter sets a footer row that is rendered after all data rows, separated
// by a border line. It accepts the same value types as AddRow. Pass nil to
// clear an existing footer.
//
// Example:
//
//	t.SetFooter("Total", 339, "—")
func (t *Table) SetFooter(values ...interface{}) *Table {
    if len(values) == 0 {
        t.footer = nil
        return t
    }

    row := make([][]byte, len(t.headers))
    for i, val := range values {
        if i >= len(t.headers) {
            break
        }
        switch v := val.(type) {
        case []byte:
            row[i] = make([]byte, len(v))
            copy(row[i], v)
        case string:
            row[i] = []byte(v)
        case int:
            row[i] = strconv.AppendInt(nil, int64(v), 10)
        case int64:
            row[i] = strconv.AppendInt(nil, v, 10)
        case float64:
            row[i] = strconv.AppendFloat(nil, v, 'f', -1, 64)
        case bool:
            row[i] = strconv.AppendBool(nil, v)
        default:
            row[i] = []byte(fmt.Sprintf("%v", v))
        }
    }
    for i := len(values); i < len(t.headers); i++ {
        row[i] = []byte{}
    }

    t.footer = row
    return t
}

// SetFooterColor sets the ANSI color/style applied to every cell in the footer
// row. Pass nil to clear.
func (t *Table) SetFooterColor(c *Color) *Table {
    t.footerColor = c
    return t
}

// ClearFooter removes the footer row.
func (t *Table) ClearFooter() *Table {
    t.footer = nil
    t.footerColor = nil
    return t
}
