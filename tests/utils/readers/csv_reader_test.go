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

func TestReadCSV(t *testing.T) {
	// Create a temporary CSV file for testing
	file, err := os.CreateTemp("", "test.csv")
	if err != nil {
		t.Fatalf("Error creating temp file: %v", err)
	}
	defer os.Remove(file.Name())

	file.WriteString("col1,col2,col3\nval1,val2,val3\n")
	file.Close()

	records, err := readers.ReadCSV(file.Name())
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if len(records) != 2 {
		t.Errorf("Expected 2 records, but got %d", len(records))
	}

	_, err = readers.ReadCSV("invalid_path.csv")
	if err == nil {
		t.Error("Expected an error for invalid path")
	}
}
