package main

import (
	"fmt"
	"io"
	"strings"
)

// Global flag to disable all colors.
// This is useful when the output is not a terminal (e.g., a file or a pipe).
var DisableColors = false

// Reset code to clear all formatting.
const Reset = "\033[0m"

// ANSI escape codes for text styles.
const (
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Underline = "\033[4m"
	Blink     = "\033[5m" // Note: Blink is not widely supported.
	Reverse   = "\033[7m"
	Hidden    = "\033[8m"
	Strike    = "\033[9m"
)

// ANSI escape codes for standard foreground colors.
const (
	FgBlack   = "\033[30m"
	FgRed     = "\033[31m"
	FgGreen   = "\033[32m"
	FgYellow  = "\033[33m"
	FgBlue    = "\033[34m"
	FgMagenta = "\033[35m"
	FgCyan    = "\033[36m"
	FgWhite   = "\033[37m"
)

// ANSI escape codes for standard background colors.
const (
	BgBlack   = "\033[40m"
	BgRed     = "\033[41m"
	BgGreen   = "\033[42m"
	BgYellow  = "\033[43m"
	BgBlue    = "\033[44m"
	BgMagenta = "\033[45m"
	BgCyan    = "\033[46m"
	BgWhite   = "\033[47m"
)

// Color is a struct that holds the ANSI escape codes for a full style.
// This makes it easy to define and reuse custom styles.
type Color struct {
	fg     string
	bg     string
	styles []string
}

// NewColor creates a new Color struct.
func NewColor() *Color {
	return &Color{
		fg:     "",
		bg:     "",
		styles: make([]string, 0),
	}
}

// WithFg sets the foreground color of the Color struct.
func (c *Color) WithFg(code string) *Color {
	c.fg = code
	return c
}

// WithBg sets the background color of the Color struct.
func (c *Color) WithBg(code string) *Color {
	c.bg = code
	return c
}

// WithStyle adds a style code (e.g., Bold, Underline) to the Color struct.
func (c *Color) WithStyle(codes ...string) *Color {
	c.styles = append(c.styles, codes...)
	return c
}

// Colorize wraps a string with the specified ANSI escape codes and then resets
// the formatting. This is the core function that makes the output compatible with
// fmt.Printf and fmt.Sprintf.
func Colorize(text string, codes ...string) string {
	if DisableColors {
		return text
	}
	startCode := strings.Join(codes, "")
	return fmt.Sprintf("%s%s%s", startCode, text, Reset)
}

// Sprint formats a string with the specified codes and returns it.
func Sprint(text string, codes ...string) string {
	return Colorize(text, codes...)
}

// Sprintln formats a string with the specified codes and returns it,
// with a newline at the end.
func Sprintln(text string, codes ...string) string {
	return Colorize(text, codes...) + "\n"
}

// Print prints the provided string with the specified codes to standard output.
func Print(text string, codes ...string) {
	fmt.Print(Colorize(text, codes...))
}

// Println prints the provided string with the specified codes to standard output,
// followed by a newline.
func Println(text string, codes ...string) {
	fmt.Println(Colorize(text, codes...))
}

// Printf is a wrapper for fmt.Printf that applies styling to the first argument.
func Printf(format string, a ...interface{}) {
	if len(a) > 1 {
		// If there are multiple arguments, assume the first one is the text
		// to be formatted and the rest are the style codes.
		text, ok := a[0].(string)
		if ok {
			codes, ok := a[1].([]string)
			if ok {
				fmt.Printf(Colorize(text, codes...), a[2:]...)
				return
			}
		}
	}
	// Fallback to standard Printf if the arguments don't match the expected format.
	fmt.Printf(format, a...)
}

// Fprint writes the formatted string with the specified codes to the provided writer.
func Fprint(w io.Writer, text string, codes ...string) (n int, err error) {
	return fmt.Fprint(w, Colorize(text, codes...))
}

// Fprintln writes the formatted string with the specified codes to the provided writer,
// followed by a newline.
func Fprintln(w io.Writer, text string, codes ...string) (n int, err error) {
	return fmt.Fprintln(w, Colorize(text, codes...))
}

// --- Dynamic Color Functions ---

// Color256 generates a 256-color code from a number (0-255).
// The code is a prefix for the string.
func Color256(code int) string {
	return fmt.Sprintf("\033[38;5;%dm", code)
}

// BgColor256 generates a 256-color background code from a number (0-255).
func BgColor256(code int) string {
	return fmt.Sprintf("\033[48;5;%dm", code)
}

// TrueColor generates a 24-bit RGB color code.
// R, G, B should be between 0 and 255.
func TrueColor(r, g, b int) string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
}

// BgTrueColor generates a 24-bit RGB background color code.
func BgTrueColor(r, g, b int) string {
	return fmt.Sprintf("\033[48;2;%d;%d;%dm", r, g, b)
}

// --- Convenience Functions ---

func Info(text string) string { return Sprint(text, FgBlue) }
func Success(text string) string { return Sprint(text, FgGreen, Bold) }
func Warning(text string) string { return Sprint(text, FgYellow) }
func Error(text string) string { return Sprint(text, FgRed, Bold) }

