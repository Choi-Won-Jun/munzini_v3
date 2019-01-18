// loader
package question

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv" // string 관련 형변환
)

func LoadData() qDataConst {
	qcwp_file, _ := os.Open("resources/data/QCWP.csv")   // open QCWP.csv
	ptoc_file, _ := os.Open("resources/data/cutoff.csv") // open cutoff.csv

	qcwp_reader := csv.NewReader(bufio.NewReader(qcwp_file)) // create csv reader for QCWP.csv
	ptoc_reader := csv.NewReader(bufio.NewReader(ptoc_file)) // create csv reader for cutoff.csv

	qcwp, _ := qcwp_reader.ReadAll() // read QCWP.csv
	ptoc, _ := ptoc_reader.ReadAll() // read cutoff.csv

	// test code to show how qcwp data looks like
	/*
		printArray(qcwp)
		printArray(ptoc)
	*/

	// make map from slice ptoc
	var ptocMap map[string]int
	ptocMap = make(map[string]int)
	for i := 1; i < len(ptoc); i++ {
		ptocMap[ptoc[i][PTOC_PATTERN]], _ = strconv.Atoi(ptoc[i][PTOC_CUTOFF])
	}

	// test code to show how ptocMap looks like
	// fmt.Print(ptocMap)

	qdatacon := qDataConst{
		QCWP: qcwp,
		PtoC: ptocMap,
	}

	return qdatacon
}

func printArray(arr [][]string) {
	for i := 1; i < len(arr); i++ {
		for j := 0; j < len(arr[i]); j++ {
			fmt.Printf("%s ", arr[i][j])
		}
		fmt.Println()
	}
}

func QIdxInit(qdata QData) QData {
	for i := 1; i < len(qdata.RawData.QCWP); i++ {
		qdata.QIdx = append(qdata.QIdx, i)
	}
	return qdata
}

func PrintStruct(qdata QData) {
	fmt.Print(qdata)
}
