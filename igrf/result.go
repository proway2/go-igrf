package igrf

type IGRFresults struct {
	// Declination (D):  14.175 °
	Declination float32
	// Inclination (I):  74.229 °
	Inclination float32
	// Horizontal intensity (H):  14626.6 nT
	HzIntensity float32
	// Total intensity (F):  53814.3 nT
	TotalIntensity float32
	// North component (X):  14181.2 nT
	NorthComponent float32
	// East component (Y):  3581.8 nT
	EastComponent float32
	// Vertical component (Z):  51788.4 nT
	VerticalComponent float32
	// Declination SV (D): -8.92 arcmin/yr
	DeclinationSV float32
	// Inclination SV (I): -2.59 arcmin/yr
	InclinationSV float32
	// Horizontal SV (H):  -16.8 nT/yr
	HorizontalSV float32
	// Total SV (F):  71.8 nT/yr
	TotalSV float32
	// North SV (X): -28.2 nT/yr
	NorthSV float32
	// East SV (Y):  32.0 nT/yr
	EastSV float32
	// Vertical SV (Z):  80.1 nT/yr
	VerticalSV float32
}
