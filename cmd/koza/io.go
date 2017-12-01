package koza

import (
	"bufio"
	"encoding/csv"
	"os"

	"github.com/kniren/gota/dataframe"
	"github.com/kniren/gota/series"
)

// ReadCSV loads a CSV file.
func ReadCSV(path string) (dataframe.DataFrame, error) {
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
