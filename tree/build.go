package tree

import "dt/internal"

// buildTree recursively builds the decision tree
func (b *Builder) BuildTree(data *internal.Dataset, usedAttributes []string, depth int) (*Node, error) {
	// Check if we've reached maximum depth
	if depth >= b.maxDepth {
		return b.CreateLeafNode(data)
	}

	// Check if the dataset is too small
	if data.NumRows < b.minInstancesPerLeaf {
		return b.CreateLeafNode(data)
	}

	// Get class distribution
	classCounts, err := b.getClassCounts(data)
	if err != nil {
		return nil, err
	}

	// If information gain is too small, make a leaf
	entropy, _ := b.CalculateEntropy(data)
	if entropy < 0.05 { // Minimal entropy threshold
		return b.CreateLeafNode(data)
	}

	// If dataset is almost pure (>95% one class), make a leaf
	if len(classCounts) > 0 {
		total := 0
		maxCount := 0
		for _, count := range classCounts {
			total += count
			if count > maxCount {
				maxCount = count
			}
		}
		if float64(maxCount)/float64(total) > 0.95 {
			return b.CreateLeafNode(data)
		}
	}

	// If all instances belong to the same class, create a leaf node
	if len(classCounts) == 1 {
		return b.CreateLeafNode(data)
	}

	// Find the best attribute to split on
	bestAttr, bestSplit, bestGainRatio, err := b.findBestSplit(data, usedAttributes)
	if err != nil {
		return nil, err
	}

	// If no good split found, create a leaf node
	if bestGainRatio <= 0 {
		return b.CreateLeafNode(data)
	}

	// Create a new node
	if data.ColumnTypes[bestAttr] == internal.CategoricalType {
		// Categorical split
		node := &Node{
			Type:          CategoricalNode,
			AttributeName: bestAttr,
			Children:      make(map[string]*Node),
			ClassCounts:   classCounts,
		}

		// Add the attribute to used attributes
		newUsedAttributes := append(usedAttributes, bestAttr)

		// Get unique values for the attribute
		uniqueValues, err := data.GetUniqueValues(bestAttr)
		if err != nil {
			return nil, err
		}

		// Create a child node for each attribute value
		for _, value := range uniqueValues {
			subset, _, err := data.SplitDataset(bestAttr, value)
			if err != nil {
				return nil, err
			}

			if subset.NumRows == 0 {
				// If no instances match this value, create a leaf node with the majority class
				majority, err := b.getMajorityClass(data)
				if err != nil {
					return nil, err
				}
				conf, err := b.getClassConfidence(data, majority)
				if err != nil {
					return nil, err
				}
				node.Children[value] = &Node{
					Type:        LeafNode,
					Prediction:  majority,
					Confidence:  conf,
					ClassCounts: classCounts,
				}
			} else {
				// Recursively build the tree for this subset
				childNode, err := b.BuildTree(subset, newUsedAttributes, depth+1)
				if err != nil {
					return nil, err
				}
				node.Children[value] = childNode
			}
		}

		return node, nil
	} else if data.ColumnTypes[bestAttr] == internal.NumericType {
		// Numeric split
		threshold := bestSplit.(float64)
		node := &Node{
			Type:          NumericNode,
			AttributeName: bestAttr,
			Threshold:     threshold,
			ClassCounts:   classCounts,
		}

		// Split the dataset based on the threshold
		lowerSubset, greaterSubset, err := data.SplitNumericDataset(bestAttr, threshold)
		if err != nil {
			return nil, err
		}

		// Add the attribute to used attributes (for numeric we don't exclude it from future splits)
		newUsedAttributes := append([]string{}, usedAttributes...)

		// Build left subtree
		if lowerSubset.NumRows >= b.minInstancesPerLeaf {
			leftChild, err := b.BuildTree(lowerSubset, newUsedAttributes, depth+1)
			if err != nil {
				return nil, err
			}
			node.LeftChild = leftChild
		} else {
			// If subset is too small, create a leaf node
			majority, err := b.getMajorityClass(data)
			if err != nil {
				return nil, err
			}
			conf, err := b.getClassConfidence(data, majority)
			if err != nil {
				return nil, err
			}
			node.LeftChild = &Node{
				Type:        LeafNode,
				Prediction:  majority,
				Confidence:  conf,
				ClassCounts: classCounts,
			}
		}

		// Build right subtree
		if greaterSubset.NumRows >= b.minInstancesPerLeaf {
			rightChild, err := b.BuildTree(greaterSubset, newUsedAttributes, depth+1)
			if err != nil {
				return nil, err
			}
			node.RightChild = rightChild
		} else {
			// If subset is too small, create a leaf node
			majority, err := b.getMajorityClass(data)
			if err != nil {
				return nil, err
			}
			conf, err := b.getClassConfidence(data, majority)
			if err != nil {
				return nil, err
			}
			node.RightChild = &Node{
				Type:        LeafNode,
				Prediction:  majority,
				Confidence:  conf,
				ClassCounts: classCounts,
			}
		}

		return node, nil
	}

	// Default case, create a leaf node
	return b.CreateLeafNode(data)
}
