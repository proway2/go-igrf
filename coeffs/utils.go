package coeffs

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
