package test

import (
	"testing"
    "json/format/parse"
)

func TestNum2(t *testing.T) {
    testNum("123 ","123",t)
    testNum("0 ","0",t)
    testNum("1 ","1",t)
    testNum("-1 ","-1",t)
    testNum("12324242 ","12324242",t)
    testNum("0.123 ","0.123",t)
    testNum("-0.123 ","-0.123",t)
}

func testNum(str string,expected string,t *testing.T) {
    r := []rune(str)
    ctx := parse.JsonContex{Runes: r,Index: 0}
    numV,err := parse.ReadNum(&ctx)
    if err != nil {
        t.Fail()
    }
    if numV.Value != expected {
        t.Fail()
    }
}

