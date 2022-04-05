package coeffs

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const file_name string = "../coeffs/igrf13coeffs.txt"
const coeffs_lines = 195
const years_before_sv = 5

var space_re *regexp.Regexp = regexp.MustCompile(`\s+`)
var year_re *regexp.Regexp = regexp.MustCompile(`\d{4}`)
var year_sv_re *regexp.Regexp = regexp.MustCompile(`\d{4}-\d{2}`)

type IGRFcoeffs struct {
	names  *[]string
	epochs *[]float64
	coeffs *[]lineData
}

type lineData struct {
	g_h    bool // the model contains either "g" or "h" g_h == true if there is "g" (false if "h")
	deg_n  int
	ord_m  int
	coeffs *[]float64
}

func NewCoeffsData() *IGRFcoeffs {
	igrf := IGRFcoeffs{}
	igrf.readCoeffs()
	return &igrf
}

func Coeffs(date float64) (*[]float64, error) {
	return &[]float64{}, nil
}

func (igrf *IGRFcoeffs) readCoeffs() {
	f, err := os.Open(file_name)
	defer f.Close()
	if err != nil {
		log.Fatal("IGRF Coeffs file not found.")
	}
	scanner := bufio.NewScanner(f)
	names, epochs, err := getEpochs(scanner)
	if err != nil {
		log.Fatal(err.Error())
	}
	igrf.names = names
	igrf.epochs = epochs
	coeffs, err := getCoeffs(scanner)
	if err != nil {
		log.Fatal(err.Error())
	}
	igrf.coeffs = coeffs
}

func getEpochs(scanner *bufio.Scanner) (*[]string, *[]float64, error) {
	cs_re := regexp.MustCompile(`^c/s.*`)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Trim(line, " ")
		if line[0] == 35 { // #
			continue
		}
		if cs_re.Match([]byte(line)) {
			scanner.Scan()
			line2 := scanner.Text()
			names, epochs := parseHeader1(line, line2)
			return &names, &epochs, nil
		}
	}
	return nil, nil, errors.New("Unable to get epochs.")
}

func parseHeader1(line1, line2 string) ([]string, []float64) {
	line1_data := space_re.Split(line1, -1)
	line2_data := space_re.Split(line2, -1)

	if len(line1_data) != len(line2_data) {
		log.Fatal("Coeffs header is incorrect.")
	}
	names := make([]string, len(line1_data))
	epochs := make([]float64, len(line1_data))
	var shift int
	for index := range line1_data {
		raw_epoch := line2_data[index]
		name := fmt.Sprintf("%v %v", line1_data[index], raw_epoch)
		names[index] = name
		if epoch, err := strconv.ParseFloat(raw_epoch, 32); err == nil {
			epochs[index] = epoch
			if shift == 0 {
				shift = index
			}
		}
		// this is the last column
		if year_sv_re.Match([]byte(raw_epoch)) {
			last_digits := raw_epoch[5:]
			decades, err := strconv.ParseFloat(last_digits, 32)
			if err != nil {
				log.Fatal("Unknown year at SV column.")
			}
			epoch := 2000.0 + decades
			epochs[index] = epoch
		}
	}
	return names, epochs[shift:]
}

func getCoeffs(scanner *bufio.Scanner) (*[]lineData, error) {
	coeffs := make([]lineData, coeffs_lines)
	for i := 0; scanner.Scan(); i++ {
		data := lineData{}
		line := scanner.Text()
		line = strings.Trim(line, " ")
		line_data := space_re.Split(line, -1)
		if line_data[0] == "g" {
			data.g_h = true
		} else {
			data.g_h = false
		}
		deg, _ := strconv.ParseInt(line_data[1], 10, 0)
		data.deg_n = int(deg)
		ord, _ := strconv.ParseInt(line_data[2], 10, 0)
		data.ord_m = int(ord)
		line_coeffs, err := parseArrayToFloat(line_data[3:])
		if err != nil {
			return nil, errors.New("Unable to parse coeffs.")
		}
		data.coeffs = line_coeffs
		coeffs[i] = data
	}
	return &coeffs, nil
}

func parseArrayToFloat(raw_data []string) (*[]float64, error) {
	data := make([]float64, len(raw_data))
	for index, token := range raw_data {
		real_data, err := strconv.ParseFloat(token, 32)
		if err != nil {
			return nil, errors.New("Unable to parse coeffs.")
		}
		if index == len(raw_data)-1 {
			real_data = data[index-1] + real_data*years_before_sv
			fmt.Println(real_data)
		}
		data[index] = real_data
	}
	return &data, nil
}
