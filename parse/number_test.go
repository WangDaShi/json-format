package parse

import (
	"testing"
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
    ctx := jsonContex{Runes: r,Index: 0}
    str2,err := readNumInner(&ctx,0)
    if err != nil {
        t.Fail()
    }
    if str2 != expected {
        t.Fail()
    }
}

