package tree

import (
	"gravion-dt/internal"
	"math"
)

// calculateEntropy calculates the entropy of the given dataset.
// Entropy is a measure of the amount of uncertainty or impurity in the dataset.
//
// Parameters:
//   data (*internal.Dataset): The dataset for which to calculate the entropy.
//
// Returns:
//   float64: The calculated entropy value.
//   error: An error if there is an issue calculating the class counts or if the dataset is empty.
func (b *Builder) CalculateEntropy(data *internal.Dataset) (float64, error) {
	// Get class counts
	classCounts, err := b.getClassCounts(data)
	if err != nil {
		return 0, err
	}

	// Calculate total instances
	total := 0
	for _, count := range classCounts {
		total += count
	}

	if total == 0 {
		return 0, nil
	}

	// Calculate entropy
	var entropy float64
	for _, count := range classCounts {
		if count > 0 {
			proportion := float64(count) / float64(total)
			entropy -= proportion * math.Log2(proportion)
		}
	}

	return entropy, nil
}
