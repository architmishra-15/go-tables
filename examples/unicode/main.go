package main

import (
    tables "github.com/architmishra-15/table"
)

func main() {
    u := tables.NewFromStrings("Lang", "Greeting", "Flag").SetStyle(tables.StyleRounded).
        AddRow("English", "Hello", "🇺🇸").
        AddRow("Japanese", "こんにちは", "🇯🇵").
        AddRow("Chinese", "你好", "🇨🇳")

    u.Print()
}
