package types

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/qnib/go-hugo-remix/lib/helper"
)

// Workshop provides the meta-data for a workshop and chapters
type Workshop struct {
	Title        string    `yaml:"title"`
	Author       string    `yaml:"author"`
	Description  string    `yaml:"description"`
	BaseURL      string    `yaml:"baseurl"`
	Flavours     []string  `yaml:"flavours"`
	BaseDir      string    // Path from which all relative path are originating (PWD of YAML file most likely)
	HugoBase     string    `yaml:"base"`          // Switch to embed hugo files later
	BaseWorkshop string    `yaml:"base-workshop"` // YAML file to extend workshop with
	Source       string    `yaml:"source"`        // source is the content of the base-url
	DstDir       string    // DstDir is used when copying files to store the destination
	Chaps        []Chapter `yaml:"chaps"`
	dLevel       int
}

// Parse takes a byte array and parses it
func (w *Workshop) Parse(yData []byte) {
	err := yaml.Unmarshal(yData, &w)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
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

func (w *Workshop) SetDebugLevel(l int) {
	w.dLevel = l
}
func (w *Workshop) GetDebugLevel() (l int) {
	return w.dLevel
}

// ExtendFromWorkshop takes w2 and extends w with it's chapters (authors)
// -> if a path (chapter || chapter/subchap) already exist in w, it WILL NOT be overwritten
func (w *Workshop) ExtendFromWorkshop(w2 Workshop) (err error) {
	oldPaths := helper.NewStrSet()
	for _, oc := range w.Chaps {
		oldPaths.Add(oc.Path)
		for _, ocs := range oc.Subchaps {
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

func (self *Workshop) Compare(other Workshop) (err error, fails []string) {
	fails = compVal(self.Title, other.Title, fails)
	fails = compVal(self.Author, other.Author, fails)
	fails = compVal(self.Description, other.Description, fails)
	fails = compVal(self.BaseURL, other.BaseURL, fails)
	fails = compVal(strings.Join(self.Flavours, ","), strings.Join(other.Flavours, ","), fails)
	fails = compVal(self.BaseDir, other.BaseDir, fails)
	fails = compVal(self.HugoBase, other.HugoBase, fails)
	fails = compVal(self.BaseWorkshop, other.BaseWorkshop, fails)
	fails = compVal(self.Source, other.Source, fails)
	fails = compVal(self.DstDir, other.DstDir, fails)
	if len(fails) > 0 {
		err = fmt.Errorf(strings.Join(fails, "\n"))
	}
	return
}
func (self *Workshop) SprintTitle() string {
	return fmt.Sprintf("Workshop.Title: %s", self.Title)
}
func (self *Workshop) SprintSource() string {
	return fmt.Sprintf("Workshop.Source: %s", self.Source)
}
func (self *Workshop) SprintAuthor() string {
	return fmt.Sprintf("Workshop.Author: %s", self.Author)
}
func (self *Workshop) SprintFlavours() string {
	return fmt.Sprintf("Workshop.Flavours: %s", strings.Join(self.Flavours, ","))
}

func (self *Workshop) String() (res []string) {
	res = append(res, self.SprintTitle())
	res = append(res, self.SprintSource())
	res = append(res, self.SprintAuthor())
	res = append(res, self.SprintFlavours())
	for _, chap := range self.Chaps {
		res = append(res, chap.String()...)
	}
	return
}

// WriteHugoConfig generates HugoConfig and puts a config file into the workshop dir
func (w *Workshop) WriteHugoConfig(tpath string) (err error) {
	err, hc := CreateHugoConfigFromWorkshop(*w)
	hc.WriteConfig(path.Join(tpath, "config.toml"))
	return
}

func (w *Workshop) CopyHugoBase(base, target string) (err error) {
	msg := fmt.Sprintf("cp -r %s/%s/* %s/", base, w.HugoBase, target)
	log.Printf(msg)
	err = CopyDir(path.Join(base, w.HugoBase), target)
	if err != nil {
		log.Printf(err.Error())

	}
	return
}

// GenerateHugo iterates over Chapter and Subchapters and copies the base, chapters and subchapters
// into a target directory
func (w *Workshop) GenerateHugo(target string) (err error, res []string) {
	log.Printf("GenerateHugo: %s", target)
	err = os.MkdirAll(target, 0755)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	err = w.CopyHugoBase(w.BaseDir, target)
	if err != nil {
		log.Print(err.Error())
		return
	}
	err = w.CopyWorkshopContent(target)
	if err != nil {
		log.Print(err.Error())
		return
	}
	for _, chap := range w.Chaps {
		chap.SetDebugLevel(w.GetDebugLevel())
		tPath := []string{target, "content", chap.Path}
		err = chap.CopyContent(w.BaseDir, tPath)
		if err != nil {
			log.Print(err.Error())
			return
		}
		for _, sub := range chap.Subchaps {
			sub.SetDebugLevel(w.GetDebugLevel())
			targetSub := append(tPath, sub.Path)
			err = sub.CopyContent(w.BaseDir, targetSub)
			if err != nil {
				log.Print(err.Error())
				return
			}
			for _, subsub := range sub.Subsubs {
				subsub.SetDebugLevel(w.GetDebugLevel())
				targetSubsub := append(targetSub, subsub.Path)
				err := subsub.CopyContent(w.BaseDir, targetSubsub)
				if err != nil {
					log.Print(err.Error())
				}
			}
		}
	}
	return
}

func (w *Workshop) CopyWorkshopContent(target string) (err error) {
	b := Base{
		Source: w.Source,
		Path:   ".",
	}
	return b.CopyContent(w.BaseDir, []string{target})
}
