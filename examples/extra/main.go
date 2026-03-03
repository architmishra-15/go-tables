// examples/extra/main.go

package main

import (
    "fmt"
    tables "github.com/architmishra-15/go-tables"
)

func main() {
    // Sorting — numeric sort descending, separators are stripped
    fmt.Println("=== Sort by Score (descending) ===")
    tables.NewFromStrings("Name", "Score", "Grade").
        SetStyle(tables.StyleRounded).
        SetHeaderColor(tables.NewColor().WithFg(tables.FgCyan).WithStyle(tables.Bold)).
        AddRow("Charlie", 74, "C").
        AddRow("Alice", 95, "A").
        AddRow("Dave", 88, "B").
        AddRow("Bob", 88, "B"). // equal to Dave — stable sort preserves order
        SortByColumn(1, false).
        Print()

    // Footer row
    fmt.Println("=== Footer ===")
    tables.NewFromStrings("Item", "Qty", "Price").
        SetStyle(tables.StyleDouble).
        SetHeaderColor(tables.NewColor().WithStyle(tables.Bold)).
        SetFooterColor(tables.NewColor().WithStyle(tables.Bold)).
        SetAlign(1, tables.AlignRight).
        SetAlign(2, tables.AlignRight).
        AddRow("Widget A", 3, "$9.00").
        AddRow("Widget B", 1, "$24.99").
        AddRow("Widget C", 5, "$3.50").
        SetFooter("Total", 9, "$51.99").
        Print()

    // Sort then footer — footer is not sorted, stays pinned at the bottom
    fmt.Println("=== Sort + Footer ===")
    tables.NewFromStrings("Country", "Population").
        SetStyle(tables.StyleSingle).
        SetAlign(1, tables.AlignRight).
        SetFooterColor(tables.NewColor().WithFg(tables.FgYellow).WithStyle(tables.Bold)).
        AddRow("Brazil", 215_000_000).
        AddRow("India", 1_400_000_000).
        AddRow("USA", 331_000_000).
        AddRow("China", 1_410_000_000).
        SortByColumn(1, false).
        SetFooter("Total", 3_356_000_000).
        Print()
}
