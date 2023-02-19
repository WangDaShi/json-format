package myjson

import (
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
	Value float32
}

func (v NumValue) ToJson(indent int) string {
	return fmt.Sprint(v.Value)
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
}

func (v JsonObject) ToJson(indent int) string {
	tab := v.Data
	s := "{"
	for k, v := range tab {
		s += "\n"
		s += genIndent(indent+1) + "\"" + k + "\" : "
		s += v.ToJson(indent + 1)
		s += ","
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
