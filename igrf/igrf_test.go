package igrf

import (
	"reflect"
	"testing"
)

func TestIGRF(t *testing.T) {
	type args struct {
		lat  float32
		lon  float32
		alt  float32
		date float32
	}
	tests := []struct {
		name string
		args args
		want IGRFresults
	}{
		{
			name: "lat: 59.9, lon: 39.9, alt: 0.0, date: 2019.123",
			args: args{lat: 59.9, lon: 39.9, alt: 0.0, date: 2019.123},
			want: IGRFresults{
				Declination:       14.033,
				Inclination:       74.197,
				HzIntensity:       14641.1,
				TotalIntensity:    53761.1,
				NorthComponent:    14204.1,
				EastComponent:     3550.1,
				VerticalComponent: 51729.1,
				DeclinationSV:     -9.55,
				InclinationSV:     -2.22,
				HorizontalSV:      -16.8,
				TotalSV:           59.5,
				NorthSV:           -25.8,
				EastSV:            35.9,
				VerticalSV:        66.7,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IGRF(tt.args.lat, tt.args.lon, tt.args.alt, tt.args.date); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IGRF() = %v, want %v", got, tt.want)
			}
		})
	}
}
