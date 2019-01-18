// main.go
package main

import (
	//"fmt"
	"testbench/question"
)

func main() {
	var qdata question.QData

	qdata.RawData = question.LoadData()
	qdata = question.QIdxInit(qdata)
	qdata = question.QRepIdxInit(qdata)

	question.PrintStruct(qdata)

}
