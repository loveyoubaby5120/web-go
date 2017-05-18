package main

import (
	"fmt"
	"web-go/util/format"
	"web-go/util/jsonutil"
)

var (
	scoreArray     []int
	jsonArray      []format.Condition
	scoreCondition []map[string]interface{}
	err            error
)

func main() {
	_, _, _ = GetData()
	// for i := 0; i < 10; i++ {
	// 	fmt.Println(rand.Intn(100))
	// }

	for _, value := range jsonArray {
		fmt.Println(value)
	}

	// WriteScore(jsonArray)
}

// WriteScore write hight school score
func WriteScore(s interface{}) {
	var outJSONArray []byte
	outJSONArray, err = format.Map2JSON(s)
	if err != nil {
		fmt.Println(err)
	}
	jsonutil.WriteJSON("static/constrain_json/文一批.json", outJSONArray)
}

// GetData get hight school config
func GetData() ([]int, []format.Condition, error) {
	var err error
	scoreArray, err := jsonutil.ReadFileArray("static/constrain_json/constrain_json/文一批_一分一段.json")
	if err != nil {
		fmt.Println("readFile: ", err.Error())
		return nil, nil, err
	}
	fmt.Println("获取全部分数成功: ", len(scoreArray))

	scoreCondition, err = jsonutil.ReadFileJSONArray("static/constrain_json/constrain_json/文一批_高考分数.json")
	if err != nil {
		fmt.Println("readFile: ", err.Error())
		return nil, nil, err
	}

	for index, value := range scoreCondition {
		if index > 2 {
			break
		}

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

		jsonC := &format.Condition{
			EnrollNum:   enrollNum,
			AvgScore:    int(value["avg_score"].(float64)),
			LowScore:    int(value["low_score"].(float64)),
			Rank:        int(value["rank"].(float64)),
			Major:       value["major"].(string),
			EnrollScore: enrollScore,
		}

		jsonArray = append(jsonArray, *jsonC)
	}

	fmt.Println("获取配置成功: ", len(jsonArray))

	return scoreArray, jsonArray, nil
}

// Sum count total score
func Sum(enrollScore []int) int {
	total := 0
	for _, value := range enrollScore {
		total += value
	}
	return total
}
