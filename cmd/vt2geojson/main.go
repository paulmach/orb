package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/maptile"
)

// tryParseTileXYZ try to get x,y,z by parse 'z/x/y.mvt' or 'z/x/y.vector.pbf' from URL
func tryParseTileXYZ(mvtSource string) (maptile.Tile, error) {
	re := regexp.MustCompile(`([0-9]{1,2})/([0-9]+)/([0-9]+)(?:\.mvt|\.vector\.pbf)`)
	sub := re.FindStringSubmatch(mvtSource)
	if len(sub) < 4 { // match fail
		return maptile.Tile{}, fmt.Errorf("can not parse z,x,y from mvt %s", mvtSource)
	}

	var z, x, y uint64
	var err error
	if z, err = strconv.ParseUint(sub[1], 10, 32); err != nil {
		return maptile.Tile{}, fmt.Errorf("parse %s to z failed, err: %v", sub[1], err)
	}
	if x, err = strconv.ParseUint(sub[2], 10, 32); err != nil {
		return maptile.Tile{}, fmt.Errorf("parse %s to x failed, err: %v", sub[2], err)
	}
	if y, err = strconv.ParseUint(sub[3], 10, 32); err != nil {
		return maptile.Tile{}, fmt.Errorf("parse %s to y failed, err: %v", sub[3], err)
	}

	tile := maptile.New(uint32(x), uint32(y), maptile.Zoom(z))
	if !tile.Valid() {
		return tile, fmt.Errorf("parsed tile %v is not valid", tile)
	}

	return tile, nil
}

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
	var tile maptile.Tile
	if flags.x == 0 && flags.y == 0 && flags.z == 0 { // if all x,y,z are NOT set, try to parse from URL
		tile, err = tryParseTileXYZ(flags.mvtSource)
		if err != nil {
			log.Fatal(err)
		}
	} else { // otherwise use the flags set directly
		tile = maptile.New(uint32(flags.x), uint32(flags.y), maptile.Zoom(flags.z))
		if !tile.Valid() {
			log.Fatalf("Invalid flags x,y,z: %d,%d,%d", flags.x, flags.y, flags.z)
		}
	}
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
