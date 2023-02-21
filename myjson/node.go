package myjson

import (
	"errors"
	"fmt"
	"strings"
)

type StringValue struct {
	Value string
}

func (v StringValue) ToJson(indent int) string {
	return "\"" + v.Value + "\""
}

type NumValue struct {
	Value string
}

func (v NumValue) ToJson(indent int) string {
	return v.Value
}

type BooleanValue struct {
	Value bool
}

func (v BooleanValue) ToJson(indent int) string {
	if v.Value {
		return "true"
	} else {
		return "false"
	}
}

type NullValue struct {
}

func (v NullValue) ToJson(indent int) string {
	return "null"
}

type ArrayValue struct {
	Arr []JsonValue
}

func (v ArrayValue) ToJson(indent int) string {
	s := "["
	for _, v := range v.Arr {
		s += "\n"
		s += genIndent(indent+1) + v.ToJson(indent+1)
		s += ","
	}
	s += "\n" + genIndent(indent) + "]"
	return s
}

type JsonObject struct {
	Data map[string]JsonValue
	Keys []string
}

func (o *JsonObject) Add(key string, v *JsonValue) error {
    if _,ok := o.Data[key]; ok {
        return errors.New(fmt.Sprintf("key %s already exist !",key))
    }
	o.Keys = append(o.Keys, key)
	o.Data[key] = *v
    return nil
}

func (v JsonObject) ToJson(indent int) string {
	s := "{"
	l := len(v.Keys)
	for i := 0; i < l; i++ {
        key := v.Keys[i]
        value := v.Data[key].ToJson(indent + 1)
        s += fmt.Sprintf("\n%s\"%s\" : %s",genIndent(indent+1),key,value)
        if i < l -1 {
            s += ","
        }
	}
	s += "\n" + genIndent(indent) + "}"
	return s
}

type JsonValue interface {
	ToJson(indent int) string
}

func genIndent(indent int) string {
	s := "    "
	return strings.Repeat(s, indent)
}
