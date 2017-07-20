package dataframe

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

func ReadCSV(path string, target string, classification bool) (*DataFrame, error) {
	var (
		df   = &DataFrame{}
		f, _ = os.Open(path)
		r    = csv.NewReader(bufio.NewReader(f))
	)
	defer f.Close()

	// Read the headers
	columns, err := r.Read()
	if err == io.EOF {
		return nil, err
	}

	// Go through column names and determine which one is the target
	var (
		targetIdx = -1
		nXCols    int
	)
	df.XNames = make([]string, len(columns)-1)
	for i, column := range columns {
		if column == target {
			df.YName = column
			targetIdx = i
			break
		} else {
			if nXCols == len(df.XNames) {
				break
			}
			df.XNames[nXCols] = column
			nXCols++
		}
	}
	if targetIdx == -1 {
		return nil, fmt.Errorf("No column named '%s'", target)
	}

	// Initialize an empty class map in case of classification
	if classification {
		df.ClassMap = MakeClassMap()
	}

	// Iterate over the rows
	for {
		var record, err = r.Read()
		if err == io.EOF {
			break
		}
		// Parse the features as float64s
		var x = make([]float64, len(record)-1)
		for i, s := range record {
			if i != targetIdx {
				x[i], _ = strconv.ParseFloat(s, 64)
			}
		}
		df.X = append(df.X, x)
		// Parse the target
		if classification {
			df.Y = append(df.Y, df.ClassMap.Get(record[targetIdx]))
		} else {
			var y, _ = strconv.ParseFloat(record[targetIdx], 64)
			df.Y = append(df.Y, y)
		}
	}
	return df, nil
}
