package igrf

// IGRFresults represents the result of IGRF calculation
//
// Fields:
// Declination (D):  14.175 ° (+ve east)
// DeclinationSV (D): -8.92 arcmin/yr
// Inclination (I):  74.229 ° (+ve down)
// InclinationSV (I): -2.59 arcmin/yr
// HzIntensity (Horizontal intensity (H)):  14626.6 nT
// HorizontalSV (H):  -16.8 nT/yr
// NorthComponent (X):  14181.2 nT
// NorthSV (X): -28.2 nT/yr
// EastComponent (Y):  3581.8 nT
// EastSV (Y):  32.0 nT/yr
// VerticalComponent (Z):  51788.4 nT
// VerticalSV (Z):  80.1 nT/yr
// TotalIntensity (F):  53814.3 nT
// TotalSV (F):  71.8 nT/yr
type IGRFresults struct {
	Declination         float32
	DeclinationSV       float32
	Inclination         float32
	InclinationSV       float32
	HorizontalIntensity float32
	HorizontalSV        float32
	NorthComponent      float32
	NorthSV             float32
	EastComponent       float32
	EastSV              float32
	VerticalComponent   float32
	VerticalSV          float32
	TotalIntensity      float32
	TotalSV             float32
}
