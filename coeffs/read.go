package coeffs

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const file_name string = "../coeffs/igrf13coeffs.txt"

var space_re *regexp.Regexp = regexp.MustCompile(`\s+`)
var year_re *regexp.Regexp = regexp.MustCompile(`\d{4}`)
var year_sv_re *regexp.Regexp = regexp.MustCompile(`\d{4}-\d{2}`)

type IGRFcoeffs struct {
	time   []float64
	coeffs [][]float64
}

func LoadCoeffsFile() {
	f, err := os.Open(file_name)
	defer f.Close()
	if err != nil {
		panic("IGRF Coeffs file not found.")
	}
	cs_re := regexp.MustCompile(`^c/s.*`)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Trim(line, " ")
		if line[0] == 35 { // #
			continue
		}
		if cs_re.Match([]byte(line)) {
			scanner.Scan()
			line2 := scanner.Text()
			a, b := parseHeader1(line, line2)
			fmt.Println(a, b)
		}
	}
}

func getEpochs(line string) []float64 {
	return []float64{}
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
