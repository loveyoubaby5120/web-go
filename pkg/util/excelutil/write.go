package excelutil

import (
	"bufio"
	"bytes"

	"github.com/tealeg/xlsx"
)

func WriteStringToXlsx(sheet [][]string) ([]byte, error) {
	file := xlsx.NewFile()
	s, err := file.AddSheet("Sheet1")
	if err != nil {
		return nil, err
	}
	for _, row := range sheet {
		r := s.AddRow()
		for _, cell := range row {
			c := r.AddCell()
			c.Value = cell
		}
	}
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	if err := file.Write(w); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
