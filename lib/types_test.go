package redux

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var data = `
title: Test
`

var DataChap = `
title: Test
chaps:
  - chap1:
     title: Chapter1
`

var testYaml = Workshop{
	Title: "Workshop1",
}

func readFile(path string) ([]byte, error) {
	yamlFile, err := ioutil.ReadFile(path)
	return yamlFile, err
}

// Note: struct fields must be public in order for unmarshal to
// correctly populate the data.

func TestParse(t *testing.T) {
	w := Workshop{}
	w.Parse([]byte(data))
	exp := "Test"
	if w.GetTitle() != exp {
		log.Printf("%v+", w)
		t.Errorf("workshop.title should be '%s', was '%s'", exp, w.GetTitle())
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
		"Chap1Sub1", "sub1", "../misc/test/sub1", "1. ", 10, []string{},
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
		Source: "../misc/test/sub1",
		Prefix: []string{},
		Weight: 10,
		Enum:   "1. ",
	}
	c1sub2 := Subchapter{
		Title:  "Chap1Sub2",
		Path:   "sub2",
		Source: "../misc/test/sub2",
		Prefix: []string{},
		Weight: 20,
		Enum:   "2. ",
	}
	c1 := Chapter{}
	c1.Title = "Chapter1"
	c1.Path = "chap1"
	c1.Source = "../misc/test/chap1"
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
	w := CreateWorkshop("Title", "../misc/hugo", []Chapter{})
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
	exp = "../misc/hugo"
	if w.BaseURL != exp {
		log.Printf("%v+", w)
		t.Errorf("workshop.BaseURL should be '%s', was '%s'", exp, w.BaseURL)
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

// TestGenerateHugo will take the test.yaml and generate the hugo workshop
func TestGenerateHugo(t *testing.T) {
	_, w := CreateWorkshopFromFile("../misc/test.yaml")
	log.Printf("w.Chaps[0]:")
	fmt.Println(w.Chaps[0].String())
	res := w.GenerateHugo("/tmp")
	fmt.Printf(strings.Join(res, "\n"))
	t.Errorf("Wait")
}
