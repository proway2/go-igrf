package coeffs

import (
	"testing"
)

const secs_per_leap_year = 31622400
const secs_per_regular_year = 31536000

func almostEqual(a, b float64, threshold int) bool {
	a_int := uint32(a * float64(threshold))
	b_int := uint32(b * float64(threshold))
	return a_int == b_int
}

func Test_secsInYear(t *testing.T) {
	type args struct {
		year int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Non leap 1700",
			args: args{year: 1700},
			want: secs_per_regular_year,
		},
		{
			name: "Non leap 1800",
			args: args{year: 1800},
			want: secs_per_regular_year,
		},
		{
			name: "Non leap 1900",
			args: args{year: 1900},
			want: secs_per_regular_year,
		},
		{
			name: "Leap 1600",
			args: args{year: 1600},
			want: secs_per_leap_year,
		},
		{
			name: "Leap 2000",
			args: args{year: 2000},
			want: secs_per_leap_year,
		},
		{
			name: "Leap 2020",
			args: args{year: 2020},
			want: secs_per_leap_year,
		},
		{
			name: "Leap 2024",
			args: args{year: 2024},
			want: secs_per_leap_year,
		},
		{
			name: "Non leap 2021",
			args: args{year: 2021},
			want: secs_per_regular_year,
		},
		{
			name: "Non leap 2022",
			args: args{year: 2022},
			want: secs_per_regular_year,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := secsInYear(tt.args.year); got != tt.want {
				t.Errorf("secsInYear() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findDateFactor(t *testing.T) {
	type args struct {
		start_epoch string
		end_epoch   string
		date        float64
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name:    "Exact match (start epoch)",
			args:    args{start_epoch: "1900.0", end_epoch: "1910.0", date: 1900.0},
			want:    0.0,
			wantErr: false,
		},
		{
			name:    "Exact match (end epoch)",
			args:    args{start_epoch: "1900.0", end_epoch: "1910.0", date: 1910.0},
			want:    0.0,
			wantErr: false,
		},
		{
			name:    "End epoch is less than start epoch",
			args:    args{start_epoch: "1910.0", end_epoch: "1900.0", date: 1910.0},
			want:    0.0,
			wantErr: false,
		},
		{
			name:    "Middle",
			args:    args{start_epoch: "1900.0", end_epoch: "1910.0", date: 1905.0},
			want:    0.5,
			wantErr: false,
		},
		{
			name:    "1950.01",
			args:    args{start_epoch: "1950.0", end_epoch: "1955.0", date: 1950.01},
			want:    0.001998904709746265,
			wantErr: false,
		},
		{
			name:    "1950.99",
			args:    args{start_epoch: "1950.0", end_epoch: "1955.0", date: 1954.99},
			want:    0.9980010952902538,
			wantErr: false,
		},
		{
			name:    "2025.5",
			args:    args{start_epoch: "2020.0", end_epoch: "2025.0", date: 2025.5},
			want:    1.1,
			wantErr: false,
		},
		{
			name:    "Incorrect start_epoch",
			args:    args{start_epoch: "start_epoch", end_epoch: "2025.0", date: 2025.5},
			want:    -999,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findDateFactor(tt.args.start_epoch, tt.args.end_epoch, tt.args.date)
			if (err != nil) && !tt.wantErr {
				t.Errorf("findDateFactor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("findDateFactor() = %v, want %v", got, tt.want)
			}
		})
	}
}
