package ensemble

import (
	"encoding/json"
	"io/ioutil"
)

// SaveEnsembleToJSON saves a Ensemble to a JSON file.
func SaveEnsembleToJSON(ensemble Ensemble, path string) error {
	var bytes, err = json.Marshal(&ensemble)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

// LoadEnsembleFromJSON loads a Ensemble from a JSON file.
func LoadEnsembleFromJSON(path string) (Ensemble, error) {
	var (
		ensemble   Ensemble
		bytes, err = ioutil.ReadFile(path)
	)
	if err != nil {
		return ensemble, err
	}
	err = json.Unmarshal(bytes, &ensemble)
	if err != nil {
		return ensemble, err
	}
	return ensemble, nil
}
