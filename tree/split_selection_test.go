package tree

import (
	"dt/internal"
	"testing"
)

func TestFindBestSplit(t *testing.T) {
	builder := &Builder{targetColumn: "class", numWorkers: 2}

	// Mock dataset with categorical & numeric attributes
	mockData := &internal.Dataset{
		Columns: []internal.Column{
			{
				Name:   "class",
				Type:   internal.CategoricalType,
				Values: []string{"A", "A", "B", "B", "B"},
				Missing: []bool{false, false, false, false, false},
			},
			{
				Name:   "color",
				Type:   internal.CategoricalType,
				Values: []string{"red", "red", "red", "blue", "blue"},
				Missing: []bool{false, false, false, false, false},
			},
			{
				Name:   "height",
				Type:   internal.NumericType,
				Values: []string{"1.5", "1.7", "1.8", "2.0", "2.2"},
				Missing: []bool{false, false, false, false, false},
			},
		},
		ColumnMap:   map[string]int{"class": 0, "color": 1, "height": 2},
		NumRows:     5,
		ColumnTypes: map[string]internal.ColumnType{"class": internal.CategoricalType, "color": internal.CategoricalType, "height": internal.NumericType},
	}

	// Case 1: Find the best attribute to split on (expecting "height")
	bestAttr, _, bestGainRatio, err := builder.findBestSplit(mockData, []string{})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if bestAttr != "height" {
		t.Errorf("Expected 'height' as the best split attribute, got '%s'", bestAttr)
	}
	if bestGainRatio <= 0 {
		t.Errorf("Expected positive gain ratio, got %f", bestGainRatio)
	}

	// Case 2: All data belongs to the same class (no split needed)
	pureData := &internal.Dataset{
		Columns: []internal.Column{
			{
				Name:   "class",
				Type:   internal.CategoricalType,
				Values: []string{"A", "A", "A", "A"},
				Missing: []bool{false, false, false, false},
			},
			{
				Name:   "color",
				Type:   internal.CategoricalType,
				Values: []string{"red", "blue", "green", "yellow"},
				Missing: []bool{false, false, false, false},
			},
		},
		ColumnMap:   map[string]int{"class": 0, "color": 1},
		NumRows:     4,
		ColumnTypes: map[string]internal.ColumnType{"class": internal.CategoricalType, "color": internal.CategoricalType},
	}

	bestAttr, _, bestGainRatio, err = builder.findBestSplit(pureData, []string{})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if bestAttr != "" || bestGainRatio != 0 {
		t.Errorf("Expected no split needed, got: attr=%s, gainRatio=%f", bestAttr, bestGainRatio)
	}
}
