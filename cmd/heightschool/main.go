package main

import (
	"encoding/json"
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
	isSuccess   bool
}

var (
	jsonArray      []Condition
	scoreCondition []map[string]interface{}
	scoreArray     []int
)

func main() {
	var err error
	scoreArray, err := jsonutil.ReadFileArray("static/predict_num_v2/constrain_json/文一批_一分一段.json")
	if err != nil {
		fmt.Println("readFile: ", err.Error())
		return
	}
	fmt.Println("获取全部分数成功: ", len(scoreArray))

	scoreCondition, err = jsonutil.ReadFileJsonArray("static/predict_num_v2/constrain_json/文一批_高考分数.json")
	if err != nil {
		fmt.Println("readFile: ", err.Error())
		return
	}

	for _, value := range scoreCondition {
		// jsonC, err := map2json2map.Map2json(value)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// fmt.Println(string(jsonC))
		// fmt.Println(value)
		enrollNum := int(value["enroll_num"].(float64))
		var enrollScore []int
		for i := 0; i < enrollNum; i++ {
			enrollScore = append(enrollScore, 0)
		}

		jsonC := &Condition{
			enrollNum:   enrollNum,
			avgScore:    int(value["avg_score"].(float64)),
			lowScore:    int(value["low_score"].(float64)),
			rank:        int(value["rank"].(float64)),
			major:       value["major"].(string),
			enrollScore: enrollScore,
		}

		jsonArray = append(jsonArray, *jsonC)

		break
	}

	fmt.Println("获取配置成功: ", len(jsonArray))
	outJsonArray, err := json.Marshal(jsonArray)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(jsonArray)
	jsonutil.WriteJson("static/predict_num_v2/文一批.json", outJsonArray)

}

func Sum(enrollScore []int) int {
	total := 0
	for _, value := range enrollScore {
		total += value
	}
	return total
}
