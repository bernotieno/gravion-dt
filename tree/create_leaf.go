package tree

import "gravion-dt/internal"

// createLeafNode creates a leaf node for the decision tree using the provided dataset.
// It calculates the majority class, class confidence, and class counts from the dataset
// and returns a Node representing the leaf node.
//
// Parameters:
//   - data: A pointer to an internal.Dataset containing the data for the leaf node.
//
// Returns:
//   - A pointer to a Node representing the leaf node.
//   - An error if any of the calculations (majority class, class confidence, class counts) fail.
func (b *Builder) createLeafNode(data *internal.Dataset) (*Node, error) {
	majorityClass, err := b.getMajorityClass(data)
	if err != nil {
		return nil, err
	}

	confidence, err := b.getClassConfidence(data, majorityClass)
	if err != nil {
		return nil, err
	}

	classCounts, err := b.getClassCounts(data)
	if err != nil {
		return nil, err
	}

	return &Node{
		Type:        LeafNode,
		Prediction:  majorityClass,
		Confidence:  confidence,
		ClassCounts: classCounts,
	}, nil
}