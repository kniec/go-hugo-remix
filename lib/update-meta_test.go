package redux

import (
	"fmt"
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

func TestUpdateMetaFile(t *testing.T) {
	t.Errorf("Wait")
}
