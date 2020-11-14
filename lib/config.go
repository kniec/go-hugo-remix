package redux

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

type HugoConfig struct {
	Title     string   `toml:"title"`
	LangCode  string   `toml:"languageCode"`
	BaseURL   string   `toml:"baseURL"`
	Theme     []string `toml:"theme"`
	Outputs   OutConfig
	Params    ParaConfig
	Languages map[string]Language `toml:"languages"`
}

type OutConfig struct {
	Home []string `toml:"home"`
}

type ParaConfig struct {
	Author           string `toml:"author"`
	Description      string `toml:"description"`
	ShowVisitedLinks bool   `toml:"showVisitedLinks"`
	EditURL          string `toml:"editURL"`
}

type Language struct {
	Title        string `toml:"title"`
	Weight       int    `toml:"weight"`
	LanguageName string `toml:"languageName"`
}

/*
[Languages]
[Languages.eng]
title = "Scientist soley focused on the application/workload"
weight = 1
languageName = "Engineer/Scientist"

[Languages.rse]
title = "Research Software Engineer: Research Sysops hybrid"
weight = 2
languageName = "RSE"

[Languages.sys]
title = "System Admin focused on the Infrastructure Setup"
weight = 3
languageName = "SysOps"
*/

var defFlavourMap = map[string]Language{
	"eng": Language{
		Title:        "Scientist soley focused on the application/workload",
		Weight:       1,
		LanguageName: "Researcher",
	},
	"rse": Language{
		Title:        "Research Software Engineer: Research Sysops hybrid",
		Weight:       2,
		LanguageName: "RSE",
	},
	"sys": Language{
		Title:        "System Admin focused on the infrastructure setup",
		Weight:       3,
		LanguageName: "SysOps",
	},
}

func CreateHugoConfigFromWorkshop(w Workshop) (err error, hc HugoConfig) {
	oc := OutConfig{
		Home: []string{"HTML", "RSS", "JSON"},
	}
	pc := ParaConfig{
		Author:           "Christian Kniep",
		Description:      w.Description,
		ShowVisitedLinks: true,
		EditURL:          "https://github.com",
	}
	hc = HugoConfig{
		Title:     w.Title,
		BaseURL:   "https://example.org",
		LangCode:  "en-us",
		Theme:     []string{"video", "learn"},
		Outputs:   oc,
		Params:    pc,
		Languages: map[string]Language{},
	}
	for _, flav := range w.Flavours {
		if val, ok := defFlavourMap[flav]; ok {
			hc.Languages[flav] = val
		} else {
			keys := []string{}
			for k, _ := range defFlavourMap {
				keys = append(keys, k)
			}
			err = fmt.Errorf("Flavor '%s' not found in defFlavourMap (%s)", flav, strings.Join(keys, ","))
			return
		}

	}
	return
}

func (hc *HugoConfig) WriteConfig(tFile string) (err error) {
	f, err := os.Create(tFile)
	if err != nil {
		// failed to create/open the file
		log.Fatal(err)
	}
	if err := toml.NewEncoder(f).Encode(hc); err != nil {
		// failed to encode
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		// failed to close the file
		log.Fatal(err)
	}
	return
}
