package types

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func Test_CopyContentBasic(t *testing.T) {
	tdir, _ := ioutil.TempDir(os.TempDir(), "base-")

	b := Base{
		Source: "../../misc/workshop1",
		Path:   ".",
	}
	fmt.Printf("CopyContent(., %s)\n  > Path: %s\n  > Source: %s\n", tdir, b.Path, b.Source)
	err := b.CopyContent(".", []string{tdir})
	if err == nil {
		fmt.Printf("base.Source should not exist: %s", b.Source)
		t.Error()
	}
	err = os.RemoveAll(tdir)
	if err != nil {
		log.Fatal(err)
	}
}

func Test_CopyContent_SimpleContent(t *testing.T) {
	tdir, _ := ioutil.TempDir(os.TempDir(), "base-")

	b := Base{
		Source: "../../misc/test/chap1",
		Path:   "chap1",
	}
	fmt.Printf("CopyContent(., %s)\n  > Path: %s\n  > Source: %s\n", tdir, b.Path, b.Source)
	err := b.CopyContent(".", []string{tdir})
	if err != nil {
		t.Error(err.Error())
	}
	// Check if path exists
	for _, cpath := range []string{"content/chap1/_index.md"} {
		if _, err := os.Stat(path.Join(tdir, cpath)); os.IsNotExist(err) {
			t.Errorf("Path '%s' should exist after creating Hugo in '%s'", cpath, tdir)
			return
		}
	}
	err = os.RemoveAll(tdir)
	if err != nil {
		log.Fatal(err)
	}
}

func Test_CopyContent_StaticContent(t *testing.T) {
	tdir, _ := ioutil.TempDir(os.TempDir(), "base-")

	b := Base{
		dLevel: 3,
		Source: "../../misc/test/chap2",
		Path:   "chap2",
	}
	//fmt.Printf("CopyContent(., %s)\n  > Path: %s\n  > Source: %s\n", tdir, b.Path, b.Source)
	err := b.CopyContent(".", []string{tdir})
	if err != nil {
		t.Error(err.Error())
	}
	// Check if path exists
	for _, cpath := range []string{"content/chap2/_index.md", "static/images/wikipedia.png"} {
		if _, err := os.Stat(path.Join(tdir, cpath)); os.IsNotExist(err) {
			t.Errorf("Path '%s' should exist after creating Hugo in '%s'", cpath, tdir)
			return
		}
	}
	err = os.RemoveAll(tdir)
	if err != nil {
		t.Error(err.Error())
	}
}

func Test_ReadMeta_Chap1(t *testing.T) {
	c := CreateBase(3)
	c.ReadMetaData("../../misc/test/chap1/_index.md")
	c.PrintMeta()
	fails := []string{}
	fails = compVal(c.Meta["title"], "Chapter1", fails)
	fails = compVal(c.Meta["weight"], 10, fails)
	fails = compVal(c.Meta["pre"], "<b>I. </b>", fails)
	fails = compVal(c.Meta["chapter"], true, fails)
	fails = compVal(c.Meta["include_toc"], true, fails)
	if len(fails) > 0 {
		t.Error(strings.Join(fails, "\n"))
	}
}

func Test_ToMetaLines(t *testing.T) {
	c := CreateBase(3)
	c.ReadMetaData("../../misc/test/chap1/_index.md")
	exp := []string{
		"chapter: true",
		`include_toc: true`,
		`pre: "<b>I. </b>"`,
		`title: "Chapter1"`,
		"weight: 10",
	}
	got := c.ToMetaLines()
	sort.Strings(got)
	if len(exp) != len(got) {
		fmt.Printf("\nGOT:\n%s", strings.Join(got, "\n"))
		t.Errorf("Length of exp (%d) and got (%d) differs", len(exp), len(got))
	}
	if !reflect.DeepEqual(exp, got) {
		fmt.Printf("\n#####GOT:\n%s", strings.Join(got, "\n"))
		fmt.Printf("\n#####Exp:\n%s", strings.Join(exp, "\n"))
		t.Error("Got !- exp")
	}
}

func Test_ReplaceHeader_Simple(t *testing.T) {
	tdir, _ := ioutil.TempDir(os.TempDir(), "base-")

	b := Base{
		Source: "../../misc/test/chap1",
		Path:   "chap1",
	}
	fmt.Printf("CopyContent(., %s)\n  > Path: %s\n  > Source: %s\n", tdir, b.Path, b.Source)
	err := b.CopyContent(".", []string{tdir})
	if err != nil {
		t.Error(err.Error())
	}
	err = b.ReadMetaData(path.Join(tdir, "content/chap1/_index.md"))
	if err != nil {
		t.Error(err.Error())
	}
	b.Meta["title"] = "Chapter2"
	b.Meta["pre"] = "<b>II. </b>"
	b.Meta["weight"] = 20
	// Update Markdown file
	b.ReplaceHeader(path.Join(tdir, "content/chap1/_index.md"))
	err = b.ReadMetaData(path.Join(tdir, "content/chap1/_index.md"))
	if err != nil {
		t.Error(err.Error())
	}
	got := b.ToMetaLines()
	exp := []string{
		"chapter: true",
		`include_toc: true`,
		`pre: "<b>II. </b>"`,
		`title: "Chapter2"`,
		"weight: 20",
	}
	sort.Strings(got)
	if !reflect.DeepEqual(exp, got) {
		fmt.Printf("\n#####GOT:\n%s", strings.Join(got, "\n"))
		fmt.Printf("\n#####Exp:\n%s", strings.Join(exp, "\n"))
		t.Error("MetaData after replace header is not what we expect")
	}
	err = os.RemoveAll(tdir)
	if err != nil {
		log.Fatal(err)
	}
}
