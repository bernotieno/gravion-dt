package tree

import (
	"dt/internal"
	"testing"
)

func TestCalculateCategoricalGainRatio(t *testing.T) {
	builder := &Builder{targetColumn: "class"}

	// Mock dataset for categorical attribute
	mockData := &internal.Dataset{
		Columns: []internal.Column{
			{
				Name:    "class",
				Type:    internal.CategoricalType,
				Values:  []string{"A", "A", "B", "B", "B"},
				Missing: []bool{false, false, false, false, false},
			},
			{
				Name:    "color",
				Type:    internal.CategoricalType,
				Values:  []string{"red", "red", "blue", "blue", "blue"},
				Missing: []bool{false, false, false, false, false},
			},
		},
		ColumnMap:   map[string]int{"class": 0, "color": 1},
		NumRows:     5,
		ColumnTypes: map[string]internal.ColumnType{"class": internal.CategoricalType, "color": internal.CategoricalType},
	}

	// Calculate class entropy
	classEntropy, err := builder.CalculateEntropy(mockData)
	if err != nil {
		t.Fatalf("Failed to calculate class entropy: %v", err)
	}

	// Compute gain ratio for categorical attribute
	_, gainRatio, err := builder.calculateCategoricalGainRatio(mockData, "color", classEntropy)
	if err != nil {
		t.Fatalf("Error in calculating gain ratio: %v", err)
	}

	// Validate gain ratio (should be between 0 and 1)
	if gainRatio < 0 || gainRatio > 1 {
		t.Errorf("Gain ratio out of expected range: got %f", gainRatio)
	}
}

func TestCalculateNumericGainRatio(t *testing.T) {
	builder := &Builder{targetColumn: "class"}

	// Mock dataset for numeric attribute
	mockData := &internal.Dataset{
		Columns: []internal.Column{
			{
				Name:    "class",
				Type:    internal.CategoricalType,
				Values:  []string{"A", "A", "B", "B", "B"},
				Missing: []bool{false, false, false, false, false},
			},
			{
				Name:    "height",
				Type:    internal.NumericType,
				Values:  []string{"1.5", "1.7", "1.8", "2.0", "2.2"},
				Missing: []bool{false, false, false, false, false},
			},
		},
		ColumnMap:   map[string]int{"class": 0, "height": 1},
		NumRows:     5,
		ColumnTypes: map[string]internal.ColumnType{"class": internal.CategoricalType, "height": internal.NumericType},
	}

	// Calculate class entropy
	classEntropy, err := builder.CalculateEntropy(mockData)
	if err != nil {
		t.Fatalf("Failed to calculate class entropy: %v", err)
	}

	// Compute gain ratio for numeric attribute
	_, gainRatio, err := builder.calculateNumericGainRatio(mockData, "height", classEntropy)
	if err != nil {
		t.Fatalf("Error in calculating gain ratio: %v", err)
	}

	// Validate gain ratio
	if gainRatio < 0 || gainRatio > 1 {
		t.Errorf("Gain ratio out of expected range: got %f", gainRatio)
	}
}

func TestCalculateGainRatioUnsupportedType(t *testing.T) {
	builder := &Builder{targetColumn: "class"}

	// Mock dataset with unsupported attribute type
	mockData := &internal.Dataset{
		Columns: []internal.Column{
			{
				Name:    "class",
				Type:    internal.CategoricalType,
				Values:  []string{"A", "A", "B", "B"},
				Missing: []bool{false, false, false, false},
			},
			{
				Name:    "timestamp",
				Type:    internal.TimestampType, // Unsupported type
				Values:  []string{"1625097600", "1625184000", "1625270400", "1625356800"},
				Missing: []bool{false, false, false, false},
			},
		},
		ColumnMap:   map[string]int{"class": 0, "timestamp": 1},
		NumRows:     4,
		ColumnTypes: map[string]internal.ColumnType{"class": internal.CategoricalType, "timestamp": internal.TimestampType},
	}

	// Calculate class entropy
	classEntropy, err := builder.CalculateEntropy(mockData)
	if err != nil {
		t.Fatalf("Failed to calculate class entropy: %v", err)
	}

	// Attempt gain ratio calculation (should return an error)
	_, _, err = builder.calculateGainRatio(mockData, "timestamp", classEntropy)
	if err == nil {
		t.Errorf("Expected error for unsupported attribute type, but got nil")
	}
}

