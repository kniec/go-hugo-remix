package types

import (
	"fmt"
	"strings"
)

// Chapter is the highest level content
// Chapter -> Subchap -> Subsub
type Chapter struct {
	Title      string `yaml:"title"`
	Path       string `yaml:"path"`
	Source     string `yaml:"source"`
	Weight     int    `yaml:"weight"`
	IncludeTOC bool   `yaml:"include_toc"`
	Enum       string `yaml:"enum"`
	Author     string
	Flavour    string
	Subchaps   []Subchapter `yaml:"subchaps"`
}

// CreateChapter build a chapter
func CreateChapter(t, p, s, e string, w int, sub []Subchapter) Chapter {
	//b := CreateBase(t, p, s, e, w)
	res := Chapter{
		Title:    t,
		Path:     p,
		Source:   s,
		Weight:   w,
		Enum:     e,
		Flavour:  "eng",
		Subchaps: sub,
	}
	return res
}

// Compare takes a second chapter and compares the two
func (self *Chapter) Compare(other Chapter) (err error, fails []string) {
	fails = compVal(self.Title, other.Title, fails)
	fails = compVal(self.Path, other.Path, fails)
	fails = compVal(self.Source, other.Source, fails)
	fails = compVal(self.Weight, other.Weight, fails)
	fails = compVal(self.IncludeTOC, other.IncludeTOC, fails)
	fails = compVal(self.Enum, other.Enum, fails)
	fails = compVal(self.Author, other.Author, fails)
	fails = compVal(self.Flavour, other.Flavour, fails)
	if len(fails) > 0 {
		err = fmt.Errorf(strings.Join(fails, "\n"))
	}
	if len(self.Subchaps) != len(other.Subchaps) {
		fails = append(fails, fmt.Sprintf("Subchaps have different length (self:%d | other:%d)", len(self.Subchaps), len(other.Subchaps)))
	}
	for i, schap := range self.Subchaps {
		e, sFails := schap.Compare(other.Subchaps[i])
		fails = append(fails, sFails...)
		err = e
	}
	return
}

func (self *Chapter) SprintTitle() string {
	return fmt.Sprintf("Chapter.Title: %s", self.Title)
}
func (self *Chapter) SprintPath() string {
	return fmt.Sprintf("Chapter.Path: %s", self.Path)
}
func (self *Chapter) SprintSource() string {
	return fmt.Sprintf("Chapter.Source: %s", self.Source)
}
func (self *Chapter) SprintWeight() string {
	return fmt.Sprintf("Chapter.Weight: %d", self.Weight)
}
func (self *Chapter) SprintIncludeTOC() string {
	return fmt.Sprintf("Chapter.IncludeTOC: %t", self.IncludeTOC)
}
func (self *Chapter) SprintEnum() string {
	return fmt.Sprintf("Chapter.Enum: %s", self.Enum)
}
func (self *Chapter) SprintAuthor() string {
	return fmt.Sprintf("Chapter.Author: %s", self.Author)
}
func (self *Chapter) SprintFlavour() string {
	return fmt.Sprintf("Chapter.Flavour: %s", self.Flavour)
}

func (self *Chapter) String() (res []string) {
	res = append(res, self.SprintTitle())
	res = append(res, self.SprintPath())
	res = append(res, self.SprintSource())
	res = append(res, self.SprintWeight())
	res = append(res, self.SprintIncludeTOC())
	res = append(res, self.SprintEnum())
	res = append(res, self.SprintAuthor())
	res = append(res, self.SprintFlavour())
	for _, schap := range self.Subchaps {
		res = append(res, schap.String()...)
	}
	return
}
