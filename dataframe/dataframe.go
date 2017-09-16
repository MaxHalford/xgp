package dataframe

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"text/tabwriter"
)

type DataFrame struct {
	X        [][]float64
	XNames   []string
	Y        []float64
	YName    string
	ClassMap *ClassMap
}

func (df DataFrame) NRows() int {
	return len(df.X)
}

func (df DataFrame) NFeatures() int {
	return len(df.X[0])
}

func (df DataFrame) Shape() (int, int) {
	return df.NRows(), df.NFeatures() + 1
}

func (df DataFrame) NClasses() (int, error) {
	if df.ClassMap == nil {
		return 0, errors.New("Target is not discrete")
	}
	return len(df.ClassMap.Map), nil
}

func (df DataFrame) Sample(k int, rng *rand.Rand) DataFrame {
	var (
		indices = randomInts(k, 0, len(df.X), rng)
		sample  = DataFrame{
			X:        make([][]float64, k),
			XNames:   df.XNames,
			Y:        make([]float64, k),
			YName:    df.YName,
			ClassMap: df.ClassMap,
		}
	)
	for i, idx := range indices {
		sample.X[i] = df.X[idx]
		sample.Y[i] = df.Y[idx]
	}
	return sample
}

func (df DataFrame) String() string {

	// Determine the length of the longest column name
	var colSize int
	for _, col := range df.XNames {
		if len(col) != colSize {
			colSize = len(col)
		}
	}

	var buffer bytes.Buffer
	//w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, '.', tabwriter.AlignRight|tabwriter.Debug)

	var w = new(tabwriter.Writer)
	w.Init(&buffer, 0, 8, 0, '\t', 0)

	// Display the column names
	fmt.Fprint(w, "\t")
	for _, name := range df.XNames {
		fmt.Fprint(w, fmt.Sprintf("\t%s", name))
	}
	fmt.Fprint(w, fmt.Sprintf("\t%s\n", df.YName))

	// Iterate over the rows
	var n = df.NRows()
	for i, X := range df.X {
		// Display the row number
		fmt.Fprintf(w, "\t%d", i)
		// Display the row content
		for _, x := range X {
			fmt.Fprintf(w, "\t%.3f", x)
		}
		// Display the target
		if df.ClassMap.N == 0 {
			fmt.Fprintf(w, "\t%.3f", df.Y[i])
		} else {
			fmt.Fprintf(w, "\t%s", df.ClassMap.ReverseMap[df.Y[i]])
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
