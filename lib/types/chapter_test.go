package types

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestCompareChap1(t *testing.T) {
	c1sub1 := CreateSubchapter("Chap1Sub1", "sub1", "../../misc/test/sub1", "1. ", 10, []Subsub{})
	c1 := CreateChapter("Chap1", "chap1", "../../misc/test/chap1", "I. ", 10, []Subchapter{c1sub1})
	err, fails := c1.Compare(C1)
	if err != nil {
		t.Error(strings.Join(fails, "\n"))
	}
}

func TestCompareChap11(t *testing.T) {
	ssub := CreateSubsub("Chap1Sub2ExtA", "extA", "../../misc/test/sub4", "a. ", 10)
	c1sub1 := CreateSubchapter("Chap1Sub1", "sub1", "../../misc/test/sub1", "1. ", 10, []Subsub{})
	c1sub2 := CreateSubchapter("Chap1Sub2", "sub2", "../../misc/test/sub2", "2. ", 20, []Subsub{ssub})
	c1 := CreateChapter("Chap1", "chap1", "../../misc/test/chap1", "I. ", 10, []Subchapter{c1sub1, c1sub2})
	err, fails := c1.Compare(C12)
	if err != nil {
		t.Error(strings.Join(fails, "\n"))
	}
}

func Test_CreateChapterDeepEqual(t *testing.T) {
	sub1 := CreateSubchapter("Chap1Sub1", "sub1", "../../misc/test/sub1", "1. ", 10, []Subsub{})
	got := CreateChapter("Chap1", "chap1", "../../misc/test/chap1", "I. ", 10, []Subchapter{sub1})
	exp := C1
	if !reflect.DeepEqual(got, exp) {
		t.Errorf("\n# Got: %v\n# Exp: %v", got, exp)
	}
}

func Test_CreateChapter(t *testing.T) {
	c1 := CreateChapter("Chapter1", "chap1", "./test/chap1", "I. ", 10, []Subchapter{C1s1, C1s2})
	fails := []string{}
	if c1.Title != "Chapter1" {
		fails = append(fails, fmt.Sprintf("Title >> Exp:'%s' // Got: '%s'", "Chapter1", c1.Title))
	}
	if c1.Path != "chap1" {
		fails = append(fails, fmt.Sprintf("Path >> Exp:'%s' // Got: '%s'", "chap1", c1.Path))
	}
	if c1.Source != "./test/chap1" {
		fails = append(fails, fmt.Sprintf("Source >> Exp:'%s' // Got: '%s'", "./test/chap1", c1.Source))
	}
	if c1.Enum != "I. " {
		fails = append(fails, fmt.Sprintf("Enum >> Exp:'%s' // Got: '%s'", "I. ", c1.Enum))
	}
	if c1.Weight != 10 {
		fails = append(fails, fmt.Sprintf("Weight >> Exp:'%d' // Got: '%d'", 10, c1.Weight))
	}
	if len(fails) > 0 {
		fmt.Printf(strings.Join(fails, "\n"))
		t.Error()
	}
}
