package types

import "testing"

func Test_compVal(t *testing.T) {
	emptySlice := []string{}
	got := compVal(1, 1, emptySlice)
	if len(got) != 0 {
		t.Error(got)
	}
	got = compVal(1, 2, emptySlice)
	if len(got) != 1 {
		t.Error(got)
	}
}
