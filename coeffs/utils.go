package coeffs

import (
	"bufio"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func secsInYear(year int) int {
	var days_per_year int = 365
	if isLeapYear(year) {
		days_per_year = 366
	}
	secs_per_day := 3600 * 24
	return days_per_year * secs_per_day
}

func isLeapYear(year int) bool {
	isDivisibleBy4 := year%4 == 0
	isDivisibleBy100 := year%100 == 0
	isDivisibleBy400 := year%400 == 0
	return isDivisibleBy400 || (isDivisibleBy4 && !isDivisibleBy100)
}

func findDateFraction(start_epoch, end_epoch string, date float64) float64 {
	start_year, _ := strconv.ParseFloat(start_epoch, 32)
	end_year, _ := strconv.ParseFloat(end_epoch, 32)
	if end_year <= start_year {
		log.Fatalf("End epoch %v is less than start epoch %v", end_epoch, start_epoch)
	}
	loc_interval := int(end_year) - int(start_year)
	var total_secs, fraction_secs float64
	for i := 0; i < loc_interval; i++ {
		year := int(start_year) + i
		secs_in_year := secsInYear(year)
		if year == int(date) {
			fraction_coeff := date - float64(int(date))
			fraction_secs = total_secs + fraction_coeff*float64(secs_in_year)
		}
		total_secs += float64(secs_in_year)
	}
	fraction := fraction_secs / total_secs
	return fraction
}

func coeffsLineProvider() <-chan string {
	ch := make(chan string)
	coeffs_reader := strings.NewReader(igrf13coeffs)
	scanner := bufio.NewScanner(coeffs_reader)
	go func() {
		defer close(ch)
		for scanner.Scan() {
			line := scanner.Text()
			line = strings.Trim(line, " ")
			ch <- line
		}
	}()
	return ch
}

func epoch2string(epoch float64) string {
	return fmt.Sprintf("%.1f", epoch)
}
