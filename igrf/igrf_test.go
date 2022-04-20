package igrf

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
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
const max_allowed_error = 0.2 // %

func TestIGRFDataCases(t *testing.T) {
	tests := getTestData()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IGRF(tt.args.lat, tt.args.lon, tt.args.alt, tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("IGRF() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			compare_x := compareFloats(float64(got.NorthComponent), float64(tt.want.NorthComponent), max_allowed_error)
			if !compare_x {
				t.Errorf("IGRF() NorthComponent = %v, want %v", got.NorthComponent, tt.want.NorthComponent)
			}
			// Fortran results are rounded!!!
			new_east_got := math.Round(float64(got.EastComponent))
			new_east_want := math.Round(float64(tt.want.EastComponent))
			compare_y := compareFloats(new_east_got, new_east_want, max_allowed_error)
			if !compare_y {
				t.Errorf("IGRF() EastComponent = %v, want %v", got.EastComponent, tt.want.EastComponent)
			}
			compare_z := compareFloats(float64(got.VerticalComponent), float64(tt.want.VerticalComponent), max_allowed_error)
			if !compare_z {
				t.Errorf("IGRF() VerticalComponent = %v, want %v", got.VerticalComponent, tt.want.VerticalComponent)
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("IGRF() = %v, want %v", got, tt.want)
			// }
		})
	}
}

func compareFloats(check, base, allowable_error float64) bool {
	value1 := math.Abs(check)
	value2 := math.Abs(base)
	calc_err := math.Abs(100 * ((value1 - value2) / value2))
	// allowable error percent
	if calc_err > allowable_error {
		return false
	}
	return true
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
	tests := make([]testsData, 625)
	var index int
	for _, file := range test_data_files {
		f, err := os.Open(file)
		defer f.Close()
		check(err)
		current_file_tests := produceTestsDataFromFile(f)
		for _, test := range current_file_tests {
			tests[index] = test
			index++
		}
	}
	return tests
}

func produceTestsDataFromFile(file_descr *os.File) []testsData {
	scanner := bufio.NewScanner(file_descr)
	split_regex := regexp.MustCompile(`\s+`)
	var lat, lon, alt float64
	tests := make([]testsData, 125)
	for num, i := 0, 0; scanner.Scan(); num++ {
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
		tests[i] = current_test
		i++
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

func toFloat32(str string) float32 {
	val, err := strconv.ParseFloat(str, 32)
	check(err)
	return float32(val)
}

func getIGRFresults(line []string) IGRFresults {
	return IGRFresults{
		Declination:         toFloat32(line[1]),
		DeclinationSV:       toFloat32(line[2]),
		Inclination:         toFloat32(line[3]),
		InclinationSV:       toFloat32(line[4]),
		HorizontalIntensity: toFloat32(line[5]),
		HorizontalSV:        toFloat32(line[6]),
		NorthComponent:      toFloat32(line[7]),
		NorthSV:             toFloat32(line[8]),
		EastComponent:       toFloat32(line[9]),
		EastSV:              toFloat32(line[10]),
		VerticalComponent:   toFloat32(line[11]),
		VerticalSV:          toFloat32(line[12]),
		TotalIntensity:      toFloat32(line[13]),
		TotalSV:             toFloat32(line[14]),
	}
}
