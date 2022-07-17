/* Logic to extract and extract local tax data into an in memory data structure for use in ETL */

package extract

import (
	"math"
	"strings"
	"sync"
	"sync/atomic"

	fuzzy "github.com/paul-mannino/go-fuzzywuzzy"

	sourcefileutils "github.com/Matthew-Curry/re-region-etl/sourceFileUtils"
)

const LOCAL_TAX_FILE = "data/Local_Income_Tax_Rates_2019.xlsx"

// helper method to format 2D array holding local tax data
func GetLocalTaxData(censusData [][]string, matchThresh int, nullString string) ([][]string, error) {

	// get data from sourcefileutils
	localTaxData, err := sourcefileutils.OpenExcelSheet(LOCAL_TAX_FILE, "Local Income Tax Rates")
	if err != nil {
		return nil, err
	}

	// iterate over the data and format each row
	// keep track of the current state
	state := ""
	// number of unmatched
	var unmatched uint64
	// processed data to return
	var processedLocalTaxData [][]string
	// wait group to manage syncing go routines
	var wg sync.WaitGroup
	// mutex for writing to data structure to return
	var mt sync.Mutex
	for i, row := range localTaxData {
		// skip empty rows and first row
		if len(row) < 3 || i == 0 {
			continue
		}
		// set the state if this record includes state. This way counties with same name in other state are not matched.
		if len(strings.TrimSpace(row[0])) != 0 {
			state = row[0]
		}

		juris := row[1]
		resident := row[2]
		nonresident := row[3]

		// state can be modified through other go routines advancing through the loop,
		// so make a copy for this loop to use
		state := state

		// process record in a go routine to exploit paralellism
		// add one to the wait for this record
		wg.Add(1)

		go func() {
			// use fuzzy matching to retrieve a county id
			county_id := getCountyId(censusData, state, juris, matchThresh, nullString)

			// increment unmatched atomically
			if county_id == nullString {
				atomic.AddUint64(&unmatched, 1)
			}

			// get the components of resident tax description
			resident_attr := getLocalTaxComponents(resident, nullString)

			// get the components of nonresident tax description
			nonresident_attr := getLocalTaxComponents(nonresident, nullString)

			processedRecord := append(append([]string{juris, county_id}, resident_attr...), nonresident_attr...)
			mt.Lock()
			processedLocalTaxData = append(processedLocalTaxData, processedRecord)
			mt.Unlock()
			// decrement the counter
			wg.Done()
		}()

	}

	// all go routines must finish before returining the data
	wg.Wait()

	if unmatched > 0 {
		logger.Warn("%v local tax jurisdictions were not able to be fuzzy matched out of %v. (%v %s).", unmatched, len(processedLocalTaxData), math.Round(float64(unmatched)/float64(len(processedLocalTaxData))*100), "%")
	}

	return processedLocalTaxData, nil

}

// helper method to decompose a description of local taxes into the component attributes
func getLocalTaxComponents(taxDesc, nullString string) []string {
	rate := nullString
	month := nullString
	year := nullString
	payPeriod := nullString
	stateLiability := nullString

	if strings.Contains(taxDesc, "-") {
		logger.Warn("%s contains range, cannot parse componenets", taxDesc)
	} else if strings.Contains(taxDesc, "month") {
		// split by delimiter, value is the first part, there is a leading dollar sign so get second character
		month = string(strings.TrimSpace(strings.Split(taxDesc, " / ")[0])[1:])
	} else if strings.Contains(taxDesc, "year") {
		// split by delimiter, value is the first part, there is a leading dollar sign so get second character
		year = string(strings.TrimSpace(strings.Split(taxDesc, " a ")[0])[1:])
	} else if strings.Contains(taxDesc, "per pay period") {
		// split by delimiter, value is the first part, there is a leading dollar sign so get second character
		payPeriod = string(strings.TrimSpace(strings.Split(taxDesc, " per pay ")[0])[1:])
	} else if strings.Contains(taxDesc, "of state liability") || strings.Contains(taxDesc, "of state tax") {
		// split by delimiter, value is the first part, remove the % sign
		stateLiability = strings.TrimSuffix(string(strings.TrimSpace(strings.Split(taxDesc, " of state ")[0])), "%")
	} else if strings.Contains(taxDesc, "+") {
		// first portion is rate, second is yearly fee
		splitDesc := strings.Split(taxDesc, " + ")
		rate = strings.TrimSpace(strings.TrimSuffix(splitDesc[0], "%"))
		// there is a value for the yearly fee in the second part
		year = strings.Split(splitDesc[1], " ")[0][1:]
	} else if strings.Contains(taxDesc, "no LST") {
		// first portion has the rate
		rate = strings.TrimSpace(strings.TrimSuffix(strings.Split(taxDesc, " (")[0], "%"))
	} else if strings.Contains(taxDesc, "%") && !strings.Contains(taxDesc, "dividends") && !strings.Contains(taxDesc, "to") {
		// else is rate, remove the % at the end
		rate = strings.TrimSuffix(strings.TrimSpace(taxDesc), "%")
	}

	return []string{taxDesc, rate, month, year, payPeriod, stateLiability}
}

// helper method that uses fuzzy matching to return the county id that matches the tax jurisdiction
func getCountyId(censusData [][]string, state string, juris string, matchThresh int, nullString string) string {
	// pass over census data, return id of greatest match above threshold
	max := 0
	match := nullString
	for _, row := range censusData {
		rowCounty := row[0]
		// parse the county portion of Juris if it exists to increase matches
		if strings.Contains(juris, "Co.") {
			split := strings.Split(juris, " (")
			if len(split) > 1 {
				juris = strings.TrimSuffix(split[1], ")")
			}
		}
		// try to match each ratio in order of decreasing liklihood of a match. County
		// id is the census data's state field concated with county, as a census county id is
		// only unique witihn a state
		ratio := fuzzy.PartialRatio(rowCounty, juris)
		if isNewFuzzyMatch(ratio, matchThresh, max, state, row[9]) {
			max = ratio
			match = row[7] + row[8]
		}

		ratio = fuzzy.TokenSortRatio(rowCounty, juris)
		if isNewFuzzyMatch(ratio, matchThresh, max, state, row[9]) {
			max = ratio
			match = row[7] + row[8]
		}

		ratio = fuzzy.TokenSetRatio(rowCounty, juris)
		if isNewFuzzyMatch(ratio, matchThresh, max, state, row[9]) {
			max = ratio
			match = row[7] + row[8]
		}

		ratio = fuzzy.Ratio(rowCounty, juris)
		if isNewFuzzyMatch(ratio, matchThresh, max, state, row[9]) {
			max = ratio
			match = row[7] + row[8]
		}

	}

	return match
}

// helper method, returns condition checking whether or not we have a new fuzzy match
func isNewFuzzyMatch(ratio int, thresh int, currentMax int, eq1 string, eq2 string) bool {
	return ratio > thresh && ratio > currentMax && strings.ToLower(strings.TrimSpace(eq1)) == strings.ToLower(strings.TrimSpace(eq2))
}
