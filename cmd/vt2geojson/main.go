package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/maptile"
)

func main() {
	flag.Parse()

	if len(flags.mvtSource) == 0 {
		log.Fatalf("Please specify the mvt file or URI by '-mvt'")
	}

	content, err := loadMVT(flags.mvtSource)
	if err != nil {
		log.Fatal(err)
	}

	layers, err := unmarshalMVT(content, flags.gzipped)
	if err != nil {
		log.Fatal(err)
	}

	if flags.summary {
		printLayersSummary(layers)
		return
	}

	// project all the geometries in all the layers backed to WGS84 from the extent and mercator projection.
	tile := maptile.New(uint32(flags.x), uint32(flags.y), maptile.Zoom(flags.z))
	layers.ProjectToWGS84(tile)

	// convert to geojson FeatureCollection
	featureCollections := layers.ToFeatureCollections()
	newFeatureCollection := geojson.NewFeatureCollection()
	if len(flags.layer) > 0 { // only specified layer
		v, found := featureCollections[flags.layer]
		if found {
			newFeatureCollection.Features = append(newFeatureCollection.Features, v.Features...)
		}
	} else { // all layers
		for _, v := range featureCollections {
			newFeatureCollection.Features = append(newFeatureCollection.Features, v.Features...)
		}
	}
	geojsonContent, err := newFeatureCollection.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", geojsonContent)
}
