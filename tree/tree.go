package tree

import (
	"dt/internal"
	"fmt"
)

// Build constructs a decision tree from the dataset.
// It verifies that the target column exists, determines feature types,
// gathers unique values for categorical features, and then builds the tree.
func (b *Builder) Build() (*DecisionTree, error) {
	if !b.data.HasColumn(b.targetColumn) {
		return nil, fmt.Errorf("target column '%s' not found in dataset", b.targetColumn)
	}

	// Prepare feature types for the model
	featureTypes := make(map[string]internal.ColumnType)
	for colName, colType := range b.data.ColumnTypes {
		if colName != b.targetColumn {
			featureTypes[colName] = colType
		}
	}

	// Prepare feature values for categorical features
	featureValues := make(map[string][]string)
	for colName, colType := range featureTypes {
		if colType == internal.CategoricalType {
			uniqueValues, err := b.data.GetUniqueValues(colName)
			if err != nil {
				return nil, err
			}
			featureValues[colName] = uniqueValues
		}
	}

	// Build the tree
	root, err := b.BuildTree(b.data, make([]string, 0), 0)
	if err != nil {
		return nil, err
	}

	return &DecisionTree{
		Root:          root,
		TargetColumn:  b.targetColumn,
		FeatureTypes:  featureTypes,
		FeatureValues: featureValues,
	}, nil
}

// SetMinInstancesPerLeaf sets the minimum number of instances per leaf node.
//
// Parameters:
// - minInstances (int): The minimum number of data points required at a leaf node.
func (b *Builder) SetMinInstancesPerLeaf(minInstances int) *Builder {
	b.minInstancesPerLeaf = minInstances
	return b
}

// SetMaxDepth sets the maximum depth of the decision tree.
//
// Parameters:
// - depth (int): The maximum depth of the tree.
func (b *Builder) SetMaxDepth(depth int) *Builder {
	b.maxDepth = depth
	return b
}

// SetNumWorkers sets the number of worker goroutines for parallel tree building.
//
// Parameters:
// - numWorkers (int): The number of worker goroutines to use.
func (b *Builder) SetNumWorkers(numWorkers int) *Builder {
	b.numWorkers = numWorkers
	return b
}
