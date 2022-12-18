package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {

	path := "D://code/go/json-format/test.json"
	path2, err := Format(path)
	if err != nil {
		panic(err)
	} else {
		println(path2)
	}
}

func Format(filePath string) (string, error) {

	content := readFileToString(filePath)
	var start int
	for i := 0; i < len(content); i++ {
		// fmt.Println(string(content[i]))
		if content[i] == ' ' {
			continue
		} else if content[i] == '{' {
			start = i
			break
		} else {
			return "", errors.New("Format 格式错误")
		}
	}
	solveObject(content, start)
	return "nothing", nil
}

func solveObject(content string, start int) (int, error) {
	fmt.Printf("start is %d \n", start)
	i := start + 1
	for i < len(content) {
		curr := content[i]
		if curr == ' ' || curr == ',' || curr == '\n' || curr == '\r' {
			i++
			continue
		} else if curr == '"' {
			i, err := solveValuePair(content, i)
			if err != nil {
				return -1, err
			} else {
				i++
			}
		} else if curr == '}' {
			return i, nil
		} else {
			return -1, errors.New("object 格式错误")
		}
	}
	return i, nil
}

func solveValuePair(content string, start int) (int, error) {
	println(content[start:])
	i, err := solveString(content, start)
	fmt.Printf("name: %s\n", content[start:i])

	if err != nil {
		return -1, err
	}
	for i < len(content) {
		if content[i] == ' ' {
			i++
		}
		break
	}
	if content[i] != ':' {
		return -1, errors.New("value pair 格式错误")
	}
	i++
	for i < len(content) {
		if content[i] == ' ' {
			i++
		}
		break
	}
	j, err := solveValue(content, i)
	if err != nil {
		return -1, err
	}
    println(content[18:])
	fmt.Printf("i:%d,j:%d", i, j)

	fmt.Printf("value: %s", content[i:j])
	return j, nil
}

func solveValue(content string, start int) (int, error) {
	i := start + 1
	for i < len(content) {
		if content[i] == ' ' {
			i++
		} else {
			break
		}
	}
	if content[i] == '"' {
		return solveString(content, i)
	}
	// if content[i] >= '0' && content[i] <= 9 {
	//     return solveNum(content,i)
	// }
	// if(content[i:i+4] == "true") {
	//
	// }
	return -1, nil
}

func solveNum(content string, start int) (int, error) {
	return -1, nil
}

func solveString(content string, start int) (int, error) {
	if content[start] != '"' {
		return -1, errors.New("not a legal string")
	}
	i := start + 1
	for i < len(content) {
		if content[i] == '"' {
			break
		} else {
			i++
		}
	}
	if content[i] == '"' {
		return i + 1, nil
	} else {
		return -1, errors.New("not a legal string")
	}
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
