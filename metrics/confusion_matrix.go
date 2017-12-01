package metrics

import (
	"bytes"
	"fmt"
	"sort"
	"text/tabwriter"
)

// A ConfusionMatrix stores true positives (TP), true negatives (TN), false
// positives (FP) and false negatives (FN).
type ConfusionMatrix map[float64]map[float64]float64

// NClasses returns the number of classes in a ConfusionMatrix.
func (cm ConfusionMatrix) NClasses() int {
	return len(cm)
}

// Classes returns a slice of classes included in a ConfusionMatrix. The
// resulting slice of classes is ordered increasingly.
func (cm ConfusionMatrix) Classes() []float64 {
	var (
		classes = make([]float64, len(cm))
		i       int
	)
	for class := range cm {
		classes[i] = class
		i++
	}
	sort.Float64s(classes)
	return classes
}

// TruePositives returns the number of times a class was correctly predicted.
func (cm ConfusionMatrix) TruePositives(class float64) float64 {

	if _, ok := cm[class]; !ok {
		return 0
	}

	return cm[class][class]
}

// FalsePositives returns the number of times a class was wrongly predicted.
func (cm ConfusionMatrix) FalsePositives(class float64) float64 {

	if _, ok := cm[class]; !ok {
		return 0
	}

	var FP float64
	for tc := range cm {
		if tc != class {
			FP += cm[tc][class]
		}
	}
	return FP
}

// FalseNegatives returns the number of times a class was wrongly not predicted.
func (cm ConfusionMatrix) FalseNegatives(class float64) float64 {

	if _, ok := cm[class]; !ok {
		return 0
	}

	var FN float64
	for pc := range cm[class] {
		if pc != class {
			FN += cm[class][pc]
		}
	}
	return FN
}

// TrueNegatives returns the number of times a class was correctly not
// predicted.
func (cm ConfusionMatrix) TrueNegatives(class float64) float64 {

	if _, ok := cm[class]; !ok {
		return 0
	}

	var TN float64
	for tc := range cm {
		if tc != class {
			for pc := range cm[tc] {
				if pc != class {
					TN += cm[tc][pc]
				}
			}
		}
	}
	return TN
}

// String returns a string that can easily be read by a human in a terminal.
func (cm ConfusionMatrix) String() string {
	var (
		buffer bytes.Buffer
		w      = tabwriter.NewWriter(&buffer, 0, 8, 0, '\t', 0)
	)

	var classes = cm.Classes()

	// Display one column for each predicted class
	for _, class := range classes {
		fmt.Fprint(w, fmt.Sprintf("\tPredicted %0.f", class))
	}
	fmt.Fprint(w, "\t\n")

	// Display one row for each true class
	for i, tc := range classes {
		fmt.Fprintf(w, "True %0.f", tc)
		for _, pc := range classes {
			fmt.Fprintf(w, "\t%0.f", cm[tc][pc])
		}
		// Only add a carriage return if the current class is not the last one
		if i != len(classes)-1 {
			fmt.Fprint(w, "\t\n")
		} else {
			fmt.Fprint(w, "\t")
		}
	}

	w.Flush()
	return buffer.String()
}

// MakeConfusionMatrix returns a ConfusionMatrix from a slice of true classes
// and another slice of predicted classes.
func MakeConfusionMatrix(yTrue, yPred, weights []float64) (ConfusionMatrix, error) {

	if len(yTrue) != len(yPred) {
		return nil, &errMismatchedLengths{len(yTrue), len(yPred)}
	}
	if weights != nil && len(yTrue) != len(weights) {
		return nil, &errMismatchedLengths{len(yTrue), len(weights)}
	}

	var cm = make(ConfusionMatrix)

	for i, yt := range yTrue {
		if _, ok := cm[yt]; ok {
			if weights != nil {
				cm[yt][yPred[i]] += weights[i]
			} else {
				cm[yt][yPred[i]]++
			}
		} else {
			cm[yt] = make(map[float64]float64)
			if weights != nil {
				cm[yt][yPred[i]] = weights[i]
			} else {
				cm[yt][yPred[i]] = 1
			}
		}
	}

	return cm, nil
}
