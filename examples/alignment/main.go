package main

import (
    tables "github.com/architmishra-15/table"
)

func main() {
    t := tables.NewFromStrings("Left", "Center", "Right").
        SetAlign(0, tables.AlignLeft).
        SetAlign(1, tables.AlignCenter).
        SetAlign(2, tables.AlignRight).
        AddRow("short", "centered text", "12345").
        AddRow("a very long left cell", "mid", "9")

    t.SetStyle(tables.StyleDouble).Print()
}
