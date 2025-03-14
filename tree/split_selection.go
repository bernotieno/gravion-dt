package tree

import (
	"gravion-dt/internal"
)

// findBestSplit finds the attribute with the highest gain ratio
func (b *Builder) findBestSplit(data *internal.Dataset, usedAttributes []string) (string, interface{}, float64, error) {
	// Calculate information content of the current data
	classEntropy, err := b.calculateEntropy(data)
	if err != nil {
		return "", nil, 0, err
	}

	if classEntropy == 0 {
		// All instances belong to the same class, no need to split
		return "", nil, 0, nil
	}

	return "", nil, 0, nil
}
