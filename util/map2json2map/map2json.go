package map2json2map

import (
	"encoding/json"
	"fmt"
)

func Map2json(m map[string]interface{}) ([]byte, error) {
	str, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
	}

	return str, err
}
