/* Logic to extract and store census data into an in memory data structure for use in ETL */

package extract

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Matthew-Curry/re-region-etl/logging"
)

// Census API constants
const (
	CENSUS_URL        = "https://api.census.gov/data/2019/acs/acs1"
	CENSUS_GET_PARAMS = "NAME,B01003_001E,B01001_002E,B01001_026E,B19013_001E,B25031_001E,C08536_001E"
	CENSUS_GEO        = "COUNTY"
	CENSUS_API_KEY    = "CENSUS_API_KEY"
)

// only one file within the package needs to define the logger and this
// is the one arbitrarily chosen
var logger, _ = logging.GetLogger("file.log")

// public method to retrieve census data from Census API for given
// number of attempts to make on a failed response
func GetCensusData(attempts int) ([][]string, error) {
	// retrieve the secret key from the environment
	key := os.Getenv(CENSUS_API_KEY)

	// format path with parameters and execute the request
	params := "?get=" + CENSUS_GET_PARAMS + "&" +
		"for=" + CENSUS_GEO + "&" +
		"key=" + key
	path := fmt.Sprintf(CENSUS_URL+"%s", params)

	body, err := executeGetRequest(path, attempts)
	if err != nil {
		return nil, err
	}

	// format the response as a slice of string slices
	var censusResp [][]string
	json.Unmarshal(body, &censusResp)

	// run business logic to process the response
	censusResp = processApiResponse(censusResp)

	return censusResp, nil
}

// helper method to connect to the given Census API path in an increasing retry count
func executeGetRequest(path string, attempts int) ([]byte, error) {
	logger.Info("Connecting to census API")
	// make the request in an increasing retry count
	var resp *http.Response
	var err error = nil
	for i := 1; i <= attempts; i++ {
		resp, err = http.Get(path)

		if err != nil && i == attempts {
			return nil, fmt.Errorf("Exceeded %v attempts trying to connect to Census API", attempts)
		} else if err != nil && i < attempts {
			sleepTime := 10 * i
			logger.Warn("Connection to Census API failed. Sleeping for %v and trying again", sleepTime)
			time.Sleep(time.Duration(sleepTime) * time.Second)
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && i == attempts {
			return nil, fmt.Errorf("Exceeded %v attempts trying to connect to Census API. Recieved status code %v from Census API",
				attempts, resp.StatusCode)
		} else if resp.StatusCode != http.StatusOK && i < attempts {
			sleepTime := 10 * i
			logger.Warn("Connection to Census API recieved status code %v failed. Sleeping for %v and trying again",
				resp.StatusCode, sleepTime)
			time.Sleep(time.Duration(sleepTime))
		}
	}

	logger.Info("Loading API Response")
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("An error has occured reading in the response: %s", err)
	}

	return body, nil
}

// holds business logic to process response from the Census API
func processApiResponse(censusResp [][]string) [][]string {
	var censusResp2 [][]string
	for _, row := range censusResp[1:] {
		// the GEO field will have the county and the state concated split by a ",".
		splitGeo := strings.Split(row[0], ", ")
		county := splitGeo[0]
		state := splitGeo[1]

		row[0] = county
		row = append(row, state)

		// Exclude DC and Puerto Rico to ensure matches to other datasets succeed
		if state == "District of Columbia" || state == "Puerto Rico" {
			continue
		}

		// divide the aggregated commute by population to get average commute
		agCommute, _ := strconv.Atoi(row[6])
		pop, _ := strconv.Atoi(row[1])

		row[6] = strconv.Itoa(agCommute / pop)

		// append altered row to the response
		censusResp2 = append(censusResp2, row)
	}

	return censusResp2
}
