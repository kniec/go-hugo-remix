package types

import (
	"fmt"
	"strings"
)

// Base holds the common variables
type Base struct {
	Title      string `yaml:"title"`
	Path       string `yaml:"path"`
	Source     string `yaml:"source"`
	Weight     int    `yaml:"weight"`
	IncludeTOC bool   `yaml:"include_toc"`
	Enum       string `yaml:"enum"`
	Author     string
	Flavour    string
}

// Compare checks each value to be equal
func (self *Base) Compare(other Base) (err error, fails []string) {
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

// CreateBase returns a filled Base
func CreateBase(t, p, s, e string, w int) Base {
	b := Base{
		Title:   t,
		Path:    p,
		Source:  s,
		Weight:  w,
		Enum:    e,
		Flavour: "eng",
	}
	return b
}
