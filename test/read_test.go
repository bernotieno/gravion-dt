package test

import (
	"gravion-dt/internal"
	"os"
	"testing"
)

// TestReadCSVFromFile verifies that the ReadCSV function correctly reads data from a CSV file.
func TestReadCSVFromFile(t *testing.T) {
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}

	err = os.Chdir("../")
	if err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	defer os.Chdir(originalDir)

	filePath := "./datasets/dataset.csv"

	// Call ReadCSV function to read the file
	dataset, err := internal.ReadCSV(filePath)
	if err != nil {
		t.Fatalf("ReadCSV failed: %v", err)
	}

	// Ensure that the dataset is not nil
	if dataset == nil {
		t.Fatalf("expected dataset, got nil")
	}

	// Check if there are any rows in the dataset
	if dataset.NumRows == 0 {
		t.Errorf("expected at least one row, got %d", dataset.NumRows)
	}

	// Check if column headers exist
	if len(dataset.Columns) == 0 {
		t.Errorf("expected at least one column, got %d", len(dataset.Columns))
	}
}