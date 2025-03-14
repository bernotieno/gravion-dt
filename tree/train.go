package tree

import (
	"encoding/json"
	"fmt"
	"os"

	"dt/internal"
)

// Train trains a decision tree model using the dataset provided in the input file.
// The target column specifies the column to be predicted. The trained model is saved
// to the output file.
//
// Parameters:
//   - inputFile: Path to the CSV file containing the dataset.
//   - targetColumn: Name of the column to be used as the target for prediction.
//   - outputFile: Path to the file where the trained model will be saved.
//
// Returns:
//   - error: An error if the training process fails, otherwise nil.
func Train(inputFile, targetColumn, outputFile string) error {
	data, err := internal.ReadCSV(inputFile)
	if err != nil {
		return fmt.Errorf("failed to load dataset: %w", err)
	}

	if !data.HasColumn(targetColumn) {
		return fmt.Errorf("target column '%s' not found in dataset", targetColumn)
	}

	trainData, testData := data.SplitIntoTrainTest(0.8)

	builder := NewBuilder(trainData, targetColumn)

	// Build the decision tree
	decisionTree, err := builder.Build()
	if err != nil {
		return fmt.Errorf("failed to build decision tree: %w", err)
	}

	// Evaluate the model on the test set
	predictions, err := decisionTree.Predict(testData)
	if err != nil {
		return fmt.Errorf("failed to make predictions: %w", err)
	}

	// Compute accuracy
	actualLabels, _, _ := testData.GetColValues(targetColumn)
	correct := 0
	for i, pred := range predictions {
		if pred == actualLabels[i] {
			correct++
		}
	}

	accuracy := float64(correct) / float64(len(predictions))
	fmt.Printf("Model Accuracy: %.2f%%\n", accuracy*100)

	// Save the model to file
	if err := decisionTree.SaveToFile(outputFile); err != nil {
		return fmt.Errorf("failed to save model: %w", err)
	}

	fmt.Printf("Decision tree model successfully trained and saved to %s\n", outputFile)
	return nil
}

// SaveToFile saves the decision tree model to a JSON file
func (dt *DecisionTree) SaveToFile(filename string) error {
	data, err := json.MarshalIndent(dt, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}
