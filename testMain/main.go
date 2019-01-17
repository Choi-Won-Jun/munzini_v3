// main.go
package main

import (
	"testbench/question"
)

func main() {
	var data question.QData

	data.RawData = question.LoadData()
	data = question.QIdxInit(data)

	question.PrintStruct(data)
}
