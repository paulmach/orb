package geojson_test

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

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

func ExampleUnmarshal() {
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
