package internal

import (
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
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

const (
	sampleSizeForEstimation = 100   // Number of rows to read for estimating row size
	maxRowsThreshold        = 10000 // Threshold for sampling
	sampleFraction          = 0.1   // Fraction of rows to sample if the file is large
)

// ReadCSV reads a CSV file from the specified filepath and returns a Dataset.
// It first opens the file and retrieves its size. Then, it reads the header
// and estimates the average row size to determine if sampling is needed based
// on the estimated number of rows. If sampling is required, it processes a
// fraction of the rows; otherwise, it processes all rows. After reading the
// rows, it detects the column types and updates the dataset accordingly.
//
// Parameters:
//   - filepath: The path to the CSV file.
//   - isPredict: A boolean indicating whether the CSV file is for prediction.
//
// Returns:
//   - *Dataset: A pointer to the Dataset containing the CSV data.
//   - error: An error if any occurred during the reading or processing of the file.
func ReadCSV(filepath string, isPredict bool) (*Dataset, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %v", err)
	}
	fileSize := fileInfo.Size()

	if fileSize == 0 {
		return nil, fmt.Errorf("the csv file provided is empty")
	}

	newReader := csv.NewReader(file)
	// Initialize dataset

	header, err := newReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %v", err)
	}

	dataset := initializeDataset(header)

	if isPredict {
		// For prediction mode, we need to read all rows
		rowCount, err := processAllRows(newReader, header, dataset)
		if err != nil {
			return nil, fmt.Errorf("failed to process rows in prediction mode: %v", err)
		}

		dataset.NumRows = rowCount

		// Detect column types
		for i, col := range dataset.Columns {
			columnType := DetectColumnType(col.Values)
			dataset.ColumnTypes[col.Name] = columnType
			dataset.Columns[i].Type = columnType
		}

		// Return the populated dataset immediately for prediction
		return dataset, nil
	}

	avgRowSize, err := estimateAverageRowSize(filepath, sampleSizeForEstimation)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate row size: %v", err)
	}
	estimatedRows := int(float64(fileSize) / avgRowSize)

	// Determine if sampling is needed
	sample := estimatedRows > maxRowsThreshold

	var rowCount int
	if sample {
		// Sample rows
		rowCount, err = processSampledRows(filepath, header, dataset, sampleFraction)
		if err != nil {
			return nil, fmt.Errorf("failed to process sampled rows: %v", err)
		}
	} else {
		// Read all rows
		rowCount, err = processAllRows(newReader, header, dataset)
		if err != nil {
			return nil, fmt.Errorf("failed to process all rows: %v", err)
		}
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

// initializeDataset initializes the dataset with columns based on the header.
func initializeDataset(header []string) *Dataset {
	dataset := NewDataset()
	for _, columnName := range header {
		dataset.Columns = append(dataset.Columns,
			Column{
				Name:     columnName,
				Values:   make([]string, 0),
				Missing:  make([]bool, 0),
				Metadata: make(map[string]any),
			})
		dataset.ColumnMap[columnName] = len(dataset.Columns) - 1
	}
	return dataset
}

// processSampledRows processes a sampled subset of rows from the CSV file.
func processSampledRows(filepath string, header []string, dataset *Dataset, sampleFraction float64) (int, error) {
	sampledRows, err := sampleRows(filepath, sampleFraction)
	if err != nil {
		return 0, fmt.Errorf("failed to sample rows: %v", err)
	}

	rowCount := 0
	for _, record := range sampledRows {
		if len(record) != len(header) {
			return 0, fmt.Errorf("row has %d fields, expected %d", len(record), len(header))
		}

		for i, value := range record {
			isMissing := value == ""
			dataset.Columns[i].Values = append(dataset.Columns[i].Values, value)
			dataset.Columns[i].Missing = append(dataset.Columns[i].Missing, isMissing)
		}
		rowCount++
	}

	return rowCount, nil
}

// processAllRows processes all rows from the CSV file.
func processAllRows(reader *csv.Reader, header []string, dataset *Dataset) (int, error) {
	rowCount := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, fmt.Errorf("failed to read file: %v", err)
		}

		if len(record) != len(header) {
			return 0, fmt.Errorf("row %d has %d fields, expected %d", rowCount+1, len(record), len(header))
		}

		for i, value := range record {
			isMissing := value == ""
			dataset.Columns[i].Values = append(dataset.Columns[i].Values, value)
			dataset.Columns[i].Missing = append(dataset.Columns[i].Missing, isMissing)
		}
		rowCount++
	}

	return rowCount, nil
}

// estimateAverageRowSize estimates the average row size by reading a sample of rows.
func estimateAverageRowSize(filepath string, sampleSize int) (float64, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	newReader := csv.NewReader(file)

	// Skip header
	_, err = newReader.Read()
	if err != nil {
		return 0, fmt.Errorf("failed to read header: %v", err)
	}

	var totalSize int64
	for i := 0; i < sampleSize; i++ {
		record, err := newReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, fmt.Errorf("failed to read row: %v", err)
		}

		// Calculate row size as the sum of the lengths of all fields
		for _, field := range record {
			totalSize += int64(len(field))
		}
	}

	// Calculate average row size
	if sampleSize == 0 {
		return 0, fmt.Errorf("no rows sampled")
	}
	return float64(totalSize) / float64(sampleSize), nil
}

// sampleRows randomly samples rows from the CSV file.
func sampleRows(filepath string, fraction float64) ([][]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	newReader := csv.NewReader(file)

	// Skip header
	_, err = newReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %v", err)
	}

	// Read all rows into memory (for simplicity, but not memory-efficient for very large files)
	var allRows [][]string
	for {
		record, err := newReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read row: %v", err)
		}
		allRows = append(allRows, record)
	}

	// Calculate the number of rows to sample
	sampleSize := int(float64(len(allRows)) * fraction)
	if sampleSize < 1 {
		sampleSize = 1
	}

	// Create a local random number generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Shuffle the rows using the local generator
	r.Shuffle(len(allRows), func(i, j int) {
		allRows[i], allRows[j] = allRows[j], allRows[i]
	})

	sampledRows := allRows[:sampleSize]
	return sampledRows, nil
}
