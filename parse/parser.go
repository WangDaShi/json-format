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

func Parse(content string) (myjson.JsonObject, error) {
	runes := []rune(content)
	ctx := jsonContex{Runes: runes, Index: 0}
	return readObj(&ctx)
}

// read a json string from given position,
// and generate the first json part
type reader interface {
	read(ctx *jsonContex) myjson.JsonValue
}

type jsonContex struct {
	Runes []rune
	Index int64
}

func (v *jsonContex) getCurr() (rune, error) {
	if v.Index >= int64(len(v.Runes)) {
		return ' ', errors.New("string end unexpected")
	}
	return v.Runes[v.Index], nil
}

func (v *jsonContex) incr() {
	v.Index++
}

func (v *jsonContex) subString(a int64, b int64) string {
	return string(v.Runes[a:b])
}

func readObj(ctx *jsonContex) (myjson.JsonObject, error) {
	m := make(map[string]myjson.JsonValue)
	obj := myjson.JsonObject{Data: m}
	if err := readObjStr(ctx, 0, obj); err != nil {
		return obj, errors.Join(err)
	}
	return obj, nil
}

func readArr(ctx *jsonContex) (myjson.ArrayValue, error) {
	ctx.incr()
	if err := readWhite(ctx); err != nil {
		return myjson.ArrayValue{}, errors.Join(err)
	}
	arr := []myjson.JsonValue{}
	curr, _ := ctx.getCurr()
	if curr == ']' {
		return myjson.ArrayValue{Arr: arr}, nil
	}
	for curr != ']' {
		value, err := readValue(ctx, 0)
		if err != nil {
			return myjson.ArrayValue{}, errors.Join(err)
		}
		arr = append(arr, value)
		if err := readWhite(ctx); err != nil {
			return myjson.ArrayValue{}, errors.Join(err)
		}
		curr, _ = ctx.getCurr()
        if curr == ',' {
            ctx.incr()
        }
		if err := readWhite(ctx); err != nil {
			return myjson.ArrayValue{}, errors.Join(err)
		}
	}
	if curr == ']' {
        ctx.incr()
		return myjson.ArrayValue{Arr: arr}, nil
	} else {
		return myjson.ArrayValue{}, errors.New("arr end upexpected")
	}
}

func readObjStr(ctx *jsonContex, state int64, obj myjson.JsonObject) error {
	if state == 0 {
		if err := readWhite(ctx); err != nil {
			return errors.Join(err)
		}
		r, err := ctx.getCurr()
		if r != '{' {
			return errors.Join(err)
		} else {
			return readObjStr(ctx, 1, obj)
		}
	}
	if state == 1 {
		ctx.incr()
		if err := readWhite(ctx); err != nil {
			return errors.Join(err)
		}
		curChar, _ := ctx.getCurr()
		for curChar != '}' {
			k, v, err := readKeyValue(ctx)
			if err != nil {
				return err
			}
			obj.Data[k] = v
			if err := readWhite(ctx); err != nil {
				return errors.Join(err)
			}
			if k, _ := ctx.getCurr(); k == ',' {
				ctx.incr()
				if err := readWhite(ctx); err != nil {
					return errors.Join(err)
				}
			}
			curChar, _ = ctx.getCurr()
		}
		return nil
	}
	return nil
}

func readKeyValue(ctx *jsonContex) (string, myjson.JsonValue, error) {
	if err := readWhite(ctx); err != nil {
		return "", myjson.JsonObject{}, errors.Join(err)
	}
	key, err := readString(ctx, 0, ctx.Index)
	if err != nil {
		return "", myjson.JsonObject{}, errors.Join(err)
	}
	if err := readWhite(ctx); err != nil {
		return "", myjson.JsonObject{}, errors.Join(err)
	}
	if r, _ := ctx.getCurr(); r != ':' {
		return "", myjson.JsonObject{}, errors.New("errr while parse :")
	}
	ctx.incr()
	value, err := readValue(ctx, 0)
	if err != nil {
		return "", myjson.JsonObject{}, errors.Join(err)
	}
	return key, value, nil
}

func readValue(ctx *jsonContex, state int64) (myjson.JsonValue, error) {
	if err := readWhite(ctx); err != nil {
		return myjson.JsonObject{}, errors.Join(err)
	}
	c, _ := ctx.getCurr()
	if c == '"' {
		v, err := readStr(ctx)
		if err == nil {
			return v, nil
		} else {
			return v, errors.Join(err)
		}
	}
	if c == '{' {
		return readObj(ctx)
	}
	if c == '[' {
		return readArr(ctx)
	}
	next := ctx.subString(ctx.Index, ctx.Index+4)
	if next == "true" {
		ctx.Index += 4
		return myjson.BooleanValue{Value: true}, nil
	}
	if next == "null" {
		ctx.Index += 4
		return myjson.NullValue{}, nil
	}
	next = ctx.subString(ctx.Index, ctx.Index+5)
	if next == "false" {
		ctx.Index += 5
		return myjson.BooleanValue{Value: false}, nil
	}
    return readNum(ctx)
}

func readNum(ctx *jsonContex) (myjson.NumValue,error) {
    return myjson.NumValue{},errors.New("not supported type")
}

// walk through all the white space from given point,until a non white character show
func readWhite(ctx *jsonContex) error {
	for {
		r, err := ctx.getCurr()
		if err != nil {
			return errors.Join(err)
		}
		if whiteChar[r] {
			ctx.incr()
		} else {
			break
		}
	}
	return nil
}

func readStr(ctx *jsonContex) (myjson.StringValue, error) {
	str, err := readString(ctx, 0, ctx.Index)
	if err != nil {
		return myjson.StringValue{}, errors.Join(err)
	}
	return myjson.StringValue{Value: str}, nil
}

// 状态机的状态，
// 0 初始状态
// 1 读取第一个"之后
// 2 遇到转义字符
// 3 遇到结束的"
// MAX 出错
func readString(ctx *jsonContex, state int64, start int64) (string, error) {
	for state == 0 {
		r, err := ctx.getCurr()
		if err != nil {
			return "", errors.Join(err)
		}
		if whiteChar[r] {
			ctx.incr()
		} else if r == '"' {
			ctx.incr()
			return readString(ctx, 1, ctx.Index)
		} else {
			return readString(ctx, math.MaxInt64, start)
		}
	}
	for state == 1 {
		r, err := ctx.getCurr()
		if err != nil {
			return "", errors.Join(err)
		}
		if r == '\\' {
			ctx.incr()
			return readString(ctx, 2, start)
		} else if r == '"' {
			return readString(ctx, 3, start)
		} else {
			ctx.incr()
		}
	}
	if state == 2 {
		r, err := ctx.getCurr()
		if err != nil {
			return "", errors.Join(err)
		}
		if specChar[r] {
			ctx.incr()
			return readString(ctx, 1, start)
		} else {
			return "", errors.New("solving string type fail, unexpected char")
		}
	}
	if state == 3 {
		str := ctx.subString(start, ctx.Index)
		ctx.incr()
		return str, nil
	}
	return "", errors.New("solving string type fail, unexpected char")

}

var whiteChar = map[rune]bool{
	' ':  true,
	'\n': true,
	'\t': true,
	'\r': true,
}

var specChar = map[rune]bool{
	'"':  true,
	'\\': true,
	'/':  true,
	'b':  true,
	'f':  true,
	'n':  true,
	'r':  true,
	't':  true,
}

func SolveString(ctx *jsonContex) (string, error) {
	fmt.Println()
	str, error := readString(ctx, 0, ctx.Index)
	if error != nil {
		return "", errors.New("error while try solving string")
	}
	return str, nil
}
