package types

import (
	"testing"
)

func TestWriteConfig(t *testing.T) {
	oc := OutConfig{
		Home: []string{"HTML", "RSS", "JSON"},
	}
	pc := ParaConfig{
		Author:           "Christian Kniep",
		Description:      "My first workshop to test out the go redux script",
		ShowVisitedLinks: false,
		EditURL:          "https://github.com",
	}
	eng := Language{
		Title:        "Scientist soley focused on the application/workload",
		Weight:       1,
		LanguageName: "Researcher",
	}
	rse := Language{
		Title:        "Research Software Engineer: Research Sysops hybrid",
		Weight:       2,
		LanguageName: "RSE",
	}
	sys := Language{
		Title:        "System Admin focused on the infrastructure setup",
		Weight:       3,
		LanguageName: "SysOps",
	}
	hc := HugoConfig{
		Title:    "My Workshop",
		BaseURL:  "https://example.org",
		LangCode: "en-us",
		Theme:    []string{"hugo-theme-learn"},
		Outputs:  oc,
		Params:   pc,
		Languages: map[string]Language{
			"eng": eng,
			"rse": rse,
			"sys": sys,
		},
		DefLangInSub: false,
		DefContLang:  "eng",
	}
	hc.WriteConfig("./test_write_config.toml")
}
