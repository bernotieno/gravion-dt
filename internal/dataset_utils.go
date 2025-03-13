package internal

import "fmt"

const TypeCategorical ColumnType = iota // 0

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

	if d.Columns[idx].Type != TypeCategorical {
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