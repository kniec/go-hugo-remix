package types

import (
	"log"
	"strings"
	"testing"
)

/**********
Test Data
*/
var shortData = `
title: Test
description: "Super workshop"
`

var dataChap = `
title: Test
description: "Super workshop"
chaps:
  - chap1:
    title: Chapter1
`
var c1sub1extA = CreateSubsub("Chap1Sub1ExtA", "extA", "./test/sub4", "a. ", 10)
var testC1sub2 = CreateSubchapter("Chap1Sub2", "sub2", "./test/sub2", "2. ", 20, []Subsub{c1sub1extA})
var testC1sub1 = CreateSubchapter("Chap1Sub1", "sub1", "./test/sub1", "1. ", 10, []Subsub{})
var testC2sub1 = CreateSubchapter("Chap2Sub1", "sub1", "./test/sub3", "1. ", 10, []Subsub{})
var testChap1 = CreateChapter("Chapter1", "chap1", "./test/chap1", "I. ", 10, []Subchapter{testC1sub1, testC1sub2})
var testChap2 = CreateChapter("Chapter2", "chap2", "./test/chap2", "II. ", 20, []Subchapter{testC2sub1})

var testWorkshop = Workshop{
	Title:       "Workshop1",
	Author:      "Christian Kniep",
	HugoBase:    "./hugo",
	Description: "Workshop description",
	Flavours:    []string{"eng", "rse", "sys"},
	Source:      "./test/workshop1",
	Chaps:       []Chapter{testChap1, testChap2},
}

// Note: struct fields must be public in order for unmarshal to
// correctly populate the data.
func Test_Parse(t *testing.T) {
	w := Workshop{}
	w.Parse([]byte(shortData))
	exp := "Test"
	if w.Title != exp {
		t.Errorf("workshop.title should be '%s', was '%s'", exp, w.Title)
	}
	exp = "Super workshop"
	if w.Description != exp {
		t.Errorf("workshop.Description workshop should be '%s', was '%s'", exp, w.Description)
	}
}

func Test_ParseData2(t *testing.T) {
	w := Workshop{}
	w.Parse([]byte(dataChap))
	exp := "Test"
	if w.Title != exp {
		log.Printf("%v+", w)
		t.Errorf("workshop.title should be '%s', was '%s'", exp, w.Title)
	}
}

func TestParseFile_Title(t *testing.T) {
	_, w := CreateWorkshopFromFile("../../misc/workshop-c1s1eA.yaml")
	exp := "Workshop1"
	if w.Title != exp {
		t.Errorf("workshop.Title should be '%s', was '%s'", exp, w.Title)
	}
	if len(w.Chaps) != 1 {
		t.Errorf("workshop.Chaps should have 1 chapters, it has '%d'", len(w.Chaps))
	}
	for _, chap := range w.Chaps {
		if len(chap.Subchaps) != 1 {
			t.Errorf("workshop.Chap1.Subchaps should have 1 Subchapter, it has '%d'", len(chap.Subchaps))
		}
		exp := "Chapter1"
		got := chap.Title
		if got != exp {
			t.Errorf("chapter1.Title should be '%s', was '%s'", exp, got)
		}
		for _, schap := range chap.Subchaps {
			if len(schap.Subsubs) != 1 {
				t.Errorf("workshop.Chap1.Sub1.Subsubs should have 1 Subchapter, it has '%d'", len(chap.Subchaps))
			}
		}
	}
}

func TestParseFile_CompareC1s1eA(t *testing.T) {
	yData, err := readFile("../../misc/workshop-c1s1eA.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	w := Workshop{}
	w.Parse(yData)
	exp := testWorkshop
	err, fails := w.Compare(exp)
	if err != nil {
		t.Error(strings.Join(fails, "\n"))
	}
}
