package main

import (
	"fmt"
	"gravion-dt/internal"
	"log"
)


func main() {
	filename := "datasets/train.csv"
	content, err := internal.ReadingCSV(filename)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(content)

}