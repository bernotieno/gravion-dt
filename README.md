# Fast & Scalable Decision Tree

## Overview
This project implements a high-performance and scalable Decision Tree (C4.5) classifier in Go. The solution is designed to handle large datasets efficiently with minimal memory overhead and includes support for parallelization.

## Features
- Implements the **C4.5 Decision Tree** algorithm.
- Supports categorical and numerical data.
- Handles missing values and noisy data.
- Outputs trained models in JSON format.
- CLI-based execution for training and prediction.
- Optimized for speed and scalability.

## Installation
To build and use the decision tree tool, follow these steps:

```sh
# Clone the repository
git clone <repo-url>
cd <repo-folder>

# Build the project
go build -o dt .
```

## Usage

### Training a Decision Tree
To train a decision tree model from a dataset:
```sh
dt -c train -i <input_data_file.csv> -t <target_column> -o <output_tree.dt>
```
#### Arguments:
- `-c train` : Specifies training mode.
- `-i <input_data_file.csv>` : Path to the input CSV file.
- `-t <target_column>` : Name of the column containing target labels.
- `-o <output_tree.dt>` : Output file path for the trained model.

#### Example:
```sh
dt -c train -i datasets/train.csv -t class -o model.dt
```

### Making Predictions
To apply the trained decision tree model to new data:
```sh
dt -c predict -i <prediction_data_file.csv> -m <model_file.dt> -o <predictions.csv>
```
#### Arguments:
- `-c predict` : Specifies prediction mode.
- `-i <prediction_data_file.csv>` : Path to the input file for prediction.
- `-m <model_file.dt>` : Path to the trained model.
- `-o <predictions.csv>` : Output file for storing predictions.

#### Example:
```sh
dt -c predict -i datasets/test.csv -m model.dt -o predictions.csv
```

## Error Handling
| Error | Cause | Suggested Fix |
|--------|--------|---------------|
| Missing input file | Incorrect or missing CSV path | Verify the file exists and path is correct |
| Target column not found | Target column missing in dataset | Check column name in CSV |
| Model file not found | Model file path incorrect | Train a model first or verify the path |
| Output path not specified | Missing `-o` argument | Specify output file path |

## Licensing
This project is licensed under the [MIT license](LICENSE)
## References
- [Implementing Decision Trees (C4.5)](https://www.elementsofcomputerscience.com/posts/implementing-decision-trees-c45-algorithm-01/)
- [Decision Tree Classification Explained](https://www.youtube.com/watch?v=ZVR2Way4nwQ)
- [Decision Tree Classification in Python (from scratch!)](https://www.youtube.com/watch?v=sgQAhG5Q7iY)

---
