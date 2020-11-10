package content

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCreateAuthor(t *testing.T) {
	_, got := CreateAuthor(0, "main", "Christian Kniep", "Framework, Baseline workshop")
	exp := Author{
		Order:        0,
		Section:      "main",
		Name:         "Christian Kniep",
		Contribution: "Framework, Baseline workshop",
	}
	if !reflect.DeepEqual(exp, got) {
		fmt.Printf("Got: %v", got)
		fmt.Printf("Exp: %v", exp)
		t.Errorf("Authors do not match")
	}
}

func TestCreateAuthors(t *testing.T) {
	_, a1 := CreateAuthor(0, "main", "Christian Kniep", "Framework, Baseline workshop")
	_, a2 := CreateAuthor(0, "review", "Christian Kniep", "Initial Reviewer")
	_, got := CreateAuthors(a1, a2)
	exp := Authors{
		Main: []Author{{
			Order:        0,
			Section:      "main",
			Name:         "Christian Kniep",
			Contribution: "Framework, Baseline workshop",
		}},
		Review: []Author{{
			Order:        0,
			Section:      "review",
			Name:         "Christian Kniep",
			Contribution: "Initial Reviewer",
		}},
	}
	if !reflect.DeepEqual(exp, got) {
		fmt.Printf("Got: %v", got)
		fmt.Printf("Exp: %v", exp)
		t.Errorf("Authors do not match")
	}
}
