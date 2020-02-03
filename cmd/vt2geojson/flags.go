package main

import "flag"

var flags struct {
	mvtSource string
	layer     string
}

func init() {
	flag.StringVar(&flags.mvtSource, "mvt", "", "Mapbox Vector Tile file path or URL, e.g. 'xxx.mvt' or 'xxx.vector.pbf' or 'https://api.mapbox.com/v4/mapbox.mapbox-streets-v6/9/150/194.mvt?access_token=YOUR_MAPBOX_ACCESS_TOKEN'.")
	flag.StringVar(&flags.layer, "layer", "", "Include only the specified layer of the Mapbox Vector Tile.")
}
