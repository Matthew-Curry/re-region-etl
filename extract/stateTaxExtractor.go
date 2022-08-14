/* Logic to extract and store state tax data into an in memory data structure for use in ETL */

package extract

import (
	"regexp"
	"strconv"
	"strings"

	sourcefileutils "github.com/Matthew-Curry/re-region-etl/sourceFileUtils"
)

const STATE_TAX_FILE = "data/State-Individual-Income-Tax-Rates-and-Brackets-for-2022-v.xlsx"

// helper method to build data structures for state tax brackets and exemptions
func GetStateTaxData(censusData [][]string, nullString string) ([][]string, [][]string, error) {

	// build hashmap of lower state to id
	mp := make(map[string]string)
	for _, row := range censusData {
		state := strings.ToLower(row[9])
		// add mapping for state if it does not already exist
		if _, ok := mp[state]; !ok {
			mp[state] = row[7]
		}

	}

	// read in the state individual file
	stateTaxData, err := sourcefileutils.OpenExcelSheet(STATE_TAX_FILE, "2022")
	if err != nil {
		return nil, nil, err
	}

	// parse data structures
	stateRates := [][]string{}
	stateExcemptions := [][]string{}
	stateId := nullString
	for _, row := range stateTaxData {
		// rows of length 12 are initial row for state, contain exemption
		if len(row) == 12 {
			// update state id, exemptions
			if newStateId, ok := mp[strings.TrimSpace(strings.ToLower(row[0]))]; ok {

				stateId = newStateId
				stateExcemptions = append(stateExcemptions, []string{stateId, strings.TrimSpace(row[0]),
					processDollarValue(row[7], nullString),
					processDollarValue(row[8], nullString),
					processDollarValue(row[9], nullString),
					processDollarValue(row[10], nullString),
					processDollarValue(row[11], nullString)})
			}
		}

		// rows of length 12 also contain the first bracket information, and rows of length 8 are successive bracket information
		if (len(row) == 12 || len(row) == 7) && stateId != nullString {
			stateRates = append(stateRates, []string{stateId,
				processRate(row[1], nullString),
				processDollarValue(row[3], nullString),
				processRate(row[4], nullString),
				processDollarValue(row[6], nullString)})

		}

	}

	return stateRates, stateExcemptions, nil

}

// helper method with logic to process an exemption
func processDollarValue(ex, nullString string) string {
	// get rid of everything not a number
	reg, _ := regexp.Compile("[^0-9]+")
	rep := reg.ReplaceAllString(ex, "")

	if rep == "" {
		return nullString
	} else {
		return rep
	}
}

// process a tax rate
func processRate(r, nullString string) string {
	// trim spaces
	r = strings.TrimSpace(r)
	// trim %
	r = strings.Trim(r, "%")
	// If non number, return null string
	_, err := strconv.ParseFloat(r, 64)
	if err != nil {
		return nullString
		// if r is empty return null
	} else if r == "" {
		return nullString
	} else {
		return r
	}
}
