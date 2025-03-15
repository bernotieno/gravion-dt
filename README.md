# Fast & Scalable Decision Tree

## Overview
This project implements a high-performance and scalable **C4.5 Decision Tree** classifier in Go. Designed for efficiency, it can handle large datasets with minimal memory overhead and supports parallelization for faster computations.

## Features
- **C4.5 Algorithm**: Implements the improved successor of ID3 with pruning and support for continuous data.
- **Handles Both Categorical & Numerical Data**: Automatically determines optimal splits.
- **Missing Value Handling**: Effectively processes datasets with missing values.
- **Model Export in JSON**: Trained models are stored in a structured JSON format.
- **Command-Line Interface (CLI)**: Simple CLI-based execution for both training and prediction.
- **Optimized for Speed & Scalability**: Supports parallel execution for enhanced performance.

## Project Structure
```sh
fast-decision-tree/
├── datasets/         # Sample training & test datasets
├── internal/         # Core implementation of the decision tree
├── tree/             # Decision tree implementation
├── Makefile          # Build automation script
├── README.md         # Project documentation
├── gitignore         # Git ignore file to exclude unnecessary files
├── go.mod            # Go module dependencies
└── main.go           # CLI entry point
```

## Installation
### Prerequisites
Ensure you have Go installed on your system. You can install Go from [golang.org](https://go.dev/dl/).

### Clone the Repository
```sh
git clone https://github.com/bernotieno/gravion-dt.git
cd gravion-dt
```

### Building the Project
The project uses a `Makefile` for building and running the application efficiently.

To compile the project:
```sh
make
```
This will generate an executable named `dt`.


## Usage
### Training a Decision Tree
Train a decision tree model from a dataset:
```sh
./dt -c train -i <input_data_file.csv> -t <target_column> -o <output_tree.dt>
```
#### Arguments:
- `-c train` : Specifies training mode.
- `-i <input_data_file.csv>` : Path to the input CSV file.
- `-t <target_column>` : Column name containing target labels.
- `-o <output_tree.dt>` : Output file path for the trained model.

#### Example:
```sh
./dt -c train -i datasets/train.csv -t class -o models/model.dt
```

### Making Predictions
Use a trained decision tree model to make predictions:
```sh
./dt -c predict -i <prediction_data_file.csv> -m <model_file.dt> -o <predictions.csv>
```
#### Arguments:
- `-c predict` : Specifies prediction mode.
- `-i <prediction_data_file.csv>` : Path to the input file for prediction.
- `-m <model_file.dt>` : Path to the trained model.
- `-o <predictions.csv>` : Output file for storing predictions.

#### Example:
```sh
./dt -c predict -i datasets/test.csv -m models/model.dt -o predictions.csv
```

## Hyperparameter Tuning Function

This function performs hyperparameter tuning using the `tune` command. It searches for the optimal combination of parameters by evaluating different configurations. The function:

- Defines a search space for hyperparameters.
- Executes multiple training runs with varying hyperparameters.
- Logs performance metrics to identify the best settings.
- Returns the optimal hyperparameters for model training.

This approach automates tuning, improving model accuracy and efficiency.

If you want to train your model with the best parameter, use:
```sh
./dt -c tune -i datasets/train.csv -t class -o models/model.dt
```
## Makefile Commands
The `Makefile` includes the following automation commands:
- **Build the project:** `make build`
- **Run the CLI:** `make run`
- **Clean compiled files:** `make clean`

## Error Handling
- **Missing Values**: Automatically handled by ignoring or imputing missing data.
- **Invalid Inputs**: Error messages provided for incorrect arguments or file formats.
- **Large Datasets**: Optimized to process large datasets efficiently without excessive memory usage.

## Authors
- [Bernard Okumu](https://github.com/bernotieno)
- [Cynthia Oketch](https://github.com/CynthiaOketch)
- [Denil Anyonyi](https://github.com/denilany)
- [Joab Owala](https://github.com/joabowala)
- [Hilary Okello](https://github.com/hilaryokello)

