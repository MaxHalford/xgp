package cmd

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/MaxHalford/xgp"
	"github.com/MaxHalford/xgp/meta"
	"github.com/kniren/gota/dataframe"
	"github.com/kniren/gota/series"
)

const perm = 0644

// readFile loads a file with the appropriate method based on the file's
// extension.
func readFile(path string) (dataframe.DataFrame, time.Duration, error) {
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

func writeProgram(prog xgp.Program, path string) error {
	bytes, err := json.Marshal(serialModel{
		Flavor: "vanilla",
		Model:  prog,
	})
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, bytes, perm)
	return err
}

func writeGradientBoosting(gb *meta.GradientBoosting, path string) error {
	bytes, err := json.Marshal(serialModel{
		Flavor: "boosting",
		Model:  gb,
	})
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, bytes, perm)
	return err
}

func readModel(path string) (sm serialModel, err error) {
	// Read the model file
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	// Extract its keys
	var raw map[string]*json.RawMessage
	err = json.Unmarshal(bytes, &raw)
	if err != nil {
		return
	}
	// Extract the flavor
	err = json.Unmarshal(*raw["flavor"], &sm.Flavor)
	if err != nil {
		return
	}
	switch sm.Flavor {
	case "vanilla":
		var prog xgp.Program
		err = json.Unmarshal(*raw["model"], &prog)
		if err != nil {
			return
		}
		sm.Model = prog
		return
	case "boosting":
		var gb meta.GradientBoosting
		err = json.Unmarshal(*raw["model"], &gb)
		if err != nil {
			return
		}
		sm.Model = gb
		return
	}
	err = errUnknownFlavor{sm.Flavor}
	return
}
