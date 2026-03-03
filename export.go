package tables

import "strings"

// ToCSV returns the table as a RFC 4180-compliant CSV string.
// ANSI escape sequences are stripped from all cells since CSV is plain text.
// Column alignment and border style are not applied — those are terminal-only
// concepts. Use WriteTo or Print for styled terminal output.
// https://www.rfc-editor.org/rfc/rfc4180.html
func (t *Table) ToCSV() string {
	if len(t.headers) == 0 {
		return ""
	}

	var sb strings.Builder

	// header row
	for i, h := range t.headers {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(csvField(StripANSI(string(h))))	// Damn quite functional style coding right here, lol.
	}

	sb.WriteByte('\n')

	// data rows
	for i, row := range t.rows {

		// skip row separator
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
				// apply max width truncation if set
				if t.maxWidths[j] > 0 {
					cell = TruncateToWidth(cell, t.maxWidths[j])
				}
			}

			sb.WriteString(csvField(cell))
		}
		sb.WriteByte('\n')
	}

	return sb.String()
}

// csvField wraps a field in quotes if it contains a comma, double-quote, or
// newline, escaping any inner double-quotes by doubling them (RFC 4180 §2.7).
func csvField(s string) string {
	if !strings.ContainsAny(s, ",\"\n\r") {
		return s
	}
	return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
}

// ToHTML returns a self-contained HTML <table> block representing the table.
// ANSI escape sequences are stripped — use inline CSS or classes externally
// if you need styling. The fragment can be embedded directly in any HTML page.
// Alignment is rendered via the style attribute on each <td>/<th>.
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
	sb.WriteString("    </tr>\n  </thead>\n  <tbody>\n")

	for i, row := range t.rows {
		if t.rowKinds[i] == rowSeparator {
			// represent a separator as a visually distinct row via a class;
			// callers can style tr.separator however they like
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
// Alignment colons are placed in the separator row per GFM spec.
// ANSI sequences are stripped. Separators added via AddSeparator are omitted
// since GFM tables have no equivalent concept.
func (t *Table) ToMarkdown() string {
	if len(t.headers) == 0 {
		return ""
	}

	// First pass — measure column widths for pretty alignment.
	// We want each column wide enough for its header and all its cells so the
	// Markdown source is human-readable, not just spec-valid.
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
				w := len(StripANSI(string(row[j])))
				if w > colWidths[j] {
					colWidths[j] = w
				}
			}
		}
	}
	// enforce maxWidths
	for i, mw := range t.maxWidths {
		if mw > 0 && colWidths[i] > mw {
			colWidths[i] = mw
		}
	}
	// GFM separator dashes need at least 3 chars to be valid
	for i := range colWidths {
		if colWidths[i] < 3 {
			colWidths[i] = 3
		}
	}

	var sb strings.Builder

	// header row
	sb.WriteByte('|')
	for i, h := range t.headers {
		cell := StripANSI(string(h))
		sb.WriteByte(' ')
		sb.WriteString(mdPad(cell, colWidths[i], t.aligns[i]))
		sb.WriteString(" |")
	}
	sb.WriteByte('\n')

	// separator row with alignment colons
	sb.WriteByte('|')
	for i := range t.headers {
		sb.WriteString(mdSeparator(colWidths[i], t.aligns[i]))
		sb.WriteByte('|')
	}
	sb.WriteByte('\n')

	// data rows
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

	return sb.String()
}

// --- helpers -----------------------------------------------------------------

// htmlAlign maps an Align constant to a CSS text-align value.
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

// htmlEscape escapes the five predefined XML/HTML entities.
func htmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&#34;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}

// mdPad pads or truncates a string to width for a Markdown cell,
// respecting alignment so the source table looks neat.
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

// mdSeparator produces a GFM alignment separator cell, e.g. ":---:", "---:", etc.
func mdSeparator(width int, align Align) string {
	dashes := strings.Repeat("-", width)
	switch align {
	case AlignCenter:
		return ":" + dashes + ":"
	case AlignRight:
		return dashes + ":"
	default:
		return dashes + "-" // left or default: plain dashes (at least 3+1 = 4)
	}
}

// itoa is a tiny int-to-string helper to avoid importing strconv here.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}
