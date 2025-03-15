
## Package Documentation: `internal`

The `internal` package provides a robust set of utilities for handling tabular datasets, including reading, manipulating, and splitting datasets. It is designed to work with CSV files and supports various data types such as categorical, numeric, date, and timestamp. The package also includes functionality for splitting datasets into training and testing subsets, as well as saving predictions to a CSV file.

---

## Key Features

### 1. **Efficient CSV Reading with Sampling**
   - **Reason:** When dealing with large datasets, reading the entire file can be time-consuming and memory-intensive. To address this, the package estimates the number of rows in the CSV file by analyzing the first 100 lines. If the estimated number of rows exceeds a threshold (e.g., 10,000), the program samples only 10% of the rows for processing.
   - **Advantages:**
     - **Efficiency:** Reduces the computational load and training time.
     - **Accuracy:** Ensures that the sampled data is representative of the entire dataset, reducing the risk of overfitting.
     - **Scalability:** Adapts to the size of the dataset, using the entire dataset for small files and sampling for large files.

   **Function:** `ReadCSV`

---

### 2. **Dataset Splitting for Training and Testing**
   - **Reason:** To evaluate the performance of a machine learning model, it is essential to split the dataset into training and testing subsets. By default, 80% of the data is used for training, and 20% is used for testing. The split is randomized to ensure a fair distribution of data.
   - **Advantages:**
     - **Fair Evaluation:** Ensures that the model is evaluated on unseen data, providing a more accurate measure of its performance.
     - **Flexibility:** Allows users to specify a custom train-test ratio.

   **Function:** `SplitIntoTrainTest`

---

### 3. **Column Type Detection**
   - **Reason:** Automatically detecting the type of data in each column (e.g., numeric, categorical, date, timestamp) simplifies data preprocessing and ensures that the correct operations are applied to each column.
   - **Advantages:**
     - **Automation:** Reduces the need for manual data type specification.
     - **Accuracy:** Ensures that data is processed correctly based on its type.

   **Function:** `DetectColumnType`

---

### 4. **Handling Missing Values**
   - **Reason:** Missing values are common in real-world datasets. The package tracks missing values in each column, allowing users to handle them appropriately during data processing.
   - **Advantages:**
     - **Robustness:** Ensures that missing values do not disrupt data processing.
     - **Flexibility:** Allows users to decide how to handle missing values (e.g., imputation, removal).

   **Functions:** `GetColValues`, `GetNumericValues`

---

### 5. **Saving Predictions to CSV**
   - **Reason:** After making predictions, it is often necessary to save the results for further analysis or reporting. The package provides a utility to save predictions to a CSV file.
   - **Advantages:**
     - **Convenience:** Simplifies the process of saving and sharing prediction results.
     - **Compatibility:** Ensures that predictions are saved in a widely-used format.

   **Function:** `SavePredictions`

---

### 6. **Subset Creation for Efficient Data Handling**
   - **Reason:** When working with large datasets, it is often necessary to create subsets of the data for specific tasks (e.g., training, testing, or analysis). The package provides a utility to create subsets based on specified row indices.
   - **Advantages:**
     - **Efficiency:** Reduces memory usage by working with smaller subsets of data.
     - **Flexibility:** Allows the program to create custom subsets for specific tasks.

   **Function:** `createSubset`

---

## Example Usage

### Reading a CSV File
```go
dataset, err := ReadCSV("data.csv")
if err != nil {
    log.Fatalf("Failed to read CSV: %v", err)
}
```

### Splitting into Training and Testing Subsets
```go
trainDataset, testDataset := dataset.SplitIntoTrainTest(0.8)
```

### Saving Predictions
```go
predictions := []string{"prediction1", "prediction2", "prediction3"}
err := SavePredictions(predictions, "predictions.csv")
if err != nil {
    log.Fatalf("Failed to save predictions: %v", err)
}
```

---

## Summary

The `internal` package is designed to handle tabular data efficiently and effectively. By providing utilities for reading, splitting, and processing datasets, it simplifies the workflow for data preprocessing and machine learning tasks. The package's focus on efficiency, accuracy, and flexibility makes it suitable for both small and large datasets.
