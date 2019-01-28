// util
package question

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"math"      // 반올림 관련
	"math/rand" // 임의 추출 관련
	"os"
	"strconv" // string 관련 형변환
	"time"    // 임의 추출 관련
)

// 1. load data to initalize the structure qDataConst
func loadData() qDataConst {
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

// 2. initialize QRepIdx of the structure QData
func qRepIdxInit(qdata QData) QData {
	rand.Seed(time.Now().UTC().UnixNano())
	for i := FIRST_IDX; i < len(qdata.RawData.QCWP); {
		catNum, _ := strconv.Atoi(qdata.RawData.QCWP[i][WEIGHT])
		randVal := rand.Intn(catNum) + i
		qdata.QRepIdx = append(qdata.QRepIdx, randVal)
		i += catNum
	}
	return qdata
}

// 3. initialize QDetailIdx of the structure QData
func qDetailIdxInit(qdata QData) QData {
	// make a restQIdx slice which includes all questions but the questions related to repIdx
	var restQIdx []int
	var qRepIdxIdx int = 0
	for i := FIRST_IDX; i < len(qdata.RawData.QCWP); i++ {
		if i == qdata.QRepIdx[qRepIdxIdx] { // if a value to put into restQIdx is same with the value of qRepIdx
			//fmt.Println(qdata.QRepIdx[qRepIdxIdx])
			if qRepIdxIdx+1 != len(qdata.QRepIdx) { // if qRepIdx[qRepIdxIdx] is not the very last value
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

// 4. shuffle QRepIdx of the structure QData
func qRepIdxShuffle(qdata QData) QData { //qdata의 qRepIdx를 섞는다.
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

// 5. calculate SQSProbPatternIdx, caculate FinalScore (only about the representative questions)
func calculateSQS(qdata QData) QData {
	qreplength := len(qdata.QRepIdx)                // 대표질문 개수
	var score map[string]int = make(map[string]int) // 변증 : (기준치점수 * 가중치)
	var biScore int                                 // binary Score = { 0 ,1 }

	for i := 0; i < len(PATTERN_NAME); i++ { // Initialize score map
		score[PATTERN_NAME[i]] = 0
	}

	for i := 0; i < len(PATTERN_NAME); i++ { // initialize FinalScore
		qdata.FinalScore = append(qdata.FinalScore, 0)
	}

	for i := 0; i < qreplength; i++ {
		// make binary score
		if qdata.Answer[qdata.QRepIdx[i]] > BI_CRITERIA {
			biScore = 1
		} else {
			biScore = 0
		}
		weight, _ := strconv.Atoi(qdata.RawData.QCWP[qdata.QRepIdx[i]][WEIGHT])
		score[qdata.RawData.QCWP[qdata.QRepIdx[i]][PATTERN]] += biScore * weight                                                  // 기준치점수 * 가중치
		qdata.FinalScore[PATTERN_INDEX[qdata.RawData.QCWP[qdata.QRepIdx[i]][PATTERN]]] += float64(qdata.Answer[qdata.QRepIdx[i]]) // 대표질문에 대한 패턴별 총점
	}

	for i := 0; i < len(PATTERN_NAME); i++ {
		if score[qdata.RawData.QCWP[qdata.QRepIdx[i]][PATTERN]] > qdata.RawData.PtoC[PATTERN_NAME[i]] {
			qdata.SQSProbPatternIdx = append(qdata.SQSProbPatternIdx, i)
		}
	}

	return qdata
}

// 6. shuffle QDetIdx of the structure QData
func qDetailIdxShuffle(qdata QData) QData {

	patlength := len(qdata.SQSProbPatternIdx)

	rand_seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(rand_seed) // To pick SliceIdx Randomly

	for i := 0; i < patlength; i++ { // 패턴 슬라이스의 길이만큼 이터레이트
		qdetlength := len(qdata.QDetailIdx[qdata.SQSProbPatternIdx[i]]) // 패턴 인덱스에 해당하는 질문 슬라이스의 길이
		for j := 0; j < qdetlength; j++ {                               // Using knuth Shuffle
			idxpicker := r.Intn(qdetlength-j) + j
			// index range : j ~ qdetlength-1
			temp := qdata.QDetailIdx[qdata.SQSProbPatternIdx[i]][j]
			qdata.QDetailIdx[qdata.SQSProbPatternIdx[i]][j] = qdata.QDetailIdx[qdata.SQSProbPatternIdx[i]][idxpicker]
			qdata.QDetailIdx[qdata.SQSProbPatternIdx[i]][idxpicker] = temp
			// Swap
		}
	}
	return qdata
}

// 7. calculate and make complete FinalScore
func calculateFinalScore(qdata QData) QData {
	// make FinalScore
	for i := 0; i < len(qdata.SQSProbPatternIdx); i++ {
		for j := 0; j < len(qdata.QDetailIdx[i]); j++ {
			qdata.FinalScore[qdata.SQSProbPatternIdx[i]] += float64(qdata.QDetailIdx[i][j]) // 문제가 있는 패턴에 대한 총점을 구하기 위해 점수를 더해나감
		}
	}

	// alter FinalScore as standard score
	for i := 0; i < len(qdata.SQSProbPatternIdx); i++ {
		maxScore := (len(qdata.QDetailIdx[i]) + 1) * 5                                                                                            // 한 변증의 만점은 해당 변증에 대한 QDetailIdx 에 있는 질문 개수에 QRepIdx 에 있는 질문 개수 (1개)를 더한 값에 5를 곱한 값임
		qdata.FinalScore[qdata.SQSProbPatternIdx[i]] = math.Round((qdata.FinalScore[qdata.SQSProbPatternIdx[i]]*100/float64(maxScore))*100) / 100 // 표준점수로 변환, 소수점 2자리 반올림 (백분위 점수)
	}

	return qdata
}

// DATA PREPARE 1 (Representative Questions): execute 1 ~ 4.
func PrepareRep(qdata QData) QData {
	qdata.RawData = loadData()    // 1.
	qdata = qRepIdxInit(qdata)    // 2.
	qdata = qDetailIdxInit(qdata) // 3.
	qdata = qRepIdxShuffle(qdata) // 4.
	qdata.RepIdx = 0
	qdata.DetPat = 0
	qdata.DetIdx = 0
	qdata.Answer = make(map[int]int)
	qdata.Answer[-1] = -1

	return qdata
}

// DATA PREPARE 2 (Detail Questions): execute 5 ~ 6.
func PrepareDet(qdata QData) QData {
	qdata = calculateSQS(qdata)      // 5.
	qdata = qDetailIdxShuffle(qdata) // 6.

	return qdata
}

// DATA PREPARE 3 (Final Score): execite 7.
func PrepareFin(qdata QData) QData {
	qdata = calculateFinalScore(qdata) // 7.

	return qdata
}

// debug: print 2-dimensional slice
func printArray(arr [][]string) {
	for i := FIRST_IDX; i < len(arr); i++ {
		for j := 0; j < len(arr[i]); j++ {
			fmt.Printf("%s ", arr[i][j])
		}
		fmt.Println()
	}
}

// debug: print the contents of struct QData
func PrintStruct(qdata QData) {
	fmt.Println(qdata.RawData.QCWP)
	fmt.Println(qdata.RawData.PtoC)
	fmt.Println(qdata.QRepIdx)
	fmt.Println(qdata.QDetailIdx)
	fmt.Println(qdata.Answer)
}
