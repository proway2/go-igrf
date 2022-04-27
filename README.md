[![Checking the project](https://github.com/proway2/go-igrf/actions/workflows/main.yml/badge.svg)](https://github.com/proway2/go-igrf/actions/workflows/main.yml)
[![golangci-lint](https://github.com/proway2/go-igrf/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/proway2/go-igrf/actions/workflows/golangci-lint.yml)

# go-igrf
Pure Go IGRF (International Geomagnetic Reference Field). This is based on the existing `C` implementation.

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

## References

- [Alken, P., Thébault, E., Beggan, C.D. et al. International Geomagnetic Reference Field: the thirteenth generation. Earth Planets Space 73, 49 (2021).](https://rdcu.be/cKqZv) https://doi.org/10.1186/s40623-020-01288-x

- Implementation is based on [geomag70.c](https://www.ngdc.noaa.gov/IAGA/vmod/geomag70_linux.tar.gz). [License information](https://www.ngdc.noaa.gov/IAGA/vmod/geomag70_license.html).