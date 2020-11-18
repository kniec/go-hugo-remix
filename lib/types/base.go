package types

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

// Base holds the common variables
type Base struct {
	Source    string
	Path      string
	dLevel    int
	Flavour   string
	Meta      map[string]interface{}
	MetaLines []string
}

func CreateBase(dlevel int) Base {
	var m map[string]interface{}
	return Base{
		dLevel: dlevel,
		Meta:   m,
	}
}

func (self *Base) SetDebugLevel(l int) {
	self.dLevel = l
}

func (self *Base) Debug(level int, msg string) {
	if level <= self.dLevel {
		fmt.Printf("%-5s: %s\n", strings.Repeat(">", level), msg)
	}
}
func (self *Base) Error(level int, msg string) {
	if level <= self.dLevel {
		fmt.Printf("%-5s: %s\n", strings.Repeat("!", level), msg)
	}
}

func (self *Base) CopyContent(baseDir string, tPath []string) (err error) {
	// TODO: we are assuming that 'static' is present once 'content' is!
	self.Debug(2, fmt.Sprintf("CopyContent(%s, %s)\n  > Path: %s\n  > Source: %s", ".", baseDir, strings.Join(tPath, "/"), self.Source))
	srcContent := path.Join(baseDir, self.Source, "content")
	if _, er := os.Stat(srcContent); os.IsNotExist(er) {
		self.Debug(2, fmt.Sprintf("Source '%s' does not exists", srcContent))
		// If it does not exists, we expect to have flat Markdown files within the source
		// and implictly no static content (like screenshots)
		//targetContentPath := append(tPath, "content", self.Path)
		srcContentPath := path.Join(baseDir, self.Source)
		self.Debug(2, fmt.Sprintf("cp -r %s %s", srcContentPath, path.Join(tPath...)))
		err = CopyDir(srcContentPath, path.Join(tPath...))
		if err != nil {
			return
		}
	} else {
		// if it DOES exists, we expect to have a content and a static dir which needs to be copied seperately
		// while removing the dirs from the source
		self.Debug(2, fmt.Sprintf("Source '%s' does exists, copying content/static seperatly", srcContent))
		srcContentPath := fmt.Sprintf("%s/content", path.Join(baseDir, self.Source))
		self.Debug(2, fmt.Sprintf("cp -r %s %s", srcContentPath, path.Join(tPath...)))
		err = CopyDir(srcContentPath, path.Join(tPath...))
		if err != nil {
			return
		}
		srcStaticPath := path.Join(baseDir, self.Source, "static")

		self.Debug(2, fmt.Sprintf("cp -r %s %s", srcStaticPath, "static"))
		err = CopyDir(srcStaticPath, "static")
		if err != nil {
			return
		}
	}
	return
}

func (self *Base) ReadMetaData(fpath string) (err error) {
	self.Debug(1, fmt.Sprintf("Reading: %s", fpath))
	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)
	var buf bytes.Buffer
	context := parser.NewContext()
	source, err := ioutil.ReadFile(fpath)
	if err != nil {
		self.Debug(0, fmt.Sprintf("Error reading file: %s", err.Error()))
		return err
	}
	if err := markdown.Convert([]byte(source), &buf, parser.WithContext(context)); err != nil {
		self.Debug(0, fmt.Sprintf("Error convert to markdown: %s", err.Error()))
		return err
	}
	self.Meta = meta.Get(context)
	return
}

func (self *Base) PrintMeta() {
	fmt.Println()
	for k, v := range self.Meta {
		fmt.Printf("%-20s: %v\n", k, v)
	}
}

func (self *Base) ToMetaLines() (res []string) {
	res = append(res, fmt.Sprintf(`title: "%s"`, self.Meta["title"]))
	res = append(res, fmt.Sprintf("weight: %d", self.Meta["weight"]))
	res = append(res, fmt.Sprintf(`chapter: %t`, self.Meta["chapter"]))
	res = append(res, fmt.Sprintf(`pre: "%s"`, self.Meta["pre"]))
	res = append(res, fmt.Sprintf(`include_toc: %t`, self.Meta["include_toc"]))
	return
}

// ReplaceHeader read a markdown file, update the meta-data info
// and replaces the markdown header within the file
func (self *Base) ReplaceHeader(fpath string) (err error) {
	// c.Meta is what is read from the file
	source, err := ioutil.ReadFile(fpath)
	inputLines := strings.Split(string(source), "\n")
	err, outputLines := updateMetaLines(inputLines, self.MetaLines)
	if err != nil {
		log.Fatalln(err)
	}
	output := strings.Join(outputLines, "\n")
	err = ioutil.WriteFile(fpath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
	return
}

func (self *Base) WalkContentDir(cDir string, mLines []string) (err error) {
	self.Debug(2, fmt.Sprintf("Start walking down '%s'", cDir))
	self.MetaLines = mLines
	err = filepath.Walk(cDir, self.WalkFunc)
	return
}

func (self *Base) WalkFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
		return err
	}
	if !info.IsDir() && info.Name() == "_index.md" {
		self.Debug(3, fmt.Sprintf("Found %s", path))
		self.ReplaceHeader(path)
	}
	return nil
}
