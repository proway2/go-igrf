[![Checking the project](https://github.com/proway2/go-igrf/actions/workflows/main.yml/badge.svg)](https://github.com/proway2/go-igrf/actions/workflows/main.yml)
[![golangci-lint](https://github.com/proway2/go-igrf/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/proway2/go-igrf/actions/workflows/golangci-lint.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/proway2/go-igrf.svg)](https://pkg.go.dev/github.com/proway2/go-igrf)
[![Go Report Card](https://goreportcard.com/badge/github.com/proway2/go-igrf)](https://goreportcard.com/report/github.com/proway2/go-igrf)
[![IGRF badge](https://badgen.net/static/IGRF%20model/14th/green)](https://www.ncei.noaa.gov/products/international-geomagnetic-reference-field)

# go-igrf
Pure Go IGRF (International Geomagnetic Reference Field). This is based on the existing `C` implementation. This package computes values for the geomagnetic field and secular variation for a given set of coordinates and date.

## Inputs

- `lat` geodetic latitude (WGS84 latitude) in decimal degrees, valid values -90.0 < lat < 90.0;
- `lon` geodetic longitude (WGS84 longitude) in decimal degrees, valid values -180.0 < lon < 180.0;
- `alt` geodetic altitude above mean sea level in km (-1.00 to 600.00);
- `date` decimal date, starting from 1900.00.

## Output

The output is of type `type IGRFresults struct`. Fields are:

|Field name|Fortran name|Unit of measurement|Sample value|Notes|
|-|-|-|-|-|
|Declination|(D)|decimal degrees(°)|14.175 °|(+ve east)|
|DeclinationSV|(D)|arcmin/yr|-8.92||
|Inclination|(I)|decimal degrees(°)|74.229|(+ve down)
|InclinationSV|(I)|arcmin/yr|-2.59||
|HzIntensity|Horizontal intensity (H)|nT|14626.6||
|HorizontalSV|(H)|nT/yr|-16.8||
|NorthComponent|(X)|nT|14181.2||
|NorthSV|(X)|nT/yr|-28.2||
|EastComponent|(Y)|nT|3581.8||
|EastSV|(Y)|nT/yr|32.0||
|VerticalComponent|(Z)|nT|51788.4|downward|
|VerticalSV|(Z)|nT/yr|80.1||
|TotalIntensity|(F)|nT|53814.3||
|TotalSV|(F)|nT/yr|71.8||

### Annual changes (SV values)

The SV values are computed by subtracting the values for the desired input date from corresponding values one year later.

### Values near geographic poles

Unlike in `C` implementation, this software calculates values near geographic poles, e.g. for latitudes higher than 89.999 and -89.999. This is the same approach the reference `FORTRAN` implementaion does have. Be adviced that near pole values are much less accurate.

## Overall accuracy

There are far more than 1000+ unittests and results are compared against those generated from FORTRAN. SV values are not covered with tests due to the initial low accuracy of `FORTRAN` values. For most values these tolerances are used (whichever is higher):
- relative tolerance: 0.005
- absolute tolerance: 0.15
- absolute tolerance for declination and inclination: 0.005

Near pole values are tested with relaxed accuracies. Since FORTRAN rounds almost all values, except `D` and `I`, actual results are of much higher accuracy.

## How to use

- Add `"github.com/proway2/go-igrf/igrf"` to your import clause.

```go
package main

import (
	"fmt"
	"github.com/proway2/go-igrf/igrf"
)

func main() {
	igrf_data := igrf.New()
	res, err := igrf_data.IGRF(46.9, 39.9, 0.0, 2021.5)
	fmt.Println(res, err)
}
```

- Run `go mod tidy`, this brings the latest version. Fix the version at `go.mod` if you need a different one.

## References

- [Alken, P., Thébault, E., Beggan, C.D. et al. International Geomagnetic Reference Field: the thirteenth generation. Earth Planets Space 73, 49 (2021).](https://rdcu.be/cKqZv) https://doi.org/10.1186/s40623-020-01288-x

- Implementation is based on [geomag70.c](https://www.ngdc.noaa.gov/IAGA/vmod/geomag70_linux.tar.gz). [License information](https://www.ngdc.noaa.gov/IAGA/vmod/geomag70_license.html).
