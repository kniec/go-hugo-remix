package redux

import (
	"fmt"
	"strings"
	"testing"
)

func TestReadMeta(t *testing.T) {
	c := Checker{}
	fpath := "../misc/test/chap1/_index.md"
	err := c.ReadMeta(fpath)
	if err != nil {
		t.Errorf(err.Error())
	}
	cur := c.Meta["title"]
	exp := "Chapter1"
	if cur != exp {
		t.Errorf(fmt.Sprintf("Title should be '%s'; found '%s", exp, cur))
	}
	curInt := c.Meta["weight"]
	expInt := 10
	if curInt != expInt {
		t.Errorf(fmt.Sprintf("Weight should be '%d'; found '%d", expInt, curInt))
	}
	curBool := c.Meta["chapter"]
	expBool := true
	if curBool != expBool {
		t.Errorf(fmt.Sprintf("Chapter should be '%t'; found '%t", expBool, curBool))
	}
	curBool = c.Meta["include_toc"]
	expBool = true
	if curBool != expBool {
		t.Errorf(fmt.Sprintf("include_toc should be '%t'; found '%t", expBool, curBool))
	}
	cur = c.Meta["pre"]
	exp = "<b>I. </b>"
	if cur != exp {
		t.Errorf(fmt.Sprintf("Prefix should be '%s'; found '%s", exp, cur))
	}
}

func TestUpdateMeta(t *testing.T) {
	c := Checker{}
	fpath := "../misc/test/sub1/_index.md"
	err := c.ReadMeta(fpath)
	if err != nil {
		t.Errorf(err.Error())
	}
	bm := BaseMeta{
		Title:      "My Sub1",
		Weight:     20,
		Chapter:    false,
		Pre:        "2. ",
		IncludeTOC: true,
	}
	_, cur := bm.UpdateDict(c.Meta)
	expStr := "My Sub1"
	curStr := cur["title"]
	if expStr != curStr {
		t.Errorf("Title: oldDict: '%s', bm: '%s', newDict: '%s'", c.Meta["title"], bm.Title, curStr)
	}
	expStr = "2. "
	curStr = cur["pre"]
	if expStr != curStr {
		t.Errorf("Pre: oldDict: '%s', bm: '%s', newDict: '%s'", c.Meta["pre"], bm.Title, curStr)
	}
	expInt := 20
	curInt := cur["weight"]
	if expInt != curInt {
		t.Errorf("Weight: oldDict: '%d', bm: '%d', newDict: '%d'", c.Meta["weight"], bm.Weight, curInt)
	}
}
func TestUpdateMetaLines(t *testing.T) {
	inputText := `---
title: "Chapter1"
weight: 10
chapter: true
pre: "<b>I. </b>"
include_toc: true
---

### Chapter 1
`
	metaLines := []string{
		`title: "My Chapter1"`,
		"weight: 20",
		`chapter: true`,
		`pre: "<b>1. </b>"`,
		`include_toc: true`,
	}
	inputLines := strings.Split(inputText, "\n")
	err, outputLines := updateMetaLines(inputLines, metaLines)
	if err != nil {
		t.Errorf(err.Error())
	}
	exp := []string{
		"---",
		`title: "My Chapter1"`,
		`weight: 20`,
		`chapter: true`,
		`pre: "<b>1. </b>"`,
		`include_toc: true`,
		`---`,
		"",
		"### Chapter 1",
		"",
	}
	if len(outputLines) != len(exp) {
		t.Errorf("len expected %d != %d len output", len(exp), len(outputLines))
	}
	fails := false
	for i, v := range outputLines {
		if v != exp[i] {
			fmt.Printf(" [%d] NOK got: %s // exp: %s\n", i, v, exp[i])
			fails = true
		} else {
			fmt.Printf(" [%d]  OK got: %s // exp: %s\n", i, v, exp[i])
		}
	}
	if fails {
		t.Errorf("Something went wrong")
	}
}

func TestToMetaLines(t *testing.T) {
	c := Checker{}
	fpath := "../misc/test/sub1/_index.md"
	err := c.ReadMeta(fpath)
	if err != nil {
		t.Errorf(err.Error())
	}
	bm := BaseMeta{
		Title:      "My Sub1",
		Weight:     20,
		Chapter:    false,
		Pre:        "2. ",
		IncludeTOC: true,
	}
	exp := []string{
		`title: "My Sub1"`,
		"weight: 20",
		`chapter: false`,
		`pre: "2. "`,
		`include_toc: true`,
	}
	c.UpdateMeta(bm)
	got := c.ToMetaLines()
	if len(got) != len(exp) {
		t.Errorf("len expected %d != %d len output", len(exp), len(got))

	}
	fails := false
	for i, v := range got {
		if v != exp[i] {
			fmt.Printf(" [%d] NOK got: %s // exp: %s\n", i, v, exp[i])
			fails = true
		} else {
			fmt.Printf(" [%d]  OK got: %s // exp: %s\n", i, v, exp[i])
		}
	}
	if fails {
		t.Errorf("Something went wrong")
	}
}
