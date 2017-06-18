package geo

import "math"

// UseHaversineGeoDistanceByDefault indicates if the more complicated
// Haversine formula should be used for distances.
var UseHaversineDistanceByDefault = false

//MinLatitude is the minimum possible latitude
var minLatitude = deg2rad(-90)

//MaxLatitude is the maxiumum possible latitude
var maxLatitude = deg2rad(90)

//MinLongitude is the minimum possible longitude
var minLongitude = deg2rad(-180)

//MaxLongitude is the maxiumum possible longitude
var maxLongitude = deg2rad(180)

// GeoHashPrecision is the number of characters of a encoded GeoHash.
var GeoHashPrecision = 12

func yesHaversine(haversine []bool) bool {
	return (len(haversine) != 0 && haversine[0]) || (UseHaversineDistanceByDefault && len(haversine) == 0)
}

func deg2rad(d float64) float64 {
	return d * math.Pi / 180.0
}

func rad2deg(r float64) float64 {
	return 180.0 * r / math.Pi
}
