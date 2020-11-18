package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/qnib/go-hugo-remix/lib/types"
)

var (
	cfgFlag    = flag.String("config", "", "The workshop config file")
	targetFlag = flag.String("target", "", "where to assemble the workshop")
	debugLevel = flag.Int("d", 3, "Debug level (0-3) [default: 0]")
)

func main() {
	// Read Workshop config from file
	flag.Parse()
	err, w := types.CreateWorkshopFromFile(*cfgFlag)
	if err != nil {
		panic(err)
	}
	w.SetDebugLevel(*debugLevel)
	fmt.Println(w.String())
	_, genOut := w.GenerateHugo(*targetFlag)
	fmt.Printf(strings.Join(append(genOut, ""), "\n"))
	fmt.Println()
	// WriteHugoConfig
	w.WriteHugoConfig(*targetFlag)
}
