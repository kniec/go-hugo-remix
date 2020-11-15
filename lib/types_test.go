package redux

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var data = `
title: Test
description: "Super workshop"
`

var DataChap = `
title: Test
description: "Super workshop"
chaps:
  - chap1:
    title: Chapter1
`

var testYaml = Workshop{
	Title:       "Workshop1",
	Description: "Super workshop",
}

// Note: struct fields must be public in order for unmarshal to
// correctly populate the data.
func TestParse(t *testing.T) {
	w := Workshop{}
	w.Parse([]byte(data))
	exp := "Test"
	if w.GetTitle() != exp {
		t.Errorf("workshop.title should be '%s', was '%s'", exp, w.GetTitle())
	}
	exp = "Super workshop"
	if w.Description != exp {
		t.Errorf("workshop.Description workshop should be '%s', was '%s'", exp, w.Description)
	}
}

func TestParseData2(t *testing.T) {
	w := Workshop{}
	w.Parse([]byte(DataChap))
	exp := "Test"
	if w.GetTitle() != exp {
		log.Printf("%v+", w)
		t.Errorf("workshop.title should be '%s', was '%s'", exp, w.GetTitle())
	}
}

func TestCompateSubchap(t *testing.T) {
	s := CreateSubchapter(
		"Chap1Sub1", "sub1", "./test/sub1", "1. ", 10, []Subsub{},
	)
	err, _ := s.CompareSubchap(testC1sub1)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestCompateChap1(t *testing.T) {
	c1sub1 := Subchapter{
		Title:  "Chap1Sub1",
		Path:   "sub1",
		Source: "./test/sub1",
		Prefix: []string{},
		Weight: 10,
		Enum:   "1. ",
	}
	c1sub2 := Subchapter{
		Title:  "Chap1Sub2",
		Path:   "sub2",
		Source: "./test/sub2",
		Prefix: []string{},
		Weight: 20,
		Enum:   "2. ",
	}
	c1 := Chapter{}
	c1.Title = "Chapter1"
	c1.Path = "chap1"
	c1.Source = "./test/chap1"
	c1.Prefix = []string{}
	c1.Weight = 10
	c1.Enum = "I. "
	c1.Subchap = []Subchapter{c1sub1, c1sub2}
	err, _ := c1.CompareChap(testChap1)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestCreateWorkshop(t *testing.T) {
	w := CreateWorkshop("Title", "Description", "../misc/hugo", "./test/workshop1", []Chapter{})
	if w.Title != "Title" && len(w.Chaps) == 0 {
		t.Errorf("Title should be 'Title' w/o chapters")
	}
}

func TestParseDate_TitleBase(t *testing.T) {
	yData, err := readFile("../misc/test.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	w := Workshop{}
	w.Parse(yData)
	exp := "Workshop1"
	if w.GetTitle() != exp {
		log.Printf("%v+", w)
		t.Errorf("workshop.Title should be '%s', was '%s'", exp, w.GetTitle())
	}
	exp = "./hugo"
	if w.HugoBase != exp {
		log.Printf("%v+", w)
		t.Errorf("workshop.HugoBase should be '%s', was '%s'", exp, w.HugoBase)
	}
}

func TestParseFile_Title(t *testing.T) {
	_, w := CreateWorkshopFromFile("../misc/test.yaml")
	exp := "Workshop1"
	if w.GetTitle() != exp {
		log.Printf("%v+", w)
		t.Errorf("workshop.Title should be '%s', was '%s'", exp, w.GetTitle())
	}
}

func TestParseFile_EqualTestYAML(t *testing.T) {
	yData, err := readFile("../misc/test.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	w := Workshop{}
	w.Parse(yData)
	exp := GetTestWorkshop()
	if cmp.Equal(w, exp) {
		fmt.Printf("CUR: %v+\n", w)
		fmt.Printf("EXP: %v+\n", exp)
		t.Errorf("workshop and test_workshop are not equal")
	}
}

func TestParseFile_CompareTestYAML(t *testing.T) {
	yData, err := readFile("../misc/test.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	w := Workshop{}
	w.Parse(yData)
	exp := GetTestWorkshop()

	err, fails := w.CompareWorkshops(exp)
	if err != nil {
		t.Errorf(strings.Join(fails, "\n"))
	}
}

// Extending test-ext.yaml with test.yaml
func Test_ExtendFromWorkshop(t *testing.T) {
	_, wExt := CreateWorkshopFromFile("../misc/test-ext.yaml")
	_, w := CreateWorkshopFromFile("../misc/test.yaml")
	wExt.ExtendFromWorkshop(w)
	c1In := false
	c2In := false
	for _, c := range wExt.Chaps {
		if c.Path == "chap1" {
			c1In = true
		}
		if c.Path == "chap2" {
			c2In = true
		}
	}
	if !c1In || !c2In {
		t.Error("chap3 should be in ")
	}
}

// TestGenerateHugo will take the test.yaml and generate the hugo workshop
func TestGenerateHugo(t *testing.T) {
	_, w := CreateWorkshopFromFile("../misc/test.yaml")
	tdir, _ := ioutil.TempDir(os.TempDir(), "hugo-")
	_, res := w.GenerateHugo(tdir)
	fmt.Printf(strings.Join(res, "\n"))
	for _, cpath := range []string{"content/chap1/sub2/extA/_index.md"} {
		if _, err := os.Stat(path.Join(tdir, cpath)); os.IsNotExist(err) {
			t.Errorf("Path '%s' should exist after creating Hugo in '%s'", cpath, tdir)
			return
		}
	}
	if false {
		log.Printf("TempDir: %s", tdir)
	} else {
		err := os.RemoveAll(tdir)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// TestUpdateChapIndex will create a Hugo Workshop and
// update the Chapter index files based on the workshop YAML
func TestUpdateIndex(t *testing.T) {
	_, w := CreateWorkshopFromFile("../misc/test.yaml")
	tdir, _ := ioutil.TempDir(os.TempDir(), "hugo-")
	log.Printf("TempDir: %s", tdir)
	w.GenerateHugo(tdir)
	// CleanUP
	c := Checker{}
	fpath := path.Join(tdir, "content/chap2/_index.md")
	err := c.ReadMeta(fpath)
	if err != nil {
		t.Errorf(err.Error())
	}
	if c.Meta["title"] != "Chapter2" {
		t.Errorf("Wrong title: '%s'", c.Meta["title"])
	}
	c2c1 := Checker{}
	fpath = path.Join(tdir, "content/chap2/sub1/_index.md")
	err = c2c1.ReadMeta(fpath)
	if err != nil {
		t.Errorf(err.Error())
	}
	if c2c1.Meta["title"] != "Chap2Sub1" {
		t.Errorf("Wrong title: %s", c2c1.Meta["title"])
	}
	// CleanUP
	err = os.RemoveAll(tdir)
	if err != nil {
		t.Errorf(err.Error())
	}
}

// BASE META
func TestUpdateDict(t *testing.T) {
	bm := BaseMeta{
		Title:      "Title1",
		Weight:     10,
		Chapter:    true,
		Pre:        "1. ",
		IncludeTOC: true,
	}
	dict := mType{
		"title":       "TitleDict",
		"weight":      20,
		"chapter":     true,
		"pre":         "I. ",
		"include_toc": false,
	}
	_, cur := bm.UpdateDict(dict)
	expStr := "Title1"
	curStr := cur["title"]
	if expStr != curStr {
		t.Errorf("Title: oldDict: '%s', bm: '%s', newDict: '%s'", dict["title"], bm.Title, curStr)
	}
	expInt := 10
	curInt := cur["weight"]
	if expInt != curInt {
		t.Errorf("Weight: oldDict: '%d', bm: '%d', newDict: '%d'", dict["weight"], bm.Weight, curInt)
	}
	expStr = "1. "
	curStr = cur["pre"]
	if expStr != curStr {
		t.Errorf("Pre: oldDict: '%s', bm: '%s', newDict: '%s'", dict["pre"], bm.Pre, curStr)
	}
	expBool := true
	curBool := cur["include_toc"]
	if expBool != curBool {
		t.Errorf("IncludeTOC: oldDict: '%t', bm: '%t', newDict: '%t', expected: '%t'", dict["include_toc"], bm.IncludeTOC, curBool, expBool)
	}
}

func TestToStrings(t *testing.T) {
	dict := mType{
		"title":       "TitleDict",
		"weight":      20,
		"chapter":     true,
		"pre":         "I. ",
		"include_toc": false,
	}
	exp := []string{
		`title: "TitleDict"`,
		`weight: 20`,
		`chapter: true`,
		`pre: "I. "`,
		`include_toc: false`,
	}
	cur := dict.ToStrings()
	sort.Strings(exp)
	sort.Strings(cur)
	for i := range cur {
		if cur[i] != exp[i] {
			t.Errorf("Exp: '%v' / Cur '%v'", exp[i], cur[i])
		}
	}
}
