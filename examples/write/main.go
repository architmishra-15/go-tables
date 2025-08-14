package main

import (
    "fmt"
    "os"
    tables "github.com/architmishra-15/table"
)

func main() {
    file, err := os.Create("table_example_output.txt")
    if err != nil {
        fmt.Println("create file error", err)
        return
    }
    defer file.Close()

    t := tables.NewFromStrings("Name", "Val").AddRow("x", 1).AddRow("y", 2)
    written, err := t.WriteTo(file)
    if err != nil {
        fmt.Println("write error", err)
        return
    }
    fmt.Println("bytes written:", written)
}
