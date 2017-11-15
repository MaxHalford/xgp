package dataset

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

func NewFromXY(X [][]float64, Y []float64, XNames []string, classification bool) (*Dataset, error) {
	var dataset = &Dataset{
		X:      X,
		XNames: XNames,
		Y:      Y,
	}
	// If the task is classification then the number of classes has to be determined
	var seen = make(map[float64]bool)
	for _, y := range Y {
		if _, ok := seen[y]; ok {
			seen[y] = true
		}
	}
	if classification {
		dataset.classMap = &classMap{N: len(seen)}
	}
	return dataset, nil
}

func ReadCSV(path string, target string, classification bool) (*Dataset, error) {
	var (
		dataset = &Dataset{}
		f, err  = os.Open(path)
	)
	if err != nil {
		return nil, err
	}
	var r = csv.NewReader(bufio.NewReader(f))
	defer f.Close()

	// Read the headers
	columns, err := r.Read()
	if err != nil {
		return nil, err
	}

	// Go through column names and determine which one is the target
	var (
		targetIdx = -1
		p         int
	)
	dataset.XNames = make([]string, len(columns)-1)
	for i, column := range columns {
		if column == target {
			targetIdx = i
			continue
		}
		if p == len(dataset.XNames) {
			break
		}
		dataset.XNames[p] = column
		p++
	}
	if targetIdx == -1 {
		return nil, fmt.Errorf("No column named '%s'", target)
	}

	// Initialize the columns
	dataset.X = make([][]float64, p)

	// Initialize an empty class map in case of classification
	if classification {
		dataset.classMap = makeClassMap()
	}

	// Iterate over the rows
	for {
		var record, err = r.Read()
		if err == io.EOF {
			break
		}
		// Parse the features as float64s
		var j int
		for i, s := range record {
			if i == targetIdx {
				continue
			}
			x, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return nil, err
			}
			dataset.X[j] = append(dataset.X[j], x)
			j++
		}
		// Parse the target
		if classification {
			dataset.Y = append(dataset.Y, dataset.classMap.Get(record[targetIdx]))
		} else {
			var y, _ = strconv.ParseFloat(record[targetIdx], 64)
			dataset.Y = append(dataset.Y, y)
		}
	}
	return dataset, nil
}
