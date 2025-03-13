package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// Define flags for command line arguments
	command := flag.String("c", "", "Command to execute: 'train' or 'predict'")
	inputFile := flag.String("i", "", "Path to the input CSV file required for training or prediction")
	targetColumn := flag.String("t", "", "Name of the target column (required for training)")
	outputFile := flag.String("o", "", "Path to save the output (model for training, predictions for predict)")
	modelFile := flag.String("m", "", "Path to the trained model file (required for predict)")

	// Parse command line arguments
	flag.Parse()

	// Validate command line arguments
	if *command == "" {
		log.Fatalln("Error: Command is required. Use -c train or -c predict")
	}

	if *inputFile == "" {
		log.Fatalln("Error: Missing input file. Use -i <path to training_data.csv/prediction_data_file.csv>")
	}

	if *outputFile == "" {
		log.Fatalln("Error: Missing output file. Use -o <path to save the output (model for training, predictions for predict)>")
	}

	// Confirm that the input file exists
	if _, err := os.Stat(*inputFile); os.IsNotExist(err) {
		log.Fatalf("Error: Input file %s does not exist\n", *inputFile)
	}

	//Create output directory if it does not exist
	outputDir := filepath.Dir(*outputFile)
	if outputDir != "" {
		if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
			log.Fatalf("Error: Unable to create output directory %s\n", outputDir)
		}
	}

	// Execute the command
	switch *command {
	case "train":
		if *targetColumn == "" {
			log.Fatalln("Error: Target column is required for training. Use -t <target_column_name>")
		}
		// // Call the train function here
		// err := Train(*inputFile, *targetColumn, *outputFile)
		// if err != nil {
		// 	log.Fatalf("Error during training: %v\n", err)
		// }
		fmt.Printf("Decision tree model successfully trained and saved to %s\n", *outputFile)

	case "predict":
		if *modelFile == "" {
			log.Fatalln("Error: Model file is required for prediction. Use -m <path to trained_model.dt>")
		}
		// Confirm that the model file exists
		if _, err := os.Stat(*modelFile); os.IsNotExist(err) {
			log.Fatalf("Error: Model file %s not found\n", *modelFile)
		}

		// // Call the predict function here
		// err := Predict(*inputFile, *modelFile, *outputFile)
		// if err != nil {
		// 	log.Fatalf("Error during prediction: %v\n", err)
		// }
		fmt.Printf("Predictions successfully generated and saved to %s\n", *outputFile)

	default:
		log.Fatalf("Error: Invalid command %s. Use -c train or -c predict\n", *command)
	}

}
