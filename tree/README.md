# Decision Tree Package Documentation

## Overview

This package provides a comprehensive implementation of a decision tree algorithm for classification tasks. The decision tree is a supervised machine learning model that recursively splits the dataset based on feature values to create a tree-like structure, where each internal node represents a decision based on a feature, and each leaf node represents a class label.

The package is designed to handle both categorical and numeric features, and it includes functionalities for building, training, evaluating, and predicting using the decision tree model. The implementation is optimized for performance, leveraging parallel processing to speed up the tree-building process.

## Key Features

1. **Handling of Categorical and Numeric Features**: The package can handle datasets with both categorical and numeric features, automatically determining the type of each feature and applying the appropriate splitting criteria.

2. **Parallel Tree Building**: The tree-building process is parallelized using goroutines, allowing for faster construction of the decision tree, especially with large datasets.

3. **Customizable Parameters**: Users can customize various parameters such as the maximum depth of the tree and the minimum number of instances required to create a leaf node.

4. **Model Persistence**: The trained decision tree model can be saved to a file and later loaded for making predictions on new data.

5. **Evaluation and Prediction**: The package includes functions for evaluating the model's accuracy on a test dataset and making predictions on new instances.

## Package Structure

The package is organized into several files, each containing related functionality:

- **build.go**: Contains the core logic for recursively building the decision tree.
- **create_leaf.go**: Implements the creation of leaf nodes in the decision tree.
- **dataset_utility.go**: Provides utility functions for handling datasets, such as calculating class counts, determining the majority class, and calculating class confidence.
- **entropy.go**: Implements the calculation of entropy, a measure of impurity in the dataset.
- **gain_ratio.go**: Contains functions for calculating the gain ratio, which is used to determine the best feature to split on.
- **models.go**: Defines the data structures used in the package, including the `Node` and `DecisionTree` types.
- **predict.go**: Implements functions for making predictions using a trained decision tree model.
- **split_selection.go**: Contains the logic for finding the best feature to split on, using parallel processing.
- **train.go**: Implements the training process, including dataset splitting, model building, and evaluation.
- **tree.go**: Provides the main interface for building the decision tree, including setting parameters and initiating the tree-building process.

## Detailed Explanation

### Building the Decision Tree (`build.go`)

The `BuildTree` function is the core of the decision tree algorithm. It recursively builds the tree by selecting the best feature to split on at each node, based on the gain ratio. The function handles both categorical and numeric features, creating different types of nodes accordingly.

- **Stopping Conditions**: The recursion stops when one of the following conditions is met:
  - The maximum depth of the tree is reached.
  - The number of instances in the dataset is below the minimum required for a leaf node.
  - The dataset is almost pure (i.e., one class dominates).
  - No further splits provide a significant gain in information.

- **Node Creation**: Depending on the type of the best feature (categorical or numeric), the function creates either a `CategoricalNode` or a `NumericNode`. For categorical features, the function creates a child node for each unique value of the feature. For numeric features, the function splits the dataset based on a threshold and creates two child nodes.

### Leaf Node Creation (`create_leaf.go`)

The `CreateLeafNode` function creates a leaf node by determining the majority class in the dataset and calculating the confidence of the prediction. The function also calculates the class counts, which are stored in the leaf node for later use during prediction.

### Dataset Utilities (`dataset_utility.go`)

This file contains utility functions for handling datasets:

- **getClassCounts**: Calculates the frequency of each class in the target column.
- **getMajorityClass**: Determines the majority class in the dataset.
- **getClassConfidence**: Calculates the confidence level of a given class.
- **ValidateFeatures**: Checks if all required feature columns are present in the dataset.
- **Predict**: Makes predictions for all instances in the dataset.
- **PredictInstance**: Makes a prediction for a single instance by traversing the decision tree.

### Entropy Calculation (`entropy.go`)

The `CalculateEntropy` function calculates the entropy of the dataset, which is a measure of impurity. Entropy is used to determine the information gain when splitting the dataset on a feature.

### Gain Ratio Calculation (`gain_ratio.go`)

The `calculateGainRatio` function calculates the gain ratio for a given feature, which is used to determine the best feature to split on. The gain ratio is calculated differently for categorical and numeric features:

- **Categorical Features**: The function calculates the split information entropy and the information gain, then divides the information gain by the split information entropy to get the gain ratio.
- **Numeric Features**: The function considers each unique value of the feature as a potential threshold, calculates the information gain for each threshold, and selects the threshold with the highest gain ratio.

### Data Structures (`models.go`)

This file defines the data structures used in the package:

- **Node**: Represents a node in the decision tree. It can be a leaf node, a categorical node, or a numeric node.
- **DecisionTree**: Represents a trained decision tree model, including the root node, target column, feature types, and feature values.
- **Builder**: Responsible for building the decision tree. It includes parameters such as the minimum number of instances per leaf, the maximum depth of the tree, and the number of worker goroutines.

### Prediction (`predict.go`)

The `Predict` function uses a trained decision tree model to make predictions on new data. It loads the model from a file, validates the feature columns in the input data, makes predictions, and saves the results to an output file.

### Split Selection (`split_selection.go`)

The `findBestSplit` function finds the best feature to split on by calculating the gain ratio for each available feature in parallel. The function uses goroutines to speed up the calculation, with the number of concurrent goroutines controlled by a semaphore.

### Training (`train.go`)

The `Train` function trains a decision tree model using the dataset provided in the input file. It splits the dataset into training and test sets, builds the decision tree, evaluates the model's accuracy on the test set, and saves the trained model to a file.

### Tree Building Interface (`tree.go`)

The `Build` function in this file provides the main interface for building the decision tree. It verifies that the target column exists, determines the feature types, gathers unique values for categorical features, and initiates the tree-building process. The file also includes functions for setting parameters such as the minimum number of instances per leaf, the maximum depth of the tree, and the number of worker goroutines.

## Usage Example

### Training a Decision Tree Model

```go
package main

import (
	"fmt"
	"tree"
)

func main() {
	// Train a decision tree model
	err := tree.Train("data.csv", "target_column", "model.json")
	if err != nil {
		fmt.Println("Error training model:", err)
	}
}
```

### Making Predictions

```go
package main

import (
	"fmt"
	"tree"
)

func main() {
	// Make predictions using a trained decision tree model
	err := tree.Predict("test_data.csv", "model.json", "predictions.csv")
	if err != nil {
		fmt.Println("Error making predictions:", err)
	}
}
```

## Conclusion

This package provides a robust and efficient implementation of a decision tree algorithm for classification tasks. It is designed to handle both categorical and numeric features, and it includes functionalities for building, training, evaluating, and predicting using the decision tree model. The package is highly customizable, allowing users to fine-tune various parameters to optimize the performance of the model.
