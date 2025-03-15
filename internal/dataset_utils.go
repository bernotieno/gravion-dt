package internal

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

// GetColIndx returns the index of a column with the given name
func (d *Dataset) GetColIndex(name string) (int, error) {
	index, exists := d.ColumnMap[name]
	if !exists {
		return -1, fmt.Errorf("column '%s' not found", name)
	}
	return index, nil
}

// HasColumn checks if the dataset has a column with the given name
func (d *Dataset) HasColumn(name string) bool {
	_, exists := d.ColumnMap[name]
	return exists
}

// GetColType returns the type of a column
func (d *Dataset) GetColType(name string) (ColumnType, error) {
	if !d.HasColumn(name) {
		return 0, fmt.Errorf("column '%s' not found", name)
	}
	return d.ColumnTypes[name], nil
}

// GetColValues returns all values in a column
func (d *Dataset) GetColValues(name string) ([]string, []bool, error) {
	index, err := d.GetColIndex(name)
	if err != nil {
		return nil, nil, err
	}
	return d.Columns[index].Values, d.Columns[index].Missing, nil
}

// GetUniqueValues returns the unique values in a categorical column
func (d *Dataset) GetUniqueValues(columnName string) ([]string, error) {
	idx, err := d.GetColIndex(columnName)
	if err != nil {
		return nil, err
	}

	if d.Columns[idx].Type != CategoricalType {
		return nil, fmt.Errorf("column '%s' is not categorical", columnName)
	}

	uniqueMap := make(map[string]bool)
	for i, val := range d.Columns[idx].Values {
		if !d.Columns[idx].Missing[i] {
			uniqueMap[val] = true
		}
	}

	uniqueValues := make([]string, 0, len(uniqueMap))
	for val := range uniqueMap {
		uniqueValues = append(uniqueValues, val)
	}

	return uniqueValues, nil
}

// GetNumericValue converts a string value to a float64
func GetNumericValue(value string) (float64, error) {
	if value == "" {
		return 0, fmt.Errorf("missing value")
	}
	return strconv.ParseFloat(value, 64)
}

// GetNumericValues returns all numeric values in a column
func (d *Dataset) GetNumericValues(columnName string) ([]float64, []bool, error) {
	idx, err := d.GetColIndex(columnName)
	if err != nil {
		return nil, nil, err
	}

	if d.Columns[idx].Type != NumericType {
		return nil, nil, fmt.Errorf("column '%s' is not numeric", columnName)
	}

	values := make([]float64, len(d.Columns[idx].Values))
	validValues := make([]bool, len(d.Columns[idx].Values))

	for i, val := range d.Columns[idx].Values {
		if d.Columns[idx].Missing[i] {
			validValues[i] = false
			continue
		}

		numVal, err := GetNumericValue(val)
		if err != nil {
			validValues[i] = false
			continue
		}

		values[i] = numVal
		validValues[i] = true
	}

	return values, validValues, nil
}

// GetRow returns all values in a row
func (d *Dataset) GetRow(rowIndex int) (map[string]string, error) {
	if rowIndex < 0 || rowIndex >= d.NumRows {
		return nil, fmt.Errorf("row index %d out of bounds (0-%d)", rowIndex, d.NumRows-1)
	}

	row := make(map[string]string)
	for _, col := range d.Columns {
		if !col.Missing[rowIndex] {
			row[col.Name] = col.Values[rowIndex]
		}
	}

	return row, nil
}

// SplitDataset divides a dataset into two parts based on a categorical attribute value
func (d *Dataset) SplitDataset(attributeName, attributeValue string) (*Dataset, *Dataset, error) {
	attrIndex, err := d.GetColIndex(attributeName)
	if err != nil {
		return nil, nil, err
	}

	matchingIndices := make([]int, 0)
	nonMatchingIndices := make([]int, 0)

	for i := 0; i < d.NumRows; i++ {
		if !d.Columns[attrIndex].Missing[i] && d.Columns[attrIndex].Values[i] == attributeValue {
			matchingIndices = append(matchingIndices, i)
		} else {
			nonMatchingIndices = append(nonMatchingIndices, i)
		}
	}

	matchingDataset := createSubset(d, matchingIndices)
	nonMatchingDataset := createSubset(d, nonMatchingIndices)

	return matchingDataset, nonMatchingDataset, nil
}

// SplitDatasetNumeric divides a dataset into two parts based on a numeric attribute threshold
func (d *Dataset) SplitNumericDataset(attributeName string, threshold float64) (*Dataset, *Dataset, error) {
	attrIndex, err := d.GetColIndex(attributeName)
	if err != nil {
		return nil, nil, err
	}

	if d.Columns[attrIndex].Type != NumericType {
		return nil, nil, fmt.Errorf("column '%s' is not numeric", attributeName)
	}

	lowerIndices := make([]int, 0)
	greaterIndices := make([]int, 0)

	for i := 0; i < d.NumRows; i++ {
		if d.Columns[attrIndex].Missing[i] {
			// Handle missing values by adding to both datasets
			lowerIndices = append(lowerIndices, i)
			greaterIndices = append(greaterIndices, i)
			continue
		}

		val, err := GetNumericValue(d.Columns[attrIndex].Values[i])
		if err != nil {
			continue
		}

		if val <= threshold {
			lowerIndices = append(lowerIndices, i)
		} else {
			greaterIndices = append(greaterIndices, i)
		}
	}

	lowerDataset := createSubset(d, lowerIndices)
	greaterDataset := createSubset(d, greaterIndices)

	return lowerDataset, greaterDataset, nil
}

// SavePredictions writes a list of predictions to a CSV file.
//
// Parameters:
//
//	predictions ([]string): A slice containing the predictions to be saved.
//	filename (string): The name of the CSV file where predictions will be written.
//
// Returns:
//
//	(error): An error if any issue occurs during file creation or writing; otherwise, nil.

func SavePredictions(predictions []string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"prediction"}); err != nil {
		return err
	}

	for _, prediction := range predictions {
		if err := writer.Write([]string{prediction}); err != nil {
			return err
		}
	}

	return nil
}

// Combine merges two datasets together and returns a new combined dataset.
// The datasets must have the same column structure.
//
// Parameters:
//   - other: Another dataset to combine with this one
//
// Returns:
//   - *Dataset: A new dataset containing rows from both datasets
func (d *Dataset) Combine(other *Dataset) *Dataset {
	// Check if the datasets have the same column structure
	if len(d.Columns) != len(other.Columns) {
		panic("datasets must have the same column structure")
	}

	for i, col := range d.Columns {
		if col.Name != other.Columns[i].Name || col.Type != other.Columns[i].Type {
			panic("datasets must have the same column structure")
		}
	}

	// Create a new dataset with the same columns
	combined := NewDataset()

	// Copy columns and metadata
	for _, col := range d.Columns {
		combined.Columns = append(combined.Columns, Column{
			Name:     col.Name,
			Type:     col.Type,
			Values:   make([]string, 0, len(col.Values)+len(other.Columns[combined.ColumnMap[col.Name]].Values)),
			Missing:  make([]bool, 0, len(col.Missing)+len(other.Columns[combined.ColumnMap[col.Name]].Missing)),
			Metadata: col.Metadata,
		})
		combined.ColumnMap[col.Name] = len(combined.Columns) - 1
		combined.ColumnTypes[col.Name] = col.Type
	}

	// Add rows from the first dataset
	for i := 0; i < d.NumRows; i++ {
		for j, col := range d.Columns {
			combined.Columns[j].Values = append(combined.Columns[j].Values, col.Values[i])
			combined.Columns[j].Missing = append(combined.Columns[j].Missing, col.Missing[i])
		}
		combined.NumRows++
	}

	// Add rows from the second dataset
	for i := 0; i < other.NumRows; i++ {
		for j, col := range other.Columns {
			combined.Columns[j].Values = append(combined.Columns[j].Values, col.Values[i])
			combined.Columns[j].Missing = append(combined.Columns[j].Missing, col.Missing[i])
		}
		combined.NumRows++
	}

	return combined
}
