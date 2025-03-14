# Internal Package

This package provides tools for reading, processing, and splitting CSV datasets for machine learning model training. It includes functions to handle data extraction, structuring, and dataset preparation.

## Package Structure

The package consists of the following files:

1. **`read.go`** - Reads data from a CSV file and stores it in a struct.
2. **`models.go`** - Defines the structs and constants used in the package.
3. **`dataset_utils.go`** - Contains utility functions to enhance data extraction.
4. **`splitting.go`** - Splits the dataset into training and test sets.

## Installation

To use this package, clone the repository and import it in your Go project.

```sh
# Clone the repository
git clone https://github.com/bernotieno/dt.git
```

## Usage

### Reading Data
Use the `read.go` functions to load CSV data into structs.
```go
import "dt/internal"
data, err := internal.ReadCSV("dataset")
if err != nil {
    log.Fatal(err)
}
```

### Data Processing Utilities
Use the functions in `dataset_utils.go` to clean and preprocess data.

### Splitting Dataset
Use `splitting.go` to divide the dataset into training and test sets.
```go
trainData, testData := internal.SplitDataset(data, 0.8)
```

