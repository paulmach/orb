package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/maptile"

	"github.com/paulmach/orb/encoding/mvt"
)

func loadMVT(mvtSource string) ([]byte, error) {

	if strings.Index(mvtSource, "http") >= 0 { // download from URL
		tr := &http.Transport{
			DisableCompression: true, // disable the silently decompression
		}
		client := &http.Client{Transport: tr}
		resp, err := client.Get(mvtSource)
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

func unmarshalMVT(mvtContent []byte, gzipped bool) (mvt.Layers, error) {
	if gzipped {
		return mvt.UnmarshalGzipped(mvtContent)
	}
	return mvt.Unmarshal(mvtContent)
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
		for _, l := range layers {
			fmt.Printf("layer %s, version %d, extent %d, features %d\n", l.Name, l.Version, l.Extent, len(l.Features))
		}
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
