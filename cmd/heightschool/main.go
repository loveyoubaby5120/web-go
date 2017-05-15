package main

import (
	"fmt"
	"web-go/util/jsonutil"
)

type Condition struct {
	enrollNum   int
	avgScore    int
	lowScore    int
	rank        int
	major       string
	enrollScore []int
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

	jsonMap := []Condition{}

	for _, value := range scoreCondition {
		fmt.Println(value)
		fmt.Println(value["enroll_num"])
		var json Condition
		switch v := value["enroll_num"].(type) {
		case int:
			json.enrollNum = v
		}

		switch v := value["avg_score"].(type) {
		case int:
			json.avgScore = v
		}

		switch v := value["low_score"].(type) {
		case int:
			json.lowScore = v
		}

		switch v := value["rank"].(type) {
		case int:
			json.rank = v
		}

		switch v := value["major"].(type) {
		case string:
			json.major = v
		}
		var enrollScore []int
		json.enrollScore = enrollScore
		jsonMap := append(jsonMap, json)
		fmt.Println(jsonMap)
		break
	}

	fmt.Printf("111")

	// fmt.Println(scoreArray)
	// fmt.Println(scoreCondition)
}
