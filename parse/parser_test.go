package parse

import (
	"json/format/tool"
	"testing"
)

const successFile = "test_string_success.txt"
const failFile = "test_string_fail.txt"
const parseFile = "test_parse.json"

// TODO 逗号没有检查的问题
// 顺序的问题
// 数字的问题
func TestParse(t *testing.T) {
    folder := tool.GetProjectFolder() + "/parse/"
    content := tool.ReadString(folder + parseFile) 
    obj,err := Parse(content)
    if err != nil {
        panic(err)
    }
    formatJson := obj.ToJson(0)
    println(formatJson)
}

func TestString(t *testing.T) {
    folder := tool.GetProjectFolder() + "/parse/"
    line := tool.ReadLine(folder + successFile) 
    for _,str := range line {
        testString(str,true,t)
    }
    line2 := tool.ReadLine(folder + failFile) 
    for _,str := range line2 {
        testString(str,false,t)
    }
}

func testString(str string,expected bool, t *testing.T) {
    rr := []rune(str)
    ctx := jsonContex{Runes: rr,Index : 0}
    value,err := SolveString(&ctx)
    acc := err == nil
    if err == nil {
        println(value)
    }
    if  acc != expected {
        t.Errorf("failed, expected: %t,acctual: %t,str:%s", expected, acc, str)
    }
}


