package format

import (
	"encoding/json"
	"fmt"
)

// JSON2map is format json to map
func JSON2map(str []byte) (interface{}, error) {
	var m interface{}
	err := json.Unmarshal(str, &m)
	if err != nil {
		fmt.Println(err)
	}

	return m, err
}
