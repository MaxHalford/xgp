package boosting

import (
	"encoding/json"
	"io/ioutil"
)

// SaveBoosterToJSON saves a Booster to a JSON file.
func SaveBoosterToJSON(booster Booster, path string) error {
	var bytes, err = json.Marshal(&booster)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

// LoadBoosterFromJSON loads a Booster from a JSON file.
func LoadBoosterFromJSON(path string) (Booster, error) {
	var (
		booster    Booster
		bytes, err = ioutil.ReadFile(path)
	)
	if err != nil {
		return booster, err
	}
	err = json.Unmarshal(bytes, &booster)
	if err != nil {
		return booster, err
	}
	return booster, nil
}
