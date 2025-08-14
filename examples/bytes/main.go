package main

import (
    "fmt"
    tables "github.com/architmishra-15/table"
)

func main() {
    b := tables.New([]byte("ID"), []byte("Value")).SetStyle(tables.StyleSingle)
    for i := 1; i <= 3; i++ {
        id := []byte(fmt.Sprintf("%d", i))
        val := []byte(fmt.Sprintf("Item_%d", i))
        b.AddRowBytes(id, val)
    }
    b.Print()
}
