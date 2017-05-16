package map2json2map

import (
	"encoding/json"
	"fmt"
)

func Json2map(str []byte) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal(str, &m)
	if err != nil {
		fmt.Println(err)
	}

	return m, err
}
