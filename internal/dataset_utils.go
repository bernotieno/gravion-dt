package internal

import (
	"fmt"
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
