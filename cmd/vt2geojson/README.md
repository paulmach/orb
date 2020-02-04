# vt2geojson
Dump vector tiles to GeoJSON.     
Similar with [mapbox/vt2geojson](https://github.com/mapbox/vt2geojson), but implemented by `Golang` and based on [paulmach/orb](github.com/paulmach/orb/) library.     

## Build
```bash
$ cd cmd/vt2geojson
$ go build 
```

## Usage


- CLI helper
```bash
./vt2geojson -h
Usage of ./vt2geojson:
  -gzipped
    	Whether uncompress the '.mvt' by gzip or not. '.mvt' comes from mapbox server is always gzipped, whatever with the 'Accept-Encoding: gzip' or not. (default true)
  -layer string
    	Include only the specified layer of the Mapbox Vector Tile.
  -mvt string
    	Mapbox Vector Tile file path or URL, e.g. 'xxx.mvt' or 'xxx.vector.pbf' or 'https://api.mapbox.com/v4/mapbox.mapbox-streets-v6/9/150/194.mvt?access_token=YOUR_MAPBOX_ACCESS_TOKEN'.
  -summary
    	Print layers summary only.
  -x uint
    	Tile x coordinate (normally inferred from the URL, e.g. 'z/x/y.mvt' or 'z/x/y.vector.pbf')
  -y uint
    	Tile x coordinate (normally inferred from the URL, e.g. 'z/x/y.mvt' or 'z/x/y.vector.pbf')
  -z uint
    	Tile zoom level (normally inferred from the URL, e.g. 'z/x/y.mvt' or 'z/x/y.vector.pbf')
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
$ ./vt2geojson -mvt https://api.mapbox.com/v4/mapbox.mapbox-streets-v6/9/150/194.mvt?access_token=YOUR_MAPBOX_ACCESS_TOKEN >geojson.json
$ ### if by mapbox/vt2geojson, the command is: 
$ ###   vt2geojson https://api.mapbox.com/v4/mapbox.mapbox-streets-v6/9/150/194.mvt?access_token=YOUR_MAPBOX_ACCESS_TOKEN >geojson.json


```