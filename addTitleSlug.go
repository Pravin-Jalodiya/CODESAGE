package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

// Function to extract title slug from a URL
func extractSlug(url string) string {
	parts := strings.Split(url, "/")
	slug := parts[len(parts)-1]
	return strings.TrimSuffix(strings.ToLower(slug), ".")
}

func main() {
	// Open the CSV file
	inputFile := "/Users/pravin/Desktop/CODESAGE/csv/questions_with_company_tags.csv" // Replace with your CSV file path
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Create a new CSV file to write the updated data
	outputFile := "/Users/pravin/Desktop/CODESAGE/csv/questions_with_slug.csv" // Output CSV file path
	outFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	// Write the header with the new column
	header := append([]string{"Title Slug"}, records[0]...)
	if err := writer.Write(header); err != nil {
		fmt.Println("Error writing header:", err)
		return
	}

	// Process each record and write to the new file
	for _, record := range records[1:] {
		if len(record) > 3 {
			slug := extractSlug(record[3])
			updatedRecord := append([]string{slug}, record...)
			if err := writer.Write(updatedRecord); err != nil {
				fmt.Println("Error writing record:", err)
				return
			}
		}
	}

	fmt.Println("CSV file processed and saved as:", outputFile)
}
