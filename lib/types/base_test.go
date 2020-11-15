package types

import (
	"reflect"
	"testing"
)

var (
	b1 = Base{
		Title:   "Workshop1",
		Path:    ".",
		Source:  "../../misc/test/workshop1",
		Weight:  1,
		Flavour: "eng",
		Enum:    "I. ",
	}
	c1s1eA = Subsub{
		Base: Base{
			Title:   "Chap1Sub1ExtA",
			Path:    "extA",
			Source:  "../../misc/test/sub4",
			Weight:  1000,
			Flavour: "eng",
			Enum:    "a. ",
		},
	}
	c1s1 = Subchapter{
		Base: Base{
			Title:   "Chap1Sub1",
			Path:    "sub1",
			Source:  "../../misc/test/sub1",
			Weight:  100,
			Flavour: "eng",
			Enum:    "1. ",
		},
		Subsub: []Subsub{c1s1eA},
	}
	c1 = Chapter{
		Base: Base{
			Title:   "Chap1",
			Path:    "chap1",
			Source:  "../../misc/test/chap1",
			Weight:  10,
			Flavour: "eng",
			Enum:    "I. ",
		},
		Subchap: []Subchapter{c1s1},
	}
)

func Test_CreateBase(t *testing.T) {
	got := CreateBase("Workshop1", ".", "../../misc/test/workshop1", "I. ", 1)
	if !reflect.DeepEqual(got, b1) {
		t.Errorf("CreateBase //\n# Got: %v\n# Exp: %v", got, b1)
	}
}

func Test_CreateSubsub(t *testing.T) {
	got := CreateSubsub("Chap1Sub1ExtA", "extA", "../../misc/test/sub4", "a. ", 1000)
	exp := c1s1eA
	if !reflect.DeepEqual(got, exp) {
		t.Errorf("CreateSubsub //\n# Got: %v\n# Exp: %v", got, exp)
	}
}

func Test_CreateSubchapter(t *testing.T) {
	ssub := CreateSubsub("Chap1Sub1ExtA", "extA", "../../misc/test/sub4", "a. ", 1000)
	got := CreateSubchapter("Chap1Sub1", "sub1", "../../misc/test/sub1", "1. ", 100, []Subsub{ssub})
	exp := c1s1
	if !reflect.DeepEqual(got, exp) {
		t.Errorf("\n# Got: %v\n# Exp: %v", got, exp)
	}
}

func Test_CreateChapter(t *testing.T) {
	ssub := CreateSubsub("Chap1Sub1ExtA", "extA", "../../misc/test/sub4", "a. ", 1000)
	sub := CreateSubchapter("Chap1Sub1", "sub1", "../../misc/test/sub1", "1. ", 100, []Subsub{ssub})
	got := CreateChapter("Chap1", "chap1", "../../misc/test/chap1", "I. ", 10, []Subchapter{sub})
	exp := c1
	if !reflect.DeepEqual(got, exp) {
		t.Errorf("\n# Got: %v\n# Exp: %v", got, exp)
	}
}
