package tree

import (
	"dt/internal"
	"fmt"
	"sync"
)

// Tune performs hyperparameter tuning by training multiple decision trees with
// different hyperparameter combinations and returns the best model.
//
// Parameters:
//   - inputFile: Path to the CSV file containing the dataset.
//   - targetColumn: Name of the column to be used as the target for prediction.
//   - outputFile: Path to the file where the best model will be saved.
//
// Returns:
//   - error: An error if the tuning process fails, otherwise nil.
func Tune(inputFile, targetColumn, outputFile string) error {
	// Load the dataset
	data, err := internal.ReadCSV(inputFile, false)
	if err != nil {
		return fmt.Errorf("failed to load dataset: %w", err)
	}

	if !data.HasColumn(targetColumn) {
		return fmt.Errorf("target column '%s' not found in dataset", targetColumn)
	}

	// Split the data into train, validation, and test sets (60%, 20%, 20%)
	trainData, tempData := data.SplitIntoTrainTest(0.6)
	validationData, testData := tempData.SplitIntoTrainTest(0.5)

	// Define hyperparameter combinations to try
	hyperparamSets := []HyperparameterSet{
		{MinInstancesPerLeaf: 5, MaxDepth: 5},
		{MinInstancesPerLeaf: 5, MaxDepth: 10},
		{MinInstancesPerLeaf: 5, MaxDepth: 15},
		{MinInstancesPerLeaf: 10, MaxDepth: 5},
		{MinInstancesPerLeaf: 10, MaxDepth: 10},
		{MinInstancesPerLeaf: 10, MaxDepth: 15},
		{MinInstancesPerLeaf: 15, MaxDepth: 5},
		{MinInstancesPerLeaf: 15, MaxDepth: 10},
		{MinInstancesPerLeaf: 15, MaxDepth: 15},
	}

	// Train models with different hyperparameter combinations in parallel
	var wg sync.WaitGroup
	var mutex sync.Mutex

	fmt.Println("Starting hyperparameter tuning...")
	fmt.Println("Testing", len(hyperparamSets), "hyperparameter combinations")

	for i := range hyperparamSets {
		wg.Add(1)
		go func(paramSet *HyperparameterSet) {
			defer wg.Done()

			// Build a model with the current hyperparameter set
			builder := NewBuilder(trainData, targetColumn)
			builder.SetMinInstancesPerLeaf(paramSet.MinInstancesPerLeaf)
			builder.SetMaxDepth(paramSet.MaxDepth)

			// Train the model
			tree, err := builder.Build()
			if err != nil {
				fmt.Printf("Error training with MinInst=%d, MaxDepth=%d: %v\n",
					paramSet.MinInstancesPerLeaf, paramSet.MaxDepth, err)
				return
			}

			// Evaluate on validation set
			predictions, err := tree.Predict(validationData)
			if err != nil {
				fmt.Printf("Error predicting with MinInst=%d, MaxDepth=%d: %v\n",
					paramSet.MinInstancesPerLeaf, paramSet.MaxDepth, err)
				return
			}

			// Calculate accuracy
			actualLabels, _, _ := validationData.GetColValues(targetColumn)
			correct := 0
			for i, pred := range predictions {
				if pred == actualLabels[i] {
					correct++
				}
			}
			accuracy := float64(correct) / float64(len(predictions))

			// Store the accuracy
			mutex.Lock()
			paramSet.Accuracy = accuracy
			fmt.Printf("MinInst=%d, MaxDepth=%d: Validation Accuracy=%.4f\n",
				paramSet.MinInstancesPerLeaf, paramSet.MaxDepth, accuracy)
			mutex.Unlock()
		}(&hyperparamSets[i])
	}

	// Wait for all models to be trained and evaluated
	wg.Wait()

	// Find the best hyperparameter set
	bestIdx := 0
	for i := 1; i < len(hyperparamSets); i++ {
		if hyperparamSets[i].Accuracy > hyperparamSets[bestIdx].Accuracy {
			bestIdx = i
		} else if hyperparamSets[i].Accuracy == hyperparamSets[bestIdx].Accuracy {
			// If accuracy is the same, prefer lower depth
			if hyperparamSets[i].MaxDepth < hyperparamSets[bestIdx].MaxDepth {
				bestIdx = i
			} else if hyperparamSets[i].MaxDepth == hyperparamSets[bestIdx].MaxDepth {
				// If depth is the same, prefer higher MinInstancesPerLeaf
				if hyperparamSets[i].MinInstancesPerLeaf > hyperparamSets[bestIdx].MinInstancesPerLeaf {
					bestIdx = i
				}
			}
		}
	}

	bestParams := hyperparamSets[bestIdx]
	fmt.Printf("\nBest hyperparameters found: MinInst=%d, MaxDepth=%d with validation accuracy=%.4f\n",
		bestParams.MinInstancesPerLeaf, bestParams.MaxDepth, bestParams.Accuracy)

	// Train the final model with the best hyperparameters using the combined train+validation set
	combinedData := trainData.Combine(validationData)
	builder := NewBuilder(combinedData, targetColumn)
	builder.SetMinInstancesPerLeaf(bestParams.MinInstancesPerLeaf)
	builder.SetMaxDepth(bestParams.MaxDepth)

	// Build the final decision tree
	decisionTree, err := builder.Build()
	if err != nil {
		return fmt.Errorf("failed to build final decision tree: %w", err)
	}

	// Evaluate the model on the test set
	predictions, err := decisionTree.Predict(testData)
	if err != nil {
		return fmt.Errorf("failed to make predictions: %w", err)
	}

	// Compute test accuracy
	actualLabels, _, _ := testData.GetColValues(targetColumn)
	correct := 0
	for i, pred := range predictions {
		if pred == actualLabels[i] {
			correct++
		}
	}

	testAccuracy := float64(correct) / float64(len(predictions))
	fmt.Printf("Final Model Test Accuracy: %.2f%%\n", testAccuracy*100)

	// Save the model to file
	if err := decisionTree.SaveToFile(outputFile); err != nil {
		return fmt.Errorf("failed to save model: %w", err)
	}

	fmt.Printf("Best decision tree model successfully trained and saved to %s\n", outputFile)
	return nil
}
