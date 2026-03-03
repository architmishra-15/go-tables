package tables

import (
	"strings"
)

// ToCSV returns the table as a RFC 4180-compliant CSV string.
// ANSI escape sequences are stripped from all cells since CSV is plain text.
// Column alignment and border style are not applied — those are terminal-only
// concepts. Separator rows are skipped. The footer row, if set, is appended last.
func (t *Table) ToCSV() string {
	if len(t.headers) == 0 {
		return ""
	}

	var sb strings.Builder

	for i, h := range t.headers {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(csvField(StripANSI(string(h))))
	}
	sb.WriteByte('\n')

	for i, row := range t.rows {
		if t.rowKinds[i] == rowSeparator {
			continue
		}
		for j := range t.headers {
			if j > 0 {
				sb.WriteByte(',')
			}
			var cell string
			if j < len(row) {
				cell = StripANSI(string(row[j]))
				if t.maxWidths[j] > 0 {
					cell = TruncateToWidth(cell, t.maxWidths[j])
				}
			}
			sb.WriteString(csvField(cell))
		}
		sb.WriteByte('\n')
	}

	if t.footer != nil {
		for j := range t.headers {
			if j > 0 {
				sb.WriteByte(',')
			}
			var cell string
			if j < len(t.footer) {
				cell = StripANSI(string(t.footer[j]))
			}
			sb.WriteString(csvField(cell))
		}
		sb.WriteByte('\n')
	}

	return sb.String()
}

func csvField(s string) string {
	if !strings.ContainsAny(s, ",\"\n\r") {
		return s
	}
	return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
}

// ToHTML returns a self-contained HTML <table> block.
// ANSI sequences are stripped. Separator rows become <tr class="separator">
// so you can style them with CSS. The footer, if set, goes in a <tfoot> block.
func (t *Table) ToHTML() string {
	if len(t.headers) == 0 {
		return ""
	}

	var sb strings.Builder

	sb.WriteString("<table>\n  <thead>\n    <tr>\n")
	for i, h := range t.headers {
		align := htmlAlign(t.aligns[i])
		cell := htmlEscape(StripANSI(string(h)))
		sb.WriteString("      <th style=\"text-align:")
		sb.WriteString(align)
		sb.WriteString("\">")
		sb.WriteString(cell)
		sb.WriteString("</th>\n")
	}
	sb.WriteString("    </tr>\n  </thead>\n")

	if t.footer != nil {
		sb.WriteString("  <tfoot>\n    <tr>\n")
		for j := range t.headers {
			align := htmlAlign(t.aligns[j])
			var cell string
			if j < len(t.footer) {
				cell = htmlEscape(StripANSI(string(t.footer[j])))
			}
			sb.WriteString("      <td style=\"text-align:")
			sb.WriteString(align)
			sb.WriteString("\"><strong>")
			sb.WriteString(cell)
			sb.WriteString("</strong></td>\n")
		}
		sb.WriteString("    </tr>\n  </tfoot>\n")
	}

	sb.WriteString("  <tbody>\n")
	for i, row := range t.rows {
		if t.rowKinds[i] == rowSeparator {
			sb.WriteString("    <tr class=\"separator\"><td colspan=\"")
			sb.WriteString(itoa(len(t.headers)))
			sb.WriteString("\"></td></tr>\n")
			continue
		}
		sb.WriteString("    <tr>\n")
		for j := range t.headers {
			align := htmlAlign(t.aligns[j])
			var cell string
			if j < len(row) {
				cell = StripANSI(string(row[j]))
				if t.maxWidths[j] > 0 {
					cell = TruncateToWidth(cell, t.maxWidths[j])
				}
			}
			sb.WriteString("      <td style=\"text-align:")
			sb.WriteString(align)
			sb.WriteString("\">")
			sb.WriteString(htmlEscape(cell))
			sb.WriteString("</td>\n")
		}
		sb.WriteString("    </tr>\n")
	}
	sb.WriteString("  </tbody>\n</table>")

	return sb.String()
}

// ToMarkdown returns the table in GitHub Flavored Markdown pipe-table format.
// Alignment colons are placed in the separator row per the GFM spec.
// ANSI sequences are stripped. AddSeparator rows are omitted — GFM has no
// equivalent. The footer row, if set, is appended as a plain data row.
func (t *Table) ToMarkdown() string {
	if len(t.headers) == 0 {
		return ""
	}

	colWidths := make([]int, len(t.headers))
	for i, h := range t.headers {
		colWidths[i] = len(StripANSI(string(h)))
	}
	for i, row := range t.rows {
		if t.rowKinds[i] == rowSeparator {
			continue
		}
		for j := range t.headers {
			if j < len(row) {
				if w := len(StripANSI(string(row[j]))); w > colWidths[j] {
					colWidths[j] = w
				}
			}
		}
	}
	if t.footer != nil {
		for j := range t.headers {
			if j < len(t.footer) {
				if w := len(StripANSI(string(t.footer[j]))); w > colWidths[j] {
					colWidths[j] = w
				}
			}
		}
	}
	for i, mw := range t.maxWidths {
		if mw > 0 && colWidths[i] > mw {
			colWidths[i] = mw
		}
	}
	for i := range colWidths {
		if colWidths[i] < 3 {
			colWidths[i] = 3
		}
	}

	var sb strings.Builder

	sb.WriteByte('|')
	for i, h := range t.headers {
		sb.WriteByte(' ')
		sb.WriteString(mdPad(StripANSI(string(h)), colWidths[i], t.aligns[i]))
		sb.WriteString(" |")
	}
	sb.WriteByte('\n')

	sb.WriteByte('|')
	for i := range t.headers {
		sb.WriteString(mdSeparator(colWidths[i], t.aligns[i]))
		sb.WriteByte('|')
	}
	sb.WriteByte('\n')

	for i, row := range t.rows {
		if t.rowKinds[i] == rowSeparator {
			continue
		}
		sb.WriteByte('|')
		for j := range t.headers {
			var cell string
			if j < len(row) {
				cell = StripANSI(string(row[j]))
				if t.maxWidths[j] > 0 {
					cell = TruncateToWidth(cell, t.maxWidths[j])
				}
			}
			sb.WriteByte(' ')
			sb.WriteString(mdPad(cell, colWidths[j], t.aligns[j]))
			sb.WriteString(" |")
		}
		sb.WriteByte('\n')
	}

	if t.footer != nil {
		sb.WriteByte('|')
		for j := range t.headers {
			var cell string
			if j < len(t.footer) {
				cell = StripANSI(string(t.footer[j]))
				if t.maxWidths[j] > 0 {
					cell = TruncateToWidth(cell, t.maxWidths[j])
				}
			}
			sb.WriteByte(' ')
			sb.WriteString(mdPad(cell, colWidths[j], t.aligns[j]))
			sb.WriteString(" |")
		}
		sb.WriteByte('\n')
	}

	return sb.String()
}

func htmlAlign(a Align) string {
	switch a {
	case AlignCenter:
		return "center"
	case AlignRight:
		return "right"
	default:
		return "left"
	}
}

func htmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&#34;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}

func mdPad(s string, width int, align Align) string {
	cur := len(s)
	if cur >= width {
		return s
	}
	pad := strings.Repeat(" ", width-cur)
	switch align {
	case AlignRight:
		return pad + s
	case AlignCenter:
		half := (width - cur) / 2
		return strings.Repeat(" ", half) + s + strings.Repeat(" ", width-cur-half)
	default:
		return s + pad
	}
}

func mdSeparator(width int, align Align) string {
	dashes := strings.Repeat("-", width)
	switch align {
	case AlignCenter:
		return ":" + dashes + ":"
	case AlignRight:
		return dashes + ":"
	default:
		return dashes + "-"
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}
