package redux

/*****
Checks the MD files within a subchap or chapter and updates
the metadata as well as the language suffix
*****/

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

type m = map[string]interface{}

func GetDefaultChapIndex() m {
	return m{
		"chapter":     true,
		"pre":         "",
		"include_toc": true,
	}
}
func GetDefaultSubChapIndex() m {
	return m{
		"chapter":     false,
		"pre":         "",
		"include_toc": true,
	}
}

func GetDefaultSubChapNonIndex() m {
	return m{
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
	err, c.Meta = bm.UpdateDict(c.Meta)
	return
}

func (c *Checker) ReplaceHeader(md mType, fpath string) (err error) {
	out := []string{}
	// c.Meta is what is read from the file
	source, err := ioutil.ReadFile(fpath)
	lines := strings.Split(string(source), "\n")
	sawDashes := 0
	for _, line := range lines {
		if strings.Contains(line, "---") {
			out = append(out, line)
			sawDashes++
		}

		if sawDashes > 1 {
			out = append(out, line)
			continue
		}
		sp := strings.Split(line, ":")
		if len(sp) == 2 {

		}

	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(fpath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
	return
}
