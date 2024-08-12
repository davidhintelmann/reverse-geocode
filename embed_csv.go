//go:generate go run embed_csv.go

package main

import (
	_ "embed"
	"fmt"
	"os"
)

//go:embed geonames.csv
var csvData string

func main() {
	output := "embedded_csv.go"
	data, err := os.ReadFile("geonames.csv")
	if err != nil {
		fmt.Println("Error reading geonames.csv:", err)
		os.Exit(1)
	}

	// Escape the CSV data as a Go string literal
	escapedData := fmt.Sprintf("%q", data)

	// Generate the Go code that embeds the CSV data
	code := fmt.Sprintf(`package main

var csvdata string = %s`, escapedData)

	// var csvData string = %s`, escapedData)

	err = os.WriteFile(output, []byte(code), 0644)
	if err != nil {
		fmt.Println("Error writing embedded CSV file:", err)
		os.Exit(1)
	}
}
