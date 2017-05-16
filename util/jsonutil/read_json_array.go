package jsonutil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var mapJSONMap = []map[string]interface{}{}

// ReadFileJSONArray can read JSONArray
func ReadFileJSONArray(filename string) ([]map[string]interface{}, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("ReadFile: ", err.Error())
		return nil, err
	}

	if err := json.Unmarshal(bytes, &mapJSONMap); err != nil {
		fmt.Println("Unmarshal: ", err.Error())
		return nil, err
	}

	return mapJSONMap, nil
}
