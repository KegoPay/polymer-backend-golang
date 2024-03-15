package geospatialcalculation

import "github.com/jftuga/geodist"

type GeoDist struct {}

func (gd *GeoDist) CalculateDistanceInMiles(long1 float64, long2 float64, lat1 float64, lat2 float64) *float64 {
	cord1 := geodist.Coord{Lat: lat1, Lon: long1}
	cord2 := geodist.Coord{Lat: lat2, Lon: long2}
	milesApart, _ := geodist.HaversineDistance(cord1, cord2)
	return &milesApart
}