package tree

import "gravion-dt/internal"

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
