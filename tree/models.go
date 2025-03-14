package tree

import (
	"dt/internal"
)

// NodeType represents the type of a decision tree node
type NodeType int

const (
	LeafNode NodeType = iota
	CategoricalNode
	NumericNode
)

// Node represents a node in the decision tree
type Node struct {
	Type           NodeType                       `json:"type"`
	AttributeName  string                         `json:"attribute,omitempty"`
	AttributeValue string                         `json:"value,omitempty"`
	Threshold      float64                        `json:"threshold,omitempty"`
	Prediction     string                         `json:"prediction,omitempty"`
	Confidence     float64                        `json:"confidence,omitempty"`
	Children       map[string]*Node               `json:"children,omitempty"`
	LeftChild      *Node                          `json:"left,omitempty"`
	RightChild     *Node                          `json:"right,omitempty"`
	ClassCounts    map[string]int                 `json:"class_counts,omitempty"`
	Metadata       map[string]any                 `json:"metadata,omitempty"`
	FeatureTypes   map[string]internal.ColumnType `json:"feature_types,omitempty"`
}

// DecisionTree represents a trained decision tree model
type DecisionTree struct {
	Root          *Node                          `json:"root"`
	TargetColumn  string                         `json:"target_column"`
	FeatureTypes  map[string]internal.ColumnType `json:"feature_types"`
	FeatureValues map[string][]string            `json:"feature_values,omitempty"`
}

// Builder is responsible for building a decision tree
type Builder struct {
	data                *internal.Dataset
	targetColumn        string
	minInstancesPerLeaf int
	maxDepth            int
	numWorkers          int
}

// NewBuilder creates a new instance of Builder with the provided dataset and target column.
// It initializes the Builder with default values for minimum instances per leaf, maximum depth,
// and number of worker goroutines.
//
// Parameters:
//   - data: A pointer to an internal.Dataset containing the data to be used for building the model.
//   - targetColumn: A string specifying the name of the target column in the dataset.
//
// Returns:
//   - A pointer to a newly created Builder instance.
func NewBuilder(data *internal.Dataset, targetColumn string) *Builder {
	return &Builder{
		data:                data,
		targetColumn:        targetColumn,
		minInstancesPerLeaf: 2,  // Default minimum instances per leaf
		maxDepth:            20, // Default maximum depth
		numWorkers:          4,  // Default number of worker goroutines
	}
}
