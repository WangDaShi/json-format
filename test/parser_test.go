package test

import (
"testing"
    "json/format/parse"
)

func jsonFileTest(t *testing.T)() {
	path := "D://code/go/json-format/test.json"
	path2, err := parse.ParseFile(path)
	if err != nil {
		panic(err)
	} else {
		println(path2)
	}
}

