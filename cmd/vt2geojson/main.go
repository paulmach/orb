package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/paulmach/orb/encoding/mvt"
)

func readMVTContent(mvtSource string) ([]byte, error) {

	if strings.Index(mvtSource, "http") >= 0 { // download from URL
		resp, err := http.Get(mvtSource)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return content, nil
	}

	// otherwise read from file
	content, err := ioutil.ReadFile(mvtSource)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func main() {
	flag.Parse()

	if len(flags.mvtSource) == 0 {
		log.Fatalf("Please specify the mvt source by '-mvt'")
	}

	//TODO: parse mvt contents to geojson

	content, err := readMVTContent(flags.mvtSource)
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
