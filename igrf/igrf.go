package igrf

import (
	"errors"
	"fmt"
)

// IGRF computes values for the geomagnetic field and secular variation for a given set of coordinates and date.
// lat, lon - geodetic latitude and longitude (WGS84 latitude and altitude above mean sea level).
// Valid values -90.0 < lat < 90.0, -180.0 < lon < 180.0.
// alt - geodetic altitude above mean sea level in km (-1.00 to 600.00).
// date - decimal date (1900.00 to 2025).
func IGRF(lat, lon, alt, date float32) (IGRFresults, error) {
	if err := checkInitialConditions(lat, lon, alt); err != nil {
		return IGRFresults{}, err
	}
	res := IGRFresults{}
	return res, nil
}

func checkInitialConditions(lat, lon, alt float32) error {
	var error_msg string
	if lat < -90.0 || lat > 90.0 {
		error_msg = fmt.Sprintf("Latitude %v° is out of range (-90.0, 90.0)", lat)
	}
	if lon < -180.0 || lon > 180.0 {
		error_msg = fmt.Sprintf("Latitude %v° is out of range (-90.0, 90.0)", lat)
	}
	if alt < -1.0 || alt > 600.0 {
		error_msg = fmt.Sprintf("Altitude %v km is out of range (-1.0, 600.0)", alt)
	}
	if len(error_msg) != 0 {
		return errors.New(error_msg)
	}
	return nil
}
