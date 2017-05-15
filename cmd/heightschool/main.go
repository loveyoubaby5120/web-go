package main

import "fmt"

func main() {
	xxxMap, err := readjson.readFileArray("static/predict_num_v2/constrain_json/文一批_一分一段.json")
	if err != nil {
		fmt.Println("readFile: ", err.Error())
		// return nil, err
	}

	fmt.Println(xxxMap)
}
