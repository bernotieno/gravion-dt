package tree

import (
	"dt/internal"
	"fmt"
	"testing"
)

// MockDataset is a testable implementation of internal.Dataset
type MockDataset struct {
	internal.Dataset
	columnTypes     map[string]internal.ColumnType
	uniqueValuesMap map[string][]string
}

func (m *MockDataset) HasColumn(colName string) bool {
	_, exists := m.columnTypes[colName]
	return exists
}

func (m *MockDataset) GetUniqueValues(colName string) ([]string, error) {
	values, exists := m.uniqueValuesMap[colName]
	if !exists {
		return nil, fmt.Errorf("column %s not found or has no unique values", colName)
	}
	return values, nil
}

func createMockDataset() *MockDataset {
	return &MockDataset{
		columnTypes: map[string]internal.ColumnType{
			"feature1": internal.NumericType,
			"feature2": internal.CategoricalType,
			"target":   internal.CategoricalType,
		},
		uniqueValuesMap: map[string][]string{
			"feature2": {"value1", "value2", "value3"},
			"target":   {"class1", "class2"},
		},
	}
}

func TestNewBuilder(t *testing.T) {
	// Create a mock dataset
	mockData := createMockDataset()

	targetColumn := "target"

	// Test builder creation
	builder := NewBuilder(&mockData.Dataset, targetColumn)

	// Check if builder is properly initialized
	if builder.data != &mockData.Dataset {
		t.Errorf("Expected builder.data to be %v, got %v", &mockData.Dataset, builder.data)
	}

	if builder.targetColumn != targetColumn {
		t.Errorf("Expected builder.targetColumn to be %s, got %s", targetColumn, builder.targetColumn)
	}

	// Check default values
	if builder.minInstancesPerLeaf != 10 {
		t.Errorf("Expected default minInstancesPerLeaf to be 10, got %d", builder.minInstancesPerLeaf)
	}

	if builder.maxDepth != 10 {
		t.Errorf("Expected default maxDepth to be 10, got %d", builder.maxDepth)
	}
}

func TestBuilderSetters(t *testing.T) {
	// Create a mock dataset
	mockData := createMockDataset()

	builder := NewBuilder(&mockData.Dataset, "target")

	// Test SetMinInstancesPerLeaf
	builder.SetMinInstancesPerLeaf(5)
	if builder.minInstancesPerLeaf != 5 {
		t.Errorf("Expected minInstancesPerLeaf to be 5, got %d", builder.minInstancesPerLeaf)
	}

	// Test SetMaxDepth
	builder.SetMaxDepth(15)
	if builder.maxDepth != 15 {
		t.Errorf("Expected maxDepth to be 15, got %d", builder.maxDepth)
	}

	// Test SetNumWorkers
	builder.SetNumWorkers(4)
	if builder.numWorkers != 4 {
		t.Errorf("Expected numWorkers to be 4, got %d", builder.numWorkers)
	}

	// Test method chaining
	builder.SetMinInstancesPerLeaf(20).SetMaxDepth(25).SetNumWorkers(8)
	if builder.minInstancesPerLeaf != 20 || builder.maxDepth != 25 || builder.numWorkers != 8 {
		t.Errorf("Method chaining failed. Expected values 20, 25, 8, got %d, %d, %d",
			builder.minInstancesPerLeaf, builder.maxDepth, builder.numWorkers)
	}
}

// TestBuild requires implementation or mocking of BuildTree
// Here's a revised version that works with our mock approach
type BuilderWithMockTree struct {
	Builder
	mockTreeNode *Node
	mockTreeErr  error
}

func (bm *BuilderWithMockTree) BuildTree(dataset *internal.Dataset, usedFeatures []string, depth int) (*Node, error) {
	return bm.mockTreeNode, bm.mockTreeErr
}
