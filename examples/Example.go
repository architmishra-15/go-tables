// main.go

package tables

import (
	"fmt"
	"os"
	"strconv"
	// "github.com/architmishra-15/table"
)

func main() {
	// Example 1: Basic table with different styles
	fmt.Println("=== Example 1: Different Table Styles ===")

	// Create a basic table
	t := NewFromStrings("Name", "Age", "City", "Status").
		AddRow("Alice Johnson", 25, "New York", "Active").
		AddRow("Bob Smith", 30, "Los Angeles", "Inactive").
		AddRow("Charlie Brown", 35, "Chicago", "Active")

	// Show different styles
	styles := []struct {
		name  string
		style Style
	}{
		{"Single", StyleSingle},
		{"Double", StyleDouble},
		{"Rounded", StyleRounded},
		{"ASCII", StyleASCII},
		{"None", StyleNone},
	}

	for _, s := range styles {
		fmt.Printf("--- %s Style ---\n", s.name)
		t.SetStyle(s.style).Print()
		fmt.Println()
	}

	// Example 2: Colored table with ANSI sequences
	fmt.Println("=== Example 2: Colored Table ===")

	colorTable := NewFromStrings("Product", "Price", "Stock", "Status").
		SetStyle(StyleRounded).
		AddRow(
			Success("Laptop"),
			Sprint("$999", Bold),
			Info("15"),
			Success("✓ In Stock"),
		).
		AddRow(
			Warning("Mouse"),
			Sprint("$25", Bold),
			Error("0"),
			Error("✗ Out of Stock"),
		).
		AddRow(
			Info("Keyboard"),
			Sprint("$75", Bold),
			Success("8"),
			Success("✓ In Stock"),
		)

	colorTable.Print()
	fmt.Println()

	// Example 3: Alignment demonstration
	fmt.Println("=== Example 3: Column Alignment ===")

	alignTable := NewFromStrings("Left", "Center", "Right", "Mixed").
		SetAlign(0, AlignLeft).
		SetAlign(1, AlignCenter).
		SetAlign(2, AlignRight).
		SetAlign(3, AlignLeft).
		SetStyle(StyleDouble).
		AddRow("Short", "Medium Text", "A Very Long Text", "Normal").
		AddRow("Very Long Text Here", "Mid", "Short", "Another").
		AddRow("Text", "Another Medium", "Mid", "Last Row")

	alignTable.Print()
	fmt.Println()

	// Example 4: Performance test with byte inputs
	fmt.Println("=== Example 4: Performance Test (Byte-first) ===")

	byteTable := New(
		[]byte("ID"),
		[]byte("Data"),
		[]byte("Value"),
	).SetStyle(StyleSingle)

	// Add rows using byte slices for maximum performance
	for i := 0; i < 5; i++ {
		id := []byte(strconv.Itoa(i + 1))
		data := []byte(fmt.Sprintf("Item_%d", i+1))
		value := []byte(fmt.Sprintf("%.2f", float64(i+1)*1.5))

		byteTable.AddRowBytes(id, data, value)
	}

	byteTable.Print()
	fmt.Println()

	// Example 5: Unicode and wide characters
	fmt.Println("=== Example 5: Unicode Support ===")

	unicodeTable := NewFromStrings("Language", "Greeting", "Flag").
		SetStyle(StyleRounded).
		AddRow("English", "Hello World", "🇺🇸").
		AddRow("Japanese", "こんにちは世界", "🗾").
		AddRow("Chinese", "你好世界", "🇨🇳").
		AddRow("Korean", "안녕하세요 세계", "🇰🇷").
		AddRow("Arabic", "مرحبا بالعالم", "🇸🇦")

	unicodeTable.Print()
	fmt.Println()

	// Example 6: Getting table as string instead of printing
	fmt.Println("=== Example 6: Table as String ===")

	stringTable := NewFromStrings("Item", "Count").
		SetStyle(StyleASCII).
		AddRow("Apples", 42).
		AddRow("Bananas", 17).
		AddRow("Oranges", 8)

	tableString := stringTable.String()

	fmt.Println("Table stored in variable:")
	fmt.Print(tableString)
	fmt.Printf("Table string length: %d characters\n\n", len(tableString))

	// Example 7: Writing to file
	fmt.Println("=== Example 7: Writing to File ===")

	file, err := os.Create("table_output.txt")
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	fileTable := NewFromStrings("Filename", "Size", "Modified").
		SetStyle(StyleSingle).
		AddRow("document.pdf", "2.1MB", "2025-01-15").
		AddRow("image.jpg", "5.8MB", "2025-01-14").
		AddRow("data.csv", "102KB", "2025-01-13")

	// Write table to file
	written, err := fileTable.WriteTo(file)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}

	fmt.Printf("Successfully wrote %d bytes to table_output.txt\n", written)

	// Example 8: Advanced customization
	fmt.Println("\n=== Example 8: Advanced Customization ===")

	advanced := NewFromStrings("Component", "Status", "Uptime", "Load").
		SetStyle(StyleDouble).
		SetAlign(1, AlignCenter).
		SetAlign(2, AlignRight).
		SetAlign(3, AlignRight).
		SetMaxWidth(0, 15). // Limit first column width
		AddRow(
			"Authentication Service",
			Success("ONLINE"),
			"99.9%",
			"12%",
		).
		AddRow(
			"Database Connection Pool",
			Warning("DEGRADED"),
			"87.2%",
			"89%",
		).
		AddRow(
			"Cache Layer",
			Error("OFFLINE"),
			"0.0%",
			"0%",
		)

	advanced.Print()

	fmt.Println("\n=== Library Demo Complete! ===")

	// Fast byte-first approach
	New([]byte("Name"), []byte("Age")).
		AddRowBytes([]byte("Alice"), []byte("25")).
		SetStyle(StyleRounded).
		Print()

	// Convenient string approach
	NewFromStrings("Name", "Age").
		AddRow("Alice", 25).
		SetStyle(StyleDouble).
		SetAlign(1, AlignRight).
		Print()

	// Colored output
	NewFromStrings("Status", "Message").
		AddRow(Success("OK"), "All systems running").
		AddRow(Error("FAIL"), "Connection timeout").
		Print()
}
