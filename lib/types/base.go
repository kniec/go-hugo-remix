package types

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

// Chapter is the highest level content
// Chapter -> Subchap -> Subsub
type Chapter struct {
	Base
	Subchap []Subchapter
}

// CreateChapter build a chapter
func CreateChapter(t, p, s, e string, w int, sub []Subchapter) Chapter {
	b := CreateBase(t, p, s, e, w)
	res := Chapter{
		Base:    b,
		Subchap: sub,
	}
	return res
}

// Subchapter is the mid-level content
// Chapter -> Subchap -> Subsub
type Subchapter struct {
	Base
	Subsub []Subsub
}

// CreateSubchap build a Subchapter
func CreateSubchapter(t, p, s, e string, w int, subsub []Subsub) Subchapter {
	b := CreateBase(t, p, s, e, w)
	res := Subchapter{
		Base:   b,
		Subsub: subsub,
	}
	return res
}

// Subsub is the lowest level content
// Chapter -> Subchap -> Subsub
type Subsub struct {
	Base
}

// CreateSubsub build a Subsub
func CreateSubsub(t, p, s, e string, w int) Subsub {
	b := CreateBase(t, p, s, e, w)
	res := Subsub{
		Base: b,
	}
	return res
}

// Workshop provides the meta-data for a workshop and chapters
type Workshop struct {
	Title        string   `yaml:"title"`
	Author       string   `yaml:"author"`
	Description  string   `yaml:"description"`
	BaseURL      string   `yaml:"baseurl"`
	Flavours     []string `yaml:"flavours"`
	BaseDir      string
	HugoBase     string `yaml:"base"`          // Switch to embed hugo files later
	BaseWorkshop string `yaml:"base-workshop"` // YAML file to extend workshop with
	Source       string `yaml:"source"`        // source is the content of the base-url
	DstDir       string // DstDir is used when copying files to store the destination
	Chaps        []Chapter
}
