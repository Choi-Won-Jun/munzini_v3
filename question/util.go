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
	for i := FIRST_IDX; i < len(ptoc); i++ {
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
	for i := FIRST_IDX; i < len(qdata.RawData.QCWP); i++ {
		qdata.QIdx = append(qdata.QIdx, i)
	}
	return qdata
}

// 3. initialize QRepIdx of the structure QData
func QRepIdxInit(qdata QData) QData {
	rand.Seed(time.Now().UTC().UnixNano())
	for i := FIRST_IDX; i < len(qdata.RawData.QCWP); {
		catNum, _ := strconv.Atoi(qdata.RawData.QCWP[i][WEIGHT])
		randVal := rand.Intn(catNum) + i
		qdata.QRepIdx = append(qdata.QRepIdx, randVal)
		i += catNum
	}
	return qdata
}

// 4. initialize QDetailIdx of the structure QData
func QDetailIdxInit(qdata QData) QData {
	// make a restQIdx slice which includes all questions but the questions related to repIdx
	var restQIdx []int
	var qRepIdxIdx int = 0
	for i := FIRST_IDX; i < len(qdata.RawData.QCWP); i++ {
		if i == qdata.QRepIdx[qRepIdxIdx] { // if a value to put into restQIdx is same with the value of qRepIdx
			//fmt.Println(qdata.QRepIdx[qRepIdxIdx])
			if qRepIdxIdx+1 != len(qdata.QRepIdx) { // if not qRepIdx[qRepIdxIdx] is the very last value
				qRepIdxIdx++
			}
			continue
		}
		restQIdx = append(restQIdx, i)
	}

	// make QDetailIdx
	var prevP string = qdata.RawData.QCWP[FIRST_IDX][PATTERN]
	var curP string
	var startPoint int = 0
	for i := 0; i < PATTERN_NUM; i++ {
		var qPIdx []int
		for j := startPoint; j < len(restQIdx); j++ {
			curP = qdata.RawData.QCWP[restQIdx[j]][PATTERN]
			if prevP != curP { // if pattern changed
				prevP = curP
				startPoint = j + 1
				break
			}
			prevP = curP
			qPIdx = append(qPIdx, restQIdx[j])
		}
		qdata.QDetailIdx = append(qdata.QDetailIdx, qPIdx)
	}
	return qdata
}

// 5. shuffle QRepIdx of the structure QData

func QRepIdxShuffle(qdata QData) QData { //qdata의 qRepIdx를 섞는다.
	qreplength := len(qdata.QRepIdx)

	rand_seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(rand_seed) // To pick SliceIdx Randomly

	for i := 0; i < qreplength; i++ { // Using knuth Shuffle
		idxpicker := r.Intn(qreplength-i) + i
		// index range : i ~ qreplength-1
		temp := qdata.QRepIdx[i]
		qdata.QRepIdx[i] = qdata.QRepIdx[idxpicker]
		qdata.QRepIdx[idxpicker] = temp
		// Swap
	}
	return qdata
}

// shuffle QDetIdx of the structure QData

func QDetailIdxShuffle(qdata QData, pattern []int) QData {

	patlength := len(pattern)

	rand_seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(rand_seed) // To pick SliceIdx Randomly

	for i := 0; i < patlength; i++ { // 패턴 슬라이스의 길이만큼 이터레이트
		qdetlength := len(qdata.QDetailIdx[pattern[i]]) // 패턴 인덱스에 해당하는 질문 슬라이스의 길이
		for j := 0; j < qdetlength; j++ {               // Using knuth Shuffle
			idxpicker := r.Intn(qdetlength-j) + j
			// index range : j ~ qdetlength-1
			temp := qdata.QDetailIdx[pattern[i]][j]
			qdata.QDetailIdx[pattern[i]][j] = qdata.QDetailIdx[pattern[i]][idxpicker]
			qdata.QDetailIdx[pattern[i]][idxpicker] = temp
			// Swap
		}
	}
	return qdata
}

func printArray(arr [][]string) {
	for i := FIRST_IDX; i < len(arr); i++ {
		for j := 0; j < len(arr[i]); j++ {
			fmt.Printf("%s ", arr[i][j])
		}
		fmt.Println()
	}
}

func PrintStruct(qdata QData) {
	fmt.Println(qdata.RawData.QCWP)
	fmt.Println(qdata.RawData.PtoC)
	fmt.Println(qdata.QIdx)
	fmt.Println(qdata.QRepIdx)
	fmt.Println(qdata.QDetailIdx)
	fmt.Println(qdata.Answer)
}

/* Old QRepIdxShuffle Algorithm
func QRepIdxShuffle(qdata QData) QData {

	Copy_QRep := make([]int, len(qdata.QRepIdx))

	qreplength := len(qdata.QRepIdx)
	isFilled := make([]bool, qreplength) // 새로 만들 qRep이 채워졌는지 ,모두 false 값

	rand_seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(rand_seed)

	n := 0
	for n < qreplength {
		if new_qrepIdx := r.Intn(qreplength); !isFilled[new_qrepIdx] {
			Copy_QRep[n] = qdata.QRepIdx[new_qrepIdx]
			isFilled[new_qrepIdx] = true
			n++
		}
	}
	copy(qdata.QRepIdx, Copy_QRep)
	return qdata
}
*/
/* Old QDetailIdxShuffle Algorithm
func QDetailIdxShuffle(qdata QData, pattern []int) QData { // int pattern 배열을 받는다.
	// 대표질문의 컷오프값을 넘긴 변증의 index값들을 담고있다.
	Copy_QDet := make([][]int, len(qdata.QDetailIdx))
	copy(Copy_QDet, qdata.QDetailIdx)                  // QDetailIdx 전체를 복사해온다.
	rand_seed := rand.NewSource(time.Now().UnixNano()) // Random Shuffle을 위함
	r := rand.New(rand_seed)                           // Random Shuffle을 위함

	for i := 0; i < len(pattern); i++ { // 해당하는 패턴들에 대해 Shuffle 반복 진행
		qdetlength := len(Copy_QDet[pattern[i]])
		isFilled := make([]bool, qdetlength)
		n := 0
		for n < qdetlength { // 해당하는 패턴의 길이만큼 Shuffle 진
			if new_qdetIdx := r.Intn(qdetlength); !isFilled[new_qdetIdx] {
				Copy_QDet[pattern[i]][n] = qdata.QDetailIdx[pattern[i]][new_qdetIdx]
				isFilled[new_qdetIdx] = true
				n++
			}
		}
	}
	copy(qdata.QDetailIdx, Copy_QDet)
	return qdata
}
*/
