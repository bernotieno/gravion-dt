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
