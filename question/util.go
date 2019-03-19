// util
package question

import (
	"bufio"        // csv data load 관련
	"encoding/csv" // csv data load 관련
	"fmt"          // 출력 ( 디버그 )
	"log"
	"math"      // 반올림 관련
	"math/rand" // 임의 추출 관련

	"os"      // csv data load 관련
	"strconv" // string 관련 형변환
	"strings"
	"time" // 임의 추출 관련
)

func SaveResult_and_CurationDataAtDB() {
	rc_file, _ := os.Open("resources/data/CDI_AISpeaker_ResultAndCuration0317.csv") //result&curation file
	rc_reader := csv.NewReader(bufio.NewReader(rc_file))
	rows, _ := rc_reader.ReadAll()

	for i, row := range rows {
		for j := range row {
			log.Printf("%s", rows[i][j])
		}
		log.Println()
		break
	}

	for i := FIRST_IDX; i < len(rows); i++ {

		pattern := rows[i][0]

		//복합 질환인 경우 pattern 변수하나에 두 질환을 합쳐 저
		if rows[i][1] != "" {
			pattern += rows[i][1]
		}

		description := rows[i][2]
		explanation := []string{rows[i][3], rows[i][4], rows[i][5], rows[i][6]}
		var curation []string
		for j := 7; j < len(rows[i]); j++ {
			append(curation, rows[i][j])
		}
		temp := DB.ResultAndCuration{
			Pattern:     pattern,     // Pattern     []string `bson:"pattern"`
			Description: description, // Description string   `bson:"description"`
			Explanation: explanation, // Explanation []string `bson:"explanation"`
			Curation:    curation,    // Curation    []string `bson:"curation"`
		}
	}
}

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
	for i := FIRST_IDX; i < len(RAW_DATA.QCWP); {
		qNum, _ := strconv.Atoi(RAW_DATA.QCWP[i][WEIGHT]) // question number per one category
		randVal := rand.Intn(qNum) + i
		qdata.QRepIdx = append(qdata.QRepIdx, randVal)
		i += qNum
	}
	return qdata
}

// 3. initialize QDetailIdx of the structure QData
func qDetailIdxInit(qdata QData) QData {
	// make a restQIdx slice which includes all questions but the questions related to repIdx
	var restQIdx []int
	var qRepIdxIdx int = 0
	for i := FIRST_IDX; i < len(RAW_DATA.QCWP); i++ {
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
	var prevP string = RAW_DATA.QCWP[FIRST_IDX][PATTERN]
	var curP string
	var startPoint int = 0
	for i := 0; i < PATTERN_NUM; i++ {
		var qPIdx []int
		for j := startPoint; j < len(restQIdx); j++ {
			curP = RAW_DATA.QCWP[restQIdx[j]][PATTERN]
			if prevP != curP { // if pattern changed
				prevP = curP
				startPoint = j
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

// 5. calculate SQSProbPatternIdx or NoSQSProbPatternIdx and caculate FinalScore (only about the representative questions)
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
		weight, _ := strconv.Atoi(RAW_DATA.QCWP[qdata.QRepIdx[i]][WEIGHT])
		score[RAW_DATA.QCWP[qdata.QRepIdx[i]][PATTERN]] += biScore * weight                                                  // 기준치점수 * 가중치
		qdata.FinalScore[PATTERN_INDEX[RAW_DATA.QCWP[qdata.QRepIdx[i]][PATTERN]]] += float64(qdata.Answer[qdata.QRepIdx[i]]) // 대표질문에 대한 패턴별 총점
	}

	rand_seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(rand_seed) // qdata.SQSProbPatternIdx/NoSQSProbPatternIdx 를 랜덤으로 섞기 위함.

	for i := 0; i < len(PATTERN_NAME); i++ {
		if score[PATTERN_NAME[i]] > RAW_DATA.PtoC[PATTERN_NAME[i]] {
			qdata.SQSProbPatternIdx = append(qdata.SQSProbPatternIdx, i) // initialize SQSProbPatternIdx
		}
	}

	if len(qdata.SQSProbPatternIdx) != 0 {
		qdata.SQSProb = true
		for i := 0; i < len(qdata.SQSProbPatternIdx); i++ { // qdata.SQSProbPatternIdx의 순서를 섞음
			idxpicker := r.Intn(len(qdata.SQSProbPatternIdx)-i) + i
			// index range : i ~ qdata.SQSProbPatternIdx -1
			temp := qdata.SQSProbPatternIdx[i]
			qdata.SQSProbPatternIdx[i] = qdata.SQSProbPatternIdx[idxpicker]
			qdata.SQSProbPatternIdx[idxpicker] = temp
		}
	} else {
		qdata.SQSProb = false
		for i := 0; i < PATTERN_NUM; i++ { // 간단문진 결과 문제가 되는 패턴이 없을 때, NoSQSProbPatternIdx를 준비한다.
			qdata.NoSQSProbPatternIdx = append(qdata.NoSQSProbPatternIdx, i) // initialize NoSQSProbPatternIdx
		}
		for i := 0; i < len(qdata.NoSQSProbPatternIdx); i++ { // Shuffle NoSQSProbPatternIdx
			idxpicker := r.Intn(len(qdata.NoSQSProbPatternIdx)-i) + i
			// index range : i ~ qdata.NoSQSProbPatternIdx -1
			temp := qdata.NoSQSProbPatternIdx[i]
			qdata.NoSQSProbPatternIdx[i] = qdata.NoSQSProbPatternIdx[idxpicker]
			qdata.NoSQSProbPatternIdx[idxpicker] = temp
		}
	}

	return qdata
}

// 6. shuffle QDetIdx of the structure QData
func qDetailIdxShuffle(qdata QData) QData {

	if qdata.SQSProb == true { //  간단문진에서 문제가 있는 패턴이 정해져 있는 경우.
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
	} else { //  간단문진에서 문제가 있는 패턴이 없지만 정밀검진을 진행한다고 하는 경우.
		patlength := len(qdata.NoSQSProbPatternIdx)

		rand_seed := rand.NewSource(time.Now().UnixNano())
		r := rand.New(rand_seed) // To pick SliceIdx Randomly

		for i := 0; i < patlength; i++ { // 패턴 슬라이스의 길이만큼 이터레이트
			qdetlength := len(qdata.QDetailIdx[qdata.NoSQSProbPatternIdx[i]]) // 패턴 인덱스에 해당하는 질문 슬라이스의 길이
			for j := 0; j < qdetlength; j++ {                                 // Using knuth Shuffle
				idxpicker := r.Intn(qdetlength-j) + j
				// index range : j ~ qdetlength-1
				temp := qdata.QDetailIdx[qdata.NoSQSProbPatternIdx[i]][j]
				qdata.QDetailIdx[qdata.NoSQSProbPatternIdx[i]][j] = qdata.QDetailIdx[qdata.NoSQSProbPatternIdx[i]][idxpicker]
				qdata.QDetailIdx[qdata.NoSQSProbPatternIdx[i]][idxpicker] = temp
				// Swap
			}
		}
	}
	return qdata
}

// 7. calculate and make complete FinalScore
func calculateFinalScore(qdata QData) QData {

	if qdata.SQSProb == true {
		// make FinalScore when SQSProbPatternIdx exists
		for i := 0; i < len(qdata.SQSProbPatternIdx); i++ {
			for j := 0; j < len(qdata.QDetailIdx[qdata.SQSProbPatternIdx[i]]); j++ {
				qdata.FinalScore[qdata.SQSProbPatternIdx[i]] += float64(qdata.Answer[qdata.QDetailIdx[qdata.SQSProbPatternIdx[i]][j]]) // 문제가 있는 패턴에 대한 총점을 구하기 위해 점수를 더해나감
			}
		}

		// alter FinalScore as standard score
		for i := 0; i < len(qdata.SQSProbPatternIdx); i++ {
			maxScore := (len(qdata.QDetailIdx[qdata.SQSProbPatternIdx[i]]) + CATEGORY_NUM[qdata.SQSProbPatternIdx[i]]) * SCORE_MAX                    // 한 변증의 만점은 해당 변증에 대한 QDetailIdx 에 있는 질문 개수에 QRepIdx 에 있는 질문 개수 (1개)를 더한 값에 5를 곱한 값임
			qdata.FinalScore[qdata.SQSProbPatternIdx[i]] = math.Round((qdata.FinalScore[qdata.SQSProbPatternIdx[i]]*100/float64(maxScore))*100) / 100 // 표준점수로 변환, 소수점 2자리 반올림 (백분위 점수)
		}
	} else {
		// make FinalScore when NoSQSProbPatternIdx exists
		for i := 0; i < len(qdata.NoSQSProbPatternIdx); i++ {
			for j := 0; j < len(qdata.QDetailIdx[qdata.NoSQSProbPatternIdx[i]]); j++ {
				qdata.FinalScore[qdata.NoSQSProbPatternIdx[i]] += float64(qdata.Answer[qdata.QDetailIdx[qdata.NoSQSProbPatternIdx[i]][j]]) // 문제가 있는 패턴에 대한 총점을 구하기 위해 점수를 더해나감
			}
		}

		// alter FinalScore as standard score
		for i := 0; i < len(qdata.NoSQSProbPatternIdx); i++ {
			maxScore := (len(qdata.QDetailIdx[qdata.NoSQSProbPatternIdx[i]]) + CATEGORY_NUM[qdata.NoSQSProbPatternIdx[i]]) * SCORE_MAX                    // 한 변증의 만점은 해당 변증에 대한 QDetailIdx 에 있는 질문 개수에 QRepIdx 에 있는 질문 개수 (1개)를 더한 값에 5를 곱한 값임
			qdata.FinalScore[qdata.NoSQSProbPatternIdx[i]] = math.Round((qdata.FinalScore[qdata.NoSQSProbPatternIdx[i]]*100/float64(maxScore))*100) / 100 // 표준점수로 변환, 소수점 2자리 반올림 (백분위 점수)
		}
	}

	return qdata
}

// 8. Sort SQSProbPatternIdx / NoSQSProbPatternIdx
func sortProbPatternIdx(qdata QData) QData {

	if qdata.SQSProb == true { // Sort SQSProbPatternIdx using SelectionSort ( 오름차순 )

		if len(qdata.SQSProbPatternIdx) == 1 { // Sort할 필요가 없는 경우
			return qdata
		}

		for i := 0; i < len(qdata.SQSProbPatternIdx); i++ {
			var minIdx int = i
			for j := i; j < len(qdata.SQSProbPatternIdx); j++ {
				if qdata.SQSProbPatternIdx[minIdx] > qdata.SQSProbPatternIdx[j] {
					minIdx = j
				}
			}
			temp := qdata.SQSProbPatternIdx[i]
			qdata.SQSProbPatternIdx[i] = qdata.SQSProbPatternIdx[minIdx]
			qdata.SQSProbPatternIdx[minIdx] = temp
		}
	} else { // Sort NoSQSProbPatternIdx using SelectionSort ( 오름차순 )
		for i := 0; i < len(qdata.NoSQSProbPatternIdx); i++ {
			var minIdx int = i
			for j := i; j < len(qdata.NoSQSProbPatternIdx); j++ {
				if qdata.NoSQSProbPatternIdx[minIdx] > qdata.NoSQSProbPatternIdx[j] {
					minIdx = j
				}
			}
			temp := qdata.NoSQSProbPatternIdx[i]
			qdata.NoSQSProbPatternIdx[i] = qdata.NoSQSProbPatternIdx[minIdx]
			qdata.NoSQSProbPatternIdx[minIdx] = temp
		}
	}

	for i := 0; i < len(qdata.SQSProbPatternIdx); i++ {
		fmt.Println(qdata.SQSProbPatternIdx[i])
	}

	return qdata
}

// DATA PREPARE 1 (Representative Questions): execute 2 ~ 4.
func PrepareRep(qdata QData) QData {
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
	qdata.QDetailCount = 0

	return qdata
}

// DATA PREPARE 3 (Final Score): execite 7 ~ 8
func PrepareFin(qdata QData) QData {
	qdata = calculateFinalScore(qdata) // 7.
	qdata = sortProbPatternIdx(qdata)  // 8.
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

/*
// debug: print the contents of struct QData
func PrintStruct(qdata QData) {
	fmt.Println(RAW_DATA.QCWP)
	fmt.Println(RAW_DATA.PtoC)
	fmt.Println(qdata.QRepIdx)
	fmt.Println(qdata.QDetailIdx)
	fmt.Println(qdata.Answer)
}
*/
