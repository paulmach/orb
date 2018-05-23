package geojson_test

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/quadtree"
)

func ExampleFeature_Point() {
	f := geojson.NewFeature(orb.Point{1, 1})
	f.Properties["key"] = "value"

	qt := quadtree.New(f.Geometry.Bound().Pad(1))
	qt.Add(f) // add the feature to a quadtree

	// type assert the feature back into a Feature from
	// the orb.Pointer interface.
	feature := qt.Find(orb.Point{0, 0}).(*geojson.Feature)
	fmt.Printf("key=%s", feature.Properties["key"])

	// Output:
	// key=value
}

func ExampleFeatureCollection_foreignMembers() {
	rawJSON := []byte(`
  { "type": "FeatureCollection",
	"features": [
	  { "type": "Feature",
		"geometry": {"type": "Point", "coordinates": [102.0, 0.5]},
		"properties": {"prop0": "value0"}
	  }
	],
	"title": "Title as Foreign Member"
  }`)

	type MyFeatureCollection struct {
		geojson.FeatureCollection
		Title string `json:"title"`
	}

	fc := &MyFeatureCollection{}
	json.Unmarshal(rawJSON, &fc)

	fmt.Println(fc.FeatureCollection.Features[0].Geometry)
	fmt.Println(fc.Features[0].Geometry)
	fmt.Println(fc.Title)
	// Output:
	// [102 0.5]
	// [102 0.5]
	// Title as Foreign Member
}

func ExampleUnmarshalFeatureCollection() {
	rawJSON := []byte(`
	  { "type": "FeatureCollection",
	    "features": [
	      { "type": "Feature",
	        "geometry": {"type": "Point", "coordinates": [102.0, 0.5]},
	        "properties": {"prop0": "value0"}
	      }
	    ]
	  }`)

	fc, _ := geojson.UnmarshalFeatureCollection(rawJSON)

	// Geometry will be unmarshalled into the correct geo.Geometry type.
	point := fc.Features[0].Geometry.(orb.Point)
	fmt.Println(point)

	// Output:
	// [102 0.5]
}

func Example_unmarshal() {
	rawJSON := []byte(`
	  { "type": "FeatureCollection",
	    "features": [
	      { "type": "Feature",
	        "geometry": {"type": "Point", "coordinates": [102.0, 0.5]},
	        "properties": {"prop0": "value0"}
	      }
	    ]
	  }`)

	fc := geojson.NewFeatureCollection()
	err := json.Unmarshal(rawJSON, &fc)
	if err != nil {
		log.Fatalf("invalid json: %v", err)
	}

	// Geometry will be unmarshalled into the correct geo.Geometry type.
	point := fc.Features[0].Geometry.(orb.Point)
	fmt.Println(point)

	// Output:
	// [102 0.5]
}

func ExampleFeatureCollection_MarshalJSON() {
	fc := geojson.NewFeatureCollection()
	fc.Append(geojson.NewFeature(orb.Point{1, 2}))

	data, err := fc.MarshalJSON()
	if err != nil {
		log.Fatalf("marshal error: %v", err)
	}

	// standard lib encoding/json package will also work
	data, err = json.MarshalIndent(fc, "", " ")
	if err != nil {
		log.Fatalf("marshal error: %v", err)
	}

	fmt.Println(string(data))

	// Output:
	// {
	//  "type": "FeatureCollection",
	//  "features": [
	//   {
	//    "type": "Feature",
	//    "geometry": {
	//     "type": "Point",
	//     "coordinates": [
	//      1,
	//      2
	//     ]
	//    },
	//    "properties": null
	//   }
	//  ]
	// }
}
