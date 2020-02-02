package main

import "flag"

var flags struct {
	mvtFile string
}

func init() {
	flag.StringVar(&flags.mvtFile, "f", "", "Mapbox Vector Tile file path, e.g. 'xxx.mvt' or 'xxx.vector.pbf'.")
}
