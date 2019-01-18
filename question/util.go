// util
package question

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"math/rand" // 임의 추출 관련
	"os"
	"strconv" // string 관련 형변환
	"time"    // 임의 추출 관련
)

// 1. load data to initalize the structure qDataConst
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

// 2. initialize QIdx of the structure QData
func QIdxInit(qdata QData) QData {
	for i := 1; i < len(qdata.RawData.QCWP); i++ {
		qdata.QIdx = append(qdata.QIdx, i)
	}
	return qdata
}

// 3. initialize QRepIdx of the structure QData
func QRepIdxInit(qdata QData) QData {
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < len(qdata.QIdx); {
		catNum, _ := strconv.Atoi(qdata.RawData.QCWP[qdata.QIdx[i]][WEIGHT])
		randVal := rand.Intn(catNum) + i
		qdata.QRepIdx = append(qdata.QRepIdx, randVal)
		i += catNum
	}
	return qdata
}

func printArray(arr [][]string) {
	for i := 1; i < len(arr); i++ {
		for j := 0; j < len(arr[i]); j++ {
			fmt.Printf("%s ", arr[i][j])
		}
		fmt.Println()
	}
}

func PrintStruct(qdata QData) {
	fmt.Println(qdata.RawData)
	fmt.Println(qdata.QIdx)
	fmt.Println(qdata.QRepIdx)
	fmt.Println(qdata.QDetailIdx)
	fmt.Println(qdata.Answer)
}
