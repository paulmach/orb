package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/paulmach/orb/encoding/mvt"
)

func main() {
	flag.Parse()

	if len(flags.mvtFile) > 0 {
		//TODO: parse mvt file to geojson

		content, err := ioutil.ReadFile(flags.mvtFile)
		if err != nil {
			log.Fatal(err)
		}

		layers, err := mvt.Unmarshal(content)
		if err != nil {
			log.Fatal(err)
		}

		if len(flags.layer) == 0 {
			for _, l := range layers {
				fmt.Printf("layer %s, version %d, extent %d, features %d\n", l.Name, l.Version, l.Extent, len(l.Features))
			}
			return
		}
		//fmt.Print(layers)
	}
}
