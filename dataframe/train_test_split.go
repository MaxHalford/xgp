package dataframe

import "math/rand"

func (df DataFrame) TrainTestSplit(testRatio float64, rng *rand.Rand) (DataFrame, DataFrame) {
	var (
		n, _     = df.Shape()
		testSize = int(float64(len(df.X)) * testRatio)
		indexes  = rng.Perm(n)
		test     = DataFrame{
			X:        make([][]float64, testSize),
			XNames:   df.XNames,
			Y:        make([]float64, testSize),
			YName:    df.YName,
			ClassMap: df.ClassMap,
		}
		train = DataFrame{
			X:        make([][]float64, n-testSize),
			XNames:   df.XNames,
			Y:        make([]float64, n-testSize),
			YName:    df.YName,
			ClassMap: df.ClassMap,
		}
	)
	for i := 0; i < testSize; i++ {
		test.X[i] = df.X[indexes[i]]
		test.Y[i] = df.Y[indexes[i]]
	}
	for i := testSize; i < n; i++ {
		train.X[i] = df.X[indexes[i]]
		train.Y[i] = df.Y[indexes[i]]
	}
	return train, test
}
