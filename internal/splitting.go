package internal

import "math/rand/v2"

// SplitIntoTrainTest splits the dataset into training and testing subsets based on the given trainRatio.
//
// Parameters:
//   - trainRatio: A float64 value representing the proportion of data to be used for training.
//     Must be in the range (0,1). If an invalid value is provided, it defaults to 0.8 (80% train, 20% test).
//
// Returns:
// - *Dataset: The training dataset.
// - *Dataset: The testing dataset.
//
// This function randomizes row indices before splitting to ensure a fair distribution.
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

// createSubset generates a new dataset containing only the rows specified by the given indices.
//
// Parameters:
// - original: Pointer to the original dataset from which the subset is created.
// - indices: A slice of integers representing the row indices to be included in the subset.
//
// Returns:
// - *Dataset: A new dataset containing only the specified rows, with the same column structure as the original dataset.
//
// This function maintains column metadata, including types and additional attributes. It also ensures that missing values
// and metadata for each column are preserved in the new dataset.
func createSubset(original *Dataset, indices []int) *Dataset {
	if len(indices) == 0 {
		return NewDataset()
	}

	subset := NewDataset()
	subset.ColumnTypes = make(map[string]ColumnType)

	// Copy column metadata
	for _, col := range original.Columns {
		newCol := Column{
			Name:     col.Name,
			Type:     col.Type,
			Values:   make([]string, len(indices)),
			Missing:  make([]bool, len(indices)),
			Metadata: make(map[string]interface{}),
		}

		for k, v := range col.Metadata {
			newCol.Metadata[k] = v
		}

		subset.Columns = append(subset.Columns, newCol)
		subset.ColumnMap[col.Name] = len(subset.Columns) - 1
		subset.ColumnTypes[col.Name] = col.Type
	}

	// Copy selected rows
	for i, rowIdx := range indices {
		for colIdx, col := range original.Columns {
			subset.Columns[colIdx].Values[i] = col.Values[rowIdx]
			subset.Columns[colIdx].Missing[i] = col.Missing[rowIdx]
		}
	}

	subset.NumRows = len(indices)
	return subset
}
