# Fast & Scalable Decision Tree

## Overview
This project implements a high-performance and scalable Decision Tree (C4.5) classifier in Go. The solution is designed to handle large datasets efficiently with minimal memory overhead and includes support for parallelization.

## Features
- Implements the **C4.5 Decision Tree** algorithm.
- Supports categorical and numerical data.
- Handles missing values.
- Outputs trained models in JSON format.
- CLI-based execution for training and prediction.
- Optimized for speed and scalability.

## Installation
To build and use the decision tree tool, follow these steps:

```sh
# Clone the repository
git clone https://github.com/bernotieno/gravion-dt.git


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



## Licensing
This project is licensed under the [MIT license](LICENSE)

## Authors
[Bernard Okumu](https://github.com/bernotieno)  
[Cynthia Oketch](https://github.com/CynthiaOketch)  
[Denil Anyonyi](https://github.com/denilany)  
[Joab Owala](https://github.com/joabowala)  
[Hilary Okello](https://github.com/hilaryokello)

