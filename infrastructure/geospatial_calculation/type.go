package geospatialcalculation

type GeospatialCalculatorType interface {
	CalculateDistanceInMiles(long1 float64, long2 float64, lat1 float64, lat2 float64) (*float64)
}