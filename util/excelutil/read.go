package excelutil

import (
	"fmt"

	"bytes"

	"github.com/extrame/xls"
	"github.com/tealeg/xlsx"
)

const InvalidString = "A2_INVALID_STRING"

func ReadToString(b []byte, format string) ([][][]string, error) {
	var data [][][]string
	switch format {
	case "xlsx":
		f, err := xlsx.OpenBinary(b)
		if err != nil {
			return nil, err
		}
		for _, sheet := range f.Sheets {
			var sheetData [][]string
			for _, row := range sheet.Rows {
				var rowData []string
				for _, cell := range row.Cells {
					if v, err := cell.String(); err == nil {
						rowData = append(rowData, v)
					} else {
						rowData = append(rowData, InvalidString)
					}
				}
				sheetData = append(sheetData, rowData)
			}
			data = append(data, sheetData)
		}
	case "xls":
		r := bytes.NewReader(b)
		wb, err := xls.OpenReader(r, "utf-8")
		if err != nil {
			return data, err
		}
		// Read multiple sheets into the first element in data
		data = append(data, wb.ReadAllCells(20000))
	default:
		return nil, fmt.Errorf("%s is not a valid excel format", format)
	}
	return data, nil
}

func ReadToJson(file string) {
	excelFileName := file
	// excelFileName := "static/predict_num_v2/constrain/文一批.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	json := []map[string]string{}
	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			for _, cell := range row.Cells {
				text, _ := cell.String()
				json = append(json, map[string]string{
					"cell": "cell",
				})
				fmt.Printf("%s\n", text)
			}
		}
	}
	fmt.Println(json)
}
