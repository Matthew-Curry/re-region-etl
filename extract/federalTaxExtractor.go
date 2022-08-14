/* Logic to extract and store federal tax data into an in memory data structure for use in ETL */

package extract

import (
	"regexp"
	"strings"

	sourcefileutils "github.com/Matthew-Curry/re-region-etl/sourceFileUtils"
)

const FEDERAL_TAX_FILE = "data/2022-Federal-Income-Tax-Rates-and-Brackets-Tax-Foundation.xlsx"

// public method to build data structures for federal tax brackets and exemptions
func GetFederalTaxData() ([][]string, []string, error) {
	// read in the federal individual sheets
	federalBrackets, err := sourcefileutils.OpenExcelSheet(FEDERAL_TAX_FILE, "Table 1")
	if err != nil {
		return nil, nil, err
	}

	federalDeductions, err := sourcefileutils.OpenExcelSheet(FEDERAL_TAX_FILE, "Table 2")
	if err != nil {
		return nil, nil, err
	}

	// pass over brackets and format data structure
	federalBrackets = federalBrackets[2 : len(federalBrackets)-1]

	for i, row := range federalBrackets {

		row[0] = "0." + strings.TrimSuffix(strings.TrimSpace(row[0]), "%")
		row[1] = processFederalBracket(row[1])
		row[2] = processFederalBracket(row[2])
		row[3] = processFederalBracket(row[3])

		federalBrackets[i] = row

	}

	formattedFederalDeductions := []string{federalDeductions[2][1], federalDeductions[3][1], federalDeductions[4][1]}

	return federalBrackets, formattedFederalDeductions, nil

}

// helper method with logic to federal bracket values
func processFederalBracket(bracket string) string {

	bracket = strings.TrimSpace(strings.Split(bracket, " to ")[0])

	reg, _ := regexp.Compile("[^0-9]+")
	rep := reg.ReplaceAllString(bracket, "")

	return rep
}
