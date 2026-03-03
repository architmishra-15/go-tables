// examples/separator/main.go

package main

import (
	tables "github.com/architmishra-15/go-tables"
)

func main() {
	// Basic separator between row groups
	tables.NewFromStrings("Name", "Score", "Grade").
		SetStyle(tables.StyleRounded).
		AddRow("Alice", 95, "A").
		AddRow("Bob", 92, "A").
		AddSeparator().
		AddRow("Charlie", 78, "C").
		AddRow("Dave", 74, "C").
		AddSeparator().
		AddRow("Total", 339, "-").
		Print()

	// Separators work seamlessly with colors and alignment
	tables.NewFromStrings("Service", "Status", "Uptime").
		SetStyle(tables.StyleDouble).
		SetAlign(1, tables.AlignCenter).
		SetAlign(2, tables.AlignRight).
		AddRow("API Gateway", tables.Success("ONLINE"), "99.9%").
		AddRow("Auth Service", tables.Success("ONLINE"), "99.7%").
		AddSeparator().
		AddRow("Cache", tables.Warning("DEGRADED"), "87.2%").
		AddSeparator().
		AddRow("Database", tables.Error("OFFLINE"), "0.0%").
		Print()
}
