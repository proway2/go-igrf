package coeffs

import "testing"

const secs_per_leap_year = 31622400
const secs_per_regular_year = 31536000

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
