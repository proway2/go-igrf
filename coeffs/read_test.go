package coeffs

import (
	"math"
	"testing"
)

func TestIGRFcoeffs_findEpochs(t *testing.T) {
	type fields struct {
		names  *[]string
		epochs *[]float64
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

func TestIGRFcoeffs_Coeffs(t *testing.T) {
	igrf, _ := NewCoeffsData()
	type args struct {
		date float64
	}
	tests := []struct {
		name    string
		args    args
		want    *[]float64
		wantErr bool
	}{
		{
			name:    "1920.5: Coeffs when degrees are equal (10), but coeffs are not definitive",
			args:    args{date: 1920.5},
			want:    &[]float64{-31046.599999999999, -2317.0999999999999, 5842.1999999999998, -844.39999999999998, 2960, -1266.5, 1413.4000000000001, 813.5, 1113.9000000000001, -1604.5, -446.69999999999999, 1204.7, 104.59999999999999, 843.20000000000005, 286.60000000000002, 889.20000000000005, 696.60000000000002, 219.59999999999999, 614.5, -136.90000000000001, -424.19999999999999, -150.69999999999999, 200.80000000000001, -58.299999999999997, -221.90000000000001, 326, -119.40000000000001, 235, 58, -23.5, -38.600000000000001, -119.59999999999999, -124.7, -62.700000000000003, 43.799999999999997, 61, 54.899999999999999, 0.30000000000000004, -9.9000000000000004, 96.299999999999997, -233.5, 11.300000000000001, -45.399999999999999, -21.600000000000001, 43.5, 17.5, -101.2, -56.5, 73, -54, -49.100000000000001, 2.1000000000000001, -14, 28.800000000000001, -13.1, -36.799999999999997, 4.0999999999999996, -15.800000000000001, 28.100000000000001, 19, -16.100000000000001, 6, -21.899999999999999, 11, 7, 8, -3, -15, -9, 6, 2, -14, 4, 5, -7, 17, 6.0999999999999996, -5, 8, -19, 8, 10, -20, 1, 14, -11, 5, 12, -3, 1, -2, -2, 9, 2, 10, 0, -2, -1, 2, -3, -4, 2, 2, 1, -5, 2, -2, 6, 6, -4, 4, 0, 0, -2, 1, 4, 3, 0, 0, -6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name:    "1944.5: Coeffs when degrees are equal (10), and coeffs are between non-definitive and definitive",
			args:    args{date: 1944.5},
			want:    &[]float64{-30600, -2285.6999999999998, 5811.1000000000004, -1230.2, 2989.0999999999999, -1693.2, 1576.8, 482.10000000000002, 1277.8, -1829.5999999999999, -499, 1252.7, 183.69999999999999, 913.29999999999995, -5.6000000000000014, 941, 774.60000000000002, 146.5, 544.60000000000002, -273.60000000000002, -419.39999999999998, -56.700000000000003, 300.10000000000002, -174.30000000000001, -251.80000000000001, 344.80000000000001, -14.099999999999998, 195.40000000000001, 92.599999999999994, -21.299999999999997, -67.799999999999997, -141.90000000000001, -118.40000000000001, -81.400000000000006, 80.700000000000003, 58.799999999999997, 56.700000000000003, 5.7999999999999998, 4.7000000000000011, 100.5, -246.30000000000001, 17.699999999999999, -24.300000000000001, -9.5999999999999996, 20.699999999999999, -14.4, -104.3, -38.399999999999999, 70.400000000000006, -41.299999999999997, -45.700000000000003, 0.39999999999999991, -18, 2, 0.40000000000000036, -29.199999999999999, 6.0999999999999996, -9.9000000000000004, 28.100000000000001, 15.199999999999999, -17.300000000000001, 26.600000000000001, -21.699999999999999, 12.800000000000001, 7, 11.6, -7.5, -20.300000000000001, -5.5, -10.300000000000001, 8.1999999999999993, -7.7999999999999998, 6.9000000000000004, 2.2999999999999998, -9.5, 18.100000000000001, 7.2000000000000002, 2.2000000000000002, 2.5, -11.800000000000001, 5.2999999999999998, -17.900000000000002, -26.399999999999999, 1, 16.800000000000001, -11.1, 26.600000000000001, 3.7999999999999998, -8.4000000000000004, 14.5, 3.2999999999999998, -2.8999999999999999, 9, -3.2999999999999998, 6.5, -2.6000000000000001, 0.70000000000000018, -3.7999999999999998, 7.4000000000000004, -3, 9.5, 4.7000000000000002, 1.1000000000000001, 1, 1.2999999999999998, -17.800000000000001, -4.7000000000000002, -0.29999999999999982, -0.29999999999999982, -5.7999999999999998, 7.5999999999999996, 5.4000000000000004, -0.90000000000000002, -3.7000000000000002, -2.5, -1.4000000000000004, 4.7999999999999998, 0, -1.8, -2.3999999999999999, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name:    "1955.5: Coeffs when degrees are equal (10), and coeffs are definitive - 1",
			args:    args{date: 1955.5},
			want:    &[]float64{-30492.099999999999, -2210.4000000000001, 5817.1000000000004, -1451.5, 3002.9000000000001, -1904.9000000000001, 1581.9000000000001, 282.5, 1302, -1948.8, -457.19999999999999, 1288.0999999999999, 216.80000000000001, 881.60000000000002, -87.700000000000003, 957.89999999999998, 796.39999999999998, 133.19999999999999, 509.39999999999998, -274.39999999999998, -396.69999999999999, -20.399999999999999, 287.89999999999998, -232.5, -228.30000000000001, 360.19999999999999, 15.1, 231.19999999999999, 111.5, -23.300000000000001, -99.900000000000006, -152.40000000000001, -120.3, -68.400000000000006, 78.299999999999997, 46.899999999999999, 57.100000000000001, -9.0999999999999996, 2.7999999999999998, 96.299999999999997, -246, 49.200000000000003, -7.2999999999999998, -16.399999999999999, 6.0999999999999996, -11.9, -107.59999999999999, -23.300000000000001, 65.200000000000003, -56, -50.5, 2.2999999999999998, -24.399999999999999, 10.5, -4.2000000000000002, -32, 7.9000000000000004, -10.6, 27.5, 9.8000000000000007, -19.800000000000001, 17, -17.899999999999999, 11.4, 8.6999999999999993, 10.1, -5.7999999999999998, -14.9, -13.699999999999999, 5.2000000000000002, 5.5999999999999996, -22.5, 10, 3.1000000000000001, -6.7999999999999998, 23, 6.4000000000000004, -3.5, 8.9000000000000004, -13.699999999999999, 4, 8.6999999999999993, -11.699999999999999, -3.6000000000000001, 12, -5.4000000000000004, 6.5, 1.8999999999999999, 5.4000000000000004, 4, -2.1000000000000001, 0.80000000000000004, 9.9000000000000004, 1.6000000000000001, 7.0999999999999996, 2.1000000000000001, -5.4000000000000004, 4.4000000000000004, 5, -2.6000000000000001, -4.7999999999999998, -3.2000000000000002, -0.5, 0.10000000000000001, 1.8, -7.2000000000000002, -2.7999999999999998, -1.6000000000000001, 6.7000000000000002, -4.0999999999999996, 4.2000000000000002, 1, -1.7, -2.7999999999999998, 5.2999999999999998, 6.9000000000000004, -1.6000000000000001, -0.90000000000000002, 0, -3.3999999999999999, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name:    "2017.5: Coeffs when degrees are equal (13), and coeffs are definitive - 2",
			args:    args{date: 2017.5},
			want:    &[]float64{-29423.129999999997, -1476.335, 4724.2449999999999, -2472.7399999999998, 2997.0999999999999, -2918.5050000000001, 1676.675, -688.38499999999999, 1356.7649999999999, -2366.73, -98.694999999999993, 1231.0250000000001, 243.47, 553.69500000000005, -541.04999999999995, 905.21000000000004, 811.58999999999992, 282.72000000000003, 103.395, -173.41500000000002, -322.125, 190.32499999999999, 59.189999999999998, -339.46500000000003, -233.60500000000002, 361.66999999999996, 47.340000000000003, 190.07499999999999, 202.63999999999999, -140.81999999999999, -120.17, -154.30000000000001, 24.140000000000001, 8.8999999999999986, 99.510000000000005, 67.775000000000006, 66.534999999999997, -19.855, 72.844999999999999, 29.199999999999999, -125.675, 55.769999999999996, -32.564999999999998, -65.569999999999993, 13.32, 8.125, -67.775000000000006, 65.254999999999995, 80.944999999999993, -76.344999999999999, -52.885000000000005, -7.4949999999999992, -18.215, 54.159999999999997, 3.895, 15.435, 23.975000000000001, 7.8600000000000003, 0.5349999999999997, -5.04, -27.350000000000001, 8.2050000000000001, -2.0600000000000001, 23.84, 9.2949999999999999, 9.2199999999999989, -17.190000000000001, -16.780000000000001, -1.8300000000000001, 12.99, -20.829999999999998, -13.149999999999999, 14.315000000000001, 15.530000000000001, 12.73, 4.6450000000000005, -16.240000000000002, -8, -1.1600000000000001, 2.5299999999999998, 5.165, 8.6150000000000002, -22.585000000000001, 2.96, 10.879999999999999, -2.3600000000000003, 10.77, -0.21499999999999997, -5.9199999999999999, -13.199999999999999, -6.5899999999999999, 0.50000000000000011, 7.7949999999999999, 8.7400000000000002, 0.71999999999999997, -9.1799999999999997, -2.645, -11.219999999999999, 9.0199999999999996, -1.9549999999999998, -6.2300000000000004, 3.3399999999999999, 0.035000000000000003, -0.30000000000000004, 1.125, 4.0750000000000002, -0.72500000000000009, 4.5999999999999996, 1.2, -8.2599999999999998, -0.78500000000000003, -0.35499999999999998, 2.0149999999999997, -4.2300000000000004, 1.865, -3.125, -2.1000000000000001, -0.6100000000000001, -3.6949999999999998, -8.7600000000000016, 3, -1.3999999999999999, 0, -2.3999999999999999, 2.3049999999999997, 2.1899999999999999, -0.59999999999999998, -0.84499999999999997, -0.72500000000000009, 0.43999999999999995, 0.67999999999999994, -0.69999999999999996, -0.20000000000000001, 0.020000000000000004, -1.9100000000000001, 1.5499999999999998, -1.52, -0.41000000000000003, -2.7850000000000001, 0.32000000000000001, -2.0049999999999999, 3.2949999999999999, -2.4699999999999998, -2.0449999999999999, -0.13, -1.1400000000000001, 0.47999999999999998, 0.435, 1.2650000000000001, 1.575, -1.0449999999999999, -1.9950000000000001, 0.77499999999999991, 0.185, 0.20000000000000001, 0.76000000000000001, 0.52000000000000002, -0.14500000000000002, -0.33499999999999996, 0.44499999999999995, -0.46499999999999997, 0.21500000000000002, 0.16, -0.89500000000000002, -1.02, -0.080000000000000002, -0.16500000000000001, 0.60999999999999999, 0.040000000000000008, -0.91000000000000003, -0.89000000000000001, 0.45999999999999996, 0.54499999999999993, 0.66500000000000004, 1.48, -0.35999999999999999, -0.45000000000000001, 0.88, -1.27, -0.095000000000000001, -0.10000000000000001, 0.80500000000000005, 0.35999999999999999, -0.065000000000000002, -0.070000000000000007, 0.39000000000000001, 0.48999999999999999, 0.089999999999999997, 0.48999999999999999, 0.47999999999999998, -0.34999999999999998, -0.42499999999999999, -0.41500000000000004, -0.38, -0.65500000000000003},
			wantErr: false,
		},
		{
			name:    "1999.5: Coeffs when degrees are not equal (10 > 13)",
			args:    args{date: 1999.5},
			want:    &[]float64{-29626.66, -1733.78, 5198.0900000000001, -2260.9299999999998, 3068.5599999999999, -2470.04, 1671.9100000000001, -453.5, 1339.1399999999999, -2285.9000000000001, -231.03999999999999, 1251.79, 294.25999999999999, 718.95000000000005, -484.69, 933.06999999999994, 786.12, 271.54000000000002, 254, -232.31, -404.5, 117.52, 112.37, -304.01999999999998, -218.32000000000002, 351.45999999999998, 44.019999999999996, 223.57000000000002, 171.21000000000001, -129.16, -134.09, -168.34, -40.869999999999997, -13.31, 106.37, 71.870000000000005, 68.079999999999998, -17.359999999999999, 73.579999999999998, 64.530000000000001, -161.81, 65.289999999999992, -5.4100000000000001, -60.880000000000003, 17.109999999999999, 0.72999999999999998, -90.660000000000011, 43.019999999999996, 78.799999999999997, -73.799999999999997, -65.039999999999992, 0.099999999999999978, -24.280000000000001, 32.769999999999996, 5.9800000000000004, 8.6899999999999995, 24, 6.6100000000000003, 15.020000000000001, 7.3700000000000001, -25.259999999999998, -1.2799999999999998, -5.8200000000000003, 24.459999999999997, 6.54, 11.81, -8.879999999999999, -21.449999999999999, -8.0099999999999998, 8.4499999999999993, -16.34, -21.649999999999999, 9.0899999999999999, 15.449999999999999, 6.9000000000000004, 9.1099999999999994, -7.6100000000000003, -15.01, -7, -2.29, 4.9000000000000004, 9.3599999999999994, -19.73, 3, 13.56, -8.5600000000000005, 12.449999999999999, 6.4699999999999998, -6.1799999999999997, -8.8100000000000005, -8.3599999999999994, -1.45, 8.3599999999999994, 9.370000000000001, 3.9199999999999999, -4.0700000000000003, -8.1799999999999997, -8.1799999999999997, 4.6200000000000001, -2.6400000000000001, -6, 1.6299999999999999, 1.73, 0, -3.1899999999999999, 4, -0.55000000000000004, 4.9100000000000001, 3.73, -5.8100000000000005, 1.1000000000000001, -1.1799999999999999, 2, -2.8100000000000001, 4.2800000000000002, 0.27999999999999992, 0.37, -2.1800000000000002, -0.9900000000000001, -7.3600000000000003, 2.4300000000000002, -1.53, 0.090000000000000011, -1.71, 1.1700000000000002, 1.3500000000000001, -0.81000000000000005, -0.090000000000000011, -2.3400000000000003, 0.090000000000000011, 0.81000000000000005, -0.63, -0.63, 0.63, -2.52, 1.53, -0.81000000000000005, 0.090000000000000011, -1.0800000000000001, 1.0800000000000001, -1.71, 3.6000000000000001, -0.81000000000000005, -1.9800000000000002, -0.27000000000000002, -0.36000000000000004, 0.18000000000000002, 0.27000000000000002, 0.81000000000000005, 2.25, -0.18000000000000002, -2.3400000000000003, 0.81000000000000005, 0.63, -0.45000000000000001, 0.27000000000000002, 0.27000000000000002, 0, -0.27000000000000002, 0, -0.36000000000000004, 0.27000000000000002, -0.090000000000000011, -0.81000000000000005, -0.18000000000000002, -0.36000000000000004, -0.36000000000000004, 0.72000000000000008, -0.18000000000000002, -0.81000000000000005, -0.81000000000000005, 0.27000000000000002, 0.18000000000000002, 0.090000000000000011, 1.6200000000000001, -0.36000000000000004, -0.36000000000000004, 1.1700000000000002, -0.90000000000000002, -0.36000000000000004, -0.090000000000000011, 0.63, 0.63, -0.36000000000000004, 0.27000000000000002, 0.27000000000000002, 0.54000000000000004, -0.090000000000000011, 0.27000000000000002, 0.36000000000000004, -0.18000000000000002, 0, -0.45000000000000001, 0.090000000000000011, -0.81000000000000005},
			wantErr: false,
		},
		{
			name:    "2022.5: Coeffs when degrees are not equal (13 > 8)",
			args:    args{date: 2022.5},
			want:    &[]float64{-29390.549999999999, -1432.4000000000001, 4587.75, -2527.0999999999999, 2964.5, -3067.0999999999999, 1671.75, -790.60000000000002, 1368.7, -2395.9499999999998, -67.099999999999994, 1243.95, 239.15000000000001, 495.70000000000005, -542.14999999999998, 900, 805.5, 281.64999999999998, 71.549999999999997, -142.15000000000001, -296.39999999999998, 208.69999999999999, 35.25, -362.19999999999999, -235.05000000000001, 364.44999999999999, 47.700000000000003, 186.30000000000001, 214.55000000000001, -140.19999999999999, -122.7, -147.94999999999999, 39.799999999999997, 15.75, 99.650000000000006, 64.75, 64.75, -19.100000000000001, 73.900000000000006, 21.100000000000001, -118.25, 49.549999999999997, -39.700000000000003, -62.5, 13.5, 8.9000000000000004, -62.450000000000003, 70.599999999999994, 80.349999999999994, -77.200000000000003, -50, -8.1999999999999993, -15.399999999999999, 58.25, 0.20000000000000018, 16.050000000000001, 23, 5.1500000000000004, -4.9500000000000002, -9.1999999999999993, -26.949999999999999, 11.800000000000001, -1.05, 23.699999999999999, 9.9499999999999993, 7.9000000000000004, -17.850000000000001, -13.800000000000001, 0.5, 12.300000000000001, -21.350000000000001, -10.449999999999999, 16.300000000000001, 14.15, 14.449999999999999, 2.6000000000000001, -16.75, -5.6500000000000004, 0.69999999999999996, 2.7999999999999998, 5, 8.4000000000000004, -23.399999999999999, 2.8999999999999999, 11, -1.5, 9.8000000000000007, -1.1000000000000001, -5.0999999999999996, -13.199999999999999, -6.2999999999999998, 1.1000000000000001, 7.7999999999999998, 8.8000000000000007, 0.40000000000000002, -9.3000000000000007, -1.3999999999999999, -11.9, 9.5999999999999996, -1.8999999999999999, -6.2000000000000002, 3.3999999999999999, -0.10000000000000001, -0.20000000000000001, 1.7, 3.6000000000000001, -0.90000000000000002, 4.7999999999999998, 0.69999999999999996, -8.5999999999999996, -0.90000000000000002, -0.10000000000000001, 1.8999999999999999, -4.2999999999999998, 1.3999999999999999, -3.3999999999999999, -2.3999999999999999, -0.10000000000000001, -3.7999999999999998, -8.8000000000000007, 3, -1.3999999999999999, 0, -2.5, 2.5, 2.2999999999999998, -0.59999999999999998, -0.90000000000000002, -0.40000000000000002, 0.29999999999999999, 0.59999999999999998, -0.69999999999999996, -0.20000000000000001, -0.10000000000000001, -1.7, 1.3999999999999999, -1.6000000000000001, -0.59999999999999998, -3, 0.20000000000000001, -2, 3.1000000000000001, -2.6000000000000001, -2, -0.10000000000000001, -1.2, 0.5, 0.5, 1.3, 1.3999999999999999, -1.2, -1.8, 0.69999999999999996, 0.10000000000000001, 0.29999999999999999, 0.80000000000000004, 0.5, -0.20000000000000001, -0.29999999999999999, 0.59999999999999998, -0.5, 0.20000000000000001, 0.10000000000000001, -0.90000000000000002, -1.1000000000000001, 0, -0.29999999999999999, 0.5, 0.10000000000000001, -0.90000000000000002, -0.90000000000000002, 0.5, 0.59999999999999998, 0.69999999999999996, 1.3999999999999999, -0.29999999999999999, -0.40000000000000002, 0.80000000000000004, -1.3, 0, -0.10000000000000001, 0.80000000000000004, 0.29999999999999999, 0, -0.10000000000000001, 0.40000000000000002, 0.5, 0.10000000000000001, 0.5, 0.5, -0.40000000000000002, -0.5, -0.40000000000000002, -0.40000000000000002, -0.59999999999999998},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := igrf.Coeffs(tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("IGRFcoeffs.Coeffs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for index, value := range *got {
				ref_value := math.Abs((*tt.want)[index])
				calc_err := math.Abs(100 * ((ref_value - math.Abs(value)) / ref_value))
				// allowable error percent
				max_error := 0.7
				if calc_err > max_error {
					t.Errorf("Calculated value = %v, reference %v, error %v is more than %v%%", value, ref_value, calc_err, max_error)
				}

			}
		})
	}
}
