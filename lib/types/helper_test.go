package types

import (
	"fmt"
	"strings"
	"testing"
)

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
