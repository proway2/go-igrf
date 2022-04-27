package igrf

import (
	"errors"
	"fmt"
	"log"
	"math"

	"github.com/proway2/go-igrf/calc"
	"github.com/proway2/go-igrf/coeffs"
)

// IGRF computes values for the geomagnetic field and secular variation for a given set of coordinates and date.
// lat, lon - geodetic latitude and longitude (WGS84 latitude and altitude above mean sea level).
// Valid values -90.0 < lat < 90.0, -180.0 < lon < 180.0.
// alt - geodetic altitude above mean sea level in km (-1.00 to 600.00).
// date - decimal date (1900.00 to 2025).
func IGRF(lat, lon, alt, date float64) (IGRFresults, error) {
	if err := checkInitialConditions(lat, lon, alt); err != nil {
		return IGRFresults{}, err
	}
	shc, err := coeffs.NewCoeffsData()
	if err != nil {
		log.Fatal(err.Error())
	}
	start_coeffs, end_coeffs, nmax, err := shc.Coeffs(date)
	if err != nil {
		return IGRFresults{}, err
	}
	x, y, z, xtemp, ytemp, ztemp := calc.Shval3(lat, lon, alt, nmax, start_coeffs, end_coeffs)
	d, i, h, f := calc.Dihf(x, y, z)

	dtemp, itemp, htemp, ftemp := calc.Dihf(xtemp, ytemp, ztemp)

	ddot := rad2deg(dtemp - d)
	if ddot > 180.0 {
		ddot -= 360.0
	}
	if ddot <= -180.0 {
		ddot += 360.0
	}
	ddot *= 60.0

	idot := rad2deg(itemp-i) * 60
	d = rad2deg(d)
	i = rad2deg(i)
	hdot := htemp - h
	xdot := xtemp - x
	ydot := ytemp - y
	zdot := ztemp - z
	fdot := ftemp - f

	// deal with geographic and magnetic poles

	// at magnetic poles
	if h < 100.0 {
		d = math.NaN()
		ddot = math.NaN()
		/* while rest is ok */
	}
	// warn_H := 0
	warn_H_val := 99999.0
	var warn_H_strong int
	warn_H_strong_val := 99999.0
	// warn_P := 0
	if h < 1000.0 {
		// warn_H = 0
		warn_H_strong = 1
		if h < warn_H_strong_val {
			warn_H_strong_val = h
		}
		// } else if h < 5000.0 && !warn_H_strong {
	} else if h < 5000.0 && warn_H_strong != 0 {
		// warn_H = 1
		if h < warn_H_val {
			warn_H_val = h
		}
	}

	// at geographic poles
	if 90.0-math.Abs(lat) <= 0.001 {
		x = math.NaN()
		y = math.NaN()
		d = math.NaN()
		xdot = math.NaN()
		ydot = math.NaN()
		ddot = math.NaN()
		// warn_P = 1
		// warn_H = 0
		warn_H_strong = 0
		/* while rest is ok */
	}

	res := IGRFresults{
		Declination:         float32(d),
		DeclinationSV:       float32(ddot),
		Inclination:         float32(i),
		InclinationSV:       float32(idot),
		HorizontalIntensity: float32(h),
		HorizontalSV:        float32(hdot),
		NorthComponent:      float32(x),
		NorthSV:             float32(xdot),
		EastComponent:       float32(y),
		EastSV:              float32(ydot),
		VerticalComponent:   float32(z),
		VerticalSV:          float32(zdot),
		TotalIntensity:      float32(f),
		TotalSV:             float32(fdot),
	}
	return res, nil
}

func checkInitialConditions(lat, lon, alt float64) error {
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

// rad2deg - converts `radians` into degrees.
func rad2deg(radians float64) float64 {
	deg := radians * 180.0 / math.Pi
	return deg
}
