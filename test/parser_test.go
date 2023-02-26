package test

import (
	"json/format/tool"
	"json/format/parse"
	"testing"
)

const successFile = "test_string_success.txt"
const failFile = "test_string_fail.txt"
const parseFile = "test_parse.json"

// TODO 逗号没有检查的问题
// 顺序的问题
// 数字的问题
func TestParse2(t *testing.T) {
	folder := tool.GetProjectFolder() + "/test/file/"
	content := tool.ReadString(folder + parseFile)
	obj, err := parse.Parse(content)
	if err != nil {
		panic(err)
	}
	formatJson := obj.ToJson(0)
	println(formatJson)
}

func testString2(t *testing.T) {
	folder := tool.GetProjectFolder() + "/test/file/"
	line := tool.ReadLine(folder + successFile)
	for _, str := range line {
		testString(str, true, t)
	}
	line2 := tool.ReadLine(folder + failFile)
	for _, str := range line2 {
		testString(str, false, t)
	}
}

func testString(str string, expected bool, t *testing.T) {
	rr := []rune(str)
	ctx := parse.JsonContex{Runes: rr, Index: 0}
	value, err := parse.SolveString(&ctx)
	acc := err == nil
	if err == nil {
		println(value)
	}
	if acc != expected {
		t.Errorf("failed, expected: %t,acctual: %t,str:%s", expected, acc, str)
	}
}

func testParse(t *testing.T, input string, expectedOutput string) {
	output, err := parse.Parse(input)
	if err != nil {
		panic(err)
	}
	if output.ToJson(0) != expectedOutput {
		t.Errorf("Parse(%q) = %q; expected %q", input, output, expectedOutput)
	}
}

func testParseChatGPT(t *testing.T) {
	// Test case 1: valid JSON input with object
	input1 := `{"name": "John", "age": 30, "city": "New York"}`
	expectedOutput1 := `{
    "name": "John",
    "age": 30,
    "city": "New York"
}`
	testParse(t, input1, expectedOutput1)

	// Test case 2: invalid JSON input with object
	input2 := `{"name": "John", "age": 30, "city": "New York"`
	expectedOutput2 := ""
	testParse(t, input2, expectedOutput2)

	// Test case 3: empty input with object
	input3 := ``
	expectedOutput3 := ""
	testParse(t, input3, expectedOutput3)

	// Test case 4: valid JSON input with array
	input4 := `["apple", "banana", "cherry"]`
	expectedOutput4 := `[
    "apple",
    "banana",
    "cherry"
]`
	testParse(t, input4, expectedOutput4)

	// Test case 5: valid JSON input with integer
	input5 := `42`
	expectedOutput5 := `42`
	testParse(t, input5, expectedOutput5)

	// Test case 6: valid JSON input with float
	input6 := `3.14159`
	expectedOutput6 := `3.14159`
	testParse(t, input6, expectedOutput6)

	// Test case 7: valid JSON input with negative integer
	input7 := `-42`
	expectedOutput7 := `-42`
	testParse(t, input7, expectedOutput7)

	// Test case 8: valid JSON input with scientific notation
	input8 := `2.998e8`
	expectedOutput8 := `2.998e+08`
	testParse(t, input8, expectedOutput8)

	// Test case 9: valid JSON input with nested object
	input9 := `{"name": "John", "age": 30, "address": {"city": "New York", "state": "NY"}}`
	expectedOutput9 := `{
    "name": "John",
    "age": 30,
    "address": {
        "city": "New York",
        "state": "NY"
    }
}`
	testParse(t, input9, expectedOutput9)

	// Test case 10: valid JSON input with nested array
	input10 := `{"fruits": ["apple", "banana", "cherry"], "veggies": ["carrot", "pepper", "lettuce"]}`
	expectedOutput10 := `{
    "fruits": [
        "apple",
        "banana",
        "cherry"
    ],
    "veggies": [
        "carrot",
        "pepper",
        "lettuce"
    ]
}`
	testParse(t, input10, expectedOutput10)
}
