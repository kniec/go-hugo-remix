package types

import (
	"fmt"
	"strings"
)

// Subsub is the lowest level content
// Chapter -> Subchap -> Subsub
type Subsub struct {
	Title      string `yaml:"title"`
	Path       string `yaml:"path"`
	Source     string `yaml:"source"`
	Weight     int    `yaml:"weight"`
	IncludeTOC bool   `yaml:"include_toc"`
	Enum       string `yaml:"enum"`
	Author     string
	Flavour    string
}

// CreateSubsub build a Subsub
func CreateSubsub(t, p, s, e string, w int) Subsub {
	res := Subsub{
		Title:   t,
		Path:    p,
		Source:  s,
		Weight:  w,
		Enum:    e,
		Flavour: "eng",
	}
	return res
}

// Compare takes a second subchapter and compares the two
func (self *Subsub) Compare(other Subsub) (err error, fails []string) {
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
	return
}

func (self *Subsub) SprintTitle() string {
	return fmt.Sprintf("Subsub.Title: %s", self.Title)
}
func (self *Subsub) SprintPath() string {
	return fmt.Sprintf("Subsub.Path: %s", self.Path)
}
func (self *Subsub) SprintSource() string {
	return fmt.Sprintf("Subsub.Source: %s", self.Source)
}
func (self *Subsub) SprintWeight() string {
	return fmt.Sprintf("Subsub.Weight: %d", self.Weight)
}
func (self *Subsub) SprintIncludeTOC() string {
	return fmt.Sprintf("Subsub.IncludeTOC: %t", self.IncludeTOC)
}
func (self *Subsub) SprintEnum() string {
	return fmt.Sprintf("Subsub.Enum: %s", self.Enum)
}
func (self *Subsub) SprintAuthor() string {
	return fmt.Sprintf("Subsub.Author: %s", self.Author)
}
func (self *Subsub) SprintFlavour() string {
	return fmt.Sprintf("Subsub.Flavour: %s", self.Flavour)
}

func (self *Subsub) String() (res []string) {
	res = append(res, self.SprintTitle())
	res = append(res, self.SprintPath())
	res = append(res, self.SprintSource())
	res = append(res, self.SprintWeight())
	res = append(res, self.SprintIncludeTOC())
	res = append(res, self.SprintEnum())
	res = append(res, self.SprintAuthor())
	res = append(res, self.SprintFlavour())
	return
}
