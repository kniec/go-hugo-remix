package main

import (
	"flag"
	"fmt"
	"strings"

	remix "github.com/qnib/go-hugo-remix/lib"
)

var (
	cfgFlag    = flag.String("config", "", "The workshop config file")
	targetFlag = flag.String("target", "", "where to assemble the workshop")
)

func main() {
	// Read Workshop config from file
	flag.Parse()
	err, w := remix.CreateWorkshopFromFile(*cfgFlag)
	if err != nil {
		panic(err)
	}
	fmt.Println(w.String())
	_, genOut := w.GenerateHugo(*targetFlag)
	fmt.Printf(strings.Join(append(genOut, ""), "\n"))
	fmt.Println()
	// WriteHugoConfig
	//hc := remix.HugoConfig{}

}
