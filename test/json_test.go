package test

import (
    "testing"
    "json/format/myjson"
)

func TestJson1(t *testing.T) {
    m := make(map[string]myjson.JsonValue)
    m["key1"] = myjson.StringValue{Value : "lalala"}
    m["key2"] = myjson.BooleanValue{Value : true}
    m["key3"] = myjson.BooleanValue{Value : false}
    m["key4"] = myjson.NullValue{}
    m["key5"] = myjson.NumValue{Value : 123}
    m["key6"] = myjson.NumValue{Value : 123.456}

    k := make(map[string]myjson.JsonValue)
    k["key1"] = myjson.StringValue{Value : "lalala"}
    k["key2"] = myjson.BooleanValue{Value : true}
    k["key3"] = myjson.BooleanValue{Value : false}
    k["key4"] = myjson.NullValue{}
    k["key5"] = myjson.NumValue{Value : 123}
    k["key6"] = myjson.NumValue{Value : 123.456}

    m["key7"] = myjson.JsonObject{Data : k}

    var j []myjson.JsonValue
    j = append(j, myjson.StringValue{Value : "lalala"})
    j = append(j, myjson.BooleanValue{Value : true})
    j = append(j, myjson.BooleanValue{Value : false})
    j = append(j, myjson.NullValue{})
    j = append(j, myjson.NumValue{Value : 123})
    j = append(j, myjson.NumValue{Value : 123.456})
    j = append(j, myjson.JsonObject{Data : k})

    m["key8"] = myjson.ArrayValue{Arr : j}

    s := myjson.JsonObject{Data : m}
    println(s.ToJson(0))
}

