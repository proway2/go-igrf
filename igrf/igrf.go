package igrf

import (
	"errors"
	"fmt"
	"log"

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
	// colat := float64(90.0 - lat)
	// alt64, colat, sd, cd := gg2geo(float64(alt), float64(colat))
	shc, err := coeffs.NewCoeffsData()
	if err != nil {
		log.Fatal(err.Error())
	}
	start_coeffs, end_coeffs, nmax, err := shc.Coeffs(date)
	if err != nil {
		return IGRFresults{}, err
	}
	x, y, z, xtemp, ytemp, ztemp := calc.Shval3(lat, lon, alt, nmax, start_coeffs, end_coeffs)
	_ = xtemp
	_ = ytemp
	_ = ztemp
	res := IGRFresults{NorthComponent: float32(x), EastComponent: float32(y), VerticalComponent: float32(z)}
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

// // gg2geo - computes geocentric colatitude and radius from geodetic colatitude and height. Uses WGS-84 ellipsoid parameters.
// //
// // Inputs:
// // h - altitude in kilometers
// // gdcolat - geodetic colatitude
// //
// // Outputs:
// // radius - Geocentric radius in kilometers.
// // theta - Geocentric colatitude in degrees.
// // sd - rotate B_X to gd_lat
// // cd - rotate B_Z to gd_lat
// //
// // References:
// // Equations (51)-(53) from "The main field" (chapter 4) by Langel, R. A. in: "Geomagnetism", Volume 1, Jacobs, J. A., Academic Press, 1987.
// // Malin, S.R.C. and Barraclough, D.R., 1981. An algorithm for synthesizing the geomagnetic field. Computers & Geosciences, 7(4), pp.401-405.
// func gg2geo(h, gdcolat float64) (radius, theta, sd, cd float64) {
// 	eqrad := 6378.137 // equatorial radius
// 	flat := 1 / 298.257223563
// 	plrad := eqrad * (1 - flat) // polar radius
// 	ctgd := math.Cos(deg2rad(gdcolat))
// 	stgd := math.Sin(deg2rad(gdcolat))

// 	a2 := eqrad * eqrad
// 	a4 := a2 * a2
// 	b2 := plrad * plrad
// 	b4 := b2 * b2
// 	c2 := ctgd * ctgd
// 	s2 := 1 - c2
// 	rho := math.Sqrt(a2*s2 + b2*c2)

// 	rad := math.Sqrt(h*(h+2*rho) + (a4*s2+b4*c2)/math.Pow(rho, 2))

// 	cd = (h + rho) / rad
// 	sd = (a2 - b2) * ctgd * stgd / (rho * rad)

// 	cthc := ctgd*cd - stgd*sd        // Also: sthc = stgd*cd + ctgd*sd
// 	theta = rad2deg(math.Acos(cthc)) // arccos returns values in [0, pi]

// 	return rad, theta, sd, cd
// }

// // deg2rad - converts `degrees` into radians.
// func deg2rad(degrees float64) float64 {
// 	rad := degrees * math.Pi / 180.0
// 	return rad
// }

// // rad2deg - converts `radians` into degrees.
// func rad2deg(radians float64) float64 {
// 	deg := radians * 180.0 / math.Pi
// 	return deg
// }
