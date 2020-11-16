package types

import (
	"reflect"
	"strings"
	"testing"
)

func Test_CreateSubsub(t *testing.T) {
	got := CreateSubsub("Chap1Sub2ExtA", "extA", "../../misc/test/sub4", "a. ", 10)
	exp := C1s2eA
	if !reflect.DeepEqual(got, exp) {
		t.Errorf("CreateSubsub //\n# Got: %v\n# Exp: %v", got, exp)
	}
}

func Test_Compare(t *testing.T) {
	got := CreateSubsub("Chap1Sub2ExtA", "extA", "../../misc/test/sub4", "a. ", 10)
	exp := C1s2eA
	err, fails := exp.Compare(got)
	if err != nil {
		t.Error(strings.Join(fails, "\n"))
	}
}
