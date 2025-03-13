package test

import (
	"gravion-dt/internal"
	"reflect"
	"testing"
)

func TestDataset_GetColIndex(t *testing.T) {
	type fields struct {
		Columns     []internal.Column
		ColumnMap   map[string]int
		NumRows     int
		ColumnTypes map[string]internal.ColumnType
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Column exists",
			fields: fields{
				ColumnMap: map[string]int{"age": 0, "name": 1},
			},
			args:    args{name: "age"},
			want:    0,
			wantErr: false,
		},
		{
			name: "Column does not exist",
			fields: fields{
				ColumnMap: map[string]int{"age": 0, "name": 1},
			},
			args:    args{name: "salary"},
			want:    -1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &internal.Dataset{
				Columns:     tt.fields.Columns,
				ColumnMap:   tt.fields.ColumnMap,
				NumRows:     tt.fields.NumRows,
				ColumnTypes: tt.fields.ColumnTypes,
			}
			got, err := d.GetColIndex(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dataset.GetColIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Dataset.GetColIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataset_GetColValues(t *testing.T) {
	type fields struct {
		Columns     []internal.Column
		ColumnMap   map[string]int
		NumRows     int
		ColumnTypes map[string]internal.ColumnType
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		want1   []bool
		wantErr bool
	}{
		{
			name: "Column exists",
			fields: fields{
				Columns: []internal.Column{
					{Name: "age", Values: []string{"25", "30"}, Missing: []bool{false, false}},
				},
				ColumnMap: map[string]int{"age": 0},
			},
			args:    args{name: "age"},
			want:    []string{"25", "30"},
			want1:   []bool{false, false},
			wantErr: false,
		},
		{
			name: "Column does not exist",
			fields: fields{
				ColumnMap: map[string]int{"age": 0},
			},
			args:    args{name: "salary"},
			want:    nil,
			want1:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &internal.Dataset{
				Columns:     tt.fields.Columns,
				ColumnMap:   tt.fields.ColumnMap,
				NumRows:     tt.fields.NumRows,
				ColumnTypes: tt.fields.ColumnTypes,
			}
			got, got1, err := d.GetColValues(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dataset.GetColValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Dataset.GetColValues() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Dataset.GetColValues() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDataset_GetUniqueValues(t *testing.T) {
	type fields struct {
		Columns     []internal.Column
		ColumnMap   map[string]int
		NumRows     int
		ColumnTypes map[string]internal.ColumnType
	}
	type args struct {
		columnName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "Categorical column with unique values",
			fields: fields{
				Columns: []internal.Column{
					{Name: "gender", Values: []string{"Male", "Female", "Male"}, Missing: []bool{false, false, false}, Type: internal.CategoricalType},
				},
				ColumnMap: map[string]int{"gender": 0},
			},
			args:    args{columnName: "gender"},
			want:    []string{"Male", "Female"},
			wantErr: false,
		},
		{
			name: "Non-categorical column",
			fields: fields{
				Columns: []internal.Column{
					{Name: "age", Values: []string{"25", "30"}, Missing: []bool{false, false}, Type: internal.NumericType},
				},
				ColumnMap: map[string]int{"age": 0},
			},
			args:    args{columnName: "age"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &internal.Dataset{
				Columns:     tt.fields.Columns,
				ColumnMap:   tt.fields.ColumnMap,
				NumRows:     tt.fields.NumRows,
				ColumnTypes: tt.fields.ColumnTypes,
			}
			got, err := d.GetUniqueValues(tt.args.columnName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dataset.GetUniqueValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Dataset.GetUniqueValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNumericValue(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name:    "Valid numeric value",
			args:    args{value: "25.5"},
			want:    25.5,
			wantErr: false,
		},
		{
			name:    "Empty value",
			args:    args{value: ""},
			want:    0,
			wantErr: true,
		},
		{
			name:    "Invalid numeric value",
			args:    args{value: "abc"},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := internal.GetNumericValue(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNumericValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetNumericValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataset_GetNumericValues(t *testing.T) {
	type fields struct {
		Columns     []internal.Column
		ColumnMap   map[string]int
		NumRows     int
		ColumnTypes map[string]internal.ColumnType
	}
	type args struct {
		columnName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []float64
		want1   []bool
		wantErr bool
	}{
		{
			name: "Numeric column with valid values",
			fields: fields{
				Columns: []internal.Column{
					{Name: "age", Values: []string{"25", "30"}, Missing: []bool{false, false}, Type: internal.NumericType},
				},
				ColumnMap: map[string]int{"age": 0},
			},
			args:    args{columnName: "age"},
			want:    []float64{25, 30},
			want1:   []bool{true, true},
			wantErr: false,
		},
		{
			name: "Non-numeric column",
			fields: fields{
				Columns: []internal.Column{
					{Name: "gender", Values: []string{"Male", "Female"}, Missing: []bool{false, false}, Type: internal.CategoricalType},
				},
				ColumnMap: map[string]int{"gender": 0},
			},
			args:    args{columnName: "gender"},
			want:    nil,
			want1:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &internal.Dataset{
				Columns:     tt.fields.Columns,
				ColumnMap:   tt.fields.ColumnMap,
				NumRows:     tt.fields.NumRows,
				ColumnTypes: tt.fields.ColumnTypes,
			}
			got, got1, err := d.GetNumericValues(tt.args.columnName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dataset.GetNumericValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Dataset.GetNumericValues() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Dataset.GetNumericValues() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDataset_GetRow(t *testing.T) {
	type fields struct {
		Columns     []internal.Column
		ColumnMap   map[string]int
		NumRows     int
		ColumnTypes map[string]internal.ColumnType
	}
	type args struct {
		rowIndex int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "Valid row index",
			fields: fields{
				Columns: []internal.Column{
					{Name: "age", Values: []string{"25", "30"}, Missing: []bool{false, false}},
					{Name: "name", Values: []string{"John", "Jane"}, Missing: []bool{false, false}},
				},
				ColumnMap: map[string]int{"age": 0, "name": 1},
				NumRows:   2,
			},
			args:    args{rowIndex: 0},
			want:    map[string]string{"age": "25", "name": "John"},
			wantErr: false,
		},
		{
			name: "Row index out of bounds",
			fields: fields{
				NumRows: 2,
			},
			args:    args{rowIndex: 2},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &internal.Dataset{
				Columns:     tt.fields.Columns,
				ColumnMap:   tt.fields.ColumnMap,
				NumRows:     tt.fields.NumRows,
				ColumnTypes: tt.fields.ColumnTypes,
			}
			got, err := d.GetRow(tt.args.rowIndex)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dataset.GetRow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Dataset.GetRow() = %v, want %v", got, tt.want)
			}
		})
	}
}
