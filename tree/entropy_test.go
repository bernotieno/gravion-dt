package tree

import (
	"dt/internal"
	"testing"
)

func TestCalculateEntropy(t *testing.T) {
	builder := &Builder{targetColumn: "class"}

	// Case 1: Pure dataset (all same class, entropy should be 0)
	pureData := &internal.Dataset{
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

	entropy, err := builder.CalculateEntropy(pureData)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if entropy != 0 {
		t.Errorf("Expected entropy 0, got %f", entropy)
	}

	// Case 2: Evenly mixed dataset (entropy should be 1 for binary classes)
	mixedData := &internal.Dataset{
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

	entropy, err = builder.CalculateEntropy(mixedData)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	expectedEntropy := 1.0 // log2(0.5) * 2 = -0.5 * 2 = 1
	if entropy != expectedEntropy {
		t.Errorf("Expected entropy %f, got %f", expectedEntropy, entropy)
	}

	// Case 3: Skewed dataset (entropy should be lower but > 0)
	skewedData := &internal.Dataset{
		Columns: []internal.Column{
			{
				Name:    "class",
				Type:    internal.CategoricalType,
				Values:  []string{"A", "A", "A", "B"},
				Missing: []bool{false, false, false, false},
			},
		},
		ColumnMap:   map[string]int{"class": 0},
		NumRows:     4,
		ColumnTypes: map[string]internal.ColumnType{"class": internal.CategoricalType},
	}

	entropy, err = builder.CalculateEntropy(skewedData)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if entropy <= 0 || entropy >= 1 {
		t.Errorf("Expected entropy between 0 and 1, got %f", entropy)
	}
}
