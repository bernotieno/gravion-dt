package tree

import (
	"dt/internal"
	"fmt"
)

// getClassCounts calculates the frequency of each class in the target column of the dataset.
// It returns a map where the keys are class labels and the values are the counts of each class.
//
// Parameters:
//
//	data (*internal.Dataset): The dataset from which to extract the target column values.
//
// Returns:
//
//	(map[string]int, error): A map containing the counts of each class and an error if any occurred during processing.
func (b *Builder) getClassCounts(data *internal.Dataset) (map[string]int, error) {
	targetValues, missingTarget, err := data.GetColValues(b.targetColumn)
	if err != nil {
		return nil, err
	}

	counts := make(map[string]int)
	for i, val := range targetValues {
		if !missingTarget[i] {
			counts[val]++
		}
	}

	return counts, nil
}

// getMajorityClass determines the majority class in the given dataset.
// It returns the class with the highest count and an error if any occurs during the process.
//
// Parameters:
//   - data: A pointer to an internal.Dataset containing the data to be analyzed.
//
// Returns:
//   - string: The majority class in the dataset.
//   - error: An error if any issue occurs while getting the class counts.
func (b *Builder) getMajorityClass(data *internal.Dataset) (string, error) {
	classCounts, err := b.getClassCounts(data)
	if err != nil {
		return "", err
	}

	var majorityClass string
	var maxCount int

	for class, count := range classCounts {
		if count > maxCount {
			maxCount = count
			majorityClass = class
		}
	}

	return majorityClass, nil
}

// getClassConfidence calculates the confidence level of a given class within the dataset.
// It returns the ratio of the count of the specified class to the total count of all classes.
//
// Parameters:
//   - data: A pointer to an internal.Dataset containing the data to analyze.
//   - class: A string representing the class for which the confidence level is to be calculated.
//
// Returns:
//   - A float64 representing the confidence level of the specified class.
//   - An error if there is an issue retrieving class counts or if the total count is zero.
func (b *Builder) getClassConfidence(data *internal.Dataset, class string) (float64, error) {
	classCounts, err := b.getClassCounts(data)
	if err != nil {
		return 0, err
	}

	totalCount := 0
	for _, count := range classCounts {
		totalCount += count
	}

	if totalCount == 0 {
		return 0, nil
	}

	return float64(classCounts[class]) / float64(totalCount), nil
}

// ValidateFeatures checks if all required feature columns are present in the dataset
func (dt *DecisionTree) ValidateFeatures(data *internal.Dataset) []string {
	missingColumns := make([]string, 0)
	for featureName := range dt.FeatureTypes {
		if !data.HasColumn(featureName) {
			missingColumns = append(missingColumns, featureName)
		}
	}
	return missingColumns
}

// Predict makes predictions for all instances in the dataset
func (dt *DecisionTree) Predict(data *internal.Dataset) ([]string, error) {
	predictions := make([]string, data.NumRows)

	for i := 0; i < data.NumRows; i++ {
		row, err := data.GetRow(i)
		if err != nil {
			return nil, err
		}

		prediction, err := dt.PredictInstance(row)
		if err != nil {
			return nil, err
		}

		predictions[i] = prediction
	}

	return predictions, nil
}

// PredictInstance makes a prediction for a single instance
func (dt *DecisionTree) PredictInstance(instance map[string]string) (string, error) {
	return dt.traverseTree(dt.Root, instance)
}

// traverseTree recursively traverses the decision tree to make a prediction
func (dt *DecisionTree) traverseTree(node *Node, instance map[string]string) (string, error) {
	if node == nil {
		return "", fmt.Errorf("encountered nil node during prediction")
	}

	// If this is a leaf node, return its prediction
	if node.Type == LeafNode {
		return node.Prediction, nil
	}

	// Get the attribute value from the instance
	attrValue, exists := instance[node.AttributeName]
	if !exists || attrValue == "" {
		// Handle missing value - return majority class at this node
		var maxCount int
		var majorityClass string
		for class, count := range node.ClassCounts {
			if count > maxCount {
				maxCount = count
				majorityClass = class
			}
		}
		return majorityClass, nil
	}

	if node.Type == CategoricalNode {
		// For categorical attributes, find the matching branch
		child, exists := node.Children[attrValue]
		if !exists {
			// If no branch matches, use the majority class at this node
			var maxCount int
			var majorityClass string
			for class, count := range node.ClassCounts {
				if count > maxCount {
					maxCount = count
					majorityClass = class
				}
			}
			return majorityClass, nil
		}

		// Continue traversing with the child node
		return dt.traverseTree(child, instance)
	} else if node.Type == NumericNode {
		// For numeric attributes, compare with the threshold
		numValue, err := internal.GetNumericValue(attrValue)
		if err != nil {
			// If value can't be converted to number, use majority class
			var maxCount int
			var majorityClass string
			for class, count := range node.ClassCounts {
				if count > maxCount {
					maxCount = count
					majorityClass = class
				}
			}
			return majorityClass, nil
		}

		if numValue <= node.Threshold {
			return dt.traverseTree(node.LeftChild, instance)
		} else {
			return dt.traverseTree(node.RightChild, instance)
		}
	}

	return "", fmt.Errorf("unknown node type during prediction")
}
