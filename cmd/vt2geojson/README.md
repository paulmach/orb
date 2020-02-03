# vt2geojson
Dump vector tiles to GeoJSON.     
Similar with [mapbox/vt2geojson](https://github.com/mapbox/vt2geojson), but implemented by `Golang` and based on [paulmach/orb](github.com/paulmach/orb/) library.     

## Usage

```bash
$ cd cmd/vt2geojson
$ go build 
```

- CLI helper
```bash
./vt2geojson -h
Usage of ./vt2geojson:
  -layer string
    	Include only the specified layer of the Mapbox Vector Tile.
  -mvt string
    	Mapbox Vector Tile file path or URL, e.g. 'xxx.mvt' or 'xxx.vector.pbf' or 'https://api.mapbox.com/v4/mapbox.mapbox-streets-v6/9/150/194.mvt?access_token=YOUR_MAPBOX_ACCESS_TOKEN'.
  -summary
    	Print layers summary only.
  -x uint
    	Tile x coordinate
  -y uint
    	Tile x coordinate
  -z uint
    	Tile zoom level.
```

- example
```bash
$ # print mvt layers summary 
$ ./vt2geojson -mvt https://api.mapbox.com/v4/mapbox.mapbox-streets-v6/9/150/194.mvt?access_token=YOUR_MAPBOX_ACCESS_TOKEN -summary
layer landuse, version 2, extent 4096, features 202
layer waterway, version 2, extent 4096, features 24
layer water, version 2, extent 4096, features 27
layer aeroway, version 2, extent 4096, features 5
layer landuse_overlay, version 2, extent 4096, features 38
layer road, version 2, extent 4096, features 787
layer admin, version 2, extent 4096, features 2
layer state_label, version 2, extent 4096, features 1
layer place_label, version 2, extent 4096, features 64
layer poi_label, version 2, extent 4096, features 17
layer road_label, version 2, extent 4096, features 11
$ 
$ # convert to geojson and dump to file
$ ./vt2geojson -x 150 -y 194 -z 9 -mvt https://api.mapbox.com/v4/mapbox.mapbox-streets-v6/9/150/194.mvt?access_token=YOUR_MAPBOX_ACCESS_TOKEN >geojson.json
$ ### if by mapbox/vt2geojson, the command is: 
$ ###   vt2geojson -x 150 -y 194 -z 9 https://api.mapbox.com/v4/mapbox.mapbox-streets-v6/9/150/194.mvt?access_token=YOUR_MAPBOX_ACCESS_TOKEN >geojson.json


```