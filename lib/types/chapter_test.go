package types

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
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
	c1.SetDebugLevel(3)
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

func Test_CopyContent_SimpleChapter(t *testing.T) {
	c1 := CreateChapter("Chapter1", "chap1", "./test/chap1", "I. ", 10, []Subchapter{})
	c1.SetDebugLevel(3)
	tdir, _ := ioutil.TempDir(os.TempDir(), "chap-")
	fmt.Printf("CopyContent(., %s)\n  > Path: %s\n  > Source: %s\n", tdir, c1.Path, c1.Source)
	c1.CopyContent("../../misc", []string{tdir})
	// Check if path exists
	for _, cpath := range []string{"content/chap2/_index.md"} {
		if _, err := os.Stat(path.Join(tdir, cpath)); os.IsNotExist(err) {
			t.Errorf("Path '%s' should exist after creating Hugo in '%s'", cpath, tdir)
			return
		}
	}
	err := os.RemoveAll(tdir)
	if err != nil {
		fmt.Println(err.Error())
		t.Error()
	}
}

func Test_CopyContent_StaticChapter(t *testing.T) {
	c := CreateChapter("Chapter2", "chap2", "./test/chap2", "III. ", 30, []Subchapter{})
	c.SetDebugLevel(3)
	tdir, _ := ioutil.TempDir(os.TempDir(), "chap-")
	c.CopyContent("../../misc", []string{tdir})
	// Check if path exists
	for _, cpath := range []string{"content/chap2/_index.md", "static/images/wikipedia.png"} {
		if _, err := os.Stat(path.Join(tdir, cpath)); os.IsNotExist(err) {
			t.Errorf("Path '%s' should exist after creating Hugo in '%s'", cpath, tdir)
			return
		}
	}
	err := os.RemoveAll(tdir)
	if err != nil {
		fmt.Println(err.Error())
		t.Error()
	}
}
