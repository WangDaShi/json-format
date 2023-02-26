package test

import (
    "testing"
    "json/format/parse"
)

func TestJson1(t *testing.T) {
    m := make(map[string]parse.JsonValue)
    m["key1"] = parse.StringValue{Value : "lalala"}
    m["key2"] = parse.BooleanValue{Value : true}
    m["key3"] = parse.BooleanValue{Value : false}
    m["key4"] = parse.NullValue{}
    // m["key5"] = parse.NumValue{Value : 123}
    // m["key6"] = parse.NumValue{Value : 123.456}

    k := make(map[string]parse.JsonValue)
    k["key1"] = parse.StringValue{Value : "lalala"}
    k["key2"] = parse.BooleanValue{Value : true}
    k["key3"] = parse.BooleanValue{Value : false}
    k["key4"] = parse.NullValue{}
    // k["key5"] = parse.NumValue{Value : 123}
    // k["key6"] = parse.NumValue{Value : 123.456}

    m["key7"] = parse.JsonObject{Data : k}

    var j []parse.JsonValue
    j = append(j, parse.StringValue{Value : "lalala"})
    j = append(j, parse.BooleanValue{Value : true})
    j = append(j, parse.BooleanValue{Value : false})
    j = append(j, parse.NullValue{})
    // j = append(j, parse.NumValue{Value : 123})
    // j = append(j, parse.NumValue{Value : 123.456})
    j = append(j, parse.JsonObject{Data : k})

    m["key8"] = parse.ArrayValue{Arr : j}

    s := parse.JsonObject{Data : m}
    println(s.ToJson(0))
}

