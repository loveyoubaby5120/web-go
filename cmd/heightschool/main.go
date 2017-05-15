package main

import (
	"fmt"
	"web-go/util/jsonutil"
)

type Condition struct {
	enroll_num int
	avg_score  int
	low_score  int
	rank       int
	major      string
}

func main() {
	// scoreArray, err := jsonutil.ReadFileArray("static/predict_num_v2/constrain_json/文一批_一分一段.json")
	// if err != nil {
	// 	fmt.Println("readFile: ", err.Error())
	// 	// return nil, err
	// }
	scoreCondition, err := jsonutil.ReadFileJsonArray("static/predict_num_v2/constrain_json/文一批_高考分数.json")
	if err != nil {
		fmt.Println("readFile: ", err.Error())
		// return nil, err
	}

	jsonMap := []map[string]interface{}{}

	for _, value := range scoreCondition {
		fmt.Println(value)
		var json Condition
		json.enroll_num = value["enroll_num"]
		json.avg_score = value["avg_score"]
		json.low_score = value["low_score"]
		json.rank = value["rank"]
		json.major = value["major"]
		jsonMap := append(jsonMap, value)
		break
	}

	fmt.Printf("111")

	// fmt.Println(scoreArray)
	// fmt.Println(scoreCondition)
}
