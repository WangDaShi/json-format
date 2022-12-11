package main

import (
	"fmt"
	"os"
)

func main() {
	path := "D://code/go/json-format/test.json"
	println(Format(path))
}

func Format(filePath string) string {
	content := readFileToString(filePath)
	return content
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// read file and return the content
func readFileToString(filePath string) string {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		fmt.Printf("file %s not exist !\n", filePath)
		panic(err)
	}
	data, err := os.ReadFile(filePath)
	check(err)
	return string(data)
}
