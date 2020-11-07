package redux

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type HugoConfig struct {
	Title    string `toml:"title"`
	LangCode string `toml:"languageCode"`
	BaseURL  string `toml:"baseURL"`
	Theme    string `toml:"theme"`
	Outputs  OutConfig
	Params   ParaConfig
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