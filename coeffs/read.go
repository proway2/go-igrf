package coeffs

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
)

// number of lines (with data) in the coeffs file
const coeffs_lines = 195

// an interval between epochs in the coeffs file, e.g. (1950.0 - 1945.0 = 5)
const interval = 5

var (
	space_re   *regexp.Regexp = regexp.MustCompile(`\s+`)
	year_sv_re *regexp.Regexp = regexp.MustCompile(`\d{4}-\d{2}`)
)

type IGRFcoeffs struct {
	names  *[]string
	epochs *[]float64
	// is `lines` needed at all?
	lines  *[]lineData
	coeffs *map[string]*[]float64
}

type lineData struct {
	g_h    bool // Gauss coefficients the model contains either "g" or "h" g_h == true if there is "g" (false if "h")
	deg_n  int  // spherical harmonic degree
	ord_m  int  // order
	coeffs *[]float64
}

func NewCoeffsData() (*IGRFcoeffs, error) {
	igrf := IGRFcoeffs{coeffs: &map[string]*[]float64{}}
	if err := igrf.readCoeffs(); err != nil {
		return nil, err
	}
	return &igrf, nil
}

func (igrf *IGRFcoeffs) Coeffs(date float64) (*[]float64, error) {
	start, end, err := igrf.findEpochs(date)
	if err != nil {
		return nil, err
	}
	coeffs := igrf.interpolateCoeffs(start, end, date)
	if err != nil {
		return nil, err
	}
	return coeffs, nil
}

func (igrf *IGRFcoeffs) interpolateCoeffs(start_epoch, end_epoch string, date float64) *[]float64 {
	fraction := findDateFraction(start_epoch, end_epoch, date)
	coeffs_start := (*igrf.coeffs)[start_epoch]
	coeffs_end := (*igrf.coeffs)[end_epoch]
	values := make([]float64, len(*coeffs_start))
	for index, coeff_start := range *coeffs_start {
		coeff_end := (*coeffs_end)[index]
		value := ((coeff_end - coeff_start) * fraction) + coeff_start
		values[index] = value
	}
	return &values
}

func (igrf *IGRFcoeffs) findEpochs(date float64) (string, string, error) {
	max_column := len(*igrf.epochs)
	min_epoch := (*igrf.epochs)[0]
	max_epoch := (*igrf.epochs)[max_column-1]
	if date < min_epoch || date >= max_epoch {
		return "", "", errors.New(fmt.Sprintf("Date %v is out of range (%v, %v).", date, min_epoch, max_epoch))
	}
	col1 := min_epoch + float64(int((date-min_epoch)/interval))*interval
	start_epoch := epoch2string(col1)
	col2 := col1 + interval
	end_epoch := epoch2string(col2)
	return start_epoch, end_epoch, nil
}

// func (igrf *IGRFcoeffs) findColumns(date float64) (int, int, error) {
// 	max_column := len(*igrf.epochs)
// 	min_epoch := (*igrf.epochs)[0]
// 	max_epoch := (*igrf.epochs)[max_column-1]
// 	if date < min_epoch || date > max_epoch {
// 		return -1, -1, errors.New(fmt.Sprintf("Date %v is out of range (%v, %v).", date, min_epoch, max_epoch))
// 	}
// 	col1 := int((date - min_epoch) / interval)
// 	col2 := col1 + 1
// 	return col1, col2, nil
// }

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
		(*igrf.coeffs)[epoch2string(epoch)] = &local_arr
	}
	// TODO: decide whether this igrf.lines is needed at all
	// igrf.lines, err = getCoeffs(line_provider)
	// if err != nil {
	// 	return err
	// }
	igrf.getCoeffsForEpochs(line_provider)
	return nil
}

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
	return nil, nil, errors.New("Unable to get epochs.")
}

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

// func getCoeffs(reader <-chan string) (*[]lineData, error) {
// 	coeffs := make([]lineData, coeffs_lines)
// 	var i int = 0
// 	for line := range reader {
// 		data := lineData{}
// 		line_data := space_re.Split(line, -1)
// 		if line_data[0] == "g" {
// 			data.g_h = true
// 		} else {
// 			data.g_h = false
// 		}
// 		deg, _ := strconv.ParseInt(line_data[1], 10, 0)
// 		data.deg_n = int(deg)
// 		ord, _ := strconv.ParseInt(line_data[2], 10, 0)
// 		data.ord_m = int(ord)
// 		line_coeffs, err := parseArrayToFloat(line_data[3:])
// 		if err != nil {
// 			return nil, errors.New("Unable to parse coeffs.")
// 		}
// 		data.coeffs = line_coeffs
// 		coeffs[i] = data
// 		i++
// 	}
// 	return &coeffs, nil
// }

func parseArrayToFloat(raw_data []string) (*[]float64, error) {
	data := make([]float64, len(raw_data))
	for index, token := range raw_data {
		real_data, err := strconv.ParseFloat(token, 32)
		if err != nil {
			return nil, errors.New("Unable to parse coeffs.")
		}
		if index == len(raw_data)-1 {
			// real value calculated for the SV column
			real_data = data[index-1] + real_data*interval
		}
		data[index] = real_data
	}
	return &data, nil
}

func (igrf *IGRFcoeffs) getCoeffsForEpochs(provider <-chan string) (*[]float64, error) {
	var i int = 0
	for line := range provider {
		data := lineData{}
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
		igrf.loadCoeffs(i, line_coeffs)
		i++
	}
	return &[]float64{}, nil
}

func (igrf *IGRFcoeffs) loadCoeffs(line_num int, line_coeffs *[]float64) {
	for index, coeff := range *line_coeffs {
		epoch := (*igrf.epochs)[index]
		epoch_str := epoch2string(epoch)
		(*(*igrf.coeffs)[epoch_str])[line_num] = coeff
	}
}
