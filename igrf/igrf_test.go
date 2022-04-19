package igrf

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
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

func TestIGRFDataCases(t *testing.T) {
	tests := getTestData()
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
