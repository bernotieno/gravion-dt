package main

import (
	"fmt"
	"gravion-dt/internal"
	"log"
)


func main() {
	filename := "input.csv"
	content, err := internal.ReadingCSV(filename)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(content)

}