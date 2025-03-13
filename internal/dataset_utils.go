package internal

import "fmt"

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
	return d.Columns[index].Values, d.Columns[idx].Missing, nil
}
