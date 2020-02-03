package main

import "flag"

var flags struct {
	mvtFile string
	layer   string
}

func init() {
	flag.StringVar(&flags.mvtFile, "mvt", "", "Mapbox Vector Tile file path, e.g. 'xxx.mvt' or 'xxx.vector.pbf'.")
	flag.StringVar(&flags.layer, "layer", "", "Include only the specified layer of the Mapbox Vector Tile.")
}
