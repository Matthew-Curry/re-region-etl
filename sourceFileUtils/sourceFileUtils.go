/* Holds utility function for reading in data from excel files */

package sourcefileutils

import (
	"fmt"

	"github.com/Matthew-Curry/re-region-etl/logging"
	"github.com/xuri/excelize/v2"
)

// logger for the package
var logger, _ = logging.GetLogger("file.log")

func OpenExcelSheet(filePath string, sheet string) ([][]string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("There was an error reading in the excel file %s: %s", filePath, err)
	}

	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("There was an error reading in the rows of the excel file %s: %s", filePath, err)
	}

	logger.Info("Loaded the file %s successfully", filePath)

	return rows, nil

}
