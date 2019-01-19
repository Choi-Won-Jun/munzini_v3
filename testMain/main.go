// main.go
package main

import (
	//"fmt"
	"testbench/question"
)

func main() {
	var qdata question.QData

	p := []int{2, 3, 4} // pattern idx : index range - 0~4(Number of Patterns)

	qdata.RawData = question.LoadData()
	qdata = question.QIdxInit(qdata)
	qdata = question.QRepIdxInit(qdata)
	qdata = question.QDetailIdxInit(qdata)
	qdata = question.QRepIdxShuffle(qdata)
	qdata = question.QDetailIdxShuffle(qdata, p)
	question.PrintStruct(qdata)
}
