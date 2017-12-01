package koza

import "github.com/kniren/gota/dataframe"

// removeString is a pure function to remove a slice from a slice.
func removeString(s []string, e string) []string {
	var c = make([]string, len(s))
	copy(c, s)
	for i := range s {
		if s[i] == e {
			c = append(c[:i], c[i+1:]...)
		}
	}
	return c
}

// containsString checks if a string is in a slice.
func containsString(s []string, e string) bool {
	for i := range s {
		if s[i] == e {
			return true
		}
	}
	return false
}

func dataFrameToFloat64(df dataframe.DataFrame) [][]float64 {
	var X = make([][]float64, df.Ncol())
	for i, col := range df.Names() {
		X[i] = df.Col(col).Float()
	}
	return X
}
