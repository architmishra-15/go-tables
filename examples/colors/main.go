package main

import (
    tables "github.com/architmishra-15/table"
)

func main() {
    ct := tables.NewFromStrings("Item", "Price", "Status").SetStyle(tables.StyleRounded).
        AddRow(tables.Success("Laptop"), tables.Sprint("$999", tables.Bold), tables.Success("In Stock")).
        AddRow(tables.Warning("Mouse"), tables.Sprint("$25", tables.Bold), tables.Error("Out"))

    ct.Print()
}
