package format

import (
	"encoding/json"
	"fmt"
)

// Map2JSON is format map to json
func Map2JSON(m interface{}) ([]byte, error) {
	str, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
	}

	return str, err
}
