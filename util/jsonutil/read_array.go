package jsonutil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var array = []int{}

func ReadFileArray(filename string) ([]int, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("ReadFile: ", err.Error())
		return nil, err
	}

	if err := json.Unmarshal(bytes, &array); err != nil {
		fmt.Println("Unmarshal: ", err.Error())
		return nil, err
	}

	return array, nil
}
