package parse

import (
    "errors"
	"fmt"
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

func Parse(content string) (JsonObject, error) {
	runes := []rune(content)
	ctx := JsonContex{Runes: runes, Index: 0}
	return readObj(&ctx)
}

// read a json string from given position,
// and generate the first json part
type reader interface {
	read(ctx *JsonContex) JsonValue
}

type JsonContex struct {
	Runes []rune
	Index int64
}

func (v *JsonContex) getCurr() (rune, error) {
	if v.Index >= int64(len(v.Runes)) {
		return ' ', errors.New("string end unexpected")
	}
	return v.Runes[v.Index], nil
}

func (v *JsonContex) incr() {
	v.Index++
}

func (v *JsonContex) subString(a int64, b int64) string {
	return string(v.Runes[a:b])
}

func readObj(ctx *JsonContex) (JsonObject, error) {
	m := make(map[string]JsonValue)
	obj := JsonObject{Data: m}
	if err := readObjInner(ctx, &obj); err != nil {
        
		return obj, errors.Join(err)
	}
	return obj, nil
}

func readArr(ctx *JsonContex) (ArrayValue, error) {
	ctx.incr()
	if err := readWhite(ctx); err != nil {
		return ArrayValue{}, errors.Join(err)
	}
	arr := []JsonValue{}
	curr, _ := ctx.getCurr()
	if curr == ']' {
		return ArrayValue{Arr: arr}, nil
	}
	for curr != ']' {
		value, err := readValue(ctx, 0)
		if err != nil {
			return ArrayValue{}, errors.Join(err)
		}
		arr = append(arr, value)
		if err := readWhite(ctx); err != nil {
			return ArrayValue{}, errors.Join(err)
		}
		curr, _ = ctx.getCurr()
		if curr == ',' {
			ctx.incr()
		}
		if err := readWhite(ctx); err != nil {
			return ArrayValue{}, errors.Join(err)
		}
	}
	if curr == ']' {
		ctx.incr()
		return ArrayValue{Arr: arr}, nil
	} else {
		return ArrayValue{}, errors.New("arr end upexpected")
	}
}

func readObjInner(ctx *JsonContex, obj *JsonObject) error {
	if err := readWhite(ctx); err != nil {
		return errors.Join(err)
	}
	r, err := ctx.getCurr()
	if r != '{' {
		return errors.Join(err)
	}
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
		err2 := obj.Add(k, &v)
		if err2 != nil {
			return errors.Join(err2)
		}
		// obj.Data[k] = v
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
	ctx.incr()
	return nil
}

func readKeyValue(ctx *JsonContex) (string, JsonValue, error) {
	if err := readWhite(ctx); err != nil {
		return "", JsonObject{}, errors.Join(err)
	}
	key, err := readString(ctx, 0, ctx.Index)
	if err != nil {
		return "", JsonObject{}, errors.Join(err)
	}
	if err := readWhite(ctx); err != nil {
		return "", JsonObject{}, errors.Join(err)
	}
	if r, _ := ctx.getCurr(); r != ':' {
		return "", JsonObject{}, errors.New("errr while parse :")
	}
	ctx.incr()
	value, err := readValue(ctx, 0)
	if err != nil {
		return "", JsonObject{}, errors.Join(err)
	}
	return key, value, nil
}

func readValue(ctx *JsonContex, state int64) (JsonValue, error) {
	if err := readWhite(ctx); err != nil {
		return JsonObject{}, errors.Join(err)
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
		return BooleanValue{Value: true}, nil
	}
	if next == "null" {
		ctx.Index += 4
		return NullValue{}, nil
	}
	next = ctx.subString(ctx.Index, ctx.Index+5)
	if next == "false" {
		ctx.Index += 5
		return BooleanValue{Value: false}, nil
	}
	return ReadNum(ctx)
}

func ReadNum(ctx *JsonContex) (NumValue, error) {
	str, err := readNumInner(ctx, 0)
	if err != nil {
		return NumValue{}, errors.Join(err)
	}
	return NumValue{Value: str}, nil
}

func readNumInner(ctx *JsonContex, state int64) (string, error) {
	start := ctx.Index
	defaultErr := errors.New("number end unexpected")
	for state != 7 {
		curr, err := ctx.getCurr()
		if err != nil {
			return "", defaultErr
		}
		if state == 0 {
			if curr == '-' {
				ctx.incr()
			}
			state = 1
			continue
		}
		if state == 1 {
			if curr == '0' {
				ctx.incr()
				state = 2
				continue
			} else if curr >= '1' && curr <= '9' {
				ctx.incr()
				state = 3
				continue
			} else {
				return "", defaultErr
			}
		}
		if state == 3 {
			for curr >= '0' && curr <= '9' {
				ctx.incr()
				curr, err = ctx.getCurr()
				if err != nil {
					return "", defaultErr
				}
			}
			state = 2
			continue
		}
		if state == 2 {
			if curr == '.' {
				ctx.incr()
				state = 4
			} else {
				state = 5
			}
			continue
		}
		if state == 4 {
			if curr < '0' && curr > '9' {
				return "", defaultErr
			}
			ctx.incr()
			for curr >= '0' && curr <= '9' {
				ctx.incr()
				curr, err = ctx.getCurr()
				if err != nil {
					return "", defaultErr
				}
			}
			state = 5
			continue
		}
		if state == 5 {
			if curr == 'e' || curr == 'E' {
    			ctx.incr()
				state = 6
			} else {
				state = 7
			}
			continue
		}
		if state == 6 {
			if curr != '+' && curr != '-' {
				return "", defaultErr
			}
			ctx.incr()
			curr, err = ctx.getCurr()
			if err != nil {
				return "", defaultErr
			}
			if curr < '0' && curr > '9' {
				return "", defaultErr
			}
			ctx.incr()
			curr, err = ctx.getCurr()
			if err != nil {
				return "", defaultErr
			}
			for curr >= '0' && curr <= '9' {
				ctx.incr()
				curr, err = ctx.getCurr()
				if err != nil {
					return "", defaultErr
				}
			}
			state = 7
			continue
		}
	}
	return ctx.subString(start, ctx.Index), nil
}

// walk through all the white space from given point,until a non white character show
func readWhite(ctx *JsonContex) error {
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

func readStr(ctx *JsonContex) (StringValue, error) {
	str, err := readString(ctx, 0, ctx.Index)
	if err != nil {
		return StringValue{}, errors.Join(err)
	}
	return StringValue{Value: str}, nil
}

// 状态机的状态，
// 0 初始状态
// 1 读取第一个"之后
// 2 遇到转义字符
// 3 遇到结束的"
// MAX 出错
func readString(ctx *JsonContex, state int64, start int64) (string, error) {
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

func SolveString(ctx *JsonContex) (string, error) {
	fmt.Println()
	str, error := readString(ctx, 0, ctx.Index)
	if error != nil {
		return "", errors.New("error while try solving string")
	}
	return str, nil
}
