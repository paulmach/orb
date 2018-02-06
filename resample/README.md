orb/resample [![Godoc Reference](https://godoc.org/github.com/paulmach/orb/resample?status.svg)](https://godoc.org/github.com/paulmach/orb/resample)
============

Package orb/resample has a couple functions for resampling line geometry
into more or less evenly spaces points.

	func Resample(ls orb.LineString, df orb.DistanceFunc, totalPoints int) orb.LineString
	func ToInterval(ls orb.LineString, df orb.DistanceFunc, dist float64) orb.LineString

For example, resampling a line string so the points are 1 planar unit apart:

	ls := resample.ToInterval(ls, planar.Distance, 1.0)
