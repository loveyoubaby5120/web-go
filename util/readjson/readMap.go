package json

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var array = []int{}

func readFileArray(filename string) ([]int, error) {
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
