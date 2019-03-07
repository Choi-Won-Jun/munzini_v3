package nlp

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"        // csv load 관련
	"math/rand" // 임의 추출 관련
	"munzinis_project/question"
	"os"
	"strconv" // 문자열 함수 관련
	"time"    // 임의 추출 관련
)

func ConvertInquiryScore(str string) string {
	for i := 0; i < len(exStrArr); i++ {
		if strIn(str, exStrArr[i]) {
			return strconv.Itoa(i + 1)
		}
	}
	return str
}

func strIn(str string, arr []string) bool {
	for i := 0; i < len(arr); i++ {
		if str == arr[i] {
			return true
		}
	}
	return false
}

// load Data to initialize PlayUptoMessage
func loadData() PlayUptoConst {

	playUptoLow_file, _ := os.Open("resources/data/PlayUptoLow.csv")   // open PlayUptoLow.csv
	playUptoMid_file, _ := os.Open("resources/data/PlayUptoMid.csv")   // open PlayUptoMid.csv
	playUptoHigh_file, _ := os.Open("resources/data/PlayUptoHigh.csv") // open PlayUptoHigh.csv

	playUptoLow_reader := csv.NewReader(bufio.NewReader(playUptoLow_file))   // create csv reader for PlayUptoLow.csv
	playUptoMid_reader := csv.NewReader(bufio.NewReader(playUptoMid_file))   // create csv reader for PlayUptoMid.csv
	playUptoHigh_reader := csv.NewReader(bufio.NewReader(playUptoHigh_file)) // create csv reader for PlayUptoHigh.csv

	var playUptolow [][]string  // make 2-dimensional array to save PlayUptoLow.csv
	var playUptomid [][]string  // make 2-dimensional array to save PlayUptoMid.csv
	var playUptohigh [][]string // make 2-dimensional array to save PlayUptoHigh.csv

	// load playUptoLow.csv
	playUptolow = make([][]string, 0) // make slice
	var row int = 0                   // row value
	for {
		playUptoLow, err := playUptoLow_reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		playUptolow = append(playUptolow, playUptoLow)
		row++

	}

	// load PlayUptoMid.csv
	playUptomid = make([][]string, 0) // make slice
	row = 0                           // reset row to zero

	for {
		playUptoMid, err := playUptoMid_reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		playUptomid = append(playUptomid, playUptoMid)
		row++
	}

	// load PlayUptoHigh.csv
	playUptohigh = make([][]string, 0) // make slice
	row = 0
	for {
		playUptoHigh, err := playUptoHigh_reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		playUptohigh = append(playUptohigh, playUptoHigh)
		row++
	}

	playUptoconst := PlayUptoConst{
		PlayUptoLowPoint:  playUptolow,
		PlayUptoMidPoint:  playUptomid,
		PlayUptoHighPoint: playUptohigh,
	}

	return playUptoconst
}

// Choose Random Answer according to current_score
func GetPlayUptoMessage(current_score int, current_idx int) string {

	var playUptoMessage string
	//fmt.Println(current_idx)	// debug

	if current_score < question.BI_CRITERIA && current_score > 0 { // 1~2점 응답
		playUptoMessage = GetPlayUptoLow(current_idx)
	} else if current_score == question.BI_CRITERIA {
		playUptoMessage = GetPlayUptoMid(current_idx)
	} else if current_score > question.BI_CRITERIA && current_score <= question.SCORE_MAX {
		playUptoMessage = GetPlayUptoHigh(current_idx)
	} else {
		playUptoMessage = " "
	}
	return playUptoMessage
}

func GetPlayUptoLow(current_idx int) string {
	rand_seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(rand_seed) // Randomly Choose Answer
	random_idx := r.Intn(len(PlayUptoMessage.PlayUptoLowPoint[current_idx])-FIRST_IDX_C) + FIRST_IDX_C
	// index range : FIRST_IDX_C ~ len(PlayUptoMessage.PlayUptoLowPoint[current_idx]) - FIRST_IDX_C

	return PlayUptoMessage.PlayUptoLowPoint[current_idx][random_idx]
}

func GetPlayUptoMid(current_idx int) string {
	rand_seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(rand_seed) // Randomly Choose Answer
	random_idx := r.Intn(len(PlayUptoMessage.PlayUptoMidPoint[current_idx])-FIRST_IDX_C) + FIRST_IDX_C
	// index range : FIRST_IDX_C ~ len(PlayUptoMessage.PlayUptoMidPoint[current_idx]) - FIRST_IDX_C

	return PlayUptoMessage.PlayUptoMidPoint[current_idx][random_idx]
}

func GetPlayUptoHigh(current_idx int) string {
	rand_seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(rand_seed) // Randomly Choose Answer

	random_idx := r.Intn(len(PlayUptoMessage.PlayUptoHighPoint[current_idx])-FIRST_IDX_C) + FIRST_IDX_C
	// index range : FIRST_IDX_C ~ len(PlayUptoMessage.PlayUptoHighPoint[current_idx]) - FIRST_IDX_C
	return PlayUptoMessage.PlayUptoHighPoint[current_idx][random_idx]
}

// debug: print 2-dimensional slice
func printArray(arr [][]string) {
	for i := 0; i < len(arr); i++ {
		for j := 0; j < len(arr[i]); j++ {
			fmt.Printf("%s ", arr[i][j])
		}
		fmt.Println()
	}
}

/*
func PlayUpto(score int) string { // 사용자가 1~5점 사이로 대답을 하였을 때, 맞장구를 쳐주는 말
	scoreIdx := score - 1

	rand_seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(rand_seed) // Randomly Choose Answer
	answerPicker := r.Intn(len(playUptoPoint[scoreIdx]))

	playUptoMessage := playUptoPoint[scoreIdx][answerPicker] + "다음 질문입니다. " // Choose Answer

	return playUptoMessage
}
*/
