package tree

import (
	"dt/internal"
	"testing"
)

func TestGetClassCounts(t *testing.T) {
	builder := &Builder{targetColumn: "class"}

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

	classCounts, err := builder.getClassCounts(mockData)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedCounts := map[string]int{"A": 2, "B": 3}
	for class, count := range expectedCounts {
		if classCounts[class] != count {
			t.Errorf("Expected count for class '%s' to be %d, got %d", class, count, classCounts[class])
		}
	}
}

func TestGetMajorityClass(t *testing.T) {
	builder := &Builder{targetColumn: "class"}

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

	majorityClass, err := builder.getMajorityClass(mockData)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if majorityClass != "B" {
		t.Errorf("Expected majority class 'B', got '%s'", majorityClass)
	}
}

func TestGetClassConfidence(t *testing.T) {
	builder := &Builder{targetColumn: "class"}

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

	// Since `getClassConfidence` expects a class name, we pass "B"
	confidence, err := builder.getClassConfidence(mockData, "B")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedConfidence := 3.0 / 5.0
	if confidence != expectedConfidence {
		t.Errorf("Expected confidence %.2f, got %.2f", expectedConfidence, confidence)
	}
}

func TestPredict(t *testing.T) {
	tree := &DecisionTree{
		Root: &Node{
			Type:       LeafNode,
			Prediction: "A",
			ClassCounts: map[string]int{
				"A": 3,
				"B": 2,
			},
		},
	}

	mockData := &internal.Dataset{
		NumRows: 3,
	}

	predictions, err := tree.Predict(mockData)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedPredictions := []string{"A", "A", "A"}
	for i, pred := range predictions {
		if pred != expectedPredictions[i] {
			t.Errorf("Expected prediction '%s', got '%s'", expectedPredictions[i], pred)
		}
	}
}

func TestTraverseTree(t *testing.T) {
	tree := &DecisionTree{
		Root: &Node{
			Type:       LeafNode,
			Prediction: "A",
			ClassCounts: map[string]int{
				"A": 3,
				"B": 2,
			},
		},
	}

	instance := map[string]string{
		"feature1": "value1",
	}

	prediction, err := tree.traverseTree(tree.Root, instance)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if prediction != "A" {
		t.Errorf("Expected prediction 'A', got '%s'", prediction)
	}
}
