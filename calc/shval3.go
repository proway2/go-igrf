package calc

import "math"

// Calculates field components from spherical harmonic (sh) models.
// The calculation is performed for two sets of coeffs for a single location,
// thus it returns two sets of X, Y, Z.
//
// X - northward component
//
// Y - eastward component
//
// Z - vertically-downward component
func Shval3(flat, flon, elev float64, nmax int, gha, ghb *[]float64) (float64, float64, float64, float64, float64, float64) {
	// similar to shval3 from C implementation
	var earths_radius float64 = 6371.2
	var dtr float64 = 0.01745329
	/*
		a2,b2     - squares of semi-major and semi-minor axes of
		the reference spheroid used for transforming
		between geodetic and geocentric coordinates or components
	*/
	var a2 float64 = 40680631.59 /* WGS84 */
	var b2 float64 = 40408299.98 /* WGS84 */
	var x, y, z, xtemp, ytemp, ztemp, aa, aa_temp, argument, clat, slat, sd, bb, cc, dd, r, ratio, power, rr, fn, fm float64
	var l, n, m, npq int
	var sl, cl [14]float64
	var p, q [119]float64
	argument = flat * dtr
	slat = math.Sin(argument)
	if (90.0 - flat) < 0.001 {
		//  300 ft. from North pole
		aa = 89.999
	} else {
		if (90.0 + flat) < 0.001 {
			//  300 ft. from South pole
			aa = -89.999
		} else {
			aa = flat
		}
	}
	argument = aa * dtr
	clat = math.Cos(argument)
	argument = flon * dtr
	// TODO: Why this start from 1?
	sl[1] = math.Sin(argument)
	cl[1] = math.Cos(argument)
	l = 0 // in C index starts from 1
	n = 0
	m = 1
	npq = (nmax * (nmax + 3)) / 2

	// this block is for geodetic coordinate system ->
	aa = a2 * clat * clat
	bb = b2 * slat * slat
	cc = aa + bb
	argument = cc
	dd = math.Sqrt(argument)
	argument = elev*(elev+2.0*dd) + (a2*aa+b2*bb)/cc
	r = math.Sqrt(argument)
	cd := (elev + dd) / r
	sd = (a2 - b2) / dd * slat * clat / r
	aa = slat
	slat = slat*cd - clat*sd
	clat = clat*cd + aa*sd
	// <- this block is for geodetic coordinate system
	ratio = earths_radius / r
	argument = 3.0
	aa = math.Sqrt(argument)
	// TODO: Why all these starts with 1?
	p[1] = 2.0 * slat
	p[2] = 2.0 * clat
	p[3] = 4.5*slat*slat - 1.5
	p[4] = 3.0 * aa * clat * slat
	q[1] = -clat
	q[2] = slat
	q[3] = -3.0 * clat * slat
	q[4] = aa * (slat*slat - clat*clat)

	for k := 1; k <= npq; k++ {
		if n < m {
			m = 0
			n = n + 1
			argument = ratio
			power = float64(n + 2)
			rr = math.Pow(argument, power)
			fn = float64(n)
		}
		fm = float64(m)
		if k >= 5 {
			if m == n {
				argument = (1.0 - 0.5/fm)
				aa = math.Sqrt(argument)
				j := k - n - 1
				p[k] = (1.0 + 1.0/fm) * aa * clat * p[j]
				q[k] = aa * (clat*q[j] + slat/fm*p[j])
				sl[m] = sl[m-1]*cl[1] + cl[m-1]*sl[1]
				cl[m] = cl[m-1]*cl[1] - sl[m-1]*sl[1]
			} else {
				argument = fn*fn - fm*fm
				aa = math.Sqrt(argument)
				argument = ((fn - 1.0) * (fn - 1.0)) - (fm * fm)
				bb = math.Sqrt(argument) / aa
				cc = (2.0*fn - 1.0) / aa
				ii := k - n
				j := k - 2*n + 1
				p[k] = (fn + 1.0) * (cc*slat/fn*p[ii] - bb/(fn-1.0)*p[j])
				q[k] = cc*(slat*q[ii]-clat/fn*p[ii]) - bb*q[j]
			}
		}
		aa = rr * (*gha)[l]
		aa_temp = rr * (*ghb)[l]
		if m == 0 {
			x = x + aa*q[k]
			z = z - aa*p[k]
			xtemp = xtemp + aa_temp*q[k]
			ztemp = ztemp - aa_temp*p[k]
			l++
		} else {
			// ->
			bb = rr * (*gha)[l+1]
			cc = aa*cl[m] + bb*sl[m]
			x = x + cc*q[k]
			z = z - cc*p[k]
			if clat > 0 {
				y = y + (aa*sl[m]-bb*cl[m])*
					fm*p[k]/((fn+1.0)*clat)
			} else {
				y = y + (aa*sl[m]-bb*cl[m])*q[k]*slat
			}
			// <-

			// ->
			bb_temp := rr * (*ghb)[l+1]
			cc_temp := aa_temp*cl[m] + bb_temp*sl[m]
			xtemp = xtemp + cc_temp*q[k]
			ztemp = ztemp - cc_temp*p[k]
			if clat > 0 {
				ytemp = ytemp + (aa_temp*sl[m]-bb_temp*cl[m])*fm*p[k]/((fn+1.0)*clat)
			} else {
				ytemp = ytemp + (aa_temp*sl[m]-bb_temp*cl[m])*q[k]*slat
			}
			// <-
			l = l + 2
		}
		m++
	}

	aa = x
	x = x*cd + z*sd
	z = z*cd - aa*sd
	aa = xtemp
	xtemp = xtemp*cd + ztemp*sd
	ztemp = ztemp*cd - aa*sd
	return x, y, z, xtemp, ytemp, ztemp
}

// Computes the geomagnetic D, I, H, and F from X, Y, and Z.
//
// D  - declination
//
// I  - inclination
//
// H  - horizontal intensity
//
// F  - total intensity
func Dihf(x, y, z float64) (float64, float64, float64, float64) {
	var d, i, h, f float64
	sn := 0.0001
	for j := 1; j <= 1; j++ {
		h2 := x*x + y*y
		argument := h2
		// calculate horizontal intensity
		h = math.Sqrt(argument)
		argument = h2 + z*z
		// calculate total intensity
		f = math.Sqrt(argument)
		if f < sn {
			// If d and i cannot be determined
			d = math.NaN()
			// set equal to NaN
			i = math.NaN()
		} else {
			argument = z
			argument2 := h
			i = math.Atan2(argument, argument2)
			if h < sn {
				d = math.NaN()
			} else {
				hpx := h + x
				if hpx < sn {
					d = math.Pi
				} else {
					argument = y
					argument2 = hpx
					d = 2.0 * math.Atan2(argument, argument2)
				}
			}
		}
	}
	return d, i, h, f
}
