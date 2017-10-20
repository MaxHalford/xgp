package dataset

import (
	"bytes"
	"errors"
	"fmt"
	"text/tabwriter"
)

type Dataset struct {
	X        [][]float64
	XNames   []string
	Y        []float64
	YName    string
	ClassMap *ClassMap
}

func (dataset Dataset) NRows() int {
	return len(dataset.X)
}

func (dataset Dataset) NFeatures() int {
	return len(dataset.X[0])
}

func (dataset Dataset) Shape() (int, int) {
	return dataset.NRows(), dataset.NFeatures() + 1
}

func (dataset Dataset) NClasses() (int, error) {
	if dataset.ClassMap == nil {
		return 0, errors.New("Target is not discrete")
	}
	return len(dataset.ClassMap.Map), nil
}

func (dataset Dataset) String() string {

	// Determine the length of the longest column name
	var colSize int
	for _, col := range dataset.XNames {
		if len(col) != colSize {
			colSize = len(col)
		}
	}

	var (
		buffer bytes.Buffer
		w      = new(tabwriter.Writer)
	)
	w.Init(&buffer, 0, 8, 0, '\t', 0)

	// Display the column names
	fmt.Fprint(w, "\t")
	for _, name := range dataset.XNames {
		fmt.Fprint(w, fmt.Sprintf("\t%s", name))
	}
	fmt.Fprint(w, fmt.Sprintf("\t%s\n", dataset.YName))

	// Iterate over the rows
	var n = dataset.NRows()
	for i, X := range dataset.X {
		// Display the row number
		fmt.Fprintf(w, "\t%d", i)
		// Display the row content
		for _, x := range X {
			fmt.Fprintf(w, "\t%.3f", x)
		}
		// Display the target
		if dataset.ClassMap == nil {
			fmt.Fprintf(w, "\t%.3f", dataset.Y[i])
		} else {
			fmt.Fprintf(w, "\t%s", dataset.ClassMap.ReverseMap[dataset.Y[i]])
		}
		// Only add a carriage return if the current class is not the last one
		if i < n-1 {
			fmt.Fprint(w, "\t\n")
		} else {
			fmt.Fprint(w, "\t")
		}
	}

	w.Flush()
	return buffer.String()
}

// NewDatasetXY returns a Dataset from a set of features X and a target Y.
func NewDatasetXY(X [][]float64, Y []float64, classification bool) (*Dataset, error) {
	return &Dataset{
		X:        X,
		XNames:   []string{},
		Y:        Y,
		YName:    "y",
		ClassMap: &ClassMap{},
	}, nil
}
