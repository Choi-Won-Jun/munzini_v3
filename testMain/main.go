// main.go
package main

import (
	"fmt"
	"munzini/nlp"
	"munzini/question"
)

const FIRST_IDX = 1

func main() {

	var playUpto nlp.PlayUptoConst

	playUpto = nlp.LoadData()

	for i := FIRST_IDX; i < len(playUpto.PlayUptoLowPoint); i++ {
		for j := 0; j < len(playUpto.PlayUptoLowPoint[i]); j++ {
			fmt.Printf("%s ", playUpto.PlayUptoLowPoint[i][j])
		}
		fmt.Print("hi")
		fmt.Println()
	}
	question.LoadData()
	//var qdata question.QData

	//	p := []int{2, 3, 4} // pattern idx : index range - 0~4(Number of Patterns)

	//	qdata.RawData = question.LoadData()
	/*
		qdata = question.QRepIdxInit(qdata)
		qdata = question.QDetailIdxInit(qdata)
		qdata = question.QRepIdxShuffle(qdata)
		qdata = question.QDetailIdxShuffle(qdata, p)
		question.PrintStruct(qdata)
	*/
}
