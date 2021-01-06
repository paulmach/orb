orb/geojson [![Godoc Reference](https://godoc.org/github.com/paulmach/orb/geojson?status.svg)](https://godoc.org/github.com/paulmach/orb/geojson)
===========

This package **encodes and decodes** [GeoJSON](http://geojson.org/) into Go structs
using the geometries in the [orb](https://github.com/paulmach/orb) package.
Supports both the [json.Marshaler](http://golang.org/pkg/encoding/json/#Marshaler) and
[json.Unmarshaler](http://golang.org/pkg/encoding/json/#Unmarshaler) interfaces.
The package also provides helper functions such as `UnmarshalFeatureCollection` and `UnmarshalFeature`.

## Examples

#### Unmarshalling  (JSON -> Go)

```go
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

// or

fc := geojson.NewFeatureCollection()
err := json.Unmarshal(rawJSON, &fc)

// Geometry will be unmarshalled into the correct geo.Geometry type.
point := fc.Features[0].Geometry.(orb.Point)
```

#### Marshalling (Go -> JSON)

```go
fc := geojson.NewFeatureCollection()
fc.Append(geojson.NewFeature(orb.Point{1, 2}))

rawJSON, _ := fc.MarshalJSON()

// or
blob, _ := json.Marshal(fc)
```

#### Foreign/extra members in a feature collection

```go
rawJSON := []byte(`
  { "type": "FeatureCollection",
    "generator": "myapp",
    "timestamp": "2020-06-15T01:02:03Z",
    "features": [
      { "type": "Feature",
        "geometry": {"type": "Point", "coordinates": [102.0, 0.5]},
        "properties": {"prop0": "value0"}
      }
    ]
  }`)

fc, _ := geojson.UnmarshalFeatureCollection(rawJSON)

fc.ExtraMembers["generator"] // == "myApp"
fc.ExtraMembers["timestamp"] // == "2020-06-15T01:02:03Z"

// Marshalling will include values in `ExtraMembers` in the
// base featureCollection object.
```

## Feature Properties

GeoJSON features can have properties of any type. This can cause issues in a statically typed
language such as Go. Included is a `Properties` type with some helper methods that will try to
force convert a property. An optional default, will be used if the property is missing or the wrong
type.

	f.Properties.MustBool(key string, def ...bool) bool
	f.Properties.MustFloat64(key string, def ...float64) float64
	f.Properties.MustInt(key string, def ...int) int
	f.Properties.MustString(key string, def ...string) string
