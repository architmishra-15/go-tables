package main

import (
    "fmt"
    tables "github.com/architmishra-15/table"
)

func main() {
    fmt.Println("=== Styles Demo ===")
    t := tables.NewFromStrings("Name", "Age", "City").
        AddRow("Alice", 30, "New York").
        AddRow("Bob", 25, "Los Angeles").
        AddRow("Charlie", 22, "Delhi")

    styles := []struct {
        name string
        sty  tables.Style
    }{
        {"Single", tables.StyleSingle},
        {"Double", tables.StyleDouble},
        {"Rounded", tables.StyleRounded},
        {"ASCII", tables.StyleASCII},
        {"None", tables.StyleNone},
    }

    for _, s := range styles {
        fmt.Printf("--- %s ---\n", s.name)
        t.SetStyle(s.sty).Print()
        fmt.Println()
    }
}
