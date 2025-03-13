package internal

import (
	"fmt"
	"golang.org/x/exp/mmap"
)

// ReadingCSV reads a file using memory-mapped I/O
func ReadingCSV(filename string) (string, error){
	reader, err := mmap.Open(filename)
	if err != nil {
		return "", fmt.Errorf("error opening file: %v", err)
	}
	defer reader.Close()

	data := make([]byte, reader.Len())
	_, err = reader.ReadAt(data, 0)

	if err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}

	return string(data), nil

}