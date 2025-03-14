package tree

import (
	"fmt"
	"math"
	"sort"

	"dt/internal"
)

// calculateGainRatio calculates the gain ratio for a given attribute in the dataset.
// It determines the type of the attribute (categorical or numeric) and calls the appropriate
// method to calculate the gain ratio based on the attribute type.
//
// Parameters:
//   - data: A pointer to the internal.Dataset containing the data.
//   - attrName: The name of the attribute for which the gain ratio is to be calculated.
//   - classEntropy: The entropy of the class attribute.
//
// Returns:
//   - An interface{} representing the result of the gain ratio calculation.
//   - A float64 representing the gain ratio value.
//   - An error if the attribute type is unsupported or if there is an issue retrieving the attribute type.
func (b *Builder) calculateGainRatio(data *internal.Dataset, attrName string, classEntropy float64) (interface{}, float64, error) {
	attrType, err := data.GetColType(attrName)
	if err != nil {
		return nil, 0, err
	}

	if attrType == internal.CategoricalType {
		return b.calculateCategoricalGainRatio(data, attrName, classEntropy)
	} else if attrType == internal.NumericType {
		return b.calculateNumericGainRatio(data, attrName, classEntropy)
	}

	// Other types - not handled
	return nil, 0, fmt.Errorf("unsupported attribute type for '%s'", attrName)
}

// calculateCategoricalGainRatio calculates the gain ratio for a categorical attribute in a dataset.
// It first computes the split information entropy of the attribute, then calculates the information gain
// by subtracting the weighted entropy of each subset (based on the attribute values) from the class entropy.
// Finally, it computes the gain ratio by dividing the information gain by the split information entropy.
//
// Parameters:
// - data: The dataset containing the attribute and class information.
// - attrName: The name of the attribute for which the gain ratio is to be calculated.
// - classEntropy: The entropy of the class variable.
//
// Returns:
// - An interface{} (currently nil) which can be used to return additional information if needed.
// - A float64 representing the gain ratio of the attribute.
// - An error if any issues occur during the calculation.
func (b *Builder) calculateCategoricalGainRatio(data *internal.Dataset, attrName string, classEntropy float64) (interface{}, float64, error) {
	uniqueValues, err := data.GetUniqueValues(attrName)
	if err != nil {
		return nil, 0, err
	}

	// Calculate split info (entropy of the attribute)
	attrValues, attrMissing, err := data.GetColValues(attrName)
	if err != nil {
		return nil, 0, err
	}

	// Count occurrences of each value
	valueCounts := make(map[string]int)
	validCount := 0
	for i, val := range attrValues {
		if !attrMissing[i] {
			valueCounts[val]++
			validCount++
		}
	}

	// Calculate split information entropy
	var splitInfo float64
	for _, count := range valueCounts {
		if count > 0 {
			proportion := float64(count) / float64(validCount)
			splitInfo -= proportion * math.Log2(proportion)
		}
	}

	if splitInfo == 0 {
		return nil, 0, nil
	}

	var infoGain float64 = classEntropy
	totalValidInstances := 0

	for i := range attrValues {
		if !attrMissing[i] {
			totalValidInstances++
		}
	}

	// Calculate entropy for each attribute value
	for _, value := range uniqueValues {
		// Get subset of data with this attribute value
		subset, _, err := data.SplitDataset(attrName, value)
		if err != nil {
			return nil, 0, err
		}

		if subset.NumRows > 0 {
			// Calculate entropy of this subset
			subsetEntropy, err := b.CalculateEntropy(subset)
			if err != nil {
				return nil, 0, err
			}

			// Weighted subtraction from information gain
			weight := float64(subset.NumRows) / float64(totalValidInstances)
			infoGain -= weight * subsetEntropy
		}
	}

	gainRatio := infoGain / splitInfo

	return nil, gainRatio, nil
}

// calculateNumericGainRatio calculates the best threshold and gain ratio for a numeric attribute
// in a dataset to split the data for decision tree construction.
//
// Parameters:
// - data: A pointer to the internal.Dataset containing the data.
// - attrName: The name of the numeric attribute to evaluate.
// - classEntropy: The entropy of the class distribution in the dataset.
//
// Returns:
// - The best threshold value for splitting the numeric attribute.
// - The best gain ratio achieved by the split.
// - An error if any occurs during the calculation.
//
// The function first retrieves the numeric values and their validity for the given attribute.
// It then sorts these values and identifies unique values to consider as potential thresholds.
// For each potential threshold, it splits the dataset and calculates the entropy of the resulting subsets.
// It computes the information gain and split information to determine the gain ratio.
// The threshold with the highest gain ratio is returned.
func (b *Builder) calculateNumericGainRatio(data *internal.Dataset, attrName string, classEntropy float64) (interface{}, float64, error) {
	numValues, validValues, err := data.GetNumericValues(attrName)
	if err != nil {
		return nil, 0, err
	}

	_, targetMissing, err := data.GetColValues(b.targetColumn)
	if err != nil {
		return nil, 0, err
	}

	// Create sorted list of unique values
	type valueIndex struct {
		value float64
		index int
	}

	sortedValues := make([]valueIndex, 0)
	for i, val := range numValues {
		if validValues[i] && !targetMissing[i] {
			sortedValues = append(sortedValues, valueIndex{val, i})
		}
	}

	if len(sortedValues) < 2 {
		// Not enough values to create a meaningful split
		return nil, 0, nil
	}

	sort.Slice(sortedValues, func(i, j int) bool {
		return sortedValues[i].value < sortedValues[j].value
	})

	var bestThreshold float64
	var bestGainRatio float64 = -1

	uniqueValues := make([]float64, 0)
	lastValue := sortedValues[0].value
	uniqueValues = append(uniqueValues, lastValue)

	for i := 1; i < len(sortedValues); i++ {
		if sortedValues[i].value != lastValue {
			lastValue = sortedValues[i].value
			uniqueValues = append(uniqueValues, lastValue)
		}
	}

	// Consider each unique value (except the last one) as a potential threshold
	for i := 0; i < len(uniqueValues)-1; i++ {
		// Compute the threshold as the midpoint between two consecutive values
		threshold := (uniqueValues[i] + uniqueValues[i+1]) / 2

		// Split the dataset based on this threshold
		lowerSubset, greaterSubset, err := data.SplitNumericDataset(attrName, threshold)
		if err != nil {
			return nil, 0, err
		}

		// Skip if either subset is empty
		if lowerSubset.NumRows == 0 || greaterSubset.NumRows == 0 {
			continue
		}

		// Calculate entropy of each subset
		lowerEntropy, err := b.CalculateEntropy(lowerSubset)
		if err != nil {
			return nil, 0, err
		}

		greaterEntropy, err := b.CalculateEntropy(greaterSubset)
		if err != nil {
			return nil, 0, err
		}

		// Calculate information gain
		totalInstances := lowerSubset.NumRows + greaterSubset.NumRows
		infoGain := classEntropy - (float64(lowerSubset.NumRows)/float64(totalInstances)*lowerEntropy +
			float64(greaterSubset.NumRows)/float64(totalInstances)*greaterEntropy)

		// Calculate split information
		lowerProp := float64(lowerSubset.NumRows) / float64(totalInstances)
		greaterProp := float64(greaterSubset.NumRows) / float64(totalInstances)
		splitInfo := -lowerProp*math.Log2(lowerProp) - greaterProp*math.Log2(greaterProp)

		// Calculate gain ratio
		if splitInfo > 0 {
			gainRatio := infoGain / splitInfo
			if gainRatio > bestGainRatio {
				bestGainRatio = gainRatio
				bestThreshold = threshold
			}
		}
	}

	return bestThreshold, bestGainRatio, nil
}
