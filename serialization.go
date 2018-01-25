package koza

import (
	"encoding/json"
	"io/ioutil"
)

// SaveProgramToJSON saves a Program to a JSON file.
func SaveProgramToJSON(program Program, path string) error {
	var bytes, err = json.Marshal(&program)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

// LoadProgramFromJSON loads a Program from a JSON file.
func LoadProgramFromJSON(path string) (Program, error) {
	var (
		program    Program
		bytes, err = ioutil.ReadFile(path)
	)
	if err != nil {
		return program, err
	}
	err = json.Unmarshal(bytes, &program)
	if err != nil {
		return program, err
	}
	return program, nil
}
