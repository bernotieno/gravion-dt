package internal

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

// DetectColumnType determines the type of data contained in a slice of strings.
// It samples up to 10 values from the input slice and attempts to classify them
// as numeric, date, timestamp, or categorical.
//
// The function returns one of the following ColumnType values:
// - NumericType: if at least half of the sampled values can be parsed as floats.
// - DateType: if at least half of the sampled values can be parsed as dates in the format "YYYY-MM-DD".
// - TimestampType: if at least half of the sampled values can be parsed as timestamps in the RFC3339 format.
// - CategoricalType: if none of the above conditions are met.
//
// Parameters:
// - values: A slice of strings representing the data to be classified.
//
// Returns:
// - ColumnType: The determined type of the data.
func DetectColumnType(values []string) ColumnType {
	if len(values) == 0 {
		return CategoricalType
	}

	numericCount := 0
	dateCount := 0
	timestampCount := 0
	sampleSize := min(len(values), 10)

	for i := 0; i < sampleSize; i++ {
		if values[i] == "" {
			continue
		}

		// Numeric
		if _, err := strconv.ParseFloat(values[i], 64); err == nil {
			numericCount++
			continue
		}

		// Date (YYYY-MM-DD)
		if _, err := time.Parse("2006-01-02", values[i]); err == nil {
			dateCount++
			continue
		}

		// Timestamp (YYYY-MM-DD HH:MM:SS)
		if _, err := time.Parse(time.RFC3339, values[i]); err == nil {
			timestampCount++
			continue
		}
	}

	// Determine the most likely type
	if numericCount >= sampleSize/2 {
		return NumericType
	} else if dateCount >= sampleSize/2 {
		return DateType
	} else if timestampCount >= sampleSize/2 {
		return TimestampType
	}

	return CategoricalType
}

// ReadCSV reads a CSV file from the given filepath and returns a Dataset.
// It expects the first row of the CSV to be the header containing column names.
// Each subsequent row is treated as a data record.
//
// The function performs the following steps:
// 1. Opens the CSV file.
// 2. Reads the header row to determine column names.
// 3. Initializes a Dataset with columns based on the header.
// 4. Reads each data row, appending values to the corresponding columns in the Dataset.
// 5. Detects the type of each column and updates the Dataset accordingly.
//
// If the file cannot be opened or read, or if there is a mismatch in the number of fields
// in any row compared to the header, an error is returned.
//
// Parameters:
// - filepath: The path to the CSV file to be read.
//
// Returns:
// - A pointer to the Dataset containing the CSV data.
// - An error if any issues are encountered during file reading or processing.
func ReadCSV(filepath string) (*Dataset, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	newReader := csv.NewReader(file)

	// Get the header for each column
	header, err := newReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to header file: %v", err)
	}

	dataset := NewDataset()
	for _, columnName := range header {
		dataset.Columns = append(dataset.Columns, 
		Column{
			Name: columnName,
			Values: make([]string, 0),
			Missing: make([]bool, 0),
			Metadata: make(map[string]any),
		})
		dataset.ColumnMap[columnName] = len(dataset.Columns) - 1
	}

	// Read data in rows
	rowCount := 0
	for {
		record, err := newReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("failed to read file: %v", err)
		}

		if len(record) != len(header) {
			return nil, fmt.Errorf("row %d has %d fields, expected %d", rowCount+1, len(record), len(header))
		}

		for i, value := range record {
			isMissing := value == ""
			dataset.Columns[i].Values = append(dataset.Columns[i].Values, value)
			dataset.Columns[i].Missing = append(dataset.Columns[i].Missing, isMissing)
		}

		rowCount++
	}

	dataset.NumRows = rowCount

	// Detect column types
	for i, col := range dataset.Columns {
		columnType := DetectColumnType(col.Values)
		dataset.ColumnTypes[col.Name] = columnType
		dataset.Columns[i].Type = columnType
	}

	return dataset, nil
}

// SplitIntoTrainTest splits the dataset into training and testing sets (80% train, 20% test)
func (d *Dataset) SplitIntoTrainTest(trainRatio float64) (*Dataset, *Dataset) {
	if trainRatio <= 0 || trainRatio >= 1 {
		trainRatio = 0.8 // Default to 80-20 split if the ratio is invalid
	}

	trainSize := int(float64(d.NumRows) * trainRatio)
	indices := rand.Perm(d.NumRows) // Randomize row indices

	trainIndices := indices[:trainSize]
	testIndices := indices[trainSize:]

	return createSubset(d, trainIndices), createSubset(d, testIndices)
}