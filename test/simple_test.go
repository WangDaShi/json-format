package test

import "testing"

func testSliceTestSlice(t *testing.T) {
    k := make([]string, 1)
    k = append(k,"1")
    k = append(k,"1")
    k = append(k,"1")
    k = append(k,"1")

    for _,v := range k {
        print(v + ",")
    }

}

