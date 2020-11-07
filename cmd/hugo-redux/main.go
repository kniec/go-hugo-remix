package main

import (
	"flag"
	"fmt"
	"strings"

	redux "github.com/qnib/go-hugo-redux/lib"
)

var (
	cfgFlag    = flag.String("config", "", "The workshop config file")
	targetFlag = flag.String("target", "", "where to assemble the workshop")
)

func main() {
	// Read Workshop config from file
	flag.Parse()
	err, w := redux.CreateWorkshopFromFile(*cfgFlag)
	if err != nil {
		panic(err)
	}
	fmt.Println(w.String())
	fmt.Printf(strings.Join(append(w.GenerateHugo(*targetFlag), ""), "\n"))
	fmt.Println()
}
