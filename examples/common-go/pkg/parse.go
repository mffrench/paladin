package examples_commongo

import (
	"encoding/json"
	"os"
	"io"
)

func ParseJSONFileToByteArray(filepath string) ([]byte, error) {
	jsonFile, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()


	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	return byteValue, nil
}

func ParseJSONFileToStruct(filepath string, target interface{}) error {
	jsonFile, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	return json.Unmarshal(byteValue, target)
}
