// width.go

package tables

import "unicode/utf8"

// unicodeRange represents a range of Unicode code points with their display width
type unicodeRange struct {
	start rune
	end   rune
	width int
}

// Embedded Unicode width tables - compact ranges for common wide characters
// Based on Unicode 15.0 East Asian Width property and common emoji ranges
var wideRanges = []unicodeRange{
	// CJK Unified Ideographs
	{0x4E00, 0x9FFF, 2},   // CJK Unified Ideographs
	{0x3400, 0x4DBF, 2},   // CJK Extension A
	{0x20000, 0x2A6DF, 2}, // CJK Extension B
	{0x2A700, 0x2B73F, 2}, // CJK Extension C
	{0x2B740, 0x2B81F, 2}, // CJK Extension D
	{0x2B820, 0x2CEAF, 2}, // CJK Extension E
	{0x2CEB0, 0x2EBEF, 2}, // CJK Extension F

	// Hangul
	{0xAC00, 0xD7AF, 2}, // Hangul Syllables
	{0x1100, 0x115F, 2}, // Hangul Jamo
	{0x1160, 0x11FF, 2}, // Hangul Jamo Extended-A
	{0xA960, 0xA97F, 2}, // Hangul Jamo Extended-B

	// Hiragana and Katakana
	{0x3040, 0x309F, 2}, // Hiragana
	{0x30A0, 0x30FF, 2}, // Katakana
	{0x31F0, 0x31FF, 2}, // Katakana Phonetic Extensions

	// CJK Symbols and Punctuation
	{0x3000, 0x303F, 2}, // CJK Symbols and Punctuation
	{0x2E80, 0x2EFF, 2}, // CJK Radicals Supplement
	{0x2F00, 0x2FDF, 2}, // Kangxi Radicals
	{0x2FF0, 0x2FFF, 2}, // Ideographic Description Characters

	// Full-width Forms
	{0xFF01, 0xFF60, 2}, // Fullwidth ASCII variants
	{0xFFE0, 0xFFE6, 2}, // Fullwidth symbol variants

	// Common Emoji ranges (width 2 for display purposes)
	{0x1F300, 0x1F5FF, 2}, // Miscellaneous Symbols and Pictographs
	{0x1F600, 0x1F64F, 2}, // Emoticons
	{0x1F680, 0x1F6FF, 2}, // Transport and Map Symbols
	{0x1F700, 0x1F77F, 2}, // Alchemical Symbols
	{0x1F780, 0x1F7FF, 2}, // Geometric Shapes Extended
	{0x1F800, 0x1F8FF, 2}, // Supplemental Arrows-C
	{0x1F900, 0x1F9FF, 2}, // Supplemental Symbols and Pictographs
	{0x1FA00, 0x1FA6F, 2}, // Chess Symbols
	{0x1FA70, 0x1FAFF, 2}, // Symbols and Pictographs Extended-A

	// Additional wide characters
	{0x2460, 0x24FF, 2}, // Enclosed Alphanumerics
	{0x25A0, 0x25FF, 2}, // Geometric Shapes
	{0x2600, 0x26FF, 2}, // Miscellaneous Symbols
	{0x2700, 0x27BF, 2}, // Dingbats
}

// Zero-width and combining characters (width 0)
var zeroWidthRanges = []unicodeRange{
	{0x0300, 0x036F, 0}, // Combining Diacritical Marks
	{0x1AB0, 0x1AFF, 0}, // Combining Diacritical Marks Extended
	{0x1DC0, 0x1DFF, 0}, // Combining Diacritical Marks Supplement
	{0x20D0, 0x20FF, 0}, // Combining Diacritical Marks for Symbols
	{0xFE20, 0xFE2F, 0}, // Combining Half Marks
	{0x200B, 0x200F, 0}, // Zero Width Space, ZWNJ, ZWJ, etc.
	{0x2028, 0x2029, 0}, // Line/Paragraph Separators
	{0x202A, 0x202E, 0}, // Bidirectional format characters
	{0x2060, 0x2064, 0}, // Word Joiner, etc.
}

// RuneWidth returns the display width of a single rune
// Returns 0 for zero-width, 1 for normal width, 2 for wide characters
func RuneWidth(r rune) int {
	// Fast path for ASCII
	if r < 0x80 {
		if r >= 0x20 {
			return 1 // Printable ASCII
		}
		return 0 // Control characters
	}

	// Check zero-width ranges first (most common for combining marks)
	for _, rang := range zeroWidthRanges {
		if r >= rang.start && r <= rang.end {
			return 0
		}
	}

	// Check wide character ranges
	for _, rang := range wideRanges {
		if r >= rang.start && r <= rang.end {
			return rang.width
		}
	}

	// Default to width 1 for everything else
	return 1
}

// StringWidth calculates the display width of a string
// This version does NOT handle ANSI escape sequences
func StringWidth(s string) int {
	width := 0
	for _, r := range s {
		width += RuneWidth(r)
	}
	return width
}

// StringWidthBytes calculates display width from byte slice
func StringWidthBytes(b []byte) int {
	width := 0
	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		if r == utf8.RuneError {
			// Invalid UTF-8, count as 1
			width += 1
			b = b[1:]
		} else {
			width += RuneWidth(r)
			b = b[size:]
		}
	}
	return width
}

// StringWidthANSI calculates display width ignoring ANSI escape sequences
func StringWidthANSI(s string) int {
	return MeasureWidthIgnoreANSI(s)
}

// StringWidthBytesANSI calculates display width from bytes ignoring ANSI
func StringWidthBytesANSI(b []byte) int {
	return MeasureWidthIgnoreANSIBytes(b)
}

// TruncateToWidth truncates a string to fit within specified display width
// Adds ellipsis (...) if truncated and there's room
func TruncateToWidth(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}

	width := 0
	var result []rune

	for _, r := range s {
		runeWidth := RuneWidth(r)
		if width+runeWidth > maxWidth {
			break
		}
		result = append(result, r)
		width += runeWidth
	}

	// Add ellipsis if truncated and there's room
	if len(result) < len([]rune(s)) && width <= maxWidth-3 {
		result = append(result, '.', '.', '.')
	}

	return string(result)
}

// TruncateToWidthBytes truncates byte slice to fit within display width
func TruncateToWidthBytes(b []byte, maxWidth int) []byte {
	if maxWidth <= 0 {
		return []byte{}
	}

	width := 0
	result := make([]byte, 0, len(b))

	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		runeWidth := RuneWidth(r)

		if width+runeWidth > maxWidth {
			break
		}

		result = append(result, b[:size]...)
		width += runeWidth
		b = b[size:]
	}

	// Add ellipsis if truncated and there's room
	if len(b) > 0 && width <= maxWidth-3 {
		result = append(result, '.', '.', '.')
	}

	return result
}

// PadToWidth pads a string to reach the specified display width
func PadToWidth(s string, width int, align Align) string {
	currentWidth := StringWidth(s)
	if currentWidth >= width {
		return s
	}

	padding := width - currentWidth
	spaces := make([]byte, padding)
	for i := range spaces {
		spaces[i] = ' '
	}

	switch align {
	case AlignCenter:
		leftPad := padding / 2
		rightPad := padding - leftPad
		leftSpaces := make([]byte, leftPad)
		rightSpaces := make([]byte, rightPad)
		for i := range leftSpaces {
			leftSpaces[i] = ' '
		}
		for i := range rightSpaces {
			rightSpaces[i] = ' '
		}
		return string(leftSpaces) + s + string(rightSpaces)

	case AlignRight:
		return string(spaces) + s

	default: // AlignLeft
		return s + string(spaces)
	}
}

// PadToWidthBytes pads byte slice to reach specified display width
func PadToWidthBytes(b []byte, width int, align Align) []byte {
	currentWidth := StringWidthBytes(b)
	if currentWidth >= width {
		return b
	}

	padding := width - currentWidth

	switch align {
	case AlignCenter:
		leftPad := padding / 2
		rightPad := padding - leftPad
		result := make([]byte, len(b)+padding)

		// Fill left padding
		for i := 0; i < leftPad; i++ {
			result[i] = ' '
		}

		// Copy original bytes
		copy(result[leftPad:], b)

		// Fill right padding
		for i := 0; i < rightPad; i++ {
			result[leftPad+len(b)+i] = ' '
		}

		return result

	case AlignRight:
		result := make([]byte, len(b)+padding)

		// Fill left padding
		for i := 0; i < padding; i++ {
			result[i] = ' '
		}

		// Copy original bytes
		copy(result[padding:], b)
		return result

	default: // AlignLeft
		result := make([]byte, len(b)+padding)

		// Copy original bytes
		copy(result, b)

		// Fill right padding
		for i := 0; i < padding; i++ {
			result[len(b)+i] = ' '
		}

		return result
	}
}

// WidthFunc is a pluggable function type for calculating character widths
// Allows users to provide custom width calculation logic
type WidthFunc func(rune) int

// DefaultWidthFunc is the default width calculation function
var DefaultWidthFunc WidthFunc = RuneWidth

// StringWidthCustom calculates string width using a custom width function
func StringWidthCustom(s string, widthFunc WidthFunc) int {
	width := 0
	for _, r := range s {
		width += widthFunc(r)
	}
	return width
}

// StringWidthBytesCustom calculates byte slice width using custom width function
func StringWidthBytesCustom(b []byte, widthFunc WidthFunc) int {
	width := 0
	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		if r == utf8.RuneError {
			width += 1 // Treat invalid UTF-8 as width 1
			b = b[1:]
		} else {
			width += widthFunc(r)
			b = b[size:]
		}
	}
	return width
}

// IsWideRune returns true if the rune has display width > 1
func IsWideRune(r rune) bool {
	return RuneWidth(r) > 1
}

// IsZeroWidthRune returns true if the rune has display width 0
func IsZeroWidthRune(r rune) bool {
	return RuneWidth(r) == 0
}
