package main

import (
    tables "github.com/architmishra-15/table"
)

func main() {
    t := tables.NewFromStrings("Name", "Role", "Location", "Status").
        SetStyle(tables.StyleRounded).
        SetAlign(2, tables.AlignCenter).
        SetAlign(3, tables.AlignRight).
        SetMaxWidth(0, 20).
        AddRow(tables.Success("Archit"), "Developer", "India", tables.Success("OK")).
        AddRow("Miyuki 🌸", "Engineer", "Tokyo", tables.Warning("WARN")).
        AddRow("李雷", "Student", "北京", tables.Error("FAIL"))

    t.Print()
    // show as string
    _ = t.String()
}
