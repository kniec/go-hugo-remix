package types

import (
	"reflect"
	"strings"
	"testing"
)

func Test_CompareSubchap(t *testing.T) {
	s := CreateSubchapter("Chap1Sub1", "sub1", "../../misc/test/sub1", "1. ", 10, []Subsub{})
	err, fails := s.Compare(C1s1)
	if err != nil {
		t.Error(strings.Join(fails, "\n"))
	}
}

func Test_CompareSubchapWithSubsub(t *testing.T) {
	s := CreateSubchapter("Chap1Sub2", "sub2", "../../misc/test/sub2", "2. ", 20, []Subsub{C1s2eA})
	err, fails := s.Compare(C1s2)
	if err != nil {
		t.Error(strings.Join(fails, "\n"))
	}
}

func Test_CreateSubchapter(t *testing.T) {
	ssub := CreateSubsub("Chap1Sub2ExtA", "extA", "../../misc/test/sub4", "a. ", 10)
	got := CreateSubchapter("Chap1Sub2", "sub2", "../../misc/test/sub2", "2. ", 20, []Subsub{ssub})
	exp := C1s2
	if !reflect.DeepEqual(got, exp) {
		t.Errorf("\n# Got: %v\n# Exp: %v", got, exp)
	}
}
