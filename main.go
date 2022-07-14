/* CLI used to run defined ETL stages to setup database of ReRegion */

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/Matthew-Curry/re-region-etl/db"
	"github.com/Matthew-Curry/re-region-etl/logging"
	"github.com/Matthew-Curry/re-region-etl/transform"
)

// logger for the app
var logger logging.Logger
var logFile *os.File

func main() {
	// the logger and file to close
	logger, logFile = logging.GetLogger("file.log")
	defer logFile.Close()
	// read in config
	yfile, _ := ioutil.ReadFile("config.yml")
	configData := make(map[string]map[string]interface{})
	yaml.Unmarshal(yfile, &configData)
	censusAttempts := configData["census"]["attempts"]
	matchThresh := configData["localTax"]["threshold"]
	nullString := configData["general"]["nullString"]
	// CLI interface of app
	c := flag.Bool("c", false, "c, Clears existing tables resulting from the ETL stages to run. Will only take effect if l option is provided to trigger the ETL")
	l := flag.Bool("l", false, "l, Runs the ETL code to load the tables")
	v := flag.Bool("v", false, "v, Runs SQL to define the views.")
	flag.Parse()
	// remaining args define stages
	var stages []string
	if stages = flag.Args(); len(stages) == 0 {
		stages = []string{"1", "2", "3", "4"}
	}

	// the db engine to run SQL queries
	engine, err := db.NewDbEngine(nullString.(string))
	if err != nil {
		logger.Error("Unable to create the db engine. Recieved error: %s", err)
	}

	// run the ETL with the provided parameters if l option provided
	if *l == true {
		runETL(*c, stages, censusAttempts.(int), matchThresh.(int), nullString.(string), engine)
	}

	// refresh the views if the v option is provided
	if *v == true {
		err = engine.RefreshViews()
		if err != nil {
			logger.Error("Unable to refresh views. Recieved error: %s", err)
		}
	}

}

func runETL(c bool, stages []string, censusAttempts int, matchThresh int, nullString string, engine *db.DbEngine) {
	// initialized in memory data structures to load to tables
	var censusData [][]string
	var localTaxData [][]string
	var stateBrackets [][]string
	var stateExemptions [][]string
	var federalBrackets [][]string
	var federalDeductions []string

	var err error
	// load data in order of descending geography. This is the order dictated by the required database dependencies.
	// The first stage for the federal data is independent, but the next 3 are linked and will load the prior dependent
	// stage if specified (i.e passing stage 4 will load 2, 3 and 4 out of necessity)

	if contains(stages, "1") {
		logger.Info("RUNNING STAGE 1, LOAD TO FEDERAL TABLES")
		// retrieve 2d array of federal tax data
		federalBrackets, federalDeductions, err = transform.GetFederalTaxData()

		if err != nil {
			logger.Error(getDataErrorStr("federal", err))
		}

		// load the 2 federal tables
		err = engine.LoadFederalDeductionsTable(federalDeductions, c)
		if err != nil {
			logger.Error(getLoadErrorStr("federal", err))
		}

		err = engine.LoadFederalBracketTable(federalBrackets, c)
		if err != nil {
			logger.Error(getLoadErrorStr("federal bracket", err))
		}

	}

	// load if stage 2 is requested or any more granular geography
	if contains(stages, "2") || contains(stages, "3") || contains(stages, "4") {
		logger.Info("RUNNING STAGE 2, LOAD TO STATE TABLE")
		// get census data at the county level as 2D array
		censusData, err = transform.GetCensusData(censusAttempts)
		if err != nil {
			logger.Error(getDataErrorStr("census", err))
		}

		// get the state data as a 2d array
		stateBrackets, stateExemptions, err = transform.GetStateTaxData(censusData, nullString)

		// use the state data to load the state tables in order of dependencies
		err = engine.LoadStateTable(stateExemptions, c)
		if err != nil {
			logger.Error(getLoadErrorStr("state", err))
		}

		err = engine.LoadStateBracketTable(stateBrackets, c)
		if err != nil {
			logger.Error(getLoadErrorStr("state bracket", err))
		}
	}

	// load stage 3 if requested or any more granular geography
	if contains(stages, "3") || contains(stages, "4") {
		logger.Info("RUNNING STAGE 3, LOAD TO COUNTY TABLE")
		// state table is loaded, so census data can now be used to load the county table
		err = engine.LoadCountyTable(censusData, c)

		if err != nil {
			logger.Error(getLoadErrorStr("county", err))
		}
	}

	if contains(stages, "4") {
		logger.Info("RUNNING STAGE 4, LOAD TO LOCAL TAX JURISDICTION TABLE")
		// retrieve 2d array of state tax data
		localTaxData, err = transform.GetLocalTaxData(censusData, matchThresh, nullString)

		if err != nil {
			logger.Error(getDataErrorStr("local tax", err))
		}

		err = engine.LoadLocalTaxTable(localTaxData, c)

		if err != nil {
			logger.Error(getLoadErrorStr("local tax", err))
		}

	}

}

// helper method to return error string for a load error
func getLoadErrorStr(table string, e error) string {
	return fmt.Sprintf("Load to the %s table failed. Error: %s", table, e)
}

// helper method to return error string for a data error
func getDataErrorStr(data string, e error) string {
	return fmt.Sprintf("Unable to get %s data. Recieved error: %s", data, e)
}

// helper method to check if string slice contains a given string
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
