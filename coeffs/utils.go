package coeffs

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var comment_line *regexp.Regexp = regexp.MustCompile(`^\s*#.*`)

// Calculates the number of seconds per year, respects leap years.
func secsInYear(year int) int {
	var days_per_year int = 365
	if isLeapYear(year) {
		days_per_year = 366
	}
	secs_per_day := 3600 * 24
	return days_per_year * secs_per_day
}

// Returns whether the given year is leap or not.
func isLeapYear(year int) bool {
	isDivisibleBy4 := year%4 == 0
	isDivisibleBy100 := year%100 == 0
	isDivisibleBy400 := year%400 == 0
	return isDivisibleBy400 || (isDivisibleBy4 && !isDivisibleBy100)
}

// Finds the factor for a given `date` between two epochs.
// In the first approximation the factor is calculated like:
//
// factor = (date - start_epoch) / (end_epoch - start_epoch)
//
// This is the coarse approach and the actual factor is calculated with respect to the leap years,
// unless `date` is beyond the `end_epoch`. In this case the above formula is used.
//
// If `end_epoch` is less or equal to `start_epoch` - 0 is returned, no negative values returned.
//
// In case of no correct epochs are provided, error is returned.
func findDateFactor(start_epoch, end_epoch string, date float64) (float64, error) {
	parser := errParser{}
	dte1 := parser.parseFloat(start_epoch)
	dte2 := parser.parseFloat(end_epoch)
	if parser.err != nil {
		return -999, fmt.Errorf("Epoch(s) cannot be parsed, start:%v, end:%v", start_epoch, end_epoch)
	}
	if dte2 <= dte1 {
		return 0, nil
	}
	if date > dte2 {
		return (date - dte1) / (dte2 - dte1), nil
	}
	loc_interval := int(dte2) - int(dte1)
	var total_secs, fraction_secs float64
	for i := 0; i < loc_interval; i++ {
		year := int(dte1) + i
		secs_in_year := secsInYear(year)
		if year == int(date) {
			fraction_coeff := date - float64(int(date))
			fraction_secs = total_secs + fraction_coeff*float64(secs_in_year)
		}
		total_secs += float64(secs_in_year)
	}
	factor := fraction_secs / total_secs
	return factor, nil
}

// Reads lines from raw coeffs data and writes a srting into a channel, drops comments.
func coeffsLineProvider() <-chan string {
	ch := make(chan string)
	coeffs_reader := strings.NewReader(igrf13coeffs)
	scanner := bufio.NewScanner(coeffs_reader)
	go func() {
		defer close(ch)
		for scanner.Scan() {
			line := scanner.Text()
			if comment_line.Match([]byte(line)) {
				continue
			}
			line = strings.Trim(line, " ")
			ch <- line
		}
	}()
	return ch
}

// epoch2string - converts `epoch` of type `float64` into string.
func epoch2string(epoch float64) string {
	return fmt.Sprintf("%.1f", epoch)
}

// Parses an array of strings into an array of floats.
func parseArrayToFloat(raw_data []string) (*[]float64, error) {
	data := make([]float64, len(raw_data))
	for index, token := range raw_data {
		real_data, err := strconv.ParseFloat(token, 32)
		if err != nil {
			return nil, errors.New("unable to parse coeffs")
		}
		if index == len(raw_data)-1 {
			// real value calculated for the SV column
			real_data = data[index-1] + real_data*interval
		}
		data[index] = real_data
	}
	return &data, nil
}

type errParser struct {
	err error
}

func (p *errParser) parseFloat(v string) float64 {
	if p.err != nil {
		return 0
	}
	var value float64
	value, p.err = strconv.ParseFloat(v, 64)
	return value
}
