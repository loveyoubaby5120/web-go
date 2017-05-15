package jsonutil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var mapJsonMap = []map[string]interface{}{}

func ReadFileJsonArray(filename string) ([]map[string]interface{}, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("ReadFile: ", err.Error())
		return nil, err
	}

	if err := json.Unmarshal(bytes, &mapJsonMap); err != nil {
		fmt.Println("Unmarshal: ", err.Error())
		return nil, err
	}

	return mapJsonMap, nil
}
