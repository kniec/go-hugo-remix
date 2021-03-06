package types

import (
	"fmt"
	"strings"
)

// Subchapter is the mid-level content
// Chapter -> Subchap -> Subsub
type Subchapter struct {
	Title      string `yaml:"title"`
	Path       string `yaml:"path"`
	Source     string `yaml:"source"`
	Weight     int    `yaml:"weight"`
	IncludeTOC bool   `yaml:"include_toc"`
	Enum       string `yaml:"enum"`
	Author     string
	Flavour    string
	Subsubs    []Subsub
}

// CreateSubchap build a Subchapter
func CreateSubchapter(t, p, s, e string, w int, subsub []Subsub) Subchapter {
	res := Subchapter{
		Title:   t,
		Path:    p,
		Source:  s,
		Weight:  w,
		Enum:    e,
		Flavour: "eng",
		Subsubs: subsub,
	}
	return res
}

// Compare takes a second subchapter and compares the two
func (self *Subchapter) Compare(other Subchapter) (err error, fails []string) {
	fails = compVal(self.Title, other.Title, fails)
	fails = compVal(self.Path, other.Path, fails)
	fails = compVal(self.Source, other.Source, fails)
	fails = compVal(self.Weight, other.Weight, fails)
	fails = compVal(self.IncludeTOC, other.IncludeTOC, fails)
	fails = compVal(self.Enum, other.Enum, fails)
	fails = compVal(self.Author, other.Author, fails)
	fails = compVal(self.Flavour, other.Flavour, fails)
	if len(self.Subsubs) != len(other.Subsubs) {
		fails = append(fails, fmt.Sprintf("Subsubs have different length (self:%d | other:%d)", len(self.Subsubs), len(other.Subsubs)))
	}
	for i, subs := range self.Subsubs {
		e, sFails := subs.Compare(other.Subsubs[i])
		fails = append(fails, sFails...)
		err = e
	}
	if len(fails) > 0 {
		err = fmt.Errorf(strings.Join(fails, "\n"))
	}
	return
}

func (self *Subchapter) SprintTitle() string {
	return fmt.Sprintf("Subchap.Title: %s", self.Title)
}
func (self *Subchapter) SprintPath() string {
	return fmt.Sprintf("Subchap.Path: %s", self.Path)
}
func (self *Subchapter) SprintSource() string {
	return fmt.Sprintf("Subchap.Source: %s", self.Source)
}
func (self *Subchapter) SprintWeight() string {
	return fmt.Sprintf("Subchap.Weight: %d", self.Weight)
}
func (self *Subchapter) SprintIncludeTOC() string {
	return fmt.Sprintf("Subchap.IncludeTOC: %t", self.IncludeTOC)
}
func (self *Subchapter) SprintEnum() string {
	return fmt.Sprintf("Subchap.Enum: %s", self.Enum)
}
func (self *Subchapter) SprintAuthor() string {
	return fmt.Sprintf("Subchap.Author: %s", self.Author)
}
func (self *Subchapter) SprintFlavour() string {
	return fmt.Sprintf("Subchap.Flavour: %s", self.Flavour)
}

func (self *Subchapter) String() (res []string) {
	res = append(res, self.SprintTitle())
	res = append(res, self.SprintPath())
	res = append(res, self.SprintSource())
	res = append(res, self.SprintWeight())
	res = append(res, self.SprintIncludeTOC())
	res = append(res, self.SprintEnum())
	res = append(res, self.SprintAuthor())
	res = append(res, self.SprintFlavour())
	for _, ssub := range self.Subsubs {
		res = append(res, ssub.String()...)
	}
	return
}
