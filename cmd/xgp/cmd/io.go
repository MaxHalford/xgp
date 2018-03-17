package cmd

import (
	"bufio"
	"encoding/csv"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/kniren/gota/dataframe"
	"github.com/kniren/gota/series"
)

// ReadFile loads a file with the appropriate method based on the file's
// extension.
func ReadFile(path string) (dataframe.DataFrame, time.Duration, error) {
	var (
		df  dataframe.DataFrame
		t   = time.Now()
		err error
	)
	// If the file is a CSV file
	if strings.HasSuffix(path, ".csv") {
		df, err = readCSV(path)
		return df, time.Since(t), err
	}
	return dataframe.DataFrame{}, 0, errors.New("Unknown file extension")
}

// readCSV loads a CSV file.
func readCSV(path string) (dataframe.DataFrame, error) {
	var f, err = os.Open(path)
	if err != nil {
		return dataframe.DataFrame{}, err
	}
	var r = csv.NewReader(bufio.NewReader(f))
	records, err := r.ReadAll()
	if err != nil {
		return dataframe.DataFrame{}, err
	}
	var df = dataframe.LoadRecords(
		records,
		dataframe.DefaultType(series.Float),
	)
	return df, nil
}
