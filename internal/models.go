package internal

type ColumnType int

const (
	CategoricalType ColumnType = iota
	NumericType
	DateType
	TimestampType
)

// Column represents a data column in a dataset
type Column struct {
	Name     string
	Type     ColumnType
	Values   []string
	Missing  []bool
	Metadata map[string]any
}

// Dataset represents a tabular dataset with columns
type Dataset struct {
	Columns     []Column
	ColumnMap   map[string]int
	NumRows     int
	ColumnTypes map[string]ColumnType
}

// NewDataset creates a new empty dataset with initialized maps
func NewDataset() *Dataset {
	return &Dataset{
		Columns:     []Column{},
		ColumnMap:   make(map[string]int),
		NumRows:     0,
		ColumnTypes: make(map[string]ColumnType),
	}
}
