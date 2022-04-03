package igrf

import (
	"reflect"
	"testing"
)

func TestIGRFEdgeCases(t *testing.T) {
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
			name:    "Latitude below -90.0°",
			args:    args{lat: -90.1},
			want:    IGRFresults{},
			wantErr: true,
		},
		{
			name:    "Latitude above 90.0°",
			args:    args{lat: 90.1},
			want:    IGRFresults{},
			wantErr: true,
		},
		{
			name:    "Longitude below -180.0°",
			args:    args{lon: -180.1},
			want:    IGRFresults{},
			wantErr: true,
		},
		{
			name:    "Longitude above 180.0°",
			args:    args{lon: 180.1},
			want:    IGRFresults{},
			wantErr: true,
		},
		{
			name:    "Altitude below -1.0 km",
			args:    args{alt: -1.1},
			want:    IGRFresults{},
			wantErr: true,
		},
		{
			name:    "Altitude above 600.0 km",
			args:    args{alt: 600.1},
			want:    IGRFresults{},
			wantErr: true,
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
