package helper

import (
	"sort"
	"strings"
)

type StrSet struct {
	List []string
}

func NewStrSet() StrSet {
	return StrSet{}
}

func (ss *StrSet) Add(e string) {
	// Super expensive to do this way, but....
	if !ss.Contains(e) {
		ss.List = append(ss.List, e)
	}

}

func (ss *StrSet) Contains(e string) bool {
	for _, v := range ss.List {
		if v == e {
			return true
		}
	}
	return false
}

func (ss *StrSet) String() (res []string) {
	sort.Strings(ss.List)
	return []string{strings.Join(ss.List, ",")}
}
