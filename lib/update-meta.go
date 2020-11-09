package redux

/*****
Checks the MD files within a subchap or chapter and updates
the metadata as well as the language suffix
*****/

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

type m = map[string]interface{}

func GetDefaultChapIndex() m {
	return m{
		"title":       "",
		"weight":      0,
		"chapter":     true,
		"pre":         "",
		"include_toc": true,
	}
}
func GetDefaultSubChapIndex() m {
	return m{
		"title":       "",
		"weight":      0,
		"chapter":     false,
		"pre":         "",
		"include_toc": true,
	}
}

func GetDefaultSubChapNonIndex() m {
	return m{
		"title":       "",
		"weight":      0,
		"chapter":     false,
		"pre":         "",
		"include_toc": false,
	}
}

// Checker holds config we might need
type Checker struct {
	IndexRegex *regexp.Regexp
	Meta       map[string]interface{}
}

// ProcessFile will be applied to each item found in WalkPath
func (c *Checker) ProcessFile(path string, fi os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if fi.IsDir() {
		// We don't care about directories.
		// A chapter will have subchapter dirs
		return nil
	}
	if c.IndexRegex.MatchString(path) {
		c.ReadMeta(path)
	}
	return nil
}

// ReadMeta Reads the file content and changes the meta-data
func (c *Checker) ReadMeta(fpath string) (err error) {
	log.Printf("Reading: %s", fpath)
	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)
	var buf bytes.Buffer
	context := parser.NewContext()
	source, err := ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}
	if err := markdown.Convert([]byte(source), &buf, parser.WithContext(context)); err != nil {
		return err
	}
	c.Meta = meta.Get(context)
	return
}

// UpdateMeta update the file metadata (c.Meta) with what comes from the workshop.yaml
func (c *Checker) UpdateMeta(bm BaseMeta) (err error) {
	// TODO: that's the place to follow up
	log.Printf("UpdateMeta.Before: %v+", c.Meta)
	err, c.Meta = bm.UpdateDict(c.Meta)
	log.Printf("UpdateMeta.After: %v+", c.Meta)
	return
}

func (c *Checker) ToMetaLines() (res []string) {
	res = append(res, fmt.Sprintf(`title: "%s"`, c.Meta["title"]))
	res = append(res, fmt.Sprintf("weight: %d", c.Meta["weight"]))
	res = append(res, fmt.Sprintf(`chapter: %t`, c.Meta["chapter"]))
	res = append(res, fmt.Sprintf(`pre: "%s"`, c.Meta["pre"]))
	res = append(res, fmt.Sprintf(`include_toc: %t`, c.Meta["include_toc"]))
	return
}
func (c *Checker) ReplaceHeader(fpath string) (err error) {
	// c.Meta is what is read from the file
	source, err := ioutil.ReadFile(fpath)
	inputLines := strings.Split(string(source), "\n")
	log.Printf("inputLines: %v", inputLines)
	metaLines := c.ToMetaLines()
	log.Printf("ToMetaLines: %v", metaLines)

	err, outputLines := updateMetaLines(inputLines, metaLines)
	log.Printf("outputLines: %v", outputLines)
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

func updateMetaLines(inputLines, metaLines []string) (err error, outputLines []string) {
	sawDashes := 0
	for _, line := range inputLines {
		if strings.Contains(line, "---") {
			if sawDashes == 0 {
				outputLines = append(outputLines, line)
				outputLines = append(outputLines, metaLines...)
			}
			sawDashes++
		}
		switch {
		case sawDashes == 1:
			continue
		case sawDashes > 1:
			outputLines = append(outputLines, line)
			continue
		}
	}
	return
}
