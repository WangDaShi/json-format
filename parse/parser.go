package parse

import (
	"errors"
	"fmt"
	"json/format/myjson"
	"math"
	"os"
)

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

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// func ParseFile(path string) (myjson.JsonObject, error) {
//     content := readFileToString(path)
//     return Parse(content)
// }

// func Parse(content string) myjson.JsonObject {
//     return 
// }


// read a json string from given position, 
// and generate the first json part
type reader interface {
    read(ctx *jsonContex) myjson.JsonValue
}




type jsonContex struct {
    Runes []rune
    Index int32
}

func (v jsonContex) getCurr() rune {
    return v.Runes[v.Index]
}

func (v jsonContex) incr() {
    v.Index ++
}

// 状态机的状态，
// 0 初始状态
// 1 读取第一个"之后
// 2 遇到转义字符
// 3 遇到结束的"
// MAX 出错
func readString(ctx *jsonContex, state int32) (string, error) {
    if state == math.MaxInt32 {
        return  "", errors.New("solving string type faile, unexpected char")
    }
    if state == 0 {
        for ctx.getCurr() == ' ' {
            ctx.incr()
        }
        if ctx.getCurr() != '"' {
            readString(ctx, math.MaxInt32)
        }else {
            ctx.incr()
            readString(ctx, 1)
        }
    }
    for state == 1 {
        if ctx.getCurr() == '\\' {
            ctx.incr()
            readString(ctx,2)
        } else if ctx.getCurr() == '"' {
            ctx.incr()
            readString(ctx,3)
        } else {
            ctx.incr()
        }
    }

    if state == 3 {
        ctx.incr()
        return "",nil
    }

}

























