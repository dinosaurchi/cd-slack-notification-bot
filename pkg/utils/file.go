package utils

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

func DumpToFile(filepath string, data interface{}) error {
	file, err := os.Create(filepath)
	if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ") // Set indent for pretty print
	err = encoder.Encode(data)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func LoadFromFile(filepath string, data interface{}) error {
	file, err := os.Open(filepath)
	if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(data)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
