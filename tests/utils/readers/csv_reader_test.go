package readers

import (
	"cli-project/pkg/utils/readers"
	"os"
	"testing"
)

// Helper function to create a temporary CSV file for testing
func createTempCSV(t *testing.T, content string) string {
	file, err := os.CreateTemp("", "*.csv")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		t.Fatalf("failed to write to temporary file: %v", err)
	}

	return file.Name()
}

func TestReadCSV_Success(t *testing.T) {
	csvContent := "name,age,city\nJohn Doe,30,New York\nJane Smith,25,Los Angeles"
	filePath := createTempCSV(t, csvContent)
	defer os.Remove(filePath)

	expected := [][]string{
		{"name", "age", "city"},
		{"John Doe", "30", "New York"},
		{"Jane Smith", "25", "Los Angeles"},
	}

	records, err := readers.ReadCSV(filePath)
	if err != nil {
		t.Fatalf("ReadCSV returned an error: %v", err)
	}

	if len(records) != len(expected) {
		t.Fatalf("expected %d rows, got %d", len(expected), len(records))
	}

	for i, record := range records {
		if len(record) != len(expected[i]) {
			t.Fatalf("expected %d columns in row %d, got %d", len(expected[i]), i, len(record))
		}
		for j, value := range record {
			if value != expected[i][j] {
				t.Errorf("expected %s at row %d column %d, got %s", expected[i][j], i, j, value)
			}
		}
	}
}

func TestReadCSV_FileNotFound(t *testing.T) {
	filePath := "nonexistent_file.csv"

	_, err := readers.ReadCSV(filePath)
	if err == nil {
		t.Fatalf("expected an error for nonexistent file, got nil")
	}
}

func TestReadCSV_FileCannotBeOpened(t *testing.T) {
	filePath := "/root/forbidden_file.csv" // Adjust this path as necessary for your environment

	_, err := readers.ReadCSV(filePath)
	if err == nil {
		t.Fatalf("expected an error for file that cannot be opened, got nil")
	}
}
