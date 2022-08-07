/* Holds database engine struct to create, insert data, and define views on the database */

package load

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
	"strconv"
	"strings"

	_"github.com/lib/pq"

	"github.com/Matthew-Curry/re-region-etl/logging"
)

var logger, _ = logging.GetLogger("file.log")

// string constants
const (
	// table names
	COUNTY             string = "county"
	FEDERAL_DEDUCTIONS string = "federal_deductions"
	FEDERAL_BRACKETS   string = "federal_brackets"
	STATE_BRACKETS     string = "state_brackets"
	STATE              string = "states"
	TAX_JURISDICTION   string = "tax_locale"
	// common sql file names
	COUNTY_SQL            string = "county.sql"
	FEDERAL_DEDUCTION_SQL string = "federal_deductions.sql"
	FEDERAL_BRACKETS_SQL  string = "federal_brackets.sql"
	STATE_BRACKETS_SQL    string = "state_brackets.sql"
	STATE_SQL             string = "state.sql"
	TAX_JURISDICION_SQL   string = "tax_locale.sql"
	// directories holding each type of SQL
	DDL_DIR    string = "ddl"
	INSERT_DIR string = "insert"
	UPDATE_DIR string = "update"
	VIEW_DIR   string = "view"
	// ids of null county and state records
	countyNullId string = "32767"
	stateNullId  string = "32767"
)

type DbEngine struct {
	// the database connection
	con *sql.DB
	// mapping of table names to sql script names common across all sql types (other than view)
	sqlMap map[string]string
	// mapping of table names to tables they are dependent on.
	depMap map[string]string
	// the null string expected for nulls in data inputs
	nullString string
}

func NewDbEngine(nullString, dbUser, dbPassword, dbName, dbHost, dbPort string) (*DbEngine, error) {
	// define the DDL map
	sqlMap := map[string]string{
		COUNTY:             COUNTY_SQL,
		FEDERAL_DEDUCTIONS: FEDERAL_DEDUCTION_SQL,
		FEDERAL_BRACKETS:   FEDERAL_BRACKETS_SQL,
		STATE_BRACKETS:     STATE_BRACKETS_SQL,
		STATE:              STATE_SQL,
		TAX_JURISDICTION:   TAX_JURISDICION_SQL,
	}

	// the dependency table. Map of tables to tables needed
	// to create the table due to foriegn key constraints.
	// at this time map to strings, may expand to map of
	// lists if database expands in the future and needs
	// a schema re-design
	depMap := map[string]string{
		COUNTY:           STATE,
		STATE_BRACKETS:   STATE,
		TAX_JURISDICTION: COUNTY,
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	engine := DbEngine{con: db, sqlMap: sqlMap, depMap: depMap, nullString: nullString}
	logger.Info("New db engine successfully created")

	return &engine, nil
}

// setup method used by each load. Checks if table exists, if not creates it. Also, will clear
// given table if the provided flag is true.
func (d *DbEngine) loadSetup(table string, c bool) error {
	logger.Info("Checking if %s exists", table)
	tableExists, err := d.doesTableExist(table)
	if err != nil {
		return err
	}

	// if the table exists and a clear is called, drop all rows in the table, else make the table
	if tableExists && c {
		logger.Info("Table %s exists and a clear flags was passed, so records are being dropped from this table", table)
		err = d.deleteTable(table)
		if err != nil {
			return err
		}
	} else if !tableExists {
		logger.Info("Table %s does not exist. Creating the table", table)
		err = d.createTable(table)
		if err != nil {
			return err
		}
	}

	return nil

}

// helper method to create a given table.
func (d *DbEngine) createTable(table string) error {
	// check if this table has depdencies, raise error if dependency does not already exist
	if depTable, ok := d.depMap[table]; ok {
		tableExists, err := d.doesTableExist(depTable)
		if err != nil {
			logger.Warn("Unable to retrieve if the depedency table %s exists. Proceeding with load", depTable)
		}

		if !tableExists {
			return errors.New(fmt.Sprintf("The dependent table %s does not exist, so %s cannot be created.", depTable, table))
		}
	}
	// get the table's ddl
	ddl, err := d.readSQLFileAsString(table, "ddl")
	if err != nil {
		return err
	}

	// execute the ddl
	logger.Info("Executing DDL for %s table", table)
	_, err = d.con.Exec(ddl)
	if err != nil {
		return err
	}

	return nil

}

// helper method to delete all rows in a given table
func (d *DbEngine) deleteTable(table string) error {
	query := fmt.Sprintf("DELETE FROM %s", table)

	_, err := d.con.Exec(query)

	return err
}

// helper method to check if a given table exists in the database
func (d *DbEngine) doesTableExist(table string) (bool, error) {
	query :=
		`SELECT EXISTS (
		SELECT FROM 
			pg_tables
		WHERE 
			schemaname = 'public' AND 
			tablename  = $1
		);`
	var result bool
	row := d.con.QueryRow(query, table)
	err := row.Scan(&result)

	if err != nil {
		return false, err
	}

	return result, nil

}

// helper method to read given sql type in for a given table
func (d *DbEngine) readSQLFileAsString(table, sqlType string) (string, error) {
	var filePath string
	if sql, ok := d.sqlMap[table]; ok {
		filePath = sqlType + "/" + sql
	} else {
		return "", fmt.Errorf("There is no SQL file for table %s", table)
	}

	return d.readFileFromSqlDir(filePath)
}

// helper method to read a SQL file in starting from the SQL directory
func (d *DbEngine) readFileFromSqlDir(filePath string) (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("No caller information")
	}

	sqlDir := path.Dir(filename)

	s, e := ioutil.ReadFile(sqlDir + "/sql/" + filePath)
	if e != nil {
		return "", e
	}
	return string(s), nil
}

// helper method to execute an insert query
func (d *DbEngine) executeInsertStatement(query string, vals []interface{}, records int) error {
	// ? -> $n for postgres
	paramCount := strings.Count(query, "?")
	for n := 1; n <= paramCount; n++ {
		query = strings.Replace(query, "?", "$"+strconv.Itoa(n), 1)
	}

	// format all vals at once
	_, err := d.con.Exec(query, vals...)
	if err != nil {
		return err
	}

	logger.Info("Successfully upserted data for %v records.", records)

	return nil
}

// helper method to convert empty strings + null strings to string representing nil int
func (d *DbEngine) NewNullIntStr(s string) string {
	nilInt := "0"
	if len(s) == 0 {
		return nilInt
	} else if s == d.nullString {
		return nilInt
	}

	return s
}

// helper method to convert empty strings + null strings to nulls
func (d *DbEngine) NewNullDecStr(s string) string {
	nilDec := "0.00"
	if len(s) == 0 {
		return nilDec
	} else if s == d.nullString {
		return nilDec
	}
	
	return s
}

// helper method to convert empty strings + null strings to nulls
func (d *DbEngine) NewNullFloat(s string) sql.NullFloat64 {
	if len(s) == 0 {
		return sql.NullFloat64{}
	} 
	
	f, _ := strconv.ParseFloat(s, 10)

	return sql.NullFloat64{
		Float64: f,
		Valid:  true,
	}
}

// public method to create the county table
func (d *DbEngine) LoadCountyTable(data [][]string, c bool) error {
	logger.Info("Executing insert for county table")
	err := d.loadSetup(COUNTY, c)
	if err != nil {
		return err
	}

	query, err := d.readSQLFileAsString(COUNTY, "insert")
	if err != nil {
		return err
	}
	vals := []interface{}{countyNullId, countyNullId, countyNullId, countyNullId, countyNullId, countyNullId, countyNullId, countyNullId, countyNullId}

	for _, row := range data {
		// county id is state concated with the county id from the input data
		county_id := row[7] + row[8]

		// if state is none, set state id to the null record
		var state_id string
		if row[7] == d.nullString {
			state_id = stateNullId
		} else {
			state_id = row[7]
		}

		query += "(?, ?, ?, ?, ?, ?, ?, ?, ?), "
		vals = append(vals, county_id, row[0], state_id, row[1], row[2], row[3], row[4], row[5], d.NewNullIntStr(row[6]))
	}

	// append the update SQL
	updateSql, err := d.readSQLFileAsString(COUNTY, "update")
	if err != nil {
		return err
	}

	query = strings.TrimSuffix(query, ", ")
	query += " "
	query += updateSql

	// execute the formed insert statement
	return d.executeInsertStatement(query, vals, len(data))

}

// method to create the local tax jurisdiction table
func (d *DbEngine) LoadLocalTaxTable(data [][]string, c bool) error {
	logger.Info("Executing insert for local tax table")
	err := d.loadSetup(TAX_JURISDICTION, c)
	if err != nil {
		return err
	}

	query, err := d.readSQLFileAsString(TAX_JURISDICTION, "insert")
	if err != nil {
		return err
	}

	// Psql has a max 65535 params per query. With 15 params per row, a max of 4369 rows can be inserted per
	// insert query.
	maxDataPartSize := 4369
	start := 0
	moreData := true
	dataSize := len(data)
	var dataPart [][]string
	for moreData {
		if dataSize-start > maxDataPartSize {
			dataPart = data[start : start+maxDataPartSize]
		} else {
			dataPart = data[start:]
			moreData = false
		}
		// process this part
		err = d.loadLocalTaxPart(dataPart, query, start)
		if err != nil {
			return err
		}
		// increment start
		start = start + maxDataPartSize
	}

	return nil
}

// helper method to load a portion of the local tax table due to Postgresql parameter constraints
func (d *DbEngine) loadLocalTaxPart(data [][]string, query string, pos int) error {
	vals := []interface{}{}
	for _, row := range data {
		query += "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?), "
		// if county id is NONE, set to null county id
		var county_id string
		if row[1] == d.nullString {
			county_id = countyNullId
		} else {
			county_id = row[1]
		}

		vals = append(vals, fmt.Sprint(pos), row[0], county_id, row[2], d.NewNullDecStr(row[3]), d.NewNullDecStr(row[4]), d.NewNullDecStr(row[5]), d.NewNullDecStr(row[6]),
			d.NewNullDecStr(row[7]), row[8], d.NewNullDecStr(row[9]), d.NewNullDecStr(row[10]),
			d.NewNullDecStr(row[11]), d.NewNullDecStr(row[12]), d.NewNullDecStr(row[13]))

		// increment the position
		pos = pos + 1

	}

	// append the update
	updateSql, err := d.readSQLFileAsString(TAX_JURISDICTION, "update")
	if err != nil {
		return err
	}

	query = strings.TrimSuffix(query, ", ")
	query += " "
	query += updateSql

	return d.executeInsertStatement(query, vals, len(data))
}

// method to create the state table
func (d *DbEngine) LoadStateTable(data [][]string, c bool) error {
	logger.Info("Executing insert for state table")
	err := d.loadSetup(STATE, c)
	if err != nil {
		return err
	}

	query, err := d.readSQLFileAsString(STATE, "insert")
	if err != nil {
		return err
	}
	vals := []interface{}{stateNullId, stateNullId, stateNullId, stateNullId, stateNullId, stateNullId, stateNullId}

	for _, row := range data {
		query += "(?, ?, ?, ?, ?, ?, ?), "
		vals = append(vals, row[0], row[1], d.NewNullIntStr(row[2]), d.NewNullIntStr(row[3]), d.NewNullIntStr(row[4]), d.NewNullIntStr(row[5]), d.NewNullIntStr(row[6]))
	}

	updateSql, err := d.readSQLFileAsString(STATE, "update")
	if err != nil {
		return err
	}

	query = strings.TrimSuffix(query, ", ")
	query += " "
	query += updateSql

	// execute the formed insert statement
	return d.executeInsertStatement(query, vals, len(data))
}

// method to create the state bracket table
func (d *DbEngine) LoadStateBracketTable(data [][]string, c bool) error {
	logger.Info("Executing insert for state bracket table")
	err := d.loadSetup(STATE_BRACKETS, c)
	if err != nil {
		return err
	}

	query, err := d.readSQLFileAsString(STATE_BRACKETS, "insert")
	if err != nil {
		return err
	}
	vals := []interface{}{}

	for _, row := range data {
		query += "(?, ?, ?, ?, ?), "
		vals = append(vals, row[0], d.NewNullDecStr(row[1]), d.NewNullIntStr(row[2]), d.NewNullDecStr(row[3]), d.NewNullIntStr(row[4]))
	}

	updateSql, err := d.readSQLFileAsString(STATE_BRACKETS, "update")
	if err != nil {
		return err
	}

	query = strings.TrimSuffix(query, ", ")
	query += " "
	query += updateSql

	// execute the formed insert statement
	return d.executeInsertStatement(query, vals, len(data))
}

// method to create the federal deductions table
func (d *DbEngine) LoadFederalDeductionsTable(data []string, c bool) error {
	logger.Info("Executing insert for federal deductions table")
	err := d.loadSetup(FEDERAL_DEDUCTIONS, c)
	if err != nil {
		return err
	}

	query, err := d.readSQLFileAsString(FEDERAL_DEDUCTIONS, "insert")
	if err != nil {
		return err
	}
	vals := []interface{}{}

	query += "(?, ?, ?), "
	vals = append(vals, data[0], data[1], data[2])

	updateSql, err := d.readSQLFileAsString(FEDERAL_DEDUCTIONS, "update")
	if err != nil {
		return err
	}
	query = strings.TrimSuffix(query, ", ")
	query += " "
	query += updateSql

	// execute the formed insert statement
	return d.executeInsertStatement(query, vals, len(data))
}

// method to create the federal brackets table
func (d *DbEngine) LoadFederalBracketTable(data [][]string, c bool) error {
	logger.Info("Executing insert for federal bracket table")
	err := d.loadSetup(FEDERAL_BRACKETS, c)
	if err != nil {
		return err
	}

	query, err := d.readSQLFileAsString(FEDERAL_BRACKETS, "insert")
	if err != nil {
		return err
	}
	vals := []interface{}{}

	for _, row := range data {
		query += "(?, ?, ?, ?), "
		vals = append(vals, row[0], row[1], row[2], row[3])
	}

	updateSql, err := d.readSQLFileAsString(FEDERAL_BRACKETS, "update")
	if err != nil {
		return err
	}

	query = strings.TrimSuffix(query, ", ")
	query += " "
	query += updateSql

	// execute the formed insert statement
	return d.executeInsertStatement(query, vals, len(data))
}

// public method used to refresh all views
func (d *DbEngine) RefreshViews() error {

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("No caller information")
	}

	viewDir := path.Dir(filename) + "/sql/" + VIEW_DIR + "/"

	views, err := ioutil.ReadDir(viewDir)
	if err != nil {
		return err
	}

	for _, view := range views {
		viewQuery, err := d.readFileFromSqlDir("view/" + view.Name())
		if err != nil {
			return err
		}
		_, err = d.con.Exec(viewQuery)

		if err != nil {
			return err
		}

		logger.Info("Successfully defined view %s", view.Name())
	}

	return nil
}
