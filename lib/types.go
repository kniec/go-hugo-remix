package redux

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/qnib/go-hugo-remix/lib/helper"
)

// GetTitle returns the title
func (w *Workshop) GetTitle() string {
	return w.Title
}

type BaseMeta struct {
	Title      string
	Weight     int
	Chapter    bool
	Pre        string
	IncludeTOC bool
}

/********
* BaseMeta Update MD files
 */
func (bm *BaseMeta) ChapIndexWalker(tpath string, info os.FileInfo, err error) error {
	log.Printf("ChapIndexWalker(tpath=%s) // BM: %v", tpath, bm)
	switch {
	case info.IsDir():
		log.Printf("SKIP: %s", tpath)
		return filepath.SkipDir
	default:
		match, _ := regexp.MatchString(`_index.*\.md$`, info.Name())
		if match {
			c := Checker{}
			log.Printf("c.ReadMeta(%s)", tpath)
			err := c.ReadMeta(tpath)
			if err != nil {
				return err
			}
			err = c.UpdateMeta(*bm)
			if err != nil {
				log.Println(err.Error())
				return err
			}
			err = c.ReplaceHeader(tpath)
			if err != nil {
				log.Println(err.Error())
				return err
			}
		}
	}
	return nil
}

func (bm *BaseMeta) FileIndexWalker(root string, files []os.FileInfo) error {
	for _, file := range files {
		fpath := path.Join(root, file.Name())
		match, _ := regexp.MatchString(`_index.*\.md$`, file.Name())
		if match {
			c := Checker{}
			log.Printf("c.ReadMeta(%s)", fpath)
			err := c.ReadMeta(fpath)
			if err != nil {
				return err
			}
			err = c.UpdateMeta(*bm)
			if err != nil {
				log.Println(err.Error())
				return err
			}
			err = c.ReplaceHeader(fpath)
			if err != nil {
				log.Println(err.Error())
				return err
			}
		}

	}
	return nil
}

// UpdateIndex takes a directory and finds `_index\*.md` files to update
// Skipping subdirs
func (bm *BaseMeta) UpdateIndex(tpath string) (err error) {
	log.Printf(">>> Walking %s", tpath)
	files, err := ioutil.ReadDir(tpath)
	bm.FileIndexWalker(tpath, files)
	//err = filepath.Walk(tpath, bm.ChapIndexWalker)
	return
}

type mType map[string]interface{}

func (mt mType) ToStrings() (res []string) {
	for k, v := range mt {
		switch v.(type) {
		case string:
			res = append(res, fmt.Sprintf(`%s: "%s"`, k, v))
		case int:
			res = append(res, fmt.Sprintf(`%s: %d`, k, v))
		case bool:
			res = append(res, fmt.Sprintf(`%s: %t`, k, v))
		}
	}
	return
}

func (bm *BaseMeta) UpdateDict(md mType) (err error, res mType) {
	res = mType{}
	for k, v := range md {
		res[k] = v
	}
	res["title"] = bm.Title
	res["weight"] = bm.Weight
	res["chapter"] = bm.Chapter
	res["pre"] = bm.Pre
	res["include_toc"] = bm.IncludeTOC
	return
}

/*****************
-> Subsub
Pages under the subchapter that are l3 section (have an _index.md)
But within the menue they are at the same all align underneath the subchapter
I. Chapter     (/content/chap1/index.md)
  1. Subchap   (/content/chap1/sub1/_index.md)
    a. Subsub  (/content/chap1/sub1/suba.md || /content/chap1/sub1/suba/_index.md)
To make rearrangement easier, we'll use the latter
*****************/
type Subsub struct {
	Title      string   `yaml:"title"`
	Path       string   `yaml:"path"`
	Source     string   `yaml:"source"`
	Prefix     []string `yaml:"prefix"`
	Weight     int      `yaml:"weight"`
	IncludeTOC bool     `yaml:"include_toc"`
	Enum       string   `yaml:"enum"`
}

// CreateSubchapter returns a subchapter
func CreateSubsub(t, p, s, e string, w int, pre []string) (res Subsub) {
	res.Title = t
	res.Path = p
	res.Source = s
	res.Enum = e
	res.Weight = w
	res.Prefix = pre
	return res
}

func (s *Subsub) ToBaseMeta() BaseMeta {
	return BaseMeta{
		Title:      s.Title,
		Weight:     s.Weight,
		Chapter:    false,
		Pre:        s.Enum,
		IncludeTOC: s.IncludeTOC,
	}
}
func (s *Subsub) GetTitle() string {
	return s.Title
}

func (s *Subsub) CompareSubchap(s2 Subsub) (err error, fails []string) {
	if s.Path != s2.Path {
		fails = append(fails, fmt.Sprintf("Path: '%s' != '%s'", s.Path, s2.Source))
	}
	if s.Source != s2.Source {
		fails = append(fails, fmt.Sprintf("Source: '%s' != '%s'", s.Source, s2.Source))
	}
	/*
		// TODO: Might need to compare booth string slices
		if s.Prefix != s2.Prefix {
			fails = append(fails, fmt.Sprintf("Prefix: '%s' != '%s'", s.Prefix, s2.Prefix))
		}
	*/
	if s.Title != s2.Title {
		fails = append(fails, fmt.Sprintf("Title: '%s' != '%s'", s.Title, s2.Title))
	}
	if s.Weight != s2.Weight {
		fails = append(fails, fmt.Sprintf("Weight: '%d' != '%d'", s.Weight, s2.Weight))
	}
	if s.Enum != s2.Enum {
		fails = append(fails, fmt.Sprintf("Enum: '%s' != '%s'", s.Enum, s2.Enum))
	}
	if len(fails) > 0 {
		return fmt.Errorf(strings.Join(fails, "\n")), fails
	}
	return nil, fails
}

/***************
-> Subchap

***************/
type Subchapter struct {
	Title      string   `yaml:"title"`
	Path       string   `yaml:"path"`
	Source     string   `yaml:"source"`
	Prefix     []string `yaml:"prefix"`
	Weight     int      `yaml:"weight"`
	IncludeTOC bool     `yaml:"include_toc"`
	Enum       string   `yaml:"enum"`
	Subsub     []Subsub `yaml:"subsub"`
}

func (s *Subchapter) ToBaseMeta() BaseMeta {
	return BaseMeta{
		Title:      s.Title,
		Weight:     s.Weight,
		Chapter:    false,
		Pre:        s.Enum,
		IncludeTOC: s.IncludeTOC,
	}
}

// CreateSubchapter returns a subchapter
func CreateSubchapter(t, p, s, e string, w int, subs []Subsub) (res Subchapter) {
	res.Title = t
	res.Path = p
	res.Source = s
	res.Enum = e
	res.Weight = w
	res.Subsub = subs
	return res
}

func (s *Subchapter) GetTitle() string {
	return s.Title
}

func (s *Subchapter) CompareSubchap(s2 Subchapter) (err error, fails []string) {
	if s.Path != s2.Path {
		fails = append(fails, fmt.Sprintf("Path: '%s' != '%s'", s.Path, s2.Source))
	}
	if s.Source != s2.Source {
		fails = append(fails, fmt.Sprintf("Source: '%s' != '%s'", s.Source, s2.Source))
	}
	/*
		// TODO: Might need to compare booth string slices
		if s.Prefix != s2.Prefix {
			fails = append(fails, fmt.Sprintf("Prefix: '%s' != '%s'", s.Prefix, s2.Prefix))
		}
	*/
	if s.Title != s2.Title {
		fails = append(fails, fmt.Sprintf("Title: '%s' != '%s'", s.Title, s2.Title))
	}
	if s.Weight != s2.Weight {
		fails = append(fails, fmt.Sprintf("Weight: '%d' != '%d'", s.Weight, s2.Weight))
	}
	if s.Enum != s2.Enum {
		fails = append(fails, fmt.Sprintf("Enum: '%s' != '%s'", s.Enum, s2.Enum))
	}
	if len(fails) > 0 {
		return fmt.Errorf(strings.Join(fails, "\n")), fails
	}
	return nil, fails
}

// Chapter references a chapter in the hugo workshop
type Chapter struct {
	Title   string   `yaml:"title"`
	Path    string   `yaml:"path"`
	Source  string   `yaml:"source"`
	Prefix  []string `yaml:"prefix"`
	Weight  int      `yaml:"weight"`
	Enum    string   `yaml:"enum"`
	Subchap []Subchapter
}

// ToBaseMeta converts a chapter to update the metadata of a file
func (c *Chapter) ToBaseMeta() BaseMeta {
	return BaseMeta{
		Title:      c.Title,
		Weight:     c.Weight,
		Chapter:    true,
		Pre:        c.Enum,
		IncludeTOC: true,
	}
}

func (c *Chapter) String() (res []string) {
	res = append(res, fmt.Sprintf("Title: %s", c.Title))
	res = append(res, fmt.Sprintf("Path: %s", c.Path))
	res = append(res, fmt.Sprintf("Source: %s", c.Source))
	res = append(res, fmt.Sprintf("Enum: %s", c.Enum))
	res = append(res, fmt.Sprintf("Weight: %d", c.Weight))
	return
}

// CreateChapter build a chapter
func CreateChapter(t, p, s, e string, w int, pre []string, sub []Subchapter) Chapter {
	res := Chapter{}
	res.Title = t
	res.Path = p
	res.Source = s
	res.Enum = e
	res.Weight = w
	res.Prefix = pre
	res.Subchap = sub
	return res
}

func (c *Chapter) CompareChap(c2 Chapter) (err error, fails []string) {
	if c.Path != c2.Path {
		fails = append(fails, fmt.Sprintf("Path: '%s' != '%s'", c.Path, c2.Source))
	}
	if c.Source != c2.Source {
		fails = append(fails, fmt.Sprintf("Source: '%s' != '%s'", c.Source, c2.Source))
	}
	/*
		// TODO: Might need to compare booth string slices
		if s.Prefix != s2.Prefix {
			fails = append(fails, fmt.Sprintf("Prefix: '%s' != '%s'", s.Prefix, s2.Prefix))
		}
	*/
	if c.Title != c2.Title {
		fails = append(fails, fmt.Sprintf("Title: '%s' != '%s'", c.Title, c2.Title))
	}
	if c.Weight != c2.Weight {
		fails = append(fails, fmt.Sprintf("Weight: '%d' != '%d'", c.Weight, c2.Weight))
	}
	if c.Enum != c2.Enum {
		fails = append(fails, fmt.Sprintf("Enum: '%s' != '%s'", c.Enum, c2.Enum))
	}
	for i, sub := range c.Subchap {
		_, subFails := sub.CompareSubchap(c2.Subchap[i])
		fails = append(fails, subFails...)
	}
	if len(fails) > 0 {
		return fmt.Errorf(strings.Join(fails, "\n")), fails
	}
	return nil, fails
}

// Workshop references a workshop with all its metadata
type Workshop struct {
	Title       string   `yaml:"title"`
	Author      string   `yaml:"author"`
	Description string   `yaml:"description"`
	DstDir      string   // DstDir is used when copying files to store the destination
	BaseURL     string   `yaml:"baseurl"`
	Flavours    []string `yaml:"flavours`
	BaseDir     string
	// HugoBase path to copy from
	// Switch to embed hugo files later
	// -> https://stackoverflow.com/questions/17796043/how-to-embed-files-into-golang-binaries
	HugoBase string `yaml:"base"`
	// YAML file to extend workshop with
	BaseWorkshop string `yaml:"base-workshop"`
	// source is the content of the base-url
	Source string `yaml:"source"`
	Chaps  []Chapter
}

// CreateWorkshop creates a workshop
func CreateWorkshop(t, desc, burl, src string, c []Chapter) Workshop {
	return Workshop{
		Title:       t,
		Description: desc,
		BaseURL:     burl,
		HugoBase:    "../misc/hugo/",
		Source:      src,
		Chaps:       c,
	}
}

func (w *Workshop) CompareWorkshops(w2 Workshop) (err error, fails []string) {
	if w.Title != w2.Title {
		fails = append(fails, fmt.Sprintf("'%s' :w.Title != w2.Title: %s", w.Title, w2.Title))
	}
	if w.HugoBase != w2.HugoBase {
		fails = append(fails, fmt.Sprintf("'%s' :w.HugoBase != w2.HugoBase: %s", w.HugoBase, w2.HugoBase))
	}
	if w.HugoBase != w2.HugoBase {
		fails = append(fails, fmt.Sprintf("'%s' :w.HugoBase != w2.HugoBase: %s", w.HugoBase, w2.HugoBase))
	}
	for i, chap := range w.Chaps {
		_, chapFails := chap.CompareChap(w2.Chaps[i])
		fails = append(fails, chapFails...)
	}

	if len(fails) > 0 {
		return fmt.Errorf("Fail!"), fails
	}
	return
}

// ExtendFromWorkshop takes w2 and extends w with it's chapters (authors)
// -> if a path (chapter || chapter/subchap) already exist in w, it WILL NOT be overwritten
func (w *Workshop) ExtendFromWorkshop(w2 Workshop) (err error) {
	oldPaths := helper.NewStrSet()
	for _, oc := range w.Chaps {
		oldPaths.Add(oc.Path)
		for _, ocs := range oc.Subchap {
			oldPaths.Add(path.Join(oc.Path, ocs.Path))
		}
	}
	for _, chap := range w2.Chaps {
		//fmt.Printf("> Check w2.Chap.Path '%s'", chap.Path)
		if !oldPaths.Contains(chap.Path) {
			//fmt.Printf(" - NOT FOUND in old chaps\n")
			w.Chaps = append(w.Chaps, chap)
		} else {
			//fmt.Printf(" - FOUND in old chaps\n")
			log.Printf("!! Chap.Path '%s' already exists; won't extend the subchaps", chap.Path)
		}
	}
	return
}

// Parse takes a byte array and parses it
func (w *Workshop) Parse(yData []byte) {
	err := yaml.Unmarshal(yData, &w)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

func (w *Workshop) String() (res []string) {
	res = append(res, fmt.Sprintf("Title: %s", w.Title))
	return res
}

func readFile(path string) ([]byte, error) {
	yamlFile, err := ioutil.ReadFile(path)
	return yamlFile, err
}

// CreateWorkshopFromFile parses a file and returns a workshop
func CreateWorkshopFromFile(fpath string) (err error, w Workshop) {
	log.Printf("Reading file: %s", fpath)
	yData, err := readFile(fpath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	w.Parse(yData)
	w.BaseDir = filepath.Dir(fpath)
	if w.BaseWorkshop != "" {
		e, wExt := CreateWorkshopFromFile(path.Join(w.BaseDir, w.BaseWorkshop))
		if e != nil {
			return e, w
		}
		w.ExtendFromWorkshop(wExt)
	}
	return
}

// WriteHugoConfig generates HugoConfig and puts a config file into the workshop dir
func (w *Workshop) WriteHugoConfig(tpath string) (err error) {
	err, hc := CreateHugoConfigFromWorkshop(*w)
	hc.WriteConfig(path.Join(tpath, "config.toml"))
	return
}

//
///// CopyDirWalk copies files and directories
func (w *Workshop) CopyDirWalk(src, dst string) (err error) {
	w.DstDir = dst
	err = filepath.Walk(src, w.WalkCopy)
	return
}

///// WalkCopy takes multi-flavor (language) into account
// by extending md-files sufix (_index.md -> {_index.oyo.md, _index.ee.md})
// and 'files' folders (_index.files -> {_index.oyo.files, _index.ee.files})
//// BUT when a folder already has a language sufix, it should not be duplicated
// -> this means that the author explicitly want it to be only for the language in question
func (w *Workshop) WalkCopy(path string, info os.FileInfo, err error) error {
	return nil
}

// GenerateHugo iterates over Chapter and Subchapters and copies the base, chapters and subchapters
// into a target directory
func (w *Workshop) GenerateHugo(t string) (err error, res []string) {
	log.Printf("GenerateHugo: %s", t)
	err = os.MkdirAll(t, 0755)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	msg := fmt.Sprintf("cp -r %s/%s/* %s/", w.BaseDir, w.HugoBase, t)
	log.Printf(msg)
	res = append(res, msg)
	err = CopyDir(path.Join(w.BaseDir, w.HugoBase), t)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	// Workshop Base FIles
	srcContent := path.Join(w.BaseDir, w.Source, "content")
	log.Printf("Check if '%s' exists", srcContent)
	if _, er := os.Stat(srcContent); os.IsNotExist(er) {
		log.Printf("%s does not exists, so we copy the flat MD files", srcContent)
		// If it does not exists, we expect to have flat Markdown files within the source
		// and implictly no static content (like screenshots)
		targetContentPath := path.Join(t, "content")
		srcContentPath := path.Join(w.BaseDir, w.Source)
		res = append(res, fmt.Sprintf("cp -r %s %s", srcContentPath, targetContentPath))
		err = CopyDir(srcContentPath, targetContentPath)
		if err != nil {
			log.Printf(">> %s -> %s", srcContentPath, targetContentPath)
			log.Printf(">> %s", err.Error())
			return
		}
	} else {
		// if it DOES exists, we expect to have a content and a static dir which needs to be copied seperately
		// while removing the dirs from the source
		srcContentPath := path.Join(w.BaseDir, w.Source, "content")
		targetContentPath := fmt.Sprintf("%s/content", t)
		res = append(res, fmt.Sprintf("cp -r %s %s", srcContentPath, targetContentPath))
		err = CopyDir(srcContentPath, targetContentPath)
		if err != nil {
			log.Printf(">> %s -> %s", srcContentPath, targetContentPath)
			log.Printf(">> %s", err.Error())
			return
		}
		srcStaticPath := path.Join(w.BaseDir, w.Source, "static")
		targetStaticPath := path.Join(t, "static")
		res = append(res, fmt.Sprintf("cp -r %s %s", srcStaticPath, targetStaticPath))
		err = CopyDir(srcStaticPath, targetStaticPath)
		if err != nil {
			log.Printf(">> %s -> %s", srcStaticPath, targetStaticPath)
			log.Printf(">> %s", err.Error())
			log.Println(err.Error())
			return
		}
	}
	for _, chap := range w.Chaps {
		// TODO: we are assuming that 'static' is present once 'content' is!
		srcContent := path.Join(w.BaseDir, chap.Source, "content")
		log.Printf("Check if '%s' exists", srcContent)
		if _, er := os.Stat(srcContent); os.IsNotExist(er) {
			log.Printf("%s does not exists, so we copy the flat MD files", srcContent)
			// If it does not exists, we expect to have flat Markdown files within the source
			// and implictly no static content (like screenshots)
			targetContentPath := fmt.Sprintf("%s/content/%s", t, chap.Path)
			srcContentPath := path.Join(w.BaseDir, chap.Source)
			res = append(res, fmt.Sprintf("cp -r %s %s", srcContentPath, targetContentPath))
			err = CopyDir(srcContentPath, targetContentPath)
			if err != nil {
				log.Println(err.Error())
				return
			}
		} else {
			// if it DOES exists, we expect to have a content and a static dir which needs to be copied seperately
			// while removing the dirs from the source
			srcContentPath := fmt.Sprintf("%s/content", path.Join(w.BaseDir, chap.Source))
			targetContentPath := fmt.Sprintf("%s/content/%s", t, chap.Path)
			res = append(res, fmt.Sprintf("cp -r %s %s", srcContentPath, targetContentPath))
			err = CopyDir(srcContentPath, targetContentPath)
			if err != nil {
				log.Println(err.Error())
				return
			}
			srcStaticPath := fmt.Sprintf("%s/static", path.Join(w.BaseDir, chap.Source))
			targetStaticPath := fmt.Sprintf("%s/static", t)
			res = append(res, fmt.Sprintf("cp -r %s %s", srcStaticPath, targetStaticPath))
			err = CopyDir(srcStaticPath, targetStaticPath)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
		bm := chap.ToBaseMeta()
		bm.UpdateIndex(path.Join(t, "content", w.BaseDir, chap.Path, ""))
		for _, sub := range chap.Subchap {
			if _, er := os.Stat(path.Join(sub.Source, "content")); os.IsNotExist(er) {
				targetContentPath := path.Join(t, "content", chap.Path, sub.Path)
				srcContentPath := path.Join(w.BaseDir, sub.Source)
				res = append(res, fmt.Sprintf("cp -r %s %s", srcContentPath, targetContentPath))
				err = CopyDir(srcContentPath, targetContentPath)
				if err != nil {
					log.Println(err.Error())
					return
				}
			} else {
				srcContentPath := path.Join(w.BaseDir, sub.Source, "content")
				targetContentPath := path.Join(t, "content", chap.Path, sub.Path)
				res = append(res, fmt.Sprintf("cp -r %s %s", srcContentPath, targetContentPath))
				err = CopyDir(srcContentPath, targetContentPath)
				if err != nil {
					log.Println(err.Error())
					return
				}
				srcStaticPath := fmt.Sprintf("%s/static", path.Join(w.BaseDir, sub.Source))
				targetStaticPath := fmt.Sprintf("%s/static", t)
				res = append(res, fmt.Sprintf("cp -r %s %s", srcStaticPath, targetStaticPath))
				err = CopyDir(srcStaticPath, targetStaticPath)
				if err != nil {
					log.Println(err.Error())
					return
				}
			}
			bm = sub.ToBaseMeta()
			bm.UpdateIndex(path.Join(t, "content", chap.Path, sub.Path, ""))
			for _, subsub := range sub.Subsub {
				srcContentPath := path.Join(w.BaseDir, subsub.Source)
				if _, er := os.Stat(path.Join(srcContentPath, "content")); os.IsNotExist(er) {
					targetContentPath := path.Join(t, "content", chap.Path, sub.Path, subsub.Path)
					msg := fmt.Sprintf("cp -r %s %s", srcContentPath, targetContentPath)
					res = append(res, msg)
					err = CopyDir(srcContentPath, targetContentPath)
					if err != nil {
						log.Println(err.Error())
						return
					}
				} else {
					srcContentPath := fmt.Sprintf("%s/content", path.Join(w.BaseDir, subsub.Source))
					targetContentPath := path.Join(t, "content", chap.Path, sub.Path, subsub.Path)
					res = append(res, fmt.Sprintf("cp -r %s %s", srcContentPath, targetContentPath))
					err = CopyDir(srcContentPath, targetContentPath)
					if err != nil {
						log.Println(err.Error())
						return
					}
					srcStaticPath := fmt.Sprintf("%s/static", path.Join(w.BaseDir, subsub.Source))
					targetStaticPath := path.Join(t, "static")
					res = append(res, fmt.Sprintf("cp -r %s %s", srcStaticPath, targetStaticPath))
					err = CopyDir(srcStaticPath, targetStaticPath)
					if err != nil {
						log.Println(err.Error())
						return
					}
				}
				bm = subsub.ToBaseMeta()
				upath := path.Join(t, "content", chap.Path, sub.Path, subsub.Path)
				log.Printf("UpdateIndex(%s)", upath)
				bm.UpdateIndex(upath)
			}
		}
	}
	return
}

/*********
Test Objects
***********/

// TestObj
var testC1sub1 = CreateSubchapter("Chap1Sub1", "sub1", "./test/sub1", "1. ", 10, []Subsub{})
var testC1sub2 = CreateSubchapter("Chap1Sub2", "sub2", "./test/sub2", "2. ", 20, []Subsub{})
var testC2sub1 = CreateSubchapter("Chap2Sub1", "sub1", "./test/sub3", "1. ", 10, []Subsub{})
var testChap1 = CreateChapter("Chapter1", "chap1", "./test/chap1", "I. ", 10, []string{},
	[]Subchapter{testC1sub1, testC1sub2})
var testChap2 = CreateChapter("Chapter2", "chap2", "./test/chap2", "II. ", 20, []string{},
	[]Subchapter{testC2sub1})

var testWorkshop = Workshop{
	Title:    "Workshop1",
	HugoBase: "./hugo",
	Chaps:    []Chapter{testChap1, testChap2},
}

// GetTestWorkshop returns the testwoerkshop
func GetTestWorkshop() Workshop {
	return testWorkshop
}
