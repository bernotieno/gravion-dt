package internal

import (
	"os"
	"testing"
)

func TestDetectColumnType(t *testing.T) {
	tests := []struct {
		name     string
		values   []string
		expected ColumnType
	}{
		{
			name:     "Empty slice",
			values:   []string{},
			expected: CategoricalType,
		},
		{
			name:     "All numeric values",
			values:   []string{"1.23", "4.56", "7.89"},
			expected: NumericType,
		},
		{
			name:     "All date values",
			values:   []string{"2023-01-01", "2023-02-02", "2023-03-03"},
			expected: DateType,
		},
		{
			name:     "All timestamp values",
			values:   []string{"2023-01-01T12:00:00Z", "2023-02-02T12:00:00Z", "2023-03-03T12:00:00Z"},
			expected: TimestampType,
		},
		{
			name:     "Mixed values",
			values:   []string{"1.23", "2023-01-01", "2023-01-01T12:00:00Z", "category"},
			expected: CategoricalType,
		},
		{
			name:     "Mostly numeric values",
			values:   []string{"1.23", "4.56", "7.89", "category"},
			expected: NumericType,
		},
		{
			name:     "Mostly date values",
			values:   []string{"2023-01-01", "2023-02-02", "2023-03-03", "category"},
			expected: DateType,
		},
		{
			name:     "Mostly timestamp values",
			values:   []string{"2023-01-01T12:00:00Z", "2023-02-02T12:00:00Z", "2023-03-03T12:00:00Z", "category"},
			expected: TimestampType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectColumnType(tt.values)
			if result != tt.expected {
				t.Errorf("DetectColumnType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestReadCSV(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		isPredict bool
		expected  *Dataset
		wantErr   bool
	}{
		{
			name: "Simple CSV with numeric data",
			content: `col1,col2
1.23,4.56
7.89,10.11`,
			isPredict: false,
			expected: &Dataset{
				Columns: []Column{
					{Name: "col1", Values: []string{"1.23", "7.89"}, Missing: []bool{false, false}, Metadata: map[string]any{}},
					{Name: "col2", Values: []string{"4.56", "10.11"}, Missing: []bool{false, false}, Metadata: map[string]any{}},
				},
				ColumnMap:   map[string]int{"col1": 0, "col2": 1},
				ColumnTypes: map[string]ColumnType{"col1": NumericType, "col2": NumericType},
				NumRows:     2,
			},
			wantErr: false,
		},
		{
			name: "CSV with missing values",
			content: `col1,col2
1.23,
,10.11`,
			isPredict: false,
			expected: &Dataset{
				Columns: []Column{
					{Name: "col1", Values: []string{"1.23", ""}, Missing: []bool{false, true}, Metadata: map[string]any{}},
					{Name: "col2", Values: []string{"", "10.11"}, Missing: []bool{true, false}, Metadata: map[string]any{}},
				},
				ColumnMap:   map[string]int{"col1": 0, "col2": 1},
				ColumnTypes: map[string]ColumnType{"col1": NumericType, "col2": NumericType},
				NumRows:     2,
			},
			wantErr: false,
		},
		{
			name: "Empty CSV",
			content: `col1,col2
`,
			isPredict: false,
			expected: &Dataset{
				Columns: []Column{
					{Name: "col1", Values: []string{}, Missing: []bool{}, Metadata: map[string]any{}},
					{Name: "col2", Values: []string{}, Missing: []bool{}, Metadata: map[string]any{}},
				},
				ColumnMap:   map[string]int{"col1": 0, "col2": 1},
				ColumnTypes: map[string]ColumnType{"col1": CategoricalType, "col2": CategoricalType},
				NumRows:     0,
			},
			wantErr: false,
		},
		{
			name:      "Non-existent file",
			content:   "",
			isPredict: false,
			expected:  nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file for testing
			var filepath string
			if tt.content != "" {
				tmpfile, err := os.CreateTemp("", "test.csv")
				if err != nil {
					t.Fatalf("Failed to create temp file: %v", err)
				}
				defer os.Remove(tmpfile.Name()) // Clean up

				if _, err := tmpfile.Write([]byte(tt.content)); err != nil {
					t.Fatalf("Failed to write to temp file: %v", err)
				}
				tmpfile.Close()
				filepath = tmpfile.Name()
			} else {
				filepath = "non_existent_file.csv"
			}

			// Call the function
			dataset, err := ReadCSV(filepath, tt.isPredict)

			// Check for errors
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadCSV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Compare the dataset
			if !tt.wantErr {
				if dataset.NumRows != tt.expected.NumRows {
					t.Errorf("ReadCSV() NumRows = %v, want %v", dataset.NumRows, tt.expected.NumRows)
				}

				for i, col := range dataset.Columns {
					if col.Name != tt.expected.Columns[i].Name {
						t.Errorf("Column name mismatch: got %v, want %v", col.Name, tt.expected.Columns[i].Name)
					}
					if len(col.Values) != len(tt.expected.Columns[i].Values) {
						t.Errorf("Column values length mismatch: got %v, want %v", len(col.Values), len(tt.expected.Columns[i].Values))
					}
					for j, val := range col.Values {
						if val != tt.expected.Columns[i].Values[j] {
							t.Errorf("Column value mismatch: got %v, want %v", val, tt.expected.Columns[i].Values[j])
						}
					}
				}
			}
		})
	}
}
