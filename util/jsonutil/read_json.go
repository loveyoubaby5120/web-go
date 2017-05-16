package jsonutil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var mapJSON = map[string]string{}

// ReadFileJSON can read JSON
func ReadFileJSON(filename string) (map[string]string, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("ReadFile: ", err.Error())
		return nil, err
	}

	if err := json.Unmarshal(bytes, &mapJSON); err != nil {
		fmt.Println("Unmarshal: ", err.Error())
		return nil, err
	}

	return mapJSON, nil
}
