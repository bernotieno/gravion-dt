package tree

import (
	"dt/internal"
	"testing"
)

func TestCreateLeafNode(t *testing.T) {
	builder := &Builder{targetColumn: "class"}

	// Case 1: Standard dataset with a majority class
	mockData := &internal.Dataset{
		Columns: []internal.Column{
			{
				Name:    "class",
				Type:    internal.CategoricalType,
				Values:  []string{"A", "A", "B", "B", "B"},
				Missing: []bool{false, false, false, false, false},
			},
		},
		ColumnMap:   map[string]int{"class": 0},
		NumRows:     5,
		ColumnTypes: map[string]internal.ColumnType{"class": internal.CategoricalType},
	}

	node, err := builder.CreateLeafNode(mockData)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Majority class should be "B"
	if node.Prediction != "B" {
		t.Errorf("Expected majority class 'B', got '%s'", node.Prediction)
	}

	// Confidence should be 3/5 = 0.6
	expectedConfidence := 0.6
	if node.Confidence != expectedConfidence {
		t.Errorf("Expected confidence %.2f, got %.2f", expectedConfidence, node.Confidence)
	}

	// Class counts should be correct
	expectedCounts := map[string]int{"A": 2, "B": 3}
	for class, count := range expectedCounts {
		if node.ClassCounts[class] != count {
			t.Errorf("Expected count for class '%s' to be %d, got %d", class, count, node.ClassCounts[class])
		}
	}
}

func TestCreateLeafNode_EqualDistribution(t *testing.T) {
	builder := &Builder{targetColumn: "class"}

	// Case 2: Dataset with equal distribution
	mockData := &internal.Dataset{
		Columns: []internal.Column{
			{
				Name:    "class",
				Type:    internal.CategoricalType,
				Values:  []string{"A", "B", "A", "B"},
				Missing: []bool{false, false, false, false},
			},
		},
		ColumnMap:   map[string]int{"class": 0},
		NumRows:     4,
		ColumnTypes: map[string]internal.ColumnType{"class": internal.CategoricalType},
	}

	node, err := builder.CreateLeafNode(mockData)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Since distribution is equal, the function may choose either class
	if node.Prediction != "A" && node.Prediction != "B" {
		t.Errorf("Expected majority class 'A' or 'B', got '%s'", node.Prediction)
	}

	// Confidence should be 2/4 = 0.5
	expectedConfidence := 0.5
	if node.Confidence != expectedConfidence {
		t.Errorf("Expected confidence %.2f, got %.2f", expectedConfidence, node.Confidence)
	}

	// Class counts should be correct
	expectedCounts := map[string]int{"A": 2, "B": 2}
	for class, count := range expectedCounts {
		if node.ClassCounts[class] != count {
			t.Errorf("Expected count for class '%s' to be %d, got %d", class, count, node.ClassCounts[class])
		}
	}
}

func TestCreateLeafNode_PureClass(t *testing.T) {
	builder := &Builder{targetColumn: "class"}

	// Case 3: Pure class dataset (all values are the same)
	mockData := &internal.Dataset{
		Columns: []internal.Column{
			{
				Name:    "class",
				Type:    internal.CategoricalType,
				Values:  []string{"A", "A", "A", "A"},
				Missing: []bool{false, false, false, false},
			},
		},
		ColumnMap:   map[string]int{"class": 0},
		NumRows:     4,
		ColumnTypes: map[string]internal.ColumnType{"class": internal.CategoricalType},
	}

	node, err := builder.CreateLeafNode(mockData)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// The only class should be "A"
	if node.Prediction != "A" {
		t.Errorf("Expected majority class 'A', got '%s'", node.Prediction)
	}

	// Confidence should be 100% (1.0)
	expectedConfidence := 1.0
	if node.Confidence != expectedConfidence {
		t.Errorf("Expected confidence %.2f, got %.2f", expectedConfidence, node.Confidence)
	}

	// Class counts should be correct
	expectedCounts := map[string]int{"A": 4}
	for class, count := range expectedCounts {
		if node.ClassCounts[class] != count {
			t.Errorf("Expected count for class '%s' to be %d, got %d", class, count, node.ClassCounts[class])
		}
	}
}

func TestCreateLeafNode_EmptyDataset(t *testing.T) {
	builder := &Builder{targetColumn: "class"}

	// Case 4: Empty dataset
	mockData := &internal.Dataset{
		Columns:     []internal.Column{},
		ColumnMap:   map[string]int{},
		NumRows:     0,
		ColumnTypes: map[string]internal.ColumnType{},
	}

	node, err := builder.CreateLeafNode(mockData)
	if err == nil {
		t.Errorf("Expected error for empty dataset, but got nil")
	}

	if node != nil {
		t.Errorf("Expected nil node for empty dataset, but got %+v", node)
	}
}
