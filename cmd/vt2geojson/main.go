package main

import (
	"flag"

	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	defer glog.Flush()

	if len(flags.mvtFile) > 0 {
		//TODO: parse mvt file to geojson
	}
}
