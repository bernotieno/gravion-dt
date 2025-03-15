package tree

import (
	"encoding/json"
	"fmt"
	"os"

	"dt/internal"
)

// Predict uses a trained decision tree model to make predictions on new data.
// It loads a model from a file, reads input data, validates feature columns,
// makes predictions, and saves the results to an output file.
//
// Parameters:
//   - inputFile: Path to the CSV file containing data for prediction.
//   - modelFile: Path to the JSON file containing the trained decision tree model.
//   - outputFile: Path to the file where predictions will be saved.
//
// Returns:
//   - error: An error message if any step fails; otherwise, nil.
func Predict(inputFile, modelFile, outputFile string) error {
	// Load the model from file
	fmt.Printf("Loading model %s...\n", modelFile)
	decisionTree, err := LoadFromFile(modelFile)
	if err != nil {
		return fmt.Errorf("failed to load model: %w", err)
	}

	// Load the dataset for prediction
	fmt.Printf("Reading %s file...\n", inputFile)
	data, err := internal.ReadCSV(inputFile, true)
	if err != nil {
		return fmt.Errorf("failed to load dataset: %w", err)
	}

	// Validate that all required feature columns exist
	missingColumns := decisionTree.ValidateFeatures(data)
	if len(missingColumns) > 0 {
		return fmt.Errorf("missing required feature columns in test dataset: %v", missingColumns)
	}

	// Perform predictions using the loaded model
	fmt.Printf("\nMaking predictions from %s...\n", modelFile)
	predictions, err := decisionTree.Predict(data)
	if err != nil {
		return fmt.Errorf("prediction error: %w", err)
	}

	// Save the predictions to the output file
	if err := internal.SavePredictions(predictions, outputFile); err != nil {
		return fmt.Errorf("failed to save predictions: %w", err)
	}

	fmt.Printf("Prediction made successfully and saved to %s\n", outputFile)
	return nil
}

// LoadFromFile loads a trained decision tree model from a JSON file.
//
// Parameters:
//   - filename: Path to the JSON file containing the decision tree model.
//
// Returns:
//   - *DecisionTree: A pointer to the loaded DecisionTree instance.
//   - error: An error message if loading or parsing fails; otherwise, nil.
func LoadFromFile(filename string) (*DecisionTree, error) {
	// Read the file content
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into a DecisionTree struct
	var tree DecisionTree
	if err := json.Unmarshal(data, &tree); err != nil {
		return nil, err
	}

	return &tree, nil
}
