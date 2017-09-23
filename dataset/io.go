package dataset

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

func ReadCSV(path string, target string, classification bool) (*Dataset, error) {
	var (
		dataset = &Dataset{}
		f, _    = os.Open(path)
		r       = csv.NewReader(bufio.NewReader(f))
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
	dataset.XNames = make([]string, len(columns)-1)
	for i, column := range columns {
		if column == target {
			dataset.YName = column
			targetIdx = i
			break
		} else {
			if nXCols == len(dataset.XNames) {
				break
			}
			dataset.XNames[nXCols] = column
			nXCols++
		}
	}
	if targetIdx == -1 {
		return nil, fmt.Errorf("No column named '%s'", target)
	}

	// Initialize an empty class map in case of classification
	if classification {
		dataset.ClassMap = MakeClassMap()
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
		dataset.X = append(dataset.X, x)
		// Parse the target
		if classification {
			dataset.Y = append(dataset.Y, dataset.ClassMap.Get(record[targetIdx]))
		} else {
			var y, _ = strconv.ParseFloat(record[targetIdx], 64)
			dataset.Y = append(dataset.Y, y)
		}
	}
	return dataset, nil
}
