package main

import (
    "json/format/myjson"
)

func main() {
    m := make(map[string]myjson.JsonValue)
    m["key1"] = myjson.StringValue{Value : "lalala"}
    s := myjson.JsonObject{Data : m}
    println(s.ToJson(0))
}

