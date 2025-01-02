package igrf

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/proway2/go-igrf/coeffs"
)

type args struct {
	lat  float64
	lon  float64
	alt  float64
	date float64
}

type testsData struct {
	name    string
	args    args
	want    IGRFresults
	wantErr bool
}

const dir_path string = "../testdata"

// allowed errors
const (
	max_rel_tol           = 0.005 // max relative tolerance
	near_pole_max_rel_tol = 0.03  // max relative tolerance is increased for near pole values
	d_i_abs_tol           = 0.005 // D and I are tested with much higher accuracy
)

const near_pole_tolerance = 0.001

func TestNew(t *testing.T) {
	shc, _ := coeffs.NewCoeffsData()
	tests := []struct {
		name string
		want *IGRFdata
	}{
		{
			name: "Creating a IGRF data structure",
			want: &IGRFdata{shc: shc},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

// SV fields are just integers in FORTRAN, so there might be situations where:
// calculated value 16.47, reference 17
// calculated value 2.52, reference 3
// ...
// const max_sv_error = 50 // %

func TestIGRFdata_IGRF(t *testing.T) {
	tests := getTestData()
	igrf_data := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := igrf_data.IGRF(tt.args.lat, tt.args.lon, tt.args.alt, tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("IGRFdata.IGRF() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			allowed_relative_error := getMaxAllowedRelativeError(tt.args.lat)
			// MAIN VALUES

			// Declination
			compare := isClose(got.Declination, tt.want.Declination, allowed_relative_error, d_i_abs_tol)
			if !compare {
				t.Errorf("IGRF() Declination = %v, want %v", got.Declination, tt.want.Declination)
			}
			// Inclination
			compare = isClose(got.Inclination, tt.want.Inclination, allowed_relative_error, d_i_abs_tol)
			if !compare {
				t.Errorf("IGRF() Inclination = %v, want %v", got.Inclination, tt.want.Inclination)
			}
			// Horizontal intensity
			allowed_abs_error := getMaxAllowedAbsoluteTolerance(tt.want.HorizontalIntensity)
			compare = isClose(got.HorizontalIntensity, tt.want.HorizontalIntensity, allowed_relative_error, allowed_abs_error)
			if !compare {
				t.Errorf("IGRF() Horizontal intensity = %v, want %v", got.HorizontalIntensity, tt.want.HorizontalIntensity)
			}
			// North component
			allowed_abs_error = getMaxAllowedAbsoluteTolerance(tt.want.NorthComponent)
			compare = isClose(got.NorthComponent, tt.want.NorthComponent, allowed_relative_error, allowed_abs_error)
			if !compare {
				t.Errorf("IGRF() NorthComponent = %v, want %v", got.NorthComponent, tt.want.NorthComponent)
			}
			// East Component - FORTRAN results are rounded!!!
			allowed_abs_error = getMaxAllowedAbsoluteTolerance(tt.want.EastComponent)
			compare = isClose(got.EastComponent, tt.want.EastComponent, allowed_relative_error, allowed_abs_error)
			if !compare {
				t.Errorf("IGRF() EastComponent = %v, want %v", got.EastComponent, tt.want.EastComponent)
			}
			// Vertical component
			allowed_abs_error = getMaxAllowedAbsoluteTolerance(tt.want.VerticalComponent)
			compare = isClose(got.VerticalComponent, tt.want.VerticalComponent, allowed_relative_error, allowed_abs_error)
			if !compare {
				t.Errorf("IGRF() VerticalComponent = %v, want %v", got.VerticalComponent, tt.want.VerticalComponent)
			}
			// Total intensity
			allowed_abs_error = getMaxAllowedAbsoluteTolerance(tt.want.TotalIntensity)
			compare = isClose(got.TotalIntensity, tt.want.TotalIntensity, allowed_relative_error, allowed_abs_error)
			if !compare {
				t.Errorf("IGRF() Total = %v, want %v", got.TotalIntensity, tt.want.TotalIntensity)
			}

			// SECULAR VALUES

			// // Declination SV
			// compare = compareFloats(math.Round(float64(got.DeclinationSV)), float64(tt.want.DeclinationSV), max_sv_error, max_absolute_error)
			// if !compare {
			// 	t.Errorf("IGRF() DeclinationSV = %v, want %v", got.DeclinationSV, tt.want.DeclinationSV)
			// }
			// // Inclination SV
			// compare = compareFloats(math.Round(float64(got.InclinationSV)), float64(tt.want.InclinationSV), max_sv_error, max_absolute_error)
			// if !compare {
			// 	t.Errorf("IGRF() InclinationSV = %v, want %v", got.InclinationSV, tt.want.InclinationSV)
			// }
			// // Horizontal SV
			// compare = compareFloats(math.Round(float64(got.HorizontalSV)), float64(tt.want.HorizontalSV), max_sv_error, max_absolute_error)
			// if !compare {
			// 	t.Errorf("IGRF() HorizontalSV = %v, want %v", got.HorizontalSV, tt.want.HorizontalSV)
			// }
			// // North SV
			// compare = compareFloats(math.Round(float64(got.NorthSV)), float64(tt.want.NorthSV), max_sv_error, max_absolute_error)
			// if !compare {
			// 	t.Errorf("IGRF() NorthSV = %v, want %v", got.NorthSV, tt.want.NorthSV)
			// }
			// // EastSV
			// compare = compareFloats(math.Round(float64(got.EastSV)), float64(tt.want.EastSV), max_sv_error, max_absolute_error)
			// if !compare {
			// 	t.Errorf("IGRF() EastSV = %v, want %v", got.EastSV, tt.want.EastSV)
			// }
			// // VerticalSV
			// compare = compareFloats(math.Round(float64(got.VerticalSV)), float64(tt.want.VerticalSV), max_sv_error, max_absolute_error)
			// if !compare {
			// 	t.Errorf("IGRF() VerticalSV = %v, want %v", got.VerticalSV, tt.want.VerticalSV)
			// }
			// // TotalSV
			// compare = compareFloats(math.Round(float64(got.TotalSV)), float64(tt.want.TotalSV), max_sv_error, max_absolute_error)
			// if !compare {
			// 	t.Errorf("IGRF() TotalSV = %v, want %v", got.TotalSV, tt.want.TotalSV)
			// }
		})
	}
}

func isNearPole(lat float64) bool {
	return 90.0-math.Abs(lat) <= near_pole_tolerance
}

func getMaxAllowedRelativeError(lat float64) float64 {
	if isNearPole(lat) {
		return near_pole_max_rel_tol
	}
	return max_rel_tol
}

func getMaxAllowedAbsoluteTolerance(value float64) float64 {
	abs_tol := 0.15
	// For some values (H, X, Y, Z) that are really small it's hard to calculate the accuracy due to the fact that
	// results from `FORTRAN` are rounded. That's why for some small values (close to the magnetic pole) absolute and relative accuracies
	// are increased.
	if value < 70 {
		abs_tol *= 4
		return abs_tol
	}
	return abs_tol
}

func isClose(a, b, rel_tol, abs_tol float64) bool {
	abs_diff := math.Abs(a - b)
	a_abs := math.Abs(a)
	b_abs := math.Abs(b)
	lhs := math.Max(rel_tol*math.Max(a_abs, b_abs), abs_tol)
	return abs_diff <= lhs
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func discoverTestData() []string {
	files, err := filepath.Glob(dir_path + "/set*")
	check(err)
	return files
}

func getTestData() []testsData {
	test_data_files := discoverTestData()
	tests := make([]testsData, 0)
	for _, file := range test_data_files {
		f, err := os.Open(file)
		check(err)
		defer f.Close()
		current_file_tests := produceTestsDataFromFile(f)
		tests = append(tests, current_file_tests...)
	}
	return tests
}

func produceTestsDataFromFile(file_descr *os.File) []testsData {
	scanner := bufio.NewScanner(file_descr)
	split_regex := regexp.MustCompile(`\s+`)
	var lat, lon, alt float64
	var tests []testsData
	for num := 0; scanner.Scan(); num++ {
		line := scanner.Text()
		if num == 1 {
			// this is just a column names
			continue
		}
		line = strings.Trim(line, " ")
		line_data := split_regex.Split(line, -1)
		if num == 0 {
			lat, lon, alt = getArgs(line_data)
			continue
		}
		date := getDate(line_data)
		igrf_res := getIGRFresults(line_data)
		// num > 0
		current_test := testsData{
			name:    fmt.Sprintf("Lat:%v Lon:%v Alt:%v Date:%v", lat, lon, alt, date),
			args:    args{lat: lat, lon: lon, alt: alt, date: date},
			want:    igrf_res,
			wantErr: false,
		}
		tests = append(tests, current_test)
	}
	return tests
}

func getArgs(line []string) (float64, float64, float64) {
	lat := toFloat64(line[1])
	lon := toFloat64(line[4])
	alt := toFloat64(line[5])
	return lat, lon, alt
}

func getDate(line []string) float64 {
	date := toFloat64(line[0])
	return date
}

func toFloat64(str string) float64 {
	val, err := strconv.ParseFloat(str, 64)
	check(err)
	return float64(val)
}

func getIGRFresults(line []string) IGRFresults {
	return IGRFresults{
		Declination:         toFloat64(line[1]),
		DeclinationSV:       toFloat64(line[2]),
		Inclination:         toFloat64(line[3]),
		InclinationSV:       toFloat64(line[4]),
		HorizontalIntensity: toFloat64(line[5]),
		HorizontalSV:        toFloat64(line[6]),
		NorthComponent:      toFloat64(line[7]),
		NorthSV:             toFloat64(line[8]),
		EastComponent:       toFloat64(line[9]),
		EastSV:              toFloat64(line[10]),
		VerticalComponent:   toFloat64(line[11]),
		VerticalSV:          toFloat64(line[12]),
		TotalIntensity:      toFloat64(line[13]),
		TotalSV:             toFloat64(line[14]),
	}
}
