package main

import (
	"fmt"
	"web-go/util/jsonutil"
	"web-go/util/map2json2map"
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

	// jsonArray := []Condition{}

	for _, value := range scoreCondition {
		fmt.Println(map2json2map.Map2json(value))
		// fmt.Println(value)

		// var jsonC Condition
		// jsonC.enrollNum = int(value["enroll_num"].(float64))
		// jsonC.avgScore = int(value["avg_score"].(float64))
		// jsonC.lowScore = int(value["low_score"].(float64))
		// jsonC.rank = int(value["rank"].(float64))
		// jsonC.major = value["major"].(string)

		// var enrollScore []int
		// for i := 0; i < jsonC.enrollNum; i++ {
		// 	enrollScore = append(enrollScore, 0)
		// }

		// jsonC.enrollScore = enrollScore
		// fmt.Println(jsonC)
		// jsonArray := append(jsonArray, jsonC)
		// fmt.Println(len(jsonArray[0].enrollScore))
		break
	}
}

func Sum(enrollScore []int) int {
	total := 0
	for _, value := range enrollScore {
		total += value
	}
	return total
}
