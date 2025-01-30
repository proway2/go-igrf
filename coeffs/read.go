package coeffs

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
)

// max possible spherical harmonic degree
const N_MAX = 13

// number of lines (with data) in the coeffs file
const coeffs_lines = N_MAX * (N_MAX + 2)

// an interval between epochs in the coeffs file, e.g. (1950.0 - 1945.0 = 5)
const interval = 5

var (
	space_re   *regexp.Regexp = regexp.MustCompile(`\s+`)
	year_sv_re *regexp.Regexp = regexp.MustCompile(`\d{4}-\d{2}`)
)

type IGRFcoeffs struct {
	names  *[]string
	epochs *[]float64
	data   *map[string]*epochData
}

type epochData struct {
	nmax   int
	coeffs *[]float64
}

// Returns an initialized IGRF SHC data structure.
func NewCoeffsData() (*IGRFcoeffs, error) {
	igrf := IGRFcoeffs{data: &map[string]*epochData{}}
	if err := igrf.readCoeffs(); err != nil {
		return nil, err
	}
	return &igrf, nil
}

// Returns two sets of SH coeffs for the given `date`,
// as well as for `date` plus one year. Also returns the maximal spherical harmonic degree.
func (igrf *IGRFcoeffs) Coeffs(date float64) (*[]float64, *[]float64, int, error) {
	max_column := len(*igrf.epochs)
	min_epoch := (*igrf.epochs)[0]
	max_epoch := (*igrf.epochs)[max_column-1]
	if date < min_epoch || date > max_epoch {
		return nil, nil, 0, fmt.Errorf("date %v is out of range (%v, %v)", date, min_epoch, max_epoch)
	}
	// calculate coeffs for the requested date
	start, end := igrf.findEpochs(date)
	var nmax int
	coeffs_start, nmax := igrf.interpolateCoeffs(start, end, date)
	// in order to calculate yearly SV add 1 year to the date
	date = date + 1
	var coeffs_end *[]float64
	if date < max_epoch {
		coeffs_end, nmax = igrf.interpolateCoeffs(start, end, date)
	} else {
		coeffs_end = igrf.extrapolateCoeffs(start, end, date)
	}
	return coeffs_start, coeffs_end, nmax, nil
}

// Computes a set of SH coeffs and the maximal spherical harmonic degree
// for a given `date` and `start_epoch`, `end_epoch`.
//
// `date` must be less than the maximal possible epoch.
func (igrf *IGRFcoeffs) interpolateCoeffs(start_epoch, end_epoch string, date float64) (*[]float64, int) {
	factor, err := findDateFactor(start_epoch, end_epoch, date)
	if err != nil {
		log.Fatal("Epochs are incorrect!")
	}
	coeffs_start := (*igrf.data)[start_epoch].coeffs
	coeffs_end := (*igrf.data)[end_epoch].coeffs
	values := make([]float64, len(*coeffs_start))
	nmax1 := (*igrf.data)[start_epoch].nmax
	nmax2 := (*igrf.data)[end_epoch].nmax
	var k, l, nmax int
	var interp func(float64, float64, float64) float64
	if nmax1 == nmax2 {
		// before 2000.0
		k = nmax1 * (nmax1 + 2)
		l = -100
		nmax = nmax1
	} else {
		if nmax1 > nmax2 {
			// the last column has degree of 8
			// now it's anything after 2025.0
			k = nmax2 * (nmax2 + 2)
			l = nmax1 * (nmax1 + 2)
			interp = func(start, end, f float64) float64 {
				return start
			}
			nmax = nmax1
		} else {
			// between 1995.0 and 2000.0
			k = nmax1 * (nmax1 + 2)
			l = nmax2 * (nmax2 + 2)
			interp = func(_, end, f float64) float64 {
				return f * end
			}
			nmax = nmax2
		}
	}
	for i := 0; i < coeffs_lines; i++ {
		coeff_start := (*coeffs_start)[i]
		coeff_end := (*coeffs_end)[i]
		var value float64
		if i >= k && i < l {
			value = interp(coeff_start, coeff_end, factor)
		} else {
			value = coeff_start + factor*(coeff_end-coeff_start)
		}
		values[i] = value
	}
	return &values, nmax
}

// Computes a set of SH coeffs and the maximal spherical harmonic degree
// for a given `date` and `start_epoch`, `end_epoch`.
//
// `date` is beyond than the maximal possible epoch.
func (igrf *IGRFcoeffs) extrapolateCoeffs(start_epoch, end_epoch string, date float64) *[]float64 {
	dte1, _ := strconv.ParseFloat(start_epoch, 32)
	factor := date - dte1
	coeffs_start := (*igrf.data)[start_epoch].coeffs
	coeffs_end := (*igrf.data)[end_epoch].coeffs
	nmax1 := (*igrf.data)[start_epoch].nmax
	nmax2 := (*igrf.data)[end_epoch].nmax
	if nmax1 <= nmax2 {
		return nil // error here?
	}
	var k, l int
	k = nmax2 * (nmax2 + 2)
	l = nmax1 * (nmax1 + 2)
	values := make([]float64, len(*coeffs_start))
	for i := 0; i < coeffs_lines; i++ {
		coeff_start := (*coeffs_start)[i]
		coeff_end := (*coeffs_end)[i]
		var value float64
		if i >= k && i < l {
			value = coeff_start
		} else {
			sv := (coeff_end - coeff_start) / interval
			value = coeff_start + factor*sv
		}
		values[i] = value
	}
	return &values
}

// Calculates start and end epochs for a given date,
// default `interval` is used for borders and for multiplicity
func (igrf *IGRFcoeffs) findEpochs(date float64) (string, string) {
	max_column := len(*igrf.epochs)
	min_epoch := (*igrf.epochs)[0]
	max_epoch := (*igrf.epochs)[max_column-1]
	var start_epoch, end_epoch string
	if date >= max_epoch {
		start_epoch = epoch2string(max_epoch)
		end_epoch = epoch2string(max_epoch + interval)
		return start_epoch, end_epoch
	}
	col1 := min_epoch + float64(int((date-min_epoch)/interval))*interval
	start_epoch = epoch2string(col1)
	col2 := col1 + interval
	end_epoch = epoch2string(col2)
	return start_epoch, end_epoch
}

// The main function that populates the existing `IGRFcoeffs` structure.
func (igrf *IGRFcoeffs) readCoeffs() error {
	line_provider := coeffsLineProvider()

	var err error
	igrf.names, igrf.epochs, err = getEpochs(line_provider)
	if err != nil {
		return err
	}
	// initializing the map
	for _, epoch := range *igrf.epochs {
		local_arr := make([]float64, coeffs_lines)
		(*igrf.data)[epoch2string(epoch)] = &epochData{coeffs: &local_arr}
	}
	err = igrf.getCoeffsForEpochs(line_provider)
	if err != nil {
		return err
	}
	for epoch := range *igrf.data {
		nmax, err := nMaxForEpoch(epoch)
		if err != nil {
			return err
		}
		(*igrf.data)[epoch].nmax = nmax
	}
	return nil
}

// Finds lines with epochs, parses these lines and returns arrays of epochs: names and floats.
func getEpochs(reader <-chan string) (*[]string, *[]float64, error) {
	cs_re := regexp.MustCompile(`^c/s.*`)
	for line := range reader {
		if !cs_re.Match([]byte(line)) {
			continue
		}
		line2 := <-reader
		names, epochs := parseHeader(line, line2)
		return &names, &epochs, nil
	}
	return nil, nil, errors.New("unable to get epochs")
}

// Parses the header of the coeffs. Usually it's the first two non-comment lines.
func parseHeader(line1, line2 string) ([]string, []float64) {
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
				log.Fatal(err.Error())
			}
			epoch := 2000.0 + decades
			epochs[index] = epoch
		}
	}
	return names, epochs[shift:]
}

// Parses lines with coefficients and populates the `IGRFcoeffs` structure.
func (igrf *IGRFcoeffs) getCoeffsForEpochs(provider <-chan string) error {
	var i int = 0
	for line := range provider {
		line_data := space_re.Split(line, -1)
		line_coeffs, err := parseArrayToFloat(line_data[3:])
		if err != nil {
			return errors.New("unable to parse coeffs")
		}
		igrf.loadCoeffs(i, line_coeffs)
		i++
	}
	return nil
}

// Populates the corresponding fields inside the `IGRFcoeffs` structure with actual coeffs.
func (igrf *IGRFcoeffs) loadCoeffs(line_num int, line_coeffs *[]float64) {
	for index, coeff := range *line_coeffs {
		epoch := (*igrf.epochs)[index]
		epoch_str := epoch2string(epoch)
		(*(*igrf.data)[epoch_str].coeffs)[line_num] = coeff
	}
}

// Returns max spherical harmonic degree for a certain epoch.
func nMaxForEpoch(epoch string) (int, error) {
	// this is hardcoded
	var nmax int
	epoch_f, err := strconv.ParseFloat(epoch, 32)
	if err != nil {
		return 0, err
	}
	if epoch_f < 2000.0 {
		nmax = 10
	} else if epoch_f > 2025.0 {
		nmax = 8
	} else {
		nmax = 13
	}
	return nmax, nil
}
