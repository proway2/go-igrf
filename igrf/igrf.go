package igrf

// IGRF computes values for the geomagnetic field and secular variation for a given set of coordinates and date.
// lat, lon - geodetic latitude and longitude (WGS84 latitude and altitude above mean sea level).
// Valid values -90.0 < lat < 90.0, -180.0 < lon < 180.0.
// alt - geodetic altitude above mean sea level in km (-1.00 to 600.00).
// date - decimal date (1900.00 to 2025).
func IGRF(lat, lon, alt, date float32) (IGRFresults, error) {
	res := IGRFresults{}
	return res, nil
}
