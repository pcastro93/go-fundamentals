package main

import (
	"bufio"
	"fmt"
	"os"
)

func readFile() {
	path := "./main.go"
	file, openErr := os.Open(path)
	if openErr != nil {
		fmt.Println("File can't be opened", openErr)
		return
	}
	// Closing the file might fail, so we defer the closing operation
	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			fmt.Println("Error closing file", closeErr)
		}
	}()
	// Print the file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if scannerErr := scanner.Err(); scannerErr != nil {
		fmt.Println("Error reading the file", scannerErr)
	}
}

func multipleDefers() {
	for i := range 3 {
		defer fmt.Printf("Call %d\n", i+1)
	}
}

func main() {
	readFile()
	multipleDefers()
}
