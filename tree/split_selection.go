package tree

import (
	"gravion-dt/internal"
	"sync"
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

	// Filter out attributes that have been used
	availableAttrs := make([]string, 0)
	for _, col := range data.Columns {
		// Skip the target column
		if col.Name == b.targetColumn {
			continue
		}

		// Skip used attributes for categorical features
		if col.Type == internal.CategoricalType {
			isUsed := false
			for _, usedAttr := range usedAttributes {
				if col.Name == usedAttr {
					isUsed = true
					break
				}
			}
			if isUsed {
				continue
			}
		}

		availableAttrs = append(availableAttrs, col.Name)
	}

	// If no attributes available, return empty result
	if len(availableAttrs) == 0 {
		return "", nil, 0, nil
	}
	// Use goroutines to calculate gain ratio for each attribute in parallel
	type attrResult struct {
		attrName  string
		splitVal  interface{}
		gainRatio float64
		err       error
	}

	resultChan := make(chan attrResult, len(availableAttrs))
	var wg sync.WaitGroup
	return "", nil, 0, nil
}
