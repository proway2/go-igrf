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
		name    string
		args    args
		want    IGRFresults
		wantErr bool
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
				DeclinationSV:     9.6,
				InclinationSV:     2.2,
				HorizontalSV:      -16.3,
				TotalSV:           59.8,
				NorthSV:           -25.8,
				EastSV:            35.9,
				VerticalSV:        66.7,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IGRF(tt.args.lat, tt.args.lon, tt.args.alt, tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("IGRF() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IGRF() = %v, want %v", got, tt.want)
			}
		})
	}
}
