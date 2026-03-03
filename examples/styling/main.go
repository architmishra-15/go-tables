// examples/styling/main.go

package main

import (
	tables "github.com/architmishra-15/go-tables"
)

func main() {
	// Header color only
	tables.NewFromStrings("Name", "Role", "Status").
		SetStyle(tables.StyleRounded).
		SetHeaderColor(tables.NewColor().WithFg(tables.FgCyan).WithStyle(tables.Bold)).
		AddRow("Alice", "Engineer", tables.Success("Active")).
		AddSeparator().
		AddRow("Bob", "Designer", tables.Success("Active")).
		AddSeparator().
		AddRow("Charlie", "Manager", tables.Warning("Away")).
		Print()

	// Column color (column 1 = Score dimmed, overridden by cell color on [0][1])
	tables.NewFromStrings("Name", "Score", "Grade").
		SetStyle(tables.StyleDouble).
		SetHeaderColor(tables.NewColor().WithStyle(tables.Bold, tables.Underline)).
		SetColumnColor(1, tables.NewColor().WithFg(tables.FgYellow)).
		SetRowColor(2, tables.NewColor().WithFg(tables.FgRed).WithStyle(tables.Dim)).
		SetCellColor(0, 1, tables.NewColor().WithFg(tables.FgGreen).WithStyle(tables.Bold)).
		AddRow("Alice", 99, "A+").   // row 0, col 1 → cell color (green bold)
		AddSeparator().
		AddRow("Bob", 87, "B").      // row 1, col 1 → column color (yellow)
		AddSeparator().
		AddRow("Charlie", 45, "F").  // row 2 → row color (red dim) wins over column
		Print()
}
