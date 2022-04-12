package coeffs

import "testing"

func TestIGRFcoeffs_findEpochs(t *testing.T) {
	type fields struct {
		names  *[]string
		epochs *[]float64
		lines  *[]lineData
		coeffs *map[string]*[]float64
	}
	type args struct {
		date float64
	}
	epochs := &[]float64{1900.0, 1905.0, 1910.0, 1915.0, 1920.0, 1925.0, 2025.0}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		{
			name:    "Below minimum",
			fields:  fields{epochs: epochs},
			args:    args{date: 1899.0},
			want:    "",
			want1:   "",
			wantErr: true,
		},
		{
			name:    "Above maximum",
			fields:  fields{epochs: epochs},
			args:    args{date: 1899.0},
			want:    "",
			want1:   "",
			wantErr: true,
		},
		{
			name:    "Exact match",
			fields:  fields{epochs: epochs},
			args:    args{date: 1905.0},
			want:    "1905.0",
			want1:   "1910.0",
			wantErr: false,
		},
		{
			name:    "Exact beginning",
			fields:  fields{epochs: epochs},
			args:    args{date: 1900.0},
			want:    "1900.0",
			want1:   "1905.0",
			wantErr: false,
		},
		{
			name:    "Close to ending",
			fields:  fields{epochs: epochs},
			args:    args{date: 2024.99},
			want:    "2020.0",
			want1:   "2025.0",
			wantErr: false,
		},
		{
			name:    "Middle match 1",
			fields:  fields{epochs: epochs},
			args:    args{date: 1923.49},
			want:    "1920.0",
			want1:   "1925.0",
			wantErr: false,
		},
		{
			name:    "Middle match 2",
			fields:  fields{epochs: epochs},
			args:    args{date: 1907.01},
			want:    "1905.0",
			want1:   "1910.0",
			wantErr: false,
		},
		{
			name:    "Middle match 3",
			fields:  fields{epochs: epochs},
			args:    args{date: 1924.99},
			want:    "1920.0",
			want1:   "1925.0",
			wantErr: false,
		},
		{
			name:    "Middle match 4",
			fields:  fields{epochs: epochs},
			args:    args{date: 1911.01},
			want:    "1910.0",
			want1:   "1915.0",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			igrf := &IGRFcoeffs{
				names:  tt.fields.names,
				epochs: tt.fields.epochs,
				lines:  tt.fields.lines,
				coeffs: tt.fields.coeffs,
			}
			got, got1, err := igrf.findEpochs(tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("IGRFcoeffs.findEpochs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IGRFcoeffs.findEpochs() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("IGRFcoeffs.findEpochs() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
