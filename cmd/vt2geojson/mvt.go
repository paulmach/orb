package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

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

func printLayersSummary(layers mvt.Layers) {
	for _, l := range layers {
		fmt.Printf("layer %s, version %d, extent %d, features %d\n", l.Name, l.Version, l.Extent, len(l.Features))
	}
}
