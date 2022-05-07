package igrf

// IGRFresults represents the result of IGRF calculation
//
// Fields:
//
// Declination (D):  14.175 ° (+ve east)
//
// DeclinationSV (D): -8.92 arcmin/yr
//
// Inclination (I):  74.229 ° (+ve down)
//
// InclinationSV (I): -2.59 arcmin/yr
//
// HzIntensity (Horizontal intensity (H)):  14626.6 nT
//
// HorizontalSV (H):  -16.8 nT/yr
//
// NorthComponent (X):  14181.2 nT
//
// NorthSV (X): -28.2 nT/yr
//
// EastComponent (Y):  3581.8 nT
//
// EastSV (Y):  32.0 nT/yr
//
// VerticalComponent (Z):  51788.4 nT
//
// VerticalSV (Z):  80.1 nT/yr
//
// TotalIntensity (F):  53814.3 nT
//
// TotalSV (F):  71.8 nT/yr
type IGRFresults struct {
	Declination         float64
	DeclinationSV       float64
	Inclination         float64
	InclinationSV       float64
	HorizontalIntensity float64
	HorizontalSV        float64
	NorthComponent      float64
	NorthSV             float64
	EastComponent       float64
	EastSV              float64
	VerticalComponent   float64
	VerticalSV          float64
	TotalIntensity      float64
	TotalSV             float64
}
