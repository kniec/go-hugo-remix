package config

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
	hc := HugoConfig{
		Title:    "My Workshop",
		BaseURL:  "https://example.org",
		LangCode: "en-us",
		Theme:    []string{"hugo-theme-learn"},
		Outputs:  oc,
		Params:   pc,
	}
	hc.WriteConfig("./test_write_config.toml")
}
