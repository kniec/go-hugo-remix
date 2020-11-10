package content

import (
	"fmt"
	"strings"
)

/****
Generate standard content items like 'authors.md'
****/

const authorTempl = `---
title: Credits
disableToc: true
---

This workshop has been developed by:

#### Main Authors

{{ with .Main }}{{ range . }}
- **{{ .Name }}** // {{ .Contribution }}
{{ end }}{{ end }} 

#### Reviewer

{{ with .Review }}{{ range . }}
- **{{ .Name }}** // {{ .Contribution }}
{{ end }}{{ end }} 

#### Assists by

{{ with .Assist }}{{ range . }}
- **{{ .Name }}** // {{ .Contribution }}
{{ end }}{{ end }} 

`

var availSec = []string{"main", "review", "assist", "feedback"}

type Author struct {
	Order        int
	Section      string
	Name         string
	Contribution string
}

func CreateAuthor(id int, sec, name, contri string) (err error, a Author) {
	invalidSec := true
	for _, v := range availSec {
		if v == sec {
			invalidSec = false
		}
	}
	if invalidSec {
		return fmt.Errorf("'%s' not a valid section (%s)", sec, strings.Join(availSec, ",")), a
	}
	return nil, Author{
		Order:        id,
		Section:      sec,
		Name:         name,
		Contribution: contri,
	}
}

/*******
Authors
***/

type Authors struct {
	Main   []Author
	Review []Author
	Assist []Author
}

func CreateAuthors(as ...Author) (err error, res Authors) {
	for _, a := range as {
		switch a.Section {
		case "main":
			res.Main = append(res.Main, a)
		case "review":
			res.Review = append(res.Review, a)
		case "assist":
			res.Assist = append(res.Assist, a)
		}
	}
	return
}
