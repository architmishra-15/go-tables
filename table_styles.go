// table_styles.go

package tables

import (
	"unicode/utf8"
)

type Select struct {
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

var (
	// StyleSingle uses single line Unicode box drawing characters
	// Pattern: ┌─┬┐│├┼┤└┴┘
	StyleSingle = Style{
		TopLeft:     '┌',
		TopRight:    '┐',
		BottomLeft:  '└',
		BottomRight: '┘',
		Horizontal:  '─',
		Vertical:    '│',
		Cross:       '┼',
		TopTee:      '┬',
		BottomTee:   '┴',
		LeftTee:     '├',
		RightTee:    '┤',
	}

	// StyleDouble uses double line Unicode box drawing characters
	// Pattern: ╔═╦╗║╠╬╣╚╩╝
	StyleDouble = Style{
		TopLeft:     '╔',
		TopRight:    '╗',
		BottomLeft:  '╚',
		BottomRight: '╝',
		Horizontal:  '═',
		Vertical:    '║',
		Cross:       '╬',
		TopTee:      '╦',
		BottomTee:   '╩',
		LeftTee:     '╠',
		RightTee:    '╣',
	}

	// StyleRounded uses rounded corner Unicode characters
	// Pattern: ╭─┬╮│├┼┤╰┴╯
	StyleRounded = Style{
		TopLeft:     '╭',
		TopRight:    '╮',
		BottomLeft:  '╰',
		BottomRight: '╯',
		Horizontal:  '─',
		Vertical:    '│',
		Cross:       '┼',
		TopTee:      '┬',
		BottomTee:   '┴',
		LeftTee:     '├',
		RightTee:    '┤',
	}

	// StyleASCII uses only ASCII characters for maximum compatibility
	// Pattern: +--+|+|+-+
	StyleASCII = Style{
		TopLeft:     '+',
		TopRight:    '+',
		BottomLeft:  '+',
		BottomRight: '+',
		Horizontal:  '-',
		Vertical:    '|',
		Cross:       '+',
		TopTee:      '+',
		BottomTee:   '+',
		LeftTee:     '+',
		RightTee:    '+',
	}

	// StyleNone provides no borders, only spacing for clean text output
	StyleNone = Style{
		TopLeft:     ' ',
		TopRight:    ' ',
		BottomLeft:  ' ',
		BottomRight: ' ',
		Horizontal:  ' ',
		Vertical:    ' ',
		Cross:       ' ',
		TopTee:      ' ',
		BottomTee:   ' ',
		LeftTee:     ' ',
		RightTee:    ' ',
	}
)

func (s Style) GetBorderChar(position string) rune {
	switch position {
	case "top-left", "tl", "topleft", "topl", "tleft", "tlft":
		return s.TopLeft
	case "top-right", "tr", "topright", "topr":
		return s.TopRight
	case "bottom-left", "btm-lft", "bl", "btmleft", "btml":
		return s.BottomLeft
	case "bottom-right", "br", "bottomright", "botr":
		return s.BottomRight
	case "horizontal", "h":
		return s.Horizontal
	case "vertical", "v":
		return s.Vertical
	case "cross":
		return s.Cross
	case "top-tee", "tt", "t-tee":
		return s.TopTee
	case "bottom-tee", "btm-tee", "btm-t", "btmt":
		return s.BottomTee
	case "left-tee", "lt", "l-tee":
		return s.LeftTee
	case "right-tee", "rt", "r-tee":
		return s.RightTee
	default:
		return ' '
	}
}

func (s Style) IsNone() bool {
	return s.TopLeft == ' ' && s.Horizontal == ' ' && s.Vertical == ' '
}

// render a complete border line using the style
func (s Style) renderBorderLine(widths []int, lineType string) []byte {
	if len(widths) == 0 {
		return []byte{}
	}

	// Calculation of total length needed
	totalLength := 1
	for i, width := range widths {
		totalLength += width + 2 // width + 2 padding space
		if i < len(widths)-1 {
			totalLength += 1 // separator char
		}
	}

	totalLength += 1 // end char
	totalLength += 1 // newline

	result := make([]byte, 0, totalLength)

	// type of to use based on line type
	var startChar, endChar, sepChar, fillChar rune

	switch lineType {
	case "top":
		startChar = s.TopLeft
		endChar = s.TopRight
		sepChar = s.TopTee
		fillChar = s.Horizontal
	case "middle", "header":
		startChar = s.LeftTee
		endChar = s.RightTee
		sepChar = s.Cross
		fillChar = s.Horizontal
	case "bottom":
		startChar = s.BottomLeft
		endChar = s.BottomRight
		sepChar = s.BottomTee
		fillChar = s.Horizontal
	default: // fallback
		startChar = s.LeftTee
		endChar = s.RightTee
		sepChar = s.Cross
		fillChar = s.Horizontal
	}

	result = appendRune(result, startChar)
	for i, width := range widths {
		// Add padding and content spaces
		for j := 0; j < width+2; j++ {
			result = appendRune(result, fillChar)
		}

		// Add separator (except for last column)
		if i < len(widths)-1 {
			result = appendRune(result, sepChar)
		}
	}

	result = appendRune(result, endChar)
	result = append(result, '\n')

	return result

}

func appendRune(b []byte, r rune) []byte {
	if r < 0x80 {
		// ASCII fast path
		return append(b, byte(r))
	}
	// UTF-8 encoding for non-ASCII
	var buf [4]byte
	n := utf8.EncodeRune(buf[:], r)
	return append(b, buf[:n]...)
}
