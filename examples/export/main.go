package main

import (
	"fmt"
	tables "github.com/architmishra-15/go-tables"
)

func main() {
	t := tables.NewFromStrings("Name", "Score", "Grade").
		SetAlign(1, tables.AlignRight).
		SetAlign(2, tables.AlignCenter).
		AddRow("Alice", 95, "A").
		AddRow("Bob", 87, "B").
		AddSeparator().
		AddRow(tables.Warning("Charlie"), 74, "C")

	fmt.Println("=== CSV ===")
	fmt.Println(t.ToCSV())

	fmt.Println("=== Markdown ===")
	fmt.Println(t.ToMarkdown())

	fmt.Println("=== HTML ===")
	fmt.Println(t.ToHTML())
}
