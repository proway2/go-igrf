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
	lat  float32
	lon  float32
	alt  float32
	date float32
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
	tests := make([]testsData, 500)
	for _, file := range test_data_files {
		f, err := os.Open(file)
		defer f.Close()
		check(err)
		current_file_tests := produceTestsDataFromFile(f)
		tests = append(tests, current_file_tests...)
	}
	return tests
}

func produceTestsDataFromFile(file_descr *os.File) []testsData {
	scanner := bufio.NewScanner(file_descr)
	split_regex := regexp.MustCompile(`\s+`)
	var lat, lon, alt float32
	tests := make([]testsData, 110)
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

func getArgs(line []string) (float32, float32, float32) {
	lat := toFloat32(line[1])
	lon := toFloat32(line[4])
	alt := toFloat32(line[5])
	return lat, lon, alt
}

func getDate(line []string) float32 {
	date := toFloat32(line[0])
	return date
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
